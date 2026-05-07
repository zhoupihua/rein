---
name: git-workflow
description: Use when making any code change — committing, branching, resolving conflicts, completing a development branch. Git is your safety net: commits are save points, branches are sandboxes, history is documentation.
---

# Git Workflow and Versioning

## Core Principles

### Trunk-Based Development (Recommended)

Keep `main` always deployable. Work in short-lived feature branches that merge back within 1-3 days. Long-lived development branches are hidden costs.

```
main ──●──●──●──●──●──●──●──●──●──  (always deployable)
        ╲      ╱  ╲    ╱
         ●──●─╱    ●──╱    ← short-lived feature branches (1-3 days)
```

### 1. Commit Early, Commit Often

Each successful increment gets its own commit. Commits are save points.

```
Work pattern:
  Implement slice → Test → Verify → Commit → Next slice
```

### 2. Atomic Commits

Each commit does one logical thing:

```
# Good
a1b2c3d Add task creation endpoint with validation
d4e5f6g Add task creation form component
h7i8j9k Connect form to API and add loading state

# Bad
x1y2z3a Add task feature, fix sidebar, update deps, refactor utils
```

### 3. Descriptive Messages

Commit messages explain the *why*, not just the *what*:

```
<type>: <short description>

<optional body explaining why, not what>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`, `ci`

### 4. Keep Concerns Separate

Don't combine formatting changes with behavior changes. Don't combine refactors with features.

### 5. Size Your Changes

```
~100 lines  → Easy to review, easy to revert
~300 lines  → Acceptable for a single logical change
~1000 lines → Split into smaller changes
```

## Branching Strategy

### Feature Branches

```
main (always deployable)
  │
  ├── feature/task-creation    ← One feature per branch
  ├── feature/user-settings    ← Parallel work
  └── fix/duplicate-tasks      ← Bug fixes
```

- Branch from `main`
- Keep branches short-lived (merge within 1-3 days)
- Delete branches after merge
- Prefer feature flags over long-lived branches

### Branch Naming

```
feature/<short-description>
fix/<short-description>
chore/<short-description>
refactor/<short-description>
```

## Working with Worktrees

For parallel AI agent work, use git worktrees:

```bash
git worktree add ../project-feature-a feature/task-creation
git worktree add ../project-feature-b feature/user-settings

# Each worktree is a separate directory with its own branch
# When done, merge and clean up
git worktree remove ../project-feature-a
```

## The Save Point Pattern

```
Agent starts work
    │
    ├── Makes a change
    │   ├── Test passes? → Commit → Continue
    │   └── Test fails? → Revert to last commit → Investigate
    │
    └── Feature complete → All commits form a clean history
```

If an agent goes off the rails, `git reset --hard HEAD` takes you back to the last successful state.

## Change Summaries

After any modification, provide a structured summary:

```
CHANGES MADE:
- src/routes/tasks.ts: Added validation middleware to POST endpoint
- src/lib/validation.ts: Added TaskCreateSchema using Zod

THINGS I DIDN'T TOUCH (intentionally):
- src/routes/auth.ts: Has similar validation gap but out of scope

POTENTIAL CONCERNS:
- The Zod schema is strict — rejects extra fields
```

## Pre-Commit Hygiene

Before every commit:

```bash
git diff --staged                                    # Check what you're committing
git diff --staged | grep -i "password\|secret\|api_key"  # No secrets
npm test                                             # Tests pass
npm run lint                                         # Linting passes
```

## Completing a Development Branch

When implementation is complete, all tests pass, and you need to integrate the work:

### Step 1: Verify Tests

```bash
npm test / cargo test / pytest / go test ./...
```

If tests fail: Stop. Don't proceed.

### Step 2: Detect Environment

Determine workspace state before presenting options:

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
```

| State | Menu | Cleanup |
|-------|------|---------|
| `GIT_DIR == GIT_COMMON` (normal repo) | Standard 4 options | No worktree to clean up |
| `GIT_DIR != GIT_COMMON`, named branch | Standard 4 options | Provenance-based (see Step 6) |
| `GIT_DIR != GIT_COMMON`, detached HEAD | Reduced 3 options (no merge) | No cleanup (externally managed) |

### Step 3: Determine Base Branch

```bash
git merge-base HEAD main 2>/dev/null || git merge-base HEAD master 2>/dev/null
```

### Step 4: Present Options

**Normal repo and named-branch worktree — present exactly these 4 options:**

```
1. Merge back to <base-branch> locally
2. Push and create a Pull Request
3. Keep the branch as-is (I'll handle it later)
4. Discard this work
```

**Detached HEAD — present exactly these 3 options:**

```
1. Push as new branch and create a Pull Request
2. Keep as-is (I'll handle it later)
3. Discard this work
```

Don't add explanation — keep options concise.

### Step 5: Execute Choice

#### Option 1: Merge Locally

```bash
# Get main repo root for CWD safety
MAIN_ROOT=$(git -C "$(git rev-parse --git-common-dir)/.." rev-parse --show-toplevel)
cd "$MAIN_ROOT"

# Merge first — verify success before removing anything
git checkout <base-branch>
git pull
git merge <feature-branch>

# Verify tests on merged result
<test command>
```

Then: Cleanup worktree (Step 6), then delete branch: `git branch -d <feature-branch>`

#### Option 2: Push and Create PR

```bash
git push -u origin <feature-branch>
gh pr create --title "<title>" --body "<description>"
```

**Do NOT clean up worktree** — user needs it alive to iterate on PR feedback.

#### Option 3: Keep As-Is

Report branch and worktree location. Don't cleanup.

#### Option 4: Discard

**Confirm first:**
```
This will permanently delete:
- Branch <name>
- All commits: <commit-list>
- Worktree at <path>

Type 'discard' to confirm.
```

If confirmed, cleanup worktree (Step 6), then: `git branch -D <feature-branch>`

### Step 6: Cleanup Workspace

**Only runs for Options 1 and 4.** Options 2 and 3 always preserve the worktree.

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
WORKTREE_PATH=$(git rev-parse --show-toplevel)
```

**If `GIT_DIR == GIT_COMMON`:** Normal repo, no worktree to clean up. Done.

**If worktree path is under `.worktrees/`, `worktrees/`, or `.claude/worktrees/`:** rein created this worktree — we own cleanup.

```bash
MAIN_ROOT=$(git -C "$(git rev-parse --git-common-dir)/.." rev-parse --show-toplevel)
cd "$MAIN_ROOT"
git worktree remove "$WORKTREE_PATH"
git worktree prune  # Self-healing: clean up any stale registrations
```

**Otherwise:** The host environment owns this workspace. Do NOT remove it.

## Using Git for Debugging

```bash
git bisect start                  # Find which commit introduced a bug
git bisect bad HEAD
git bisect good <known-good>

git log --oneline -20             # View recent changes
git blame src/file.ts             # Find who last changed a line
git log --grep="keyword" --oneline  # Search commit messages
```

## Common Rationalizations

| Rationalization | Reality |
|---|---|
| "I'll commit when the feature is done" | One giant commit is impossible to review or revert. Commit each slice. |
| "The message doesn't matter" | Messages are documentation. Future you needs to understand why. |
| "I'll squash it all later" | Squashing destroys the development narrative. Clean incremental commits from the start. |
| "I'll split this change later" | Large changes are harder to review, riskier to deploy. Split before submitting. |

## Red Flags

- Large uncommitted changes accumulating
- Commit messages like "fix", "update", "misc"
- Formatting changes mixed with behavior changes
- No `.gitignore` in the project
- Committing `.env` or build artifacts
- Long-lived branches diverging from main
- Force-pushing to shared branches
- Proceeding with failing tests
- Deleting work without confirmation
- Merging without verifying tests on result
- Removing a worktree you didn't create (provenance check)
- Running `git worktree remove` from inside the worktree
- Cleaning up worktree for PR option (user needs it for iteration)

## Verification

For every commit:

- [ ] Commit does one logical thing
- [ ] Message explains the why, follows type conventions
- [ ] Tests pass before committing
- [ ] No secrets in the diff
- [ ] No formatting-only changes mixed with behavior changes
