package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hasnat/govel/framework/validation"
)

// ContextHandler is a handler function that receives a Context.
type ContextHandler func(ctx *Context) error

// NewHandler wraps a ContextHandler to implement http.Handler interface.
func NewHandler(h ContextHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := New(w, r)
		if err := h(ctx); err != nil {
			_ = ctx.ServerError("Internal Server Error")
		}
	})
}

// Context represents an HTTP request/response context.
type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	container  interface{} // Service container reference
	user       interface{} // Authenticated user
	params     map[string]interface{}
	middleware []Middleware
	index      int
}

// Middleware is a request/response interceptor function.
type Middleware func(ctx *Context, next Handler) error

// Handler is an HTTP request handler.
type Handler func(ctx *Context) error

// New creates a new HTTP context.
func New(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:    r,
		Response:   w,
		params:     make(map[string]interface{}),
		middleware: make([]Middleware, 0),
		index:      -1,
	}
}

// JSON sends a JSON response.
func (c *Context) JSON(statusCode int, data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(statusCode)

	encoder := json.NewEncoder(c.Response)
	return encoder.Encode(data)
}

// String sends a plain text response.
func (c *Context) String(statusCode int, data string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(statusCode)
	_, err := c.Response.Write([]byte(data))
	return err
}

// HTML sends an HTML response.
func (c *Context) HTML(statusCode int, html string) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(statusCode)
	_, err := c.Response.Write([]byte(html))
	return err
}

// Redirect redirects to a given location.
func (c *Context) Redirect(location string, statusCode ...int) error {
	code := http.StatusFound
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	c.Response.Header().Set("Location", location)
	c.Response.WriteHeader(code)
	return nil
}

// Input retrieves a query or form parameter.
func (c *Context) Input(key string) string {
	// Check query parameters first
	if val := c.Request.URL.Query().Get(key); val != "" {
		return val
	}

	// Then form data
	if err := c.Request.ParseForm(); err == nil {
		return c.Request.FormValue(key)
	}

	return ""
}

// Param retrieves a URL parameter.
func (c *Context) Param(key string) string {
	if val, ok := c.params[key].(string); ok {
		return val
	}
	return ""
}

// SetParam sets a URL parameter.
func (c *Context) SetParam(key string, value interface{}) {
	c.params[key] = value
}

// GetParam retrieves a URL parameter as interface.
func (c *Context) GetParam(key string) interface{} {
	return c.params[key]
}

// QueryParams returns all query parameters.
func (c *Context) QueryParams() map[string][]string {
	return c.Request.URL.Query()
}

// Body reads the request body.
func (c *Context) Body() ([]byte, error) {
	defer c.Request.Body.Close()
	return io.ReadAll(c.Request.Body)
}

// BindJSON decodes JSON request body.
func (c *Context) BindJSON(data interface{}) error {
	defer c.Request.Body.Close()
	return json.NewDecoder(c.Request.Body).Decode(data)
}

// Cookie retrieves a cookie value.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// SetCookie sets a response cookie.
func (c *Context) SetCookie(name, value string, maxAge int) {
	http.SetCookie(c.Response, &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: maxAge,
		Path:   "/",
	})
}

// Header retrieves a request header.
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

// Status sets the HTTP status code (should be called before writing body).
func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
}

// Validate validates request data against rules.
func (c *Context) Validate(data interface{}, rules map[string]string) (validation.Errors, error) {
	validator := validation.New(c.Request.Context(), data, rules)
	return validator.Validate()
}

// SetUser sets the authenticated user.
func (c *Context) SetUser(user interface{}) {
	c.user = user
}

// User retrieves the authenticated user.
func (c *Context) User() interface{} {
	return c.user
}

// IsAuthenticated checks if a user is authenticated.
func (c *Context) IsAuthenticated() bool {
	return c.user != nil
}

// SetContainer sets the service container reference.
func (c *Context) SetContainer(container interface{}) {
	c.container = container
}

// Container retrieves the service container.
func (c *Context) Container() interface{} {
	return c.container
}

type contextKey string

const httpContextKey contextKey = "http_context"

// WithContext returns a new context with the given context.Context.
func (c *Context) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, httpContextKey, c)
}

// Next calls the next middleware in the chain.
func (c *Context) Next(ctx *Context) error {
	c.index++
	if c.index >= len(c.middleware) {
		return nil
	}

	return c.middleware[c.index](c, c.Next)
}

// UseMiddleware adds middleware to the context chain.
func (c *Context) UseMiddleware(m ...Middleware) {
	c.middleware = append(c.middleware, m...)
}

// Message is a generic response message.
type Message struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Success returns a success response.
func (c *Context) Success(statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Message{
		Message: message,
		Status:  statusCode,
		Data:    data,
	})
}

// Error returns an error response.
func (c *Context) Error(statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, Message{
		Message: message,
		Status:  statusCode,
		Errors:  errors,
	})
}

// BadRequest returns a 400 response.
func (c *Context) BadRequest(message string) error {
	return c.Error(http.StatusBadRequest, message, nil)
}

// Unauthorized returns a 401 response.
func (c *Context) Unauthorized(message string) error {
	return c.Error(http.StatusUnauthorized, message, nil)
}

// Forbidden returns a 403 response.
func (c *Context) Forbidden(message string) error {
	return c.Error(http.StatusForbidden, message, nil)
}

// NotFound returns a 404 response.
func (c *Context) NotFound(message string) error {
	return c.Error(http.StatusNotFound, message, nil)
}

// ServerError returns a 500 response.
func (c *Context) ServerError(message string) error {
	return c.Error(http.StatusInternalServerError, message, nil)
}

// Abort aborts the request with a status code.
func (c *Context) Abort(statusCode int) error {
	c.Response.WriteHeader(statusCode)
	return fmt.Errorf("aborted with status %d", statusCode)
}
