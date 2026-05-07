# Persuasion Principles for Skill Design

Research-backed principles for writing skills that agents actually follow. Based on Meincke et al. (2025), N=28,000 AI conversations, where persuasion more than doubled compliance (33% to 72%).

## The Seven Principles

### 1. Authority — Imperative Language

Use non-negotiable framing for mandatory steps.

**Apply with:**
- "MUST", "NEVER", "ALWAYS" for critical rules
- Non-negotiable framing: "This is not optional"
- No hedging: not "consider", but "do this"

**Example:**
```
NEVER skip the failing test. If the test doesn't fail, you haven't proven anything.
```

### 2. Commitment — Require Announcements

Force explicit choices before proceeding.

**Apply with:**
- Require agents to announce what they're about to do
- Use tracking mechanisms (checkboxes, status markers)
- Force explicit "I will / I won't" decisions at branching points

**Example:**
```
Before starting implementation, state which task you're working on and what the acceptance criteria are.
```

### 3. Scarcity — Time-Bound Requirements

Create urgency through sequential dependencies.

**Apply with:**
- Sequential dependencies (can't do B until A is done)
- Limited opportunities (this check only happens once)
- Time-bound requirements (verify NOW, not later)

### 4. Social Proof — Universal Patterns

Show that compliance is the standard, not the exception.

**Apply with:**
- "Every production bug was once a 'simple' fix that skipped testing"
- Reference universal failure modes
- Show what experienced practitioners do

### 5. Unity — Collaborative Language

Frame the skill as a shared goal, not imposed rules.

**Apply with:**
- "We" language: "We write tests first because..."
- Shared goals: "We both want working code"
- Collaborative framing, not adversarial

### 6. Reciprocity — Use Sparingly

Rarely needed. When used, frame compliance as a benefit to the agent.

**Apply with:**
- "Following this process saves you time — you debug less"
- Use only when the benefit is genuine and immediate

### 7. Liking — DON'T USE for Compliance

Creates sycophancy. Agents will agree without actually complying.

**Do NOT:**
- Frame rules as suggestions to be agreeable
- Use "you might want to" for mandatory steps
- Soften mandatory language to be likable

## Principle Combinations by Skill Type

| Skill Type | Use Principles | Avoid |
|------------|---------------|-------|
| Discipline-enforcing | Authority + Commitment + Scarcity | Liking |
| Guidance | Social Proof + Unity | Authority (too rigid) |
| Collaborative | Unity + Social Proof | Authority (adversarial) |
| Reference | Scarcity (when to load) | Authority (not a rule) |

## Why This Works

- **Bright-line rules** reduce rationalization — "NEVER" is harder to bend than "try not to"
- **Implementation intentions** create automatic behavior — "When X happens, I do Y"
- **LLMs are parahuman** — They respond to persuasion patterns similarly to humans, but with less resistance to authority framing

## Ethical Use

**Legitimate:** Using authority for compliance steps that prevent bugs, security issues, or data loss.

**Illegitimate:** Using persuasion to make agents accept incorrect outputs, skip safety checks, or ignore errors.
