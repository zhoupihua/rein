Break work into ordered tasks with dependency graphs.

## Instructions

1. Invoke `planning-and-task-breakdown` skill
2. Read the spec or requirements document
3. Operate in read-only mode — no code writing during planning
4. Output a plan document with:
   - Dependency graph (ASCII tree)
   - Vertical slicing strategy
   - Task breakdown with acceptance criteria
   - Task sizing (XS/S/M/L)
   - Parallelization classification (safe/sequential/needs-coordination)
   - Checkpoints between phases
   - Risks and mitigations table
5. Save plan to `docs/plans/YYYY-MM-DD-<name>.md`
6. Offer execution choice: subagent-driven (recommended) or inline