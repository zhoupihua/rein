# User Auth Feature

## 1. Define

- [ ] 1.1 Clarify auth requirements `docs/auth-requirements.md`
- [ ] 1.2 Define user model `internal/models/user.go`
- [ ] 1.3 Specify API endpoints `api/auth.go`
- [x] 1.4 Review with team

## 2. Plan

- [ ] 2.1 Design database schema `migrations/001_users.sql`
- [ ] 2.2 Plan API structure
- [x] 2.3 Choose auth library

## 3. Build

- [ ] 3.1 Implement user model `internal/models/user.go`
  - [ ] RED: Test user model fields
  - [x] GREEN: Implement user struct
  - [ ] REFACTOR: Extract validation helpers
- [ ] 3.2 Implement auth handler `internal/handlers/auth.go`
- [ ] 3.3 Write tests `internal/handlers/auth_test.go`
- [ ] 3.4 Add middleware `internal/middleware/auth.go`
