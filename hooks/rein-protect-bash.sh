#!/bin/bash
# rein Bash Protection Hook (PreToolUse → Bash)
# Prevents destructive commands targeting rein-managed files

MANIFEST="${CLAUDE_PROJECT_DIR}/.claude/.rein-manifest"
[ -f "$MANIFEST" ] || exit 0

INPUT="$CLAUDE_TOOL_INPUT"

# Only check destructive commands
if ! echo "$INPUT" | grep -qE '(rm |rmdir |del |mv |sed -i |truncate|>\s*/|Remove-Item|Move-Item)'; then
    exit 0
fi

while IFS= read -r entry; do
    [[ "$entry" =~ ^[[:space:]]*# ]] && continue
    [[ -z "${entry// /}" ]] && continue
    entry=$(echo "$entry" | tr '\\' '/')
    if echo "$INPUT" | grep -qF "$entry"; then
        echo "{\"decision\":\"block\",\"reason\":\"Command targets rein-managed file: $entry. Use rein commands to update.\"}"
        exit 2
    fi
done < "$MANIFEST"

exit 0
