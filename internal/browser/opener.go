// Package browser provides functionality to open files in a browser.
package browser

import "os/exec"

// Opener opens files using a specified browser command.
type Opener struct {
	command string
}

// NewOpener creates a new Opener with the given browser command.
func NewOpener(command string) *Opener {
	return &Opener{command: command}
}

// Open opens the specified file in the browser.
func (o *Opener) Open(filePath string) error {
	cmd := exec.Command(o.command, filePath)
	return cmd.Run()
}
