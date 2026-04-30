#!/bin/bash
# checkbox-guard Hook (PostToolUse → Edit|Write|MultiEdit)
# Warns when a task.md file is edited without toggling a checkbox

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

# Only trigger on task.md files in docs/rein/tasks/
echo "$TARGET" | grep -qE 'docs/rein/tasks/.*task\.md$' || exit 0

# Check if the file exists
[ -f "$TARGET" ] || exit 0

# Check tool input for checkbox toggle evidence ([x] in the edit content)
echo "$TOOL_INPUT" | grep -qE '\[x\]' && exit 0

# Also check tool result
TOOL_RESULT="$CLAUDE_TOOL_RESULT"
if [ -n "$TOOL_RESULT" ] && echo "$TOOL_RESULT" | grep -qE '\[x\]'; then
    exit 0
fi

# Task.md was edited but no checkbox was toggled - inject warning
MSG="WARNING: You edited a task file but did not toggle any checkbox from [ ] to [x]. If you completed a task, you MUST update its checkbox NOW. The /do loop will re-find the same task until its checkbox is updated."
# Escape for JSON
MSG_ESCAPED=$(echo "$MSG" | sed 's/\\/\\\\/g; s/"/\\"/g')
echo "{\"hookSpecificOutput\": {\"additionalContext\": \"$MSG_ESCAPED\"}}"
