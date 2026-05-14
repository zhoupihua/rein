# 统一项目方案：rein

## Context

当前 AI 编程工作流需要配置 3 个独立项目（OpenSpec + Superpowers + Agent Skills），安装流程复杂（npm 包 + 插件市场 + 手动配置），且技能大量重叠。目标：参考三者设计，独立实现 **一个零外部依赖的项目**，拉代码改一改，一条脚本完成安装。

---

## 项目结构

```
rein/
├── README.md
│
├── skills/                          # 21 个统一技能
│   │
│   │── # ===== 元技能 =====
│   ├── using-rein/              # 技能发现与操作行为
│   │   ├── SKILL.md
│   │   └── references/         # 平台适配文档
│   ├── writing-skills/          # 技能创建方法
│   │   ├── SKILL.md
│   │   ├── anthropic-best-practices.md
│   │   ├── persuasion-principles.md
│   │   ├── testing-skills-with-subagents.md
│   │   └── graphviz-conventions.dot
│   │
│   │── # ===== DEFINE 阶段 =====
│   ├── define/                  # 合并自 refine + spec-driven
│   │   ├── SKILL.md
│   │   ├── examples.md
│   │   ├── frameworks.md
│   │   ├── refinement-criteria.md
│   │   ├── visual-thinking.md
│   │   ├── spec-reviewer-prompt.md
│   │   ├── frame-template.html
│   │   └── helper.js
│   │
│   │── # ===== PLAN 阶段 =====
│   ├── planning/
│   │   ├── SKILL.md
│   │   └── plan-reviewer-prompt.md
│   │
│   │── # ===== BUILD 阶段 =====
│   ├── executing-plans/         # 合并自 executing-plans + incremental
│   │   └── SKILL.md
│   ├── subagent/                # 合并自 subagent + parallel-dispatch
│   │   ├── SKILL.md
│   │   ├── implementer-prompt.md
│   │   ├── spec-reviewer-prompt.md
│   │   └── code-quality-reviewer-prompt.md
│   ├── tdd/
│   │   └── SKILL.md
│   ├── refactor/                # Ralph loop 结构化重构
│   │   └── SKILL.md
│   ├── context-engineering/     # 合并自 context-engineering + source-driven
│   │   └── SKILL.md
│   ├── frontend/
│   │   └── SKILL.md
│   │
│   │── # ===== VERIFY 阶段 =====
│   ├── debugging/
│   │   ├── SKILL.md
│   │   ├── root-cause-tracing.md
│   │   ├── defense-in-depth.md
│   │   ├── condition-based-waiting.md
│   │   └── testing-anti-patterns.md
│   ├── browser-testing/
│   │   └── SKILL.md
│   ├── integration-testing/
│   │   └── SKILL.md
│   ├── verify/
│   │   └── SKILL.md
│   │
│   │── # ===== REVIEW 阶段 =====
│   ├── code-review/             # 合并自 code-review + simplify
│   │   └── SKILL.md
│   ├── security/
│   │   └── SKILL.md
│   ├── performance/
│   │   └── SKILL.md
│   │
│   │── # ===== SHIP 阶段 =====
│   ├── git-workflow/            # 合并自 git-workflow + git-worktrees
│   │   └── SKILL.md
│   ├── shipping/                # 合并自 shipping + cicd
│   │   └── SKILL.md
│   ├── migration/
│   │   └── SKILL.md
│   └── docs-and-adrs/
│       └── SKILL.md
│
├── agents/                          # 3 个专家代理
│   ├── code-reviewer.md             # Senior Staff Engineer
│   ├── test-engineer.md             # QA Specialist
│   └── security-auditor.md         # Security Engineer
│
├── commands/                        # 12 个统一斜杠命令
│   ├── quick.md                     # L1 轻量变更
│   ├── fix.md                       # L2 标准变更
│   ├── feature.md                   # L3 完整变更（6 步流程）
│   ├── continue.md                  # 断点恢复
│   ├── spec.md                      # 定义规格
│   ├── plan.md                      # 拆解任务
│   ├── do.md                        # 增量构建
│   ├── code-review.md               # 5 轴审查
│   ├── ship.md                      # 并行专家审查 + GO/NO-GO
│   ├── status.md                    # 任务进度 & 漂移检测
│   └── archive.md                   # 归档已完成工件
│   └── refactor.md                  # Ralph loop 结构化重构
│
├── hooks/                           # 会话钩子
│   ├── session-start.sh/ps1         # 注入 using-rein 元技能
│   ├── guard.sh/ps1                 # 文件保护
│   ├── guard-bash.sh/ps1            # Bash 保护
│   ├── gate.sh/ps1                  # 测试网关
│   ├── format.sh/ps1                # 自动 Prettier 格式化
│   ├── checkbox-guard.sh/ps1        # checkbox 警告
│   ├── leak-guard.sh/ps1            # 密钥泄露拦截
│   ├── inject.sh/ps1                # 注入审查清单
│   └── artifact-validate            # Go-only，验证制品阶段完整性
│
├── references/                      # 参考清单
│   ├── testing-patterns.md
│   ├── security-checklist.md
│   ├── performance-checklist.md
│   ├── accessibility-checklist.md
│   ├── orchestration-patterns.md
│   └── api-design.md              # 降级自 api-design 技能
│
├── install/                         # 安装脚本
│   ├── install.sh
│   └── install.ps1
│
└── templates/                       # 制品模板
    ├── proposal.md
    ├── spec.md
    ├── tasks.md
    └── checklists/review.md
```

---

## 制品目录结构

```
<project-root>/
└── docs/
    └── rein/
        ├── changes/               # 活跃的功能变更
        │   └── <name>/            # 每个功能一个目录
        │       ├── proposal.md    # DEFINE 阶段（Why, What Changes, Goals, Non-Goals, Assumptions, Open Questions）— L2 可选
        │       ├── spec.md        # DEFINE 阶段（Requirements, Decisions, Risks）
        │       ├── plan.md        # PLAN 阶段（Architecture, Dependency Graph, 实现计划）
        │       ├── task.md        # PLAN 阶段（任务清单，执行层，唯一任务源）
        │       └── review.md      # REVIEW 阶段（代码审查报告）
        ├── specs/                 # 主规格文件（可选，delta 合并目标）
        │   └── <domain>/spec.md
        ├── schema.json            # Artifact Graph 配置（可选）
        └── archive/               # 已归档变更
            └── <name>/
```

> 每个功能一个目录，所有工件集中管理。proposal.md 是 refine 阶段的产出（Why、Goals、NonGoals、Assumptions），L3 必需，L2 可选。spec.md 是 PRD 工件（需求、决策、风险）。plan.md 遵循 Superpowers 规范（架构、依赖图、切片策略、风险缓解、并行化、自审、交接）。task.md 是执行层唯一任务源。

### 生成流程

```
/spec   → changes/<name>/proposal.md + spec.md  (发散/收敛思考 → 动机/范围 → 需求/决策/风险)
/plan   → changes/<name>/plan.md + task.md
/do     → 读 task.md 逐项执行，勾选 [x]
```

---

## 技能设计策略

### 1. define ← refine + spec-driven

**合并逻辑：** refine 的发散/收敛思维和 spec-driven 的规格文档输出是同一个流程的两个阶段——先发散理解需求，再收敛写入规格。

**define/SKILL.md 结构：**
- Phase 1: Explore & Expand（原 refine Phase 1-2 的发散/收敛思维）
- Phase 2: Specify（原 spec-driven 的规格文档输出）
  - Step 2a: 写 proposal.md（Why/Goals/NonGoals/Assumptions/OpenQuestions）— L3 必需，L2 可选
  - Step 2b: 写 spec.md（Requirements/Decisions/Risks）
- Hard gate: 两个阶段都需人类审批后才继续

### 2. planning ← SP:writing-plans + AS:planning

**骨架**：SP writing-plans（执行就绪的计划格式、无占位符规则、自审、显式交接）
**注入 AS 元素**：依赖图可视化、垂直切片原则、任务大小表、风险/缓解表、并行化分类

### 3. executing-plans ← executing-plans + incremental

**合并逻辑：** 两者都是"执行计划"——executing-plans 是轻量内联执行，incremental 是切片策略和实现规则。合并为一个完整的执行技能。

**合并后结构：**
- Checkbox 循环（原 executing-plans 核心）
- 切片策略：垂直切片、Contract-First、Risk-First（原 incremental）
- 实现规则：Simplicity First、Scope Discipline 等（原 incremental）

### 4. tdd ← SP:TDD + AS:TDD

**骨架**：SP TDD（铁律）+ AS 元素（测试金字塔、DAMP over DRY、Prove-It 模式）

### 5. debugging ← SP:systematic-debugging + AS:debugging

**骨架**：AS debugging + SP 元素（铁律、3 次修复失败 → 质疑架构、反合理化表）

### 6. code-review ← code-review + simplify

**合并逻辑：** 简化是代码审查的一个输出——当审查发现复杂度问题时，简化是修复方式。

---

## 命令设计

### L1/L2/L3 分级命令

**`/quick`** — L1 轻量变更（≤5 行，无逻辑影响）
**`/fix`** — L2 标准变更（1-3 文件，需求明确）
**`/feature`** — L3 完整变更（6 步流程）

```
1. define → define（发散/收敛 → 生成 proposal.md + spec.md）
2. branch → git-workflow 分支隔离（worktree 或直接分支）
3. plan → planning 细化任务（plan.md + task.md）
4. implement → executing-plans + tdd 逐任务实现
5. review → code-review 5 轴审查（含简化）
6. ship → verify 验证 + git-workflow 提交 + shipping 发布检查
```

### 工作流命令

| 命令 | 用途 |
|------|------|
| `/spec` | 生成 proposal.md + spec.md |
| `/plan` | 拆解 spec 为任务 |
| `/do` | 逐项执行 task.md |
| `/code-review` | 5 轴审查 |
| `/ship` | 并行 fan-out + GO/NO-GO |
| `/continue` | 断点恢复 |
| `/status` | 任务进度 |
| `/archive` | 归档 |

---

## 安装流程

```bash
# Linux/Mac
cd your-project
bash /path/to/rein/install/install.sh

# Windows
cd your-project
powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1
```

### install 脚本做的事

1. 创建制品目录：`docs/rein/changes/`、`docs/rein/archive/`
2. 复制 `commands/*.md` 到 `.rein/commands/`
3. 复制 `skills/` 到 `.rein/skills/`
4. 复制 `agents/*.md` 到 `.rein/agents/`
5. 配置 `.claude/settings.json` 中的 hooks
6. 生成 `.rein/.rein-manifest` 保护清单

---

## 验证方式

1. 运行 install 脚本，确认目录和文件全部创建
2. 启动 Claude Code 新会话，确认 session-start hook 注入了 using-rein 元技能
3. 测试 `/spec test-feature`：确认生成 `docs/rein/changes/test-feature/proposal.md` + `spec.md`
4. 测试 `/do`：确认读取 tasks.md 并执行
5. 测试 `/ship`：确认 fan-out 3 专家代理
6. 测试 `/continue`：中断后确认能恢复
7. 测试 hooks：编辑 rein 管理文件，确认被阻断
