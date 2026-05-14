# Getting Started with rein

rein works with any AI coding agent that accepts Markdown instructions. This guide covers the universal approach. For tool-specific setup, see the dedicated guides.

## How Skills Work

Each skill is a Markdown file (`SKILL.md`) that describes a specific engineering workflow. When loaded into an agent's context, the agent follows the workflow — including verification steps, anti-patterns to avoid, and exit criteria.

**Skills are not reference docs.** They're step-by-step processes the agent follows.

## Quick Start (Any Agent)

### 1. Install rein

```bash
# Option A: Go install + init (CLI only, no skills/commands)
go install github.com/zhoupihua/rein/cmd/rein@latest
cd /your/project
rein init

# Option B: Full install (skills, commands, agents into .rein/)
cd /your/project
bash /path/to/rein/install/install.sh       # Linux/Mac
powershell -File \path\to\rein\install\install.ps1  # Windows
```

`rein init` creates `docs/rein/{changes,archive}` directories. The install script copies skills, commands, agents into your project's `.rein/` directory and configures hooks in `.claude/settings.json`.

### 2. Choose a skill

Browse the `skills/` directory. Each subdirectory contains a `SKILL.md` with:
- **When to use** — triggers that indicate this skill applies
- **Process** — step-by-step workflow
- **Verification** — how to confirm the work is done
- **Common rationalizations** — excuses the agent might use to skip steps
- **Red flags** — signs the skill is being violated

### 3. Load the skill into your agent

Copy the relevant `SKILL.md` content into your agent's system prompt, rules file, or conversation:

**System prompt:** Paste the skill content at the start of the session.

**Rules file:** Add skill content to your project's rules file (CLAUDE.md, .cursorrules, etc.).

**Conversation:** Reference the skill when giving instructions: "Follow the tdd process for this change."

### 4. Use the meta-skill for discovery

Start with the `using-rein` skill loaded. It contains a flowchart that maps task types to the appropriate skill.

## Recommended Setup

### Minimal (Start here)

Load three essential skills into your rules file:

1. **define** — For defining what to build (explore requirements + write spec)
2. **tdd** — For proving it works
3. **code-review** — For verifying quality before merge (includes simplification)

These three cover the most critical quality gaps in AI-assisted development.

### Full Lifecycle

For comprehensive coverage, load skills by phase:

```
Starting a project:  define → planning
During development:  executing-plans + tdd
Before merge:        code-review + security
Before deploy:       shipping
```

### Context-Aware Loading

Don't load all skills at once — it wastes context. Load skills relevant to the current task:

- Working on UI? Load `frontend`
- Debugging? Load `debugging`
- Setting up CI? Load `shipping` (includes CI/CD automation)

## Skill Anatomy

Every skill follows the same structure:

```
YAML frontmatter (name, description)
├── Overview — What this skill does
├── When to Use — Triggers and conditions
├── Core Process — Step-by-step workflow
├── Examples — Code samples and patterns
├── Common Rationalizations — Excuses and rebuttals
├── Red Flags — Signs the skill is being violated
└── Verification — Exit criteria checklist
```

See [skill-anatomy.md](skill-anatomy.md) for the full specification.

## Using Agents

The `agents/` directory contains pre-configured agent personas:

| Agent | Purpose |
|-------|---------|
| `code-reviewer.md` | Five-axis code review |
| `test-engineer.md` | Test strategy and writing |
| `security-auditor.md` | Vulnerability detection |

Load an agent definition when you need specialized review. For example, ask your coding agent to "review this change using the code-reviewer agent persona" and provide the agent definition.

## Using Commands

The `commands/` directory contains slash commands for Claude Code:

| Command | Purpose |
|---------|---------|
| `/quick` | L1 lightweight change (≤5 lines) |
| `/fix` | L2 bug fix or small feature (1-3 files) |
| `/feature` | L3 full feature (6-step workflow) |
| `/spec` | Generate proposal.md + spec.md |
| `/plan` | Break spec into tasks |
| `/do` | Execute tasks from task.md (includes TDD) |
| `/code-review` | Five-axis code review (includes simplification) |
| `/ship` | Pre-launch fan-out + GO/NO-GO |
| `/continue` | Resume from breakpoint |
| `/status` | Show task progress |
| `/archive` | Archive completed features |

## Using References

The `references/` directory contains supplementary checklists:

| Reference | Use With |
|-----------|----------|
| `testing-patterns.md` | tdd |
| `performance-checklist.md` | performance |
| `security-checklist.md` | security |
| `accessibility-checklist.md` | frontend |
| `orchestration-patterns.md` | subagent |
| `api-design.md` | code-review, planning |

## Spec and Task Artifacts

The `/spec` and `/plan` commands create working artifacts under `docs/rein/changes/<name>/`:

```
docs/rein/changes/<name>/
  proposal.md — DEFINE phase (Why, What Changes, Goals, Non-Goals, Assumptions, Open Questions) — optional for L2
  spec.md    — DEFINE phase (Requirements, Decisions, Risks)
  plan.md    — PLAN phase (Architecture, Dependency Graph, Task Details)
  task.md    — PLAN phase (checkbox progress tracking)
  review.md  — REVIEW phase
```

Treat them as **living documents** while the work is in progress. Update them when scope or decisions change. After shipping, archive them to `docs/rein/archive/`.

## Using the Go CLI

```bash
rein init                          # Create docs/rein/{changes,archive} directories
rein validate [feature]            # Check artifact completeness
rein status [feature]              # Show phase and task progress
rein task next [feature]           # Find next unchecked task
rein task done N.M                 # Mark task complete (e.g., rein task done 1.1)
rein task list [feature]           # List all tasks with status
rein visual start                  # Start visual brainstorming server
rein visual stop                   # Stop visual brainstorming server
rein instructions apply            # Apply rein instructions to project
rein instructions specs            # Generate agent instructions from specs
rein instructions tasks            # Generate agent instructions from tasks
rein hook <name>                   # Run a hook handler
```

All commands support `--json` for machine-readable output.

## Tips

1. **Start with define** for any non-trivial work — explore requirements before coding
2. **Always load tdd** when writing code
3. **Don't skip verification steps** — they're the whole point
4. **Load skills selectively** — more context isn't always better
5. **Use the agents for review** — different perspectives catch different issues
6. **Use `/feature`** for multi-file features — it runs the full L3 6-step workflow
