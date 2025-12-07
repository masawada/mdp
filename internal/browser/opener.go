package browser

type Opener struct {
	command string
}

func NewOpener(command string) *Opener {
	return &Opener{command: command}
}
