Run TDD workflow with browser testing support.

## Instructions

1. Invoke `test-driven-development` skill
2. For new features: Write failing tests first (RED)
3. For bug fixes: Use the Prove-It pattern — write a test reproducing the bug
4. Implement minimal code to pass (GREEN)
5. Refactor while keeping tests green (REFACTOR)
6. For browser-based changes: Also invoke `browser-testing-with-devtools` for runtime verification

## Test Pyramid Guidance

- **Unit tests (~80%):** Pure logic, no I/O, milliseconds
- **Integration tests (~15%):** API boundaries, component interactions
- **E2E tests (~5%):** Critical user flows only

## Verification

After completing:
- [ ] All new behaviors have corresponding tests
- [ ] All tests pass
- [ ] Bug fixes include reproduction tests
- [ ] No tests skipped or disabled