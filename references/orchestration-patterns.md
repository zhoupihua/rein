# Orchestration Patterns

Reference catalog of agent orchestration patterns this repo endorses, plus anti-patterns to avoid. Read this before adding a new slash command that coordinates multiple personas, or before introducing a new persona that "wraps" existing ones.

The governing rule: **the user (or a slash command) is the orchestrator. Personas do not invoke other personas.** Skills are mandatory hops inside a persona's workflow.

---

## Endorsed Patterns

### 1. Direct Invocation (No Orchestration)

Single persona, single perspective, single artifact. The default and cheapest option.

```
user → code-reviewer → report → user
```

**Use when:** the work is one perspective on one artifact and you can describe it in one sentence.

**Examples:**
- "Review this PR" → `code-reviewer`
- "Find security issues in `auth.ts`" → `security-auditor`
- "What tests are missing for the checkout flow?" → `test-engineer`

**Cost:** one round trip.

---

### 2. Single-Persona Slash Command

A slash command that wraps one persona with the project's skills.

```
/review → code-reviewer (with code-review skill) → report
```

**Use when:** the same single-persona invocation happens repeatedly with the same setup.

**Anti-signal:** if the slash command's body is mostly "decide which persona to call," delete it and let the user call the persona directly.

---

### 3. Parallel Fan-Out with Merge

Multiple personas operate on the same input concurrently, each producing an independent report. A merge step synthesizes them into a single decision.

```
                    ┌─→ code-reviewer    ─┐
/ship → fan out  ───┼─→ security-auditor ─┤→ merge → go/no-go + rollback
                    └─→ test-engineer    ─┘
```

**Use when:**
- The sub-tasks are genuinely independent (no shared mutable state, no ordering dependency)
- Each sub-agent benefits from its own context window
- The merge step is small enough to stay in the main context
- Wall-clock latency matters

**Examples in rein:** `/ship`

**Cost:** N parallel sub-agent contexts + one merge turn.

**Validation checklist before adopting:**
- [ ] Can I run all sub-agents at the same time without ordering issues?
- [ ] Does each persona produce a different *kind* of finding?
- [ ] Will the merge step fit in the main agent's remaining context?
- [ ] Is the user's wait time long enough that parallelism is actually noticeable?

If any answer is "no," fall back to direct invocation or a single-persona command.

---

### 4. Sequential Pipeline as User-Driven Slash Commands

The user runs slash commands in a defined order, carrying context between them. The user IS the orchestrator.

```
user runs:  /quick  →  /fix  →  /feature  →  /do  →  /ship
```

**Use when:** the workflow has dependencies and human judgment between steps adds value.

**Examples in rein:** the entire L1/L2/L3 lifecycle.

**Cost:** one sub-agent context per step. No orchestration overhead.

---

### 5. Research Isolation (Context Preservation)

When a task requires reading large amounts of material, spawn a research sub-agent that returns only a digest.

```
main agent → research sub-agent (reads 50 files) → digest → main agent continues
```

**Use when:**
- The main session needs to stay focused on a downstream task
- The investigation result is much smaller than the input it consumes

**On Claude Code, use the built-in `Explore` subagent** rather than defining a custom research persona.

---

## Claude Code Compatibility

### Where Personas Live

Agent persona files go in `agents/` at the project root: `agents/code-reviewer.md`, `agents/security-auditor.md`, `agents/test-engineer.md`.

### Subagents vs. Agent Teams

| | Subagents | Agent Teams |
|--|-----------|-------------|
| Coordination | Main agent fans out, sub-agents only report back | Teammates message each other, share a task list |
| Context | Own context window per subagent | Own context window per teammate |
| When to use | Independent tasks producing reports | Collaborative work needing discussion |
| Status | Stable | Experimental |

### Platform-Enforced Rules

- **Subagents cannot spawn other subagents** — Anti-pattern B and D cannot exist by construction.
- **No nested teams** — teammates cannot spawn their own teams.

### Built-in Subagents

| Built-in | Purpose |
|----------|---------|
| `Explore` | Read-only codebase search and analysis. Use for Pattern 5. |
| `Plan` | Read-only research during plan mode. |
| `general-purpose` | Multi-step tasks needing both exploration and modification. |

---

## Anti-Patterns

### A. Router Persona ("Meta-Orchestrator")

A persona whose job is to decide which other persona to call.

**Why it fails:** Pure routing layer with no domain value. Adds two paraphrasing hops → information loss + 2x token cost.

**What to do instead:** add or refine slash commands. Document intent → command mapping in `AGENTS.md`.

---

### B. Persona That Calls Another Persona

A `code-reviewer` that internally invokes `security-auditor` when it sees auth code.

**Why it fails:** Personas produce a single perspective; chaining defeats that. The summary the calling persona passes loses context. Hides cost from the user.

**What to do instead:** have the calling persona *recommend* a follow-up audit in its report. The user or a slash command runs the second pass.

---

### C. Sequential Orchestrator That Paraphrases

An agent that calls `/spec`, then `/plan`, then `/build`, etc. on the user's behalf.

**Why it fails:** Loses human checkpoints that catch wrong-direction work. Each hand-off summarizes context — accumulated drift. Doubles token cost. Removes user agency.

**What to do instead:** keep the user as the orchestrator. Document the recommended sequence.

---

### D. Deep Persona Trees

`/ship` calls a `pre-ship-coordinator` that calls a `quality-coordinator` that calls `code-reviewer`.

**Why it fails:** Each layer adds latency and tokens with no decision value. Leaf personas lose context to multiple summarization steps.

**What to do instead:** keep orchestration depth at most 1 (slash command → personas). The merge happens in the main agent.

---

## Decision Flow

```
Is the work one perspective on one artifact?
├── Yes → Direct invocation. Stop.
└── No  → Will the same composition repeat?
         ├── No  → Direct invocation, ad hoc. Stop.
         └── Yes → Are sub-tasks independent?
                  ├── No  → Sequential slash commands run by user (Pattern 4).
                  └── Yes → Parallel fan-out with merge (Pattern 3).
                           Validate against the checklist above.
                           If any check fails → fall back to single-persona command (Pattern 2).
```
