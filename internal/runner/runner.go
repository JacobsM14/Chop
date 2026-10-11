package runner

import (
	"context"
	"os"
	"os/exec"
)

// Run ejecuta un comando de shell conectado al terminal actual.
func Run(ctx context.Context, text string) error {
	c := exec.CommandContext(ctx, "sh", "-c", text)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
