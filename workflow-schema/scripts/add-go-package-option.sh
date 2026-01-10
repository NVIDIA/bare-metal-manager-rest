#!/bin/bash

# Script to add go_package option to proto files if not already present
# Usage: ./add_go_package_option.sh <proto_file>

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <proto_file>"
    exit 1
fi

PROTO_FILE="$1"
GO_PACKAGE_OPTION='option go_package = "github.com/nvidia/carbide-rest/workflow-schema/proto";'

if [ ! -f "$PROTO_FILE" ]; then
    echo "Error: File '$PROTO_FILE' not found"
    exit 1
fi

# Check if the go_package option already exists
if grep -qF "$GO_PACKAGE_OPTION" "$PROTO_FILE"; then
    echo "File '$PROTO_FILE' already has go_package option. Skipping."
    exit 0
fi

# Find the line number of the last import statement
LAST_IMPORT_LINE=$(grep -n "^import " "$PROTO_FILE" | tail -n 1 | cut -d: -f1)

if [ -z "$LAST_IMPORT_LINE" ]; then
    # No import statements found, add after package declaration
    PACKAGE_LINE=$(grep -n "^package " "$PROTO_FILE" | head -n 1 | cut -d: -f1)

    if [ -z "$PACKAGE_LINE" ]; then
        # No package statement either, add after syntax declaration
        SYNTAX_LINE=$(grep -n "^syntax " "$PROTO_FILE" | head -n 1 | cut -d: -f1)

        if [ -z "$SYNTAX_LINE" ]; then
            echo "Error: Could not find syntax, package, or import statements in '$PROTO_FILE'"
            exit 1
        fi

        INSERT_LINE=$((SYNTAX_LINE + 1))
    else
        INSERT_LINE=$((PACKAGE_LINE + 1))
    fi
else
    # Insert after the last import statement
    INSERT_LINE=$((LAST_IMPORT_LINE + 1))
fi

# Create a temporary file
TEMP_FILE=$(mktemp)

# Insert the go_package option at the appropriate line
awk -v line="$INSERT_LINE" -v option="$GO_PACKAGE_OPTION" '
NR == line {
    print ""
    print option
}
{ print }
' "$PROTO_FILE" > "$TEMP_FILE"

# Replace the original file with the modified one
mv "$TEMP_FILE" "$PROTO_FILE"

echo "Added go_package option to '$PROTO_FILE'"
