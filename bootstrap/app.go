package bootstrap
package bootstrap

import (
	"fmt"

	"github.com/hasnat/govel/framework/console"
	"github.com/hasnat/govel/framework/container"
	"github.com/hasnat/govel/framework/logging"
	"github.com/hasnat/govel/framework/routing"
	"github.com/hasnat/govel/framework/support"
)

// Application represents the Govel application instance.
type Application struct {
	container *container.Container
	config    *support.Config
	logger    *logging.Logger
	router    *routing.Router
	kernel    *console.Kernel
	basePath  string
}

// New creates a new application instance.
func New(basePath string) *Application {
	app := &Application{
		container: container.New(),
		config:    support.NewConfig(),
		logger:    logging.New(logging.INFO),
		router:    routing.New(),
		kernel:    console.NewKernel(),
		basePath:  basePath,
	}

	// Register core services
	app.registerCoreServices()

	return app
}

// registerCoreServices registers core framework services.
func (app *Application) registerCoreServices() {
	// Register config as singleton
	app.container.Singleton("config", func() interface{} {
		return app.config
	})

	// Register logger as singleton
	app.container.Singleton("logger", func() interface{} {
		return app.logger
	})

	// Register router as singleton
	app.container.Singleton("router", func() interface{} {
		return app.router
	})

	// Register container as singleton
	app.container.Singleton("container", func() interface{} {
		return app.container
	})
}

// Boot bootstraps the application.
func (app *Application) Boot() error {
	// Load environment
	if err := app.config.LoadEnv(app.basePath + "/.env"); err != nil {
		app.logger.Warning("Could not load .env file: " + err.Error())
	}

	// Boot service providers
	if err := app.container.BootProviders(); err != nil {
		return fmt.Errorf("failed to boot providers: %w", err)
	}

	app.logger.Info("Application bootstrapped successfully")
	return nil
}

// Register registers a service provider.
func (app *Application) Register(provider container.ServiceProvider) error {
	return app.container.RegisterProvider(provider)
}

// Container returns the service container.
func (app *Application) Container() *container.Container {
	return app.container
}

// Config returns the configuration instance.
func (app *Application) Config() *support.Config {
	return app.config
}

// Logger returns the logger instance.
func (app *Application) Logger() *logging.Logger {
	return app.logger
}

// Router returns the router instance.
func (app *Application) Router() *routing.Router {
	return app.router
}

// Kernel returns the console kernel.
func (app *Application) Kernel() *console.Kernel {
	return app.kernel
}

// BasePath returns the application base path.
func (app *Application) BasePath() string {
	return app.basePath
}

// Make resolves a service from the container.
func (app *Application) Make(abstract string) (interface{}, error) {
	return app.container.Make(abstract)
}

// Singleton registers a singleton service.
func (app *Application) Singleton(abstract string, concrete interface{}) {
	app.container.Singleton(abstract, concrete)
}

// Bind registers a transient service.
func (app *Application) Bind(abstract string, concrete interface{}) {
	app.container.Bind(abstract, concrete)
}

// Run starts the application HTTP server.
func (app *Application) Run(port string) error {
	if port == "" {
		port = app.Config().GetString("APP_PORT", "8000")
	}

	addr := ":" + port
	app.Logger().Info("Starting server on http://localhost:" + port)

	return nil // Would return http.ListenAndServe(addr, app.Router())
}
