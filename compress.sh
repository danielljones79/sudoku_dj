#!/bin/bash

# compress.sh - Script to zip sudoku_dj project excluding unnecessary files
# Creates a zip file in the parent directory
# Usage: ./compress.sh [suffix]
# Example: ./compress.sh v02 will create sudoku_dj_v02.zip

# Get the directory name of the current project
PROJECT_DIR=$(basename "$PWD")
PARENT_DIR=$(dirname "$PWD")

# Check if a suffix was provided
if [ -n "$1" ]; then
    OUTPUT_FILE="$PARENT_DIR/${PROJECT_DIR}_$1.zip"
else
    OUTPUT_FILE="$PARENT_DIR/$PROJECT_DIR.zip"
fi

echo "Creating zip file for $PROJECT_DIR..."
echo "Output will be saved to $OUTPUT_FILE"

# Remove existing zip file if it exists
if [ -f "$OUTPUT_FILE" ]; then
    rm "$OUTPUT_FILE"
    echo "Removed existing zip file."
fi

# Create a temporary exclusion list
EXCLUDE_FILE=$(mktemp)
cat > "$EXCLUDE_FILE" << EOF
*.DS_Store
*.log
*.iml
*.ipr
*.iws
.git/*
.gitignore
compress.sh
.vscode/*
.idea/*
.cursor/*
node_modules/*
build/*
dist/*
__pycache__/*
puzzles/*
EOF

# Create zip file with exclusions
zip -r "$OUTPUT_FILE" . -x@"$EXCLUDE_FILE"

# Check if zip was successful
if [ $? -eq 0 ]; then
    echo "Successfully created $OUTPUT_FILE"
    echo "Excluded: .git, .vscode, .idea, .cursor, node_modules, build, dist, __pycache__, puzzles, and IDE/system files"
    # Verify exclusions
    echo "Verifying zip contents..."
    unzip -l "$OUTPUT_FILE" | grep -E "node_modules/|.git/|.vscode/|puzzles/" > /dev/null
    if [ $? -eq 0 ]; then
        echo "Warning: Some excluded directories might still be in the zip file."
        echo "You can check the contents with: unzip -l \"$OUTPUT_FILE\""
    else
        echo "Verification passed. Excluded directories are not in the zip file."
    fi
else
    echo "Error creating zip file!"
    rm -f "$EXCLUDE_FILE"
    exit 1
fi

# Clean up
rm -f "$EXCLUDE_FILE" 