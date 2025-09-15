#!/bin/bash

# Check if heading argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 \"<heading_text>\""
    echo "Example: $0 \"To read 8 bytes at a time.\""
    exit 1
fi

# Get the heading from command line argument
HEADING="$1"

# Check if main.go exists
if [ ! -f "main.go" ]; then
    echo "Error: main.go not found in current directory"
    exit 1
fi

# Create exercises directory if it doesn't exist
mkdir -p exercises

# Find the highest exercise number in the exercises folder
LAST_NUM=0
if [ -d "exercises" ]; then
    for file in exercises/exercise_*.md; do
        if [ -f "$file" ]; then
            # Extract number from filename
            NUM=$(basename "$file" | sed 's/exercise_\([0-9]*\)\.md/\1/')
            if [ "$NUM" -gt "$LAST_NUM" ]; then
                LAST_NUM="$NUM"
            fi
        fi
    done
fi

# Calculate next exercise number
NEXT_NUM=$((LAST_NUM + 1))

# Create filename
FILENAME="exercises/exercise_${NEXT_NUM}.md"

# Read main.go content
MAIN_GO_CONTENT=$(cat main.go)

# Create the markdown file
cat > "$FILENAME" << EOF
##### ${HEADING}
\`\`\`
${MAIN_GO_CONTENT}
\`\`\`
EOF

echo "Created: $FILENAME"
echo "Exercise number: $NEXT_NUM"
echo "Heading: $HEADING"