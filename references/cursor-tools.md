# Cursor Tool Mapping

Maps Claude Code tool names to Cursor equivalents for skill content that references Claude Code-specific tools.

## Tool Mappings

| Claude Code Tool | Cursor Equivalent | Notes |
|-----------------|-------------------|-------|
| `Read` | Automatic | Cursor reads files automatically when referenced |
| `Edit` | Automatic | Cursor applies edits via diff |
| `Write` | Automatic | Cursor creates/overwrites files |
| `Bash` | Terminal | Use Cursor's integrated terminal |
| `Glob` | File search | Cursor searches files automatically |
| `Grep` | Search | Cursor searches content automatically |
| `Agent` | `@<rule-name>` | Reference agent rules by name in chat |
| `Skill` | `@<rule-name>` | Reference skill rules by name in chat |
| `TodoWrite` | Manual | Track tasks with `rein task done <id>` in terminal |
| `TaskCreate` | `rein task done` | Use rein CLI for task management |
| `CronCreate` | Unavailable | No recurring task support in Cursor |
| `WebSearch` | Available | Cursor can search the web |
| `LSP` | Built-in | Cursor has native code intelligence |

## Key Differences

1. **No slash commands** — Cursor doesn't have `/command` syntax. Use `@<rule-name>` to reference rules instead.
2. **No agent spawning** — Cursor doesn't support subagent isolation. Reference agent personas via `@<rule-name>` and the AI will adopt that perspective.
3. **No hooks with context injection** — Cursor hooks can only block (via exit code) or allow. They cannot inject additional context into the AI's response.
4. **No permission system** — Cursor doesn't have an equivalent to Claude Code's permission prompts. File protection relies on PreEdit hooks only.
