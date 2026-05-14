# Using rein with Codex CLI

## Setup

### Option 1: Install Script (Recommended)

```bash
git clone https://github.com/zhoupihua/rein.git
cd your-project
bash /path/to/rein/install/install.sh --ide codex
```

On Windows:

```powershell
powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1 -Ide codex
```

### Option 2: Manual Setup

1. Install the rein CLI:
   ```bash
   bash /path/to/rein/install/install.sh --global
   ```

2. Convert artifacts to CODEX.md:
   ```bash
   rein convert --ide codex --source-dir /path/to/rein --output-dir .
   ```

3. Create `.codex/config.toml` with hooks:
   ```toml
   [features]
   multi_agent = true

   [[hooks.pre_command]]
   command = "/path/to/rein hook guard"
   description = "Block edits to rein-managed files"

   [[hooks.pre_command]]
   command = "/path/to/rein hook gate"
   description = "Run tests before deploy commands"

   [[hooks.post_command]]
   command = "/path/to/rein hook format"
   description = "Auto-format web files with prettier"
   ```

## How Codex Maps to rein Artifacts

Codex CLI reads a single `CODEX.md` file for project instructions. All rein artifacts are merged into this file with clear section headings.

| rein Artifact | Codex Equivalent | Location |
|---------------|-----------------|----------|
| Skills (`SKILL.md`) | `## Skills` section in CODEX.md | Single file |
| Commands (`*.md`) | `## Commands` section in CODEX.md | Single file |
| Agents (`*.md`) | `## Agents` section in CODEX.md | Single file |
| CLAUDE.md conventions | `## Task Progress` + `## rein Workflow` in CODEX.md | Single file |

### Referencing Skills and Commands

Unlike Claude Code (slash commands) or Cursor (@rule-name), Codex loads all instructions from `CODEX.md` automatically. To activate a specific skill or agent:

- **Skills**: Describe what you want and Codex will follow the matching skill section
- **Commands**: Reference by name (e.g., "Run the feature workflow" matches the `feature` command section)
- **Agents**: Reference by name (e.g., "Adopt the code-reviewer perspective")

## Hooks

rein configures three hooks in `.codex/config.toml`:

| Hook | rein Command | What it does |
|------|-------------|-------------|
| pre_command | `rein hook guard` | Blocks edits to rein-managed files |
| pre_command | `rein hook gate` | Runs tests before deploy commands |
| post_command | `rein hook format` | Auto-formats web files with prettier |

### Multi-Agent Support

The config enables `multi_agent = true` which allows Codex to spawn sub-agents. This is required for skills like `subagent` that dispatch worker agents. See `references/codex-tools.md` for the full tool mapping including `spawn_agent`, `wait`, and `close_agent`.

### Hook Limitations

These Claude Code hooks have no Codex equivalent:

| Claude Code Hook | Why No Equivalent |
|-----------------|-------------------|
| `session-start` | Codex has no session start event |
| `leak-guard` | No post-read hook |
| `inject` | No prompt expansion hooks |
| `checkbox-guard` | No access to edit content diff |
| `artifact-validate` | No access to edit content for validation |
| `guard-bash` | No pre-command filtering by tool type |

## Usage Tips

1. **CODEX.md is always loaded** — All skills, commands, and agents are available in every session
2. **Use rein CLI for task tracking** — `rein status`, `rein task done 1.2`, `rein validate` work from Codex's terminal
3. **Reference agents by name** — "Use the code-reviewer perspective" activates the agent's framework
4. **Named agent dispatch requires workaround** — See `references/codex-tools.md` for how to dispatch `spawn_agent` with agent persona content
5. **File protection is limited** — `pre_command` hooks run before every command but lack file-path context for precise guarding
