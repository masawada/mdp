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

for md_file in "$TESTDATA_DIR"/*.md; do
    name=$(basename "$md_file" .md)
    expected_file="$TESTDATA_DIR/$name.html"

    if [[ ! -f "$expected_file" ]]; then
        echo "SKIP: $name (no expected html)"
        continue
    fi

    tmpdir=$(mktemp -d)
    trap "rm -rf $tmpdir" EXIT

    config_file="$tmpdir/config.yaml"
    output_dir="$tmpdir/output"
    cat > "$config_file" <<EOF
output_dir: $output_dir
browser_command: echo
EOF

    abs_md_path=$(cd "$(dirname "$md_file")" && pwd)/$(basename "$md_file")

    if ! "$MDP_BIN" --config "$config_file" "$abs_md_path" > /dev/null 2>&1; then
        echo "FAIL: $name (command failed)"
        FAILED=1
        continue
    fi

    path_without_ext="${abs_md_path%.md}"
    relative_path="${path_without_ext#/}"
    generated_file="$output_dir/$relative_path/index.html"

    if [[ ! -f "$generated_file" ]]; then
        echo "FAIL: $name (output not found)"
        FAILED=1
        continue
    fi

    if diff -q "$expected_file" "$generated_file" > /dev/null 2>&1; then
        echo "PASS: $name"
    else
        echo "FAIL: $name (content mismatch)"
        echo "--- Expected ---"
        cat "$expected_file"
        echo "--- Actual ---"
        cat "$generated_file"
        echo "----------------"
        FAILED=1
    fi
done

exit $FAILED
