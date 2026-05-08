---
description: Execute tasks from task.md incrementally — one at a time, test, commit, update checkbox
---

Execute tasks from task.md incrementally. Replaces /opsx:apply.

## Setup: Task Progress Watchdog

Before starting the execution loop, start a periodic watchdog that checks task.md sync:

1. Determine the feature name from `docs/rein/changes/<name>/task.md`
2. Start a CronCreate recurring job (every 5 minutes) with this prompt:

```
Task progress watchdog for <name>. Read docs/rein/changes/<name>/task.md. Count [ ] (unchecked) and [x] (checked) checkboxes. Then:
- If all tasks are checked ([x]): delete this cron job (it's done), then remind the user that all tasks are complete.
- If unchecked tasks exist but code changes suggest they should be done (e.g., files mentioned in task descriptions have been modified, tests for the task exist and pass): warn the user that task.md may be out of sync — list the suspect tasks and suggest checking them off.
- If no sync issues detected: no output needed, just wait for next check.
```

3. Store the cron job ID so you can delete it when all tasks complete (or when the user stops /do)

## Execution Loop

You MUST follow this exact loop. Each iteration starts by reading task.md fresh — you cannot work from memory or cached state.

```
LOOP:
  1. Read task.md — find the FIRST line matching `- [ ]`
  2. If no `- [ ]` found → all tasks complete → see "When All Tasks Complete" section
  3. Found task X.Y → read plan.md for X.Y details
  4. If task has RED/GREEN/REFACTOR sub-tasks, execute them in order:
     a. RED: write failing test → check off sub-task
     b. GREEN: implement to pass → check off sub-task
     c. REFACTOR: clean up → check off sub-task (parent auto-checks when all done)
  5. If task has no sub-tasks, implement task X.Y directly
  6. Verify tests pass
  7. Commit with descriptive message
  8. Edit task.md: change `- [ ] X.Y` to `- [x] X.Y` (if not auto-checked by sub-tasks)
  9. Read task.md: confirm the checkbox now shows `- [x]`
  10. GO TO STEP 1 (re-read the file — do NOT proceed from memory)
```

**Why this loop enforces checkbox updates:** Step 1 finds the next task by scanning for `- [ ]`. If step 8-9 was skipped, step 1 will find the same task again on the next iteration. The ONLY way to advance is to update the checkbox.

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

## When All Tasks Complete

- Delete the watchdog cron job (CronDelete with the stored job ID)
- Invoke `integration-testing` skill
- Then invoke `git-workflow` skill to complete

## Task Execution Rules

- One task at a time
- Re-read task.md at the start of each iteration (not from cache)
- Verify tests pass after each task
- Commit after each verified increment
- Update checkbox + read back after each commit (enforced by the loop)
- If scope expands, stop and update the plan AND task.md together
