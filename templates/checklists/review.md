# Code Review Checklist

## Correctness
- [ ] Logic matches the spec/intent
- [ ] Edge cases handled (null, empty, boundary)
- [ ] No off-by-one errors

## Security
- [ ] No hardcoded secrets or credentials
- [ ] Input validation at system boundaries
- [ ] No SQL injection / XSS / command injection vectors

## Performance
- [ ] No N+1 queries or unnecessary loops
- [ ] Appropriate data structures used
- [ ] No redundant re-renders (frontend)

## Maintainability
- [ ] Clear naming — no abbreviations without reason
- [ ] No dead code or unused imports
- [ ] Error messages are actionable

## Testing
- [ ] Happy path covered
- [ ] Failure paths covered
- [ ] Tests are deterministic (no flaky tests)
