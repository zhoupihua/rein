# rein install script (Windows)
# Project install: powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1
# Global install:  powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1 -Global

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WorkflowDir = Split-Path -Parent $ScriptDir

# Parse arguments
$Global = $false
if ($args -contains "--global" -or $args -contains "-Global") {
    $Global = $true
}

# --- Shared helper: build/download binary ---
function Install-Binary([string]$BinDir) {
    $ReinVersion = "v0.1.0"
    $Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    $BinaryName = "rein-windows-$Arch.exe"
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null

    # Check existing binary
    if (Test-Path "$BinDir\rein.exe") {
        $existingVersion = & "$BinDir\rein.exe" --version 2>$null
        Write-Host "  INFO Existing rein found: $existingVersion — upgrading" -ForegroundColor Gray
    }

    if (Test-Path "$WorkflowDir\cmd\rein\main.go") {
        # Dev mode: build from source
        Write-Host "  Building from source..." -ForegroundColor Gray
        Push-Location $WorkflowDir
        go build -o "$BinDir\rein.exe" .\cmd\rein\
        Pop-Location
    } else {
        Invoke-WebRequest -Uri "https://github.com/zhoupihua/rein/releases/download/$ReinVersion/$BinaryName" `
            -OutFile "$BinDir\rein.exe"
    }
    Write-Host "  OK rein CLI installed to $BinDir\rein.exe"
}

# --- Shared helper: generate settings.json ---
function Configure-Settings([string]$SettingsFile, [string]$HookCmd) {
    if (Test-Path $SettingsFile) {
        # Smart merge
        $data = Get-Content $SettingsFile -Raw | ConvertFrom-Json

        $reinHooks = @{}
        $reinHooks["SessionStart"] = @([PSCustomObject]@{matcher = ""; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd session-start"})})
        $reinHooks["PreToolUse"] = @(
            [PSCustomObject]@{matcher = "Edit|Write|MultiEdit"; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd guard"})},
            [PSCustomObject]@{matcher = "Bash"; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd guard-bash"}, [PSCustomObject]@{type = "command"; command = "$HookCmd gate"})}
        )
        $reinHooks["PostToolUse"] = @(
            [PSCustomObject]@{matcher = "Write|Edit|MultiEdit"; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd format"}, [PSCustomObject]@{type = "command"; command = "$HookCmd checkbox-guard"}, [PSCustomObject]@{type = "command"; command = "$HookCmd task-progress"}, [PSCustomObject]@{type = "command"; command = "$HookCmd artifact-validate"})},
            [PSCustomObject]@{matcher = "Read|Bash"; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd leak-guard"})}
        )
        $reinHooks["UserPromptExpansion"] = @([PSCustomObject]@{matcher = "code-review"; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd inject"})})

        if (-not $data.hooks) { $data | Add-Member -NotePropertyName hooks -NotePropertyValue @{} }
        foreach ($event in $reinHooks.Keys) {
            $existing = $data.hooks.$event
            if ($existing) {
                $filtered = @($existing | Where-Object {
                    -not ($_.hooks | Where-Object { $_.command -match "rein.*hook" })
                })
                $data.hooks.$event = @($filtered) + @($reinHooks.$event)
            } else {
                $data.hooks | Add-Member -NotePropertyName $event -NotePropertyValue $reinHooks.$event -Force
            }
        }

        # Add Bash(rein *) to permissions.allow
        if (-not $data.permissions) { $data | Add-Member -NotePropertyName permissions -NotePropertyValue @{} }
        if (-not $data.permissions.allow) { $data.permissions | Add-Member -NotePropertyName allow -NotePropertyValue @() }
        if ($data.permissions.allow -notcontains "Bash(rein *)") {
            $data.permissions.allow += "Bash(rein *)"
        }

        $data | ConvertTo-Json -Depth 10 | Set-Content $SettingsFile
        Write-Host "  OK settings.json merged (hooks + permissions)"
    } else {
        $Settings = @"
{
  "hooks": {
    "SessionStart": [{"matcher": "", "hooks": [{"type": "command", "command": "$HookCmd session-start"}]}],
    "PreToolUse": [
      {"matcher": "Edit|Write|MultiEdit", "hooks": [{"type": "command", "command": "$HookCmd guard"}]},
      {"matcher": "Bash", "hooks": [
        {"type": "command", "command": "$HookCmd guard-bash"},
        {"type": "command", "command": "$HookCmd gate"}
      ]}
    ],
    "PostToolUse": [
      {"matcher": "Write|Edit|MultiEdit", "hooks": [
        {"type": "command", "command": "$HookCmd format"},
        {"type": "command", "command": "$HookCmd checkbox-guard"},
        {"type": "command", "command": "$HookCmd task-progress"},
        {"type": "command", "command": "$HookCmd artifact-validate"}
      ]},
      {"matcher": "Read|Bash", "hooks": [{"type": "command", "command": "$HookCmd leak-guard"}]}
    ],
    "UserPromptExpansion": [{"matcher": "code-review", "hooks": [{"type": "command", "command": "$HookCmd inject"}]}]
  },
  "permissions": {
    "allow": ["Bash(rein *)"]
  }
}
"@
        Set-Content -Path $SettingsFile -Value $Settings
        Write-Host "  OK settings.json created with hooks + permissions"
    }
}

# --- Shared helper: generate manifest ---
function Generate-Manifest([string]$BaseDir) {
    $ManifestFile = "$BaseDir\.rein-manifest"
    $ManifestLines = @(
        "# rein Managed Files - DO NOT EDIT",
        "# To allow edits, remove the file's line from this manifest.",
        ""
    )
    $Dirs = @("bin", "commands", "agents", "checklists")
    foreach ($Dir in $Dirs) {
        $DirPath = "$BaseDir\$Dir"
        if (Test-Path $DirPath) {
            Get-ChildItem $DirPath -File | ForEach-Object {
                $RelPath = $_.FullName.Substring($BaseDir.Length + 1).Replace('\', '/')
                $ManifestLines += $RelPath
            }
        }
    }
    if (Test-Path "$BaseDir\skills") {
        Get-ChildItem "$BaseDir\skills" -Directory | ForEach-Object {
            $RelPath = $_.FullName.Substring($BaseDir.Length + 1).Replace('\', '/') + "/"
            $ManifestLines += $RelPath
        }
    }
    Set-Content -Path $ManifestFile -Value $ManifestLines
    $ManifestCount = ($ManifestLines | Where-Object { $_ -notmatch '^\s*#' -and $_ -ne '' }).Count
    Write-Host "  OK $ManifestCount entries in .rein-manifest"
}

# --- Shared helper: copy resources ---
function Copy-Resources([string]$TargetDir) {
    # Skills — clean first to remove deleted skills on upgrade
    if (Test-Path "$TargetDir\skills") { Remove-Item "$TargetDir\skills" -Recurse -Force }
    New-Item -ItemType Directory -Path "$TargetDir\skills" -Force | Out-Null
    Copy-Item -Path "$WorkflowDir\skills\*" -Destination "$TargetDir\skills\" -Recurse -Force
    $SkillCount = (Get-ChildItem "$TargetDir\skills" -Directory).Count
    Write-Host "  OK $SkillCount skills installed"

    # Commands — clean first
    if (Test-Path "$TargetDir\commands") { Remove-Item "$TargetDir\commands" -Recurse -Force }
    New-Item -ItemType Directory -Path "$TargetDir\commands" -Force | Out-Null
    Copy-Item "$WorkflowDir\commands\*.md" "$TargetDir\commands\"
    $CmdCount = (Get-ChildItem "$TargetDir\commands\*.md").Count
    Write-Host "  OK $CmdCount commands installed"

    # Agents — clean first
    if (Test-Path "$TargetDir\agents") { Remove-Item "$TargetDir\agents" -Recurse -Force }
    New-Item -ItemType Directory -Path "$TargetDir\agents" -Force | Out-Null
    Copy-Item "$WorkflowDir\agents\*.md" "$TargetDir\agents\"
    $AgentCount = (Get-ChildItem "$TargetDir\agents\*.md").Count
    Write-Host "  OK $AgentCount agents installed"

    # Checklists — clean first
    if (Test-Path "$TargetDir\checklists") { Remove-Item "$TargetDir\checklists" -Recurse -Force }
    New-Item -ItemType Directory -Path "$TargetDir\checklists" -Force | Out-Null
    if (Test-Path "$WorkflowDir\templates\checklists\review.md") {
        Copy-Item "$WorkflowDir\templates\checklists\review.md" "$TargetDir\checklists\" -Force
        Write-Host "  OK review.md checklist installed"
    }
}

# ============================================================
# Global Install
# ============================================================
if ($Global) {
    $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }
    $BinDir = "$ConfigDir\bin"

    Write-Host "=== rein Global Installer ===" -ForegroundColor Cyan
    Write-Host "Target: $ConfigDir"
    Write-Host ""

    # Check existing installation
    if (Test-Path "$ConfigDir\.rein-manifest") {
        Write-Host "INFO Existing rein installation detected — upgrading" -ForegroundColor Yellow
    }

    # [1/8] Install binary
    Write-Host "[1/8] Installing rein CLI..." -ForegroundColor Yellow
    Install-Binary $BinDir

    # [2/8] Copy resources
    Write-Host "[2/8] Installing resources..." -ForegroundColor Yellow
    Copy-Resources $ConfigDir

    # [3/8] Generate manifest
    Write-Host "[3/8] Generating protection manifest..." -ForegroundColor Yellow
    Generate-Manifest $ConfigDir

    # [4/8] Configure settings.json (hooks use $CLAUDE_CONFIG_DIR)
    Write-Host "[4/8] Configuring settings.json..." -ForegroundColor Yellow
    $HookCmd = '${CLAUDE_CONFIG_DIR:-$HOME/.claude}/bin/rein.exe hook'
    Configure-Settings "$ConfigDir\settings.json" $HookCmd

    # Clean up old bash/ps1 hooks
    if (Test-Path "$ConfigDir\hooks") {
        Remove-Item "$ConfigDir\hooks\*.sh", "$ConfigDir\hooks\*.ps1" -Force -ErrorAction SilentlyContinue
        Write-Host "  OK Cleaned old hook scripts"
    }

    # [5/8] Add to PATH
    Write-Host "[5/8] Adding to PATH..." -ForegroundColor Yellow
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$BinDir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$BinDir", "User")
        Write-Host "  OK Added $BinDir to User PATH"
        Write-Host "  INFO Restart terminal for PATH to take effect"
    } else {
        Write-Host "  OK PATH already configured"
    }

    # [6/8] Create artifact directories
    Write-Host "[6/8] Creating artifact directories..." -ForegroundColor Yellow
    $ProjectDir = Get-Location
    New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\changes" -Force | Out-Null
    New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
    Write-Host "  OK docs/rein/{changes,archive}"

    Write-Host ""
    Write-Host "=== Global Installation Complete ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Installed:"
    Write-Host "  Binary:    $BinDir\rein.exe"
    Write-Host "  Resources: $ConfigDir\skills\, commands\, agents\"
    Write-Host "  Hooks:     $ConfigDir\settings.json"
    Write-Host "  Artifacts: $ProjectDir\docs\rein\"
    Write-Host ""
    Write-Host "Run 'rein status' to check current phase"

# ============================================================
# Project Install
# ============================================================
} else {
    $ProjectDir = Get-Location
    $ClaudeDir = "$ProjectDir\.claude"
    $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }

    Write-Host "=== rein Project Installer ===" -ForegroundColor Cyan
    Write-Host "Target: $ProjectDir"
    Write-Host ""

    # Check existing installation
    if (Test-Path "$ClaudeDir\.rein-manifest") {
        Write-Host "INFO Existing rein installation detected — upgrading" -ForegroundColor Yellow
    }

    # [1/6] Install/upgrade rein CLI globally
    Write-Host "[1/6] Installing rein CLI..." -ForegroundColor Yellow
    $GlobalBinDir = "$ConfigDir\bin"
    Install-Binary $GlobalBinDir

    # [2/6] Copy resources
    Write-Host "[2/6] Installing resources..." -ForegroundColor Yellow
    Copy-Resources $ClaudeDir

    # [3/6] Generate manifest
    Write-Host "[3/6] Generating protection manifest..." -ForegroundColor Yellow
    Generate-Manifest $ClaudeDir

    # [4/6] Configure settings.json (hooks use global binary)
    Write-Host "[4/6] Configuring settings.json..." -ForegroundColor Yellow
    $HookCmd = '${CLAUDE_CONFIG_DIR:-$HOME/.claude}/bin/rein.exe hook'
    Configure-Settings "$ClaudeDir\settings.json" $HookCmd

    # Clean up old bash/ps1 hooks
    if (Test-Path "$ClaudeDir\hooks") {
        Remove-Item "$ClaudeDir\hooks\*.sh", "$ClaudeDir\hooks\*.ps1" -Force -ErrorAction SilentlyContinue
        Write-Host "  OK Cleaned old hook scripts"
    }

    # [5/6] Create artifact directories
    Write-Host "[5/6] Creating artifact directories..." -ForegroundColor Yellow
    New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\changes" -Force | Out-Null
    New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
    Write-Host "  OK docs/rein/{changes,archive}"

    # [6/6] Verification
    Write-Host "[6/6] Verifying installation..." -ForegroundColor Yellow
    if (Test-Path "$GlobalBinDir\rein.exe") {
        $ver = & "$GlobalBinDir\rein.exe" --version 2>$null
        Write-Host "  OK rein $ver available at $GlobalBinDir\" -ForegroundColor Green
    } else {
        Write-Host "  WARN rein CLI not found — installation may have failed" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "=== Project Installation Complete ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Installed:"
    Write-Host "  Resources: $ClaudeDir\skills\, commands\, agents\"
    Write-Host "  Hooks:     $ClaudeDir\settings.json"
    Write-Host "  Artifacts: $ProjectDir\docs\rein\"
    Write-Host ""
    Write-Host "Run 'rein status' to check current phase"
}
