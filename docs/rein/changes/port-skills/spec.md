# Spec: Port Skills from Matt Pocock's Skills Project

## Requirements

### Success Criteria

WHEN a bug is encountered THEN the debugging skill directs the agent to build a feedback loop first, before reproducing or hypothesizing
- **TEST** `TestDebuggingSkill_StartsWithFeedbackLoop`

WHEN a feedback loop cannot be built THEN the agent stops and asks the user for help instead of proceeding to hypothesize without a loop
- **TEST** `TestDebuggingSkill_StopsWithoutLoop`

WHEN the agent has a working feedback loop THEN it generates 3-5 ranked falsifiable hypotheses before testing any of them
- **TEST** `TestDebuggingSkill_RankedHypotheses`

WHEN the agent instruments code for debugging THEN it tags all debug logs with a unique `[DEBUG-xxxx]` prefix
- **TEST** `TestDebuggingSkill_TaggedInstrumentation`

WHEN a bug fix is complete THEN the agent asks "what would have prevented this bug?" as a post-mortem
- **TEST** `TestDebuggingSkill_PostMortem`

WHEN the user invokes the grill skill THEN the agent asks relentless one-at-a-time questions about a plan or design, providing recommended answers
- **TEST** `TestGrillSkill_OneAtATime`

WHEN a grill question can be answered by exploring the codebase THEN the agent explores the codebase instead of asking
- **TEST** `TestGrillSkill_ExploresInsteadOfAsking`

WHEN the tdd skill is invoked THEN the agent can reference deep-modules.md, mocking.md, and refactoring.md for interface design and mock discipline guidance
- **TEST** `TestTDDSkill_ReferenceFilesAvailable`

### Commands

```bash
go build ./...                        # Verify no Go code changes break build
go test ./...                         # Run all tests
```

This is a static content change only. No Go code changes required.

### Project Structure

New files:
```
skills/
  debugging/
    SKILL.md                          # Enhanced with feedback loop + hypotheses + tagged instrumentation + post-mortem
    feedback-loop-strategies.md       # NEW: 10 strategies for building feedback loops (from diagnose)
    hypothesis-framework.md           # NEW: Ranked falsifiable hypotheses framework (from diagnose)
    error-triage-patterns.md          # NEW: Error-specific patterns moved from main SKILL.md
  grill/
    SKILL.md                          # NEW: Interview-based alignment skill
  tdd/
    SKILL.md                          # Unchanged
    deep-modules.md                   # NEW: Deep vs shallow module design (from skills/tdd)
    mocking.md                        # NEW: When to mock, designing for mockability (from skills/tdd)
    refactoring.md                    # NEW: Refactor candidates checklist (from skills/tdd)
```

Modified files:
```
skills/debugging/SKILL.md             # Restructured with diagnose concepts
```

### Code Style

All ported content follows alloy's skill conventions:
- Frontmatter with `name` and `description` (max 1024 chars, includes trigger conditions)
- Iron Law section
- Common Rationalizations table
- Red Flags list
- Verification checklist
- Supporting reference files loaded on demand (progressive disclosure)

### Testing Strategy

Manual verification:
1. Read each modified/new SKILL.md and verify it follows alloy conventions
2. Verify no content is lost from the original debugging skill (moved to reference files)
3. Verify grill SKILL.md is complete and self-contained
4. Verify tdd reference files are linked from the main SKILL.md

### Boundaries

- Always: Follow alloy skill conventions (Iron Law, Rationalizations, Red Flags, Verification)
- Always: Port concepts, not verbatim text — adapt to alloy's voice and structure
- Always: Keep SKILL.md files focused (< 500 lines); move detailed content to reference files
- Ask first: If a concept from the skills project conflicts with an existing alloy concept
- Never: Modify Go CLI code, hooks, or install scripts
- Never: Port skills that duplicate existing alloy functionality (git-guardrails, to-prd, setup-*)
- Never: Add new directories or file types beyond what alloy already supports

## Decisions

**Decision:** Place grill skill in the DEFINE phase alongside define — **Rationale:** Grill is used to stress-test plans before implementation, which is the same phase as define. It's a complement to define, not a separate phase.

**Decision:** Restructure debugging skill to lead with "Build a Feedback Loop" as Step 1, keeping the current triage flow but inserting diagnose concepts at natural points (current steps become 2-7) — **Rationale:** The feedback loop is the #1 insight from diagnose and deserves top billing, but the current triage flow (Reproduce → Localize → Reduce → Fix → Guard → Verify) is solid and shouldn't be thrown away. The hybrid preserves the best of both.

**Decision:** Move error-specific patterns and safe fallback patterns from debugging SKILL.md to `error-triage-patterns.md` reference file — **Rationale:** The main SKILL.md should focus on the debugging process. Error-specific patterns are lookup material that agents can reference on demand. This keeps the main file focused while preserving the content.

**Decision:** Port tdd reference files (deep-modules, mocking, refactoring) as-is with minimal adaptation — **Rationale:** These are concise reference documents (10-30 lines each) that need no structural adaptation. The content is language-agnostic and directly useful.

**Decision:** Don't port the HITL loop template script — **Rationale:** It's a bash-specific artifact that doesn't fit alloy's skill structure. The concept of structured human-in-the-loop debugging can be described in the feedback-loop-strategies.md reference file.

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Restructuring debugging skill breaks agent muscle memory for current flow | Med | Keep step numbering consistent; current steps become 2-7 with "Build Feedback Loop" as step 1 |
| Grill skill is too brief (5 lines in original) to fill alloy's convention structure | Low | Expand with alloy-specific sections: Iron Law ("No implementation without alignment"), Common Rationalizations, Red Flags, Verification |
| TDD reference files go undiscovered by agents | Low | Reference them explicitly from tdd SKILL.md's "Supporting Techniques" section |
| Moving error patterns to reference file makes them less discoverable | Med | Add explicit reference in SKILL.md's "Supporting Techniques" section; include a "See error-triage-patterns.md for error-specific patterns" note at the relevant step |
