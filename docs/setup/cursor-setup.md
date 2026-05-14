# Using rein with Cursor

## Setup

### Option 1: Install Script (Recommended)

```bash
git clone https://github.com/zhoupihua/rein.git
cd your-project
bash /path/to/rein/install/install.sh --ide cursor
```

On Windows:

```powershell
powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1 -Ide cursor
```

### Option 2: Manual Setup

1. Install the rein CLI:
   ```bash
   bash /path/to/rein/install/install.sh --global
   ```

2. Convert skills, commands, and agents to `.mdc` rules:
   ```bash
   rein convert --ide cursor --source-dir /path/to/rein --output-dir .cursor/rules/
   ```

3. Create `.cursor/hooks.json`:
   ```json
   {
     "hooks": {
       "PreEdit": [{"command": "/path/to/rein", "args": ["hook", "guard"]}],
       "PostEdit": [{"command": "/path/to/rein", "args": ["hook", "format"]}],
       "PreCommit": [{"command": "/path/to/rein", "args": ["hook", "gate"]}]
     }
   }
   ```

## How Cursor Rules Map to rein Artifacts

| rein Artifact | Cursor Rule | Application |
|---------------|-------------|-------------|
| Skills (`SKILL.md`) | `.cursor/rules/<skill-name>.mdc` | Agent-requested (AI loads based on description) |
| Commands (`*.md`) | `.cursor/rules/<command-name>.mdc` | Agent-requested (reference with @name) |
| Agents (`*.md`) | `.cursor/rules/<agent-name>.mdc` | Agent-requested (reference with @name) |
| CLAUDE.md conventions | `.cursor/rules/rein-project.mdc` | Always applied (alwaysApply: true) |

### Rule Reference

Reference rules in Cursor chat with `@<rule-name>`:

**Skills** (on-demand expertise):
- `@define` — Starting a new project or feature
- `@tdd` — Implementing logic or fixing bugs
- `@code-review` — Reviewing code before merge
- `@planning` — Breaking work into tasks
- `@security` — Security reviews
- `@performance` — Performance work
- `@frontend` — Building UI
- `@migration` — Database migrations
- `@integration-testing` — Integration test strategy
- `@shipping` — Pre-launch checklist

**Commands** (workflow triggers):
- `@feature` — Multi-file feature workflow (6 steps)
- `@fix` — Bug fix workflow
- `@plan` — Break work into tasks
- `@do` — Execute tasks incrementally
- `@ship` — Pre-launch checklist
- `@spec` — Write a structured spec
- `@code-review` — Five-axis code review

**Agents** (expert perspectives):
- `@code-reviewer` — Structured code review with severity ratings
- `@test-engineer` — Test strategy and coverage analysis
- `@security-auditor` — Security analysis and threat modeling

## Hooks

rein configures three hooks in `.cursor/hooks.json`:

| Hook | rein Command | What it does |
|------|-------------|-------------|
| PreEdit | `rein hook guard` | Blocks edits to rein-managed files (protected by .rein-manifest) |
| PostEdit | `rein hook format` | Auto-formats web files with prettier |
| PreCommit | `rein hook gate` | Runs tests before commit |

### Hook Limitations

These Claude Code hooks have no Cursor equivalent:

| Claude Code Hook | Why No Equivalent |
|-----------------|-------------------|
| `session-start` | Cursor has no session start event |
| `leak-guard` | No post-read/post-command hooks |
| `inject` | No prompt expansion hooks |
| `checkbox-guard` | No access to edit content diff |
| `artifact-validate` | No access to edit content for validation |
| `guard-bash` | No pre-command hooks |

## Usage Tips

1. **Reference rules explicitly** — Type `@code-review` in Cursor chat to ensure the rule loads
2. **Use agents for reviews** — `@code-reviewer Review this change` applies the reviewer persona
3. **CLI commands work normally** — `rein status`, `rein task done 1.2`, `rein validate` all work from Cursor's terminal
4. **Manifest protection is partial** — PreEdit blocks edits to protected files, but Cursor has no PreBash equivalent to block destructive shell commands
5. **The rein-project rule is always active** — Task progress tracking is enforced automatically
