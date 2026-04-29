L3 full change — the 8-step iron triangle workflow.

## The Full Workflow

### Step 1: Idea Refine
Invoke `idea-refine` skill. Diverge and converge on the idea. Output a markdown one-pager with Problem Statement, Recommended Direction, Key Assumptions, MVP Scope, and Not Doing list. Save to `docs/specs/` and commit.

### Step 2: Spec-Driven Development
Invoke `spec-driven-development` skill. Write a PRD covering Objective, Commands, Project Structure, Code Style, Testing Strategy, and Boundaries. Save as SPEC.md and commit.

### Step 3: Generate Artifacts
Use `/spec <name>` to generate the full artifact set under `changes/<name>/`:
- `proposal.md` — Why and what
- `specs/<feature>/spec.md` — Delta specs
- `design.md` — Engineering decisions
- `tasks.md` — Task checklist

### Step 4: Branch Isolation
Invoke `using-git-worktrees` skill. Create an isolated worktree with a new branch. Verify clean test baseline.

### Step 5: Plan Tasks
Invoke `planning-and-task-breakdown` skill. Break the spec into verifiable tasks with dependency graph, acceptance criteria, and checkpoints. Save plan to `docs/plans/` and commit.

### Step 6: Implement
Invoke `incremental-implementation` + `test-driven-development` skills:
- Build in thin vertical slices
- RED → GREEN → REFACTOR for each slice
- Commit after each verified increment

Routing:
- Frontend work → also invoke `frontend-ui-engineering`
- API work → also invoke `api-and-interface-design`
- Parallel tasks → invoke `subagent-driven-development`
- Hit a bug → invoke `debugging-and-error-recovery`

### Step 7: Code Review
Invoke `code-review-and-quality` skill. Five-axis review: Correctness, Readability, Architecture, Security, Performance.

If issues found:
- Security concerns → invoke `security-and-hardening`
- Performance concerns → invoke `performance-optimization`

### Step 8: Verify and Ship
Invoke `verification-before-completion` skill. Verify with fresh evidence.

Then invoke `git-workflow-and-versioning` skill to:
- Verify all tests pass
- Present options: merge locally / create PR / keep / discard
- Commit with appropriate message

Post-merge:
- Invoke `shipping-and-launch` for release checks (if applicable)
- Invoke `documentation-and-adrs` for decision documentation
- Archive `changes/<name>/` to `archive/YYYY-MM-DD-<name>/`

## Resuming

If interrupted at any step, use `/resume` to pick up where you left off.