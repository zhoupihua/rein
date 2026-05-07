---
name: planning
description: Use when you have a spec or requirements for a multi-step task, before touching code. Breaks work into ordered tasks with dependency graphs, vertical slicing, and explicit acceptance criteria.
---

# Planning and Task Breakdown

## Overview

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Decompose work into small, verifiable tasks with explicit acceptance criteria. Document everything: which files to touch, code, testing, docs. Bite-sized tasks. DRY. YAGNI. TDD. Frequent commits.

**Announce at start:** "I'm using the planning skill to create the implementation plan."

**Save output to:** `docs/rein/changes/<name>/plan.md` + `docs/rein/changes/<name>/task.md`
- `plan.md` — Architecture decisions, dependency graph, slicing strategy, risks (decision layer)
- `task.md` — Ordered task checklist with acceptance criteria (execution layer, the SINGLE source of truth for task tracking — MUST be generated, do not skip)
- If no `docs/rein/changes/<name>/` directory exists, create it

## Scope Check

If the spec covers multiple independent subsystems, suggest breaking into separate plans — one per subsystem. Each plan should produce working, testable software on its own.

## Planning Process

### Step 1: Enter Plan Mode

Before writing any code, operate in read-only mode:

- Read the spec (and proposal.md if it exists in the feature directory) and relevant codebase sections
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

**TDD-Structured Sub-Tasks** — Implementation tasks should include RED/GREEN/REFACTOR sub-checkboxes so each task is test-driven. Use this format in task.md:

```
- [ ] 2.1 Add API endpoint for X
  - [ ] RED: Write failing test for endpoint returning 404
  - [ ] GREEN: Implement handler returning correct response
  - [ ] REFACTOR: Extract validation into shared middleware
```

This ensures every implementation task follows the RED (write failing test) → GREEN (make it pass) → REFACTOR (clean up) cycle.

### Step 6: Order and Checkpoint

Arrange tasks so that:
1. Dependencies are satisfied (build foundation first)
2. Each task leaves the system in a working state
3. Verification checkpoints occur after every 2-3 tasks
4. High-risk tasks are early (fail fast)

## Output: Two Files

### plan.md Template (Decision + Implementation Layer)

```markdown
# [Feature Name] Plan

> **For agentic workers:** This is your primary reference during implementation.
> Read tasks.md only for status tracking (which tasks are done).

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

---

## Architecture Overview

[2-3 paragraphs describing the high-level architecture: major components, how they interact, data flow, and key design patterns. This is the first thing an implementing agent reads — make it self-contained.]

## Architecture Decisions
- [Key decision 1 and rationale]
- [Key decision 2 and rationale]

## Dependency Graph

[ASCII tree showing component dependencies]

```
Database schema
    ├── API models/types
    │       ├── API endpoints
    │       └── Validation logic
    └── Seed data / migrations
```

## Vertical Slice Strategy

[How work is sliced into independent, testable increments. Explain which slices can ship standalone and why the chosen order minimizes integration risk.]

## File Map

| File | Purpose | New/Modified |
|------|---------|-------------|
| `src/path/to/file.ts` | [Purpose] | New/Modified |

## Task Details

### 1.1 [Short descriptive title]
- **Acceptance:** [Specific, testable condition]
- **Verification:** `npm test -- --grep "feature-name"`
- **Dependencies:** None
- **Files:** `src/path/to/file.ts`
- **Scope:** S
- **Approach:** [How to implement — key algorithm, pattern, or code structure]
- **Edge Cases:** [Boundary conditions and how to handle them]
- **Rollback:** [How to revert if this task fails without breaking other tasks]

### 1.2 [Short descriptive title]
- **Acceptance:** [Specific, testable condition]
- **Verification:** `npm run build`
- **Dependencies:** 1.1
- **Files:** `src/path/to/file.ts`, `tests/path/to/test.ts`
- **Scope:** M
- **Approach:** [How to implement — key algorithm, pattern, or code structure]
- **Edge Cases:** [Boundary conditions and how to handle them]
- **Rollback:** [How to revert if this task fails without breaking other tasks]

## Parallelization Classification

| Category | Tasks | Strategy |
|----------|-------|----------|
| Safe to parallelize | [Task numbers] | Dispatch concurrently |
| Sequential | [Task numbers] | Execute in order |
| Needs coordination | [Task numbers] | Define contract first, then parallelize |

## Risk/Mitigation Table

| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk] | [High/Med/Low] | [Strategy] |

## Self-Audit Checklist

- [ ] Every spec requirement maps to at least one task
- [ ] No task depends on a later task (no circular dependencies)
- [ ] Every task has acceptance criteria that are independently verifiable
- [ ] No placeholders (TBD, TODO, "implement later") in any task detail
- [ ] File paths are specific and accurate for this codebase
- [ ] Rollback strategy exists for every task

## Explicit Handoff Statement

**This plan is ready for implementation by [subagent / inline execution]. The implementing agent should:**
1. Read this plan.md first for context and approach
2. Update task.md checkboxes as work progresses
3. Run verification commands listed in each task after completion
4. Flag any blocking issues immediately rather than working around them

## Open Questions
- [Question needing human input]
```

### tasks.md Template (Execution Layer — Status Tracking Only)

tasks.md is a **lightweight checkbox list** for tracking progress. Implementation details live in plan.md. During execution, agents read plan.md for HOW, update tasks.md for STATUS.

```markdown
# Tasks: [Feature Name]

## 1. Foundation
- [ ] 1.1 Create database migration for X
- [ ] 1.2 Implement repository layer for X

## 2. Core Feature
- [ ] 2.1 Add API endpoint for X
- [ ] 2.2 Build UI for X

## 3. Polish
- [ ] 3.1 Add error handling
- [ ] 3.2 Update documentation
```

**Rules:**
- One line per task, no nested metadata (Acceptance/Verification/Dependencies/Files/Scope go in plan.md)
- Group by phase with `##` headings
- Number tasks for easy reference (1.1, 1.2, 2.1, ...)
- Task titles should be specific enough to identify the work, but details go in plan.md
- Sub-tasks use indented checkboxes with TDD structure (RED/GREEN/REFACTOR):
  ```
  - [ ] 2.1 Add API endpoint for X
    - [ ] RED: Write failing test for POST /x returning 201
    - [ ] GREEN: Implement route handler and create action
    - [ ] REFACTOR: Extract input validation into shared middleware
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

1. **Spec coverage:** Can you point to a task in plan.md's Task Details that implements each spec requirement? List any gaps.
2. **Placeholder scan:** Search for red flags from the "No Placeholders" section. Fix them.
3. **Type consistency:** Do types, method signatures, and property names match across tasks in plan.md?
4. **Alignment check:** Extract task numbers from both files. Every checkbox in tasks.md (e.g., `- [ ] 1.1 ...`) must have a matching `### 1.1` section in plan.md, and vice versa. If numbers differ, fix before proceeding.

Fix any issues inline. If you find a spec requirement with no task, add the task.

## Execution Handoff

After saving both files, offer execution choice:

**"Plan complete and saved to `docs/rein/changes/<name>/plan.md` and `docs/rein/changes/<name>/task.md`. Two execution options:**

**1. Subagent-Driven (recommended)** — Fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** — Execute tasks in this session using executing-plans, batch execution with checkpoints

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
- plan.md Task Details say "implement the feature" without acceptance criteria
- No verification steps in plan.md Task Details
- All tasks are XL-sized
- Dependency order isn't considered
- Placeholders or vague steps in plan.md
- tasks.md contains implementation details (should be simple checkboxes only)
- tasks.md contains architecture decisions (should be in plan.md)
- plan.md and tasks.md task numbers don't match

## Verification

Before starting implementation, confirm:

- [ ] plan.md contains architecture decisions AND task details (acceptance, verification, files, dependencies)
- [ ] tasks.md is a simple checkbox list with no nested metadata
- [ ] Every task in tasks.md has a corresponding Task Details section in plan.md
- [ ] Task numbers match between plan.md and tasks.md
- [ ] No task touches more than ~5 files
- [ ] The human has reviewed and approved both documents

## Supporting Files

- **`plan-reviewer-prompt.md`** — Plan document reviewer subagent prompt template (verify completeness, spec alignment, task decomposition, buildability)
