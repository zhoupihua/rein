Pre-launch fan-out orchestrator with GO/NO-GO decision. Replaces /opsx:archive.

## Instructions

### Phase A: Fan-Out (Parallel Review)

Spawn 3 expert agents concurrently:

1. **code-reviewer** — Five-axis code review (correctness, readability, architecture, security, performance)
2. **security-auditor** — OWASP assessment, threat modeling, vulnerability detection
3. **test-engineer** — Coverage analysis, Prove-It verification, edge case coverage

**Skip fan-out ONLY if ALL of these are true:**
- ≤2 files changed
- <50 lines changed
- No auth, payments, data handling, or config changes

### Phase B: Merge Reports

Collect findings from all three agents. Categorize by severity:
- **Critical** — Must fix before proceeding
- **High** — Should fix, but may proceed with justification
- **Medium** — Worth addressing
- **Low/Info** — No action needed

### Phase C: GO/NO-GO Decision

**GO if:**
- Zero Critical findings
- All High findings have remediation plan
- Test coverage meets project threshold
- Build succeeds

**NO-GO if:**
- Any Critical finding exists
- High findings without remediation plan
- Tests don't pass
- Build fails

**If NO-GO:** Present findings and rollback plan. Ask user whether to fix issues or rollback.

**If GO:**
1. Invoke `git-workflow` to complete the branch
2. Invoke `shipping` for release checks (if applicable)
3. Archive `docs/rein/` artifacts to `docs/rein/archive/YYYY-MM-DD-<name>/`
4. Report completion