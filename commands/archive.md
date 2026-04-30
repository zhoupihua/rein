Archive completed or abandoned rein artifacts (spec, plan, task) to `docs/rein/archive/`.

## Instructions

1. Scan `docs/rein/specs/`, `docs/rein/plans/`, `docs/rein/tasks/` for artifacts
2. List all found artifacts grouped by `YYYY-MM-DD-<name>` prefix
3. For each group, show:
   - Spec: exists/missing
   - Plan: exists/missing
   - Task: checkbox progress (e.g. 3/5 complete)
4. Ask user which group(s) to archive
5. Create archive directory: `docs/rein/archive/YYYY-MM-DD-<name>/`
6. Move matching spec, plan, task files into the archive directory
7. Confirm archived files and source directories are clean

## Arguments

- No argument: interactive mode — list and confirm
- `<name>`: archive the matching group directly (e.g. `/archive auth`)

## Output

```
Rein Artifacts:

  1. 2026-04-30-auth
     spec: ✓  plan: ✓  task: 5/5 complete

  2. 2026-05-01-export
     spec: ✓  plan: ✗  task: 2/4 complete

Archive which group(s)? (number, name, or "all"): 1

Archived 2026-04-30-auth:
  → docs/rein/archive/2026-04-30-auth/2026-04-30-auth-spec.md
  → docs/rein/archive/2026-04-30-auth/2026-04-30-auth-plan.md
  → docs/rein/archive/2026-04-30-auth/2026-04-30-auth-task.md
```
