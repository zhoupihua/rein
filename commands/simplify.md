Simplify code while preserving behavior.

## Instructions

1. Invoke `simplify` skill
2. Read project conventions and understand the code
3. Scan for simplification opportunities:
   - Structural complexity (nested conditionals, excessive abstraction)
   - Naming (unclear names, inconsistent terminology)
   - Redundancy (duplicate logic, unnecessary wrappers)
4. Apply simplifications incrementally:
   - One change at a time
   - Run tests after each change
   - If tests fail, revert immediately
5. Never simplify code you don't fully understand (Chesterton's Fence)

## Safety Rules

- **Preserve behavior** — Simplification must not change what the code does
- **Test after each change** — If tests fail, revert
- **Ask before deleting** — Don't remove code that might be needed elsewhere
- **Document decisions** — Note what was simplified and why