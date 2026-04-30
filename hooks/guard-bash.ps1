# guard-bash Hook (PreToolUse → Bash)
# Prevents destructive commands targeting rein-managed files

$ManifestPath = Join-Path $env:CLAUDE_PROJECT_DIR ".claude\.rein-manifest"
if (-not (Test-Path $ManifestPath)) { exit 0 }

$ToolInput = $env:CLAUDE_TOOL_INPUT

# Only check destructive commands
if ($ToolInput -notmatch '(rm |rmdir |del |mv |sed -i|Remove-Item|Move-Item|Set-Content|Out-File)') {
    exit 0
}

$entries = Get-Content $ManifestPath
foreach ($entry in $entries) {
    if ($entry -match '^\s*#' -or $entry -match '^\s*$') { continue }
    $normalizedEntry = $entry -replace '\\', '/'
    if ($ToolInput -like "*$normalizedEntry*") {
        Write-Output "{`"decision`":`"block`",`"reason`":`"Command targets rein-managed file: $normalizedEntry. Use rein commands to update.`"}"
        exit 2
    }
}

exit 0
