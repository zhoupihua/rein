Execute tasks from tasks.md incrementally. Replaces /opsx:apply.

## Instructions

1. Read `docs/rein/tasks/` for the current task file (latest or user-specified)
2. Read the corresponding `docs/rein/plans/` file for task details (acceptance criteria, verification, files)
3. Find the first unchecked task
4. For each task, execute this exact sequence — **do not skip any step**:
   a. Look up task details in plan.md
   b. Implement the task
   c. Verify tests pass
   d. Commit with descriptive message
   e. **Edit tasks.md: change `- [ ]` to `- [x]` for this task**
   f. **Read tasks.md back and verify the checkbox is now `- [x]`** — if not, fix it before proceeding
5. If a task is blocked:
   - Stop and report the blocker
   - Suggest using `debugging` if it's a bug
   - Ask the user for direction
6. After all tasks complete, suggest `/code-review` or `/ship`

## IRON RULE: No Checkbox Update = Task Not Done

**A task is NOT complete until its checkbox in tasks.md is updated AND verified.**

After every single task, you MUST do these two actions in order:
1. **Edit** tasks.md: change the task's `- [ ]` to `- [x]`
2. **Read** tasks.md: confirm the checkbox now shows `- [x]`

Only after both steps are done may you proceed to the next task.

This applies to ALL tasks — trivial, XS, batched, skipped-review, no exceptions.

## Task Execution Rules

- One task at a time
- Verify tests pass after each task
- Commit after each verified increment
- Edit + Read tasks.md after each commit (IRON RULE)
- Don't skip tasks or mark them complete without verification
- If scope expands, stop and update the plan
