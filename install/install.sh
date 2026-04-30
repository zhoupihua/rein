#!/bin/bash
# Alloy install script (Linux/Mac)
# Run from your project root: bash /path/to/Alloy/install/install.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKFLOW_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(pwd)"

echo "=== Alloy Installer ==="
echo "Workflow source: $WORKFLOW_DIR"
echo "Target project:  $PROJECT_DIR"
echo ""

# 1. Create artifact directories
echo "[1/9] Creating artifact directories..."
mkdir -p "$PROJECT_DIR/specs"
mkdir -p "$PROJECT_DIR/changes"
mkdir -p "$PROJECT_DIR/archive"
echo "  ✓ specs/, changes/, archive/"

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
cp "$WORKFLOW_DIR/hooks/session-start.sh" "$PROJECT_DIR/.claude/hooks/"
cp "$WORKFLOW_DIR/hooks/format.sh" "$PROJECT_DIR/.claude/hooks/"
cp "$WORKFLOW_DIR/hooks/test-gateway.sh" "$PROJECT_DIR/.claude/hooks/"
cp "$WORKFLOW_DIR/hooks/secret-scan.sh" "$PROJECT_DIR/.claude/hooks/"
cp "$WORKFLOW_DIR/hooks/context-inject.sh" "$PROJECT_DIR/.claude/hooks/"
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
  cp "$WORKFLOW_DIR/hooks/session-start.ps1" "$PROJECT_DIR/.claude/hooks/"
  cp "$WORKFLOW_DIR/hooks/format.ps1" "$PROJECT_DIR/.claude/hooks/"
  cp "$WORKFLOW_DIR/hooks/test-gateway.ps1" "$PROJECT_DIR/.claude/hooks/"
  cp "$WORKFLOW_DIR/hooks/secret-scan.ps1" "$PROJECT_DIR/.claude/hooks/"
  cp "$WORKFLOW_DIR/hooks/context-inject.ps1" "$PROJECT_DIR/.claude/hooks/"
fi
chmod +x "$PROJECT_DIR/.claude/hooks/"*.sh
echo "  ✓ All hooks installed (session-start, format, test-gateway, secret-scan, context-inject)"

# 6. Copy checklists
echo "[6/9] Installing checklists..."
mkdir -p "$PROJECT_DIR/.claude/checklists"
if [ -f "$WORKFLOW_DIR/templates/checklists/review.md" ]; then
  cp "$WORKFLOW_DIR/templates/checklists/review.md" "$PROJECT_DIR/.claude/checklists/"
  echo "  ✓ review.md checklist installed"
else
  echo "  ⚠ No review checklist template found"
fi

# 7. Configure settings.json
echo "[7/9] Configuring hooks in settings.json..."
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
    "PostToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/format.sh\""
          }
        ]
      },
      {
        "matcher": "Read|Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/secret-scan.sh\""
          }
        ]
      }
    ],
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/test-gateway.sh\""
          }
        ]
      }
    ],
    "UserPromptExpansion": [
      {
        "matcher": "review",
        "hooks": [
          {
            "type": "command",
            "command": "$HOOK_BASE/context-inject.sh\""
          }
        ]
      }
    ]
  }
}
SETTINGS
  echo "  ✓ settings.json created with all hooks"
fi

# 8. Append workflow instructions to CLAUDE.md
echo "[8/9] Updating CLAUDE.md..."
CLAUDE_MD="$PROJECT_DIR/CLAUDE.md"
WORKFLOW_BLOCK=$(cat <<'BLOCK'

## Alloy

This project uses Alloy for structured AI-assisted development.

### Commands
- `/triage` — Classify a change as L1/L2/L3
- `/quick` — L1: ≤5 lines, no logic impact
- `/fix` — L2: 1-3 files, clear requirements
- `/feature` — L3: Full 8-step workflow
- `/spec` — Generate change artifacts
- `/plan` — Task breakdown
- `/build` — Execute tasks from tasks.md
- `/test` — TDD workflow
- `/review` — Five-axis code review
- `/ship` — Fan-out review + GO/NO-GO
- `/simplify` — Code simplification
- `/resume` — Resume from breakpoint

### Artifact Directories
- `specs/` — Published specs (long-lived)
- `changes/` — Active changes (short-lived)
- `archive/` — Archived changes
BLOCK
)

if [ -f "$CLAUDE_MD" ]; then
  if ! grep -q "Alloy" "$CLAUDE_MD"; then
    echo "$WORKFLOW_BLOCK" >> "$CLAUDE_MD"
    echo "  ✓ Workflow instructions appended to CLAUDE.md"
  else
    echo "  ℹ CLAUDE.md already contains Alloy section"
  fi
else
  echo "# CLAUDE.md" > "$CLAUDE_MD"
  echo "$WORKFLOW_BLOCK" >> "$CLAUDE_MD"
  echo "  ✓ CLAUDE.md created with workflow instructions"
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
echo "  1. SessionStart   — Inject using-workflow skill"
echo "  2. Format         — Auto-format with Prettier (PostToolUse: Write|Edit|MultiEdit)"
echo "  3. Test Gateway   — Run tests before deploy (PreToolUse: Bash)"
echo "  4. Secret Scan    — Block secrets in output (PostToolUse: Read|Bash)"
echo "  5. Context Inject — Inject review checklist (UserPromptExpansion: /review)"
echo ""
echo "Verification steps:"
echo "1. Start a new Claude Code session"
echo "2. The using-workflow skill should be auto-injected"
echo "3. Try /triage to test the workflow"
echo "4. Try /review to test checklist injection"
