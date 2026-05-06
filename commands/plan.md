Break down work into ordered tasks with dependency graphs.

## Instructions

1. Invoke `planning` skill
2. Read the spec or requirements document (from `docs/rein/changes/<name>/` or user-provided)
3. Operate in read-only mode — no code writing during planning
4. **MUST output TWO files** (both are required, do not skip tasks):
   - `docs/rein/changes/<name>/plan.md` — Architecture decisions, dependency graph, vertical slicing, **task details** (acceptance criteria, verification, files, dependencies, scope, notes)
   - `docs/rein/changes/<name>/task.md` — Simple checkbox list for status tracking only (no nested metadata)
5. Offer execution choice: subagent-driven (recommended) or inline

## Two-File Responsibility Split

**plan.md** = HOW (implementation reference):
- Architecture decisions, dependency graph, risks
- Per-task details: acceptance criteria, verification commands, files, dependencies, scope, implementation notes

**task.md** = STATUS (progress tracking):
- Simple checkbox list, one line per task
- Grouped by phase with `##` headings
- Numbered for easy cross-reference with plan.md
- No Acceptance/Verification/Files/Scope — those live in plan.md

```
## task.md format (simple):
- [ ] 1.1 Create database migration for X
- [ ] 1.2 Implement repository layer for X

## plan.md format (detailed):
### 1.1 Create database migration for X
- **Acceptance:** migration runs, table exists
- **Verification:** `make db-migrate && make db-verify`
- **Dependencies:** None
- **Files:** `pkg/sqlmigration/XXX_create_x.go`
- **Scope:** S
- **Notes:** Use bun/migrate pattern, see existing migrations
```

## Self-Check

After generation, verify both files exist and are non-empty:
1. `docs/rein/changes/<name>/plan.md` has architecture decisions, dependency graph, AND Task Details sections
2. `docs/rein/changes/<name>/task.md` has simple checkbox tasks (no nested metadata)
3. Task numbers match between the two files

If either file is missing, generate it before reporting completion.

## Output

After generation, report:
```
Created plan: docs/rein/changes/<name>/plan.md
Created tasks: docs/rein/changes/<name>/task.md (N tasks, M phases)

Next step: /do to start implementing, or /continue to continue later
```
