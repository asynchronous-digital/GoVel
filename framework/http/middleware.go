package http

import (
	"strings"
)

// MiddlewareStack chains middleware together.
type MiddlewareStack struct {
	middlewares []Middleware
}

// NewMiddlewareStack creates a new middleware stack.
func NewMiddlewareStack(middlewares ...Middleware) *MiddlewareStack {
	return &MiddlewareStack{
		middlewares: append([]Middleware(nil), middlewares...),
	}
}

// Add appends middleware to the stack.
func (ms *MiddlewareStack) Add(m ...Middleware) *MiddlewareStack {
	ms.middlewares = append(ms.middlewares, m...)
	return ms
}

// Handle chains all middleware and executes the handler.
func (ms *MiddlewareStack) Handle(ctx *Context, handler Handler) error {
	if len(ms.middlewares) == 0 {
		return handler(ctx)
	}

	// Build chain
	var chain Handler = handler
	for i := len(ms.middlewares) - 1; i >= 0; i-- {
		middleware := ms.middlewares[i]
		current := chain
		chain = func(c *Context) error {
			return middleware(c, func(c *Context) error {
				return current(c)
			})
		}
	}

	return chain(ctx)
}

// Built-in middleware

// CORSMiddleware handles Cross-Origin Resource Sharing.
func CORSMiddleware(allowedOrigins ...string) Middleware {
	return func(ctx *Context, next Handler) error {
		origin := ctx.Request.Header.Get("Origin")

		// Simple check: allow all origins or specific ones
		allowed := len(allowedOrigins) == 0
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed || origin == "" {
			ctx.SetHeader("Access-Control-Allow-Origin", "*")
			ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// Handle preflight
		if ctx.Request.Method == "OPTIONS" {
			ctx.Status(200)
			return nil
		}

		return next(ctx)
	}
}

// JSONMiddleware ensures JSON content type for API requests.
func JSONMiddleware(ctx *Context, next Handler) error {
	ctx.SetHeader("Content-Type", "application/json")
	return next(ctx)
}

// TrimMiddleware trims slashes from URLs.
func TrimMiddleware(ctx *Context, next Handler) error {
	path := ctx.Request.URL.Path
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		ctx.Request.URL.Path = strings.TrimSuffix(path, "/")
	}
	return next(ctx)
}

// LoggingMiddleware logs HTTP requests (placeholder).
func LoggingMiddleware(ctx *Context, next Handler) error {
	// In production, integrate with the logging system
	err := next(ctx)
	return err
}

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware(ctx *Context, next Handler) error {
	defer func() {
		if r := recover(); r != nil {
			ctx.ServerError("Internal Server Error")
		}
	}()
	return next(ctx)
}

// SecurityHeadersMiddleware adds security headers.
func SecurityHeadersMiddleware(ctx *Context, next Handler) error {
	ctx.SetHeader("X-Content-Type-Options", "nosniff")
	ctx.SetHeader("X-Frame-Options", "SAMEORIGIN")
	ctx.SetHeader("X-XSS-Protection", "1; mode=block")
	ctx.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")
	return next(ctx)
}
