---
description: Run Ralph loop refactoring — identify smells, characterize behavior, small refactor, verify, commit, repeat
---

Structured refactoring using the Ralph Johnson loop. Improves code structure without changing behavior.

## Instructions

1. Invoke `refactor` skill
2. Follow the Ralph loop:
   - **Identify smell** — scan target code against the smell catalog, pick one
   - **Characterize behavior** — if tests don't cover the code, write characterization tests first
   - **Small refactor** — one smell, one technique, ≤30 lines changed
   - **Verify** — run tests, all must pass; if fail, revert immediately
   - **Commit** — one smell fixed per commit
   - **Loop** — go back to identify next smell
3. Exit when no more worthwhile smells remain
4. Run full test suite to confirm no regressions

## Scope

If invoked without arguments, refactor the most recently changed code. If a file or directory is specified, refactor that target.

## Output

```
Ralph loop starting on <target>

Smell 1: <smell name> → <refactoring technique>
  Tests: ✅ pass after refactor
  Commit: refactor: <technique> — <smell eliminated>

Smell 2: <smell name> → <refactoring technique>
  Tests: ✅ pass after refactor
  Commit: refactor: <technique> — <smell eliminated>

...

Ralph loop complete. N smells eliminated, M remain (not worth fixing).
Full test suite: ✅ all pass
```
