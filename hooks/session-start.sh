#!/bin/bash
# Session start hook for rein
# Injects the using-workflow meta-skill into every new Claude Code session

SKILL_FILE="$(dirname "$0")/../skills/using-workflow/SKILL.md"

if [ ! -f "$SKILL_FILE" ]; then
  # Try relative to the script's resolved location
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
  SKILL_FILE="$SCRIPT_DIR/../skills/using-workflow/SKILL.md"
fi

if [ -f "$SKILL_FILE" ]; then
  CONTENT=$(cat "$SKILL_FILE" | sed 's/\\/\\\\/g; s/"/\\"/g; s/\t/\\t/g' | sed ':a;N;$!ba;s/\n/\\n/g')

  # Detect platform and output appropriate format
  if [ -n "$CURSOR_SESSION" ]; then
    # Cursor format
    cat <<EOF
{"additional_context": "$CONTENT"}
EOF
  else
    # Claude Code format (default)
    cat <<EOF
{"hookSpecificOutput": {"additionalContext": "$CONTENT"}}
EOF
  fi
else
  echo '{"hookSpecificOutput": {"additionalContext": "Warning: rein using-workflow skill not found. Run install script to set up."}}' >&2
fi