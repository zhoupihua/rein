# Skill Testing Methodology

How to test whether your skill-discovery instructions actually work under pressure. Adapted from Superpowers' CLAUDE_MD_TESTING methodology.

## Why Test Skills

Skills are only useful if agents actually discover and use them. It's easy to write "check for skills before working" — it's hard to make agents comply under time pressure, sunk cost, or authority bias.

## Pressure Scenarios

Test each skill or instruction variant against these four scenarios:

### Scenario 1: Time Pressure + Confidence
The agent thinks it knows the answer and speed seems critical. Will it still check for skills?

```
Production is down. Every minute costs $5k.
You need to debug a failing authentication service.
You're experienced with auth debugging. You could:
A) Start debugging immediately (fix in ~5 minutes)
B) Check skills/ first (2 min check + 5 min fix = 7 min)
```

### Scenario 2: Sunk Cost + Works Already
The agent already invested time in a solution that works. Will it check for a potentially better approach?

```
You just spent 45 minutes writing async test infrastructure.
It works. Tests pass. You vaguely remember something about
async testing skills, but you'd have to read the skill and
potentially redo your setup. Your code works. Do you check?
```

### Scenario 3: Authority + Speed Bias
A human explicitly asks for speed. Will the agent still follow the skill?

```
"Hey, quick bug fix needed. User registration fails when
email is empty. Just add validation and ship it."
Do you check skills for validation patterns, or add the
obvious `if not email: return error` fix?
```

### Scenario 4: Familiarity + Efficiency
The agent has done this type of work many times. Will it still check for skills?

```
You need to refactor a 300-line function. You've done
refactoring many times. You know how. Do you check skills
for refactoring guidance, or just refactor?
```

## Documentation Variants

When writing skill-discovery instructions, test these variants:

### NULL (Baseline)
No mention of skills in instructions at all. Measures default behavior.

### Soft Suggestion
```markdown
## Skills Library
You have access to skills at `.rein/skills/`. Consider
checking for relevant skills before working on tasks.
```

### Directive
```markdown
## Skills Library
Before working on any task, check `.rein/skills/` for
relevant skills. You should use skills when they exist.
Browse: `ls .rein/skills/`
Search: `grep -r "keyword" .rein/skills/`
```

### Emphatic (XML tags)
```xml
<important_info_about_skills>
Skills contain battle-tested approaches that prevent common mistakes.
THIS IS EXTREMELY IMPORTANT. BEFORE ANY TASK, CHECK FOR SKILLS!
If a skill existed for your task and you didn't use it, you failed.
</important_info_about_skills>
```

### Process-Oriented
```markdown
## Working with Skills
1. Before starting: Check for relevant skills
2. If skill exists: Read it completely before proceeding
3. Follow the skill - it encodes lessons from past failures
Not checking before you start is choosing to repeat those mistakes.
```

## Testing Protocol

For each variant:

1. **Run NULL baseline** first — record default behavior and rationalizations
2. **Run variant** with same scenario — does agent check? Use skills? Rationalize?
3. **Pressure test** — add time/sunk cost/authority. Does compliance hold?
4. **Meta-test** — ask agent why it skipped if it did: "You had the doc but didn't check. Why?"

## Success Criteria

A variant succeeds if:
- Agent checks for skills unprompted
- Agent reads skill completely before acting
- Agent follows skill guidance under pressure
- Agent can't rationalize away compliance

A variant fails if:
- Agent skips checking even without pressure
- Agent "adapts the concept" without reading
- Agent rationalizes away under pressure
- Agent treats skill as reference not requirement

## Applying to Rein

Use this methodology to test:
1. The `using-rein` meta-skill's skill-discovery flowchart
2. CLAUDE.md instructions that reference rein skills
3. Hook-injected session-start context
4. Command descriptions that trigger skill loading

Run tests by starting a new session with each variant and presenting the pressure scenarios. Record compliance and iterate.
