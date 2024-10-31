# Use a lightweight base image
FROM golang:alpine

# Create and use a non-root user
RUN adduser -D -g '' appuser

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o decapcms-oauth2 .

# Change ownership of the application binary
RUN chown appuser:appuser decapcms-oauth2

# Expose the port the app runs on
EXPOSE 9000

# Switch to the non-root user
USER appuser

# Command to run the application
CMD ["./decapcms-oauth2"]