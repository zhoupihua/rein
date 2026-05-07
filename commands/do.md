Execute tasks from task.md incrementally. Replaces /opsx:apply.

## Execution Loop

You MUST follow this exact loop. Each iteration starts by reading task.md fresh — you cannot work from memory or cached state.

```
LOOP:
  1. Read task.md — find the FIRST line matching `- [ ]`
  2. If no `- [ ]` found → all tasks complete → invoke `integration-testing` skill → then suggest `/code-review` or `/ship`
  3. Found task X.Y → read plan.md for X.Y details
  4. Implement task X.Y
  5. Verify tests pass
  6. Commit with descriptive message
  7. Edit task.md: change `- [ ] X.Y` to `- [x] X.Y`
  8. Read task.md: confirm the checkbox now shows `- [x]`
  9. GO TO STEP 1 (re-read the file — do NOT proceed from memory)
```

**Why this loop enforces checkbox updates:** Step 1 finds the next task by scanning for `- [ ]`. If step 7-8 was skipped, step 1 will find the same task again on the next iteration. The ONLY way to advance is to update the checkbox.

## Plan-Task Consistency Check

Before the first iteration, verify plan.md and task.md are aligned:
1. Extract task numbers from task.md (e.g., `1.1`, `2.1` from `- [ ] X.Y ...`)
2. Extract task numbers from plan.md (e.g., `1.1`, `2.1` from `### X.Y ...`)
3. Both sets must be identical. If mismatch, stop and report.

## If a Task Is Blocked

- Stop and report the blocker
- Suggest using `debugging` if it's a bug
- Ask the user for direction
- Do NOT mark the task as complete

## Task Execution Rules

- One task at a time
- Re-read task.md at the start of each iteration (not from cache)
- Verify tests pass after each task
- Commit after each verified increment
- Update checkbox + read back after each commit (enforced by the loop)
- If scope expands, stop and update the plan AND task.md together
