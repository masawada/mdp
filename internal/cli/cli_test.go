package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run([]string{})
	if exitCode != 1 {
		t.Errorf("run() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "usage") {
		t.Errorf("stderr should contain usage message, got: %s", stderr.String())
	}
}
