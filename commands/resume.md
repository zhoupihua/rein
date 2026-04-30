Resume work from a breakpoint.

## Instructions

1. Scan `docs/rein/tasks/` for task files with unchecked items
2. For each task file, read checkbox status
3. Find the first unchecked task
4. Determine which phase the work is in:

| Condition | Phase | Skill to invoke |
|-----------|-------|----------------|
| No spec in `docs/rein/specs/` | DEFINE | idea-refine |
| No plan in `docs/rein/plans/` | PLAN | planning-and-task-breakdown |
| No tasks in `docs/rein/tasks/` | PLAN | planning-and-task-breakdown |
| Tasks unchecked in tasks file | BUILD | incremental-implementation + test-driven-development |
| All tasks checked, no review | REVIEW | code-review-and-quality |
| Review done, not committed | SHIP | git-workflow-and-versioning |

5. Invoke the appropriate skill and continue from that point
6. If no active changes found, suggest starting with `/triage`

## Output

```
Found active work: docs/rein/tasks/YYYY-MM-DD-<name>-tasks.md
Phase: <DEFINE|PLAN|BUILD|REVIEW|SHIP>
Next step: <description>
Resuming with <skill-name>...
```