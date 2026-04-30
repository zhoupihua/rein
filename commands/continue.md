Resume work from a breakpoint.

## Instructions

1. Scan `docs/rein/tasks/` for task files with unchecked items
2. For each task file, read checkbox status
3. Find the first unchecked task
4. Verify tasks.md and plan.md are in sync (all `[x]` items match completed work)
5. Determine which phase the work is in:

| Condition | Phase | Skill to invoke |
|-----------|-------|----------------|
| No spec in `docs/rein/specs/` | DEFINE | refine |
| No plan in `docs/rein/plans/` | PLAN | planning |
| No tasks in `docs/rein/tasks/` | PLAN | planning |
| Tasks unchecked in tasks file | BUILD | incremental + tdd |
| All tasks checked, no review | REVIEW | code-review |
| Review done, not committed | SHIP | git-workflow |

5. Invoke the appropriate skill and continue from that point
6. During BUILD phase: after each completed task, MUST update tasks.md checkbox (IRON RULE)
7. If no active changes found, suggest starting with `/triage`

## Output

```
Found active work: docs/rein/tasks/YYYY-MM-DD-<name>-task.md
Phase: <DEFINE|PLAN|BUILD|REVIEW|SHIP>
Next step: <description>
Resuming with <skill-name>...
```