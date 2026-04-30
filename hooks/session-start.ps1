# Session start hook for rein (Windows PowerShell)
# Injects the using-rein meta-skill into every new Claude Code session

$SkillFile = Join-Path $PSScriptRoot "..\skills\using-rein\SKILL.md"

if (Test-Path $SkillFile) {
    $Content = Get-Content $SkillFile -Raw
    # Escape for JSON
    $Content = $Content -replace '\\', '\\' -replace '"', '\"' -replace "`t", '\t' -replace "`r`n", '\n' -replace "`n", '\n'

    # Scan for active tasks
    $TasksDir = Join-Path $env:CLAUDE_PROJECT_DIR "docs\rein\tasks"
    if (Test-Path $TasksDir) {
        $TaskFiles = Get-ChildItem "$TasksDir\*task.md" -File
        foreach ($TaskFile in $TaskFiles) {
            $Unchecked = (Select-String -Path $TaskFile.FullName -Pattern '^\s*- \[ \]' -SimpleMatch:$false).Count
            if ($Unchecked -gt 0) {
                $FName = $TaskFile.Name
                $ActiveMsg = "\n\nACTIVE TASKS: $Unchecked unchecked task(s) in $FName. Use /continue to resume or /status to check progress."
                $Content = $Content + $ActiveMsg
                break
            }
        }
    }

    # Claude Code format (default)
    $Output = '{"hookSpecificOutput": {"additionalContext": "' + $Content + '"}}'
    Write-Output $Output
} else {
    $Warning = '{"hookSpecificOutput": {"additionalContext": "Warning: rein using-rein skill not found. Run install script to set up."}}'
    Write-Error $Warning
}