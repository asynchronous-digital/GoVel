# Govel Framework Architecture

## Overview

Govel is a production-grade Laravel-inspired Go web framework that combines the expressiveness and developer experience of Laravel with the performance and concurrency capabilities of Go.

### Design Philosophy

- **Laravel DX First**: Familiar syntax and folder structure for Laravel developers
- **Go Best Practices**: Idiomatic Go internally with interfaces and composition
- **Performance**: Built for high concurrency with goroutines and channels
- **Modularity**: Loosely coupled components via dependency injection
- **Scalability**: Enterprise-ready architecture for production systems

---

## Core Architecture

### 1. Service Container (Dependency Injection)

**Package**: `framework/container`

The foundation of Govel. Manages application bindings and dependencies.

**Key Features**:
- Singleton & transient bindings
- Automatic dependency resolution via reflection
- Service provider registration
- Bootstrapping lifecycle

**Core Interfaces**:
```go
type Container interface {
    Bind(abstract string, concrete interface{})
    Singleton(abstract string, concrete interface{})
    Make(abstract string) (interface{}, error)
    Resolve(abstract string) (interface{}, error)
}
```

### 2. HTTP Context

**Package**: `framework/http`

Wraps `*http.Request` and `*http.ResponseWriter` with Laravel-like ergonomics.

**Key Responsibilities**:
- Request/Response manipulation
- Dependency injection for handlers
- Middleware chain execution
- JSON/XML/HTML responses
- Session/Cookie management
- Authentication context

**Methods**:
- `JSON(statusCode int, data interface{})`
- `String(statusCode int, data string)`
- `HTML(statusCode int, html string)`
- `Redirect(location string)`
- `Input(key string) string`
- `Validate(rules map[string]string) error`
- `Auth() User` (when authenticated)

### 3. Router

**Package**: `framework/routing`

Expressive, Laravel-like routing with groups, middleware, and named routes.

**Features**:
- Method routing (GET, POST, PUT, DELETE, PATCH)
- Route groups with prefixes & middleware
- Named routes
- Route parameters with constraints
- Resource routes
- API routes
- Rate limiting hooks

**Syntax**:
```go
router.Get("/users", UserController{}.Index)
router.Middleware("auth").Group(func(r Router) {
    r.Get("/dashboard", DashboardController{}.Show)
})
```

### 4. Middleware Pipeline

**Package**: `framework/http`

Request/response interceptors with elegant chaining.

**Execution Flow**:
```
Request → Global Middleware → Route Middleware → Handler → Route Middleware → Global Middleware → Response
```

**Middleware Interface**:
```go
type Middleware func(ctx *Context, next Handler) error
```

### 5. Configuration Management

**Package**: `framework/support` (config module)

Environment-aware configuration with type safety.

**Features**:
- `.env` file support
- Typed config access
- Environment switching
- Config caching
- Configuration providers

**Usage**:
```go
config.Get("app.name", "Govel")
config.GetInt("app.port", 8000)
config.GetBool("app.debug", false)
```

### 6. Validation

**Package**: `framework/validation`

Laravel-inspired validation with expressive rules.

**Rules**:
- required, nullable, present
- string, integer, numeric, array, object
- email, url, ip, phone
- regex, ends_with, starts_with
- min, max, between
- unique, exists
- custom rules

**Usage**:
```go
validator := validation.New(ctx)
validator.Rule("email", "required|email")
validator.Rule("password", "required|min:8|confirmed")
errors := validator.Validate()
```

### 7. ORM Abstraction

**Package**: `framework/database`

GORM-backed with Laravel-like query builder and relationships.

**Features**:
- Model definition with associations
- Query builder interface
- Eager loading
- Scopes
- Transactions
- Pagination
- Soft deletes
- Timestamps

**Model Interface**:
```go
type Model interface {
    TableName() string
    Relationships() map[string]Relationship
}
```

### 8. Authentication

**Package**: `framework/auth`

Multi-guard authentication system.

**Guards**:
- Session guard
- Token/Bearer guard
- Custom guard support

**Features**:
- User authentication
- Authorization policies
- Password hashing (bcrypt)
- Remember tokens
- Auth context in handlers

### 9. Events

**Package**: `framework/events`

Event broadcasting with listeners.

**Features**:
- Event registration
- Listener binding
- Event dispatching
- Queued events
- Event broadcasting

### 10. Queue System

**Package**: `framework/queue`

Background job processing with multiple drivers.

**Drivers**:
- Sync (in-process)
- Database
- Redis (future)
- SQS (future)

**Job Interface**:
```go
type Job interface {
    Handle(ctx context.Context) error
}
```

### 11. Caching

**Package**: `framework/cache`

Multi-store caching system.

**Drivers**:
- In-memory
- File-based
- Database
- Redis (future)

### 12. Scheduling

**Package**: `framework/scheduling`

Cron-like job scheduling.

**Features**:
- Scheduled commands
- Periodic tasks
- Timezone support
- Callback scheduling

### 13. Logging

**Package**: `framework/logging`

Structured logging with multiple channels.

**Features**:
- Multiple log levels
- Structured logging
- Log formatting
- File rotation
- Syslog integration (future)

### 14. CLI / Artisan

**Package**: `framework/console`

Command-line interface using Cobra.

**Built-in Commands**:
- `serve` - Start development server
- `make:controller` - Generate controller
- `make:model` - Generate model
- `make:middleware` - Generate middleware
- `make:request` - Generate form request
- `make:event` - Generate event
- `make:listener` - Generate listener
- `make:job` - Generate job
- `migrate` - Run migrations
- `migrate:rollback` - Rollback migrations
- `db:seed` - Seed database
- `queue:work` - Process queue jobs

### 15. Session Management

**Package**: `framework/session`

Server-side session handling.

**Drivers**:
- Cookie-based
- Database
- File-based
- In-memory

### 16. WebSocket Support

**Package**: `framework/websocket`

WebSocket server infrastructure.

**Features**:
- Connection management
- Message broadcasting
- Rooms/channels
- Authentication
- Middleware support

---

## Application Structure

### `app/` - Application Code

Where user code lives. Organized by concern:

- `console/` - Console commands
- `exceptions/` - Exception classes
- `http/` - HTTP layer
  - `controllers/` - Route handlers
  - `middleware/` - HTTP middleware
  - `requests/` - Form request validation
- `jobs/` - Queueable jobs
- `listeners/` - Event listeners
- `models/` - Database models
- `policies/` - Authorization policies
- `providers/` - Service providers
- `services/` - Business logic
- `traits/` - Shared behavior

### `bootstrap/` - Bootstrapping

Application initialization and setup.

### `config/` - Configuration

Environment-specific configuration files:
- `app.go`
- `database.go`
- `cache.go`
- `queue.go`
- `mail.go`
- `session.go`

### `database/`

Database evolution and seeding:
- `migrations/` - Schema migrations
- `seeders/` - Database seeders
- `factories/` - Model factories for testing

### `framework/` - Core Framework

Reusable, generic framework components. Independent of application code.

### `public/` - Web Root

Publicly served static files.

### `resources/`

Application assets:
- `views/` - HTML templates
- `lang/` - Localization files
- `assets/` - CSS, JS, images

### `routes/`

Route definitions:
- `web.go` - Web routes
- `api.go` - API routes
- `console.go` - Console commands

### `storage/`

Runtime storage:
- `logs/` - Application logs
- `cache/` - Cached files
- `app/` - Application files

### `tests/` - Test Suite

Unit, feature, and integration tests.

---

## Request Lifecycle

```
HTTP Request
    ↓
HTTP Server Routes to Framework
    ↓
Global Middleware (Middleware Pipeline)
    ↓
Service Container Resolves Route Handler
    ↓
Route Middleware
    ↓
Context Created (Request wrapped)
    ↓
Handler Executed (Controller method)
    ↓
Response Returned
    ↓
After Middleware
    ↓
HTTP Response Sent
    ↓
Cleanup & Shutdown Hooks
```

---

## Dependency Injection Flow

```
Application.Boot()
    ↓
Register Service Providers
    ↓
Container Registers Bindings
    ↓
Register Routes
    ↓
Start HTTP Server
    ↓
Request comes in
    ↓
Container resolves Controller
    ↓
All dependencies auto-injected
    ↓
Handler executed with resolved dependencies
```

---

## Design Patterns

### 1. Service Provider Pattern
Providers register bindings, boot services, and publish configuration.

```go
type ServiceProvider interface {
    Register(c Container) error
    Boot(c Container) error
}
```

### 2. Middleware Pattern
Composable request/response interceptors.

### 3. Repository Pattern
Database access abstraction via interfaces.

### 4. Factory Pattern
Model factories for testing.

### 5. Observer Pattern
Event listeners and dispatchers.

### 6. Singleton Pattern
Single instance services in container.

---

## Concurrency Model

- **Goroutine per Request**: Each HTTP request handled in its own goroutine
- **Context Propagation**: `context.Context` passed through entire request chain
- **Non-blocking I/O**: Database and external service calls are non-blocking
- **Channel Communication**: Inter-goroutine communication via channels
- **Graceful Shutdown**: Context cancellation on server shutdown

---

## Performance Considerations

### Memory Management
- Connection pooling for database
- Cache for frequently accessed data
- Request object reuse (object pools)

### Concurrency
- Goroutine per request (HTTP/2 multiplexing)
- Non-blocking I/O for all I/O operations
- Worker pools for queue processing

### Caching Strategies
- HTTP caching headers
- Response caching middleware
- Database query caching
- Full-page caching

---

## Error Handling

**Hierarchy**:
```
Exception (interface)
    ├── ValidationException
    ├── AuthenticationException
    ├── AuthorizationException
    ├── NotFoundHttpException
    ├── ServerErrorHttpException
    └── Custom exceptions
```

**Error Response**:
```json
{
  "message": "Validation failed",
  "status": 422,
  "errors": {
    "email": ["Email is required"]
  }
}
```

---

## Testing Strategy

**Test Layers**:
1. **Unit Tests**: Business logic isolation
2. **Feature Tests**: Full request/response cycle
3. **Integration Tests**: External service interaction
4. **Performance Tests**: Load and stress testing

**Testing Helpers**:
- HTTP test client
- Database seeding for tests
- Faker for realistic data
- Assertion helpers

---

## Security

### Built-in Protections
- CSRF tokens (middleware)
- XSS prevention (output encoding)
- SQL injection prevention (prepared statements)
- Password hashing (bcrypt)
- Rate limiting
- CORS support
- Security headers middleware

---

## Scalability

### Horizontal Scaling
- Stateless request handling
- Shared session/cache stores (Redis)
- Database connection pooling
- Load balancer friendly

### Vertical Scaling
- Efficient goroutine usage
- Memory-efficient data structures
- Lazy loading
- Caching layers

### Future Enhancements
- Microservices integration
- Service mesh support
- Distributed tracing
- Observability integration

---

## Configuration

### Environment Variables
Read from `.env` file with fallback to OS environment.

### Configuration Files
Go files in `config/` directory with typed accessors.

### Example: `config/database.go`
```go
func DatabaseConfig(app Container) map[string]interface{} {
    return map[string]interface{}{
        "driver": os.Getenv("DB_DRIVER"),
        "host": os.Getenv("DB_HOST"),
        "port": os.Getenv("DB_PORT"),
        "database": os.Getenv("DB_DATABASE"),
        "username": os.Getenv("DB_USERNAME"),
        "password": os.Getenv("DB_PASSWORD"),
    }
}
```

---

## Extension Points

### Custom Service Providers
Register custom services and bootstrap logic.

### Custom Middleware
Request/response interception.

### Custom Commands
CLI commands via `make:command`.

### Custom Validation Rules
Extend validation with custom rules.

### Custom Guards
Extend authentication.

### Custom Exception Handlers
Application-specific error handling.

---

## Module Organization

Each `framework/*` package is independently usable:
- Import only what you need
- Minimal external dependencies
- Clear interfaces
- Easy to test
- Easy to replace

---

## Naming Conventions

- **Packages**: lowercase, single word (`routing`, `validation`)
- **Types**: PascalCase (`UserController`, `CreateUserRequest`)
- **Functions**: camelCase (`getUserByID`, `validateEmail`)
- **Constants**: UPPER_CASE (`DB_HOST`)
- **Interfaces**: ending with "er" or "or" (`Reader`, `Logger`, `Handler`)
- **Files**: snake_case.go (`user_controller.go`, `auth_middleware.go`)

---

## Development Workflow

### 1. Project Setup
```bash
go run artisan.go migrate
go run artisan.go db:seed
```

### 2. Generate Resources
```bash
go run artisan.go make:controller UserController
go run artisan.go make:model User
go run artisan.go make:migration create_users_table
```

### 3. Development
```bash
go run artisan.go serve
```

### 4. Testing
```bash
go test ./...
```

### 5. Deployment
```bash
go build -o app
./app --port=8000
```

---

## Roadmap

### Phase 1 (Current)
- ✅ Project structure
- ✅ Service container
- ✅ Router
- ✅ Middleware pipeline
- ✅ Context
- ✅ Basic CLI

### Phase 2
- Validation layer
- ORM abstraction
- Authentication
- Database migrations

### Phase 3
- Events
- Queue system
- Caching
- Scheduling

### Phase 4
- WebSocket support
- Mail system
- File storage
- Broadcasting

### Phase 5
- Advanced features
- Performance optimization
- Production hardening
- Observability integration

---

## Dependencies (Strategic)

Minimal external dependencies:
- `gorm.io/gorm` - ORM (optional, for database)
- `spf13/cobra` - CLI
- `joho/godotenv` - .env parsing
- `stretchr/testify` - Testing (dev only)
- Go standard library: `net/http`, `encoding/json`, `context`, etc.

---

## Community & Contribution

Govel prioritizes:
- Clear, well-documented code
- Test coverage
- Backward compatibility
- Performance
- Developer experience
- Production readiness
