Archive completed or abandoned rein artifacts to `docs/rein/archive/`.

## Instructions

1. Scan `docs/rein/changes/` for feature directories
2. List all found features and their artifact status
3. For each feature, show:
   - refine.md: exists/missing
   - spec.md: exists/missing
   - design.md: exists/missing
   - plan.md: exists/missing
   - task.md: checkbox progress (e.g. 3/5 complete)
   - review.md: exists/missing
4. Ask user which feature(s) to archive
5. Move the entire feature directory to `docs/rein/archive/<name>/`
6. Confirm archived files and source directories are clean

## Arguments

- No argument: interactive mode — list and confirm
- `<name>`: archive the matching feature directly (e.g. `/archive feishu-login`)

## Output

```
Rein Artifacts:

  1. feishu-login
     refine: ✓  spec: ✓  design: ✓  plan: ✓  task: 5/5  review: ✓

  2. sms-login
     refine: ✓  spec: ✓  design: ✗  plan: ✓  task: 2/4  review: ✗

Archive which feature(s)? (number, name, or "all"): 1

Archived feishu-login:
  → docs/rein/archive/feishu-login/refine.md
  → docs/rein/archive/feishu-login/spec.md
  → docs/rein/archive/feishu-login/design.md
  → docs/rein/archive/feishu-login/plan.md
  → docs/rein/archive/feishu-login/task.md
  → docs/rein/archive/feishu-login/review.md
```
