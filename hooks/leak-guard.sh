# leak-guard Hook (PostToolUse → Read|Bash)
# Block output containing potential secrets/API keys
echo "$CLAUDE_TOOL_RESULT" | grep -qiE '(AKIA|sk-|ghp_|-----BEGIN (RSA|EC|OPENSSH))' && echo '{"decision":"block","reason":"Possible secret detected in output"}' && exit 2 || true
