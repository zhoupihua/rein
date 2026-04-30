# rein

A unified, zero-dependency AI coding workflow that consolidates the best of [Superpowers](https://github.com/obra/superpowers), [Agent Skills](https://github.com/addyosmani/agent-skills), and [OpenSpec](https://github.com/Fission-AI/OpenSpec) into one project.

## Why

Setting up the "Iron Triangle" (OpenSpec + Superpowers + Agent Skills) requires:
- 3 separate projects to install and configure
- npm package (OpenSpec CLI) + plugin marketplace (Superpowers) + manual setup (Agent Skills)
- Overlapping skills with different names and slight variations

**rein** merges them into a single project with zero external dependencies. Clone, run install script, done.

## What's Inside

| Component | Count | Description |
|-----------|-------|-------------|
| Skills | 25 | Unified skills organized by SDLC phase |
| Agents | 3 | Expert personas (code-reviewer, test-engineer, security-auditor) |
| Commands | 12 | Slash commands from L1 quick fixes to L3 full features |
| Hooks | 2 | Session start injection (bash + PowerShell) |
| References | 4 | Testing, security, performance, accessibility checklists |
| Templates | 4 | Proposal, spec, design, tasks artifact templates |

## Quick Start

```bash
# Clone this repo
git clone https://github.com/<org>/rein.git

# Install into your project (Linux/Mac)
cd your-project
bash /path/to/rein/install/install.sh

# Install into your project (Windows)
cd your-project
powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1
```

No `npm install -g`, no plugin marketplace, no manual configuration.

## Commands

### Triage & Entry Points

| Command | Level | Use When |
|---------|-------|----------|
| `/triage` | — | Not sure which level? Auto-classify your change |
| `/quick` | L1 | ≤5 lines, no logic impact (typos, constants, configs) |
| `/fix` | L2 | 1-3 files, clear requirements (bug fix, small feature) |
| `/feature` | L3 | Multi-file, new feature, architecture change |

### Workflow Steps

| Command | Purpose | Replaces |
|---------|---------|----------|
| `/spec` | Generate change artifacts (proposal/specs/design/tasks) | `/opsx:propose` + `/opsx:explore` + `/opsx:continue` + `/opsx:ff` |
| `/plan` | Task breakdown with dependency graph | — |
| `/build` | Execute tasks from tasks.md | `/opsx:apply` |
| `/test` | TDD workflow + browser testing | — |
| `/review` | 5-axis code review + security + performance | `/opsx:verify` |
| `/ship` | Parallel expert fan-out → GO/NO-GO | `/opsx:archive` |
| `/simplify` | Code simplification | — |
| `/resume` | Resume from breakpoint | `openspec status` |

## Skills by Phase

### Meta
- **using-workflow** — Discovery and operating behaviors for all skills

### DEFINE
- **idea-refine** — Structured divergent/convergent thinking
- **spec-driven-development** — Write PRD before code

### PLAN
- **planning-and-task-breakdown** — Decompose specs into verifiable tasks
- **using-git-worktrees** — Isolated workspace on new branch

### BUILD
- **incremental-implementation** — Thin vertical slices, scope discipline
- **test-driven-development** — RED-GREEN-REFACTOR iron law
- **subagent-driven-development** — Dispatch parallel implementer agents
- **dispatching-parallel-agents** — Independent task parallelization
- **context-engineering** — Right info at the right time
- **source-driven-development** — Ground decisions in official docs
- **frontend-ui-engineering** — Production UI with accessibility
- **api-and-interface-design** — Stable, hard-to-misuse interfaces

### VERIFY
- **debugging-and-error-recovery** — Systematic triage, no fix without root cause
- **browser-testing-with-devtools** — Live browser data via DevTools MCP
- **verification-before-completion** — No claims without fresh evidence

### REVIEW
- **code-review-and-quality** — 5-axis review, change size control
- **code-simplification** — Reduce complexity preserving behavior
- **security-and-hardening** — OWASP Top 10 prevention
- **performance-optimization** — Measure-first optimization

### SHIP
- **git-workflow-and-versioning** — Trunk-based, atomic commits
- **shipping-and-launch** — Pre-launch checklist, staged rollout
- **ci-cd-and-automation** — Quality gate pipeline
- **deprecation-and-migration** — Strangler pattern, zombie code removal
- **documentation-and-adrs** — Document decisions, not code

## Expert Agents

| Agent | Role | Use |
|-------|------|-----|
| code-reviewer | Senior Staff Engineer | Five-axis code review |
| test-engineer | QA Specialist | Coverage analysis, Prove-It pattern |
| security-auditor | Security Engineer | OWASP assessment, threat modeling |

## Artifact Structure

```
<project-root>/
├── specs/                    # Published specs (long-lived)
│   └── <feature>/spec.md
├── changes/                  # Active changes (short-lived)
│   └── <change-name>/
│       ├── .change.yaml      # Change metadata
│       ├── proposal.md       # Why and what
│       ├── specs/
│       │   └── <feature>/spec.md  # Delta specs
│       ├── design.md         # Engineering decisions
│       └── tasks.md          # Task checklist
└── archive/                  # Archived changes
    └── YYYY-MM-DD-<name>/
```

## Workflow

### L1: Quick (`/quick`)
```
Direct edit → confirm tests pass → commit
```

### L2: Fix (`/fix`)
```
Bug: debugging → TDD → verify → commit
Feature: TDD → verify → commit
Frontend: + browser testing
```

### L3: Feature (`/feature`)
```
1. idea-refine → one-pager
2. spec-driven-development → PRD
3. /spec → artifacts
4. using-git-worktrees → branch isolation
5. planning-and-task-breakdown → task list
6. incremental-implementation + TDD → code
7. code-review-and-quality → 5-axis review
8. verification-before-completion → verify
   → commit → release check → docs → archive
```

## What's Different from the Source Projects

| Aspect | Source Projects | rein |
|--------|----------------|-----------------|
| Install | 3 projects, npm + plugin + manual | 1 script, zero dependencies |
| Skills | 14 + 20 with overlaps | 25 merged, no duplicates |
| Commands | 3 deprecated + 7 separate | 12 unified |
| Spec management | Requires OpenSpec CLI | Built into /spec command |
| Templates | Generated by CLI | Static files, AI fills in |

## License

MIT
