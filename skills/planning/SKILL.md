---
name: planning
description: Use when you have a spec or requirements for a multi-step task, before touching code. Breaks work into ordered tasks with dependency graphs, vertical slicing, and explicit acceptance criteria.
---

# Planning and Task Breakdown

## Overview

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Decompose work into small, verifiable tasks with explicit acceptance criteria. Document everything: which files to touch, code, testing, docs. Bite-sized tasks. DRY. YAGNI. TDD. Frequent commits.

**Announce at start:** "I'm using the planning skill to create the implementation plan."

**Save output to:** `docs/rein/plans/YYYY-MM-DD-<feature-name>.md` + `docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md`
- `plans/` — Architecture decisions, dependency graph, slicing strategy, risks (decision layer)
- `tasks/` — Ordered task checklist with acceptance criteria (execution layer, overwrites /spec's coarse tasks)
- If no `docs/rein/plans/` or `docs/rein/tasks/` directory exists, create it

## Scope Check

If the spec covers multiple independent subsystems, suggest breaking into separate plans — one per subsystem. Each plan should produce working, testable software on its own.

## Planning Process

### Step 1: Enter Plan Mode

Before writing any code, operate in read-only mode:

- Read the spec and relevant codebase sections
- Identify existing patterns and conventions
- Map dependencies between components
- Note risks and unknowns

**Do NOT write code during planning.** The output is two documents, not implementation.

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

### Step 6: Order and Checkpoint

Arrange tasks so that:
1. Dependencies are satisfied (build foundation first)
2. Each task leaves the system in a working state
3. Verification checkpoints occur after every 2-3 tasks
4. High-risk tasks are early (fail fast)

## Output: Two Files

### plan.md Template (Decision Layer)

```markdown
# [Feature Name] Plan

> **For agentic workers:** Read this file for architecture context, then read tasks.md for execution.

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

---

## Architecture Decisions
- [Key decision 1 and rationale]
- [Key decision 2 and rationale]

## Dependency Graph

[ASCII tree showing component dependencies]

## Vertical Slicing Strategy

[How work is sliced into independent, testable increments]

## File Map

| File | Purpose | New/Modified |
|------|---------|-------------|
| `src/path/to/file.ts` | [Purpose] | New/Modified |

## Parallelization

| Category | Tasks | Strategy |
|----------|-------|----------|
| Safe to parallelize | [Task numbers] | Dispatch concurrently |
| Must be sequential | [Task numbers] | Execute in order |
| Needs coordination | [Task numbers] | Define contract first |

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk] | [High/Med/Low] | [Strategy] |

## Open Questions
- [Question needing human input]
```

### tasks.md Template (Execution Layer)

```markdown
## 1. Foundation

- [ ] 1.1 [Short descriptive title]
  - Acceptance: [Specific, testable condition]
  - Verification: `npm test -- --grep "feature-name"`
  - Dependencies: None
  - Files: `src/path/to/file.ts`
  - Scope: [XS | S | M | L]

- [ ] 1.2 [Short descriptive title]
  - Acceptance: [Specific, testable condition]
  - Verification: `npm run build`
  - Dependencies: 1.1
  - Files: `src/path/to/file.ts`, `tests/path/to/test.ts`
  - Scope: S

### Checkpoint: Foundation
- [ ] Tests pass, builds clean

## 2. Core Features

- [ ] 2.1 [Short descriptive title]
  - Acceptance: [Specific, testable condition]
  - Verification: `npm test`
  - Dependencies: 1.2
  - Files: `src/path/to/file.ts`
  - Scope: M

### Checkpoint: Core Features
- [ ] End-to-end flow works

## 3. Polish

- [ ] 3.1 [Short descriptive title]
  - Acceptance: [Specific, testable condition]
  - Verification: `npm test && npm run build`
  - Dependencies: 2.1
  - Files: `src/path/to/file.ts`
  - Scope: XS

### Checkpoint: Complete
- [ ] All acceptance criteria met
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

## No Placeholders

Every step must contain the actual content an engineer needs. These are **plan failures**:
- "TBD", "TODO", "implement later", "fill in details"
- "Add appropriate error handling" / "add validation" / "handle edge cases"
- "Write tests for the above" (without actual test code)
- "Similar to Task N" (repeat the code — the engineer may be reading tasks out of order)
- Steps that describe what to do without showing how (code blocks required for code steps)

## Self-Review

After writing both documents:

1. **Spec coverage:** Can you point to a task that implements each spec requirement? List any gaps.
2. **Placeholder scan:** Search for red flags from the "No Placeholders" section. Fix them.
3. **Type consistency:** Do types, method signatures, and property names match across tasks?
4. **Alignment check:** Does every task in tasks.md reference files listed in plan.md's file map? Does every parallelization note in plan.md match the dependency declarations in tasks.md?

Fix any issues inline. If you find a spec requirement with no task, add the task.

## Execution Handoff

After saving both files, offer execution choice:

**"Plan complete and saved to `docs/rein/plans/YYYY-MM-DD-<feature-name>.md` and `docs/rein/tasks/YYYY-MM-DD-<feature-name>-tasks.md`. Two execution options:**

**1. Subagent-Driven (recommended)** — Fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** — Execute tasks in this session using incremental, batch execution with checkpoints

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
- plan.md contains task lists (should be in tasks.md)
- tasks.md contains architecture decisions (should be in plan.md)

## Verification

Before starting implementation, confirm:

- [ ] plan.md contains only decision-layer content (architecture, dependencies, risks)
- [ ] tasks.md contains only execution-layer content (ordered tasks with checkboxes)
- [ ] Every task has acceptance criteria
- [ ] Every task has a verification step
- [ ] Task dependencies are identified and ordered correctly
- [ ] No task touches more than ~5 files
- [ ] Checkpoints exist between major phases
- [ ] The human has reviewed and approved both documents
