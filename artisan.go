package main

import (
	"log"
	"os"

	"github.com/hasnat/govel/bootstrap"
	"github.com/hasnat/govel/framework/console"
)

// This file serves as the Artisan CLI entry point
// Run with: go run artisan.go <command> [options]

func init() {
	// Only run this if explicitly called as artisan
	if len(os.Args) == 0 || (len(os.Args) > 0 && os.Args[0] != "artisan.go") {
		return
	}
}

// main is the entry point for the Artisan CLI.
// Usage:
// go run artisan.go serve
// go run artisan.go make:controller UserController
// go run artisan.go migrate
func main() {
	// Initialize application
	app := bootstrap.New(".")

	// Bootstrap application
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	kernel := app.Kernel()

	// Register console commands
	registerArtisanCommands(kernel)

	// Execute console command
	if err := kernel.Execute(); err != nil {
		log.Fatalf("Artisan error: %v", err)
	}
}

// registerArtisanCommands registers all Artisan CLI commands.
func registerArtisanCommands(kernel *console.Kernel) {
	// HTTP Server
	kernel.Register(console.ServeCommand(func(port string) error {
		app := bootstrap.New(".")
		app.Boot()

		// Placeholder: Start HTTP server
		return nil
	}))

	// Generators
	kernel.Register(console.MakeControllerCommand())
	kernel.Register(console.MakeModelCommand())
	kernel.Register(console.MakeMigrationCommand())
	kernel.Register(console.MakeMiddlewareCommand())

	// Database
	kernel.Register(console.MigrateCommand())
	kernel.Register(console.SeedCommand())

	// Additional commands can be registered here
}
