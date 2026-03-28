package command

import "context"

// ConsoleCommand represents an executable CLI command in the exchange application.
type ConsoleCommand interface {
	// Run executes the command with the given context and command-line flags.
	Run(ctx context.Context,flags []string)
}



