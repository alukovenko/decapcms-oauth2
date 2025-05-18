package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var clientId string
var clientSecret string
var trustedOrigins []string
var serverHost string
var serverPort string

func authHandler(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=repo,user", clientId)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func init() {
	clientId = os.Getenv("OAUTH_CLIENT_ID")
	clientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
	serverHost = os.Getenv("SERVER_HOST")
	serverPort = os.Getenv("SERVER_PORT")

	rawOrigins := os.Getenv("TRUSTED_ORIGINS")
	singleOrigin := os.Getenv("TRUSTED_ORIGIN")

	switch {
	case rawOrigins != "":
		trustedOrigins = strings.Split(rawOrigins, ",")
	case singleOrigin != "":
		trustedOrigins = []string{singleOrigin}
	default:
		log.Fatalf("You must set either TRUSTED_ORIGIN or TRUSTED_ORIGINS\n")
	}

	for i := range trustedOrigins {
		trustedOrigins[i] = strings.TrimSpace(trustedOrigins[i])
	}

	if clientId == "" || clientSecret == "" || serverPort == "" || len(trustedOrigins) == 0 {
		log.Fatalf("OAUTH_CLIENT_ID, OAUTH_CLIENT_SECRET, SERVER_PORT, and at least one TRUSTED_ORIGIN(S) are required\n")
	}
}

func getAccessToken(code string) (string, error) {
	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}

	if token, ok := responseData["access_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("access token not found")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := getAccessToken(code)
	if err != nil {
		http.Error(w, "Error getting access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	postMsgContent := map[string]string{
		"token":    token,
		"provider": "github",
	}
	postMsgJSON, _ := json.Marshal(postMsgContent)
	trustedJSON, _ := json.Marshal(trustedOrigins)

	script := fmt.Sprintf(`
<html>
<body>
<script>
(function() {
	const trusted = %s;

	function receiveMessage(e) {
		if (trusted.includes(e.origin)) {
			window.opener.postMessage(
				'authorization:github:success:%s',
				e.origin
			);
		} else {
			console.warn("Origin not trusted:", e.origin);
		}
	}

	window.addEventListener("message", receiveMessage, false);
	window.opener.postMessage("authorizing:github", "*");
})();
</script>
</body>
</html>
`, trustedJSON, string(postMsgJSON))

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(script))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/callback", callbackHandler)

	bindTo := fmt.Sprintf("%s:%s", serverHost, serverPort)
	log.Printf("Server started on http://%s\n", bindTo)
	if err := http.ListenAndServe(bindTo, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
