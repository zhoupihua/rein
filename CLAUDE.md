# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o rein ./cmd/rein          # Build CLI
go build ./...                        # Build all packages
go test ./...                         # Run all tests
go test ./internal/artifact/...       # Run a single package's tests
go run ./cmd/rein                     # Run without building binary
```

The CLI uses [cobra](https://github.com/spf13/cobra) (`go.mod` has only one direct dependency). All commands support `--json` for machine-readable output.

## Architecture

rein is a zero-dependency AI coding workflow toolkit. It ships as static markdown files (skills, commands, agents, hooks, templates, references) plus a Go CLI for artifact validation and task management.

### Go Code Structure

- **`cmd/rein/main.go`** — Entry point, delegates to `cli.Execute()`
- **`internal/cli/`** — Cobra commands: `validate`, `status`, `task {next|done|list}`, `instructions {apply|specs|tasks}`, `visual {start|stop}`, `hook <name>`, `init`
- **`internal/artifact/`** — Parsers for markdown artifacts: `proposal.go` (Why/What/Goals/Assumptions), `spec.go` (PRD with requirements/scenarios/decisions/risks), `task.go` (checkbox task file), `plan.go` (goal + architecture + task details), `graph.go` (artifact DAG with topological sort), `delta.go` + `merge.go` (incremental spec operations). Each parser uses regex-based line scanning.
- **`internal/project/`** — Project resolution (`CLAUDE_PROJECT_DIR` or cwd), feature discovery under `docs/rein/changes/`, phase validation logic. `PhaseArtifact` maps phases to required files.
- **`internal/hook/`** — Hook handlers called via `rein hook <name>`. Read tool input from `CLAUDE_TOOL_INPUT` env (or `CLAUDE_TOOL_INPUT_FILE_PATH`), output JSON with `decision: block` to reject actions or `additionalContext` to inject context.
- **`internal/output/`** — JSON/human output helpers
- **`internal/visual/`** — Visual brainstorming server (HTTP + hand-rolled WebSocket), directory polling, PID monitoring

### Static Content (installed into target projects)

- **`skills/`** — 28 SKILL.md files organized by SDLC phase (DEFINE → PLAN → BUILD → VERIFY → REVIEW → SHIP)
- **`commands/`** — Slash command definitions consumed by Claude Code
- **`agents/`** — Expert persona prompts (code-reviewer, test-engineer, security-auditor)
- **`hooks/`** — Shell/PowerShell scripts + `hooks.json` wiring; each has `.sh` and `.ps1` variants
- **`references/`** — Checklists (testing, security, performance, accessibility, orchestration)
- **`templates/`** — Artifact markdown templates (proposal, spec, tasks)

### Install Flow

`install/install.sh` (or `.ps1`) copies skills/commands/agents/hooks into `<project>/.claude/`, configures `settings.json` hooks, and generates `.rein-manifest` for file protection. The guard hooks read this manifest to block AI edits to rein-managed files.

## Workflow Levels

| Level | Command | Scope | Flow |
|-------|---------|-------|------|
| L1 | `/quick` | ≤5 lines, no logic | Edit → test → commit |
| L2 | `/fix` | 1-3 files | DEFINE → BUILD → VERIFY → SHIP |
| L3 | `/feature` | Multi-file feature | 6-step: define(spec) → branch → plan → implement → review → ship |

Each level has quality gates. Phase transitions are detected by checking which artifacts exist under `docs/rein/changes/<name>/`.

## Artifact Directory

```
docs/rein/changes/<name>/
  proposal.md    # DEFINE phase (Why, What Changes, Goals, Non-Goals, Assumptions, Open Questions) — optional for L2
  spec.md        # DEFINE phase (Requirements, Decisions, Risks)
  plan.md        # PLAN phase (Architecture, Dependency Graph, ### N.N tasks with details)
  task.md        # PLAN phase (checkbox format: - [ ] 1.1 description, with optional RED/GREEN/REFACTOR sub-tasks)
  review.md      # REVIEW phase
docs/rein/archive/<name>/   # shipped features
```

`task.md` is the single source of truth for build progress. The `artifact.ParseTaskFile` parser recognizes `## N. PhaseName` headings, `- [ ] N.N description` checkboxes, and RED/GREEN/REFACTOR sub-task checkboxes (`  - [ ] RED: ...`, `  - [x] GREEN: ...`). Sub-tasks are first-class: each has an index, done state, and auto-checks the parent task when all are complete.

`spec.md` is the PRD artifact — it includes requirements (with WHEN/THEN/TEST scenarios), design decisions (`**Decision:** ... — **Rationale:** ...`), and risks. Context, Goals, and Non-Goals live in `proposal.md` (output of the refine skill). For L2 `/fix` workflows, spec.md can be generated directly without proposal.md.

`plan.md` follows Superpowers conventions: Architecture Overview, Dependency Graph (ASCII), Vertical Slice Strategy, Risk/Mitigation Table, Parallelization Classification, Self-Audit Checklist, Handoff Statement, plus Task Details with Approach/Edge Cases/Rollback fields.

## Key Conventions

- Task IDs use `phase.seq` format (e.g., `1.1`, `2.3`); sub-task IDs use `phase.seq.index` format (e.g., `1.1.0`, `1.1.2`)
- Spec scenarios use `WHEN`/`THEN`/`TEST` format parsed by regex in `spec.go`; decisions use `**Decision:** ... — **Rationale:** ...` format
- Plan has section-level fields (Architecture Overview, Dependency Graph, etc.) parsed by `## ` heading accumulation in `plan.go`; task details use bold-labeled fields (`**Acceptance:**`, `**Approach:**`, etc.) parsed by regex
- Hook communication: read `CLAUDE_TOOL_INPUT` (JSON), output `{"decision":"block","reason":"..."}` or `{"hookSpecificOutput":{"additionalContext":"..."}}`
- All hooks have both `.sh` and `.ps1` implementations
- No external dependencies beyond cobra; stdlib-only for all internal packages
