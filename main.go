package main
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hasnat/govel/bootstrap"
	"github.com/hasnat/govel/framework/console"
	"github.com/hasnat/govel/framework/http"
	"github.com/hasnat/govel/routes"
)

func main() {
	// Determine if running in console mode
	if len(os.Args) > 1 {
		runConsole()
		return
	}

	// Initialize application
	app := bootstrap.New(".")

	// Bootstrap application
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	// Register routes
	routes.RegisterWebRoutes(app.Router())
	routes.RegisterAPIRoutes(app.Router())

	// Add global middleware
	router := app.Router()
	router.NotFound(&http.ErrorHandler{})

	// Add middleware to router
	setupMiddleware(router)

	// Start server
	port := app.Config().GetString("APP_PORT", "8000")
	addr := ":" + port

	app.Logger().Info("Starting server on http://localhost:" + port)

	if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

// setupMiddleware configures middleware for the router.
func setupMiddleware(router *http.Router) {
	// In a real application, middleware would be applied here
	// via router.Middleware() or route group middleware
}

// runConsole runs the console/CLI.
func runConsole() {
	app := bootstrap.New(".")

	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	kernel := app.Kernel()

	// Register console commands
	registerCommands(kernel)

	// Execute console command
	if err := kernel.Execute(); err != nil {
		log.Fatalf("Console error: %v", err)
	}
}

// registerCommands registers console commands.
func registerCommands(kernel *console.Kernel) {
	// Register built-in commands
	kernel.Register(console.ServeCommand(func(port string) error {
		// Start server logic
		app := bootstrap.New(".")
		app.Boot()

		routes.RegisterWebRoutes(app.Router())
		routes.RegisterAPIRoutes(app.Router())

		addr := ":" + port
		return http.ListenAndServe(addr, app.Router())
	}))

	kernel.Register(console.MakeControllerCommand())
	kernel.Register(console.MakeModelCommand())
	kernel.Register(console.MakeMigrationCommand())
	kernel.Register(console.MakeMiddlewareCommand())
	kernel.Register(console.MigrateCommand())
	kernel.Register(console.SeedCommand())
}
