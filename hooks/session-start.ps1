# Session start hook for rein (Windows PowerShell)
# Injects the using-workflow meta-skill into every new Claude Code session

$SkillFile = Join-Path $PSScriptRoot "..\skills\using-workflow\SKILL.md"

if (Test-Path $SkillFile) {
    $Content = Get-Content $SkillFile -Raw
    # Escape for JSON
    $Content = $Content -replace '\\', '\\' -replace '"', '\"' -replace "`t", '\t' -replace "`r`n", '\n' -replace "`n", '\n'

    # Claude Code format (default)
    $Output = '{"hookSpecificOutput": {"additionalContext": "' + $Content + '"}}'
    Write-Output $Output
} else {
    $Warning = '{"hookSpecificOutput": {"additionalContext": "Warning: rein using-workflow skill not found. Run install script to set up."}}'
    Write-Error $Warning
}