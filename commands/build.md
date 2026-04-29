Execute tasks from tasks.md incrementally. Replaces /opsx:apply.

## Instructions

1. Read `changes/<name>/tasks.md` for the current change
2. Find the first unchecked task
3. For each task:
   a. Mark as in-progress
   b. Invoke `incremental-implementation` + `test-driven-development`
   c. Build in thin vertical slices: implement → test → verify → commit
   d. Mark the task checkbox as complete (`- [x]`)
   e. Commit with descriptive message
4. If a task is blocked:
   - Stop and report the blocker
   - Suggest using `debugging-and-error-recovery` if it's a bug
   - Ask the user for direction
5. After all tasks complete, suggest `/review` or `/ship`

## Task Execution Rules

- One task at a time
- Verify tests pass after each task
- Commit after each verified increment
- Don't skip tasks or mark them complete without verification
- If scope expands, stop and update the plan