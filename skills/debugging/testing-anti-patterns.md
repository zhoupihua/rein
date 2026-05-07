# Testing Anti-Patterns

**Load this reference when:** writing or changing tests, adding mocks, or tempted to add test-only methods to production code.

## Overview

Tests must verify real behavior, not mock behavior. Mocks are a means to isolate, not the thing being tested.

**Core principle:** Test what the code does, not what the mocks do.

**Following strict TDD prevents these anti-patterns.**

## The Iron Laws

```
1. NEVER test mock behavior
2. NEVER add test-only methods to production code
3. NEVER mock without understanding dependencies
```

## Anti-Pattern 1: Testing Mock Behavior

**The violation:**

```typescript
// BAD: Testing that the mock exists
test('renders sidebar', () => {
  render(<Page />);
  expect(screen.getByTestId('sidebar-mock')).toBeInTheDocument();
});
```

**Why this is wrong:** You're verifying the mock works, not that the component works. Test passes when mock is present, fails when it's not. Tells you nothing about real behavior.

**The fix:**

```typescript
// GOOD: Test real component or don't mock it
test('renders sidebar', () => {
  render(<Page />);  // Don't mock sidebar
  expect(screen.getByRole('navigation')).toBeInTheDocument();
});
```

### Gate Function

```
BEFORE asserting on any mock element:
  Ask: "Am I testing real component behavior or just mock existence?"

  IF testing mock existence:
    STOP - Delete the assertion or unmock the component

  Test real behavior instead
```

## Anti-Pattern 2: Test-Only Methods in Production

**The violation:**

```go
// BAD: Destroy only used in tests
type Session struct { /* ... */ }

func (s *Session) Destroy() error {  // Looks like production API!
    return s.workspaceManager.DestroyWorkspace(s.id)
}

// In tests
defer session.Destroy()
```

**Why this is wrong:** Production class polluted with test-only code. Dangerous if accidentally called in production.

**The fix:**

```go
// GOOD: Test utilities handle test cleanup
// Session has no Destroy method — it's stateless in production

func cleanupSession(t *testing.T, session *Session) {
    t.Helper()
    workspace := session.GetWorkspaceInfo()
    if workspace != nil {
        workspaceManager.DestroyWorkspace(workspace.ID)
    }
}

// In tests
defer cleanupSession(t, session)
```

## Anti-Pattern 3: Mocking Without Understanding

**The violation:**

```go
// BAD: Mock breaks test logic
func TestDetectDuplicateServer(t *testing.T) {
    // Mock prevents config write that test depends on!
    mockCatalog := &MockToolCatalog{
        DiscoverAndCacheFunc: func() error { return nil },
    }

    AddServer(config, mockCatalog)
    err := AddServer(config, mockCatalog)  // Should fail — but won't!
    // Mock didn't write config, so duplicate check never triggers
}
```

**Why this is wrong:** Mocked method had side effect test depended on. Over-mocking breaks actual behavior.

**The fix:**

```go
// GOOD: Mock at correct level
func TestDetectDuplicateServer(t *testing.T) {
    // Mock the slow part, preserve behavior test needs
    mockServer := &MockMCPServer{}  // Just mock slow server startup

    AddServer(config, mockServer)  // Config written
    err := AddServer(config, mockServer)  // Duplicate detected
}
```

## Anti-Pattern 4: Incomplete Mocks

**The violation:**

```go
// BAD: Partial mock — only fields you think you need
response := &APIResponse{
    Status: "success",
    Data:   UserData{ID: "123", Name: "Alice"},
    // Missing: Metadata that downstream code uses
}
```

**Why this is wrong:** Partial mocks hide structural assumptions. Downstream code may depend on fields you didn't include. Tests pass but integration fails.

**The Iron Rule:** Mock the COMPLETE data structure as it exists in reality, not just fields your immediate test uses.

**The fix:**

```go
// GOOD: Mirror real API completeness
response := &APIResponse{
    Status: "success",
    Data:   UserData{ID: "123", Name: "Alice"},
    Metadata: RequestMeta{
        RequestID: "req-789",
        Timestamp: 1234567890,
    },
    // All fields real API returns
}
```

## Anti-Pattern 5: Integration Tests as Afterthought

```
Implementation complete
No tests written
"Ready for testing"
```

**Why this is wrong:** Testing is part of implementation, not optional follow-up.

**The fix:**

```
TDD cycle:
1. Write failing test
2. Implement to pass
3. Refactor
4. THEN claim complete
```

## Quick Reference

| Anti-Pattern | Fix |
|--------------|-----|
| Assert on mock elements | Test real component or unmock it |
| Test-only methods in production | Move to test utilities |
| Mock without understanding | Understand dependencies first, mock minimally |
| Incomplete mocks | Mirror real API completely |
| Tests as afterthought | TDD — tests first |

## Red Flags

- Assertion checks for `*-mock` test IDs
- Methods only called in test files
- Mock setup is >50% of test
- Test fails when you remove mock
- Can't explain why mock is needed
- Mocking "just to be safe"

## The Bottom Line

**Mocks are tools to isolate, not things to test.**

If TDD reveals you're testing mock behavior, you've gone wrong. Fix: Test real behavior or question why you're mocking at all.
