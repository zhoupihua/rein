# task-progress: PostToolUse hook on Edit|Write|MultiEdit
# Auto-checks task checkboxes when edited files match task descriptions
# No AI cooperation needed — directly modifies task.md

$ToolInput = $env:CLAUDE_TOOL_INPUT
if (-not $ToolInput -and $env:CLAUDE_TOOL_INPUT_FILE_PATH -and (Test-Path $env:CLAUDE_TOOL_INPUT_FILE_PATH)) {
    $ToolInput = Get-Content $env:CLAUDE_TOOL_INPUT_FILE_PATH -Raw
}
if (-not $ToolInput) { exit 0 }

# Extract target file path
if ($ToolInput -match '"file_path"\s*:\s*"([^"]+)"') {
    $Target = ($Matches[1] -replace '\\\\', '\') -replace '\\', '/'
} else { exit 0 }

# Skip task.md edits (avoid recursive triggers)
if ($Target -match 'docs/rein/tasks/.*task\.md$') { exit 0 }

# Extract short filename for matching
$EditedFile = [System.IO.Path]::GetFileName($Target)

$TasksDir = Join-Path $env:CLAUDE_PROJECT_DIR "docs\rein\tasks"
if (-not (Test-Path $TasksDir)) { exit 0 }

$MatchedTask = ""
$MatchedTaskfile = $null

foreach ($taskfile in Get-ChildItem "$TasksDir\*task.md" -File) {
    $content = Get-Content $taskfile

    foreach ($line in $content) {
        if ($line -notmatch '^\s*- \[ \]') { continue }

        if ($line -match '^\s*- \[ \] (\d+\.\d+)') {
            $taskNum = $Matches[1]
        } else { continue }

        # Extract backtick file references
        $refs = [regex]::Matches($line, '`([^`]+)`') | ForEach-Object { $_.Groups[1].Value }

        foreach ($ref in $refs) {
            $refBase = [System.IO.Path]::GetFileName($ref)
            if ($refBase -eq $EditedFile) {
                $MatchedTask = $taskNum
                $MatchedTaskfile = $taskfile
                break
            }
        }

        if ($MatchedTask) { break }
    }

    if ($MatchedTask) { break }
}

if ($MatchedTask) {
    # Auto-check
    $pattern = "- [ ] $MatchedTask"
    $replacement = "- [x] $MatchedTask"
    $content = Get-Content $MatchedTaskfile
    $newContent = $content -replace [regex]::Escape($pattern), $replacement
    Set-Content $MatchedTaskfile $newContent

    $Msg = "Auto-checked task $MatchedTask (file match: $EditedFile)"
    $MsgEscaped = $Msg -replace '\\', '\\' -replace '"', '\"'
    Write-Output "{`"hookSpecificOutput`": {`"additionalContext`": `"$MsgEscaped`"}}"
}
