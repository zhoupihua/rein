# Format Hook (PostToolUse → Write|Edit|MultiEdit)
# Auto-format files with Prettier after Claude modifies them
npx prettier --write "$env:CLAUDE_TOOL_INPUT_FILE_PATH" 2>$null; exit 0
