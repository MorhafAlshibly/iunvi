package middleware

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
)

type CORS struct {
	AllowedOrigins []string
	MaxAge         int
}

func WithAllowedOrigins(allowedOrigins []string) func(*CORS) {
	return func(input *CORS) {
		input.AllowedOrigins = allowedOrigins
	}
}

func WithMaxAge(maxAge int) func(*CORS) {
	return func(input *CORS) {
		input.MaxAge = maxAge
	}
}

func NewCORS(options ...func(*CORS)) *CORS {
	cors := &CORS{
		MaxAge: 7200,
	}
	for _, option := range options {
		option(cors)
	}
	return cors
}

func (c *CORS) Middleware(connectHandler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   c.AllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   append(connectcors.AllowedHeaders(), "Authorization"),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		MaxAge:           c.MaxAge,
	}).Handler(connectHandler)
}
