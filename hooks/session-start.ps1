# Session start hook for rein (Windows PowerShell)
# Injects the using-rein meta-skill into every new Claude Code session

$SkillFile = Join-Path $PSScriptRoot "..\skills\using-rein\SKILL.md"

if (Test-Path $SkillFile) {
    $Content = Get-Content $SkillFile -Raw
    # Escape for JSON
    $Content = $Content -replace '\\', '\\' -replace '"', '\"' -replace "`t", '\t' -replace "`r`n", '\n' -replace "`n", '\n'

    # Scan for active tasks
    $ChangesDir = Join-Path $env:CLAUDE_PROJECT_DIR "docs\rein\changes"
    if (Test-Path $ChangesDir) {
        $FeatureDirs = Get-ChildItem $ChangesDir -Directory
        foreach ($FeatureDir in $FeatureDirs) {
            $TaskFile = Join-Path $FeatureDir.FullName "task.md"
            if (-not (Test-Path $TaskFile)) { continue }
            $Unchecked = (Select-String -Path $TaskFile -Pattern '^\s*- \[ \]' -SimpleMatch:$false).Count
            if ($Unchecked -gt 0) {
                $FName = $FeatureDir.Name
                $ActiveMsg = "\n\nACTIVE TASKS: $Unchecked unchecked task(s) in $FName. Use /continue to resume or /status to check progress."
                $Content = $Content + $ActiveMsg
                break
            }
        }
    }
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
    $Output = '{"hookSpecificOutput": {"hookEventName": "SessionStart", "additionalContext": "' + $Content + '"}}'
    Write-Output $Output
} else {
    $Warning = '{"hookSpecificOutput": {"hookEventName": "SessionStart", "additionalContext": "Warning: rein using-rein skill not found. Run install script to set up."}}'
    Write-Error $Warning
}