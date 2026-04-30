#!/bin/bash
# rein install script (Linux/Mac)
# Run from your project root: bash /path/to/rein/install/install.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKFLOW_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(pwd)"

echo "=== rein Installer ==="
echo "Workflow source: $WORKFLOW_DIR"
echo "Target project:  $PROJECT_DIR"
echo ""

# 1. Create artifact directories
echo "[1/9] Creating artifact directories..."
mkdir -p "$PROJECT_DIR/docs/rein/specs"
mkdir -p "$PROJECT_DIR/docs/rein/plans"
mkdir -p "$PROJECT_DIR/docs/rein/tasks"
mkdir -p "$PROJECT_DIR/docs/rein/archive"
echo "  ✓ docs/rein/specs/, docs/rein/plans/, docs/rein/tasks/, docs/rein/archive/"

# 2. Copy commands
echo "[2/9] Installing commands..."
mkdir -p "$PROJECT_DIR/.claude/commands"
cp "$WORKFLOW_DIR/commands/"*.md "$PROJECT_DIR/.claude/commands/"
CMD_COUNT=$(ls "$PROJECT_DIR/.claude/commands/"*.md | wc -l)
echo "  ✓ $CMD_COUNT commands installed"

# 3. Copy skills
echo "[3/9] Installing skills..."
mkdir -p "$PROJECT_DIR/.claude/skills"
cp -r "$WORKFLOW_DIR/skills/"* "$PROJECT_DIR/.claude/skills/"
SKILL_COUNT=$(ls -d "$PROJECT_DIR/.claude/skills/"*/ | wc -l)
echo "  ✓ $SKILL_COUNT skills installed"

# 4. Copy agents
echo "[4/9] Installing agents..."
mkdir -p "$PROJECT_DIR/.claude/agents"
cp "$WORKFLOW_DIR/agents/"*.md "$PROJECT_DIR/.claude/agents/"
AGENT_COUNT=$(ls "$PROJECT_DIR/.claude/agents/"*.md | wc -l)
echo "  ✓ $AGENT_COUNT agents installed"

# 5. Copy hooks
echo "[5/9] Installing hooks..."
mkdir -p "$PROJECT_DIR/.claude/hooks"
for hook in session-start format gate leak-guard inject guard guard-bash checkbox-guard task-progress; do
  cp "$WORKFLOW_DIR/hooks/${hook}.sh" "$PROJECT_DIR/.claude/hooks/"
  if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    cp "$WORKFLOW_DIR/hooks/${hook}.ps1" "$PROJECT_DIR/.claude/hooks/"
  fi
done
chmod +x "$PROJECT_DIR/.claude/hooks/"*.sh
echo "  ✓ All hooks installed"

# 6. Copy checklists
echo "[6/9] Installing checklists..."
mkdir -p "$PROJECT_DIR/.claude/checklists"
if [ -f "$WORKFLOW_DIR/templates/checklists/review.md" ]; then
  cp "$WORKFLOW_DIR/templates/checklists/review.md" "$PROJECT_DIR/.claude/checklists/"
  echo "  ✓ review.md checklist installed"
else
  echo "  ⚠ No review checklist template found"
fi

# 7. Generate manifest
echo "[7/9] Generating protection manifest..."
MANIFEST_FILE="$PROJECT_DIR/.claude/.rein-manifest"
{
  echo "# rein Managed Files - DO NOT EDIT"
  echo "# These files are protected from modification by the guard hook."
  echo "# To allow edits to a specific file, remove its line from this manifest."
  echo ""
  # List individual files in flat directories
  for dir in hooks commands agents checklists; do
    if [ -d "$PROJECT_DIR/.claude/$dir" ]; then
      find "$PROJECT_DIR/.claude/$dir" -maxdepth 1 -type f | while read -r f; do
        echo "${f#$PROJECT_DIR/}"
      done
    fi
  done
  # Skills are directories - list each skill dir as a prefix
  if [ -d "$PROJECT_DIR/.claude/skills" ]; then
    find "$PROJECT_DIR/.claude/skills" -mindepth 1 -maxdepth 1 -type d | while read -r d; do
      echo "${d#$PROJECT_DIR/}/"
    done
  fi
} > "$MANIFEST_FILE"
MANIFEST_COUNT=$(grep -cvE '^\s*#|^\s*$' "$MANIFEST_FILE")
echo "  ✓ $MANIFEST_COUNT entries in .rein-manifest"

# 8. Configure settings.json
echo "[8/9] Configuring hooks in settings.json..."
SETTINGS_FILE="$PROJECT_DIR/.claude/settings.json"
HOOK_BASE='bash "${CLAUDE_PROJECT_DIR}/.claude/hooks'
if [ -f "$SETTINGS_FILE" ]; then
  echo "  ℹ settings.json exists — merge hooks manually if needed"
  echo "  See hooks/hooks.json for the full configuration template"
else
  mkdir -p "$PROJECT_DIR/.claude"
  cat > "$SETTINGS_FILE" <<SETTINGS
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/session-start.sh\""
          }
        ]
      }
    ],
    "PreToolUse": [
      {
        "matcher": "Edit|Write|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/guard.sh\""
          }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/guard-bash.sh\""
          },
          {
            "type": "command",
            "command": "$HOOK_BASE/gate.sh\""
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/format.sh\""
          },
          {
            "type": "command",
            "command": "$HOOK_BASE/checkbox-guard.sh\""
          },
          {
            "type": "command",
            "command": "$HOOK_BASE/task-progress.sh\""
          }
        ]
      },
      {
        "matcher": "Read|Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/leak-guard.sh\""
          }
        ]
      }
    ],
    "UserPromptExpansion": [
      {
        "matcher": "code-review",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/inject.sh\""
          }
        ]
      }
    ]
  }
}
SETTINGS
  echo "  ✓ settings.json created with all hooks"
fi

# 9. Handle AGENTS.md (Codex CLI compatibility)
echo "[9/9] Checking for Codex CLI..."
AGENTS_MD="$PROJECT_DIR/AGENTS.md"
if [ -f "$AGENTS_MD" ]; then
  echo "  ℹ AGENTS.md found — Codex CLI detected"
  echo "  Append command definitions to AGENTS.md manually if needed"
else
  echo "  ℹ No AGENTS.md found — skipping Codex CLI setup"
fi

echo ""
echo "=== Installation Complete ==="
echo ""
echo "Installed hooks:"
echo "  1. session-start   — Inject using-rein skill (SessionStart)"
echo "  2. guard           — Block edits to rein-managed files (PreToolUse: Edit|Write|MultiEdit)"
echo "  3. guard-bash      — Block destructive cmds on rein files (PreToolUse: Bash)"
echo "  4. gate            — Run tests before deploy (PreToolUse: Bash)"
echo "  5. format          — Auto-format with Prettier (PostToolUse: Write|Edit|MultiEdit)"
echo "  6. checkbox-guard  — Warn when task checkbox not updated (PostToolUse: Write|Edit|MultiEdit)"
echo "  7. leak-guard      — Block secrets in output (PostToolUse: Read|Bash)"
echo "  8. inject          — Inject review checklist (UserPromptExpansion: /code-review)"
echo "  9. task-progress   — Inject task progress after code edits (PostToolUse: Write|Edit|MultiEdit)"
echo ""
echo "Protection:"
echo "  rein-managed files are listed in .claude/.rein-manifest"
echo "  Edit/Write on these files will be blocked automatically"
echo "  To allow edits, remove the file's entry from the manifest"
echo ""
echo "Verification steps:"
echo "1. Start a new Claude Code session"
echo "2. The using-rein skill should be auto-injected"
echo "3. Try /triage to test the workflow"
echo "4. Try /code-review to test checklist injection"
