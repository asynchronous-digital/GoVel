package database
# ORM abstraction layer (Phase 2)
# This package provides a high-level interface to GORM for database operations
# while keeping the framework database-agnostic

package database

import (
	"context"
	"gorm.io/gorm"
)

// DB is the main database interface
type DB interface {
	// Connection management
	Ping(ctx context.Context) error
	Close() error

	// Query operations
	Query() QueryBuilder
	First(dest interface{}, conds ...interface{}) error
	All(dest interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Count(model interface{}) int64

	// Write operations
	Create(value interface{}) error
	Update(model interface{}, updates map[string]interface{}) error
	Delete(model interface{}, conds ...interface{}) error

	// Transactions
	BeginTx(ctx context.Context) DB
	Commit() error
	Rollback() error

	// Relationships
	LoadRelation(model interface{}, relation string) error
}

// QueryBuilder provides a fluent interface for building queries
type QueryBuilder interface {
	Where(query interface{}, args ...interface{}) QueryBuilder
	OrderBy(column string, direction string) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Select(columns ...string) QueryBuilder
	Distinct() QueryBuilder
	Get(dest interface{}) error
	First(dest interface{}) error
	Count() int64
	Exists() bool
}

// Model is the base interface for database models
type Model interface {
	TableName() string
	PrimaryKey() string
}

// Relationship defines model relationships
type Relationship interface {
	Type() string
	Related() interface{}
	ForeignKey() string
	LocalKey() string
}

// HasMany relationship
type HasMany struct {
	LocalModel  interface{}
	RelatedModel interface{}
	ForeignKey string
	LocalKey   string
}

// BelongsTo relationship
type BelongsTo struct {
	LocalModel   interface{}
	RelatedModel interface{}
	ForeignKey  string
	OwnerKey    string
}

// ManyToMany relationship
type ManyToMany struct {
	LocalModel   interface{}
	RelatedModel interface{}
	Table        string
	ForeignKey  string
	RelatedKey  string
}

// Scopes - reusable query conditions
type Scope func(db DB) DB

// Hook interface for lifecycle events
type Hook interface {
	BeforeCreate(db DB) error
	AfterCreate(db DB) error
	BeforeUpdate(db DB) error
	AfterUpdate(db DB) error
	BeforeDelete(db DB) error
	AfterDelete(db DB) error
}

// Future: Implementation will wrap GORM and provide Laravel-like syntax
