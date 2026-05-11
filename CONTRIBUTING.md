# Contributing to Govel

Thank you for your interest in contributing to Govel! This document provides guidelines and directions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourname/govel.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make changes and commit: `git commit -am 'Add feature'`
5. Push to your fork: `git push origin feature/your-feature`
6. Submit a Pull Request

## Code Standards

### Go Style Guide

We follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go).

**Key principles**:
- Use `gofmt` for formatting
- Run `go vet` before committing
- Keep functions small and focused
- Use clear variable names
- Add comments for exported functions and types

### Naming Conventions

- **Packages**: lowercase, single word (`routing`, `validation`)
- **Types**: PascalCase (`UserController`, `ValidationError`)
- **Functions**: camelCase (`getUserByID`, `validateEmail`)
- **Constants**: UPPER_CASE (`MAX_CONNECTIONS`, `DEFAULT_TIMEOUT`)
- **Interfaces**: Verb ending (`Reader`, `Handler`, `Middleware`)
- **Files**: snake_case.go (`user_controller.go`, `auth_middleware.go`)

### Documentation

Every exported function, type, and constant must have a comment:

```go
// UserController handles user-related HTTP requests.
type UserController struct{}

// Index returns a list of all users.
func (c UserController) Index(ctx *http.Context) error {
    // Implementation
}
```

## Testing

All code must include tests:

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out
```

**Test file naming**: `filename_test.go`

**Test structure**:
```go
func TestUserIndex(t *testing.T) {
    // Arrange
    controller := UserController{}
    ctx := /* mock context */
    
    // Act
    err := controller.Index(ctx)
    
    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
}
```

## Commit Messages

Use clear, descriptive commit messages:

```
feat: add user authentication middleware
fix: resolve race condition in request context
docs: update routing documentation
test: add validation layer tests
refactor: improve container resolution performance
```

**Commit message format**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

## Pull Request Process

1. **Title**: Clear, descriptive title
2. **Description**: What does this PR do? Why?
3. **Related Issues**: Link to any related issues (#123)
4. **Checklist**:
   - [ ] Tests added/updated
   - [ ] Documentation updated
   - [ ] Code formatted (`gofmt`)
   - [ ] No linting errors (`go vet`)
   - [ ] All tests passing

**Example PR description**:
```markdown
## Description
Adds authentication system with session guard and token guard support.

## Related Issues
Fixes #42
Related to #38

## Changes
- Implement Auth interface
- Add SessionGuard implementation
- Add TokenGuard implementation
- Add authentication middleware
- Add tests for all components

## Testing
All tests pass locally:
```bash
go test ./framework/auth -v
```

## Areas for Contribution

### High Priority (Phase 2-3)
- [ ] ORM/Database layer
- [ ] Authentication system
- [ ] Event system
- [ ] Queue system
- [ ] Task scheduling

### Medium Priority (Phase 4)
- [ ] WebSocket support
- [ ] Mail system
- [ ] File storage
- [ ] Broadcasting

### Low Priority (Phase 5+)
- [ ] GraphQL API
- [ ] gRPC services
- [ ] Serverless support
- [ ] Multi-tenancy

## Bug Reports

Found a bug? Please report it on GitHub Issues with:

1. **Clear title**: Describe the bug briefly
2. **Reproduction steps**: How to reproduce?
3. **Expected behavior**: What should happen?
4. **Actual behavior**: What actually happens?
5. **Environment**: Go version, OS, etc.
6. **Error logs**: Any stack traces or error messages

**Example bug report**:
```
Title: Context parameters not persisted in middleware

Description:
When setting parameters on the context in middleware, they're lost in the handler.

Steps to Reproduce:
1. Create middleware that calls ctx.SetParam("user_id", 123)
2. Create handler that calls ctx.Param("user_id")
3. Run request

Expected: Parameter should be "123"
Actual: Parameter is ""

Environment:
- Go 1.23
- macOS 14.1
```

## Feature Requests

Have an idea for Govel? Submit a feature request:

1. **Title**: Clear description of feature
2. **Use case**: Why is this needed?
3. **Proposed solution**: How should it work?
4. **Alternatives**: Other approaches considered?

## Development Setup

### Prerequisites
- Go 1.23 or later
- Git
- Make (optional)

### Quick Setup
```bash
# Clone repository
git clone https://github.com/hasnat/govel.git
cd govel

# Install dependencies
go mod download

# Run tests
go test ./...

# Start development server
make dev  # or: go run main.go
```

### Useful Commands
```bash
make help        # Show all available commands
make test        # Run tests
make fmt         # Format code
make vet         # Run go vet
make lint        # Run linter (requires golangci-lint)
make clean       # Clean build artifacts
```

## Questions?

- **Discussions**: Use GitHub Discussions
- **Issues**: Use GitHub Issues for bug reports
- **Email**: Create an issue and mention @maintainers

## Code of Conduct

Please be respectful and constructive. We aim to create an inclusive and welcoming community.

**Be respectful**: Treat all contributors with respect and kindness.
**Be constructive**: Provide helpful feedback and suggestions.
**Be inclusive**: Welcome contributors of all backgrounds and experience levels.

---

Thank you for contributing to Govel! 🙏
