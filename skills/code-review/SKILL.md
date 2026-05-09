---
name: code-review
description: Conducts five-axis code review including simplification. Use before merging any change. Use when reviewing code written by yourself, another agent, or a human. Also governs how to receive and act on code review feedback.
---

# Code Review and Quality

## Overview

Multi-dimensional code review with quality gates. Every change gets reviewed before merge — no exceptions. Review covers five axes: correctness, readability, architecture, security, and performance. Simplification is a core review output — when review reveals unnecessary complexity, the fix is to simplify.

**The approval standard:** Approve when it definitely improves overall code health, even if it isn't perfect. Don't block because it isn't exactly how you would have written it.

## The Five-Axis Review

### 1. Correctness
- Does it match the spec or task requirements?
- Are edge cases handled (null, empty, boundary values)?
- Are error paths handled (not just the happy path)?
- Does it pass all tests? Are the tests testing the right things?

### 2. Readability & Simplicity
- Are names descriptive and consistent?
- Is control flow straightforward?
- Could this be done in fewer lines?
- Are abstractions earning their complexity?
- Are there dead code artifacts?

**Simplification signals to look for:**

| Pattern | Signal | Simplification |
|---------|--------|----------------|
| Deep nesting (3+ levels) | Hard to follow control flow | Extract into guard clauses or helpers |
| Long functions (50+ lines) | Multiple responsibilities | Split into focused functions |
| Nested ternaries | Requires mental stack to parse | Replace with if/else or lookup |
| Boolean parameter flags | `doThing(true, false, true)` | Replace with options objects |
| Repeated conditionals | Same check in multiple places | Extract to a named predicate |
| Generic names | `data`, `result`, `temp`, `val` | Rename to describe content |
| Duplicated logic | Same 5+ lines in multiple places | Extract to a shared function |
| Over-engineered patterns | Factory-for-a-factory | Replace with direct approach |

**Simplification principles:**
- **Preserve behavior exactly** — don't change what the code does, only how it expresses it
- **Prefer clarity over cleverness** — explicit code beats compact code when compact requires a mental pause
- **Follow project conventions** — simplification means consistency with the codebase, not imposing external preferences
- **Scope to what changed** — default to simplifying recently modified code, avoid drive-by refactors

### 3. Architecture
- Does it follow existing patterns or introduce a new one? If new, is it justified?
- Does it maintain clean module boundaries?
- Are dependencies flowing in the right direction?
- Is the abstraction level appropriate?

### 4. Security
- Is user input validated and sanitized?
- Are secrets kept out of code, logs, and version control?
- Are SQL queries parameterized?
- Are outputs encoded to prevent XSS?
- Is external data treated as untrusted?

See `security` for detailed guidance.

### 5. Performance
- Any N+1 query patterns?
- Any unbounded loops or data fetching?
- Any unnecessary re-renders in UI components?
- Any missing pagination on list endpoints?

See `performance` for detailed guidance.

## Review Process

1. **Understand the context** — What is this change trying to accomplish?
2. **Review the tests first** — Tests reveal intent and coverage
3. **Review the implementation** — Walk through code with five axes
4. **Categorize findings** — Label with severity (Critical/Important/Suggestion/Nit/FYI)
5. **Verify the verification** — Check what tests were run, build status, manual testing

## Finding Severity

| Prefix | Meaning | Author Action |
|--------|---------|---------------|
| **Critical:** | Blocks merge | Must address — security vulnerability, data loss, broken functionality |
| **Important:** | Should fix | Significant issue but not blocking |
| **Nit:** | Minor, optional | Author may ignore — formatting, style preferences |
| **Suggestion:** | Worth considering | Not required |
| **FYI** | Informational only | No action needed |

## Change Sizing

```
~100 lines changed   → Good. Reviewable in one sitting.
~300 lines changed   → Acceptable if it's a single logical change.
~1000 lines changed  → Too large. Split it.
```

Separate refactoring from feature work. A change that refactors and adds new behavior is two changes.

## Receiving Code Review Feedback

### The Response Pattern

```
1. READ: Complete feedback without reacting
2. UNDERSTAND: Restate requirement in own words (or ask)
3. VERIFY: Check against codebase reality
4. EVALUATE: Technically sound for THIS codebase?
5. RESPOND: Technical acknowledgment or reasoned pushback
6. IMPLEMENT: One item at a time, test each
```

### Forbidden Responses

**NEVER:**
- "You're absolutely right!" / "Great point!" / "Excellent feedback!"
- "Let me implement that now" (before verification)
- ANY gratitude expression — "Thanks", "Thank you", "Appreciate it"

**If you catch yourself about to write "Thanks":** DELETE IT. State the fix instead.

**INSTEAD:**
- Restate the technical requirement
- Ask clarifying questions
- Push back with technical reasoning if wrong
- Just start working (actions > words)

### Handling Unclear Feedback

```
IF any item is unclear:
  STOP - do not implement anything yet
  ASK for clarification on ALL unclear items

WHY: Items may be related. Partial understanding = wrong implementation.
```

**Example:**
```
Reviewer: "Fix items 1-6"
You understand 1,2,3,6. Unclear on 4,5.

❌ WRONG: Implement 1,2,3,6 now, ask about 4,5 later
✅ RIGHT: "I understand items 1,2,3,6. Need clarification on 4 and 5 before proceeding."
```

### Source-Specific Handling

**From your human partner:**
- Trusted — implement after understanding
- Still ask if scope unclear
- No performative agreement
- Skip to action or technical acknowledgment

**From External Reviewers:**
```
BEFORE implementing:
  1. Check: Technically correct for THIS codebase?
  2. Check: Breaks existing functionality?
  3. Check: Reason for current implementation?
  4. Check: Works on all platforms/versions?
  5. Check: Reviewer understands full context?

IF suggestion seems wrong:
  Push back with technical reasoning

IF can't easily verify:
  Say so: "I can't verify this without [X]. Should I [investigate/ask/proceed]?"

IF conflicts with prior architectural decisions:
  Stop and discuss with your human partner first
```

### YAGNI Check

If reviewer suggests implementing a "professional" feature:
1. Grep codebase for actual usage
2. If unused: "This isn't used. Remove it (YAGNI)?"
3. If used: Then implement properly

### Implementation Order

For multi-item feedback:
1. Clarify anything unclear FIRST
2. Blocking issues (breaks, security)
3. Simple fixes (typos, imports)
4. Complex fixes (refactoring, logic)
5. Test each fix individually

### When To Push Back

Push back when:
- Suggestion breaks existing functionality
- Reviewer lacks full context
- Violates YAGNI (unused feature)
- Technically incorrect for this stack
- Conflicts with architectural decisions
- Legacy/compatibility reasons exist

**How:** Technical reasoning, not defensiveness. Reference working tests/code. Involve your human partner if architectural.

### Gracefully Correcting Your Pushback

If you pushed back and were wrong:
```
✅ "You were right - I checked [X] and it does [Y]. Implementing now."
✅ "Verified this and you're correct. My initial understanding was wrong because [reason]. Fixing."

❌ Long apology
❌ Defending why you pushed back
❌ Over-explaining
```

### Acknowledging Correct Feedback

```
✅ "Fixed. [Brief description of what changed]"
✅ "Good catch - [specific issue]. Fixed in [location]."
✅ [Just fix it and show in the code]

❌ "You're absolutely right!"
❌ "Great point!"
❌ "Thanks for catching that!"
❌ "Thanks for [anything]"
```

Actions speak. Just fix it. The code itself shows you heard the feedback.

### GitHub Thread Replies

When replying to inline review comments on GitHub, reply in the comment thread (`gh api repos/{owner}/{repo}/pulls/{pr}/comments/{id}/replies`), not as a top-level PR comment.

## Multi-Model Review Pattern

```
Model A writes the code
    │
    ▼
Model B reviews for correctness and architecture
    │
    ▼
Model A addresses the feedback
    │
    ▼
Human makes the final call
```

## Dead Code Hygiene

After any refactoring, check for orphaned code. List it explicitly and ask before deleting.

## Refactoring Follow-Up

When code review identifies structural issues (deep nesting, long functions, duplicated logic, etc.), use the `refactor` skill's Ralph loop to address them systematically: identify the smell → characterize behavior → small refactor → verify → commit → repeat. This ensures every refactoring step is test-guarded and individually reversible.

## Honesty in Review

- Don't rubber-stamp. "LGTM" without evidence of review helps no one.
- Don't soften real issues.
- Quantify problems when possible.
- Push back on approaches with clear problems.
- Accept override gracefully.

## Dependency Discipline

Before adding any dependency:
1. Does the existing stack solve this?
2. How large is the dependency?
3. Is it actively maintained?
4. Does it have known vulnerabilities?
5. What's the license?

## The Review Checklist

```markdown
## Review: [PR/Change title]

### Context
- [ ] I understand what this change does and why

### Correctness
- [ ] Change matches spec/task requirements
- [ ] Edge cases handled
- [ ] Error paths handled
- [ ] Tests cover the change adequately

### Readability
- [ ] Names are clear and consistent
- [ ] Logic is straightforward
- [ ] No unnecessary complexity

### Architecture
- [ ] Follows existing patterns
- [ ] No unnecessary coupling or dependencies
- [ ] Appropriate abstraction level
- [ ] No circular dependencies

### Security
- [ ] No secrets in code
- [ ] Input validated at boundaries
- [ ] No injection vulnerabilities
- [ ] Auth checks in place
- [ ] External data treated as untrusted

### Performance
- [ ] No N+1 patterns
- [ ] No unbounded operations
- [ ] Pagination on list endpoints
- [ ] No large objects in hot paths

### Verification
- [ ] Tests pass
- [ ] Build succeeds
- [ ] Manual verification done (if applicable)

### Verdict
- [ ] **Approve** — Ready to merge
- [ ] **Request changes** — Issues must be addressed
```

## Red Flags

- PRs merged without any review
- "LGTM" without evidence of actual review
- Security-sensitive changes without security review
- Large PRs that are "too big to review" (split them)
- Accepting "I'll fix it later" — it never happens
- No regression tests with bug fix PRs
- Review comments without severity labels
- Simplification that requires modifying tests to pass (you likely changed behavior)
- "Simplified" code that is harder to follow than the original
- Removing error handling because "it makes the code cleaner"

## Review Report

After completing the review (L3 feature workflow only), save a review report to `docs/rein/changes/<name>/review.md`.

**Report template:**

```markdown
# Review: <name>

## 范围
- 分支: <branch-name>
- 提交: <commit-range>
- 文件: <N> changed
- 行数: +<N> / -<N>

## 五轴评分

| 维度 | 评分 | 说明 |
|------|------|------|
| 正确性 | ✅/⚠️/❌ | |
| 可读性 | ✅/⚠️/❌ | |
| 架构 | ✅/⚠️/❌ | |
| 安全性 | ✅/⚠️/❌ | |
| 性能 | ✅/⚠️/❌ | |

## 发现

| 级别 | 问题 | 位置 | 状态 |
|------|------|------|------|
| Critical/Important/Nit/Suggestion | <描述> | <file:line> | 已修复/待修复/忽略(原因) |

## 结论
- [ ] **通过** — 可以合并
- [ ] **需修改** — 上述问题解决后重新审查
```

Commit the report after saving.

## Verification

After review is complete:

- [ ] All Critical issues are resolved
- [ ] All Important issues are resolved or explicitly deferred
- [ ] Tests pass
- [ ] Build succeeds
