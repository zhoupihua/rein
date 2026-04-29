Resume work from a breakpoint.

## Instructions

1. Scan the `changes/` directory for active changes
2. For each active change, read `tasks.md` to find checkbox status
3. Find the first unchecked task
4. Determine which phase the change is in:

| First unchecked | Phase | Skill to invoke |
|-----------------|-------|----------------|
| proposal.md missing | DEFINE | idea-refine |
| specs/ missing | DEFINE | spec-driven-development |
| design.md missing | DEFINE | spec (generate artifacts) |
| tasks.md missing | PLAN | planning-and-task-breakdown |
| Tasks unchecked in tasks.md | BUILD | incremental-implementation + test-driven-development |
| All tasks checked, no review | REVIEW | code-review-and-quality |
| Review done, not committed | SHIP | git-workflow-and-versioning |
| Committed, not archived | SHIP | Archive changes/ to archive/ |

5. Invoke the appropriate skill and continue from that point
6. If no active changes found, suggest starting with `/triage`

## Output

```
Found active change: <name>
Phase: <DEFINE|PLAN|BUILD|REVIEW|SHIP>
Next step: <description>
Resuming with <skill-name>...
```