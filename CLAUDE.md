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
- **`internal/cli/`** — Cobra commands: `validate`, `status`, `task {next|done|list}`, `instructions {apply|specs|tasks}`, `hook <name>`, `init`
- **`internal/artifact/`** — Parsers for markdown artifacts: `task.go` (checkbox task file), `spec.go` (PRD with requirements/scenarios/decisions/risks), `plan.go` (goal + task details), `epic.go` (epic with shared context/dependencies). Each parser uses regex-based line scanning.
- **`internal/project/`** — Project resolution (`CLAUDE_PROJECT_DIR` or cwd), feature discovery under `docs/rein/changes/`, phase validation logic. `PhaseArtifact` maps phases to required files.
- **`internal/hook/`** — Hook handlers called via `rein hook <name>`. Read tool input from `CLAUDE_TOOL_INPUT` env (or `CLAUDE_TOOL_INPUT_FILE_PATH`), output JSON with `decision: block` to reject actions or `additionalContext` to inject context.
- **`internal/output/`** — JSON/human output helpers

### Static Content (installed into target projects)

- **`skills/`** — 25 SKILL.md files organized by SDLC phase (DEFINE → PLAN → BUILD → VERIFY → REVIEW → SHIP)
- **`commands/`** — Slash command definitions consumed by Claude Code
- **`agents/`** — Expert persona prompts (code-reviewer, test-engineer, security-auditor)
- **`hooks/`** — Shell/PowerShell scripts + `hooks.json` wiring; each has `.sh` and `.ps1` variants
- **`references/`** — Checklists (testing, security, performance, accessibility)
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
  spec.md        # DEFINE phase (Context, Goals, Requirements, Decisions, Risks)
  plan.md        # PLAN phase (### N.N tasks with Acceptance/Verification/Dependencies/Files)
  task.md        # PLAN phase (checkbox format: - [ ] 1.1 description)
  review.md      # REVIEW phase
docs/rein/archive/<name>/   # shipped features
```

`task.md` is the single source of truth for build progress. The `artifact.ParseTaskFile` parser recognizes `## N. PhaseName` headings and `- [ ] N.N description` checkboxes.

`spec.md` is the single DEFINE artifact — it includes refine thinking (Context, Goals, Non-Goals) and design decisions (`**Decision:** ... — **Rationale:** ...`) alongside requirements.

## Key Conventions

- Task IDs use `phase.seq` format (e.g., `1.1`, `2.3`)
- Spec scenarios use `WHEN`/`THEN` format parsed by regex in `spec.go`; decisions use `**Decision:** ... — **Rationale:** ...` format
- Plan task details use bold-labeled fields (`**Acceptance:**`, `**Verification:**`, etc.) parsed by regex in `plan.go`
- Hook communication: read `CLAUDE_TOOL_INPUT` (JSON), output `{"decision":"block","reason":"..."}` or `{"hookSpecificOutput":{"additionalContext":"..."}}`
- All hooks have both `.sh` and `.ps1` implementations
- No external dependencies beyond cobra; stdlib-only for all internal packages
