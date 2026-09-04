#!/bin/bash
set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TESTDATA_DIR="$SCRIPT_DIR/testdata"

MDP_BIN="${MDP_BIN:-$PROJECT_ROOT/mdp}"

if [[ ! -x "$MDP_BIN" ]]; then
    echo "Error: mdp binary not found at $MDP_BIN"
    exit 1
fi

FAILED=0

tmpdir=$(mktemp -d)
trap "rm -rf $tmpdir" EXIT

# generated_path prints where mdp writes the HTML for an absolute markdown path
generated_path() {
    local output_dir="$1" abs_md_path="$2"
    local path_without_ext="${abs_md_path%.md}"
    echo "$output_dir/${path_without_ext#/}/index.html"
}

# check_output compares a generated HTML file with its expected HTML.
# Prints PASS/FAIL with the given label and returns non-zero on failure.
check_output() {
    local label="$1" expected_file="$2" generated_file="$3" work_dir="$4"

    if [[ ! -f "$generated_file" ]]; then
        echo "FAIL: $label (output not found)"
        return 1
    fi

    # Expected html uses placeholders for the directories in file:// image URLs
    local resolved_expected_file="$work_dir/expected.html"
    sed -e "s|__TESTDATA_DIR__|$TESTDATA_DIR|g" -e "s|__E2E_DIR__|$SCRIPT_DIR|g" "$expected_file" > "$resolved_expected_file"

    if diff -q "$resolved_expected_file" "$generated_file" > /dev/null 2>&1; then
        echo "PASS: $label"
        return 0
    fi

    echo "FAIL: $label (content mismatch)"
    echo "--- Expected ---"
    cat "$resolved_expected_file"
    echo "--- Actual ---"
    cat "$generated_file"
    echo "----------------"
    return 1
}

for md_file in "$TESTDATA_DIR"/*.md; do
    name=$(basename "$md_file" .md)
    expected_file="$TESTDATA_DIR/$name.html"

    if [[ ! -f "$expected_file" ]]; then
        echo "SKIP: $name (no expected html)"
        continue
    fi

    test_dir="$tmpdir/$name"
    mkdir -p "$test_dir"

    config_file="$test_dir/config.yaml"
    output_dir="$test_dir/output"

    # Check if theme file exists for this test
    theme_file="$TESTDATA_DIR/$name.theme.html"
    if [[ -f "$theme_file" ]]; then
        # Setup theme
        themes_dir="$test_dir/themes"
        mkdir -p "$themes_dir"
        cp "$theme_file" "$themes_dir/test-theme.html"
        cat > "$config_file" <<EOF
output_dir: $output_dir
browser_command: echo
theme: test-theme
EOF
    else
        cat > "$config_file" <<EOF
output_dir: $output_dir
browser_command: echo
EOF
    fi

    # Append extra config if exists
    extra_config="$TESTDATA_DIR/$name.config.yaml"
    if [[ -f "$extra_config" ]]; then
        cat "$extra_config" >> "$config_file"
    fi

    abs_md_path=$(cd "$(dirname "$md_file")" && pwd)/$(basename "$md_file")

    if ! "$MDP_BIN" --config "$config_file" "$abs_md_path" > /dev/null 2>&1; then
        echo "FAIL: $name (command failed)"
        FAILED=1
        continue
    fi

    generated_file=$(generated_path "$output_dir" "$abs_md_path")
    if ! check_output "$name" "$expected_file" "$generated_file" "$test_dir"; then
        FAILED=1
    fi
done

# Multiple files: every file is converted and the browser is opened once per file
multi_name="multiple-files"
multi_dir="$tmpdir/$multi_name"
mkdir -p "$multi_dir"
multi_config="$multi_dir/config.yaml"
multi_output="$multi_dir/output"

# The browser command records every path it is asked to open
multi_opened_log="$multi_dir/opened.log"
multi_browser="$multi_dir/browser.sh"
cat > "$multi_browser" <<EOF
#!/bin/bash
echo "\$1" >> "$multi_opened_log"
EOF
chmod +x "$multi_browser"

cat > "$multi_config" <<EOF
output_dir: $multi_output
browser_command: $multi_browser
EOF

multi_files=("$TESTDATA_DIR/simple.md" "$TESTDATA_DIR/gfm.md")

if ! "$MDP_BIN" --config "$multi_config" "${multi_files[@]}" > /dev/null 2>&1; then
    echo "FAIL: $multi_name (command failed)"
    FAILED=1
else
    touch "$multi_opened_log"
    for md_file in "${multi_files[@]}"; do
        name=$(basename "$md_file" .md)
        generated_file=$(generated_path "$multi_output" "$md_file")
        if ! check_output "$multi_name ($name)" "$TESTDATA_DIR/$name.html" "$generated_file" "$multi_dir"; then
            FAILED=1
            continue
        fi

        opened=$(grep -Fxc "$generated_file" "$multi_opened_log" || true)
        if [[ "$opened" != "1" ]]; then
            echo "FAIL: $multi_name ($name) (browser opened $opened times, want 1)"
            FAILED=1
        fi
    done
fi

exit $FAILED
