#!/bin/bash
# task-progress: PostToolUse hook on Edit|Write|MultiEdit
# Auto-checks task checkboxes when edited files match task descriptions
# No AI cooperation needed — directly modifies task.md

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
echo "$TARGET" | grep -qE 'docs/rein/changes/.*task\.md$' && exit 0

# Extract short filename for matching
EDITED_FILE=$(basename "$TARGET")

CHANGES_DIR="${CLAUDE_PROJECT_DIR}/docs/rein/changes"
[ -d "$CHANGES_DIR" ] || exit 0

MATCHED_TASK=""
MATCHED_TASKFILE=""

# Extract file references from a line — backtick paths + plain filenames with code extensions
extract_refs() {
    local line="$1"
    # 1. Backtick-enclosed references
    local bt_refs=$(echo "$line" | grep -oE '`[^`]+`' | sed 's/`//g')
    # 2. Plain filenames with common code extensions (excludes task numbers like 1.1)
    local plain_refs=$(echo "$line" | grep -oE '[A-Za-z0-9_/.-]+\.(go|ts|tsx|js|jsx|py|rs|java|rb|sql|yaml|yml|json|toml|proto|graphql|css|scss|html|sh|ps1|mod|sum|env|conf|xml|dart|swift|kt|c|cpp|h|hpp|php|tf|lock|txt|md)')
    echo "$bt_refs $plain_refs"
}

# Scan each feature directory
for feature_dir in "$CHANGES_DIR"/*/; do
    [ -d "$feature_dir" ] || continue
    taskfile="$feature_dir/task.md"
    planfile="$feature_dir/plan.md"

    # Phase 1: Scan task.md for file references in task lines
    if [ -f "$taskfile" ]; then
        while IFS= read -r line; do
            # Only process unchecked tasks
            echo "$line" | grep -qE '^\s*- \[ \]' || continue

            # Extract task number (e.g., "1.1" from "- [ ] 1.1 ...")
            TASK_NUM=$(echo "$line" | sed -n 's/^\s*- \[ \] \([0-9]\+\.[0-9]\+\).*/\1/p')
            [ -n "$TASK_NUM" ] || continue

            REFS=$(extract_refs "$line")
            for ref in $REFS; do
                [ -z "$ref" ] && continue
                REF_BASE=$(basename "$ref")
                if [ "$REF_BASE" = "$EDITED_FILE" ]; then
                    MATCHED_TASK="$TASK_NUM"
                    MATCHED_TASKFILE="$taskfile"
                    break 2
                fi
            done
        done < "$taskfile"
    fi

    # Phase 2: If no match in task.md, scan plan.md **Files:** fields
    if [ -z "$MATCHED_TASK" ] && [ -f "$planfile" ] && [ -f "$taskfile" ]; then
        CURRENT_TASK=""
        while IFS= read -r line; do
            [ -n "$MATCHED_TASK" ] && break
            # Track current task section: ### 1.1 ...
            if echo "$line" | grep -qE '^### [0-9]+\.[0-9]+'; then
                CURRENT_TASK=$(echo "$line" | sed -n 's/^### \([0-9]\+\.[0-9]\+\).*/\1/p')
            fi
            # Check **Files:** line within a task section
            if [ -n "$CURRENT_TASK" ] && echo "$line" | grep -qE '\*\*Files\*\*:'; then
                REFS=$(extract_refs "$line")
                for ref in $REFS; do
                    [ -z "$ref" ] && continue
                    REF_BASE=$(basename "$ref")
                    if [ "$REF_BASE" = "$EDITED_FILE" ]; then
                        # Find this unchecked task in task.md
                        if grep -qE "^\s*- \[ \] ${CURRENT_TASK}" "$taskfile"; then
                            MATCHED_TASK="$CURRENT_TASK"
                            MATCHED_TASKFILE="$taskfile"
                            break
                        fi
                    fi
                done
            fi
        done < "$planfile"
    fi

    [ -n "$MATCHED_TASK" ] && break
done

if [ -n "$MATCHED_TASK" ] && [ -n "$MATCHED_TASKFILE" ]; then
    # Escape dots for sed pattern
    TASK_NUM_ESC=$(echo "$MATCHED_TASK" | sed 's/\./\\./g')

    # Auto-check (allow optional leading whitespace, preserve it)
    sed -i "s/^\(\s*\)- \[ \] ${TASK_NUM_ESC}/\1- [x] ${MATCHED_TASK}/" "$MATCHED_TASKFILE"

    # Inject confirmation
    MSG="Auto-checked task ${MATCHED_TASK} (file match: ${EDITED_FILE})"
    MSG_ESCAPED=$(echo "$MSG" | sed 's/\\/\\\\/g; s/"/\\"/g')
    echo "{\"hookSpecificOutput\": {\"hookEventName\": \"PostToolUse\", \"additionalContext\": \"$MSG_ESCAPED\"}}"
else
    # No auto-match found — inject open task reminder
    for feature_dir in "$CHANGES_DIR"/*/; do
        [ -d "$feature_dir" ] || continue
        taskfile="$feature_dir/task.md"
        [ -f "$taskfile" ] || continue
        FEATURE_NAME=$(basename "$feature_dir")
        OPEN_TASKS=""
        while IFS= read -r line; do
            if echo "$line" | grep -qE '^\s*- \[ \]'; then
                TASK_DESC=$(echo "$line" | sed 's/^\s*- \[ \] //')
                if [ -n "$OPEN_TASKS" ]; then
                    OPEN_TASKS="${OPEN_TASKS}; ${TASK_DESC}"
                else
                    OPEN_TASKS="${TASK_DESC}"
                fi
            fi
        done < "$taskfile"
        if [ -n "$OPEN_TASKS" ]; then
            MSG="[task-reminder] ${FEATURE_NAME}: unchecked tasks remain — ${OPEN_TASKS}. If you just completed one, update task.md now."
            MSG_ESCAPED=$(echo "$MSG" | sed 's/\\/\\\\/g; s/"/\\"/g')
            echo "{\"hookSpecificOutput\": {\"hookEventName\": \"PostToolUse\", \"additionalContext\": \"$MSG_ESCAPED\"}}"
            break
        fi
    done
fi
