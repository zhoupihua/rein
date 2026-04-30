# Format Hook (PostToolUse → Write|Edit|MultiEdit)
# Auto-format files with Prettier after Claude modifies them
npx prettier --write "$CLAUDE_TOOL_INPUT_FILE_PATH" 2>/dev/null || true
