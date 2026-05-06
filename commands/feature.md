L3 full change — the 8-step iron triangle workflow.

## The Full Workflow

### Step 1: Idea Refine
Invoke `refine` skill. Diverge and converge on the idea. Output a markdown one-pager with Problem Statement, Recommended Direction, Key Assumptions, MVP Scope, and Not Doing list. Save to `docs/rein/changes/<name>/refine.md` and commit.

### Step 2: Spec-Driven Development
Invoke `spec-driven` skill. Write a PRD covering Objective, Commands, Project Structure, Code Style, Testing Strategy, and Boundaries. Save to `docs/rein/changes/<name>/spec.md` and commit.

### Step 3: Generate Design Spec
Use `/spec <name>` to generate the design spec:
- `docs/rein/changes/<name>/design.md` — Design spec with requirements and decisions

### Step 4: Branch Setup
Create a feature branch from current branch. Ask the user whether to use worktree isolation:

- **直接开发（默认）**：在当前目录创建 feature 分支，直接开发。适合大多数场景。
- **Worktree 隔离**：调用 `git-worktrees` skill 创建隔离工作区。适合需要同时在多个分支上工作的场景。

If the user chooses worktree: invoke `git-worktrees` skill, create an isolated worktree with the new branch, verify clean test baseline.

If the user chooses direct development: create the feature branch and stay in the current directory.

### Step 5: Plan Tasks
Invoke `planning` skill. Break the spec into verifiable tasks with dependency graph, acceptance criteria, and checkpoints. Save plan to `docs/rein/changes/<name>/plan.md` and tasks to `docs/rein/changes/<name>/task.md`, then commit.

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
Invoke `code-review` skill. Five-axis review: Correctness, Readability, Architecture, Security, Performance. Save review report to `docs/rein/changes/<name>/review.md` and commit.

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
- Archive completed feature: move `docs/rein/changes/<name>/` to `docs/rein/archive/<name>/`
- If worktree was used: clean up the worktree

## Resuming

If interrupted at any step, use `/continue` to pick up where you left off.
