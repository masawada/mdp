package browser

import "testing"

func TestNewOpener(t *testing.T) {
	opener := NewOpener("firefox")
	if opener.command != "firefox" {
		t.Errorf("NewOpener().command = %q, want %q", opener.command, "firefox")
	}
}
