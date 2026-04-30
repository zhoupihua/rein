Break down work into ordered tasks with dependency graphs.

## Instructions

1. Invoke `planning-and-task-breakdown` skill
2. Read the spec or requirements document (from `docs/rein/specs/` or user-provided)
3. Operate in read-only mode — no code writing during planning
4. Output two files:
   - `docs/rein/plans/YYYY-MM-DD-<feature-name>.md` — Architecture decisions, dependency graph, vertical slicing strategy, file map, parallelization, risks (decision layer)
   - `docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md` — Ordered task checklist with acceptance criteria, verification commands, dependencies, file paths, scope (execution layer, overwrites /spec's coarse tasks)
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

## Output

After generation, report:
```
Created plan: docs/rein/plans/YYYY-MM-DD-<feature-name>.md
Created tasks: docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md (N tasks, M phases)

Next step: /build to start implementing, or /resume to continue later
```
