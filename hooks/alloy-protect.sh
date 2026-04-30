#!/bin/bash
# Alloy Protection Hook (PreToolUse → Edit|Write|MultiEdit)
# Prevents modification of Alloy-managed files

MANIFEST="${CLAUDE_PROJECT_DIR}/.claude/.alloy-manifest"
[ -f "$MANIFEST" ] || exit 0

# Extract file_path from tool input JSON
TARGET=$(echo "$CLAUDE_TOOL_INPUT" | sed -n 's/.*"file_path"\s*:\s*"\([^"]*\)".*/\1/p')
[ -n "$TARGET" ] || exit 0

# Normalize to forward slashes for cross-platform matching
TARGET=$(echo "$TARGET" | tr '\\' '/')

while IFS= read -r entry; do
    [[ "$entry" =~ ^[[:space:]]*# ]] && continue
    [[ -z "${entry// /}" ]] && continue
    entry=$(echo "$entry" | tr '\\' '/')
    if [[ "$TARGET" == *"$entry"* ]]; then
        echo '{"decision":"block","reason":"This file is managed by Alloy and cannot be modified. Use Alloy commands to update, or remove its path from .claude/.alloy-manifest to allow edits."}'
        exit 2
    fi
done < "$MANIFEST"

exit 0
