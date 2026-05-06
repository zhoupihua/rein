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

  # Scan for active features and their status
  CHANGES_DIR="${CLAUDE_PROJECT_DIR}/docs/rein/changes"
  ACTIVE_MSG=""
  if [ -d "$CHANGES_DIR" ]; then
    for feature_dir in "$CHANGES_DIR"/*/; do
      [ -d "$feature_dir" ] || continue
      FNAME=$(basename "$feature_dir")

      # Check for active tasks
      TASKFILE="$feature_dir/task.md"
      if [ -f "$TASKFILE" ]; then
        UNCHECKED=$(grep -cE '^\s*- \[ \]' "$TASKFILE" 2>/dev/null || echo "0")
        if [ "$UNCHECKED" -gt 0 ]; then
          ACTIVE_MSG="\\n\\nACTIVE TASKS: $UNCHECKED unchecked task(s) in $FNAME. Use /continue to resume or /status to check progress."
        fi
      fi

      # Check phase completeness
      MISSING_MSG=""
      # DEFINE: refine.md, spec.md, design.md
      DEFINE_MISSING=""
      [ -f "$feature_dir/refine.md" ] || DEFINE_MISSING="$DEFINE_MISSING refine.md"
      [ -f "$feature_dir/spec.md" ] || DEFINE_MISSING="$DEFINE_MISSING spec.md"
      [ -f "$feature_dir/design.md" ] || DEFINE_MISSING="$DEFINE_MISSING design.md"
      if [ -n "$DEFINE_MISSING" ] && [ -f "$feature_dir/refine.md" -o -f "$feature_dir/spec.md" -o -f "$feature_dir/design.md" ]; then
        MISSING_MSG="\\n⚠️ DEFINE 阶段不完整，缺少:$DEFINE_MISSING"
      fi

      # PLAN: plan.md, task.md
      PLAN_MISSING=""
      [ -f "$feature_dir/plan.md" ] || PLAN_MISSING="$PLAN_MISSING plan.md"
      [ -f "$feature_dir/task.md" ] || PLAN_MISSING="$PLAN_MISSING task.md"
      if [ -n "$PLAN_MISSING" ] && [ -f "$feature_dir/plan.md" -o -f "$feature_dir/task.md" ]; then
        MISSING_MSG="${MISSING_MSG}\\n⚠️ PLAN 阶段不完整，缺少:$PLAN_MISSING"
      fi

      # REVIEW: review.md (only check if tasks are all done)
      if [ -f "$TASKFILE" ]; then
        UNCHECKED_R=$(grep -cE '^\s*- \[ \]' "$TASKFILE" 2>/dev/null || echo "1")
        if [ "$UNCHECKED_R" -eq 0 ] && [ ! -f "$feature_dir/review.md" ]; then
          MISSING_MSG="${MISSING_MSG}\\n⚠️ REVIEW 阶段不完整，缺少: review.md"
        fi
      fi

      if [ -n "$MISSING_MSG" ]; then
        ACTIVE_MSG="${ACTIVE_MSG}${MISSING_MSG}"
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
