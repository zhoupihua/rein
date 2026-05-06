# Session start hook for rein (Windows PowerShell)
# Injects the using-rein meta-skill into every new Claude Code session

$SkillFile = Join-Path $PSScriptRoot "..\skills\using-rein\SKILL.md"

if (Test-Path $SkillFile) {
    $Content = Get-Content $SkillFile -Raw
    # Escape for JSON
    $Content = $Content -replace '\\', '\\' -replace '"', '\"' -replace "`t", '\t' -replace "`r`n", '\n' -replace "`n", '\n'

    # Scan for active features and their status
    $ChangesDir = Join-Path $env:CLAUDE_PROJECT_DIR "docs\rein\changes"
    if (Test-Path $ChangesDir) {
        $FeatureDirs = Get-ChildItem $ChangesDir -Directory
        foreach ($FeatureDir in $FeatureDirs) {
            $FName = $FeatureDir.Name

            # Check for active tasks
            $TaskFile = Join-Path $FeatureDir.FullName "task.md"
            if (Test-Path $TaskFile) {
                $Unchecked = (Select-String -Path $TaskFile -Pattern '^\s*- \[ \]' -SimpleMatch:$false).Count
                if ($Unchecked -gt 0) {
                    $ActiveMsg = "\n\nACTIVE TASKS: $Unchecked unchecked task(s) in $FName. Use /continue to resume or /status to check progress."
                    $Content = $Content + $ActiveMsg
                }
            }

            # Check phase completeness
            # DEFINE: refine.md, spec.md, design.md
            $DefineMissing = @()
            if (-not (Test-Path (Join-Path $FeatureDir.FullName "refine.md"))) { $DefineMissing += " refine.md" }
            if (-not (Test-Path (Join-Path $FeatureDir.FullName "spec.md"))) { $DefineMissing += " spec.md" }
            if (-not (Test-Path (Join-Path $FeatureDir.FullName "design.md"))) { $DefineMissing += " design.md" }
            $HasDefineArtifact = (Test-Path (Join-Path $FeatureDir.FullName "refine.md")) -or (Test-Path (Join-Path $FeatureDir.FullName "spec.md")) -or (Test-Path (Join-Path $FeatureDir.FullName "design.md"))
            if ($DefineMissing.Count -gt 0 -and $HasDefineArtifact) {
                $MissingStr = $DefineMissing -join ""
                $Content = $Content + "\n⚠️ DEFINE stage incomplete, missing:$MissingStr"
            }

            # PLAN: plan.md, task.md
            $PlanMissing = @()
            if (-not (Test-Path (Join-Path $FeatureDir.FullName "plan.md"))) { $PlanMissing += " plan.md" }
            if (-not (Test-Path (Join-Path $FeatureDir.FullName "task.md"))) { $PlanMissing += " task.md" }
            $HasPlanArtifact = (Test-Path (Join-Path $FeatureDir.FullName "plan.md")) -or (Test-Path (Join-Path $FeatureDir.FullName "task.md"))
            if ($PlanMissing.Count -gt 0 -and $HasPlanArtifact) {
                $MissingStr = $PlanMissing -join ""
                $Content = $Content + "\n⚠️ PLAN stage incomplete, missing:$MissingStr"
            }

            # REVIEW: review.md (only check if tasks are all done)
            if (Test-Path $TaskFile) {
                $UncheckedR = (Select-String -Path $TaskFile -Pattern '^\s*- \[ \]' -SimpleMatch:$false).Count
                if ($UncheckedR -eq 0 -and -not (Test-Path (Join-Path $FeatureDir.FullName "review.md"))) {
                    $Content = $Content + "\n⚠️ REVIEW stage incomplete, missing: review.md"
                }
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
