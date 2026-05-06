#!/bin/bash
# rein install script (Linux/Mac)
# Project install: bash /path/to/rein/install/install.sh
# Global install:  bash /path/to/rein/install/install.sh --global

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKFLOW_DIR="$(dirname "$SCRIPT_DIR")"

# Parse arguments
GLOBAL=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --global|-g) GLOBAL=true; shift ;;
        *) shift ;;
    esac
done

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
        (cd "$WORKFLOW_DIR" && go build -o "$bin_dir/rein" ./cmd/rein/)
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

# --- Shared helper: copy resources ---
copy_resources() {
    local target_dir="$1"

    # Skills
    mkdir -p "$target_dir/skills"
    cp -r "$WORKFLOW_DIR/skills/"* "$target_dir/skills/"
    local skill_count=$(ls -d "$target_dir/skills/"*/ 2>/dev/null | wc -l)
    echo "  ✓ $skill_count skills installed"

    # Commands
    mkdir -p "$target_dir/commands"
    cp "$WORKFLOW_DIR/commands/"*.md "$target_dir/commands/"
    local cmd_count=$(ls "$target_dir/commands/"*.md 2>/dev/null | wc -l)
    echo "  ✓ $cmd_count commands installed"

    # Agents
    mkdir -p "$target_dir/agents"
    cp "$WORKFLOW_DIR/agents/"*.md "$target_dir/agents/"
    local agent_count=$(ls "$target_dir/agents/"*.md 2>/dev/null | wc -l)
    echo "  ✓ $agent_count agents installed"

    # Checklists
    mkdir -p "$target_dir/checklists"
    if [ -f "$WORKFLOW_DIR/templates/checklists/review.md" ]; then
        cp "$WORKFLOW_DIR/templates/checklists/review.md" "$target_dir/checklists/"
        echo "  ✓ review.md checklist installed"
    fi
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
    } > "$manifest_file"
    local count=$(grep -cvE '^\s*#|^\s*$' "$manifest_file")
    echo "  ✓ $count entries in .rein-manifest"
}

# --- Shared helper: configure settings.json ---
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
            {"type": "command", "command": f"{hook_cmd} task-progress"}
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
        {"type": "command", "command": "$hook_cmd task-progress"}
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

# ============================================================
# Global Install
# ============================================================
if [ "$GLOBAL" = true ]; then
    CONFIG_DIR="${CLAUDE_CONFIG_DIR:-$HOME/.claude}"
    BIN_DIR="$CONFIG_DIR/bin"

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

    # [4/8] Configure settings.json (hooks use $CLAUDE_CONFIG_DIR)
    echo "[4/8] Configuring settings.json..."
    HOOK_CMD='${CLAUDE_CONFIG_DIR:-$HOME/.claude}/bin/rein hook'
    configure_settings "$CONFIG_DIR/settings.json" "$HOOK_CMD"

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
    mkdir -p "$PROJECT_DIR/docs/rein/specs"
    mkdir -p "$PROJECT_DIR/docs/rein/plans"
    mkdir -p "$PROJECT_DIR/docs/rein/tasks"
    mkdir -p "$PROJECT_DIR/docs/rein/archive"
    echo "  ✓ docs/rein/{specs,plans,tasks,archive}"

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

# ============================================================
# Project Install
# ============================================================
else
    PROJECT_DIR="$(pwd)"
    CLAUDE_DIR="$PROJECT_DIR/.claude"
    CONFIG_DIR="${CLAUDE_CONFIG_DIR:-$HOME/.claude}"

    echo "=== rein Project Installer ==="
    echo "Target: $PROJECT_DIR"
    echo ""

    # Check existing installation
    if [ -f "$CLAUDE_DIR/.rein-manifest" ]; then
        echo "ℹ Existing rein installation detected — upgrading" >&2
    fi

    # [1/6] Install/upgrade rein CLI globally
    echo "[1/6] Installing rein CLI..."
    GLOBAL_BIN_DIR="$CONFIG_DIR/bin"
    install_binary "$GLOBAL_BIN_DIR"

    # [2/6] Copy resources
    echo "[2/6] Installing resources..."
    copy_resources "$CLAUDE_DIR"

    # [3/6] Generate manifest
    echo "[3/6] Generating protection manifest..."
    generate_manifest "$CLAUDE_DIR"

    # [4/6] Configure settings.json (hooks use global binary)
    echo "[4/6] Configuring settings.json..."
    HOOK_CMD='${CLAUDE_CONFIG_DIR:-$HOME/.claude}/bin/rein hook'
    configure_settings "$CLAUDE_DIR/settings.json" "$HOOK_CMD"

    # Clean up old bash/ps1 hooks
    if [ -d "$CLAUDE_DIR/hooks" ]; then
        rm -f "$CLAUDE_DIR/hooks/"*.sh "$CLAUDE_DIR/hooks/"*.ps1 2>/dev/null
        echo "  ✓ Cleaned old hook scripts"
    fi

    # [5/6] Create artifact directories
    echo "[5/6] Creating artifact directories..."
    mkdir -p "$PROJECT_DIR/docs/rein/specs"
    mkdir -p "$PROJECT_DIR/docs/rein/plans"
    mkdir -p "$PROJECT_DIR/docs/rein/tasks"
    mkdir -p "$PROJECT_DIR/docs/rein/archive"
    echo "  ✓ docs/rein/{specs,plans,tasks,archive}"

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
    echo "  Resources: $CLAUDE_DIR/skills/, commands/, agents/"
    echo "  Hooks:     $CLAUDE_DIR/settings.json"
    echo "  Artifacts: $PROJECT_DIR/docs/rein/"
    echo ""
    echo "Run 'rein status' to check current phase"
fi
