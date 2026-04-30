# rein Protection Hook (PreToolUse → Edit|Write|MultiEdit)
# Prevents modification of rein-managed files

$ManifestPath = Join-Path $env:CLAUDE_PROJECT_DIR ".claude\.rein-manifest"
if (-not (Test-Path $ManifestPath)) { exit 0 }

$ToolInput = $env:CLAUDE_TOOL_INPUT

# Extract file_path from JSON
if ($ToolInput -match '"file_path"\s*:\s*"([^"]+)"') {
    $Target = $Matches[1] -replace '\\', '/'
} else {
    exit 0
}

# Check against manifest
$entries = Get-Content $ManifestPath
foreach ($entry in $entries) {
    if ($entry -match '^\s*#' -or $entry -match '^\s*$') { continue }
    $normalizedEntry = $entry -replace '\\', '/'
    if ($Target -like "*$normalizedEntry*") {
        Write-Output '{"decision":"block","reason":"This file is managed by rein and cannot be modified. Use rein commands to update, or remove its path from .claude/.rein-manifest to allow edits."}'
        exit 2
    }
}

exit 0
