#!/bin/bash
# task-progress: PostToolUse hook on Edit|Write|MultiEdit
# Injects task progress after code edits, making checkbox state visible to AI
# Inspired by OpenSpec's CLI feedback loop — AI sees progress and naturally corrects

TOOL_INPUT="$CLAUDE_TOOL_INPUT"
if [ -z "$TOOL_INPUT" ] && [ -n "$CLAUDE_TOOL_INPUT_FILE_PATH" ] && [ -f "$CLAUDE_TOOL_INPUT_FILE_PATH" ]; then
    TOOL_INPUT=$(cat "$CLAUDE_TOOL_INPUT_FILE_PATH")
fi
[ -n "$TOOL_INPUT" ] || exit 0

# Extract target file path
TARGET=$(echo "$TOOL_INPUT" | sed -n 's/.*"file_path"\s*:\s*"\([^"]*\)".*/\1/p')
[ -n "$TARGET" ] || exit 0
TARGET=$(echo "$TARGET" | sed 's/\\\\/\\/g' | tr '\\' '/')

# Skip when editing task.md (avoid recursive noise)
echo "$TARGET" | grep -qE 'docs/rein/tasks/.*task\.md$' && exit 0

# Parse task.md progress
TASKS_DIR="${CLAUDE_PROJECT_DIR}/docs/rein/tasks"
[ -d "$TASKS_DIR" ] || exit 0

TOTAL=0; COMPLETE=0; UNCHECKED_LIST=""
for taskfile in "$TASKS_DIR"/*task.md; do
    [ -f "$taskfile" ] || continue
    while IFS= read -r line; do
        if echo "$line" | grep -qE '^\s*- \[[xX]\]'; then
            TOTAL=$((TOTAL + 1)); COMPLETE=$((COMPLETE + 1))
        elif echo "$line" | grep -qE '^\s*- \[ \]'; then
            TOTAL=$((TOTAL + 1))
            DESC=$(echo "$line" | sed 's/^\s*- \[ \] //' | head -c 60)
            UNCHECKED_LIST="${UNCHECKED_LIST}${DESC}, "
        fi
    done < "$taskfile"
done

# No tasks or all complete → silent
[ "$TOTAL" -eq 0 ] && exit 0
REMAINING=$((TOTAL - COMPLETE))
[ "$REMAINING" -eq 0 ] && exit 0

# Trim trailing ", "
UNCHECKED_LIST=$(echo "$UNCHECKED_LIST" | sed 's/, $//')

# Inject progress
MSG="Task Progress: ${COMPLETE}/${TOTAL}. Unchecked: ${UNCHECKED_LIST}. If you completed a task, update its checkbox in task.md: \`- [ ]\` → \`- [x]\`"
MSG_ESCAPED=$(echo "$MSG" | sed 's/\\/\\\\/g; s/"/\\"/g')
echo "{\"hookSpecificOutput\": {\"additionalContext\": \"$MSG_ESCAPED\"}}"
