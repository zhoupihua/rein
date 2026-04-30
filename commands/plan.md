Break down work into ordered tasks with dependency graphs.

## Instructions

1. Invoke `planning` skill
2. Read the spec or requirements document (from `docs/rein/specs/` or user-provided)
3. Operate in read-only mode — no code writing during planning
4. **MUST output TWO files** (both are required, do not skip tasks):
   - `docs/rein/plans/YYYY-MM-DD-<feature-name>.md` — Architecture decisions, dependency graph, vertical slicing strategy, file map, parallelization, risks (decision layer)
   - `docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md` — Ordered task checklist with acceptance criteria, verification commands, dependencies, file paths, scope (execution layer)
5. Offer execution choice: subagent-driven (recommended) or inline

## Task Format in tasks.md

Each task uses checkbox format with inline metadata:

```
- [ ] 1.1 [Short descriptive title]
  - Acceptance: [Specific, testable condition]
  - Verification: [test command]
  - Dependencies: [Task numbers or "None"]
  - Files: `path/to/file.ts`
  - Scope: [XS | S | M | L]
```

## Self-Check

After generation, verify both files exist and are non-empty:
1. `docs/rein/plans/YYYY-MM-DD-<feature-name>.md` has architecture decisions and dependency graph
2. `docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md` has checkbox tasks with acceptance criteria

If either file is missing, generate it before reporting completion.

## Output

After generation, report:
```
Created plan: docs/rein/plans/YYYY-MM-DD-<feature-name>.md
Created tasks: docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md (N tasks, M phases)

Next step: /do to start implementing, or /continue to continue later
```
