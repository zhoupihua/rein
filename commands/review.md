Five-axis code review with security and performance checks. Replaces /opsx:verify.

## Instructions

1. Invoke `code-review-and-quality` skill
2. Review across five axes:
   - **Correctness** — Does it match spec? Edge cases? Error paths?
   - **Readability** — Clear names? Straightforward logic? No unnecessary complexity?
   - **Architecture** — Follows existing patterns? Clean boundaries?
   - **Security** — Input validated? No secrets? Parameterized queries?
   - **Performance** — N+1 queries? Unbounded operations? Missing pagination?
3. Invoke `security-and-hardening` skill for security-critical changes
4. Invoke `performance-optimization` skill for performance-sensitive changes
5. Categorize findings:
   - **Critical** — Blocks merge (security vulnerability, data loss)
   - **Important** — Should fix before merge
   - **Suggestion** — Worth considering
   - **Nit** — Minor, optional

## Output

```
## Review Summary

### Critical (must fix)
- [list]

### Important (should fix)
- [list]

### Suggestions
- [list]

### Nits
- [list]

### Verdict
APPROVE | REQUEST CHANGES
```