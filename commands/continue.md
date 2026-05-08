---
description: Resume work from a breakpoint — detect phase, verify consistency, continue
---

Resume work from a breakpoint.

## Instructions

1. Scan `docs/rein/changes/` for feature directories with artifacts
2. For each feature, read artifact status (which files exist, task checkbox progress)
3. Find the first incomplete phase
4. **Verify plan-task consistency** — extract task numbers from both files, confirm they match. If mismatch, report and suggest `/plan` to regenerate.
5. Verify task.md status matches reality (all `[x]` items correspond to actually completed work)
6. Determine which phase the work is in:

| Condition | Phase | Skill to invoke |
|-----------|-------|----------------|
| No proposal.md and no spec.md | DEFINE | define |
| Has proposal.md, no spec.md | DEFINE | define (resume from Step 2b) |
| No plan.md in feature dir | PLAN | planning |
| No task.md in feature dir | PLAN | planning |
| Tasks unchecked in task.md | BUILD | executing-plans + tdd |
| All tasks checked, no review.md | REVIEW | code-review |
| Review done, not committed | SHIP | git-workflow |

7. Invoke the appropriate skill and continue from that point
8. During BUILD phase: after each completed task, MUST update task.md checkbox (IRON RULE)
9. During BUILD phase: start the task progress watchdog (same as `/do` Setup section) if not already running
10. If no active changes found, suggest starting with `/quick` or `/feature`

## Plan-Task Consistency Check

Same as `/do`: extract task numbers from task.md checkboxes and plan.md `### X.Y` headings. Both sets must match. If they don't, stop and report the mismatch before resuming work.

## Output

```
Found active work: docs/rein/changes/<name>/
Phase: <DEFINE|PLAN|BUILD|REVIEW|SHIP>
Next step: <description>
Resuming with <skill-name>...
```
