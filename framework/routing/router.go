package routing

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	httpfw "github.com/hasnat/govel/framework/http"
)

// Router manages HTTP routes and dispatching.
type Router struct {
	routes          []*Route
	middlewares     []httpfw.Middleware
	prefix          string
	compiled        map[string]*compiledRoute
	notFoundHandler http.Handler
}

// Route represents an HTTP route.
type Route struct {
	Method      string
	Path        string
	Handler     http.Handler
	RouteName   string
	Middlewares []httpfw.Middleware
	pattern     *regexp.Regexp
	paramNames  []string
}

// compiledRoute represents a compiled route for matching.
type compiledRoute struct {
	pattern    *regexp.Regexp
	paramNames []string
	route      *Route
}

// New creates a new router instance.
func New() *Router {
	return &Router{
		routes:   make([]*Route, 0),
		compiled: make(map[string]*compiledRoute),
	}
}

// Get registers a GET route.
func (r *Router) Get(path string, handler http.Handler) *Route {
	return r.Add(http.MethodGet, path, handler)
}

// Post registers a POST route.
func (r *Router) Post(path string, handler http.Handler) *Route {
	return r.Add(http.MethodPost, path, handler)
}

// Put registers a PUT route.
func (r *Router) Put(path string, handler http.Handler) *Route {
	return r.Add(http.MethodPut, path, handler)
}

// Patch registers a PATCH route.
func (r *Router) Patch(path string, handler http.Handler) *Route {
	return r.Add(http.MethodPatch, path, handler)
}

// Delete registers a DELETE route.
func (r *Router) Delete(path string, handler http.Handler) *Route {
	return r.Add(http.MethodDelete, path, handler)
}

// Any registers a route for all HTTP methods.
func (r *Router) Any(path string, handler http.Handler) *Route {
	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions,
	}

	route := &Route{
		Method:  strings.Join(methods, "|"),
		Path:    r.prefix + path,
		Handler: handler,
	}

	r.routes = append(r.routes, route)
	r.compileRoute(route)

	return route
}

// Add registers a route with the given method.
func (r *Router) Add(method, path string, handler http.Handler) *Route {
	fullPath := r.prefix + path

	route := &Route{
		Method:      method,
		Path:        fullPath,
		Handler:     handler,
		Middlewares: append([]httpfw.Middleware(nil), r.middlewares...),
	}

	r.routes = append(r.routes, route)
	r.compileRoute(route)

	return route
}

// Middleware adds middleware to subsequent routes.
func (r *Router) Middleware(middlewares ...httpfw.Middleware) *RouterGroup {
	return &RouterGroup{
		router:      r,
		middlewares: middlewares,
	}
}

// Prefix adds a prefix to subsequent routes.
func (r *Router) Prefix(prefix string) *RouterGroup {
	return &RouterGroup{
		router: r,
		prefix: prefix,
	}
}

// Group creates a route group with shared prefix/middleware.
func (r *Router) Group(cb func(r *Router)) {
	group := &RouterGroup{router: r}
	group.Group(cb)
}

// Name sets the name of the last registered route.
func (route *Route) Name(name string) *Route {
	route.RouteName = name
	return route
}

// Middleware adds middleware to a specific route.
func (route *Route) Middleware(middlewares ...httpfw.Middleware) *Route {
	route.Middlewares = append(route.Middlewares, middlewares...)
	return route
}

// compileRoute compiles a route pattern for matching.
func (r *Router) compileRoute(route *Route) {
	pattern := route.Path
	params := make([]string, 0)

	// Replace {param} with regex groups
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatchIndex(pattern, -1)

	for _, match := range matches {
		paramName := pattern[match[2]:match[3]]
		params = append(params, paramName)
		pattern = pattern[:match[0]] + "([^/]+)" + pattern[match[1]:]
	}

	// Escape special regex characters
	pattern = "^" + regexp.QuoteMeta(route.Path) + "$"
	pattern = strings.ReplaceAll(pattern, `\{[^}]+\}`, "([^/]+)")

	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid route pattern: %s", route.Path))
	}

	route.pattern = compiledPattern
	route.paramNames = params

	key := route.Method + ":" + route.Path
	r.compiled[key] = &compiledRoute{
		pattern:    compiledPattern,
		paramNames: params,
		route:      route,
	}
}

// Match finds a matching route for the given method and path.
func (r *Router) Match(method, path string) (*Route, map[string]string, error) {
	for _, route := range r.routes {
		// Check method
		if !r.methodMatches(route, method) {
			continue
		}

		// Check path pattern
		matches := route.pattern.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		// Extract parameters
		params := make(map[string]string)
		for i, name := range route.paramNames {
			if i+1 < len(matches) {
				params[name] = matches[i+1]
			}
		}

		return route, params, nil
	}

	return nil, nil, fmt.Errorf("route not found: %s %s", method, path)
}

// methodMatches checks if a route matches the given method.
func (r *Router) methodMatches(route *Route, method string) bool {
	if route.Method == "" {
		return true
	}

	methods := strings.Split(route.Method, "|")
	for _, m := range methods {
		if m == method {
			return true
		}
	}

	return false
}

// ServeHTTP implements http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := httpfw.New(w, req)

	// Match route
	route, params, err := r.Match(req.Method, req.URL.Path)
	if err != nil {
		if r.notFoundHandler != nil {
			r.notFoundHandler.ServeHTTP(w, req)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 Not Found")
		}
		return
	}

	// Set URL parameters on context
	for k, v := range params {
		ctx.SetParam(k, v)
	}

	// Call handler
	route.Handler.ServeHTTP(w, req)
}

// NotFound sets a custom 404 handler.
func (r *Router) NotFound(handler http.Handler) {
	r.notFoundHandler = handler
}

// RouterGroup manages a group of routes with shared configuration.
type RouterGroup struct {
	router      *Router
	prefix      string
	middlewares []httpfw.Middleware
}

// Group creates nested route groups.
func (rg *RouterGroup) Group(cb func(r *Router)) {
	// Save current state
	oldPrefix := rg.router.prefix
	oldMiddlewares := append([]httpfw.Middleware(nil), rg.router.middlewares...)

	// Apply group configuration
	rg.router.prefix = oldPrefix + rg.prefix
	rg.router.middlewares = append(oldMiddlewares, rg.middlewares...)

	// Execute callback
	cb(rg.router)

	// Restore state
	rg.router.prefix = oldPrefix
	rg.router.middlewares = oldMiddlewares
}

// Get registers a GET route in the group.
func (rg *RouterGroup) Get(path string, handler http.Handler) *Route {
	return rg.router.Get(path, handler)
}

// Post registers a POST route in the group.
func (rg *RouterGroup) Post(path string, handler http.Handler) *Route {
	return rg.router.Post(path, handler)
}

// Put registers a PUT route in the group.
func (rg *RouterGroup) Put(path string, handler http.Handler) *Route {
	return rg.router.Put(path, handler)
}

// Patch registers a PATCH route in the group.
func (rg *RouterGroup) Patch(path string, handler http.Handler) *Route {
	return rg.router.Patch(path, handler)
}

// Delete registers a DELETE route in the group.
func (rg *RouterGroup) Delete(path string, handler http.Handler) *Route {
	return rg.router.Delete(path, handler)
}

// Prefix creates a nested prefix group.
func (rg *RouterGroup) Prefix(prefix string) *RouterGroup {
	return &RouterGroup{
		router:      rg.router,
		prefix:      rg.prefix + prefix,
		middlewares: rg.middlewares,
	}
}

// Middleware creates a nested middleware group.
func (rg *RouterGroup) Middleware(middlewares ...httpfw.Middleware) *RouterGroup {
	newMiddlewares := append([]httpfw.Middleware(nil), rg.middlewares...)
	newMiddlewares = append(newMiddlewares, middlewares...)

	return &RouterGroup{
		router:      rg.router,
		prefix:      rg.prefix,
		middlewares: newMiddlewares,
	}
}
