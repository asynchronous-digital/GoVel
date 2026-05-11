package routes
package routes

import (
	"github.com/hasnat/govel/app/http/controllers"
	"github.com/hasnat/govel/framework/http"
	"github.com/hasnat/govel/framework/routing"
)

// RegisterWebRoutes registers web routes for the application.
func RegisterWebRoutes(router *routing.Router) {
	// Home route
	router.Get("/", http.NewHandler(func(ctx *http.Context) error {
		return controllers.HomeController{}.Index(ctx)
	}))

	// Health check
	router.Get("/health", http.NewHandler(func(ctx *http.Context) error {
		return controllers.HealthController{}.Check(ctx)
	}))
}

// RegisterAPIRoutes registers API routes for the application.
func RegisterAPIRoutes(router *routing.Router) {
	// API info
	router.Get("/api/info", http.NewHandler(func(ctx *http.Context) error {
		return controllers.APIController{}.Info(ctx)
	}))

	// User routes (REST API)
	userGroup := router.Prefix("/api/users")
	userGroup.Group(func(r *routing.Router) {
		// List users
		r.Get("", http.NewHandler(func(ctx *http.Context) error {
			return controllers.UserController{}.Index(ctx)
		}))

		// Create user
		r.Post("", http.NewHandler(func(ctx *http.Context) error {
			return controllers.UserController{}.Store(ctx)
		}))

		// Show user
		r.Get("/{id}", http.NewHandler(func(ctx *http.Context) error {
			return controllers.UserController{}.Show(ctx)
		}))

		// Update user
		r.Put("/{id}", http.NewHandler(func(ctx *http.Context) error {
			return controllers.UserController{}.Update(ctx)
		}))

		// Delete user
		r.Delete("/{id}", http.NewHandler(func(ctx *http.Context) error {
			return controllers.UserController{}.Delete(ctx)
		}))
	})

	// Post routes
	postGroup := router.Prefix("/api/posts")
	postGroup.Group(func(r *routing.Router) {
		r.Get("", http.NewHandler(func(ctx *http.Context) error {
			return controllers.PostController{}.Index(ctx)
		}))

		r.Get("/{id}", http.NewHandler(func(ctx *http.Context) error {
			return controllers.PostController{}.Show(ctx)
		}))
	})
}

// RegisterConsoleRoutes registers console commands.
func RegisterConsoleRoutes(kernel interface{}) {
	// Console commands would be registered here
}
