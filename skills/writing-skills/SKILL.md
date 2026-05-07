---
name: writing-skills
description: Use when creating new skills, editing existing skills, or verifying skills work before deployment
---

# Writing Skills

Skill creation is TDD applied to process documentation. A skill that hasn't been tested is a skill that doesn't work.

## The Iron Law

```
NO SKILL WITHOUT A FAILING TEST FIRST
```

If you can't describe a scenario where the agent would fail without this skill, you don't need the skill.

## What is a Skill

A skill is a **reusable technique, pattern, tool, or reference guide** that changes agent behavior. It is NOT a narrative, NOT documentation, and NOT a story.

Skills work because they constrain behavior. A good skill makes the right thing easy and the wrong thing hard.

## Skill Types

| Type | Purpose | Test Approach | Example |
|------|---------|---------------|---------|
| **Technique** | Step-by-step process for a specific task | Pressure scenarios against each step | TDD, debugging, code review |
| **Pattern** | Reusable solution to a recurring problem | Apply to 2+ contexts, verify consistency | Error handling, API design |
| **Reference** | Lookup material for quick access | Verify agent finds correct info | Security checklist, style guide |

## TDD Mapping for Skills

| TDD Concept | Skill Equivalent |
|-------------|-----------------|
| Test case | Pressure scenario |
| Production code | SKILL.md content |
| RED | Agent violates baseline without skill |
| GREEN | Agent complies with skill |
| REFACTOR | Close loopholes, tighten language |

## When to Create a Skill

**Create when:**
- An agent repeatedly makes the same mistake
- A process has non-obvious steps agents skip
- A compliance rule needs enforcement (not just documentation)
- Multiple agents need consistent behavior

**Don't create when:**
- The behavior is obvious (agents already do it)
- A one-line instruction suffices
- It's reference material an agent can look up
- You can't describe a failing scenario

## Directory Structure

```
skills/
  skill-name/
    SKILL.md           # Required: The skill definition
    supporting-file.md # Optional: Reference material loaded on demand
```

Flat namespace. SKILL.md is required. Supporting files only when the main file would exceed 500 lines.

## SKILL.md Structure

### Frontmatter (Required)

```yaml
---
name: skill-name-with-hyphens
description: Guides agents through [task/workflow]. Use when [specific trigger conditions].
---
```

**Rules:**
- `name`: Lowercase, hyphen-separated. Must match the directory name.
- `description`: Start with what the skill does, then "Use when" triggers. Max 1024 characters.
- The description is the **search index** — it must tell the agent when to activate, not summarize the workflow.

### Standard Sections

1. **Overview** — What this skill does and why (1-2 sentences)
2. **When to Use** — Triggering conditions and exclusions
3. **Core Process** — The main workflow with numbered steps
4. **Common Rationalizations** — Table of excuses vs. reality
5. **Red Flags** — Observable signs the skill is being violated
6. **Verification** — Exit criteria checklist

## Claude Search Optimization (CSO)

Agents discover skills by reading descriptions. Optimize for discovery:

1. **Rich description field** — Include when-to-use triggers, not just what-it-does
2. **Keyword coverage** — Use terms agents would search for
3. **Descriptive naming** — Name from the agent's perspective, not the author's
4. **Cross-reference** — Mention related skills by name

**Critical rule:** description = when to use, NOT what the skill does. If the description contains process steps, the agent may follow the summary instead of reading the full skill.

## RED-GREEN-REFACTOR for Skills

### RED: Prove the Need

1. Create a pressure scenario (see `testing-skills-with-subagents.md`)
2. Run WITHOUT the skill — document the agent's choices and rationalizations verbatim
3. If the agent naturally does the right thing, you don't need the skill

### GREEN: Write the Minimal Skill

1. Write only enough to address the specific baseline failures
2. The skill should make the right action obvious and the wrong action hard
3. Test with the same pressure scenario

### REFACTOR: Close Loopholes

1. Run new pressure scenarios trying to break the skill
2. For each rationalization found:
   - Add explicit negation
   - Add rationalization table entry
   - Add red flag entry
   - Update description

## Anti-Patterns

- **Narrative example:** Skills are instructions, not stories. "Imagine you're..." is a red flag.
- **Multi-language dilution:** One excellent example beats many mediocre ones.
- **Code in flowcharts:** Flowcharts show decision points, not implementation details.
- **Generic labels:** "Handle error" is useless. "Retry 3 times with exponential backoff, then log and continue" is a skill.

## Flowchart Usage

Only for non-obvious decision points and process loops. Never for:
- Reference material
- Linear instructions
- Simple if/else logic

See `graphviz-conventions.dot` for DOT style conventions.

## Skill Creation Checklist

### RED Phase
- [ ] Pressure scenario written (what pressure, what's the wrong choice)
- [ ] Baseline test run (agent behavior WITHOUT skill documented)
- [ ] Need confirmed (agent made the wrong choice under pressure)

### GREEN Phase
- [ ] Minimal skill written addressing baseline failures
- [ ] Skill tested with same pressure scenario
- [ ] Agent now makes the right choice

### REFACTOR Phase
- [ ] Additional pressure scenarios tested
- [ ] Rationalizations captured and addressed
- [ ] Loopholes closed with explicit negations
- [ ] Red flags list updated
- [ ] Description updated for better discovery

### Quality
- [ ] No narrative or story-like content
- [ ] Code examples are specific and correct
- [ ] Flowcharts only for non-obvious decisions
- [ ] Supporting files loaded on demand (progressive disclosure)

### Deployment
- [ ] File location follows naming conventions
- [ ] Frontmatter is valid YAML with name + description
- [ ] Description optimized for agent discovery (CSO)

## Supporting Files

- `anthropic-best-practices.md` — Anthropic's official skill authoring best practices
- `persuasion-principles.md` — Research-backed persuasion principles for skill design
- `testing-skills-with-subagents.md` — Complete testing methodology for skills
- `graphviz-conventions.dot` — Graphviz DOT style guide for process diagrams
