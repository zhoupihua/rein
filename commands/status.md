Show task progress and detect stale checkboxes.

## Instructions

1. Scan `docs/rein/tasks/` for task files
2. For each task file:
   a. Count total tasks, checked (`- [x]`), and unchecked (`- [ ]`)
   b. List all unchecked tasks with task numbers and descriptions
3. Check for drift:
   a. Read task.md last modification time
   b. Run `git log --oneline -10` to see recent commits
   c. If commits exist after task.md was last modified, warn that checkboxes may be stale
4. If `--fix` argument provided:
   a. For each unchecked task, ask: "Is task X.Y [description] complete? (y/n/skip)"
   b. If yes, edit task.md: change `- [ ] X.Y` to `- [x] X.Y`
   c. After all fixes, read task.md back to confirm
5. Output summary

## Output

```
Task Progress: docs/rein/tasks/YYYY-MM-DD-<name>-task.md
  Completed: 3/8  Remaining: 5

  Remaining tasks:
  - [ ] 2.1 Add API endpoint for X
  - [ ] 2.2 Build UI for X
  - [ ] 3.1 Add error handling
  - [ ] 3.2 Update documentation
  - [ ] 4.1 End-to-end test

  Drift: 2 commits after last task.md update
  → Run /status --fix to update checkboxes interactively
```

## Drift Detection

If the task.md file was last modified before the most recent commit, it's likely that some tasks were completed but their checkboxes weren't updated. In this case:
- Warn the user about potential stale state
- Suggest running `/status --fix` to interactively update checkboxes
- Do NOT auto-fix without user confirmation
