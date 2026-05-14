#!/bin/bash
# rein install script (Linux/Mac)
# Project install: bash /path/to/rein/install/install.sh
# Global install:  bash /path/to/rein/install/install.sh --global
# Cursor install:  bash /path/to/rein/install/install.sh --ide cursor
# Codex install:   bash /path/to/rein/install/install.sh --ide codex

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKFLOW_DIR="$(dirname "$SCRIPT_DIR")"

# Parse arguments
GLOBAL=false
IDE="claude"
while [[ $# -gt 0 ]]; do
    case $1 in
        --global|-g) GLOBAL=true; shift ;;
        --ide|-i) IDE="$2"; shift 2 ;;
        --ide=*) IDE="${1#*=}"; shift ;;
        *) shift ;;
    esac
done

# Validate IDE
if [[ "$IDE" != "claude" && "$IDE" != "cursor" && "$IDE" != "codex" ]]; then
    echo "Error: unsupported IDE '$IDE'. Use 'claude', 'cursor', or 'codex'." >&2
    exit 1
fi

# --- Shared helper: install binary ---
install_binary() {
    local bin_dir="$1"
    local rein_version="v0.1.0"

    mkdir -p "$bin_dir"

    # Check existing binary
    if [ -f "$bin_dir/rein" ] || [ -f "$bin_dir/rein.exe" ]; then
        echo "  ℹ Existing rein found — upgrading" >&2
    fi

    if [ -f "$WORKFLOW_DIR/cmd/rein/main.go" ]; then
        # Dev mode: build from source
        echo "  Building from source..." >&2
        (cd "$WORKFLOW_DIR" && go build -mod=mod -o "$bin_dir/rein" ./cmd/rein/)
    else
        local platform=$(uname -s | tr '[:upper:]' '[:lower:]')
        local arch=$(uname -m)
        [ "$arch" = "x86_64" ] && arch="amd64"
        [ "$arch" = "aarch64" ] && arch="arm64"
        local binary_name="rein-${platform}-${arch}"
        if [ "$platform" = "mingw64" ] || [ "$platform" = "msys" ]; then
            binary_name="rein-windows-${arch}.exe"
        fi
        curl -sL "https://github.com/zhoupihua/rein/releases/download/${rein_version}/${binary_name}" \
            -o "$bin_dir/rein"
    fi
    chmod +x "$bin_dir/rein" 2>/dev/null || true
    echo "  ✓ rein CLI installed to $bin_dir/rein"
}

# --- Shared helper: copy resources (Claude Code) ---
copy_resources() {
    local target_dir="$1"

    # Skills — clean first to remove deleted skills on upgrade
    rm -rf "$target_dir/skills"
    mkdir -p "$target_dir/skills"
    cp -r "$WORKFLOW_DIR/skills/"* "$target_dir/skills/"
    local skill_count=$(ls -d "$target_dir/skills/"*/ 2>/dev/null | wc -l)
    echo "  ✓ $skill_count skills installed"

    # Commands — clean first
    rm -rf "$target_dir/commands"
    mkdir -p "$target_dir/commands"
    cp "$WORKFLOW_DIR/commands/"*.md "$target_dir/commands/"
    local cmd_count=$(ls "$target_dir/commands/"*.md 2>/dev/null | wc -l)
    echo "  ✓ $cmd_count commands installed"

    # Agents — clean first
    rm -rf "$target_dir/agents"
    mkdir -p "$target_dir/agents"
    cp "$WORKFLOW_DIR/agents/"*.md "$target_dir/agents/"
    local agent_count=$(ls "$target_dir/agents/"*.md 2>/dev/null | wc -l)
    echo "  ✓ $agent_count agents installed"

    # Checklists — clean first
    rm -rf "$target_dir/checklists"
    mkdir -p "$target_dir/checklists"
    if [ -f "$WORKFLOW_DIR/templates/checklists/review.md" ]; then
        cp "$WORKFLOW_DIR/templates/checklists/review.md" "$target_dir/checklists/"
        echo "  ✓ review.md checklist installed"
    fi
}

# --- Shared helper: copy resources (Cursor) ---
copy_resources_cursor() {
    local cursor_dir="$1"

    rm -rf "$cursor_dir/rules"
    mkdir -p "$cursor_dir/rules"

    local rein_bin="${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein"
    if [ ! -f "$rein_bin" ]; then
        rein_bin="${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein.exe"
    fi

    if [ -f "$rein_bin" ]; then
        # Use rein convert for proper frontmatter transformation
        "$rein_bin" convert --ide cursor --source-dir "$WORKFLOW_DIR" --output-dir "$cursor_dir/rules/" 2>&1 | tail -1
    else
        # Fallback: simple copy with frontmatter transformation
        _convert_skills_fallback "$WORKFLOW_DIR/skills/" "$cursor_dir/rules/"
        _convert_commands_fallback "$WORKFLOW_DIR/commands/" "$cursor_dir/rules/"
        _convert_agents_fallback "$WORKFLOW_DIR/agents/" "$cursor_dir/rules/"
    fi

    # Create always-apply rule for project conventions + task progress
    inject_cursor_rules "$cursor_dir/rules"
}

_convert_skills_fallback() {
    local src_dir="$1"
    local out_dir="$2"
    local count=0

    for skill_dir in "$src_dir"*/; do
        [ -d "$skill_dir" ] || continue
        name=$(basename "$skill_dir")
        skill_file="$skill_dir/SKILL.md"
        [ -f "$skill_file" ] || continue

        # Extract description from YAML frontmatter
        desc=$(grep '^description:' "$skill_file" | head -1 | sed 's/^description: *//')
        [ -z "$desc" ] && desc="$name skill"

        # Build .mdc file — skip original frontmatter, output body
        {
            echo "---"
            echo "description: $desc"
            echo "alwaysApply: false"
            echo "---"
            sed '1{/^---$/d}; /^---$/,/^---$/d; //' "$skill_file"
        } > "$out_dir/$name.mdc"
        count=$((count + 1))
    done
    echo "  ✓ $count skills converted to .mdc rules"
}

_convert_commands_fallback() {
    local src_dir="$1"
    local out_dir="$2"
    local count=0

    for cmd_file in "$src_dir"*.md; do
        [ -f "$cmd_file" ] || continue
        name=$(basename "$cmd_file" .md)

        desc=$(grep '^description:' "$cmd_file" | head -1 | sed 's/^description: *//')
        [ -z "$desc" ] && desc="$name command"

        {
            echo "---"
            echo "description: $desc"
            echo "alwaysApply: false"
            echo "---"
            sed '1{/^---$/d}; /^---$/,/^---$/d; //' "$cmd_file"
        } > "$out_dir/$name.mdc"
        count=$((count + 1))
    done
    echo "  ✓ $count commands converted to .mdc rules"
}

_convert_agents_fallback() {
    local src_dir="$1"
    local out_dir="$2"
    local count=0

    for agent_file in "$src_dir"*.md; do
        [ -f "$agent_file" ] || continue
        name=$(basename "$agent_file" .md)

        desc=$(grep '^description:' "$agent_file" | head -1 | sed 's/^description: *//')
        [ -z "$desc" ] && desc="$name agent"

        {
            echo "---"
            echo "description: $desc"
            echo "alwaysApply: false"
            echo "---"
            sed '1{/^---$/d}; /^---$/,/^---$/d; //' "$agent_file"
        } > "$out_dir/$name.mdc"
        count=$((count + 1))
    done
    echo "  ✓ $count agents converted to .mdc rules"
}

# --- Shared helper: generate manifest ---
generate_manifest() {
    local base_dir="$1"
    local manifest_file="$base_dir/.rein-manifest"

    {
        echo "# rein Managed Files - DO NOT EDIT"
        echo "# To allow edits, remove the file's line from this manifest."
        echo ""
        for dir in bin commands agents checklists; do
            if [ -d "$base_dir/$dir" ]; then
                find "$base_dir/$dir" -maxdepth 1 -type f | while read -r f; do
                    echo "${f#$base_dir/}"
                done
            fi
        done
        if [ -d "$base_dir/skills" ]; then
            find "$base_dir/skills" -mindepth 1 -maxdepth 1 -type d | while read -r d; do
                echo "${d#$base_dir/}/"
            done
        fi
        # For Cursor installs, also list .mdc rule files
        if [ -d "$base_dir/rules" ]; then
            find "$base_dir/rules" -maxdepth 1 -type f -name "*.mdc" | while read -r f; do
                echo "${f#$base_dir/}"
            done
        fi
    } > "$manifest_file"
    local count=$(grep -cvE '^\s*#|^\s*$' "$manifest_file")
    echo "  ✓ $count entries in .rein-manifest"
}

# --- Shared helper: configure settings.json (Claude Code) ---
configure_settings() {
    local settings_file="$1"
    local hook_cmd="$2"

    if [ -f "$settings_file" ]; then
        if command -v python3 &>/dev/null; then
            python3 - "$settings_file" "$hook_cmd" <<'PYEOF'
import json, sys
sf, hook_cmd = sys.argv[1], sys.argv[2]
with open(sf) as f: data = json.load(f)

rein_hooks = {
    "SessionStart": [{"matcher": "", "hooks": [{"type": "command", "command": f"{hook_cmd} session-start"}]}],
    "PreToolUse": [
        {"matcher": "Edit|Write|MultiEdit", "hooks": [{"type": "command", "command": f"{hook_cmd} guard"}]},
        {"matcher": "Bash", "hooks": [
            {"type": "command", "command": f"{hook_cmd} guard-bash"},
            {"type": "command", "command": f"{hook_cmd} gate"}
        ]}
    ],
    "PostToolUse": [
        {"matcher": "Write|Edit|MultiEdit", "hooks": [
            {"type": "command", "command": f"{hook_cmd} format"},
            {"type": "command", "command": f"{hook_cmd} checkbox-guard"},
            {"type": "command", "command": f"{hook_cmd} artifact-validate"}
        ]},
        {"matcher": "Read|Bash", "hooks": [{"type": "command", "command": f"{hook_cmd} leak-guard"}]}
    ],
    "UserPromptExpansion": [{"matcher": "code-review", "hooks": [{"type": "command", "command": f"{hook_cmd} inject"}]}]
}

existing = data.get("hooks", {})
for event, entries in rein_hooks.items():
    if event not in existing:
        existing[event] = entries
    else:
        non_rein = [e for e in existing[event] if not any("rein" in h.get("command","") and "hook" in h.get("command","") for h in e.get("hooks",[]))]
        existing[event] = non_rein + entries
data["hooks"] = existing

perms = data.get("permissions", {})
allow = perms.get("allow", [])
if "Bash(rein *)" not in allow:
    allow.append("Bash(rein *)")
perms["allow"] = allow
data["permissions"] = perms

with open(sf, 'w') as f: json.dump(data, f, indent=2)
PYEOF
            echo "  ✓ settings.json merged (hooks + permissions)"
        else
            echo "  ⚠ python3 not found — merge settings.json manually"
            echo "  Add hooks pointing to: $hook_cmd <name>"
            echo "  Add \"Bash(rein *)\" to permissions.allow"
        fi
    else
        cat > "$settings_file" <<SETTINGSJSON
{
  "hooks": {
    "SessionStart": [{"matcher": "", "hooks": [{"type": "command", "command": "$hook_cmd session-start"}]}],
    "PreToolUse": [
      {"matcher": "Edit|Write|MultiEdit", "hooks": [{"type": "command", "command": "$hook_cmd guard"}]},
      {"matcher": "Bash", "hooks": [
        {"type": "command", "command": "$hook_cmd guard-bash"},
        {"type": "command", "command": "$hook_cmd gate"}
      ]}
    ],
    "PostToolUse": [
      {"matcher": "Write|Edit|MultiEdit", "hooks": [
        {"type": "command", "command": "$hook_cmd format"},
        {"type": "command", "command": "$hook_cmd checkbox-guard"},
        {"type": "command", "command": "$hook_cmd artifact-validate"}
      ]},
      {"matcher": "Read|Bash", "hooks": [{"type": "command", "command": "$hook_cmd leak-guard"}]}
    ],
    "UserPromptExpansion": [{"matcher": "code-review", "hooks": [{"type": "command", "command": "$hook_cmd inject"}]}]
  },
  "permissions": {
    "allow": ["Bash(rein *)"]
  }
}
SETTINGSJSON
        echo "  ✓ settings.json created with hooks + permissions"
    fi
}

# --- Shared helper: configure hooks.json (Cursor) ---
configure_cursor_hooks() {
    local hooks_file="$1"
    local rein_cmd="$2"

    local hooks_content
    hooks_content=$(cat <<HOOKJSON
{
  "hooks": {
    "PreEdit": [
      {"command": "$rein_cmd", "args": ["hook", "guard"]}
    ],
    "PostEdit": [
      {"command": "$rein_cmd", "args": ["hook", "format"]}
    ],
    "PreCommit": [
      {"command": "$rein_cmd", "args": ["hook", "gate"]}
    ]
  }
}
HOOKJSON
)

    if [ -f "$hooks_file" ]; then
        if command -v python3 &>/dev/null; then
            python3 - "$hooks_file" "$hooks_content" <<'PYEOF'
import json, sys
sf, new_json = sys.argv[1], sys.argv[2]
with open(sf) as f: data = json.load(f)
new_hooks = json.loads(new_json)
if "hooks" not in data: data["hooks"] = {}
for event, entries in new_hooks["hooks"].items():
    existing = data["hooks"].get(event, [])
    non_rein = [e for e in existing if "rein" not in e.get("command", "")]
    data["hooks"][event] = non_rein + entries
with open(sf, 'w') as f: json.dump(data, f, indent=2)
PYEOF
            echo "  ✓ hooks.json merged"
        else
            echo "  ⚠ python3 not found — merge hooks.json manually"
        fi
    else
        echo "$hooks_content" > "$hooks_file"
        echo "  ✓ hooks.json created"
    fi
}

# --- Shared helper: inject task progress rule into CLAUDE.md ---
inject_claude_md() {
    local project_dir="$1"
    local claude_md="$project_dir/CLAUDE.md"

    local marker="<!-- rein:task-progress -->"
    local block="$marker
## Task Progress

When working on a feature with \`docs/rein/changes/<name>/task.md\`, after completing
any task or sub-task, you MUST immediately mark it as done:

  rein task done <id>          # e.g., rein task done 1.2
  rein task done <subtask-id>  # e.g., rein task done 1.2.0

Do NOT skip this step. Marking progress is mandatory, not optional.
$marker"

    if [ -f "$claude_md" ]; then
        if grep -qF "$marker" "$claude_md" 2>/dev/null; then
            echo "  ✓ CLAUDE.md task-progress rule already present"
        else
            echo "" >> "$claude_md"
            echo "$block" >> "$claude_md"
            echo "  ✓ CLAUDE.md injected task-progress rule"
        fi
    else
        echo "$block" > "$claude_md"
        echo "  ✓ CLAUDE.md created with task-progress rule"
    fi
}

# --- Shared helper: inject Cursor always-apply rule ---
inject_cursor_rules() {
    local rules_dir="$1"
    local rule_file="$rules_dir/rein-project.mdc"

    cat > "$rule_file" <<'RULE'
---
description: rein project conventions and task progress tracking
alwaysApply: true
---

## Task Progress

When working on a feature with `docs/rein/changes/<name>/task.md`, after completing
any task or sub-task, you MUST immediately mark it as done:

  rein task done <id>          # e.g., rein task done 1.2
  rein task done <subtask-id>  # e.g., rein task done 1.2.0

Do NOT skip this step. Marking progress is mandatory, not optional.

## rein Workflow

This project uses rein for structured development. Key commands:
- `rein status` — Check current workflow phase
- `rein task next` — Show next unchecked task
- `rein task done <id>` — Mark task complete
- `rein validate <feature>` — Validate artifact completeness

## Available Rules

Reference these rules with @<name> in Cursor chat:

| Rule | Type | When to use |
|------|------|-------------|
| @define | skill | Starting a new project or feature |
| @tdd | skill | Implementing logic or fixing bugs |
| @code-review | skill | Reviewing code before merge |
| @planning | skill | Breaking work into tasks |
| @security | skill | Security reviews |
| @performance | skill | Performance work |
| @feature | command | Multi-file feature workflow |
| @fix | command | Bug fix workflow |
| @ship | command | Pre-launch checklist |
| @code-reviewer | agent | Structured code review |
| @test-engineer | agent | Test strategy and coverage |
| @security-auditor | agent | Security analysis |
RULE
    echo "  ✓ rein-project.mdc (alwaysApply) created"
}

# --- Shared helper: copy resources (Codex) ---
copy_resources_codex() {
    local project_dir="$1"
    local codex_dir="$project_dir/.codex"

    mkdir -p "$codex_dir"

    local rein_bin="${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein"
    if [ ! -f "$rein_bin" ]; then
        rein_bin="${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein.exe"
    fi

    if [ -f "$rein_bin" ]; then
        "$rein_bin" convert --ide codex --source-dir "$WORKFLOW_DIR" --output-dir "$project_dir" 2>&1 | tail -1
    else
        # Fallback: generate CODEX.md manually
        _generate_codex_md_fallback "$project_dir"
        _generate_codex_config_fallback "$codex_dir"
    fi
}

_generate_codex_md_fallback() {
    local project_dir="$1"
    local codex_md="$project_dir/CODEX.md"

    {
        echo "# rein Project Instructions"
        echo ""
        echo "## Task Progress"
        echo ""
        echo "When working on a feature with \`docs/rein/changes/<name>/task.md\`, after completing"
        echo "any task or sub-task, you MUST immediately mark it as done:"
        echo ""
        echo "  rein task done <id>          # e.g., rein task done 1.2"
        echo "  rein task done <subtask-id>  # e.g., rein task done 1.2.0"
        echo ""
        echo "Do NOT skip this step. Marking progress is mandatory, not optional."
        echo ""
        echo "## rein Workflow"
        echo ""
        echo "This project uses rein for structured development. Key commands:"
        echo "- \`rein status\` — Check current workflow phase"
        echo "- \`rein task next\` — Show next unchecked task"
        echo "- \`rein task done <id>\` — Mark task complete"
        echo "- \`rein validate <feature>\` — Validate artifact completeness"

        # Append skills
        echo ""
        echo "## Skills"
        for skill_dir in "$WORKFLOW_DIR/skills/"*/; do
            [ -d "$skill_dir" ] || continue
            name=$(basename "$skill_dir")
            skill_file="$skill_dir/SKILL.md"
            [ -f "$skill_file" ] || continue
            echo ""
            echo "### $name"
            echo ""
            # Strip frontmatter
            sed '1{/^---$/d}; /^---$/,/^---$/d; //' "$skill_file"
        done

        # Append commands
        echo ""
        echo "## Commands"
        echo ""
        echo "Reference these commands by name when asking Codex to perform a workflow."
        for cmd_file in "$WORKFLOW_DIR/commands/"*.md; do
            [ -f "$cmd_file" ] || continue
            name=$(basename "$cmd_file" .md)
            echo ""
            echo "### $name"
            echo ""
            sed '1{/^---$/d}; /^---$/,/^---$/d; //' "$cmd_file"
        done

        # Append agents
        echo ""
        echo "## Agents"
        echo ""
        echo "Reference these agents by name when asking Codex to adopt a perspective."
        for agent_file in "$WORKFLOW_DIR/agents/"*.md; do
            [ -f "$agent_file" ] || continue
            name=$(basename "$agent_file" .md)
            echo ""
            echo "### $name"
            echo ""
            sed '1{/^---$/d}; /^---$/,/^---$/d; //' "$agent_file"
        done
    } > "$codex_md"
    echo "  ✓ CODEX.md created"
}

_generate_codex_config_fallback() {
    local codex_dir="$1"
    local config_file="$codex_dir/config.toml"
    local rein_cmd="${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein"

    cat > "$config_file" <<TOML
# rein Codex configuration
[features]
multi_agent = true

[[hooks.pre_command]]
command = "$rein_cmd hook guard"
description = "Block edits to rein-managed files"

[[hooks.pre_command]]
command = "$rein_cmd hook gate"
description = "Run tests before deploy commands"

[[hooks.post_command]]
command = "$rein_cmd hook format"
description = "Auto-format web files with prettier"
TOML
    echo "  ✓ .codex/config.toml created"
}

# --- Shared helper: configure Codex config.toml ---
configure_codex_config() {
    local config_file="$1"
    local rein_cmd="$2"
    local codex_dir
    codex_dir=$(dirname "$config_file")
    mkdir -p "$codex_dir"

    local config_content
    config_content=$(cat <<TOML
# rein Codex configuration
[features]
multi_agent = true

[[hooks.pre_command]]
command = "$rein_cmd hook guard"
description = "Block edits to rein-managed files"

[[hooks.pre_command]]
command = "$rein_cmd hook gate"
description = "Run tests before deploy commands"

[[hooks.post_command]]
command = "$rein_cmd hook format"
description = "Auto-format web files with prettier"
TOML
)

    if [ -f "$config_file" ]; then
        # Simple merge: append rein hooks if not present
        if grep -q "rein hook guard" "$config_file" 2>/dev/null; then
            echo "  ✓ config.toml already has rein hooks"
        else
            echo "" >> "$config_file"
            echo "$config_content" >> "$config_file"
            echo "  ✓ config.toml appended with rein hooks"
        fi
    else
        echo "$config_content" > "$config_file"
        echo "  ✓ config.toml created"
    fi
}

# ============================================================
# Global Install
# ============================================================
if [ "$GLOBAL" = true ]; then
    if [ "$IDE" = "cursor" ]; then
        CONFIG_DIR="${REIN_CONFIG_DIR:-$HOME/.rein}"
        BIN_DIR="$CONFIG_DIR/bin"
        CURSOR_GLOBAL_DIR="$HOME/.cursor"

        echo "=== rein Global Installer (Cursor) ==="
        echo "Target: $CURSOR_GLOBAL_DIR"
        echo ""

        # [1/6] Install binary
        echo "[1/6] Installing rein CLI..."
        install_binary "$BIN_DIR"

        # [2/6] Copy resources as .mdc rules
        echo "[2/6] Installing Cursor rules..."
        copy_resources_cursor "$CURSOR_GLOBAL_DIR"

        # [3/6] Generate manifest
        echo "[3/6] Generating protection manifest..."
        generate_manifest "$CURSOR_GLOBAL_DIR"

        # [4/6] Configure hooks.json
        echo "[4/6] Configuring hooks.json..."
        HOOK_CMD="$BIN_DIR/rein"
        configure_cursor_hooks "$CURSOR_GLOBAL_DIR/hooks.json" "$HOOK_CMD"

        # [5/6] Add to PATH
        echo "[5/6] Adding to PATH..."
        PATH_LINE="export PATH=\"$BIN_DIR:\$PATH\""
        SHELL_RC=""
        if [ -n "$ZSH_VERSION" ]; then
            SHELL_RC="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            SHELL_RC="$HOME/.bashrc"
        else
            SHELL_RC="$HOME/.profile"
        fi

        if [ -f "$SHELL_RC" ] && ! grep -qF "$BIN_DIR" "$SHELL_RC" 2>/dev/null; then
            echo "" >> "$SHELL_RC"
            echo "# rein CLI" >> "$SHELL_RC"
            echo "$PATH_LINE" >> "$SHELL_RC"
            echo "  ✓ Added to $SHELL_RC"
            echo "  ℹ Run 'source $SHELL_RC' or open a new terminal"
        else
            echo "  ✓ PATH already configured"
        fi

        # [6/6] Create artifact directories
        echo "[6/6] Creating artifact directories..."
        PROJECT_DIR="$(pwd)"
        mkdir -p "$PROJECT_DIR/docs/rein/changes"
        mkdir -p "$PROJECT_DIR/docs/rein/archive"
        echo "  ✓ docs/rein/{changes,archive}"

        echo ""
        echo "=== Global Installation Complete (Cursor) ==="
        echo ""
        echo "Installed:"
        echo "  Binary:    $BIN_DIR/rein"
        echo "  Rules:     $CURSOR_GLOBAL_DIR/rules/"
        echo "  Hooks:     $CURSOR_GLOBAL_DIR/hooks.json"
        echo "  Artifacts: $PROJECT_DIR/docs/rein/"
        echo ""
        echo "Reference rules in Cursor chat with @<rule-name>"

    elif [ "$IDE" = "codex" ]; then
        CONFIG_DIR="${REIN_CONFIG_DIR:-$HOME/.rein}"
        BIN_DIR="$CONFIG_DIR/bin"
        PROJECT_DIR="$(pwd)"
        CODEX_DIR="$PROJECT_DIR/.codex"

        echo "=== rein Global Installer (Codex) ==="
        echo "Target: $PROJECT_DIR"
        echo ""

        # [1/6] Install binary
        echo "[1/6] Installing rein CLI..."
        install_binary "$BIN_DIR"

        # [2/6] Copy resources as CODEX.md
        echo "[2/6] Generating CODEX.md and config..."
        copy_resources_codex "$PROJECT_DIR"

        # [3/6] Generate manifest
        echo "[3/6] Generating protection manifest..."
        generate_manifest "$CODEX_DIR"

        # [4/6] Configure config.toml
        echo "[4/6] Configuring .codex/config.toml..."
        HOOK_CMD="$BIN_DIR/rein"
        configure_codex_config "$CODEX_DIR/config.toml" "$HOOK_CMD"

        # [5/6] Add to PATH
        echo "[5/6] Adding to PATH..."
        PATH_LINE="export PATH=\"$BIN_DIR:\$PATH\""
        SHELL_RC=""
        if [ -n "$ZSH_VERSION" ]; then
            SHELL_RC="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            SHELL_RC="$HOME/.bashrc"
        else
            SHELL_RC="$HOME/.profile"
        fi

        if [ -f "$SHELL_RC" ] && ! grep -qF "$BIN_DIR" "$SHELL_RC" 2>/dev/null; then
            echo "" >> "$SHELL_RC"
            echo "# rein CLI" >> "$SHELL_RC"
            echo "$PATH_LINE" >> "$SHELL_RC"
            echo "  ✓ Added to $SHELL_RC"
            echo "  ℹ Run 'source $SHELL_RC' or open a new terminal"
        else
            echo "  ✓ PATH already configured"
        fi

        # [6/6] Create artifact directories
        echo "[6/6] Creating artifact directories..."
        mkdir -p "$PROJECT_DIR/docs/rein/changes"
        mkdir -p "$PROJECT_DIR/docs/rein/archive"
        echo "  ✓ docs/rein/{changes,archive}"

        echo ""
        echo "=== Global Installation Complete (Codex) ==="
        echo ""
        echo "Installed:"
        echo "  Binary:    $BIN_DIR/rein"
        echo "  Rules:     $PROJECT_DIR/CODEX.md"
        echo "  Config:    $CODEX_DIR/config.toml"
        echo "  Artifacts: $PROJECT_DIR/docs/rein/"
        echo ""
        echo "Codex will read CODEX.md automatically on session start"

    else
        CONFIG_DIR="${REIN_CONFIG_DIR:-$HOME/.rein}"
        BIN_DIR="$CONFIG_DIR/bin"
        CLAUDE_SETTINGS_DIR="$HOME/.claude"

        echo "=== rein Global Installer ==="
        echo "Target: $CONFIG_DIR"
        echo ""

        # Check existing installation
        if [ -f "$CONFIG_DIR/.rein-manifest" ]; then
            echo "ℹ Existing rein installation detected — upgrading" >&2
        fi

        # [1/8] Install binary
        echo "[1/8] Installing rein CLI..."
        install_binary "$BIN_DIR"

        # [2/8] Copy resources
        echo "[2/8] Installing resources..."
        copy_resources "$CONFIG_DIR"

        # [3/8] Generate manifest
        echo "[3/8] Generating protection manifest..."
        generate_manifest "$CONFIG_DIR"

        # [4/8] Configure settings.json (hooks use $REIN_CONFIG_DIR)
        echo "[4/8] Configuring settings.json..."
        HOOK_CMD='${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein hook'
        configure_settings "$CLAUDE_SETTINGS_DIR/settings.json" "$HOOK_CMD"

        # Clean up old bash/ps1 hooks
        if [ -d "$CONFIG_DIR/hooks" ]; then
            rm -f "$CONFIG_DIR/hooks/"*.sh "$CONFIG_DIR/hooks/"*.ps1 2>/dev/null
            echo "  ✓ Cleaned old hook scripts"
        fi

        # [5/8] Add to PATH
        echo "[5/8] Adding to PATH..."
        PATH_LINE="export PATH=\"$BIN_DIR:\$PATH\""
        SHELL_RC=""
        if [ -n "$ZSH_VERSION" ]; then
            SHELL_RC="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            SHELL_RC="$HOME/.bashrc"
        else
            SHELL_RC="$HOME/.profile"
        fi

        if [ -f "$SHELL_RC" ] && ! grep -qF "$BIN_DIR" "$SHELL_RC" 2>/dev/null; then
            echo "" >> "$SHELL_RC"
            echo "# rein CLI" >> "$SHELL_RC"
            echo "$PATH_LINE" >> "$SHELL_RC"
            echo "  ✓ Added to $SHELL_RC"
            echo "  ℹ Run 'source $SHELL_RC' or open a new terminal"
        else
            echo "  ✓ PATH already configured"
        fi

        # [6/8] Create artifact directories
        echo "[6/8] Creating artifact directories..."
        PROJECT_DIR="$(pwd)"
        mkdir -p "$PROJECT_DIR/docs/rein/changes"
        mkdir -p "$PROJECT_DIR/docs/rein/archive"
        echo "  ✓ docs/rein/{changes,archive}"

        # [7/8] Inject task progress rule into CLAUDE.md
        echo "[7/8] Injecting task-progress rule into CLAUDE.md..."
        inject_claude_md "$PROJECT_DIR"

        echo ""
        echo "=== Global Installation Complete ==="
        echo ""
        echo "Installed:"
        echo "  Binary:    $BIN_DIR/rein"
        echo "  Resources: $CONFIG_DIR/skills/, commands/, agents/"
        echo "  Hooks:     $CONFIG_DIR/settings.json"
        echo "  Artifacts: $PROJECT_DIR/docs/rein/"
        echo ""
        echo "Run 'rein status' to check current phase"
    fi

# ============================================================
# Project Install
# ============================================================
else
    PROJECT_DIR="$(pwd)"

    if [ "$IDE" = "cursor" ]; then
        CURSOR_DIR="$PROJECT_DIR/.cursor"
        CONFIG_DIR="${REIN_CONFIG_DIR:-$HOME/.rein}"

        echo "=== rein Project Installer (Cursor) ==="
        echo "Target: $PROJECT_DIR"
        echo ""

        # Check existing installation
        if [ -f "$CURSOR_DIR/.rein-manifest" ]; then
            echo "ℹ Existing rein installation detected — upgrading" >&2
        fi

        # [1/5] Install/upgrade rein CLI globally
        echo "[1/5] Installing rein CLI..."
        GLOBAL_BIN_DIR="$CONFIG_DIR/bin"
        install_binary "$GLOBAL_BIN_DIR"

        # [2/5] Copy resources as .mdc rules
        echo "[2/5] Installing Cursor rules..."
        copy_resources_cursor "$CURSOR_DIR"

        # [3/5] Generate manifest
        echo "[3/5] Generating protection manifest..."
        generate_manifest "$CURSOR_DIR"

        # [4/5] Configure hooks.json
        echo "[4/5] Configuring hooks.json..."
        HOOK_CMD="$GLOBAL_BIN_DIR/rein"
        configure_cursor_hooks "$CURSOR_DIR/hooks.json" "$HOOK_CMD"

        # [5/5] Create artifact directories
        echo "[5/5] Creating artifact directories..."
        mkdir -p "$PROJECT_DIR/docs/rein/changes"
        mkdir -p "$PROJECT_DIR/docs/rein/archive"
        echo "  ✓ docs/rein/{changes,archive}"

        echo ""
        echo "=== Project Installation Complete (Cursor) ==="
        echo ""
        echo "Installed:"
        echo "  Rules:     $CURSOR_DIR/rules/"
        echo "  Hooks:     $CURSOR_DIR/hooks.json"
        echo "  Artifacts: $PROJECT_DIR/docs/rein/"
        echo ""
        echo "Reference rules in Cursor chat with @<rule-name>"

    elif [ "$IDE" = "codex" ]; then
        CODEX_DIR="$PROJECT_DIR/.codex"
        CONFIG_DIR="${REIN_CONFIG_DIR:-$HOME/.rein}"

        echo "=== rein Project Installer (Codex) ==="
        echo "Target: $PROJECT_DIR"
        echo ""

        # Check existing installation
        if [ -f "$CODEX_DIR/.rein-manifest" ]; then
            echo "ℹ Existing rein installation detected — upgrading" >&2
        fi

        # [1/5] Install/upgrade rein CLI globally
        echo "[1/5] Installing rein CLI..."
        GLOBAL_BIN_DIR="$CONFIG_DIR/bin"
        install_binary "$GLOBAL_BIN_DIR"

        # [2/5] Copy resources as CODEX.md
        echo "[2/5] Generating CODEX.md and config..."
        copy_resources_codex "$PROJECT_DIR"

        # [3/5] Generate manifest
        echo "[3/5] Generating protection manifest..."
        generate_manifest "$CODEX_DIR"

        # [4/5] Configure config.toml
        echo "[4/5] Configuring .codex/config.toml..."
        HOOK_CMD="$GLOBAL_BIN_DIR/rein"
        configure_codex_config "$CODEX_DIR/config.toml" "$HOOK_CMD"

        # [5/5] Create artifact directories
        echo "[5/5] Creating artifact directories..."
        mkdir -p "$PROJECT_DIR/docs/rein/changes"
        mkdir -p "$PROJECT_DIR/docs/rein/archive"
        echo "  ✓ docs/rein/{changes,archive}"

        echo ""
        echo "=== Project Installation Complete (Codex) ==="
        echo ""
        echo "Installed:"
        echo "  Rules:     $PROJECT_DIR/CODEX.md"
        echo "  Config:    $CODEX_DIR/config.toml"
        echo "  Artifacts: $PROJECT_DIR/docs/rein/"
        echo ""
        echo "Codex will read CODEX.md automatically on session start"

    else
        REIN_DIR="$PROJECT_DIR/.rein"
        CONFIG_DIR="${REIN_CONFIG_DIR:-$HOME/.rein}"
        CLAUDE_DIR="$PROJECT_DIR/.claude"

        echo "=== rein Project Installer ==="
        echo "Target: $PROJECT_DIR"
        echo ""

        # Check existing installation
        if [ -f "$REIN_DIR/.rein-manifest" ]; then
            echo "ℹ Existing rein installation detected — upgrading" >&2
        fi

        # [1/6] Install/upgrade rein CLI globally
        echo "[1/6] Installing rein CLI..."
        GLOBAL_BIN_DIR="$CONFIG_DIR/bin"
        install_binary "$GLOBAL_BIN_DIR"

        # [2/6] Copy resources
        echo "[2/6] Installing resources..."
        copy_resources "$REIN_DIR"

        # [3/6] Generate manifest
        echo "[3/6] Generating protection manifest..."
        generate_manifest "$REIN_DIR"

        # [4/6] Configure settings.json (hooks use global binary)
        echo "[4/6] Configuring settings.json..."
        HOOK_CMD='${REIN_CONFIG_DIR:-$HOME/.rein}/bin/rein hook'
        configure_settings "$CLAUDE_DIR/settings.json" "$HOOK_CMD"

        # Clean up old bash/ps1 hooks
        if [ -d "$REIN_DIR/hooks" ]; then
            rm -f "$REIN_DIR/hooks/"*.sh "$REIN_DIR/hooks/"*.ps1 2>/dev/null
            echo "  ✓ Cleaned old hook scripts"
        fi

        # [5/6] Create artifact directories
        echo "[5/6] Creating artifact directories..."
        mkdir -p "$PROJECT_DIR/docs/rein/changes"
        mkdir -p "$PROJECT_DIR/docs/rein/archive"
        echo "  ✓ docs/rein/{changes,archive}"

        # [5.5/6] Inject task progress rule into CLAUDE.md
        inject_claude_md "$PROJECT_DIR"

        # [6/6] Verification
        echo "[6/6] Verifying installation..."
        if [ -f "$GLOBAL_BIN_DIR/rein" ] || [ -f "$GLOBAL_BIN_DIR/rein.exe" ]; then
            echo "  ✓ rein CLI available at $GLOBAL_BIN_DIR/"
        else
            echo "  ⚠ rein CLI not found — installation may have failed"
        fi

        echo ""
        echo "=== Project Installation Complete ==="
        echo ""
        echo "Installed:"
        echo "  Resources: $REIN_DIR/skills/, commands/, agents/"
        echo "  Hooks:     $CLAUDE_DIR/settings.json"
        echo "  Artifacts: $PROJECT_DIR/docs/rein/"
        echo ""
        echo "Run 'rein status' to check current phase"
    fi
fi
