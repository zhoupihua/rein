# guard Hook (PreToolUse → Edit|Write|MultiEdit)
# Prevents modification of rein-managed files

$ManifestPath = Join-Path $env:CLAUDE_PROJECT_DIR ".rein\.rein-manifest"
if (-not (Test-Path $ManifestPath)) { exit 0 }

# Read tool input from env or file
$ToolInput = $env:CLAUDE_TOOL_INPUT
if (-not $ToolInput -and $env:CLAUDE_TOOL_INPUT_FILE_PATH -and (Test-Path $env:CLAUDE_TOOL_INPUT_FILE_PATH)) {
    $ToolInput = Get-Content $env:CLAUDE_TOOL_INPUT_FILE_PATH -Raw
}
if (-not $ToolInput) { exit 0 }

# Extract file_path from JSON
if ($ToolInput -match '"file_path"\s*:\s*"([^"]+)"') {
    # Unescape JSON \\ to \, then normalize all \ to /
    $Target = ($Matches[1] -replace '\\\\', '\') -replace '\\', '/'
} else {
    exit 0
}

# Check against manifest
$entries = Get-Content $ManifestPath
foreach ($entry in $entries) {
    if ($entry -match '^\s*#' -or $entry -match '^\s*$') { continue }
    $normalizedEntry = $entry -replace '\\', '/'
    if ($Target -like "*$normalizedEntry*") {
        Write-Output '{"decision":"block","reason":"This file is managed by rein and cannot be modified. Use rein commands to update, or remove its path from .rein/.rein-manifest to allow edits."}'
        exit 2
    }
}

exit 0
