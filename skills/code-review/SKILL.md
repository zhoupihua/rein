---
name: code-review
description: Conducts five-axis code review. Use before merging any change. Use when reviewing code written by yourself, another agent, or a human. Also governs how to receive and act on code review feedback.
---

# Code Review and Quality

## Overview

Multi-dimensional code review with quality gates. Every change gets reviewed before merge — no exceptions. Review covers five axes: correctness, readability, architecture, security, and performance.

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

**INSTEAD:**
- Restate the technical requirement
- Ask clarifying questions
- Push back with technical reasoning if wrong
- Just start working (actions > words)

### Source-Specific Handling

**From your human partner:**
- Trusted — implement after understanding
- Still ask if scope unclear
- No performative agreement

**From External Reviewers:**
- Check: Technically correct for THIS codebase?
- Check: Breaks existing functionality?
- Check: Reviewer understands full context?
- Push back with technical reasoning if wrong

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

**How:** Technical reasoning, not defensiveness. Reference working tests/code.

### Acknowledging Correct Feedback

```
✅ "Fixed. [Brief description of what changed]"
✅ "Good catch - [specific issue]. Fixed in [location]."
✅ [Just fix it and show in the code]

❌ "You're absolutely right!"
❌ "Great point!"
❌ "Thanks for [anything]"
```

Actions speak. Just fix it.

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
- [ ] Sync operations that should be async

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
