---
name: integration-testing
description: Use after TDD implementation is complete, before code review or shipping. Runs spec traceability, integration tests, coverage analysis, and regression checks to verify testing sufficiency beyond unit tests.
---

# Integration Testing

## Overview

TDD verifies each unit in isolation. This skill verifies the system works as a whole.

**When to invoke:** After all TDD tasks are complete (each task has passing unit tests), before code review or shipping.

**Core principle:** Unit tests passing ≠ system working. Integration gaps are where bugs hide.

## The Four Gates

```
All TDD tasks done
       │
       ▼
  1. Spec Traceability ──→ Gaps? ──→ Back to implement missing tests
       │
       │ No gaps
       ▼
  2. Integration Tests ──→ Failures? ──→ Fix interfaces
       │
       │ Pass
       ▼
  3. Coverage Analysis ──→ Below threshold? ──→ Add tests
       │
       │ Meets threshold
       ▼
  4. Regression Check ──→ Failures? ──→ Investigate root cause
       │
       │ All pass
       ▼
  Integration complete
```

You must pass through all four gates in order. A failure at any gate stops the process.

---

## Gate 1: Spec Traceability

Before checking how well things are tested, check whether you're testing the right things at all.

### Build the Traceability Matrix

1. Read `spec.md` — extract every `WHEN/THEN` scenario
2. Read `plan.md` — extract every task's `**Acceptance:**` and `**Verification:**`
3. For each scenario and acceptance criterion, search the codebase for a test that verifies it
4. Produce the matrix:

```markdown
| Spec Scenario | Test Location | Status |
|---------------|---------------|--------|
| WHEN user submits form with valid data THEN record is created | test/api/handlers_test.go:42 | Covered |
| WHEN task checkbox is toggled THEN task.md is updated | (none found) | GAP |
| WHEN invalid input THEN 422 returned | test/api/handlers_test.go:67 | Covered |
```

### Rules

- **GAP = blocker.** Do not proceed past this gate with uncovered spec scenarios.
- A scenario is "Covered" only if a test asserts the THEN clause, not just the WHEN clause.
- If a scenario is intentionally out of scope, mark it `DEFERRED` with a reason. Deferred scenarios still need tracking.

### Why This Matters

TDD drives tests from implementation perspective ("what does this function do?"). Spec traceability drives tests from requirements perspective ("what should the system do?"). Both perspectives are needed. A function can have 100% coverage and still miss a requirement.

---

## Gate 2: Integration Tests

Unit tests verify components in isolation. Integration tests verify they work together.

### What Needs Integration Testing

```
Change touches single function/package     → Unit tests sufficient (skip)
Change crosses package boundaries          → Integration test required
Change crosses API/service boundaries      → API test required
Change affects end-to-end user flow        → E2E smoke test required
```

### Identify Integration Points

1. Read `git diff` to find which packages/modules were changed
2. For each changed package, find its callers and callees
3. Map the cross-boundary interfaces:

```
Package A ──calls──→ Package B ──calls──→ Package C
   │                    │
   │ changed            │ changed
   │                    │
   └── need test: A→B interface
                        └── need test: B→C interface
```

4. For each interface, write or verify a test that exercises the contract

### Language-Specific Guidance

**Go:**
```go
// Integration test: verify CLI command end-to-end
func TestValidateCommand(t *testing.T) {
    cmd := exec.Command("rein", "validate")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("validate failed: %s", output)
    }
}

// Or use testscript for CLI integration tests
// https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript
```

**Node/TypeScript:**
```typescript
// API integration test with supertest
import request from 'supertest';
import { app } from '../app';

it('creates a task and returns 201', async () => {
  const res = await request(app)
    .post('/api/tasks')
    .send({ title: 'Test task' });
  expect(res.status).toBe(201);
  expect(res.body.id).toBeDefined();
});
```

**Python:**
```python
# API integration test with test client
def test_create_task(client):
    response = client.post('/api/tasks', json={'title': 'Test task'})
    assert response.status_code == 201
    assert 'id' in response.json
```

### Vertical Slice Integration

For each vertical slice (e.g., "Create a task" = DB + API + UI):

1. Test the full path: input at the top → correct state at the bottom
2. Verify data flows correctly through each layer
3. Verify error handling at each boundary

### Integration Test Rules

- Test the contract, not the implementation — assert on inputs and outputs, not internal state
- Use real dependencies where possible (real DB with test data, real filesystem with temp dirs)
- Mock only external services you can't control (third-party APIs, payment processors)
- Each integration test should be independent — no shared mutable state
- Integration tests should be runnable in CI without manual setup

---

## Gate 3: Coverage Analysis

Numbers aren't everything, but they catch what you forgot.

### Run Coverage

```bash
# Go
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

# Node/TypeScript
nyc --reporter=text --reporter=lcov npm test

# Python
pytest --cov --cov-report=term-missing
```

### Thresholds

| Scope | Minimum Coverage | Rationale |
|-------|-----------------|-----------|
| New/changed code | **80%** | Baseline quality floor |
| Core paths (auth, data handling, error paths) | **100%** | Failure here = production incident |
| Utility/helper code | 60% | Lower risk, but still needs basic coverage |
| Test files themselves | N/A | Not measured |

### How to Evaluate Coverage Gaps

Coverage percentage alone is misleading. When you find uncovered code:

1. **Is it a core path?** → Must add test. No exceptions.
2. **Is it an error handler?** → Must add test. Error paths fail in production.
3. **Is it a simple getter/setter?** → Low priority. Note it, move on.
4. **Is it dead code?** → Delete it, don't test it.
5. **Is it a complex branch with no test?** → Add test. Untested branches hide bugs.

### Coverage Anti-Patterns

| Anti-Pattern | Problem | Fix |
|-------------|---------|-----|
| Testing implementation details to hit 100% | Refactor breaks tests even if behavior unchanged | Test behavior, not implementation |
| Ignoring uncovered error paths | "That'll never happen" → it happens in production | Every error path needs a test |
| Skipping coverage for "simple" code | Simple code has simple bugs | Simple tests take 30 seconds |
| Adding assertions that always pass | Test that never fails is useless | Verify test fails when code is broken |

---

## Gate 4: Regression Check

After all new tests pass, verify you didn't break anything that was working before.

### Full Suite Run

```bash
# Run the ENTIRE test suite, not just affected packages
go test ./...
npm test
pytest
```

### Affected Module Deep Test

Based on `git diff --stat`, identify packages that import changed packages. Run their tests too:

```bash
# Go: find packages that import changed package
go list ./... | xargs go test

# Or use go test graph to find dependents
```

### Smoke Test Checklist

Critical paths that must always work — never skip these:

- [ ] Application starts without errors
- [ ] Main user flow works end-to-end
- [ ] Authentication/authorization works
- [ ] Data can be created, read, updated, deleted
- [ ] Error responses are correct (not 500s on valid input)

### Skipped/Disabled Tests

Check for tests that were skipped or disabled:

```bash
# Go: find skipped tests
grep -r "t.Skip\|t.Skipf" --include="*_test.go"

# Node: find skipped tests
grep -r "it.skip\|describe.skip\|xit\|xdescribe" --include="*.test.*"
```

**Rule:** Every skipped test must have a comment explaining why and a tracking issue. If you introduced a skip, it's your responsibility to fix it before shipping.

---

## Common Rationalizations

| Excuse | Reality |
|--------|---------|
| "Unit tests passed, integration will be fine" | Integration bugs are the ones unit tests can't catch — that's the point |
| "It's too complex to integration test" | If it's too complex to test, it's too complex to ship |
| "We'll add integration tests later" | Later = never. Write them now while context is fresh |
| "Coverage is just a number" | Yes, and a number below 80% means you're not testing enough |
| "The spec is out of date" | Then update the spec. Don't use stale spec as an excuse to skip traceability |
| "No time for regression testing" | Shipping a regression costs 10x the time of running the suite |
| "These are just small changes" | Small changes cause big regressions. Especially small changes |

## Red Flags - STOP

- Skipping the spec traceability gate
- Running only affected package tests (not full suite)
- Coverage below 50% on new code
- Any test marked skip/disable without a tracking issue
- "We'll do integration testing after the code review"
- Changed cross-package interfaces without integration tests
- Merging without all four gates passing

## Verification Checklist

- [ ] Every spec scenario has a corresponding test (GATE 1)
- [ ] Cross-component interfaces have integration tests (GATE 2)
- [ ] Coverage meets thresholds: new code >= 80%, core paths 100% (GATE 3)
- [ ] Full test suite passes (GATE 4)
- [ ] No skipped/disabled tests without tracking issues (GATE 4)
- [ ] Smoke test checklist is green (GATE 4)

All gates must pass before proceeding to code review or ship.
