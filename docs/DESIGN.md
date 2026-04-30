# 统一项目方案：rein

## Context

当前 AI 编程工作流需要配置 3 个独立项目（OpenSpec + Superpowers + Agent Skills），安装流程复杂（npm 包 + 插件市场 + 手动配置），且技能大量重叠。目标：将三者精华合并为 **一个零外部依赖的独立项目**，拉代码改一改，一条脚本完成安装。

---

## 项目结构

```
rein/
├── README.md
│
├── skills/                          # 22 个统一技能
│   │
│   │── # ===== 元技能 =====
│   ├── using-rein/              # 合并 using-superpowers + using-agent-skills
│   │   └── SKILL.md
│   │
│   │── # ===== DEFINE 阶段 =====
│   ├── refine/                 # 合并 SP:brainstorming + AS:refine
│   │   └── SKILL.md
│   ├── spec-driven/     # AS 独有
│   │   └── SKILL.md
│   │
│   │── # ===== PLAN 阶段 =====
│   ├── planning/ # 合并 SP:writing-plans + AS:planning
│   │   └── SKILL.md
│   ├── git-worktrees/         # SP 独有
│   │   └── SKILL.md
│   │
│   │── # ===== BUILD 阶段 =====
│   ├── incremental/  # 合并 SP:executing-plans + AS:incremental
│   │   └── SKILL.md
│   ├── tdd/     # 合并 SP:TDD + AS:TDD
│   │   └── SKILL.md
│   ├── subagent/ # SP 独有
│   │   └── SKILL.md
│   ├── parallel-dispatch/ # SP 独有
│   │   └── SKILL.md
│   ├── context-engineering/         # AS 独有
│   │   └── SKILL.md
│   ├── source-driven/   # AS 独有
│   │   └── SKILL.md
│   ├── frontend/     # AS 独有
│   │   └── SKILL.md
│   ├── api-design/    # AS 独有
│   │   └── SKILL.md
│   │
│   │── # ===== VERIFY 阶段 =====
│   ├── debugging/ # 合并 SP:systematic-debugging + AS:debugging
│   │   └── SKILL.md
│   ├── browser-testing/ # AS 独有
│   │   └── SKILL.md
│   ├── verify/ # SP 独有
│   │   └── SKILL.md
│   │
│   │── # ===== REVIEW 阶段 =====
│   ├── code-review/     # 合并 SP:requesting+receiving-code-review + AS:code-review
│   │   └── SKILL.md
│   ├── simplify/         # AS 独有
│   │   └── SKILL.md
│   ├── security/      # AS 独有
│   │   └── SKILL.md
│   ├── performance/    # AS 独有
│   │   └── SKILL.md
│   │
│   │── # ===== SHIP 阶段 =====
│   ├── git-workflow/ # AS 独有（融入 SP:finishing-a-development-branch）
│   │   └── SKILL.md
│   ├── shipping/         # AS 独有
│   │   └── SKILL.md
│   ├── cicd/        # AS 独有
│   │   └── SKILL.md
│   ├── migration/   # AS 独有
│   │   └── SKILL.md
│   └── docs-and-adrs/      # AS 独有
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
│   ├── feature.md                   # L3 完整变更（8 步铁三角流程）
│   ├── triage.md                    # 自动分级判定
│   ├── resume.md                    # 断点恢复
│   ├── spec.md                      # 定义规格（替代 /opsx:propose + /opsx:explore）
│   ├── plan.md                      # 拆解任务（替代 /opsx:continue + /opsx:ff）
│   ├── build.md                     # 增量构建
│   ├── test.md                      # TDD
│   ├── review.md                    # 5 轴审查
│   ├── ship.md                      # 并行专家审查 + GO/NO-GO
│   └── simplify.md                  # 代码简化
│
├── hooks/                           # 会话钩子
│   ├── session-start.sh             # 注入 using-rein 元技能
│   └── session-start.ps1            # Windows 版
│
├── references/                      # 参考清单（来自 AS）
│   ├── testing-patterns.md
│   ├── security-checklist.md
│   ├── performance-checklist.md
│   └── accessibility-checklist.md
│
├── install/                         # 安装脚本
│   ├── install.sh                   # Linux/Mac
│   └── install.ps1                  # Windows
│
└── templates/                       # 制品模板（替代 OpenSpec CLI 的模板生成）
    ├── proposal.md
    ├── spec.md
    ├── design.md
    └── tasks.md
```

---

## OpenSpec CLI 去除方案

OpenSpec CLI 的每个功能都有替代实现：

| OpenSpec CLI | 原功能 | 替代方案 |
|-------------|--------|---------|
| `openspec init` | 创建 `openspec/changes/`、`openspec/specs/` 目录 | install 脚本创建目录 |
| `openspec update` | 生成 `.claude/skills/` 等指令文件 | 不需要——我们自己写 skills，不存在自动生成 |
| `openspec validate` | 验证制品完整性 | `spec` 命令内置验证步骤（AI 检查制品是否齐全）|
| `openspec list` | 列出变更 | `resume` 命令扫描 `changes/` 目录 |
| `openspec status` | 查看制品进度 | `resume` 命令读取 `tasks.md` checkbox 状态 |
| `openspec archive` | 归档到 archive/ | `ship` 命令结尾执行归档（mv 到 archive/） |
| `/opsx:propose` | 提出变更 + 生成制品 | `spec` 命令（读取 templates/ 模板生成制品）|
| `/opsx:explore` | 探索性对话 | `spec` 命令（内置 explore 模式）|
| `/opsx:apply` | 按 tasks.md 实现 | `build` 命令 |
| `/opsx:verify` | 验证实现 | `review` 命令内置验证 |
| `/opsx:archive` | 归档变更 | `ship` 命令结尾 |
| `/opsx:continue` | 创建下一个制品 | `spec` 命令（可逐步或一次性生成）|
| `/opsx:ff` | 快进生成所有制品 | `spec` 命令的 `--full` 模式 |

### 制品目录结构（保留 OpenSpec 的核心设计）

```
<project-root>/
├── specs/                    # 已发布的规格（长期存在）
│   └── <feature>/
│       └── spec.md
├── changes/                  # 活跃变更（短期存在）
│   └── <change-name>/
│       ├── .change.yaml      # 变更元数据
│       ├── proposal.md       # 为什么做、要改什么
│       ├── specs/
│       │   └── <feature>/spec.md  # delta specs（ADDED/MODIFIED/REMOVED）
│       ├── design.md         # 工程决策 + Open Questions
│       └── tasks.md          # 任务清单（checkbox 格式）
└── archive/                  # 已归档变更
    └── YYYY-MM-DD-<name>/
        └── (完整制品副本)
```

> 相比 OpenSpec 的 `openspec/changes/`，简化为 `changes/`（少一层嵌套）。

### 制品模板（templates/）

`spec` 命令不再依赖 OpenSpec 的 TypeScript 模板引擎，而是直接在命令中内嵌模板内容。AI 读取模板后按项目上下文填充。

---

## 技能合并策略（6 对合并）

### 1. refine ← SP:brainstorming + AS:refine

**骨架**：SP brainstorming（更强的流程关卡、自审循环、到 planning 的显式交接）
**注入 AS 元素**：
- 7 个发散视角（逆向、约束移除、受众切换、组合、简化、10x、专家视角）
- "How Might We" 问题重构
- 假设显式化 + 验证策略
- "不做清单"作为必填输出
- 收敛时的三维度压力测试（用户价值/可行性/差异化）

### 2. planning ← SP:writing-plans + AS:planning

**骨架**：SP writing-plans（执行就绪的计划格式、无占位符规则、自审、显式交接）
**注入 AS 元素**：
- 依赖图可视化（ASCII 树）
- 垂直切片 vs 水平切片原则
- 任务大小表（XS/S/M/L/XL）+ 何时进一步拆分
- 风险/缓解表
- 并行化分类（safe/sequential/needs-coordination）
- 反合理化表

### 3. incremental ← SP:executing-plans + AS:incremental

**骨架**：AS incremental（实现规则、范围纪律、feature flags、rollback-friendly）
**注入 SP 元素**：
- 执行前审阅计划步骤
- 遇阻时的 stop-and-ask 协议
- 到 git-workflow 的显式交接

### 4. tdd ← SP:TDD + AS:TDD

**骨架**：SP TDD（铁律、删除先写代码的规则、12 条反合理化）
**注入 AS 元素**：
- 测试金字塔（80/15/5）+ 测试大小（S/M/L）
- DAMP over DRY 原则
- 状态测试 > 交互测试
- mock 偏好序：real > fake > stub > mock
- Prove-It 模式（bug 修复专用）
- 浏览器测试集成（DevTools MCP）

### 5. debugging ← SP:systematic-debugging + AS:debugging

**骨架**：AS debugging（6 步分诊、错误分类决策树、不可重现 bug 处理、安全回退）
**注入 SP 元素**：
- 铁律（无根因分析不修复）
- 3 次修复失败 → 质疑架构
- 反合理化表（8 条）
- 多组件诊断仪器

### 6. code-review ← SP:requesting+receiving-code-review + AS:code-review

**骨架**：AS code-review（5 轴审查、变更大小控制、严重性分类、多模型审查模式）
**注入 SP receiving-code-review 元素**：
- 禁止回复清单（不得说"好观点！"、"你说得对！"）
- YAGNI 检查（审查建议的"专业"特性）
- 不同来源的处理方式（人 vs 外部审查者）
- 多项反馈的实现顺序
- 反对时的优雅纠正

---

## 命令设计

### L1/L2/L3 分级命令（我们的定制）

**`/quick`** — L1 轻量变更（≤5 行，无逻辑影响）
```
直接修改 → 确认测试通过 → 提交
```

**`/fix`** — L2 标准变更（单文件/2-3 文件，需求明确）
```
Bug: debugging → tdd → verify → 提交
功能: tdd → verify → 提交
前端 Bug: + browser-testing
```

**`/feature`** — L3 完整变更（8 步铁三角流程）
```
1. refine → 发散/收敛，输出 one-pager
2. spec-driven → 生成 PRD
3. spec → 生成 changes/<name>/ 全套制品（proposal/specs/design/tasks）
4. git-worktrees → 分支隔离 + baseline
5. planning → 细化任务
6. incremental + tdd → 实现
   - 前端 → frontend
   - API → api-design
   - 并行 → subagent
   - 遇 bug → debugging
7. code-review → 5 轴审查
   - 安全 → security
   - 性能 → performance
8. verify → 验证
   → git-workflow → 提交
   → shipping → 发布检查
   → docs-and-adrs → 文档
   → 归档 changes/ 到 archive/
```

**`/triage`** — 自动分级判定

**`/continue`** — 断点恢复（读取 tasks.md checkbox 状态）

### 工作流命令（替代 OpenSpec 的 /opsx:*）

**`/spec`** — 替代 `/opsx:propose` + `/opsx:explore` + `/opsx:continue` + `/opsx:ff`
- 无参数：交互式选择模式（explore → propose → 逐步生成）
- `/spec <name>`：直接生成全套制品
- `/spec --step`：逐个生成制品（替代 /opsx:continue）
- `/spec --validate`：验证当前变更的制品完整性（替代 `openspec validate`）

**`/do`** — 替代 `/opsx:apply`
- 读取 tasks.md，逐项执行，勾选 checkbox

**`/plan`** — 独立调用 planning

**`/test`** — 独立调用 tdd + browser-testing

**`/review`** — 独立调用 code-review + security + performance

**`/ship`** — 并行 fan-out（3 专家代理）→ GO/NO-GO → 归档

**`/simplify`** — 独立调用 simplify

---

## 安装流程

### 一条命令安装（零 npm 依赖）

```bash
# Linux/Mac
git clone https://github.com/<org>/rein.git
cd your-project
bash /path/to/rein/install/install.sh

# Windows
git clone https://github.com/<org>/rein.git
cd your-project
powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1
```

### install 脚本做的事

1. 创建制品目录：`specs/`、`changes/`、`archive/`
2. 创建 `.claude/commands/` 目录，将 `commands/*.md` 复制进去
3. 创建 `.claude/skills/` 目录，将 `skills/` 复制进去
4. 创建 `.claude/agents/` 目录，将 `agents/*.md` 复制进去
5. 创建 `.claude/hooks/` 目录，将 `hooks/` 复制进去
6. 配置 `.claude/settings.json` 中的 hooks（session-start + 自定义守卫）
7. 在 `CLAUDE.md` 末尾追加工作流说明
8. 如有 `AGENTS.md`（Codex CLI），追加命令定义

> 不需要 `npm install -g`，不需要 `/plugin install`，不需要任何外部依赖。

---

## 与 Codex CLI 兼容

install 脚本检测平台，如果是 Codex CLI：
- 将命令内容写入 `AGENTS.md` 的 `## /command` 段落
- 将技能关键内容内联到 `AGENTS.md`（因为 Codex 无插件机制）
- 将强制规则写入 `AGENTS.md`（替代 hooks）

---

## 验证方式

1. 运行 install 脚本，确认目录和文件全部创建
2. 启动 Claude Code 新会话，确认 session-start hook 注入了 using-rein 元技能
3. 测试 `/triage`：输入一个变更描述，确认正确分级
4. 测试 `/spec test-feature`：确认生成 `changes/test-feature/` 全套制品
5. 测试 `/do`：确认读取 tasks.md 并执行
6. 测试 `/ship`：确认 fan-out 3 专家代理
7. 测试 `/continue`：中断后确认能恢复
8. 测试 hooks：编辑 sqlmigration 文件，确认被阻断
