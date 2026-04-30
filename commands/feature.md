L3 full change — the 8-step iron triangle workflow.

## The Full Workflow

### Step 1: Idea Refine
Invoke `refine` skill. Diverge and converge on the idea. Output a markdown one-pager with Problem Statement, Recommended Direction, Key Assumptions, MVP Scope, and Not Doing list. Save to `docs/rein/specs/YYYY-MM-DD-<topic>-design.md` and commit.

### Step 2: Spec-Driven Development
Invoke `spec-driven` skill. Write a PRD covering Objective, Commands, Project Structure, Code Style, Testing Strategy, and Boundaries. Save to `docs/rein/specs/YYYY-MM-DD-<topic>-design.md` and commit.

### Step 3: Generate Design Spec
Use `/spec <name>` to generate the design spec:
- `docs/rein/specs/YYYY-MM-DD-<name>-design.md` — Design spec with requirements and decisions

### Step 4: Branch Isolation
Invoke `git-worktrees` skill. Create an isolated worktree with a new branch. Verify clean test baseline.

### Step 5: Plan Tasks
Invoke `planning` skill. Break the spec into verifiable tasks with dependency graph, acceptance criteria, and checkpoints. Save plan to `docs/rein/plans/` and tasks to `docs/rein/tasks/`, then commit.

### Step 6: Implement
Invoke `incremental` + `tdd` skills:
- Build in thin vertical slices
- RED → GREEN → REFACTOR for each slice
- Commit after each verified increment

Routing:
- Frontend work → also invoke `frontend`
- API work → also invoke `api-design`
- Parallel tasks → invoke `subagent`
- Hit a bug → invoke `debugging`

### Step 7: Code Review
Invoke `code-review` skill. Five-axis review: Correctness, Readability, Architecture, Security, Performance.

If issues found:
- Security concerns → invoke `security`
- Performance concerns → invoke `performance`

### Step 8: Verify and Ship
Invoke `verify` skill. Verify with fresh evidence.

Then invoke `git-workflow` skill to:
- Verify all tests pass
- Present options: merge locally / create PR / keep / discard
- Commit with appropriate message

Post-merge:
- Invoke `shipping` for release checks (if applicable)
- Invoke `docs-and-adrs` for decision documentation
- Archive completed `docs/rein/` artifacts to `docs/rein/archive/YYYY-MM-DD-<name>/`

## Resuming

If interrupted at any step, use `/continue` to pick up where you left off.