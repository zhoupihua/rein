# checkbox-guard Hook (PostToolUse → Edit|Write|MultiEdit)
# Warns when a task.md file is edited without toggling a checkbox

# Read tool input from env or file
$ToolInput = $env:CLAUDE_TOOL_INPUT
if (-not $ToolInput -and $env:CLAUDE_TOOL_INPUT_FILE_PATH -and (Test-Path $env:CLAUDE_TOOL_INPUT_FILE_PATH)) {
    $ToolInput = Get-Content $env:CLAUDE_TOOL_INPUT_FILE_PATH -Raw
}
if (-not $ToolInput) { exit 0 }

# Extract file_path from JSON
if ($ToolInput -match '"file_path"\s*:\s*"([^"]+)"') {
    $Target = ($Matches[1] -replace '\\\\', '\') -replace '\\', '/'
} else {
    exit 0
}

# Only trigger on task.md files in docs/rein/changes/
if ($Target -notmatch 'docs/rein/changes/.*task\.md$') { exit 0 }

# Check if the file exists
if (-not (Test-Path $Target)) { exit 0 }

# Check tool input for checkbox toggle evidence ([x] in the edit content)
if ($ToolInput -match '\[x\]') { exit 0 }

# Also check tool result
$ToolResult = $env:CLAUDE_TOOL_RESULT
if ($ToolResult -and $ToolResult -match '\[x\]') { exit 0 }

# Task.md was edited but no checkbox was toggled - inject warning
$Msg = "WARNING: You edited a task file but did not toggle any checkbox from [ ] to [x]. If you completed a task, you MUST update its checkbox NOW. The /do loop will re-find the same task until its checkbox is updated."
$MsgEscaped = $Msg -replace '\\', '\\' -replace '"', '\"'
Write-Output "{`"hookSpecificOutput`": {`"hookEventName`": `"PostToolUse`", `"additionalContext`": `"$MsgEscaped`"}}"
