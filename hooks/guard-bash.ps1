# guard-bash Hook (PreToolUse → Bash)
# Prevents destructive commands targeting rein-managed files

$ManifestPath = Join-Path $env:CLAUDE_PROJECT_DIR ".rein\.rein-manifest"
if (-not (Test-Path $ManifestPath)) { exit 0 }

# Read tool input from env or file
$ToolInput = $env:CLAUDE_TOOL_INPUT
if (-not $ToolInput -and $env:CLAUDE_TOOL_INPUT_FILE_PATH -and (Test-Path $env:CLAUDE_TOOL_INPUT_FILE_PATH)) {
    $ToolInput = Get-Content $env:CLAUDE_TOOL_INPUT_FILE_PATH -Raw
}
if (-not $ToolInput) { exit 0 }

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
