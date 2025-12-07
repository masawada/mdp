package browser

import "testing"

func TestNewOpener(t *testing.T) {
	opener := NewOpener("firefox")
	if opener.command != "firefox" {
		t.Errorf("NewOpener().command = %q, want %q", opener.command, "firefox")
	}
}

func TestOpen(t *testing.T) {
	opener := NewOpener("echo")
	err := opener.Open("/path/to/file.html")
	if err != nil {
		t.Errorf("Open() error: %v", err)
	}
}

func TestOpen_InvalidCommand(t *testing.T) {
	opener := NewOpener("nonexistent-command-12345")
	err := opener.Open("/path/to/file.html")
	if err == nil {
		t.Error("Open() should return error for invalid command")
	}
}
