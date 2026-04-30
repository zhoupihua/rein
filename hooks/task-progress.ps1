# task-progress: PostToolUse hook on Edit|Write|MultiEdit
# Injects task progress after code edits, making checkbox state visible to AI
# Inspired by OpenSpec's CLI feedback loop — AI sees progress and naturally corrects

$ToolInput = $env:CLAUDE_TOOL_INPUT
if (-not $ToolInput -and $env:CLAUDE_TOOL_INPUT_FILE_PATH -and (Test-Path $env:CLAUDE_TOOL_INPUT_FILE_PATH)) {
    $ToolInput = Get-Content $env:CLAUDE_TOOL_INPUT_FILE_PATH -Raw
}
if (-not $ToolInput) { exit 0 }

# Extract target file path
if ($ToolInput -match '"file_path"\s*:\s*"([^"]+)"') {
    $Target = ($Matches[1] -replace '\\\\', '\') -replace '\\', '/'
} else { exit 0 }

# Skip when editing task.md (avoid recursive noise)
if ($Target -match 'docs/rein/tasks/.*task\.md$') { exit 0 }

# Parse task.md progress
$TasksDir = Join-Path $env:CLAUDE_PROJECT_DIR "docs\rein\tasks"
if (-not (Test-Path $TasksDir)) { exit 0 }

$Total = 0; $Complete = 0; $UncheckedList = @()
foreach ($taskfile in Get-ChildItem "$TasksDir\*task.md" -File) {
    foreach ($line in Get-Content $taskfile) {
        if ($line -match '^\s*- \[[xX]\]') {
            $Total++; $Complete++
        } elseif ($line -match '^\s*- \[ \]') {
            $Total++
            $desc = ($line -replace '^\s*- \[ \] ', '')
            if ($desc.Length -gt 60) { $desc = $desc.Substring(0, 60) }
            $UncheckedList += $desc
        }
    }
}

# No tasks or all complete → silent
if ($Total -eq 0 -or ($Total - $Complete) -eq 0) { exit 0 }

$UncheckedStr = $UncheckedList -join ', '

# Inject progress
$Msg = "Task Progress: ${Complete}/${Total}. Unchecked: ${UncheckedStr}. If you completed a task, update its checkbox in task.md: ``- [ ]`` → ``- [x]``"
$MsgEscaped = $Msg -replace '\\', '\\' -replace '"', '\"'
Write-Output "{`"hookSpecificOutput`": {`"additionalContext`": `"$MsgEscaped`"}}"
