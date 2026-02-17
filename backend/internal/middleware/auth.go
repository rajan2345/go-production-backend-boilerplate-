package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rajan2345/go-boilerplate/internal/server"
)

// struct and initializer function should be there
type AuthMiddleware struct {
	server *server.Server
}

// Initializer function
func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

// using this server instance we can have access to everything like logger etc..
// that we are passing as a dependency (here is one of the example of dependency injection)

// now we just need one function to check if the request coming is authenticated or not
// it basically means it will have all the tokens e.g. jwt, oauth etc. and if it is not we will be passing 404 error
// we will be using clerk for all this
func (a *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {

}
