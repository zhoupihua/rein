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
| Skills | 25 | Unified skills organized by SDLC phase |
| Agents | 3 | Expert personas (code-reviewer, test-engineer, security-auditor) |
| Commands | 13 | Slash commands from L1 quick fixes to L3 full features |
| Hooks | 9 | session-start, guard, guard-bash, gate, format, checkbox-guard, task-progress, leak-guard, inject |
| References | 4 | Testing, security, performance, accessibility checklists |
| Templates | 4 | Proposal, spec, design, tasks artifact templates |

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
                    ┌─────────────┐
                    │   /triage   │
                    │  自动分级    │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         ┌────────┐  ┌────────┐  ┌────────────┐
         │  L1    │  │  L2    │  │    L3      │
         │ /quick │  │  /fix  │  │  /feature  │
         │ ≤5行   │  │ 1-3文件│  │ 多文件特性  │
         └────┬───┘  └───┬────┘  └─────┬──────┘
              │          │             │
              ▼          ▼             ▼
         直接编辑    DEFINE→BUILD   DEFINE→PLAN→BUILD→REVIEW→SHIP
         测试→提交   测试→提交      完整8步铁三角
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
     │                │                │               │
     ▼                ▼                ▼               ▼
  需求明确?        RED→GREEN?       测试全绿?       提交消息
  复现确认?        重构完成?         回归通过?       符合规范?
```

质量门：TDD (RED→GREEN→REFACTOR) + 测试全绿

### L3: Feature (`/feature`)

完整的8步铁三角工作流，每步之间都有质量门：

```
┌──────────────────────────────────────────────────────────────────────┐
│                        L3 Full Workflow                              │
│                                                                      │
│  Step 1        Step 2        Step 3        Step 4                    │
│ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐                │
│ │  DEFINE  │ │  DEFINE  │ │  DEFINE  │ │  ISOLATE │                │
│ │  refine  │ │spec-     │ │  /spec   │ │  git-    │                │
│ │          │ │driven    │ │ 生成工件  │ │worktrees │                │
│ │          │ │          │ │          │ │          │                │
│ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘                │
│      │            │            │            │                        │
│      ▼            ▼            ▼            ▼                        │
│  ┌─────────────────────────────────────────────────┐                │
│  │           Quality Gate: DEFINE                   │                │
│  │  ✓ 问题陈述明确  ✓ 假设已列出                    │                │
│  │  ✓ 需求可测试    ✓ MVP范围已界定                  │                │
│  │  ✓ 人工审阅通过  ✓ 工件已提交                     │                │
│  └─────────────────────┬───────────────────────────┘                │
│                        │                                             │
│  Step 5        Step 6  │                                             │
│ ┌──────────┐ ┌────────┴─┐                                           │
│ │   PLAN   │ │   BUILD  │                                           │
│ │ planning │ │incremen- │                                           │
│ │          │ │tal + tdd │                                           │
│ │          │ │          │                                           │
│ └────┬─────┘ └────┬─────┘                                           │
│      │            │                                                  │
│      ▼            ▼                                                  │
│  ┌─────────────────────────────────────────────────┐                │
│  │           Quality Gate: PLAN + BUILD             │                │
│  │  ✓ 每个任务有验收标准  ✓ 每个任务有验证步骤       │                │
│  │  ✓ 依赖顺序正确       ✓ 无任务超过5文件           │                │
│  │  ✓ 阶段间有检查点     ✓ TDD: RED→GREEN→REFACTOR  │                │
│  │  ✓ 每个增量可编译     ✓ 测试全绿                  │                │
│  └─────────────────────┬───────────────────────────┘                │
│                        │                                             │
│  Step 7        Step 8  │                                             │
│ ┌──────────┐ ┌────────┴─┐                                           │
│ │  REVIEW  │ │   SHIP   │                                           │
│ │code-     │ │ verify + │                                           │
│ │review    │ │git-      │                                           │
│ │          │ │workflow  │                                           │
│ └────┬─────┘ └────┬─────┘                                           │
│      │            │                                                  │
│      ▼            ▼                                                  │
│  ┌─────────────────────────────────────────────────┐                │
│  │           Quality Gate: REVIEW + SHIP            │                │
│  │  ✓ 五轴审查通过 (正确性/可读性/架构/安全/性能)   │                │
│  │  ✓ 无安全漏洞      ✓ 无性能退化                   │                │
│  │  ✓ 用新证据验证    ✓ 测试全绿                     │                │
│  │  ✓ 提交消息规范    ✓ 工件归档                     │                │
│  └─────────────────────────────────────────────────┘                │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

### 质量门详细说明

| 质量门 | 位置 | 检查项 | 不通过则 |
|--------|------|--------|----------|
| **DEFINE** | Step 1-3 后 | 需求明确、假设列出、MVP界定、人工审阅 | 回到 refine 重新定义 |
| **PLAN + BUILD** | Step 5-6 后 | 任务有验收标准、TDD循环完成、每增量可编译、测试全绿 | 回到对应任务重做 |
| **REVIEW + SHIP** | Step 7-8 后 | 五轴审查通过、安全无漏洞、性能无退化、用新证据验证 | 修复问题后重新审查 |

### 恢复断点

任何步骤中断后，使用 `/continue` 自动识别当前阶段并继续：

| 检测条件 | 阶段 | 调用技能 |
|----------|------|----------|
| `docs/rein/specs/` 无文件 | DEFINE | refine |
| `docs/rein/plans/` 无文件 | PLAN | planning |
| `docs/rein/tasks/` 有未勾选项 | BUILD | incremental + TDD |
| 任务全勾选，无审查 | REVIEW | code-review |
| 审查完成，未提交 | SHIP | git-workflow |

## Commands

### Triage & Entry Points

| Command | Level | Use When |
|---------|-------|----------|
| `/triage` | — | Not sure which level? Auto-classify your change |
| `/quick` | L1 | ≤5 lines, no logic impact (typos, constants, configs) |
| `/fix` | L2 | 1-3 files, clear requirements (bug fix, small feature) |
| `/feature` | L3 | Multi-file, new feature, architecture change |

### Workflow Steps

| Command | Purpose | Output |
|---------|---------|--------|
| `/spec` | Generate design spec | `docs/rein/specs/` |
| `/plan` | Break spec into tasks with dependency graph | `docs/rein/plans/` + `docs/rein/tasks/` |
| `/do` | Execute tasks from tasks file | Code + commits |
| `/test` | TDD workflow | Tests passing |
| `/code-review` | 5-axis code review + security + performance | Review report |
| `/ship` | Parallel expert fan-out → GO/NO-GO | Merge / PR |
| `/simplify` | Code simplification | Simpler code |
| `/continue` | Resume from breakpoint | Continue workflow |
| `/status` | Task progress & drift detection | Fix stale checkboxes |
| `/archive` | Archive spec/plan/task to archive/ | Clean up completed artifacts |

## Hooks

rein 安装后自动配置 9 个钩子，无需手动干预：

| Hook | Event | Matcher | 作用 |
|------|-------|---------|------|
| session-start | SessionStart | * | 注入 using-rein 技能 + 活跃任务提醒 |
| guard | PreToolUse | Edit\|Write\|MultiEdit | 阻止修改 rein 管理的文件 |
| guard-bash | PreToolUse | Bash | 阻止破坏性命令操作 rein 文件 |
| gate | PreToolUse | Bash | deploy/push/publish 前自动跑测试 |
| format | PostToolUse | Write\|Edit\|MultiEdit | 自动 Prettier 格式化 |
| checkbox-guard | PostToolUse | Write\|Edit\|MultiEdit | 编辑 task 文件未更新 checkbox 时警告 |
| task-progress | PostToolUse | Write\|Edit\|MultiEdit | 代码编辑后注入任务进度，让 AI 可见 checkbox 状态 |
| leak-guard | PostToolUse | Read\|Bash | 拦截密钥泄露 (AKIA/sk-/ghp_) |
| inject | UserPromptExpansion | /code-review | 注入审查清单 |

## Skills by Phase

### Meta
- **using-rein** — Discovery and operating behaviors for all skills

### DEFINE
- **refine** — Structured divergent/convergent thinking
- **spec-driven** — Write PRD before code

### PLAN
- **planning** — Decompose specs into verifiable tasks
- **git-worktrees** — Isolated workspace on new branch

### BUILD
- **incremental** — Thin vertical slices, scope discipline
- **tdd** — RED-GREEN-REFACTOR iron law
- **subagent** — Dispatch parallel implementer agents
- **parallel-dispatch** — Independent task parallelization
- **context-engineering** — Right info at the right time
- **source-driven** — Ground decisions in official docs
- **frontend** — Production UI with accessibility
- **api-design** — Stable, hard-to-misuse interfaces

### VERIFY
- **debugging** — Systematic triage, no fix without root cause
- **browser-testing** — Live browser data via DevTools MCP
- **verify** — No claims without fresh evidence

### REVIEW
- **code-review** — 5-axis review, change size control
- **simplify** — Reduce complexity preserving behavior
- **security** — OWASP Top 10 prevention
- **performance** — Measure-first optimization

### SHIP
- **git-workflow** — Trunk-based, atomic commits
- **shipping** — Pre-launch checklist, staged rollout
- **cicd** — Quality gate pipeline
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
        ├── specs/                 # Design specs only (no tasks)
        │   └── YYYY-MM-DD-<name>-spec.md
        ├── plans/                 # Implementation plans (decision layer)
        │   └── YYYY-MM-DD-<name>-plan.md
        ├── tasks/                 # Task checklists (single source of truth)
        │   └── YYYY-MM-DD-<name>-task.md
        └── archive/               # Archived artifacts
            └── YYYY-MM-DD-<name>/
```

- **specs/** 只放设计文档（需求、决策、风险），不放任务清单
- **tasks/** 是任务的唯一来源，由 `/plan` 生成，由 `/do` 执行并勾选
- 工件按日期前缀 + 主题命名，同一特性的 spec/plan/tasks 共享相同前缀

### 生成流程

```
/spec  → specs/ (设计文档，不含任务)
/plan  → plans/ (决策层) + tasks/ (执行层，唯一任务源)
/do    → 读 tasks/ 逐项执行，勾选 [x]
```

## File Protection

rein 安装的文件受 `guard` 钩子保护，AI 无法修改或删除：

- 保护清单：`.claude/.rein-manifest`（安装时自动生成）
- 保护范围：hooks、commands、skills、agents、checklists
- 用户可新增自己的文件，不受保护
- 如需修改 rein 文件，从 `.rein-manifest` 中删除对应行即可

## Compared to Similar Projects

| Aspect | Source Projects | rein |
|--------|----------------|------|
| Install | 3 projects, npm + plugin + manual | 1 script, zero dependencies |
| Skills | 14 + 20 with overlaps | 25 redesigned, no duplicates |
| Commands | 3 deprecated + 7 separate | 12 unified |
| Spec management | Requires OpenSpec CLI | /spec for design, /plan for tasks, single task source |
| Templates | Generated by CLI | Static files, AI fills in |
| File protection | None | Auto-protect via hooks |
| Quality gates | Per-project, manual | Built-in at each phase boundary |
| Hooks | Session injection only | 9 hooks (session-start, guard, guard-bash, gate, format, checkbox-guard, task-progress, leak-guard, inject) |

## License

MIT
