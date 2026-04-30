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
echo "[1/8] Creating artifact directories..."
mkdir -p "$PROJECT_DIR/specs"
mkdir -p "$PROJECT_DIR/changes"
mkdir -p "$PROJECT_DIR/archive"
echo "  ✓ specs/, changes/, archive/"

# 2. Copy commands
echo "[2/8] Installing commands..."
mkdir -p "$PROJECT_DIR/.claude/commands"
cp "$WORKFLOW_DIR/commands/"*.md "$PROJECT_DIR/.claude/commands/"
CMD_COUNT=$(ls "$PROJECT_DIR/.claude/commands/"*.md | wc -l)
echo "  ✓ $CMD_COUNT commands installed"

# 3. Copy skills
echo "[3/8] Installing skills..."
mkdir -p "$PROJECT_DIR/.claude/skills"
cp -r "$WORKFLOW_DIR/skills/"* "$PROJECT_DIR/.claude/skills/"
SKILL_COUNT=$(ls -d "$PROJECT_DIR/.claude/skills/"*/ | wc -l)
echo "  ✓ $SKILL_COUNT skills installed"

# 4. Copy agents
echo "[4/8] Installing agents..."
mkdir -p "$PROJECT_DIR/.claude/agents"
cp "$WORKFLOW_DIR/agents/"*.md "$PROJECT_DIR/.claude/agents/"
AGENT_COUNT=$(ls "$PROJECT_DIR/.claude/agents/"*.md | wc -l)
echo "  ✓ $AGENT_COUNT agents installed"

# 5. Copy hooks
echo "[5/8] Installing hooks..."
mkdir -p "$PROJECT_DIR/.claude/hooks"
cp "$WORKFLOW_DIR/hooks/session-start.sh" "$PROJECT_DIR/.claude/hooks/"
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
  cp "$WORKFLOW_DIR/hooks/session-start.ps1" "$PROJECT_DIR/.claude/hooks/"
fi
chmod +x "$PROJECT_DIR/.claude/hooks/session-start.sh"
echo "  ✓ Hooks installed"

# 6. Configure settings.json
echo "[6/8] Configuring hooks in settings.json..."
SETTINGS_FILE="$PROJECT_DIR/.claude/settings.json"
if [ -f "$SETTINGS_FILE" ]; then
  echo "  ℹ settings.json exists — merge hooks manually if needed"
  echo "  Add this to your hooks config:"
  echo '  {"hooks": {"SessionStart": [{"matcher": "", "hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR}/.claude/hooks/session-start.sh\""}]}]}}'
else
  mkdir -p "$PROJECT_DIR/.claude"
  cat > "$SETTINGS_FILE" <<'SETTINGS'
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR}/.claude/hooks/session-start.sh\""
          }
        ]
      }
    ]
  }
}
SETTINGS
  echo "  ✓ settings.json created with session-start hook"
fi

# 7. Append workflow instructions to CLAUDE.md
echo "[7/8] Updating CLAUDE.md..."
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

# 8. Handle AGENTS.md (Codex CLI compatibility)
echo "[8/8] Checking for Codex CLI..."
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
echo "Verification steps:"
echo "1. Start a new Claude Code session"
echo "2. The using-workflow skill should be auto-injected"
echo "3. Try /triage to test the workflow"
echo "4. Try /spec test-feature to test artifact generation"
echo ""
echo "Artifact directories created: specs/, changes/, archive/"
echo "Next: Start a new Claude Code session and try /triage"