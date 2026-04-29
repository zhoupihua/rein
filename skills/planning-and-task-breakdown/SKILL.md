---
name: planning-and-task-breakdown
description: Use when you have a spec or requirements for a multi-step task, before touching code. Breaks work into ordered tasks with dependency graphs, vertical slicing, and explicit acceptance criteria.
---

# Planning and Task Breakdown

## Overview

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Decompose work into small, verifiable tasks with explicit acceptance criteria. Document everything: which files to touch, code, testing, docs. Bite-sized tasks. DRY. YAGNI. TDD. Frequent commits.

**Announce at start:** "I'm using the planning-and-task-breakdown skill to create the implementation plan."

**Save plans to:** `docs/plans/YYYY-MM-DD-<feature-name>.md`
- (User preferences for plan location override this default)

## Scope Check

If the spec covers multiple independent subsystems, suggest breaking into separate plans — one per subsystem. Each plan should produce working, testable software on its own.

## Planning Process

### Step 1: Enter Plan Mode

Before writing any code, operate in read-only mode:

- Read the spec and relevant codebase sections
- Identify existing patterns and conventions
- Map dependencies between components
- Note risks and unknowns

**Do NOT write code during planning.** The output is a plan document, not implementation.

### Step 2: Identify the Dependency Graph

Map what depends on what:

```
Database schema
    │
    ├── API models/types
    │       │
    │       ├── API endpoints
    │       │       │
    │       │       └── Frontend API client
    │       │               │
    │       │               └── UI components
    │       │
    │       └── Validation logic
    │
    └── Seed data / migrations
```

Implementation order follows the dependency graph bottom-up: build foundations first.

### Step 3: Slice Vertically

Instead of building all the database, then all the API, then all the UI — build one complete feature path at a time:

**Bad (horizontal slicing):**
```
Task 1: Build entire database schema
Task 2: Build all API endpoints
Task 3: Build all UI components
Task 4: Connect everything
```

**Good (vertical slicing):**
```
Task 1: User can create an account (schema + API + UI for registration)
Task 2: User can log in (auth schema + API + UI for login)
Task 3: User can create a task (task schema + API + UI for creation)
Task 4: User can view task list (query + API + UI for list view)
```

Each vertical slice delivers working, testable functionality.

### Step 4: Map File Structure

Before defining tasks, map out which files will be created or modified and what each one is responsible for:

- Design units with clear boundaries and well-defined interfaces
- Prefer smaller, focused files over large ones
- Files that change together should live together
- In existing codebases, follow established patterns

### Step 5: Write Tasks

**Bite-Sized Task Granularity** — Each step is one action (2-5 minutes):
- "Write the failing test" - step
- "Run it to make sure it fails" - step
- "Implement the minimal code to make the test pass" - step
- "Run the tests and make sure they pass" - step
- "Commit" - step

Each task follows this structure:

```markdown
## Task [N]: [Short descriptive title]

**Description:** One paragraph explaining what this task accomplishes.

**Acceptance criteria:**
- [ ] [Specific, testable condition]
- [ ] [Specific, testable condition]

**Verification:**
- [ ] Tests pass: `npm test -- --grep "feature-name"`
- [ ] Build succeeds: `npm run build`

**Dependencies:** [Task numbers this depends on, or "None"]

**Files likely touched:**
- `src/path/to/file.ts`
- `tests/path/to/test.ts`

**Estimated scope:** [XS | S | M | L]
```

### Step 6: Order and Checkpoint

Arrange tasks so that:
1. Dependencies are satisfied (build foundation first)
2. Each task leaves the system in a working state
3. Verification checkpoints occur after every 2-3 tasks
4. High-risk tasks are early (fail fast)

## Plan Document Template

```markdown
# [Feature Name] Implementation Plan

> **For agentic workers:** Use subagent-driven-development (recommended) or incremental-implementation to execute this plan task-by-task.

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

---

## Architecture Decisions
- [Key decision 1 and rationale]
- [Key decision 2 and rationale]

## Task List

### Phase 1: Foundation
- [ ] Task 1: ...
- [ ] Task 2: ...

### Checkpoint: Foundation
- [ ] Tests pass, builds clean

### Phase 2: Core Features
- [ ] Task 3: ...
- [ ] Task 4: ...

### Checkpoint: Core Features
- [ ] End-to-end flow works

### Phase 3: Polish
- [ ] Task 5: ...

### Checkpoint: Complete
- [ ] All acceptance criteria met

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk] | [High/Med/Low] | [Strategy] |

## Open Questions
- [Question needing human input]
```

## Task Sizing Guidelines

| Size | Files | Scope | Example |
|------|-------|-------|---------|
| **XS** | 1 | Single function or config change | Add a validation rule |
| **S** | 1-2 | One component or endpoint | Add a new API endpoint |
| **M** | 3-5 | One feature slice | User registration flow |
| **L** | 5-8 | Multi-component feature | Search with filtering and pagination |
| **XL** | 8+ | **Too large — break it down further** | — |

If a task is L or larger, it should be broken into smaller tasks. An agent performs best on S and M tasks.

**When to break a task down further:**
- It would take more than one focused session
- You cannot describe the acceptance criteria in 3 or fewer bullet points
- It touches two or more independent subsystems
- You find yourself writing "and" in the task title

## Parallelization Opportunities

| Category | Examples | Strategy |
|----------|----------|----------|
| **Safe to parallelize** | Independent feature slices, tests for implemented features, docs | Dispatch concurrently |
| **Must be sequential** | Database migrations, shared state changes, dependency chains | Execute in order |
| **Needs coordination** | Features sharing an API contract | Define contract first, then parallelize |

## No Placeholders

Every step must contain the actual content an engineer needs. These are **plan failures**:
- "TBD", "TODO", "implement later", "fill in details"
- "Add appropriate error handling" / "add validation" / "handle edge cases"
- "Write tests for the above" (without actual test code)
- "Similar to Task N" (repeat the code — the engineer may be reading tasks out of order)
- Steps that describe what to do without showing how (code blocks required for code steps)

## Self-Review

After writing the complete plan:

1. **Spec coverage:** Can you point to a task that implements each spec requirement? List any gaps.
2. **Placeholder scan:** Search for red flags from the "No Placeholders" section. Fix them.
3. **Type consistency:** Do types, method signatures, and property names match across tasks?

Fix any issues inline. If you find a spec requirement with no task, add the task.

## Execution Handoff

After saving the plan, offer execution choice:

**"Plan complete and saved to `docs/plans/<filename>.md`. Two execution options:**

**1. Subagent-Driven (recommended)** — Fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** — Execute tasks in this session using incremental-implementation, batch execution with checkpoints

**Which approach?"**

## Common Rationalizations

| Rationalization | Reality |
|---|---|
| "I'll figure it out as I go" | That's how you end up with rework. 10 minutes of planning saves hours. |
| "The tasks are obvious" | Write them down anyway. Explicit tasks surface hidden dependencies. |
| "Planning is overhead" | Planning is the task. Implementation without a plan is just typing. |
| "I can hold it all in my head" | Context windows are finite. Written plans survive session boundaries. |

## Red Flags

- Starting implementation without a written task list
- Tasks that say "implement the feature" without acceptance criteria
- No verification steps in the plan
- All tasks are XL-sized
- No checkpoints between tasks
- Dependency order isn't considered
- Placeholders or vague steps in the plan

## Verification

Before starting implementation, confirm:

- [ ] Every task has acceptance criteria
- [ ] Every task has a verification step
- [ ] Task dependencies are identified and ordered correctly
- [ ] No task touches more than ~5 files
- [ ] Checkpoints exist between major phases
- [ ] The human has reviewed and approved the plan
