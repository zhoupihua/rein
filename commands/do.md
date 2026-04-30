Execute tasks from tasks.md incrementally. Replaces /opsx:apply.

## Instructions

1. Read `docs/rein/tasks/` for the current task file (latest or user-specified)
2. Read the corresponding `docs/rein/plans/` file for task details (acceptance criteria, verification, files)
3. Find the first unchecked task
4. For each task:
   a. Look up task details in plan.md
   b. Implement the task
   c. Verify tests pass
   d. Commit with descriptive message
   e. **Edit tasks.md: change `- [ ]` to `- [x]` for this task** (MANDATORY — do not skip)
5. If a task is blocked:
   - Stop and report the blocker
   - Suggest using `debugging` if it's a bug
   - Ask the user for direction
6. After all tasks complete, suggest `/code-review` or `/ship`

## IRON RULE

**A task is NOT complete until its checkbox in tasks.md is updated.**
You MUST Edit tasks.md after each completed task, BEFORE starting the next one.
This is not optional — not for trivial tasks, not for small tasks, not ever.

## Task Execution Rules

- One task at a time
- Verify tests pass after each task
- Commit after each verified increment
- Update tasks.md checkbox after each commit (mandatory)
- Don't skip tasks or mark them complete without verification
- If scope expands, stop and update the plan