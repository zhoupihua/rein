L3 full change — the 6-step workflow.

## The Full Workflow

### Step 1: Define
Invoke `refine` skill for divergent/convergent thinking, then `spec-driven` skill to write the PRD. The spec includes Context, Goals, Non-Goals, Requirements, Decisions, and Risks — all in one document. Save to `docs/rein/changes/<name>/spec.md` and commit.

No separate refine.md or design.md files. The refine thinking is internal; the design decisions are a section within spec.md.

### Step 2: Branch Setup
Create a feature branch from current branch. Ask the user whether to use worktree isolation:

- **直接开发（默认）**：在当前目录创建 feature 分支，直接开发。适合大多数场景。
- **Worktree 隔离**：调用 `git-worktrees` skill 创建隔离工作区。适合需要同时在多个分支上工作的场景。

If the user chooses worktree: invoke `git-worktrees` skill, create an isolated worktree with the new branch, verify clean test baseline.

If the user chooses direct development: create the feature branch and stay in the current directory.

### Step 3: Plan Tasks
Invoke `planning` skill. Break the spec into verifiable tasks with dependency graph, acceptance criteria, and checkpoints. Save plan to `docs/rein/changes/<name>/plan.md` and tasks to `docs/rein/changes/<name>/task.md`, then commit.

### Step 4: Implement
Invoke `incremental` + `tdd` skills:
- Build in thin vertical slices
- RED → GREEN → REFACTOR for each slice
- Commit after each verified increment

Routing:
- Frontend work → also invoke `frontend`
- API work → also invoke `api-design`
- Parallel tasks → invoke `subagent`
- Hit a bug → invoke `debugging`

### Step 5: Code Review
Invoke `code-review` skill. Five-axis review: Correctness, Readability, Architecture, Security, Performance. Save review report to `docs/rein/changes/<name>/review.md` and commit.

If issues found:
- Security concerns → invoke `security`
- Performance concerns → invoke `performance`

### Step 5.5: Integrate
Invoke `integration-testing` skill. This gate ensures testing sufficiency beyond per-task TDD:

1. **Spec traceability** — verify every spec scenario has a corresponding test
2. **Integration tests** — verify cross-component interfaces work together
3. **Coverage analysis** — new/changed code >= 80%, core paths 100%
4. **Regression check** — full test suite passes, no unexplained skips

If gaps are found at any gate, go back to Step 4 to add missing tests before proceeding.

### Step 6: Verify and Ship
Invoke `verify` skill. Verify with fresh evidence.

Then invoke `git-workflow` skill to:
- Verify all tests pass
- Present options: merge locally / create PR / keep / discard
- Commit with appropriate message

Post-merge:
- Invoke `shipping` for release checks (if applicable)
- Invoke `docs-and-adrs` for decision documentation
- Archive completed feature: move `docs/rein/changes/<name>/` to `docs/rein/archive/<name>/`
- If worktree was used: clean up the worktree

## Resuming

If interrupted at any step, use `/continue` to pick up where you left off.
