# Govel Framework - Quick Reference Guide

## Quick Start (30 seconds)

```bash
cd /Users/hasnat/_WorkSpace_/govel
go mod download
go run main.go
# Server running on http://localhost:8000
```

---

## Common Tasks

### 1. Start Server
```bash
go run main.go
make dev
go run artisan.go serve
```

### 2. Create Controller
```bash
go run artisan.go make:controller UserController
# Then add to app/http/controllers/user_controller.go
```

### 3. Register Route
```go
// In routes/api.go
router.Get("/users", UserController{}.Index)
```

### 4. Validate Input
```go
errors, _ := ctx.Validate(data, map[string]string{
    "email": "required|email",
    "password": "required|min:8",
})
```

### 5. Return Response
```go
ctx.JSON(200, users)                    // JSON
ctx.Success(201, "Created", data)       // Success message
ctx.Error(422, "Failed", errors)        // Error response
ctx.Unauthorized("Auth required")       // 401
ctx.NotFound("Not found")               // 404
```

---

## Routing Cheat Sheet

```go
// HTTP Methods
router.Get("/path", handler)
router.Post("/path", handler)
router.Put("/path", handler)
router.Patch("/path", handler)
router.Delete("/path", handler)

// Parameters
router.Get("/users/{id}", handler)
id := ctx.Param("id")

// Query Parameters
value := ctx.Input("search")

// Route Groups
router.Prefix("/api").Group(func(r *Router) {
    r.Get("/users", handler)
})

// Middleware
authMiddleware := func(ctx *Context, next Handler) error {
    // Check auth
    return next(ctx)
}

router.Middleware(authMiddleware).Group(func(r *Router) {
    r.Get("/protected", handler)
})

// Multiple Middleware
router.Middleware(cors, auth, logging).Group(func(r *Router) {
    r.Get("/data", handler)
})
```

---

## Context API

```go
// Get Parameters
ctx.Param("id")                         // URL param
ctx.Input("search")                     // Query/Form param

// Request Body
var data map[string]interface{}
ctx.BindJSON(&data)                     // Parse JSON

// Validation
errors, _ := ctx.Validate(data, rules)

// Responses
ctx.JSON(200, data)                     // JSON response
ctx.HTML(200, "<h1>Hello</h1>")        // HTML response
ctx.String(200, "Hello")                // Text response
ctx.Redirect("/path")                   // 302 redirect

// Status & Headers
ctx.SetHeader("X-Custom", "value")
ctx.SetCookie("name", "value", 3600)

// Headers & Cookies
header := ctx.Header("Authorization")
value, _ := ctx.Cookie("session_id")

// Authentication
ctx.SetUser(user)
if ctx.IsAuthenticated() {
    user := ctx.User()
}
```

---

## Validation Rules

```go
// Basic Rules
"required"              // Cannot be empty
"nullable"              // Can be null
"present"               // Field must exist

// Type Rules
"string"                // Must be string
"integer|int"           // Must be integer
"numeric"               // Must be numeric
"array"                 // Must be array

// Format Rules
"email"                 // Valid email
"url"                   // Valid URL
"regex:/pattern/"       // Regex match
"starts_with:value"     // Starts with
"ends_with:value"       // Ends with

// Length Rules
"min:8"                 // Minimum length
"max:255"               // Maximum length
"between:10,100"        // Between min-max

// Comparison Rules
"confirmed"             // Matches field_confirmation
"unique:table,column"   // Unique in database (Phase 2)
"exists:table,column"   // Exists in database (Phase 2)
```

---

## Configuration

```go
// Load environment
config.LoadEnv(".env")

// Get values
config.GetString("APP_NAME", "Govel")
config.GetInt("APP_PORT", 8000)
config.GetBool("APP_DEBUG", false)
config.GetFloat("APP_TIMEOUT", 30.0)
config.Get("CUSTOM_KEY", default)

// Set values
config.Set("custom_key", value)
```

---

## Middleware

```go
// Create middleware
func AuthMiddleware(ctx *Context, next Handler) error {
    if !isAuthenticated(ctx) {
        return ctx.Unauthorized("Unauthorized")
    }
    return next(ctx)
}

// Apply globally
router.Middleware(AuthMiddleware).Group(func(r *Router) {
    r.Get("/protected", handler)
})

// Built-in Middleware
http.CORSMiddleware("*")              // CORS
http.SecurityHeadersMiddleware()       // Security headers
http.LoggingMiddleware()               // Request logging
http.RecoveryMiddleware()              // Panic recovery
http.TrimMiddleware()                  // Trim trailing slashes
http.JSONMiddleware()                  // Set JSON content-type
```

---

## Logging

```go
logger := logging.New(logging.INFO)

// Log levels
logger.Debug("message")                 // DEBUG
logger.Info("message")                  // INFO
logger.Warning("message")               // WARNING
logger.Error("message")                 // ERROR
logger.Critical("message")              // CRITICAL

// With context
logger.Info("User logged in", map[string]interface{}{
    "user_id": 123,
    "email": "user@example.com",
})

// File output
logger.SetFile("storage/logs/app.log")
```

---

## Service Container

```go
// Register singleton (one instance)
container.Singleton("logger", func() interface{} {
    return logging.New(logging.INFO)
})

// Register transient (new each time)
container.Bind("request", func() interface{} {
    return http.New(w, r)
})

// Register instance
container.Instance("config", config)

// Resolve from container
logger, _ := container.Make("logger")

// Check if exists
if container.Has("logger") {
    // Use it
}
```

---

## CLI (Artisan)

```bash
# Server
go run artisan.go serve
go run artisan.go serve --port=3000

# Generators
go run artisan.go make:controller NameController
go run artisan.go make:model Name
go run artisan.go make:migration create_table_name
go run artisan.go make:middleware AuthMiddleware
go run artisan.go make:request StoreUserRequest

# Database
go run artisan.go migrate
go run artisan.go migrate:rollback
go run artisan.go db:seed

# Utilities
go run artisan.go --help
```

---

## Project Structure

```
govel/
├── app/                     # Your application code
│   ├── http/
│   │   ├── controllers/     # Route handlers
│   │   ├── middleware/      # Custom middleware
│   │   └── requests/        # Form validation
│   ├── models/              # Database models (Phase 2)
│   ├── services/            # Business logic
│   └── jobs/                # Async jobs (Phase 3)
│
├── routes/
│   └── api.go               # Route definitions
│
├── config/                  # Configuration files
├── database/                # Migrations & seeders (Phase 2)
├── framework/               # Core framework (DO NOT MODIFY)
├── storage/                 # Logs, cache, files
├── public/                  # Static files
│
├── main.go                  # HTTP entry point
├── artisan.go               # CLI entry point
└── .env                     # Environment variables
```

---

## Common Patterns

### REST API Endpoints
```go
// List
router.Get("/api/users", UserController{}.Index)

// Create
router.Post("/api/users", UserController{}.Store)

// Show
router.Get("/api/users/{id}", UserController{}.Show)

// Update
router.Put("/api/users/{id}", UserController{}.Update)

// Delete
router.Delete("/api/users/{id}", UserController{}.Delete)
```

### Controller Pattern
```go
type PostController struct{}

func (c PostController) Index(ctx *http.Context) error {
    // List posts
    return ctx.JSON(200, posts)
}

func (c PostController) Store(ctx *http.Context) error {
    // Create post
    var data map[string]interface{}
    ctx.BindJSON(&data)
    
    errors, _ := ctx.Validate(data, map[string]string{
        "title": "required",
        "body": "required",
    })
    
    if errors != nil {
        return ctx.Error(422, "Validation failed", errors)
    }
    
    return ctx.Success(201, "Post created", post)
}

func (c PostController) Show(ctx *http.Context) error {
    // Get post
    id := ctx.Param("id")
    return ctx.JSON(200, post)
}

func (c PostController) Update(ctx *http.Context) error {
    // Update post
    return ctx.Success(200, "Post updated", post)
}

func (c PostController) Delete(ctx *http.Context) error {
    // Delete post
    return ctx.Success(200, "Post deleted", nil)
}
```

### Error Handling
```go
// Validation errors
if errors != nil {
    return ctx.Error(422, "Validation failed", errors)
}

// Not found
if user == nil {
    return ctx.NotFound("User not found")
}

// Unauthorized
if !isAuthorized(user) {
    return ctx.Unauthorized("Unauthorized")
}

// Server error
if err != nil {
    return ctx.ServerError("Internal server error")
}

// Custom error
return ctx.Error(400, "Custom error", nil)
```

---

## Testing

```bash
# Run all tests
go test ./...

# Verbose output
go test -v ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test ./app/http/controllers
```

---

## Troubleshooting

### Port Already in Use
```bash
# Use different port
go run main.go --port=9000
# or
APP_PORT=9000 go run main.go
```

### Module Not Found
```bash
go mod tidy
go mod download
```

### Route Not Matching
- Check exact path and method
- Routes are matched in order
- Remember parameter syntax: `{id}`

### Middleware Not Running
- Check order: global → group → route
- Ensure `next(ctx)` is called
- Middleware must be registered before route

---

## Performance Tips

1. **Use caching** - For frequently accessed data
2. **Connection pooling** - Database connections (Phase 2)
3. **Batch operations** - For database (Phase 2)
4. **Eager loading** - Load relations upfront (Phase 2)
5. **Index columns** - Database indices (Phase 2)
6. **Response compression** - Add middleware
7. **Goroutine pooling** - For worker processes (Phase 3)

---

## Next Steps

1. **Read**
   - START: README.md
   - DEEP: ARCHITECTURE.md
   - ROADMAP: IMPLEMENTATION_GUIDE.md

2. **Build**
   - Create controllers
   - Define routes
   - Add middleware
   - Validate input

3. **Extend**
   - Contribute features
   - Build extensions
   - Create libraries

4. **Deploy**
   - Build binary
   - Docker container
   - Cloud platform

---

## Useful Links

- [Go Documentation](https://golang.org/doc/)
- [GORM Docs](https://gorm.io/) (Phase 2)
- [Cobra CLI](https://cobra.dev/)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [REST API Design](https://restfulapi.net/)

---

## Help & Support

- 📖 **Documentation**: See README.md, ARCHITECTURE.md
- 🐛 **Report Bugs**: GitHub Issues
- 💬 **Questions**: GitHub Discussions
- 🤝 **Contribute**: See CONTRIBUTING.md

---

**Happy coding with Govel!** 🚀

Quick Reference v1.0 | Last Updated: May 2024
