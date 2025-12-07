package browser

import "os/exec"

type Opener struct {
	command string
}

func NewOpener(command string) *Opener {
	return &Opener{command: command}
}

func (o *Opener) Open(filePath string) error {
	cmd := exec.Command(o.command, filePath)
	return cmd.Run()
}
