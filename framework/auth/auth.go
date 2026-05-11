package auth
// Authentication system (Phase 3)
// Provides multi-guard authentication with user management

package auth

import (
	"context"
)

// Guard defines an authentication method
type Guard interface {
	// User login/authentication
	Attempt(credentials map[string]interface{}) (User, error)
	
	// Check if authenticated
	Check(ctx context.Context) bool
	
	// Get authenticated user
	User(ctx context.Context) User
	
	// Login a specific user
	Login(user User) error
	
	// Logout user
	Logout(ctx context.Context) error
}

// User represents an authenticated user
type User interface {
	GetID() interface{}
	GetEmail() string
	GetName() string
}

// Authenticatable defines what a user model must implement
type Authenticatable interface {
	User
	GetPasswordHash() string
	SetPasswordHash(hash string)
}

// Provider manages authentication guards
type Provider interface {
	Guard(name string) Guard
	SetDefaultGuard(name string)
}

// Available guards:
// - session: Traditional session-based authentication
// - token: API token authentication
// - jwt: JSON Web Token authentication (future)

// Password management
type PasswordBroker interface {
	Reset(email string) error
	Verify(token string, email string, password string) error
}

// Future implementation with bcrypt and JWT support
