---
name: executing-plans
description: Use when you have a written implementation plan to execute — a lightweight inline execution skill for environments without subagent support
---

# Executing Plans

## Overview

Load plan, review critically, execute all tasks, report when complete.

This is a lightweight inline execution skill for environments where subagents are not available. If subagents are available, use `rein:subagent` instead — the quality of work will be significantly higher with subagent support.

## The Process

### Step 1: Load and Review Plan

1. Read the plan file (`docs/rein/changes/<name>/plan.md`)
2. Review critically — identify any questions or concerns about the plan
3. If concerns: Raise them with your human partner before starting
4. If no concerns: Proceed to Step 2

### Step 2: Execute Tasks

For each task in `task.md`:

1. Find the first `- [ ]` line
2. Read the corresponding task detail in `plan.md`
3. Implement the task following the plan steps
4. Run verifications as specified in the plan
5. If tests pass: Edit `task.md` to change `- [ ] X.Y` to `- [x] X.Y`
6. Re-read `task.md` to confirm the checkbox update
7. Commit with descriptive message
8. Go to the next task

**The checkbox loop enforces progress:** Step 1 finds the next task by scanning for `- [ ]`. If step 5-6 was skipped, step 1 will find the same task again. The ONLY way to advance is to update the checkbox.

### Step 3: Complete Development

After all tasks complete and verified:

1. Run the full test suite
2. Invoke `rein:git-workflow` to complete the branch
3. Follow that skill to verify tests, present options, execute choice

## When to Stop and Ask for Help

**STOP executing immediately when:**
- Hit a blocker (missing dependency, test fails, instruction unclear)
- Plan has critical gaps preventing starting
- You don't understand an instruction
- Verification fails repeatedly

**Ask for clarification rather than guessing.**

## When to Revisit Earlier Steps

**Return to Review (Step 1) when:**
- Partner updates the plan based on your feedback
- Fundamental approach needs rethinking

**Don't force through blockers** — stop and ask.

## Remember

- Review plan critically first
- Follow plan steps exactly
- Don't skip verifications
- Reference skills when plan says to
- Stop when blocked, don't guess
- Never start implementation on main/master branch without explicit user consent
- Re-read task.md at the start of each iteration — don't work from cached state

## Integration

**Related skills:**
- **rein:planning** — Creates the plan this skill executes
- **rein:git-worktrees** — Ensures isolated workspace
- **rein:git-workflow** — Complete development after all tasks
- **rein:subagent** — Preferred alternative when subagents are available
