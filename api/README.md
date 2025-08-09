# Go API with JWT Authentication

This Go API demonstrates JWT token validation for Better-Auth integration. It provides a secure backend that can validate JWT tokens issued by the Better-Auth frontend.

## Features

- JWT token validation using JWKS (JSON Web Key Sets)
- CORS support for frontend integration
- Request logging middleware
- Health check endpoint
- Protected authentication verification endpoint

## Prerequisites

- Go 1.24.2 or higher
- Better-Auth frontend application running (for JWT token generation)

## Installation

1. Clone the repository and navigate to the API directory:
```bash
cd api
```

2. Install dependencies:
```bash
go mod download
```

## Running the API

Start the server:
```bash
go run main.go
```

The API will start on port 8080 by default.

## Available Endpoints

- `GET /health` - Health check endpoint
- `GET /api/auth/verify` - Verify JWT token (requires valid Authorization header)

## Environment Variables

The API uses the following environment variables (optional):
- `PORT` - Server port (default: 8080)

## Testing

Run tests:
```bash
go test ./...
```

## How It Works

1. The API fetches public keys from the Better-Auth JWKS endpoint
2. When a request comes in with a JWT token, it validates the token using the cached public keys
3. Valid tokens allow access to protected endpoints and return user information
4. Invalid or expired tokens are rejected with appropriate error messages