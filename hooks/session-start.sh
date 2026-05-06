#!/bin/bash
# Session start hook for rein
# Injects the using-rein meta-skill into every new Claude Code session

SKILL_FILE="$(dirname "$0")/../skills/using-rein/SKILL.md"

if [ ! -f "$SKILL_FILE" ]; then
  # Try relative to the script's resolved location
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
  SKILL_FILE="$SCRIPT_DIR/../skills/using-rein/SKILL.md"
fi

if [ -f "$SKILL_FILE" ]; then
  CONTENT=$(cat "$SKILL_FILE" | sed 's/\\/\\\\/g; s/"/\\"/g; s/\t/\\t/g' | sed ':a;N;$!ba;s/\n/\\n/g')

  # Scan for active tasks
  CHANGES_DIR="${CLAUDE_PROJECT_DIR}/docs/rein/changes"
  ACTIVE_MSG=""
  if [ -d "$CHANGES_DIR" ]; then
    for taskfile in "$CHANGES_DIR"/*/task.md; do
      [ -f "$taskfile" ] || continue
      UNCHECKED=$(grep -cE '^\s*- \[ \]' "$taskfile" 2>/dev/null || echo "0")
      if [ "$UNCHECKED" -gt 0 ]; then
        FNAME=$(basename "$taskfile")
        ACTIVE_MSG="\\n\\nACTIVE TASKS: $UNCHECKED unchecked task(s) in $FNAME. Use /continue to resume or /status to check progress."
        break
      fi
    done
  fi
  CONTENT="${CONTENT}${ACTIVE_MSG}"

  # Detect platform and output appropriate format
  if [ -n "$CURSOR_SESSION" ]; then
    # Cursor format
    cat <<EOF
{"additional_context": "$CONTENT"}
EOF
  else
    # Claude Code format (default)
    cat <<EOF
{"hookSpecificOutput": {"hookEventName": "SessionStart", "additionalContext": "$CONTENT"}}
EOF
  fi
else
  echo '{"hookSpecificOutput": {"hookEventName": "SessionStart", "additionalContext": "Warning: rein using-rein skill not found. Run install script to set up."}}' >&2
fi