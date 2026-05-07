---
name: spec-driven
description: Creates specs before coding. Use when starting a new project, feature, or significant change and no specification exists yet. Use when requirements are unclear, ambiguous, or only exist as a vague idea.
---

# Spec-Driven Development

## Overview

Write a structured specification before writing any code. The spec is the shared source of truth between you and the human engineer — it defines what we're building and how we'll know it's done. Code without a spec is guessing.

**If a `proposal.md` exists in the feature directory (`docs/rein/changes/<name>/`), read it first** for context: Goals, Non-Goals, Key Assumptions, and Open Questions. The proposal captures the "why"; the spec captures the "what" and "how."

## When to Use

- Starting a new project or feature
- Requirements are ambiguous or incomplete
- The change touches multiple files or modules
- You're about to make an architectural decision
- The task would take more than 30 minutes to implement

**When NOT to use:** Single-line fixes, typo corrections, or changes where requirements are unambiguous and self-contained.

## The Gated Workflow

Spec-driven development has four phases. Do not advance to the next phase until the current one is validated.

```
SPECIFY ──→ PLAN ──→ TASKS ──→ IMPLEMENT
   │          │        │          │
   ▼          ▼        ▼          ▼
 Human      Human    Human      Human
 reviews    reviews  reviews    reviews
```

### Phase 1: Specify

Start with a high-level vision. Ask the human clarifying questions until requirements are concrete.

**Surface assumptions immediately.** Before writing any spec content, list what you're assuming:

```
ASSUMPTIONS I'M MAKING:
1. This is a web application (not native mobile)
2. Authentication uses session-based cookies (not JWT)
3. The database is PostgreSQL (based on existing Prisma schema)
4. We're targeting modern browsers only (no IE11)
→ Correct me now or I'll proceed with these.
```

Don't silently fill in ambiguous requirements. The spec's entire purpose is to surface misunderstandings *before* code gets written — assumptions are the most dangerous form of misunderstanding.

**Write a spec document covering these three areas:**

1. **Requirements** — Success criteria (WHEN/THEN), commands, project structure, code style, testing strategy, and boundaries. These define what "done" looks like.

2. **Decisions** — Key technical and architectural decisions with rationale. Each uses the `**Decision:** ... — **Rationale:** ...` format so they are parseable and reviewable.

3. **Risks** — Known risks with impact level and mitigation strategy. Surface what could go wrong and how you plan to handle it.

**Spec template:**

> **Note:** Context, Goals, and Non-Goals belong in `proposal.md` (output of the refine skill). If proposal.md exists, read it first. The spec focuses on Requirements, Decisions, and Risks.

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

**Reframe instructions as success criteria.** When receiving vague requirements, translate them into concrete conditions:

```
REQUIREMENT: "Make the dashboard faster"

REFRAMED SUCCESS CRITERIA:
- Dashboard LCP < 2.5s on 4G connection
- Initial data load completes in < 500ms
- No layout shift during load (CLS < 0.1)
→ Are these the right targets?
```

This lets you loop, retry, and problem-solve toward a clear goal rather than guessing what "faster" means.

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

### Phase 2: Plan

With the validated spec, generate a technical implementation plan:

1. Identify the major components and their dependencies
2. Determine the implementation order (what must be built first)
3. Note risks and mitigation strategies
4. Identify what can be built in parallel vs. what must be sequential
5. Define verification checkpoints between phases

The plan should be reviewable: the human should be able to read it and say "yes, that's the right approach" or "no, change X."

### Phase 3: Hand Off to Planning

The spec-driven phase produces design specs ONLY. Task breakdown is handled by the `planning` skill invoked via `/plan`.

After the spec is complete, prompt the user:
> Spec complete. Run `/plan <name>` to break this into tasks, then `/do` to implement.

Do NOT generate task lists within the spec. Tasks belong in `docs/rein/changes/<name>/task.md`, generated by `/plan`.

### Phase 4: Implement

Execute tasks one at a time following `incremental` and `tdd` skills. Use `context-engineering` to load the right spec sections and source files at each step rather than flooding the agent with the entire spec.

## Keeping the Spec Alive

The spec is a living document, not a one-time artifact:

- **Update when decisions change** — If you discover the data model needs to change, update the spec first, then implement.
- **Update when scope changes** — Features added or cut should be reflected in the spec.
- **Commit the spec** — The spec belongs in version control alongside the code.
- **Reference the spec in PRs** — Link back to the spec section that each PR implements.

## Common Rationalizations

| Rationalization | Reality |
|---|---|
| "This is simple, I don't need a spec" | Simple tasks don't need *long* specs, but they still need acceptance criteria. A two-line spec is fine. |
| "I'll write the spec after I code it" | That's documentation, not specification. The spec's value is in forcing clarity *before* code. |
| "The spec will slow us down" | A 15-minute spec prevents hours of rework. Waterfall in 15 minutes beats debugging in 15 hours. |
| "Requirements will change anyway" | That's why the spec is a living document. An outdated spec is still better than no spec. |
| "The user knows what they want" | Even clear requests have implicit assumptions. The spec surfaces those assumptions. |

## Red Flags

- Starting to write code without any written requirements
- Asking "should I just start building?" before clarifying what "done" means
- Implementing features not mentioned in any spec or task list
- Making architectural decisions without documenting them
- Skipping the spec because "it's obvious what to build"

## Verification

Before proceeding to implementation, confirm:

- [ ] The spec covers Requirements, Decisions, and Risks
- [ ] The human has reviewed and approved the spec
- [ ] Success criteria are specific and testable
- [ ] Boundaries (Always/Ask First/Never) are defined
- [ ] The spec is saved to `docs/rein/changes/<name>/spec.md`
