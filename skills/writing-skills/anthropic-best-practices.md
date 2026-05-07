# Anthropic Best Practices for Skill Authoring

Best practices distilled from Anthropic's official skill authoring guidelines.

## Core Principles

1. **Concise is key** — Context window is a public good. Every token you spend on a skill is a token the agent can't spend on the actual task.
2. **Set appropriate degrees of freedom** — High (creative work), Medium (structured tasks), Low (compliance enforcement). Match the constraint level to the task.
3. **Test with all models** — A skill that only works on Opus is a fragile skill. Test on Sonnet and Haiku too.

## Skill Structure

### Frontmatter

```yaml
---
name: skill-name          # Max 64 chars, lowercase-hyphen
description: ...           # Max 1024 chars, include "Use when" triggers
---
```

### Naming Conventions

- Use gerund form (gerund = verb-ing): `writing-skills`, `debugging`, `deploying`
- Name from the agent's perspective: what would the agent call this?

### Writing Effective Descriptions

- Write in third person
- Be specific about when to use
- Include key terms agents search for
- Do NOT summarize the workflow in the description

**Bad:** "A skill for code review that checks correctness, readability, and security"
**Good:** "Use when reviewing code changes for correctness, readability, architecture, security, or performance issues. Triggers: after writing code, before committing, during PR review."

## Progressive Disclosure Patterns

Keep SKILL.md under 500 lines. Use supporting files for deeper content.

### Pattern 1: High-Level Guide with References

SKILL.md contains the workflow and decision points. Supporting files contain detailed reference material loaded on demand.

### Pattern 2: Domain-Specific Organization

SKILL.md contains the common core. Domain-specific variations go in supporting files.

### Pattern 3: Conditional Details

SKILL.md contains the main flow. Conditional branches and edge cases go in supporting files, loaded only when the agent hits that branch.

**Keep references one level deep from SKILL.md.** Deep nesting defeats progressive disclosure.

## Content Guidelines

- **Avoid time-sensitive information** — Version numbers, pricing, current events
- **Use consistent terminology** — Same term for same concept throughout
- **One excellent example beats many mediocre ones** — Show the ideal, not variations

## Common Patterns

### Template Pattern
Provide a fill-in-the-blank template that forces the agent to follow the structure.

### Examples Pattern
Input/output pairs showing what good looks like. The agent generalizes from examples.

### Conditional Workflow Pattern
Decision tree with clear branching. Each branch is a self-contained sub-flow.

## Evaluation and Iteration

1. **Build evaluations first** — Know how you'll test before you write the skill
2. **Develop iteratively with Claude** — Claude A creates, Claude B tests
3. **Observe how Claude navigates skills** — Where does it skip? Where does it get confused?

## Anti-Patterns

- **Avoid Windows-style paths** in examples — Use forward slashes, relative paths
- **Avoid offering too many options** — The agent should follow one clear path
- **Don't nest references deeply** — One level from SKILL.md only
- **Don't duplicate content** — Reference and link instead

## Technical Notes

### YAML Frontmatter Requirements
- Must be valid YAML
- `name` and `description` are required
- No trailing whitespace
- No tab characters

### Token Budgets
- SKILL.md: Target under 500 lines
- Description: Max 1024 characters
- Name: Max 64 characters
- Supporting files: No hard limit, but each should be loadable in a single context window

## Checklist for Effective Skills

### Core Quality
- [ ] Description tells when to use, not what it does
- [ ] No narrative or story-like content
- [ ] Specific instructions, not vague advice
- [ ] Anti-rationalization table for skip-worthy steps
- [ ] Verification checklist at the end

### Code and Scripts
- [ ] Code examples are specific and runnable
- [ ] No assumed tools without checking availability
- [ ] Error handling in utility scripts
- [ ] Dependencies documented

### Testing
- [ ] Pressure scenarios defined
- [ ] Baseline behavior without skill documented
- [ ] Skill tested with multiple pressure types
- [ ] Works across model tiers (Sonnet, Haiku)
