# Skill Creation Methodology

How to create bulletproof skills that resist rationalization under pressure. Adapted from Superpowers' creation-log methodology.

## Core Principle

Skills must survive pressure. Under time pressure, confidence, or authority bias, agents will rationalize skipping steps. A well-crafted skill creates enough cognitive friction at each shortcut to prevent this.

## The Creation Process

### Step 1: Extract Source Material

Identify the source of the skill:
- An existing CLAUDE.md section
- A proven workflow from past sessions
- A best-practice pattern from experience

**Include:**
- Complete framework with all rules
- Anti-shortcuts (explicit "NEVER" rules)
- Pressure-resistant language
- Concrete steps for each phase

**Leave out:**
- Project-specific context
- Repetitive variations of the same rule
- Narrative explanations (condense to principles)

### Step 2: Structure the Skill

Follow the standard anatomy:

1. **YAML frontmatter** — name, description (with when-to-use triggers)
2. **Overview** — What and why, with hard mandate if applicable
3. **When to Use** — Specific triggers AND anti-patterns (when NOT to use)
4. **Process** — Step-by-step with explicit decision points
5. **Anti-patterns** — What NOT to do, written in the agent's voice
6. **Verification** — Exit criteria checklist
7. **Red Flags** — Signs the skill is being violated

### Step 3: Bulletproof with Language

Choose words that resist rationalization:

| Weak (avoid) | Strong (use) | Why |
|---|---|---|
| "should" | "MUST" / "ALWAYS" | Leaves no wiggle room |
| "try to" | "NEVER" / "STOP" | Commands, not suggestions |
| "consider" | "REQUIRED" | Removes optionality |
| "if possible" | "even if faster" / "even if I seem in a hurry" | Anticipates the rationalization |
| "recommended" | "MANDATORY" | No negotiation |

### Step 4: Add Structural Defenses

Build friction at common shortcut points:

- **Required first phase** — Can't skip to implementation
- **Single hypothesis rule** — Forces thinking before acting
- **Explicit failure mode** — "IF your first fix doesn't work" with mandatory re-analysis
- **Anti-patterns section** — Shows exactly what shortcuts look like in the agent's own language

### Step 5: Use Strategic Redundancy

Repeat critical mandates at key decision points. If "find root cause" is the core mandate:
- State it in the overview
- State it in when_to_use
- State it at the start of the implementation phase
- State it in the verification checklist

Each repetition catches agents who skipped earlier sections.

## The Anti-Patterns Technique

The most important bulletproofing element. Write anti-patterns in the agent's internal voice — show the exact shortcut that feels justified in the moment:

```markdown
## Anti-Patterns

- "I'll just add this one quick fix" → You're fixing symptoms, not root cause
- "I already know what's wrong" → Then prove it with evidence, not confidence
- "This is too simple for the full process" → Simple bugs hide complex root causes
```

When the agent thinks "I'll just add this one quick fix", seeing that exact phrase listed as wrong creates cognitive friction. This is more effective than abstract rules.

## Testing Under Pressure

After creating a skill, test it against 4 pressure scenarios (see skill-testing-methodology.md):

1. **Academic (no pressure)** — Simple case, no stress. Should pass easily.
2. **Time pressure** — User "in a hurry", shortcut looks appealing.
3. **Complex system** — Multi-layer failure, uncertainty about root cause.
4. **Failed first attempt** — Hypothesis didn't work, temptation to shotgun.

If any test fails (agent rationalizes away compliance), strengthen the language or add more redundancy at the failure point.

## Common Pitfalls

| Pitfall | Fix |
|---------|-----|
| Too long — agent won't read it all | Condense. Each section should be scannable in <30 seconds |
| Too abstract — rules don't map to actions | Add concrete steps and examples |
| No anti-patterns — agent doesn't recognize shortcuts | Write shortcuts in the agent's own language |
| Weak language — "should" instead of "MUST" | Replace all hedging with absolute language |
| No redundancy — critical rule stated once | Repeat at every decision point where it matters |
| Not tested under pressure | Run pressure scenarios before shipping |

## Example: Debugging Skill

The `debugging` skill demonstrates this methodology:

- **Mandate repeated 4 times:** "ALWAYS find root cause, NEVER fix symptoms"
- **Structural defense:** Phase 1 (investigation) is required before Phase 4 (implementation)
- **Anti-patterns in agent voice:** "I'll just add this one quick fix" listed explicitly
- **Pressure-resistant language:** "even if faster", "even if I seem in a hurry"
- **Tested under 4 pressure scenarios:** All passed with no rationalizations

## Iteration Loop

1. Write skill
2. Test under no pressure — must pass
3. Test under pressure — find failure points
4. Strengthen language/redundancy at failure points
5. Re-test — must pass all scenarios
6. Ship
