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
if ($Target -match 'docs/rein/changes/.*task\.md$') { exit 0 }

# Extract short filename for matching
$EditedFile = [System.IO.Path]::GetFileName($Target)

$ChangesDir = Join-Path $env:CLAUDE_PROJECT_DIR "docs\rein\changes"
if (-not (Test-Path $ChangesDir)) { exit 0 }

$MatchedTask = ""
$MatchedTaskfile = $null

# Code file extensions for plain filename matching
$CodeExts = @('go','ts','tsx','js','jsx','py','rs','java','rb','sql','yaml','yml','json','toml','proto','graphql','css','scss','html','sh','ps1','mod','sum','env','conf','xml','dart','swift','kt','c','cpp','h','hpp','php','tf','lock','txt','md')
$ExtPattern = ($CodeExts | ForEach-Object { [regex]::Escape($_) }) -join '|'

function Extract-Refs {
    param([string]$Line)
    $refs = @()
    # 1. Backtick-enclosed references
    $btRefs = [regex]::Matches($Line, '`([^`]+)`') | ForEach-Object { $_.Groups[1].Value }
    $refs += $btRefs
    # 2. Plain filenames with code extensions (excludes task numbers like 1.1)
    $plainRefs = [regex]::Matches($Line, "(?i)\b[A-Za-z0-9_/.-]+\.($ExtPattern)\b") | ForEach-Object { $_.Value }
    $refs += $plainRefs
    return $refs | Select-Object -Unique
}

# Scan each feature directory
foreach ($featureDir in Get-ChildItem $ChangesDir -Directory) {
    $taskfile = Join-Path $featureDir.FullName "task.md"
    $planfile = Join-Path $featureDir.FullName "plan.md"

    # Phase 1: Scan task.md for file references in task lines
    if (Test-Path $taskfile) {
        $content = Get-Content $taskfile

        foreach ($line in $content) {
            if ($line -notmatch '^\s*- \[ \]') { continue }

            if ($line -match '^\s*- \[ \] (\d+\.\d+)') {
                $taskNum = $Matches[1]
            } else { continue }

            $refs = Extract-Refs -Line $line

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
    }

    # Phase 2: If no match in task.md, scan plan.md **Files:** fields
    if (-not $MatchedTask -and (Test-Path $planfile) -and (Test-Path $taskfile)) {
        $planContent = Get-Content $planfile
        $currentTask = ""

        foreach ($line in $planContent) {
            # Track current task section: ### 1.1 ...
            if ($line -match '^### (\d+\.\d+)') {
                $currentTask = $Matches[1]
            }

            # Check **Files:** line within a task section
            if ($currentTask -and $line -match '\*\*Files\*\*:') {
                $refs = Extract-Refs -Line $line

                foreach ($ref in $refs) {
                    $refBase = [System.IO.Path]::GetFileName($ref)
                    if ($refBase -eq $EditedFile) {
                        # Find this unchecked task in task.md
                        $tfContent = Get-Content $taskfile
                        if ($tfContent | Where-Object { $_ -match "^\s*- \[ \] $currentTask" }) {
                            $MatchedTask = $currentTask
                            $MatchedTaskfile = $taskfile
                            break
                        }
                    }
                }
                if ($MatchedTask) { break }
            }
        }
    }

    if ($MatchedTask) { break }
}

if ($MatchedTask -and $MatchedTaskfile) {
    # Auto-check
    $pattern = "- [ ] $MatchedTask"
    $replacement = "- [x] $MatchedTask"
    $content = Get-Content $MatchedTaskfile
    $newContent = $content -replace [regex]::Escape($pattern), $replacement
    Set-Content $MatchedTaskfile $newContent

    $Msg = "Auto-checked task $MatchedTask (file match: $EditedFile)"
    $MsgEscaped = $Msg -replace '\\', '\\' -replace '"', '\"'
    Write-Output "{`"hookSpecificOutput`": {`"hookEventName`": `"PostToolUse`", `"additionalContext`": `"$MsgEscaped`"}}"
}
