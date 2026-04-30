Resume work from a breakpoint.

## Instructions

1. Scan `docs/rein/tasks/` for task files with unchecked items
2. For each task file, read checkbox status
3. Find the first unchecked task
4. **Verify plan-task consistency** — extract task numbers from both files, confirm they match. If mismatch, report and suggest `/plan` to regenerate.
5. Verify task.md status matches reality (all `[x]` items correspond to actually completed work)
6. Determine which phase the work is in:

| Condition | Phase | Skill to invoke |
|-----------|-------|----------------|
| No spec in `docs/rein/specs/` | DEFINE | refine |
| No plan in `docs/rein/plans/` | PLAN | planning |
| No tasks in `docs/rein/tasks/` | PLAN | planning |
| Tasks unchecked in task file | BUILD | incremental + tdd |
| All tasks checked, no review | REVIEW | code-review |
| Review done, not committed | SHIP | git-workflow |

7. Invoke the appropriate skill and continue from that point
8. During BUILD phase: after each completed task, MUST update task.md checkbox (IRON RULE)
9. If no active changes found, suggest starting with `/triage`

## Plan-Task Consistency Check

Same as `/do`: extract task numbers from task.md checkboxes and plan.md `### X.Y` headings. Both sets must match. If they don't, stop and report the mismatch before resuming work.

## Output

```
Found active work: docs/rein/tasks/YYYY-MM-DD-<name>-task.md
Phase: <DEFINE|PLAN|BUILD|REVIEW|SHIP>
Next step: <description>
Resuming with <skill-name>...
```
