# API and Interface Design Reference

Design stable, well-documented interfaces that are hard to misuse. Good interfaces make the right thing easy and the wrong thing hard.

## Core Principles

### Hyrum's Law

> With a sufficient number of users of an API, all observable behaviors of your system will be depended on by somebody, regardless of what you promise in the contract.

Every public behavior — including undocumented quirks, error message text, timing, and ordering — becomes a de facto contract once users depend on it.

- **Be intentional about what you expose.** Every observable behavior is a potential commitment.
- **Don't leak implementation details.** If users can observe it, they will depend on it.
- **Plan for deprecation at design time.**
- **Tests are not enough.** Even with perfect contract tests, "safe" changes can break real users.

### The One-Version Rule

Avoid forcing consumers to choose between multiple versions of the same dependency or API. Design for a world where only one version exists at a time — extend rather than fork.

### Contract First

Define the interface before implementing it. The contract is the spec — implementation follows.

### Consistent Error Semantics

Pick one error strategy and use it everywhere:

```typescript
interface APIError {
  error: {
    code: string;        // Machine-readable: "VALIDATION_ERROR"
    message: string;     // Human-readable: "Email is required"
    details?: unknown;   // Additional context when helpful
  };
}
```

Don't mix patterns. If some endpoints throw, others return null, and others return `{ error }` — the consumer can't predict behavior.

### Validate at Boundaries

Trust internal code. Validate at system edges where external input enters. Third-party API responses are untrusted data — always validate their shape before using them.

Where validation belongs: API route handlers, form submission handlers, external service response parsing, environment variable loading.

Where validation does NOT belong: Between internal functions that share type contracts, in utility functions called by already-validated code, on data from your own database.

### Prefer Addition Over Modification

Extend interfaces without breaking existing consumers — add optional fields, never change or remove existing ones.

## REST API Patterns

### Resource Design

```
GET    /api/tasks              → List tasks (with query params for filtering)
POST   /api/tasks              → Create a task
GET    /api/tasks/:id          → Get a single task
PATCH  /api/tasks/:id          → Update a task (partial)
DELETE /api/tasks/:id          → Delete a task
GET    /api/tasks/:id/comments → List comments for a task (sub-resource)
POST   /api/tasks/:id/comments → Add a comment to a task
```

### Pagination

```typescript
GET /api/tasks?page=1&pageSize=20&sortBy=createdAt&sortOrder=desc
// Response:
{ "data": [...], "pagination": { "page": 1, "pageSize": 20, "totalItems": 142, "totalPages": 8 } }
```

### Predictable Naming

| Pattern | Convention | Example |
|---------|-----------|---------|
| REST endpoints | Plural nouns, no verbs | `GET /api/tasks`, `POST /api/tasks` |
| Query params | camelCase | `?sortBy=createdAt&pageSize=20` |
| Response fields | camelCase | `{ createdAt, updatedAt, taskId }` |
| Boolean fields | is/has/can prefix | `isComplete`, `hasAttachments` |
| Enum values | UPPER_SNAKE | `"IN_PROGRESS"`, `"COMPLETED"` |

## TypeScript Interface Patterns

### Discriminated Unions for Variants

```typescript
type TaskStatus =
  | { type: 'pending' }
  | { type: 'in_progress'; assignee: string; startedAt: Date }
  | { type: 'completed'; completedAt: Date; completedBy: string }
  | { type: 'cancelled'; reason: string; cancelledAt: Date };
```

### Input/Output Separation

```typescript
interface CreateTaskInput { title: string; description?: string; }
interface Task { id: string; title: string; description: string | null; createdAt: Date; updatedAt: Date; createdBy: string; }
```

### Branded Types for IDs

```typescript
type TaskId = string & { readonly __brand: 'TaskId' };
type UserId = string & { readonly __brand: 'UserId' };
```

## Red Flags

- Endpoints that return different shapes depending on conditions
- Inconsistent error formats across endpoints
- Validation scattered throughout internal code instead of at boundaries
- Breaking changes to existing fields (type changes, removals)
- List endpoints without pagination
- Verbs in REST URLs (`/api/createTask`, `/api/getUsers`)
- Third-party API responses used without validation or sanitization
