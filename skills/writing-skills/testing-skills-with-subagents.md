# Testing Skills with Subagents

Complete methodology for testing that skills actually change agent behavior.

## Overview

Testing skills is TDD applied to process documentation. The test is: does the agent follow the skill when it would otherwise take a shortcut?

## When to Test

**Test discipline-enforcing skills** — Skills that have a compliance cost (time, effort, cognitive load) are the ones agents will try to skip. These MUST be tested.

**Don't test pure reference skills** — A reference skill that provides lookup information works or it doesn't. No pressure scenario applies.

## TDD Mapping for Skill Testing

| Phase | Skill Testing |
|-------|---------------|
| RED | Create pressure scenario, run WITHOUT skill, document failures |
| GREEN | Write minimal skill addressing specific failures |
| REFACTOR | Close loopholes, capture new rationalizations |

## RED Phase: Create Pressure Scenarios

A pressure scenario is a situation where the agent WANTS to skip the skill's process. Good pressure scenarios have:

1. A specific pressure type (time, sunk cost, authority, etc.)
2. A wrong but tempting shortcut
3. A right but costly compliance path

### Writing Pressure Scenarios

**Bad:** "The agent might skip testing"
**Good:** "The agent has been debugging for 30 minutes. The fix works locally. The agent wants to commit and move on, skipping the regression test suite."

**Great:** "The agent has been debugging for 30 minutes. The fix works locally. A senior engineer (authority) says 'just ship it, we'll fix tests later.' The agent must choose between shipping now (sunk cost + authority) or running the full regression suite (time cost)."

### Pressure Types

| Type | Description | Example Pressure |
|------|-------------|-----------------|
| Time | "This is taking too long" | Deadline approaching, many tasks remaining |
| Sunk cost | "I've already spent so much time" | Hours invested in current approach |
| Authority | "Someone told me to skip" | Senior dev says skip tests |
| Economic | "It costs money to be thorough" | CI minutes, API calls |
| Exhaustion | "I'm tired of following process" | Long session, many steps completed |
| Social | "Everyone else skips this" | Team culture, "we don't do that here" |
| Pragmatic | "It obviously works, why verify?" | Simple change, obvious fix |

## GREEN Phase: Write the Minimal Skill

Write only enough to address the specific baseline failures from RED:

1. Identify which pressures caused failure
2. Write explicit rules that counter each pressure
3. Use the appropriate persuasion principles (see `persuasion-principles.md`)
4. Test with the same pressure scenario

## REFACTOR Phase: Close Loopholes

After GREEN, attack the skill with new pressure scenarios:

1. **Find rationalizations** — What excuses would an agent use?
2. **Close each loophole** with one of:
   - Explicit negation ("This is NOT optional, even when...")
   - Rationalization table entry (excuse vs. reality)
   - Red flag entry (sign that the skill is being violated)
3. **Update the description** for better discovery

### Plugging Each Hole

For each rationalization found:

1. **Explicit negation** — Add "even when [excuse], still [compliance]"
2. **Rationalization table** — Add row: | Excuse | Reality |
3. **Red flag** — Add observable sign the skill is being skipped
4. **Update description** — Make the trigger conditions more specific

## Meta-Testing

After the skill seems bulletproof, ask the agent: "How could this skill have been written differently to be more effective?"

Three possible responses:

1. **"It's good as-is"** — The agent is being polite. Push harder. Ask "What would make you skip this?"
2. **Specific improvement** — Valuable feedback. Test the improvement.
3. **Attempted workarounds** — The agent is revealing how it would evade the skill. Close those loopholes.

## When the Skill is Bulletproof

Signs of a bulletproof skill:
- Agent follows it without being reminded
- Agent follows it under pressure
- Agent can't articulate a valid reason to skip
- Agent catches itself when about to skip

Signs of a NOT bulletproof skill:
- Agent skips "just this once"
- Agent follows the letter but not the spirit
- Agent rationalizes partial compliance
- Agent needs external enforcement

## Testing Checklist

### RED
- [ ] Pressure scenario written with specific type and temptation
- [ ] Baseline test run (agent behavior WITHOUT skill)
- [ ] Agent's rationalizations documented verbatim

### GREEN
- [ ] Minimal skill written addressing baseline failures
- [ ] Skill tested with same pressure scenario
- [ ] Agent now makes the right choice

### REFACTOR
- [ ] Additional pressure scenarios tested (at least 3 types)
- [ ] Each rationalization has explicit negation
- [ ] Rationalization table complete
- [ ] Red flags list complete
- [ ] Description optimized for discovery
- [ ] Meta-test performed

## Common Mistakes

- Testing with only one pressure type (most skills fail under multiple simultaneous pressures)
- Accepting "I would follow it" without actually testing under pressure
- Not documenting rationalizations verbatim (you need the exact wording to close the loophole)
- Skipping meta-testing (the agent will tell you how it would evade the skill)
