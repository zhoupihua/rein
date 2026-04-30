L2 standard change (1-3 files, clear requirements).

## Bug Fix Flow
1. Invoke `debugging` skill — reproduce, localize, identify root cause
2. Invoke `tdd` skill — Prove-It pattern: write failing test that reproduces the bug
3. Implement the fix (minimal change to address root cause)
4. Invoke `verify` skill — verify all tests pass, build succeeds
5. Commit with `fix:` message

## Feature Flow
1. Invoke `tdd` skill — write failing test for new behavior
2. Implement minimal code to pass the test
3. Refactor if needed (keep tests green)
4. Invoke `verify` skill — verify all tests pass, build succeeds
5. Commit with `feat:` message

## Frontend Bug
Same as Bug Fix Flow, but also invoke `browser-testing` for runtime verification.

**If the scope grows beyond 3 files or requirements become unclear, stop and use `/feature` instead.**