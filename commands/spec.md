Generate change artifacts. Replaces /opsx:propose + /opsx:explore + /opsx:continue + /opsx:ff.

## Modes

### No arguments: Interactive mode
1. Ask what the user wants to build
2. Explore the codebase for context (invoke `idea-refine` for divergent thinking)
3. Propose a change name and generate artifacts step by step
4. After each artifact, ask if they want to continue to the next

### /spec <name>: Direct generation
1. Create `docs/rein/specs/YYYY-MM-DD-<name>/` directory (if needed)
2. Generate all artifacts in dependency order:
   - `docs/rein/specs/YYYY-MM-DD-<name>-design.md` — Design spec with requirements
   - `docs/rein/plans/YYYY-MM-DD-<name>.md` — Implementation plan (decision layer)
   - `docs/rein/tasks/YYYY-MM-DD-<name>-tasks.md` — Ordered task checklist (execution layer)

### /spec --step: Step-by-step
Generate one artifact at a time, stopping after each for review.

### /spec --validate: Validate artifacts
Check that the current change has all required artifacts and they are complete.

## Artifact Templates

Read templates from the `templates/` directory:
- `templates/proposal.md` — Proposal structure
- `templates/spec.md` — Delta spec format
- `templates/design.md` — Design document structure
- `templates/tasks.md` — Task checklist format

## Artifact Dependency Order

```
proposal → specs + design (parallel) → tasks
```

Do not generate `tasks.md` until both `specs/` and `design.md` exist.

## Output

After generation, report:
```
Created artifacts:
  ✓ docs/rein/specs/YYYY-MM-DD-<name>-design.md
  ✓ docs/rein/plans/YYYY-MM-DD-<name>.md
  ✓ docs/rein/tasks/YYYY-MM-DD-<name>-tasks.md (N tasks, M phases)

Next step: /build to start implementing, or /plan to refine tasks
```