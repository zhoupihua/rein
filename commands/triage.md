Auto-triage a change to determine the right workflow level.

## Instructions

Ask the user to describe the change they want to make, then classify it:

### L1: Quick (`/quick`)
- ≤5 lines of code
- No logic impact
- Examples: typo fix, constant update, config change, CSS color tweak

### L2: Fix (`/fix`)
- 1-3 files affected
- Clear requirements
- Examples: single bug fix, small feature addition, API parameter change

### L3: Feature (`/feature`)
- Multi-file change
- New feature or architecture change
- Unclear requirements needing refinement
- Examples: new user flow, database schema change, multi-component feature

## Output

```
Based on the description, this is a [L1/L2/L3] change.

Recommended command: /quick | /fix | /feature

Reasoning: [1-2 sentences]

Proceed with /quick | /fix | /feature?
```