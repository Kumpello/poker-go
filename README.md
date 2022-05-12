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

All requests (except `/auth/*`) are required to have JWT token attached (`Header -> Authorization: Bearer <<token>>`).
When token expires a user must renew it (with `/login`).

# Articles Scrapper

The app can scrape some sources to get news that can be displayed somewhere else. Currently, supported sources are:

* pokernews.com/news

The news are scrapped from the main page, then saved to the database. Clients get data from database. To run the
scrapper use `go run cmd/cli/main.go updateArticles`.

From Article data the hashCode is calculated and duplicated news are filtered. The hashCode will be returned to the
client and used for the paging (TBD).

# Development

To run app in development, at first run MongoDB docker container:
`docker run --name mongodb -p 27017:27017 -e MONGODB_ROOT_PASSWORD=password123 bitnami/mongodb:4.4`

This command will run the mongodb container with root user: `root:password123` on port 27017

## Binding Requests

There is a very useful generic function that binds the HTTP request and validates it.

```go
package request

func echoFunc(c echo.Context) error {
	data, bindErr := binder.BindRequest[bodyType, queryType](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	// to access the body/query use:
	data.Request // binds basing on json-tags and validate-tags
	data.Query   // TBD

	return c.String(200, "ok")
}
```

# Production

TBD.
