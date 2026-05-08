# rein

A unified, zero-dependency AI coding workflow inspired by [Superpowers](https://github.com/obra/superpowers), [Agent Skills](https://github.com/addyosmani/agent-skills), and [OpenSpec](https://github.com/Fission-AI/OpenSpec).

## Why

Projects like Superpowers, Agent Skills, and OpenSpec each explored AI coding workflows independently, but using them together meant:
- 3 separate projects to install and configure
- npm package (OpenSpec CLI) + plugin marketplace (Superpowers) + manual setup (Agent Skills)
- Overlapping concepts with different names and slight variations

**rein** builds on their ideas as a single zero-dependency project. Clone, run install script, done.

## What's Inside

| Component | Count | Description |
|-----------|-------|-------------|
| Skills | 20 | Unified skills organized by SDLC phase |
| Agents | 3 | Expert personas (code-reviewer, test-engineer, security-auditor) |
| Commands | 11 | Slash commands from L1 quick fixes to L3 full features |
| Hooks | 10 | session-start, guard, guard-bash, gate, format, checkbox-guard, task-progress, leak-guard, inject, artifact-validate |
| References | 6 | Testing, security, performance, accessibility, orchestration, api-design checklists |
| Templates | 4 | Proposal, spec, tasks, review-checklist artifact templates |

## Quick Start

```bash
# Clone this repo
git clone https://github.com/zhoupihua/rein.git

# Install into your project (Linux/Mac)
cd your-project
bash /path/to/rein/install/install.sh

# Install into your project (Windows)
cd your-project
powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1
```

No `npm install -g`, no plugin marketplace, no manual configuration.

## Workflow Overview

rein 的工作流按变更复杂度分为三级，每级有对应的流程和质量门：

```
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         ┌────────┐  ┌────────┐  ┌────────────┐
         │  L1    │  │  L2    │  │    L3      │
         │ /quick │  │  /fix  │  │  /feature  │
         │ ≤5行   │  │ 1-3文件│  │ 多文件特性  │
         └────┬───┘  └───┬────┘  └─────┬──────┘
              │          │             │
              ▼          ▼             ▼
         直接编辑    DEFINE→BUILD   DEFINE→BRANCH→PLAN→BUILD→REVIEW→SHIP
         测试→提交   测试→提交      完整6步工作流
```

### L1: Quick (`/quick`)

```
编辑 → 测试通过 → 提交
```

质量门：测试通过

### L2: Fix (`/fix`)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────┐
│  DEFINE  │────▶│  BUILD   │────▶│  VERIFY  │────▶│ SHIP │
│ 理解需求  │     │ TDD实现   │     │ 测试通过  │     │ 提交  │
│ 定位问题  │     │ 修复/新增  │     │ 验证修复  │     │      │
└──────────┘     └──────────┘     └──────────┘     └──────┘
```

质量门：TDD (RED→GREEN→REFACTOR) + 测试全绿

### L3: Feature (`/feature`)

完整的6步工作流，每步之间都有质量门：

```
┌──────────────────────────────────────────────────────────────────────┐
│                        L3 Full Workflow                              │
│                                                                      │
│  Step 1              Step 2              Step 3                      │
│ ┌────────────────┐ ┌────────────────┐ ┌────────────────┐            │
│ │    DEFINE      │ │    BRANCH      │ │     PLAN       │            │
│ │    define      │ │  git-workflow  │ │   planning     │            │
│ │ proposal+spec  │ │  分支隔离       │ │  拆解任务       │            │
│ └───────┬────────┘ └───────┬────────┘ └───────┬────────┘            │
│         │                  │                   │                      │
│         ▼                  ▼                   ▼                      │
│  ┌─────────────────────────────────────────────────┐                │
│  │        Quality Gate: DEFINE → PLAN               │                │
│  └─────────────────────┬───────────────────────────┘                │
│                        │                                             │
│  Step 4              Step 5           Step 6                         │
│ ┌────────────────┐ ┌───────────────┐ ┌───────────────┐              │
│ │    BUILD       │ │    REVIEW     │ │     SHIP      │              │
│ │executing-plans │ │  code-review  │ │ verify +      │              │
│ │    + tdd       │ │  5轴审查       │ │ git-workflow  │              │
│ │ 逐任务实现      │ │               │ │ 归档           │              │
│ └───────┬────────┘ └───────┬───────┘ └───────┬───────┘              │
│         │                  │                  │                       │
│         ▼                  ▼                  ▼                       │
│  ┌─────────────────────────────────────────────────┐                │
│  │        Quality Gate: BUILD → SHIP                │                │
│  └─────────────────────────────────────────────────┘                │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

### 质量门详细说明

| 质量门 | 位置 | 检查项 | 不通过则 |
|--------|------|--------|----------|
| **DEFINE → PLAN** | Step 1-3 后 | 需求明确、MVP界定、任务有验收标准、依赖顺序正确 | 回到 define 重新定义 |
| **BUILD → SHIP** | Step 4-6 后 | TDD循环完成、每增量可编译、测试全绿、五轴审查通过 | 回到对应任务重做 |

### 恢复断点

任何步骤中断后，使用 `/continue` 自动识别当前阶段并继续：

| 检测条件 | 阶段 | 调用技能 |
|----------|------|----------|
| `docs/rein/changes/<name>/` 无 proposal.md 且无 spec.md | DEFINE | define |
| 有 proposal.md，无 spec.md | DEFINE | define (resume from Step 2b) |
| `docs/rein/changes/<name>/` 无 plan.md | PLAN | planning |
| `docs/rein/changes/<name>/task.md` 有未勾选项 | BUILD | executing-plans + TDD |
| 任务全勾选，无审查 | REVIEW | code-review |
| 审查完成，未提交 | SHIP | git-workflow |

## Commands

### Entry Points

| Command | Level | Use When |
|---------|-------|----------|
| `/quick` | L1 | ≤5 lines, no logic impact (typos, constants, configs) |
| `/fix` | L2 | 1-3 files, clear requirements (bug fix, small feature) |
| `/feature` | L3 | Multi-file, new feature, architecture change |

### Workflow Steps

| Command | Purpose | Output |
|---------|---------|--------|
| `/spec` | Generate proposal + spec | `docs/rein/changes/<name>/proposal.md` + `spec.md` |
| `/plan` | Break spec into tasks with dependency graph | `docs/rein/changes/<name>/plan.md` + `docs/rein/changes/<name>/task.md` |
| `/do` | Execute tasks from tasks file | Code + commits |
| `/code-review` | 5-axis code review + security + performance | Review report |
| `/ship` | Parallel expert fan-out → GO/NO-GO | Merge / PR |
| `/continue` | Resume from breakpoint | Continue workflow |
| `/status` | Task progress & drift detection | Fix stale checkboxes |
| `/archive` | Archive spec/plan/task to archive/ | Clean up completed artifacts |

## Hooks

rein 安装后自动配置 10 个钩子，无需手动干预：

| Hook | Event | Matcher | 作用 |
|------|-------|---------|------|
| session-start | SessionStart | * | 注入 using-rein 技能 + 活跃任务提醒 |
| guard | PreToolUse | Edit\|Write\|MultiEdit | 阻止修改 rein 管理的文件 |
| guard-bash | PreToolUse | Bash | 阻止破坏性命令操作 rein 文件 |
| gate | PreToolUse | Bash | deploy/push/publish 前自动跑测试 |
| format | PostToolUse | Write\|Edit\|MultiEdit | 自动 Prettier 格式化 |
| checkbox-guard | PostToolUse | Write\|Edit\|MultiEdit | 编辑 task 文件未更新 checkbox 时警告 |
| task-progress | PostToolUse | Write\|Edit\|MultiEdit | 编辑代码文件时自动勾选匹配的 task checkbox |
| leak-guard | PostToolUse | Read\|Bash | 拦截密钥泄露 (AKIA/sk-/ghp_) |
| inject | UserPromptExpansion | /code-review | 注入审查清单 |
| artifact-validate | PostToolUse | Write\|Edit\|MultiEdit | 编辑 changes/ 下工件时验证阶段完整性 |

## Skills by Phase

### Meta
- **using-rein** — Discovery and operating behaviors for all skills
- **writing-skills** — Skill creation as TDD, bulletproofing, persuasion

### DEFINE
- **define** — Refine ideas through divergent/convergent thinking, then write spec

### PLAN
- **planning** — Decompose specs into verifiable tasks

### BUILD
- **executing-plans** — Thin vertical slices, test each increment, commit frequently
- **subagent** — Dispatch implementer agents per task (sequential or parallel)
- **tdd** — RED-GREEN-REFACTOR iron law
- **context-engineering** — Right context at the right time, verify against official docs
- **frontend** — Production UI with accessibility

### VERIFY
- **debugging** — Systematic triage, no fix without root cause
- **browser-testing** — Live browser data via DevTools MCP
- **integration-testing** — Integration and contract testing patterns
- **verify** — No claims without fresh evidence

### REVIEW
- **code-review** — 5-axis review including simplification, change size control
- **security** — OWASP Top 10 prevention
- **performance** — Measure-first optimization

### SHIP
- **git-workflow** — Trunk-based, atomic commits, worktree isolation
- **shipping** — Pre-launch checklist, CI/CD, staged rollout
- **migration** — Strangler pattern, zombie code removal
- **docs-and-adrs** — Document decisions, not code

## Expert Agents

| Agent | Role | Use |
|-------|------|-----|
| code-reviewer | Senior Staff Engineer | Five-axis code review |
| test-engineer | QA Specialist | Coverage analysis, Prove-It pattern |
| security-auditor | Security Engineer | OWASP assessment, threat modeling |

## Artifact Structure

```
<project-root>/
└── docs/
    └── rein/
        ├── changes/               # Active feature work
        │   └── <name>/            # One directory per feature
        │       ├── proposal.md    # DEFINE phase (Why, What Changes, Goals, Non-Goals, Assumptions, Open Questions) — optional for L2
        │       ├── spec.md        # DEFINE phase (Requirements, Decisions, Risks)
        │       ├── plan.md        # PLAN phase (Architecture, Dependency Graph, Implementation plan)
        │       ├── task.md        # PLAN phase (Checkbox task list, single source of truth)
        │       └── review.md      # REVIEW phase (Code review report)
        ├── specs/                 # Master specification files (optional)
        │   └── <domain>/spec.md
        ├── schema.json            # Artifact graph configuration (optional)
        └── archive/               # Archived artifacts
            └── <name>/
```

- **changes/<name>/** 每个功能一个目录，所有工件集中管理
- **proposal.md** 是 refine 阶段的产出（Why/Goals/NonGoals/Assumptions），L3 必需，L2 可选
- **spec.md** 是 PRD 工件 — 需求（WHEN/THEN/TEST）、决策、风险
- **plan.md** 遵循 Superpowers 规范 — 架构概览、依赖图、切片策略、风险缓解、并行化分类、自审、交接
- **task.md** 是任务的唯一来源，由 `/plan` 生成，由 `/do` 执行并勾选

### 生成流程

```
/spec   → changes/<name>/proposal.md + spec.md  (发散/收敛思考 → 动机/范围 → 需求/决策/风险)
/plan   → changes/<name>/plan.md + task.md
/do     → 读 task.md 逐项执行，勾选 [x]
```

## File Protection

rein 安装的文件受 `guard` 钩子保护，AI 无法修改或删除：

- 保护清单：`.claude/.rein-manifest`（安装时自动生成）
- 保护范围：hooks、commands、skills、agents、references
- 用户可新增自己的文件，不受保护
- 如需修改 rein 文件，从 `.rein-manifest` 中删除对应行即可

## Compared to Similar Projects

| Aspect | Source Projects | rein |
|--------|----------------|------|
| Install | 3 projects, npm + plugin + manual | 1 script, zero dependencies |
| Skills | 14 + 20 with overlaps | 20 redesigned, no duplicates |
| Commands | 3 deprecated + 7 separate | 11 unified |
| Spec management | Requires OpenSpec CLI | /spec for proposal+PRD+decisions, /plan for tasks, single task source |
| Templates | Generated by CLI | Static files, AI fills in |
| File protection | None | Auto-protect via hooks |
| Quality gates | Per-project, manual | Built-in at each phase boundary |
| Hooks | Session injection only | 10 hooks (session-start, guard, guard-bash, gate, format, checkbox-guard, task-progress, leak-guard, inject, artifact-validate) |

## License

MIT
