package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Container manages application bindings and dependency resolution.
type Container struct {
	bindings  map[string]*binding
	instances map[string]interface{}
	mu        sync.RWMutex
	providers []ServiceProvider
	booted    bool
}

// binding represents a service binding configuration.
type binding struct {
	abstract  string
	concrete  interface{}
	singleton bool
}

// ServiceProvider defines the interface for service providers.
type ServiceProvider interface {
	Register(c *Container) error
	Boot(c *Container) error
}

// New creates a new container instance.
func New() *Container {
	return &Container{
		bindings:  make(map[string]*binding),
		instances: make(map[string]interface{}),
	}
}

// Bind registers a transient binding in the container.
func (c *Container) Bind(abstract string, concrete interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings[abstract] = &binding{
		abstract:  abstract,
		concrete:  concrete,
		singleton: false,
	}
}

// Singleton registers a singleton binding in the container.
func (c *Container) Singleton(abstract string, concrete interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings[abstract] = &binding{
		abstract:  abstract,
		concrete:  concrete,
		singleton: true,
	}
}

// Instance directly registers an instance in the container.
func (c *Container) Instance(abstract string, instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.instances[abstract] = instance
}

// Make resolves a binding from the container.
func (c *Container) Make(abstract string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if instance exists (singleton)
	if instance, ok := c.instances[abstract]; ok {
		return instance, nil
	}

	// Get binding
	b, ok := c.bindings[abstract]
	if !ok {
		return nil, fmt.Errorf("binding not found: %s", abstract)
	}

	return c.resolve(b)
}

// resolve resolves a binding by calling its concrete function or constructor.
func (c *Container) resolve(b *binding) (interface{}, error) {
	var instance interface{}
	var err error

	switch concrete := b.concrete.(type) {
	case func() interface{}:
		instance = concrete()
	case func(*Container) interface{}:
		instance = concrete(c)
	case func() (interface{}, error):
		instance, err = concrete()
	case func(*Container) (interface{}, error):
		instance, err = concrete(c)
	default:
		// Try to call as a function with reflection
		fn := reflect.ValueOf(b.concrete)
		if fn.Kind() != reflect.Func {
			return b.concrete, nil
		}

		// Call function with no arguments
		result := fn.Call(nil)
		if len(result) == 0 {
			return nil, fmt.Errorf("invalid concrete for %s", b.abstract)
		}

		// Check for error return
		if len(result) > 1 && !result[1].IsNil() && result[1].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return result[0].Interface(), result[1].Interface().(error)
		}

		instance = result[0].Interface()
	}

	if err != nil {
		return nil, err
	}

	// Cache singleton instance
	if b.singleton {
		c.mu.Lock()
		c.instances[b.abstract] = instance
		c.mu.Unlock()
	}

	return instance, nil
}

// RegisterProvider registers a service provider.
func (c *Container) RegisterProvider(provider ServiceProvider) error {
	c.providers = append(c.providers, provider)
	if err := provider.Register(c); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}
	return nil
}

// BootProviders boots all registered providers.
func (c *Container) BootProviders() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.booted {
		return nil
	}

	for _, provider := range c.providers {
		if err := provider.Boot(c); err != nil {
			return fmt.Errorf("failed to boot provider: %w", err)
		}
	}

	c.booted = true
	return nil
}

// Call invokes a function with automatic dependency injection.
func (c *Container) Call(fn interface{}, params ...interface{}) (interface{}, error) {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return nil, errors.New("expected a function")
	}

	fnType := fnValue.Type()
	numIn := fnType.NumIn()

	// Build arguments
	args := make([]reflect.Value, 0, numIn)

	for i := 0; i < numIn; i++ {
		paramType := fnType.In(i)

		// Check if parameter is in params slice
		if i < len(params) {
			args = append(args, reflect.ValueOf(params[i]))
			continue
		}

		// Try to resolve from container
		var resolved reflect.Value

		// Check for *Container type
		if paramType == reflect.TypeOf((*Container)(nil)) {
			resolved = reflect.ValueOf(c)
		} else {
			// Try to resolve by type name
			resolved = reflect.Zero(paramType)
		}

		args = append(args, resolved)
	}

	// Call function
	results := fnValue.Call(args)

	if len(results) == 0 {
		return nil, nil
	}

	// Handle error return
	if len(results) > 1 && !results[1].IsNil() && results[1].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return results[0].Interface(), results[1].Interface().(error)
	}

	return results[0].Interface(), nil
}

// Get resolves a binding (alias for Make).
func (c *Container) Get(abstract string) (interface{}, error) {
	return c.Make(abstract)
}

// Has checks if a binding exists.
func (c *Container) Has(abstract string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, bindings := c.bindings[abstract]
	_, instances := c.instances[abstract]

	return bindings || instances
}

// Flush removes all bindings from the container.
func (c *Container) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings = make(map[string]*binding)
	c.instances = make(map[string]interface{})
	c.booted = false
}
