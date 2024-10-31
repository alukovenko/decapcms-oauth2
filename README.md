# OAuth2 GitHub Authentication for DecapCMS

This utility provides a simple way to authenticate users via GitHub OAuth for self-hosted https://decapcms.org/ admin panels. Typical configuration suggested in https://decapcms.org/docs/configure-decap-cms/ requires creating an account on [Netlify](https://www.netlify.com/) and adding your website to the configuration. Using this configuration implies sharing access to your repositories with a third party.

This project is a drop-in replacement for the OAuth2 used in DecapCMS. It includes an HTTP server that handles the OAuth flow, redirects users to GitHub for authentication, and retrieves an access token.

## Prerequisites

- Go 1.22.7 or later; or just use Docker
- GitHub OAuth application with OAUTH_CLIENT_ID and OAUTH_CLIENT_SECRET
- Environment variables set for OAUTH_CLIENT_ID, OAUTH_CLIENT_SECRET, SERVER_HOST, SERVER_PORT, and TRUSTED_ORIGIN

## Docker

```
docker run -it --rm -e OAUTH_CLIENT_ID=<CLIENT_ID> -e OAUTH_CLIENT_SECRET=<CLIENT_SECRET> -e TRUSTED_ORIGIN=https://www.example.com -e SERVER_HOST=127.0.0.1 -e SERVER_PORT=9000 alukovenko/decapcms-oauth2:latest
```

## Installation

1. Log into your [GitHub account](https://github.com/login) and [register](https://github.com/settings/applications/new) a new OAuth app. Refer to the [GitHub documentation](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) if in doubt. Write down Client ID and Client Secret for the next steps.

2. Clone the repository:

```bash
git clone https://github.com/yourusername/oauth-github-auth.git
cd oauth-github-auth
```

Build the application:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o decapcms-oauth2
```

## Configuration

1. Set the following environment variables:

- OAUTH_CLIENT_ID: Your GitHub OAuth application's client ID
- OAUTH_CLIENT_SECRET: Your GitHub OAuth application's client secret
- SERVER_HOST: The host address for the server to bind on (e.g., localhost)
- SERVER_PORT: The port for the server to bind on (e.g., 8080)
- TRUSTED_ORIGIN: The trusted origin for CORS (e.g., https://example.com if your DecapCMS is on https://example.com/admin/; for debug purposes this can be set to "\*")

Example:

```bash
export OAUTH_CLIENT_ID=your_client_id
export OAUTH_CLIENT_SECRET=your_client_secret
export SERVER_HOST=localhost
export SERVER_PORT=8080
export TRUSTED_ORIGIN=http://localhost:3000
```

2. Update your config.yml / config.yaml configuration of DecapCMS, adding your self-hosted backend as `base_url`:

```yaml
backend:
  name: github
  repo: username/repository
  branch: main
  base_url: https://auth.example.com
```

## Usage

```bash
./oauth-github-auth
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## Contact

For any questions or issues, please open an issue on GitHub.
