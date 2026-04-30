Execute tasks from task.md incrementally. Replaces /opsx:apply.

## Instructions

1. Read `docs/rein/tasks/` for the current task file (latest or user-specified)
2. Read the corresponding `docs/rein/plans/` file for task details (acceptance criteria, verification, files)
3. **Verify plan-task consistency** (see below) — if mismatch, stop and report
4. Find the first unchecked task
5. For each task, execute this exact sequence — **do not skip any step**:
   a. Look up task details in plan.md
   b. Implement the task
   c. Verify tests pass
   d. Commit with descriptive message
   e. **Edit task.md: change `- [ ]` to `- [x]` for this task**
   f. **Read task.md back and verify the checkbox is now `- [x]`** — if not, fix it before proceeding
6. If a task is blocked:
   - Stop and report the blocker
   - Suggest using `debugging` if it's a bug
   - Ask the user for direction
7. After all tasks complete, suggest `/code-review` or `/ship`

## Plan-Task Consistency Check

Before starting execution, verify that plan.md and task.md are aligned:

1. Extract all task numbers from task.md (e.g., `1.1`, `1.2`, `2.1` from `- [ ] X.Y ...`)
2. Extract all task numbers from plan.md (e.g., `1.1`, `1.2`, `2.1` from `### X.Y ...`)
3. Both sets must be identical

**If mismatch found:**
- Tasks in task.md but not in plan.md → orphan checkboxes, no implementation reference
- Tasks in plan.md but not in task.md → work with no status tracking
- Stop and report: "Plan-task mismatch: [details]. Run `/plan` to regenerate, or fix manually."

**If consistent:** proceed with execution.

## IRON RULE: No Checkbox Update = Task Not Done

**A task is NOT complete until its checkbox in task.md is updated AND verified.**

After every single task, you MUST do these two actions in order:
1. **Edit** task.md: change the task's `- [ ]` to `- [x]`
2. **Read** task.md: confirm the checkbox now shows `- [x]`

Only after both steps are done may you proceed to the next task.

This applies to ALL tasks — trivial, XS, batched, skipped-review, no exceptions.

## Task Execution Rules

- One task at a time
- Verify tests pass after each task
- Commit after each verified increment
- Edit + Read task.md after each commit (IRON RULE)
- Don't skip tasks or mark them complete without verification
- If scope expands, stop and update the plan AND task.md together
