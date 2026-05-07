---
name: shipping
description: Prepares production launches and automates CI/CD pipelines. Use when preparing to deploy to production, setting up build/deployment pipelines, configuring quality gates, or planning staged rollouts.
---

# Shipping, Launch, and CI/CD

## Overview

Ship with confidence. The goal is not just to deploy — it's to deploy safely, with monitoring in place, a rollback plan ready, automated quality gates, and a clear understanding of what success looks like. Every launch should be reversible, observable, and incremental.

**Shift Left:** Catch problems as early in the pipeline as possible. A bug caught in linting costs minutes; the same bug caught in production costs hours.

**Faster is Safer:** Smaller batches and more frequent releases reduce risk, not increase it. Frequent releases build confidence in the release process itself.

## When to Use

- Deploying a feature to production for the first time
- Setting up a new project's CI pipeline
- Adding or modifying automated checks
- Configuring deployment pipelines
- Releasing a significant change to users
- Migrating data or infrastructure
- Any deployment that carries risk (all of them)

## The Quality Gate Pipeline

Every change goes through these gates before merge:

```
Pull Request Opened
    │
    ▼
┌─────────────────┐
│   LINT CHECK     │  eslint, prettier
│   ↓ pass         │
│   TYPE CHECK     │  tsc --noEmit
│   ↓ pass         │
│   UNIT TESTS     │  jest/vitest
│   ↓ pass         │
│   BUILD          │  npm run build
│   ↓ pass         │
│   INTEGRATION    │  API/DB tests
│   ↓ pass         │
│   E2E (optional) │  Playwright/Cypress
│   ↓ pass         │
│   SECURITY AUDIT │  npm audit
│   ↓ pass         │
│   BUNDLE SIZE    │  bundlesize check
└─────────────────┘
    │
    ▼
  Ready for review
```

**No gate can be skipped.** If lint fails, fix lint — don't disable the rule.

## GitHub Actions Configuration

### Basic CI Pipeline

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
          cache: 'npm'
      - run: npm ci
      - run: npm run lint
      - run: npx tsc --noEmit
      - run: npm test -- --coverage
      - run: npm run build
      - run: npm audit --audit-level=high
```

### With Database Integration Tests

```yaml
  integration:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: ci_user
          POSTGRES_PASSWORD: ${{ secrets.CI_DB_PASSWORD }}
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '22', cache: 'npm' }
      - run: npm ci
      - name: Run migrations
        run: npx prisma migrate deploy
        env:
          DATABASE_URL: postgresql://ci_user:${{ secrets.CI_DB_PASSWORD }}@localhost:5432/testdb
      - name: Integration tests
        run: npm run test:integration
        env:
          DATABASE_URL: postgresql://ci_user:${{ secrets.CI_DB_PASSWORD }}@localhost:5432/testdb
```

### E2E Tests

```yaml
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '22', cache: 'npm' }
      - run: npm ci
      - run: npx playwright install --with-deps chromium
      - run: npm run build
      - run: npx playwright test
      - uses: actions/upload-artifact@v4
        if: failure()
        with:
          name: playwright-report
          path: playwright-report/
```

## Feeding CI Failures Back to Agents

The power of CI with AI agents is the feedback loop:

```
CI fails → Copy failure output → Feed to agent → Agent fixes → Push → CI runs again
```

**Key patterns:**
- Lint failure → Agent runs `npm run lint --fix` and commits
- Type error → Agent reads the error location and fixes the type
- Test failure → Agent follows debugging skill
- Build error → Agent checks config and dependencies

## The Pre-Launch Checklist

### Code Quality

- [ ] All tests pass (unit, integration, e2e)
- [ ] Build succeeds with no warnings
- [ ] Lint and type checking pass
- [ ] Code reviewed and approved
- [ ] No TODO comments that should be resolved before launch
- [ ] No `console.log` debugging statements in production code
- [ ] Error handling covers expected failure modes

### Security

- [ ] No secrets in code or version control
- [ ] `npm audit` shows no critical or high vulnerabilities
- [ ] Input validation on all user-facing endpoints
- [ ] Authentication and authorization checks in place
- [ ] Security headers configured (CSP, HSTS, etc.)
- [ ] Rate limiting on authentication endpoints
- [ ] CORS configured to specific origins (not wildcard)

### Performance

- [ ] Core Web Vitals within "Good" thresholds
- [ ] No N+1 queries in critical paths
- [ ] Images optimized (compression, responsive sizes, lazy loading)
- [ ] Bundle size within budget
- [ ] Database queries have appropriate indexes
- [ ] Caching configured for static assets and repeated queries

### Accessibility

- [ ] Keyboard navigation works for all interactive elements
- [ ] Screen reader can convey page content and structure
- [ ] Color contrast meets WCAG 2.1 AA (4.5:1 for text)
- [ ] Focus management correct for modals and dynamic content
- [ ] Error messages are descriptive and associated with form fields

### Infrastructure

- [ ] Environment variables set in production
- [ ] Database migrations applied (or ready to apply)
- [ ] DNS and SSL configured
- [ ] CDN configured for static assets
- [ ] Logging and error reporting configured
- [ ] Health check endpoint exists and responds

### Documentation

- [ ] README updated with any new setup requirements
- [ ] API documentation current
- [ ] ADRs written for any architectural decisions
- [ ] Changelog updated

## Feature Flag Strategy

Ship behind feature flags to decouple deployment from release:

```typescript
const flags = await getFeatureFlags(userId);
if (flags.taskSharing) {
  return <TaskSharingPanel task={task} />;
}
return null;
```

**Feature flag lifecycle:**

```
1. DEPLOY with flag OFF     → Code is in production but inactive
2. ENABLE for team/beta     → Internal testing in production environment
3. GRADUAL ROLLOUT          → 5% → 25% → 50% → 100% of users
4. MONITOR at each stage    → Watch error rates, performance, user feedback
5. CLEAN UP                 → Remove flag and dead code path after full rollout
```

**Rules:**
- Every feature flag has an owner and an expiration date
- Clean up flags within 2 weeks of full rollout
- Don't nest feature flags (creates exponential combinations)
- Test both flag states (on and off) in CI

## Staged Rollout

### The Rollout Sequence

```
1. DEPLOY to staging → Full test suite + manual smoke test
2. DEPLOY to production (flag OFF) → Verify deployment + error monitoring
3. ENABLE for team → 24-hour monitoring window
4. CANARY rollout (5%) → Monitor error rates, latency, user behavior; 24-48h window
5. GRADUAL increase (25% → 50% → 100%) → Same monitoring at each step
6. FULL rollout → Monitor for 1 week → Clean up feature flag
```

### Rollout Decision Thresholds

| Metric | Advance (green) | Hold (yellow) | Roll back (red) |
|--------|-----------------|---------------|-----------------|
| Error rate | Within 10% of baseline | 10-100% above baseline | >2x baseline |
| P95 latency | Within 20% of baseline | 20-50% above baseline | >50% above baseline |
| Client JS errors | No new error types | New errors at <0.1% | New errors at >0.1% |
| Business metrics | Neutral or positive | Decline <5% | Decline >5% |

### When to Roll Back

Roll back immediately if:
- Error rate increases by more than 2x baseline
- P95 latency increases by more than 50%
- User-reported issues spike
- Data integrity issues detected
- Security vulnerability discovered

## Monitoring and Observability

### What to Monitor

```
Application: Error rate, Response time (p50/p95/p99), Request volume, Active users, Business metrics
Infrastructure: CPU/memory, DB connection pool, Disk space, Network latency, Queue depth
Client: Core Web Vitals (LCP, INP, CLS), JS errors, API error rates, Page load time
```

### Post-Launch Verification

In the first hour after launch:

```
1. Check health endpoint returns 200
2. Check error monitoring dashboard (no new error types)
3. Check latency dashboard (no regression)
4. Test the critical user flow manually
5. Verify logs are flowing and readable
6. Confirm rollback mechanism works (dry run if possible)
```

## Rollback Strategy

Every deployment needs a rollback plan before it happens:

```markdown
## Rollback Plan for [Feature/Release]

### Trigger Conditions
- Error rate > 2x baseline
- P95 latency > [X]ms

### Rollback Steps
1. Disable feature flag OR deploy previous version
2. Verify rollback: health check, error monitoring
3. Communicate: notify team

### Database Considerations
- Migration [X] has a rollback
- Data inserted by new feature: [preserved / cleaned up]

### Time to Rollback
- Feature flag: < 1 minute
- Redeploy: < 5 minutes
- Database rollback: < 15 minutes
```

## Environment Management

```
.env.example       → Committed (template for developers)
.env                → NOT committed (local development)
.env.test           → Committed (test environment, no real secrets)
CI secrets          → Stored in GitHub Secrets / vault
Production secrets  → Stored in deployment platform / vault
```

CI should never have production secrets.

## Deployment Strategies

### Preview Deployments

Every PR gets a preview deployment for manual testing.

### Automation Beyond CI

- **Dependabot / Renovate** for automated dependency updates
- **Build Cop** role — designated person keeps CI green
- **PR Checks:** Required reviews, required status checks, branch protection, auto-merge

## CI Optimization

When the pipeline exceeds 10 minutes, apply in order of impact:

```
Slow CI pipeline?
├── Cache dependencies
├── Run jobs in parallel (split lint, typecheck, test, build)
├── Only run what changed (path filters)
├── Use matrix builds (shard test suites)
├── Optimize the test suite (move slow tests to scheduled runs)
└── Use larger runners
```

## See Also

- For security pre-launch checks, see `references/security-checklist.md`
- For performance pre-launch checklist, see `references/performance-checklist.md`
- For accessibility verification before launch, see `references/accessibility-checklist.md`

## Common Rationalizations

| Rationalization | Reality |
|---|---|
| "CI is too slow" | Optimize the pipeline, don't skip it. A 5-minute pipeline prevents hours of debugging. |
| "This change is trivial, skip CI" | Trivial changes break builds. CI is fast for trivial changes. |
| "It works in staging, it'll work in production" | Production has different data, traffic patterns, and edge cases. |
| "We don't need feature flags for this" | Every feature benefits from a kill switch. |
| "Monitoring is overhead" | Not having monitoring means you discover problems from user complaints. |
| "Rolling back is admitting failure" | Rolling back is responsible engineering. Shipping broken is the failure. |
| "The test is flaky, just re-run" | Flaky tests mask real bugs. Fix the flakiness. |
| "We'll add CI later" | Projects without CI accumulate broken states. Set it up on day one. |

## Red Flags

- No CI pipeline in the project
- CI failures ignored or silenced
- Tests disabled in CI to make the pipeline pass
- Deploying without a rollback plan
- No monitoring or error reporting in production
- Big-bang releases (everything at once, no staging)
- Feature flags with no expiration or owner
- No one monitoring the deploy for the first hour
- Production secrets stored in code or CI config
- Long CI times with no optimization effort
- "It's Friday afternoon, let's ship it"

## Verification

Before deploying:

- [ ] Pre-launch checklist completed (all sections green)
- [ ] Feature flag configured (if applicable)
- [ ] Rollback plan documented
- [ ] Monitoring dashboards set up
- [ ] Team notified of deployment
- [ ] All quality gates present (lint, types, tests, build, audit)
- [ ] Pipeline runs on every PR and push to main
- [ ] Failures block merge (branch protection configured)

After deploying:

- [ ] Health check returns 200
- [ ] Error rate is normal
- [ ] Latency is normal
- [ ] Critical user flow works
- [ ] Logs are flowing
- [ ] Rollback tested or verified ready
