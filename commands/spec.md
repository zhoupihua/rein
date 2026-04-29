Generate change artifacts. Replaces /opsx:propose + /opsx:explore + /opsx:continue + /opsx:ff.

## Modes

### No arguments: Interactive mode
1. Ask what the user wants to build
2. Explore the codebase for context (invoke `idea-refine` for divergent thinking)
3. Propose a change name and generate artifacts step by step
4. After each artifact, ask if they want to continue to the next

### /spec <name>: Direct generation
1. Create `changes/<name>/` directory
2. Generate all artifacts in dependency order:
   - `proposal.md` — Why and what changes
   - `specs/<feature>/spec.md` — Delta specs (ADDED/MODIFIED/REMOVED)
   - `design.md` — Engineering decisions, goals/non-goals, risks
   - `tasks.md` — Ordered task checklist with checkbox format
3. Create `.change.yaml` with metadata

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
Created change: <name>
Artifacts:
  ✓ proposal.md
  ✓ specs/<feature>/spec.md
  ✓ design.md
  ✓ tasks.md (N tasks, M phases)

Next step: /plan to refine tasks, or /build to start implementing
```