package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// one constant is the name of the header where we will store the requestID and find the requestID and other requestID
const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// One thing to note -- this is middleware -- and this is only declaration -- initialization will be happening in final router
// this function return a new middleware as you see in echo.MiddlewareFunc

// inside function we are defining a handler , everytime a new request comes a new context is being created through our web framework,
// we can access this context and using this context we can get different thing or info about our request e.g. the path parameters, query parameters , the payload and everything
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Set(RequestIDKey, requestID)
			c.Response().Header().Set(RequestIDHeader, requestID)

			return next(c)
		}
	}
}

func GetRequestID(c echo.Context) string {
	if requestID, ok := c.Get(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// echo.MiddlewareFunc represents a middleware component in echo web framework
// context -- temporary storage for the whole lifecycle of the request
