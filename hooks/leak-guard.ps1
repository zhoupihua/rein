# Secret Scan Hook (PostToolUse → Read|Bash)
# Block output containing potential secrets/API keys
if ($env:CLAUDE_TOOL_RESULT -match '(AKIA|sk-|ghp_|-----BEGIN (RSA|EC|OPENSSH))') {
    Write-Output '{"decision":"block","reason":"Possible secret detected in output"}'
    exit 2
}
exit 0
