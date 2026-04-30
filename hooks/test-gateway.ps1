# Test Gateway Hook (PreToolUse → Bash)
# Run tests before deploy/push/publish commands
if ($env:CLAUDE_TOOL_INPUT -match '(deploy|push|publish)') {
    npm test
}
exit 0
