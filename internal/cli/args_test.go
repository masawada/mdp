package cli

import (
	"slices"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantArgs   *parsedArgs
		wantErrMsg string
	}{
		{
			name: "no config flag",
			args: []string{"test.md"},
			wantArgs: &parsedArgs{
				configPath: "",
				filePaths:  []string{"test.md"},
			},
		},
		{
			name: "config flag with value",
			args: []string{"--config", "/path/to/config.yaml", "test.md"},
			wantArgs: &parsedArgs{
				configPath: "/path/to/config.yaml",
				filePaths:  []string{"test.md"},
			},
		},
		{
			name: "multiple files",
			args: []string{"a.md", "b.md", "c.md"},
			wantArgs: &parsedArgs{
				filePaths: []string{"a.md", "b.md", "c.md"},
			},
		},
		{
			name: "multiple files with flags",
			args: []string{"--watch", "--config", "config.yaml", "a.md", "b.md"},
			wantArgs: &parsedArgs{
				filePaths:  []string{"a.md", "b.md"},
				configPath: "config.yaml",
				watchMode:  true,
			},
		},
		{
			name:       "config flag without value",
			args:       []string{"--config"},
			wantErrMsg: "flag needs an argument",
		},
		{
			name:       "no arguments",
			args:       []string{},
			wantErrMsg: "markdown file is required",
		},
		{
			name:       "only config flag",
			args:       []string{"--config", "/path/to/config.yaml"},
			wantErrMsg: "markdown file is required",
		},
		{
			name:       "help flag long form",
			args:       []string{"--help"},
			wantErrMsg: "help requested",
		},
		{
			name:       "help flag short form",
			args:       []string{"-h"},
			wantErrMsg: "help requested",
		},
		{
			name:       "help flag with file",
			args:       []string{"--help", "test.md"},
			wantErrMsg: "help requested",
		},
		{
			name: "list flag",
			args: []string{"--list"},
			wantArgs: &parsedArgs{
				showList: true,
			},
		},
		{
			name: "list flag with config",
			args: []string{"--config", "/path/to/config.yaml", "--list"},
			wantArgs: &parsedArgs{
				configPath: "/path/to/config.yaml",
				showList:   true,
			},
		},
		{
			name:       "help and list flags together",
			args:       []string{"--help", "--list"},
			wantErrMsg: "help requested",
		},
		{
			name: "watch flag",
			args: []string{"--watch", "test.md"},
			wantArgs: &parsedArgs{
				filePaths: []string{"test.md"},
				watchMode: true,
			},
		},
		{
			name: "watch flag with config",
			args: []string{"--watch", "--config", "config.yaml", "test.md"},
			wantArgs: &parsedArgs{
				filePaths:  []string{"test.md"},
				configPath: "config.yaml",
				watchMode:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if tt.wantErrMsg != "" {
				if err == nil {
					t.Errorf("parseArgs() expected error containing %q, got nil", tt.wantErrMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("parseArgs() error = %q, want error containing %q", err.Error(), tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("parseArgs() unexpected error = %v", err)
				return
			}
			if got.configPath != tt.wantArgs.configPath {
				t.Errorf("parseArgs() configPath = %v, want %v", got.configPath, tt.wantArgs.configPath)
			}
			if !slices.Equal(got.filePaths, tt.wantArgs.filePaths) {
				t.Errorf("parseArgs() filePaths = %v, want %v", got.filePaths, tt.wantArgs.filePaths)
			}
			if got.showList != tt.wantArgs.showList {
				t.Errorf("parseArgs() showList = %v, want %v", got.showList, tt.wantArgs.showList)
			}
			if got.watchMode != tt.wantArgs.watchMode {
				t.Errorf("parseArgs() watchMode = %v, want %v", got.watchMode, tt.wantArgs.watchMode)
			}
		})
	}
}
