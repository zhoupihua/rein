---
name: define
description: Use before any creative work - creating features, building components, adding functionality, or modifying behavior. Refines ideas through structured divergent and convergent thinking, then produces a spec that is the shared source of truth for implementation.
---

# Define: From Idea to Spec

Write a structured specification before writing any code. The spec is the shared source of truth between you and the human engineer — it defines what we're building and how we'll know it's done. Code without a spec is guessing.

<HARD-GATE>
Do NOT invoke any implementation skill, write any code, scaffold any project, or take any implementation action until you have presented a spec and the user has approved it. This applies to EVERY project regardless of perceived simplicity.
</HARD-GATE>

## Anti-Pattern: "This Is Too Simple To Need A Design"

Every project goes through this process. A todo list, a single-function utility, a config change — all of them. "Simple" projects are where unexamined assumptions cause the most wasted work. The design can be short (a few sentences for truly simple projects), but you MUST present it and get approval.

## When to Use

- Starting a new project or feature
- Requirements are ambiguous or incomplete
- The change touches multiple files or modules
- You're about to make an architectural decision
- The task would take more than 30 minutes to implement

**When NOT to use:** Single-line fixes, typo corrections, or changes where requirements are unambiguous and self-contained.

## The Gated Workflow

Define has three phases. Do not advance to the next phase until the current one is validated.

```
EXPLORE ──→ SPECIFY ──→ HAND OFF
   │           │           │
   ▼           ▼           ▼
 Human       Human       Human
 reviews     reviews     reviews
```

## Checklist

You MUST create a task for each of these items and complete them in order:

1. **Explore project context** — check files, docs, recent commits
2. **Offer visual companion** (if topic will involve visual questions) — this is its own message, not combined with a clarifying question
3. **Ask clarifying questions** — one at a time, understand purpose/constraints/success criteria
4. **Generate idea variations** — 5-8 variations using divergent thinking lenses
5. **Propose 2-3 approaches** — with trade-offs and your recommendation
6. **Present design** — in sections scaled to complexity, get user approval after each section
7. **Write proposal** — save to `docs/rein/changes/<name>/proposal.md` (L3 required, L2 optional)
8. **Write spec** — save to `docs/rein/changes/<name>/spec.md` and commit
9. **Spec self-review** — quick inline check for placeholders, contradictions, ambiguity, scope
10. **User reviews written spec** — ask user to review the spec file before proceeding
11. **Transition to implementation** — prompt user to run `/plan`

## Phase 1: Explore & Expand (Divergent → Convergent)

**Goal:** Take the raw idea and open it up, then converge on a direction.

### 1a. Explore & Understand

1. **Explore project context first** — check files, docs, recent commits
2. **Assess scope:** If the request describes multiple independent subsystems, flag this immediately. Help the user decompose into sub-projects before refining details.
3. **Restate the idea** as a crisp "How Might We" problem statement. This forces clarity on what's actually being solved.
4. **Ask 3-5 sharpening questions** — one at a time. Focus on:
   - Who is this for, specifically?
   - What does success look like?
   - What are the real constraints (time, tech, resources)?
   - What's been tried before?
   - Why now?

**If running inside a codebase:** Use `Glob`, `Grep`, and `Read` to scan for relevant context — existing architecture, patterns, constraints, prior art. Ground your variations in what actually exists.

### 1b. Generate Variations (Divergent)

Generate 5-8 idea variations using these lenses:
- **Inversion:** "What if we did the opposite?"
- **Constraint removal:** "What if budget/time/tech weren't factors?"
- **Audience shift:** "What if this were for [different user]?"
- **Combination:** "What if we merged this with [adjacent idea]?"
- **Simplification:** "What's the version that's 10x simpler?"
- **10x version:** "What would this look like at massive scale?"
- **Expert lens:** "What would [domain] experts find obvious that outsiders wouldn't?"

Push beyond what the user initially asked for. Create products people don't know they need yet.

Read `frameworks.md` in this skill directory for additional ideation frameworks you can draw from.

### 1c. Evaluate & Converge

After the user reacts to the variations, shift to convergent mode:

1. **Cluster** the ideas that resonated into 2-3 distinct directions. Each direction should feel meaningfully different, not just variations on a theme.

2. **Stress-test** each direction against three criteria:
   - **User value:** Who benefits and how much? Is this a painkiller or a vitamin?
   - **Feasibility:** What's the technical and resource cost? What's the hardest part?
   - **Differentiation:** What makes this genuinely different? Would someone switch from their current solution?

   Read `refinement-criteria.md` in this skill directory for the full evaluation rubric.

3. **Surface hidden assumptions.** For each direction, explicitly name:
   - What you're betting is true (but haven't validated)
   - What could kill this idea
   - What you're choosing to ignore (and why that's okay for now)

4. **Propose 2-3 approaches** with trade-offs and your recommendation. Lead with your recommended option and explain why.

**Be honest, not supportive.** If an idea is weak, say so with kindness. A good ideation partner is not a yes-machine.

## Phase 2: Specify

With the validated direction, produce two artifacts: **proposal.md** (why & scope) then **spec.md** (requirements, decisions, risks).

**Surface assumptions immediately.** Before writing any content, list what you're assuming:

```
ASSUMPTIONS I'M MAKING:
1. This is a web application (not native mobile)
2. Authentication uses session-based cookies (not JWT)
3. The database is PostgreSQL (based on existing Prisma schema)
4. We're targeting modern browsers only (no IE11)
→ Correct me now or I'll proceed with these.
```

Don't silently fill in ambiguous requirements. The spec's entire purpose is to surface misunderstandings *before* code gets written — assumptions are the most dangerous form of misunderstanding.

### Step 2a: Write proposal.md

The proposal captures **why** and **scope**. Required for L3 (`/feature`), optional for L2 (`/fix`).

Save to `docs/rein/changes/<name>/proposal.md`:

```markdown
# Proposal: [Project/Feature Name]

## Why
[One-sentence "How Might We" framing — the problem this solves and why now]

## What Changes
[The chosen direction and what will be different — 2-3 paragraphs max]

## Goals
- [Goal 1]
- [Goal 2]
- [Goal 3]

## Non-Goals
- [Thing 1] — [reason]
- [Thing 2] — [reason]
- [Thing 3] — [reason]

## Key Assumptions
- [ ] [Assumption 1 — how to test it]
- [ ] [Assumption 2 — how to test it]
- [ ] [Assumption 3 — how to test it]

## Open Questions
- [Question 1]
- [Question 2]
```

**The "Non-Goals" list is arguably the most valuable part.** Focus is about saying no to good ideas. Make the trade-offs explicit.

### Step 2b: Write spec.md

The spec captures **what and how** — requirements, decisions, and risks. If proposal.md exists, read it for Goals/Non-Goals context first.

Save to `docs/rein/changes/<name>/spec.md`:

```markdown
# Spec: [Project/Feature Name]

## Requirements

### Success Criteria
WHEN <condition> THEN <expected behavior>
- **TEST** `TestFunctionName` (optional)

### Commands
[Build, test, lint, dev — full commands]

### Project Structure
[Directory layout with descriptions]

### Code Style
[Example snippet + key conventions]

### Testing Strategy
[Framework, test locations, coverage requirements, test levels]

### Boundaries
- Always: [...]
- Ask first: [...]
- Never: [...]

## Decisions
**Decision:** [What was decided] — **Rationale:** [Why]
- **Decision:** [e.g., "Use session-based auth over JWT"] — **Rationale:** [e.g., "Simpler revocation, no token leakage risk for this app's threat model"]

## Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk description] | [High/Med/Low] | [Strategy] |
```

**Reframe instructions as success criteria.** When receiving vague requirements, translate them into concrete conditions. This applies to spec.md content only — proposal.md captures intent and scope, not testable criteria.

```
REQUIREMENT: "Make the dashboard faster"

REFRAMED SUCCESS CRITERIA:
- Dashboard LCP < 2.5s on 4G connection
- Initial data load completes in < 500ms
- No layout shift during load (CLS < 0.1)
→ Are these the right targets?
```

**Write scenarios in WHEN/THEN format.** Each success criterion should be expressed as a scenario:

```
WHEN <condition> THEN <expected behavior>
```

Optionally, link a scenario to the test function that verifies it:

```
WHEN user submits valid credentials THEN auth token is returned
- **TEST** `TestAuthJWT_ValidCredentials`
```

The `**TEST**` field is optional — add it when the test function name is known or planned. It creates traceability from requirement to verification.

**The "Non-Goals" list is arguably the most valuable part.** Focus is about saying no to good ideas. Make the trade-offs explicit.

## Spec Self-Review

After writing the spec document, look at it with fresh eyes:

1. **Placeholder scan:** Any "TBD", "TODO", incomplete sections, or vague requirements? Fix them.
2. **Internal consistency:** Do any sections contradict each other? Does the architecture match the feature descriptions?
3. **Scope check:** Is this focused enough for a single implementation plan, or does it need decomposition?
4. **Ambiguity check:** Could any requirement be interpreted two different ways? If so, pick one and make it explicit.

Fix any issues inline.

## User Review Gate

After the spec review loop passes, ask the user to review the written spec before proceeding:

> "Spec written and committed to `<path>`. Please review it and let me know if you want to make any changes before we start planning."

Wait for the user's response. If they request changes, make them and re-run the spec review loop. Only proceed once the user approves.

## Phase 3: Hand Off to Planning

The define phase produces specs ONLY. Task breakdown is handled by the `planning` skill invoked via `/plan`.

After the spec is complete, prompt the user:
> Spec complete. Run `/plan <name>` to break this into tasks, then `/do` to implement.

Do NOT generate task lists within the spec. Tasks belong in `docs/rein/changes/<name>/task.md`, generated by `/plan`.

## Keeping the Spec Alive

The spec is a living document, not a one-time artifact:

- **Update when decisions change** — If you discover the data model needs to change, update the spec first, then implement.
- **Update when scope changes** — Features added or cut should be reflected in the spec.
- **Commit the spec** — The spec belongs in version control alongside the code.
- **Reference the spec in PRs** — Link back to the spec section that each PR implements.

## Visual Companion

Decide per-question, not per-session. The test: **would the user understand this better by seeing it than reading it?**

**Use the browser** when the content itself is visual: UI mockups, architecture diagrams, side-by-side visual comparisons, design polish, spatial relationships.

**Use the terminal** when the content is text or tabular: requirements, conceptual choices, tradeoff lists, technical decisions, clarifying questions.

A question *about* a UI topic is not automatically a visual question. "What kind of wizard do you want?" is conceptual — use the terminal. "Which of these wizard layouts feels right?" is visual — use the browser.

**When visual questions will come up:** Offer the visual companion as its own message early in the session. Start with `rein visual start`, tell the user the URL, and use it for questions that benefit from visual presentation.

See `visual-thinking.md` in this skill directory for the full companion guide.

## Working in Existing Codebases

When ideating inside an existing project:

1. **Explore structure before proposing changes** — Use `Glob`, `Grep`, and `Read` to understand the current architecture, patterns, and constraints. Ground your variations in what actually exists.
2. **Include targeted improvements** — Where existing code directly affects the work being done, include fixes that make the change fit better (e.g., updating a shared utility that the new feature calls).
3. **Don't propose unrelated refactoring** — Stay focused on what serves the current goal. A "while we're here" refactor is scope creep, not efficiency.

## Key Principles

- **One question at a time** — Don't overwhelm with multiple questions
- **Multiple choice preferred** — Easier to answer than open-ended when possible
- **YAGNI ruthlessly** — Remove unnecessary features from all designs
- **Explore alternatives** — Always propose 2-3 approaches before settling
- **Incremental validation** — Present design, get approval before moving on
- **Surface assumptions** — Untested assumptions are the #1 killer of good ideas
- **Non-Goals list is mandatory** — Focus is about saying no to good ideas
- **Design for isolation and clarity** — Break systems into smaller units with one clear purpose, well-defined interfaces, and independent testability
- **Be flexible** — Go back and clarify when something doesn't make sense

## Common Rationalizations

| Rationalization | Reality |
|---|---|
| "This is simple, I don't need a spec" | Simple tasks don't need *long* specs, but they still need acceptance criteria. A two-line spec is fine. |
| "I'll write the spec after I code it" | That's documentation, not specification. The spec's value is in forcing clarity *before* code. |
| "The spec will slow us down" | A 15-minute spec prevents hours of rework. Waterfall in 15 minutes beats debugging in 15 hours. |
| "Requirements will change anyway" | That's why the spec is a living document. An outdated spec is still better than no spec. |
| "The user knows what they want" | Even clear requests have implicit assumptions. The spec surfaces those assumptions. |

## Anti-patterns

- Generating 20+ shallow variations instead of 5-8 considered ones
- Skipping the "who is this for" question
- No assumptions surfaced before committing to a direction
- Yes-machining weak ideas instead of pushing back with specificity
- Producing a plan without a "Non-Goals" list
- Ignoring existing codebase constraints when ideating inside a project
- Jumping to output without running the full process
- Unrelated refactoring proposals — stay focused on what serves the current goal
- Starting to write code without any written requirements
- Implementing features not mentioned in any spec or task list
- Making architectural decisions without documenting them

## Red Flags

- Starting to write code without any written requirements
- Asking "should I just start building?" before clarifying what "done" means
- Implementing features not mentioned in any spec or task list
- Making architectural decisions without documenting them
- Skipping the spec because "it's obvious what to build"

## Verification

After completing the define phase:

- [ ] A clear "How Might We" problem statement exists
- [ ] The target user and success criteria are defined
- [ ] Multiple directions were explored, not just the first idea
- [ ] Hidden assumptions are explicitly listed with validation strategies
- [ ] A "Non-Goals" list makes trade-offs explicit
- [ ] proposal.md is saved (L3 required, L2 optional)
- [ ] spec.md covers Requirements, Decisions, and Risks
- [ ] The human has reviewed and approved the spec
- [ ] Success criteria are specific and testable
- [ ] Boundaries (Always/Ask First/Never) are defined
- [ ] The spec is saved to `docs/rein/changes/<name>/spec.md`
- [ ] The spec has no placeholders, contradictions, or ambiguity

## Supporting Files

- **`frameworks.md`** — 7 ideation frameworks (SCAMPER, HMW, First Principles, JTBD, Constraint-Based, Pre-mortem, Analogous Inspiration)
- **`refinement-criteria.md`** — Evaluation rubric (User Value, Feasibility, Differentiation) + assumption audit + MVP scoping
- **`examples.md`** — 3 complete ideation session examples
- **`visual-thinking.md`** — Browser-based visual brainstorming companion guide
- **`spec-reviewer-prompt.md`** — Spec document reviewer subagent prompt template
