# Run...

Poker-GO is a tool allowing to manage poker-games.

# Authorization

- All users are required to create an account.
- For authentication a JWT token is required.
- To signup use `/auth/signup`

```go
package auth

type signUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type logInRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type authResponse struct {
	ID           string `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
```

All requests (except `/auth/*`) are required to have JWT token attached (`Header: token`).
When token expires a user must renew it (with `/login`).

# Development

To run app in development, at first run MongoDB docker container:
`docker run --name mongodb -p 27017:27017 -e MONGODB_ROOT_PASSWORD=password123 bitnami/mongodb:4.4`

This command will run the mongodb container with root user: `root:password123` on port 27017

# Production

TBD.
