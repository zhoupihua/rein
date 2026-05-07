# User Auth Feature — Plan

**Goal:** Implement user authentication with JWT tokens

## Architecture Overview

The auth system follows a layered architecture:
- HTTP handlers in `internal/handlers/auth.go`
- Repository layer in `internal/repository/user.go`
- JWT service in `internal/service/jwt.go`

## Dependency Graph

```
1.1 Clarify requirements
1.2 Define user model
   └── 2.1 Design schema
       └── 3.1 Implement model
           ├── 3.2 Implement handler
           └── 3.3 Add middleware
```

## Vertical Slice Strategy

Slice by user action: Register → Login → Token validation

## Risk/Mitigation Table

| Risk | Mitigation |
|------|------------|
| JWT secret leak | Environment variable, never in code |
| Password storage | bcrypt with cost factor 12 |

## Parallelization

| Task | Classification | Notes |
|------|---------------|-------|
| 1.1, 1.2 | safe | Independent research |
| 3.1, 3.2, 3.3 | sequential | Handler depends on model, middleware depends on handler |

## Self-Audit Checklist

- [x] All tasks have acceptance criteria
- [x] No placeholder values
- [x] Dependencies are satisfied in order
- [x] Each task leaves system working

## Handoff

Ready to execute. All tasks reference real file paths and specific functions.

## Task Details

### 1.1 Clarify auth requirements

- **Acceptance:** Auth requirements documented and reviewed
- **Verification:** Team sign-off on requirements doc
- **Dependencies:** None
- **Files:** `docs/auth-requirements.md`
- **Scope:** Research only
- **Notes:** Check OWASP guidelines
- **Approach:** Interview stakeholders and review OWASP
- **Edge Cases:** Multi-tenant auth requirements
- **Rollback:** Delete requirements doc

### 1.2 Define user model

- **Acceptance:** User struct defined with all required fields
- **Verification:** Compile check passes
- **Dependencies:** 1.1
- **Files:** `internal/models/user.go`
- **Scope:** Model definition only
- **Notes:** Include soft delete support

### 2.1 Design database schema

- **Acceptance:** Migration file creates users table
- **Verification:** Migration runs without error
- **Dependencies:** 1.2
- **Files:** `migrations/001_users.sql`
- **Scope:** Schema only, no data migration
- **Notes:** Use UUID primary key

### 3.1 Implement user model

- **Acceptance:** CRUD operations work for user model
- **Verification:** Unit tests pass
- **Dependencies:** 2.1
- **Files:** `internal/models/user.go`
- **Scope:** Model + repository
- **Notes:** Follow existing patterns
