# Using rein with Gemini CLI

## Setup

### Option 1: Install as Skills (Recommended)

Gemini CLI has a native skills system that auto-discovers `SKILL.md` files in `.gemini/skills/` directories.

**Install from a local clone:**

```bash
git clone https://github.com/zhoupihua/rein.git
gemini skills install /path/to/rein/skills/
```

**Install for a specific workspace only:**

```bash
gemini skills install /path/to/rein/skills/ --scope workspace
```

Skills installed at workspace scope go into `.gemini/skills/`. User-level skills go into `~/.gemini/skills/`.

Once installed, verify with:

```
/skills list
```

### Option 2: GEMINI.md (Persistent Context)

For skills you want always loaded as persistent project context:

```bash
# Create GEMINI.md with core skills
cat .claude/skills/incremental/SKILL.md > GEMINI.md
echo -e "\n---\n" >> GEMINI.md
cat .claude/skills/code-review/SKILL.md >> GEMINI.md
```

Or modularize by importing from separate files:

```markdown
# Project Instructions

@skills/tdd/SKILL.md
@skills/incremental/SKILL.md
```

> **Skills vs GEMINI.md:** Skills are on-demand expertise that activate only when relevant, keeping your context window clean. GEMINI.md provides persistent context loaded for every prompt.

## Recommended Configuration

### Always-On (GEMINI.md)

Add these as persistent context for every session:

- `incremental` — Build in small verifiable slices
- `code-review` — Five-axis review

### On-Demand (Skills)

Install these as skills so they activate only when relevant:

- `tdd` — Activates when implementing logic or fixing bugs
- `spec-driven` — Activates when starting a new project or feature
- `frontend` — Activates when building UI
- `security` — Activates during security reviews
- `performance` — Activates during performance work

## Slash Commands

The repo ships slash commands under `.gemini/commands/` that map to the development lifecycle:

| Command | What it does |
|---------|--------------|
| `/spec` | Write a structured spec before writing code |
| `/planning` | Break work into small, verifiable tasks |
| `/build` | Implement the next task incrementally |
| `/test` | Run TDD workflow — red, green, refactor |
| `/review` | Five-axis code review |
| `/code-simplify` | Reduce complexity without changing behavior |
| `/ship` | Pre-launch checklist via parallel persona fan-out |

> **Note:** Use `/planning` instead of `/plan` — `/plan` conflicts with a Gemini CLI internal command name.

## Usage Tips

1. **Prefer skills over GEMINI.md** — Skills activate on demand and keep your context window focused
2. **Skill descriptions matter** — Each SKILL.md has a `description` field in its frontmatter that tells agents when to activate it
3. **Use agents for review** — Copy `agents/code-reviewer.md` content when requesting structured code reviews
4. **Combine with references** — Reference checklists from `references/` when working on specific quality areas
