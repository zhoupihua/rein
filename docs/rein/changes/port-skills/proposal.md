# Proposal: Port Skills from Matt Pocock's Skills Project

## Why
How might we enhance alloy's skill set by selectively porting high-value concepts from Matt Pocock's "Skills For Real Engineers" project — specifically structured debugging feedback loops, interview-based alignment, and enriched TDD references — while maintaining alloy's workflow philosophy and avoiding duplication with existing skills?

## What Changes
Three targeted enhancements to alloy's static content:

1. **Enhance the debugging skill** by surgically inserting diagnose's key insights: a "Build a Feedback Loop" phase as the first step (the #1 insight — a fast deterministic loop makes the bug 90% fixed), ranked falsifiable hypotheses between Localize and Reduce, tagged `[DEBUG-xxxx]` instrumentation, and a post-mortem question. Current strengths (triage checklist, error-specific patterns, safe fallbacks) move to supporting reference files to keep the main SKILL.md focused.

2. **Create a new grill skill** for interview-based alignment before implementation. A relentless one-at-a-time questioning technique that challenges plans and designs, with recommended answers for each question. When a question can be answered by exploring the codebase, the skill explores instead of asking. Follows alloy conventions (Iron Law, Rationalizations, Red Flags, Verification).

3. **Enrich the tdd skill** with three reference files ported from the skills project: `deep-modules.md` (deep vs shallow module design), `mocking.md` (when to mock, designing for mockability), and `refactoring.md` (refactor candidates checklist after TDD cycle).

## Goals
- Debugging skill leads with the feedback loop insight — the single most impactful debugging concept from the skills project
- Grill skill fills the alignment gap — alloy has define for specification but no tool for stress-testing plans through interview
- TDD skill has richer reference material for interface design and mock discipline
- All ported content follows alloy's conventions (Iron Law, Common Rationalizations table, Red Flags, Verification checklist)
- No duplication with existing skills or hooks

## Non-Goals
- Porting improve-codebase-architecture — valuable but scope creep for this feature; can be a separate port later
- Porting prototype — interesting concept but requires significant adaptation; defer
- Porting zoom-out — too small for a standalone skill; could be a section in an existing skill
- Porting triage/to-issues — requires issue tracker integration that alloy doesn't have
- Porting grill-with-docs (CONTEXT.md system) — needs project-wide adoption, too heavy for now
- Porting caveman — too niche
- Changing Go CLI code or hook implementations — this is a static content change only
- Modifying the install scripts — no new file types or directories are introduced

## Key Assumptions
- [ ] The feedback loop concept is the most valuable part of diagnose and deserves top billing in debugging — testable by whether agents actually build loops before hypothesizing
- [ ] Grill-me's one-at-a-time question pattern works well with alloy's skill structure — testable by using it in a real planning session
- [ ] TDD reference files will be discovered and used by agents — testable by checking if agents reference them during implementation
- [ ] Moving error-specific patterns and safe fallbacks to reference files won't reduce their effectiveness — testable by verifying agents still apply them

## Open Questions
- Should the grill skill be placed in the DEFINE phase (alongside define) or in its own phase? It's used before implementation to stress-test plans.
- Should the debugging skill's "Guard Against Recurrence" step be renamed to "Cleanup + Post-Mortem" to match diagnose's Phase 6 naming?
