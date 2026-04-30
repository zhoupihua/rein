Execute tasks from tasks.md incrementally. Replaces /opsx:apply.

## Instructions

1. Read `docs/rein/tasks/` for the current task file (latest or user-specified)
2. Read the corresponding `docs/rein/plans/` file for task details (acceptance criteria, verification, files)
3. Find the first unchecked task
4. For each task:
   a. Look up task details in plan.md
   b. Mark as in-progress
   c. Invoke `incremental` + `tdd`
   d. Build in thin vertical slices: implement → test → verify → commit
   e. Mark the task checkbox in tasks.md as complete (`- [x]`)
   f. Commit with descriptive message
4. If a task is blocked:
   - Stop and report the blocker
   - Suggest using `debugging` if it's a bug
   - Ask the user for direction
5. After all tasks complete, suggest `/code-review` or `/ship`

## Task Execution Rules

- One task at a time
- Verify tests pass after each task
- Commit after each verified increment
- Don't skip tasks or mark them complete without verification
- If scope expands, stop and update the plan