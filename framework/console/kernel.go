package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Kernel manages CLI commands.
type Kernel struct {
	commands map[string]*cobra.Command
	root     *cobra.Command
}

// NewKernel creates a new console kernel.
func NewKernel() *Kernel {
	root := &cobra.Command{
		Use:   "artisan",
		Short: "Govel CLI Tool",
		Long:  "A Laravel-inspired CLI tool for Govel framework",
	}

	return &Kernel{
		commands: make(map[string]*cobra.Command),
		root:     root,
	}
}

// Register adds a command to the kernel.
func (k *Kernel) Register(cmd *cobra.Command) {
	k.commands[cmd.Use] = cmd
	k.root.AddCommand(cmd)
}

// Execute runs the console kernel.
func (k *Kernel) Execute() error {
	return k.root.Execute()
}

// Command creates a new console command.
type Command struct {
	Signature   string
	Description string
	Handler     func(cmd *cobra.Command, args []string) error
}

// MakeCommand creates a cobra command from a Govel command.
func MakeCommand(c Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.Signature,
		Short: c.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.Handler != nil {
				return c.Handler(cmd, args)
			}
			return nil
		},
	}
	return cmd
}

// Built-in commands

// ServeCommand starts the development server.
func ServeCommand(handler func(port string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the development server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetString("port")
			if port == "" {
				port = "8000"
			}
			fmt.Printf("Govel server running on http://localhost:%s\n", port)
			return handler(port)
		},
	}
}

// MakeControllerCommand generates a new controller.
func MakeControllerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:controller",
		Short: "Create a new controller",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			// TODO: Implement controller generation
			fmt.Printf("Created controller: %s\n", name)
			return nil
		},
	}
}

// MakeModelCommand generates a new model.
func MakeModelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:model",
		Short: "Create a new model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			// TODO: Implement model generation
			fmt.Printf("Created model: %s\n", name)
			return nil
		},
	}
}

// MakeMigrationCommand generates a new migration.
func MakeMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:migration",
		Short: "Create a new migration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			// TODO: Implement migration generation
			fmt.Printf("Created migration: %s\n", name)
			return nil
		},
	}
}

// MakeMiddlewareCommand generates a new middleware.
func MakeMiddlewareCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:middleware",
		Short: "Create a new middleware",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			// TODO: Implement middleware generation
			fmt.Printf("Created middleware: %s\n", name)
			return nil
		},
	}
}

// MigrateCommand runs database migrations.
func MigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement migration runner
			fmt.Println("Running migrations...")
			return nil
		},
	}
}

// SeedCommand seeds the database.
func SeedCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "db:seed",
		Short: "Seed the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement database seeder
			fmt.Println("Seeding database...")
			return nil
		},
	}
}
