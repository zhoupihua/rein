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
- [ ] 2.1 为X添加API端点
  - [ ] RED: 编写端点返回404的失败测试
  - [ ] GREEN: 实现返回正确响应的处理函数
  - [ ] REFACTOR: 提取验证逻辑为共享中间件
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

> **给执行代理的说明：** 这是实施过程中的主要参考文档。
> 仅在需要了解任务完成状态时读取 tasks.md。

**Goal:** [一句话描述本计划要构建什么]

**Architecture:** [2-3句话描述技术方案]

**Tech Stack:** [关键技术/库]

---

## Architecture Overview

[2-3段描述高层架构：主要组件、交互方式、数据流、关键设计模式。这是执行代理首先阅读的内容——确保自洽完整。]

## Architecture Decisions
- [关键决策1及理由]
- [关键决策2及理由]

## Dependency Graph

[ASCII树形图展示组件依赖关系]

```
Database schema
    ├── API models/types
    │       ├── API endpoints
    │       └── Validation logic
    └── Seed data / migrations
```

## Vertical Slice Strategy

[如何将工作切分为独立、可测试的增量。说明哪些切片可独立交付，以及为何所选顺序能将集成风险降到最低。]

## File Map

| File | Purpose | New/Modified |
|------|---------|-------------|
| `src/path/to/file.ts` | [用途] | 新建/修改 |

## Task Details

### 1.1 [简短描述性标题]
- **Acceptance:** [具体、可验证的验收条件]
- **Verification:** `npm test -- --grep "feature-name"`
- **Dependencies:** 无
- **Files:** `src/path/to/file.ts`
- **Scope:** S
- **Approach:** [实现方式——关键算法、模式或代码结构]
- **Edge Cases:** [边界条件及处理方式]
- **Rollback:** [如何回滚而不影响其他任务]

### 1.2 [简短描述性标题]
- **Acceptance:** [具体、可验证的验收条件]
- **Verification:** `npm run build`
- **Dependencies:** 1.1
- **Files:** `src/path/to/file.ts`, `tests/path/to/test.ts`
- **Scope:** M
- **Approach:** [实现方式——关键算法、模式或代码结构]
- **Edge Cases:** [边界条件及处理方式]
- **Rollback:** [如何回滚而不影响其他任务]

## Parallelization Classification

| Category | Tasks | Strategy |
|----------|-------|----------|
| Safe to parallelize | [任务编号] | 并行执行 |
| Sequential | [任务编号] | 按序执行 |
| Needs coordination | [任务编号] | 先定义接口，再并行 |

## Risk/Mitigation Table

| Risk | Impact | Mitigation |
|------|--------|------------|
| [风险] | [高/中/低] | [应对策略] |

## Self-Audit Checklist

- [ ] 每个规格需求都映射到至少一个任务
- [ ] 没有任务依赖后续任务（无循环依赖）
- [ ] 每个任务都有可独立验证的验收条件
- [ ] 没有占位符（TBD、TODO、"后续实现"）
- [ ] 文件路径具体且准确
- [ ] 每个任务都有回滚策略

## Explicit Handoff Statement

**此计划已准备好由 [子代理 / 内联执行] 实施。执行代理应：**
1. 首先阅读 plan.md 了解上下文和方案
2. 随着工作进展更新 task.md 复选框
3. 完成每个任务后运行验证命令
4. 遇到阻塞立即上报，不要绕过

## Open Questions
- [需要用户输入的问题]
```

### tasks.md Template (Execution Layer — Status Tracking Only)

tasks.md is a **lightweight checkbox list** for tracking progress. Implementation details live in plan.md. During execution, agents read plan.md for HOW, update tasks.md for STATUS.

```markdown
# Tasks: [Feature Name]

## 1. Foundation
- [ ] 1.1 为X创建数据库迁移
- [ ] 1.2 实现X的数据仓库层

## 2. Core Feature
- [ ] 2.1 添加X的API端点
- [ ] 2.2 构建X的UI

## 3. Polish
- [ ] 3.1 添加错误处理
- [ ] 3.2 更新文档
```

**Rules:**
- One line per task, no nested metadata (Acceptance/Verification/Dependencies/Files/Scope go in plan.md)
- Group by phase with `##` headings
- Number tasks for easy reference (1.1, 1.2, 2.1, ...)
- Task titles should be specific enough to identify the work, but details go in plan.md
- Sub-tasks use indented checkboxes with TDD structure (RED/GREEN/REFACTOR):
  ```
  - [ ] 2.1 为X添加API端点
    - [ ] RED: 编写POST /x返回201的失败测试
    - [ ] GREEN: 实现路由处理和创建操作
    - [ ] REFACTOR: 提取输入验证为共享中间件
  ```

## Task Sizing Guidelines

| Size | Files | Scope | Example |
|------|-------|-------|---------|
| **XS** | 1 | 单个函数或配置变更 | 添加一条验证规则 |
| **S** | 1-2 | 一个组件或端点 | 添加新的API端点 |
| **M** | 3-5 | 一个功能切片 | 用户注册流程 |
| **L** | 5-8 | 多组件功能 | 带筛选和分页的搜索 |
| **XL** | 8+ | **过大——需要进一步拆分** | — |

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
