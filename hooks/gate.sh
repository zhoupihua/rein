# gate Hook (PreToolUse → Bash)
# Run tests before deploy/push/publish commands
echo "$CLAUDE_TOOL_INPUT" | grep -qE '(deploy|push|publish)' && npm test || true
