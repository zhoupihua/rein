# Claude Directory Protection Hook (PreToolUse → Write|Edit|MultiEdit)
# Block manual modifications to .claude/ directory
# Users should use the Alloy install script instead

FILE_PATH="$CLAUDE_TOOL_INPUT_FILE_PATH"

if echo "$FILE_PATH" | grep -qE '(\.claude[/\\])'; then
    echo '{"decision":"block","reason":"Manual editing of .claude/ is blocked. Use the Alloy install script to update configuration: powershell -ExecutionPolicy Bypass -File install.ps1"}'
    exit 2
fi

exit 0
