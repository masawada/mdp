# mdp

A Markdown previewer that shows markdown files in your browser.

## Description

`mdp` converts markdown files to HTML and opens them in your browser. It supports GitHub Flavored Markdown and custom themes.

```console
$ mdp README.md
Generated: /Users/you/.mdp/README.html
```

You can pass more than one file. Each file is converted and opened in its own tab, and `--watch` regenerates whichever file changes.

```console
$ mdp --watch README.md CHANGELOG.md
Generated: /Users/you/.mdp/README.html
Generated: /Users/you/.mdp/CHANGELOG.html
Watching for changes... (Ctrl+C to stop)
```

## Synopsis

```
mdp [options] <markdown-file>...
```

## Options

```
--config <config-file>  path to config file
--watch                 watch for file changes and regenerate
--list                  list generated files
--help                  show help message
```

## Installation

### Homebrew (recommended)

```bash
brew install masawada/tap/mdp
```

### GitHub Releases

Download the latest archive from [GitHub Releases](https://github.com/masawada/mdp/releases) and install:

```bash
tar xzf mdp_*.tar.gz
sudo install mdp /usr/local/bin/
```

### Go

```console
$ go install github.com/masawada/mdp/cmd/mdp@latest
```

### Build from source

```bash
git clone https://github.com/masawada/mdp.git
cd mdp
make
sudo install mdp /usr/local/bin/
```

## Images

`mdp` rewrites relative image paths into `file://` URLs based on the directory of the markdown file, so the browser can show images next to the markdown without copying them. This works for both `![alt](path)` and raw `<img>` tags (with `unsafe: true`). URLs that already have a scheme, such as `https://` or `data:`, stay as they are.

## Configuration

`mdp` reads the file given by `--config`. Without the flag, it looks for a configuration file in these locations, in order:

1. `$UserConfigDir/mdp/config.yaml`
2. `$UserConfigDir/mdp/config.yml`
3. `$HOME/.config/mdp/config.yaml`
4. `$HOME/.config/mdp/config.yml`

If no file is found, it uses the defaults.

`$UserConfigDir` comes from `os.UserConfigDir()`:

- macOS: `~/Library/Application Support`
- Linux: `~/.config` (or `$XDG_CONFIG_HOME`)

```yaml
# Output directory for generated HTML files (default: ~/.mdp)
output_dir: ~/.mdp

# Command to open browser (default: open on macOS, xdg-open on Linux)
browser_command: open

# Theme name (optional, looks for themes/<name>.html in config directory)
theme: custom

# Convert newlines in paragraphs to <br> (default: false)
hard_wraps: false

# Allow raw HTML in markdown (default: false)
unsafe: false
```

## Themes

To create a custom theme, put an HTML template file in the `themes/` directory under your config directory.

For example, a theme named `custom` lives at `themes/custom.html`:

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>
  <style>
    body { font-family: sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
  </style>
</head>
<body>
  {{.Content}}
</body>
</html>
```

### Template variables

| Variable | Description |
|----------|-------------|
| `{{.Title}}` | Document title extracted from the markdown |
| `{{.Content}}` | Rendered HTML content |

### Title extraction

`mdp` picks the title from the markdown file in this order:

1. `title` field in YAML front-matter
2. First heading in the document
3. `"Untitled"` (default)

Example with front-matter:

```markdown
---
title: My Document Title
---

# Heading

Content here.
```

Here `{{.Title}}` is `"My Document Title"`.

## License

MIT
