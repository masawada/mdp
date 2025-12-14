# mdp

Markdown previewer - preview markdown files in your browser.

## Description

`mdp` converts a markdown file to HTML and opens it in your browser. It supports GitHub Flavored Markdown and custom themes.

```console
$ mdp README.md
Generated: /Users/you/.mdp/README.html
```

## Synopsis

```
mdp [options] <markdown-file>
```

## Options

```
--config <config-file>  path to config file
--list                  list generated files
--help                  show help message
```

## Installation

### Download binary

Download the latest binary from [Releases](https://github.com/masawada/mdp/releases) and place it in your `$PATH`.

### Go

```console
$ go install github.com/masawada/mdp/cmd/mdp@latest
```

### Build from source

```console
$ git clone https://github.com/masawada/mdp.git
$ cd mdp
$ make
```

## Configuration

Configuration file is loaded from the following location:

- macOS: `~/Library/Application Support/mdp/config.yaml`
- Linux: `~/.config/mdp/config.yaml` (or `$XDG_CONFIG_HOME/mdp/config.yaml`)

```yaml
# Output directory for generated HTML files (default: ~/.mdp)
output_dir: ~/.mdp

# Command to open browser (default: open on macOS, xdg-open on Linux)
browser_command: open

# Theme name (optional, looks for themes/<name>.html in config directory)
theme: custom
```

## Themes

You can create custom themes by placing HTML template files in the `themes/` directory under your config directory.

For example, to use a theme named `custom`, create `themes/custom.html` in your config directory:

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
  </style>
</head>
<body>
  {{.Content}}
</body>
</html>
```

The `{{.Content}}` placeholder will be replaced with the rendered HTML content.
