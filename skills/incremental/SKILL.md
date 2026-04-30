---
name: incremental
description: Use when implementing any feature or change from a plan. Build in thin vertical slices, test each increment, commit frequently. Use when you have a written implementation plan to execute.
---

# Incremental Implementation

## Overview

Build in thin vertical slices — implement one piece, test it, verify it, then expand. Each increment should leave the system in a working, testable state. This is the execution discipline that makes large features manageable.

**Announce at start:** "I'm using the incremental skill to execute this plan."

## Pre-Execution: Load and Review Plan

1. Read `docs/rein/plans/YYYY-MM-DD-<name>.md` for architecture context and decisions
2. Read `docs/rein/tasks/YYYY-MM-DD-<name>-tasks.md` for the ordered task list
3. Review critically — identify any questions or concerns about the plan
4. If concerns: Raise them with your human partner before starting
5. If no concerns: Proceed with the first unchecked task

## The Increment Cycle

```
┌──────────────────────────────────────┐
│                                      │
│   Implement ──→ Test ──→ Verify ──┐  │
│       ▲                           │  │
│       └───── Commit ◄─────────────┘  │
│              │                       │
│              ▼                       │
│          Next slice                  │
│                                      │
└──────────────────────────────────────┘
```

For each slice:

1. **Implement** the smallest complete piece of functionality
2. **Test** — run the test suite (or write a test if none exists)
3. **Verify** — confirm the slice works as expected (tests pass, build succeeds, manual check)
4. **Commit** — save your progress with a descriptive message
5. **Update task status** — mark the completed task in tasks.md: change `- [ ]` to `- [x]`
6. **Move to the next slice** — carry forward, don't restart

## Slicing Strategies

### Vertical Slices (Preferred)

Build one complete path through the stack:

```
Slice 1: Create a task (DB + API + basic UI)
    → Tests pass, user can create a task via the UI

Slice 2: List tasks (query + API + UI)
    → Tests pass, user can see their tasks

Slice 3: Edit a task (update + API + UI)
    → Tests pass, user can modify tasks

Slice 4: Delete a task (delete + API + UI + confirmation)
    → Tests pass, full CRUD complete
```

### Contract-First Slicing

When backend and frontend need to develop in parallel:

```
Slice 0: Define the API contract (types, interfaces, OpenAPI spec)
Slice 1a: Implement backend against the contract + API tests
Slice 1b: Implement frontend against mock data matching the contract
Slice 2: Integrate and test end-to-end
```

### Risk-First Slicing

Tackle the riskiest or most uncertain piece first:

```
Slice 1: Prove the WebSocket connection works (highest risk)
Slice 2: Build real-time task updates on the proven connection
Slice 3: Add offline support and reconnection
```

## Implementation Rules

### Rule 0: Simplicity First

Before writing any code, ask: "What is the simplest thing that could work?"

After writing code, review it against these checks:
- Can this be done in fewer lines?
- Are these abstractions earning their complexity?
- Would a staff engineer look at this and say "why didn't you just..."?
- Am I building for hypothetical future requirements, or the current task?

### Rule 0.5: Scope Discipline

Touch only what the task requires.

Do NOT:
- "Clean up" code adjacent to your change
- Refactor imports in files you're not modifying
- Remove comments you don't fully understand
- Add features not in the spec because they "seem useful"
- Modernize syntax in files you're only reading

If you notice something worth improving outside your task scope, note it — don't fix it.

### Rule 1: One Thing at a Time

Each increment changes one logical thing. Don't mix concerns.

### Rule 2: Keep It Compilable

After each increment, the project must build and existing tests must pass. Don't leave the codebase in a broken state between slices.

### Rule 3: Feature Flags for Incomplete Features

If a feature isn't ready for users but you need to merge increments:

```typescript
const ENABLE_TASK_SHARING = process.env.FEATURE_TASK_SHARING === 'true';

if (ENABLE_TASK_SHARING) {
  // New sharing UI
}
```

### Rule 4: Safe Defaults

New code should default to safe, conservative behavior.

### Rule 5: Rollback-Friendly

Each increment should be independently revertable:
- Additive changes (new files, new functions) are easy to revert
- Modifications to existing code should be minimal and focused
- Database migrations should have corresponding rollback migrations
- Avoid deleting something in one commit and replacing it in the same commit

## When to Stop and Ask for Help

**STOP executing immediately when:**
- Hit a blocker (missing dependency, test fails, instruction unclear)
- Plan has critical gaps preventing starting
- You don't understand an instruction
- Verification fails repeatedly

**Ask for clarification rather than guessing.**

## Task Status Tracking (MANDATORY)

**IRON RULE: A task is NOT complete until its checkbox in tasks.md is updated AND verified.** Moving to the next task without updating tasks.md is a process violation.

During execution, use **plan.md** as the implementation reference and **tasks.md** for status tracking:

- **plan.md** → HOW: Read task details (acceptance criteria, verification, files, notes) before implementing each task
- **tasks.md** → STATUS: After completing each task, you MUST update the checkbox

After completing each task increment, execute this two-step sequence — **both steps are required**:

1. **Edit** `docs/rein/tasks/YYYY-MM-DD-<name>-tasks.md` — change the task's `- [ ]` to `- [x]`
2. **Read** the same file back — confirm the checkbox now shows `- [x]`. If not, fix it immediately.

Only after both steps are done may you proceed to the next task.

**This applies even when:**
- The task was trivial or XS-sized
- You plan to batch multiple tasks (update each one as it completes)
- You are executing quickly and want to skip the step
- You skipped reviews for a small task (still must update tasks.md)

The tasks.md checkbox state is the **single source of truth** for progress — `/continue` relies on it to determine resume points and current phase.

## After All Tasks Complete

After all tasks are complete and verified:
- Invoke **git-workflow** skill to complete this work
- Follow that skill to verify tests, present options, execute choice

## Common Rationalizations

| Rationalization | Reality |
|---|---|
| "I'll test it all at the end" | Bugs compound. A bug in Slice 1 makes Slices 2-5 wrong. Test each slice. |
| "It's faster to do it all at once" | It feels faster until something breaks and you can't find which changed line caused it. |
| "These changes are too small to commit separately" | Small commits are free. Large commits hide bugs and make rollbacks painful. |
| "I'll add the feature flag later" | If the feature isn't complete, it shouldn't be user-visible. Add the flag now. |
| "This refactor is small enough to include" | Refactors mixed with features make both harder to review and debug. Separate them. |

## Red Flags

- Moving to the next task without updating tasks.md checkbox (IRON RULE violation)
- More than 100 lines of code written without running tests
- Multiple unrelated changes in a single increment
- "Let me just quickly add this too" scope expansion
- Skipping the test/verify step to move faster
- Build or tests broken between increments
- Large uncommitted changes accumulating
- Building abstractions before the third use case demands it
- Touching files outside the task scope "while I'm here"

## Verification

After completing all increments for a task:

- [ ] Each increment was individually tested and committed
- [ ] The full test suite passes
- [ ] The build is clean
- [ ] The feature works end-to-end as specified
- [ ] No uncommitted changes remain
