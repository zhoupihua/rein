#!/bin/bash
# task-progress: PostToolUse hook on Edit|Write|MultiEdit
# Auto-checks task checkboxes when edited files match task descriptions
# No AI cooperation needed — directly modifies task.md via sed

TOOL_INPUT="$CLAUDE_TOOL_INPUT"
if [ -z "$TOOL_INPUT" ] && [ -n "$CLAUDE_TOOL_INPUT_FILE_PATH" ] && [ -f "$CLAUDE_TOOL_INPUT_FILE_PATH" ]; then
    TOOL_INPUT=$(cat "$CLAUDE_TOOL_INPUT_FILE_PATH")
fi
[ -n "$TOOL_INPUT" ] || exit 0

# Extract target file path
TARGET=$(echo "$TOOL_INPUT" | sed -n 's/.*"file_path"\s*:\s*"\([^"]*\)".*/\1/p')
[ -n "$TARGET" ] || exit 0
TARGET=$(echo "$TARGET" | sed 's/\\\\/\\/g' | tr '\\' '/')

# Skip task.md edits (avoid recursive triggers)
echo "$TARGET" | grep -qE 'docs/rein/tasks/.*task\.md$' && exit 0

# Extract short filename for matching
EDITED_FILE=$(basename "$TARGET")

TASKS_DIR="${CLAUDE_PROJECT_DIR}/docs/rein/tasks"
[ -d "$TASKS_DIR" ] || exit 0

MATCHED_TASK=""
MATCHED_TASKFILE=""

for taskfile in "$TASKS_DIR"/*task.md; do
    [ -f "$taskfile" ] || continue

    while IFS= read -r line; do
        # Only process unchecked tasks
        echo "$line" | grep -qE '^\s*- \[ \]' || continue

        # Extract task number (e.g., "1.1" from "- [ ] 1.1 ...")
        TASK_NUM=$(echo "$line" | sed -n 's/^\s*- \[ \] \([0-9]\+\.[0-9]\+\).*/\1/p')
        [ -n "$TASK_NUM" ] || continue

        # Extract backtick-enclosed file references
        REFS=$(echo "$line" | grep -oE '`[^`]+`' | sed 's/`//g')

        for ref in $REFS; do
            REF_BASE=$(basename "$ref")
            if [ "$REF_BASE" = "$EDITED_FILE" ]; then
                MATCHED_TASK="$TASK_NUM"
                MATCHED_TASKFILE="$taskfile"
                break 2
            fi
        done
    done < "$taskfile"
done

if [ -n "$MATCHED_TASK" ]; then
    # Escape dots for sed pattern
    TASK_NUM_ESC=$(echo "$MATCHED_TASK" | sed 's/\./\\./g')

    # Auto-check the matching task
    sed -i "s/^- \[ \] ${TASK_NUM_ESC}/- [x] ${MATCHED_TASK}/" "$MATCHED_TASKFILE"

    # Inject confirmation
    MSG="Auto-checked task ${MATCHED_TASK} (file match: ${EDITED_FILE})"
    MSG_ESCAPED=$(echo "$MSG" | sed 's/\\/\\\\/g; s/"/\\"/g')
    echo "{\"hookSpecificOutput\": {\"hookEventName\": \"PostToolUse\", \"additionalContext\": \"$MSG_ESCAPED\"}}"
fi
