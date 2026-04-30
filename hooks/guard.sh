#!/bin/bash
# guard Hook (PreToolUse → Edit|Write|MultiEdit)
# Prevents modification of rein-managed files

MANIFEST="${CLAUDE_PROJECT_DIR}/.claude/.rein-manifest"
[ -f "$MANIFEST" ] || exit 0

# Read tool input from env or file
TOOL_INPUT="$CLAUDE_TOOL_INPUT"
if [ -z "$TOOL_INPUT" ] && [ -n "$CLAUDE_TOOL_INPUT_FILE_PATH" ] && [ -f "$CLAUDE_TOOL_INPUT_FILE_PATH" ]; then
    TOOL_INPUT=$(cat "$CLAUDE_TOOL_INPUT_FILE_PATH")
fi
[ -n "$TOOL_INPUT" ] || exit 0

# Extract file_path from tool input JSON
TARGET=$(echo "$TOOL_INPUT" | sed -n 's/.*"file_path"\s*:\s*"\([^"]*\)".*/\1/p')
[ -n "$TARGET" ] || exit 0

# Unescape JSON \\ to \, then normalize to forward slashes
TARGET=$(echo "$TARGET" | sed 's/\\\\/\\/g' | tr '\\' '/')

while IFS= read -r entry; do
    [[ "$entry" =~ ^[[:space:]]*# ]] && continue
    [[ -z "${entry// /}" ]] && continue
    entry=$(echo "$entry" | tr '\\' '/')
    if [[ "$TARGET" == *"$entry"* ]]; then
        echo '{"decision":"block","reason":"This file is managed by rein and cannot be modified. Use rein commands to update, or remove its path from .claude/.rein-manifest to allow edits."}'
        exit 2
    fi
done < "$MANIFEST"

exit 0
