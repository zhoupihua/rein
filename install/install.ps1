# rein install script (Windows)
# Project install: powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1
# Global install:  powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1 -Global
# Cursor install:  powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1 -Ide cursor
# Codex install:   powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1 -Ide codex

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WorkflowDir = Split-Path -Parent $ScriptDir

# Parse arguments
$Global = $false
$Ide = "claude"
for ($i = 0; $i -lt $args.Count; $i++) {
    switch ($args[$i]) {
        "--global" { $Global = $true }
        "-Global" { $Global = $true }
        "--ide" { $Ide = $args[++$i] }
        "-Ide" { $Ide = $args[++$i] }
        { $_ -match "^--ide=" } { $Ide = $_.Substring(6) }
    }
}

# Validate IDE
if ($Ide -notin @("claude", "cursor", "codex")) {
    Write-Host "Error: unsupported IDE '$Ide'. Use 'claude', 'cursor', or 'codex'." -ForegroundColor Red
    exit 1
}

# --- Shared helper: build/download binary ---
function Install-Binary([string]$BinDir) {
    $ReinVersion = "v0.1.0"
    $Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    $BinaryName = "rein-windows-$Arch.exe"
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null

    # Check existing binary
    if (Test-Path "$BinDir\rein.exe") {
        try {
            $existingVersion = & "$BinDir\rein.exe" --version 2>$null
            Write-Host "  INFO Existing rein found: $existingVersion — upgrading" -ForegroundColor Gray
        } catch {
            Write-Host "  INFO Existing rein found — upgrading (version check blocked by policy)" -ForegroundColor Gray
        }
    }

    if (Test-Path "$WorkflowDir\cmd\rein\main.go") {
        # Dev mode: build from source
        Write-Host "  Building from source..." -ForegroundColor Gray
        Push-Location $WorkflowDir
        go build -mod=mod -o "$BinDir\rein.exe" .\cmd\rein\
        Pop-Location
    } else {
        Invoke-WebRequest -Uri "https://github.com/zhoupihua/rein/releases/download/$ReinVersion/$BinaryName" `
            -OutFile "$BinDir\rein.exe"
    }
    Write-Host "  OK rein CLI installed to $BinDir\rein.exe"
}

# --- Shared helper: generate settings.json (Claude Code) ---
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
            [PSCustomObject]@{matcher = "Write|Edit|MultiEdit"; hooks = @([PSCustomObject]@{type = "command"; command = "$HookCmd format"}, [PSCustomObject]@{type = "command"; command = "$HookCmd checkbox-guard"}, [PSCustomObject]@{type = "command"; command = "$HookCmd artifact-validate"})},
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

# --- Shared helper: configure hooks.json (Cursor) ---
function Configure-CursorHooks([string]$HooksFile, [string]$ReinCmd) {
    $HooksContent = @"
{
  "hooks": {
    "PreEdit": [
      {"command": "$ReinCmd", "args": ["hook", "guard"]}
    ],
    "PostEdit": [
      {"command": "$ReinCmd", "args": ["hook", "format"]}
    ],
    "PreCommit": [
      {"command": "$ReinCmd", "args": ["hook", "gate"]}
    ]
  }
}
"@
    if (Test-Path $HooksFile) {
        $data = Get-Content $HooksFile -Raw | ConvertFrom-Json
        $newHooks = $HooksContent | ConvertFrom-Json

        if (-not $data.hooks) { $data | Add-Member -NotePropertyName hooks -NotePropertyValue @{} }
        foreach ($event in $newHooks.hooks.PSObject.Properties.Name) {
            $existing = $data.hooks.$event
            $newEntries = @($newHooks.hooks.$event)
            if ($existing) {
                $filtered = @($existing | Where-Object {
                    $_.command -notmatch "rein"
                })
                $data.hooks.$event = @($filtered) + @($newEntries)
            } else {
                $data.hooks | Add-Member -NotePropertyName $event -NotePropertyValue $newEntries -Force
            }
        }

        $data | ConvertTo-Json -Depth 10 | Set-Content $HooksFile
        Write-Host "  OK hooks.json merged"
    } else {
        Set-Content -Path $HooksFile -Value $HooksContent
        Write-Host "  OK hooks.json created"
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
    # For Cursor installs, also list .mdc rule files
    if (Test-Path "$BaseDir\rules") {
        Get-ChildItem "$BaseDir\rules" -File -Filter "*.mdc" | ForEach-Object {
            $RelPath = $_.FullName.Substring($BaseDir.Length + 1).Replace('\', '/')
            $ManifestLines += $RelPath
        }
    }
    Set-Content -Path $ManifestFile -Value $ManifestLines
    $ManifestCount = ($ManifestLines | Where-Object { $_ -notmatch '^\s*#' -and $_ -ne '' }).Count
    Write-Host "  OK $ManifestCount entries in .rein-manifest"
}

# --- Shared helper: copy resources (Claude Code) ---
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

# --- Shared helper: copy resources (Cursor) ---
function Copy-Resources-Cursor([string]$CursorDir) {
    if (Test-Path "$CursorDir\rules") { Remove-Item "$CursorDir\rules" -Recurse -Force }
    New-Item -ItemType Directory -Path "$CursorDir\rules" -Force | Out-Null

    $ReinBin = if ($env:CLAUDE_CONFIG_DIR) { "$env:CLAUDE_CONFIG_DIR\bin\rein.exe" } else { "$env:USERPROFILE\.claude\bin\rein.exe" }

    if (Test-Path $ReinBin) {
        # Use rein convert for proper frontmatter transformation
        & $ReinBin convert --ide cursor --source-dir $WorkflowDir --output-dir "$CursorDir\rules"
    } else {
        # Fallback: simple copy with frontmatter transformation
        Convert-SkillsFallback "$WorkflowDir\skills" "$CursorDir\rules"
        Convert-CommandsFallback "$WorkflowDir\commands" "$CursorDir\rules"
        Convert-AgentsFallback "$WorkflowDir\agents" "$CursorDir\rules"
    }

    # Create always-apply rule
    Inject-CursorRules "$CursorDir\rules"
}

function Convert-SkillsFallback([string]$SrcDir, [string]$OutDir) {
    $count = 0
    Get-ChildItem $SrcDir -Directory | ForEach-Object {
        $skillFile = "$($_.FullName)\SKILL.md"
        if (-not (Test-Path $skillFile)) { return }
        $name = $_.Name
        $content = Get-Content $skillFile -Raw
        $desc = ($content -split "`n" | Where-Object { $_ -match '^description:' } | Select-Object -First 1) -replace '^description:\s*', ''
        if (-not $desc) { $desc = "$name skill" }

        $body = $content -replace '(?s)^---\s*\n.*?\n---\s*\n?', ''
        $mdc = "---`ndescription: $desc`nalwaysApply: false`n---$body"
        Set-Content -Path "$OutDir\$name.mdc" -Value $mdc
        $count++
    }
    Write-Host "  OK $count skills converted to .mdc rules"
}

function Convert-CommandsFallback([string]$SrcDir, [string]$OutDir) {
    $count = 0
    Get-ChildItem "$SrcDir\*.md" | ForEach-Object {
        $name = $_.BaseName
        $content = Get-Content $_.FullName -Raw
        $desc = ($content -split "`n" | Where-Object { $_ -match '^description:' } | Select-Object -First 1) -replace '^description:\s*', ''
        if (-not $desc) { $desc = "$name command" }

        $body = $content -replace '(?s)^---\s*\n.*?\n---\s*\n?', ''
        $mdc = "---`ndescription: $desc`nalwaysApply: false`n---$body"
        Set-Content -Path "$OutDir\$name.mdc" -Value $mdc
        $count++
    }
    Write-Host "  OK $count commands converted to .mdc rules"
}

function Convert-AgentsFallback([string]$SrcDir, [string]$OutDir) {
    $count = 0
    Get-ChildItem "$SrcDir\*.md" | ForEach-Object {
        $name = $_.BaseName
        $content = Get-Content $_.FullName -Raw
        $desc = ($content -split "`n" | Where-Object { $_ -match '^description:' } | Select-Object -First 1) -replace '^description:\s*', ''
        if (-not $desc) { $desc = "$name agent" }

        $body = $content -replace '(?s)^---\s*\n.*?\n---\s*\n?', ''
        $mdc = "---`ndescription: $desc`nalwaysApply: false`n---$body"
        Set-Content -Path "$OutDir\$name.mdc" -Value $mdc
        $count++
    }
    Write-Host "  OK $count agents converted to .mdc rules"
}

# --- Shared helper: inject task progress rule into CLAUDE.md ---
function Inject-ClaudeMd([string]$ProjectDir) {
    $ClaudeMd = "$ProjectDir\CLAUDE.md"
    $Marker = "<!-- rein:task-progress -->"
    $Block = @"
$Marker
## Task Progress

When working on a feature with ``docs/rein/changes/<name>/task.md``, after completing
any task or sub-task, you MUST immediately mark it as done:

  rein task done <id>          # e.g., rein task done 1.2
  rein task done <subtask-id>  # e.g., rein task done 1.2.0

Do NOT skip this step. Marking progress is mandatory, not optional.
$Marker
"@

    if (Test-Path $ClaudeMd) {
        $content = Get-Content $ClaudeMd -Raw
        if ($content -match [regex]::Escape($Marker)) {
            Write-Host "  OK CLAUDE.md task-progress rule already present"
        } else {
            Add-Content -Path $ClaudeMd -Value "`n$Block"
            Write-Host "  OK CLAUDE.md injected task-progress rule"
        }
    } else {
        Set-Content -Path $ClaudeMd -Value $Block
        Write-Host "  OK CLAUDE.md created with task-progress rule"
    }
}

# --- Shared helper: inject Cursor always-apply rule ---
function Inject-CursorRules([string]$RulesDir) {
    $RuleFile = "$RulesDir\rein-project.mdc"
    $Content = @"
---
description: rein project conventions and task progress tracking
alwaysApply: true
---

## Task Progress

When working on a feature with ``docs/rein/changes/<name>/task.md``, after completing
any task or sub-task, you MUST immediately mark it as done:

  rein task done <id>          # e.g., rein task done 1.2
  rein task done <subtask-id>  # e.g., rein task done 1.2.0

Do NOT skip this step. Marking progress is mandatory, not optional.

## rein Workflow

This project uses rein for structured development. Key commands:
- ``rein status`` — Check current workflow phase
- ``rein task next`` — Show next unchecked task
- ``rein task done <id>`` — Mark task complete
- ``rein validate <feature>`` — Validate artifact completeness

## Available Rules

Reference these rules with @<name> in Cursor chat:

| Rule | Type | When to use |
|------|------|-------------|
| @define | skill | Starting a new project or feature |
| @tdd | skill | Implementing logic or fixing bugs |
| @code-review | skill | Reviewing code before merge |
| @planning | skill | Breaking work into tasks |
| @security | skill | Security reviews |
| @performance | skill | Performance work |
| @feature | command | Multi-file feature workflow |
| @fix | command | Bug fix workflow |
| @ship | command | Pre-launch checklist |
| @code-reviewer | agent | Structured code review |
| @test-engineer | agent | Test strategy and coverage |
| @security-auditor | agent | Security analysis |
"@
    Set-Content -Path $RuleFile -Value $Content
    Write-Host "  OK rein-project.mdc (alwaysApply) created"
}

# --- Shared helper: copy resources (Codex) ---
function Copy-Resources-Codex([string]$ProjectDir) {
    $CodexDir = "$ProjectDir\.codex"
    New-Item -ItemType Directory -Path $CodexDir -Force | Out-Null

    $ReinBin = if ($env:CLAUDE_CONFIG_DIR) { "$env:CLAUDE_CONFIG_DIR\bin\rein.exe" } else { "$env:USERPROFILE\.claude\bin\rein.exe" }

    if (Test-Path $ReinBin) {
        & $ReinBin convert --ide codex --source-dir $WorkflowDir --output-dir $ProjectDir
    } else {
        Generate-CodexMDFallback $ProjectDir
        Generate-CodexConfigFallback $CodexDir
    }
}

function Generate-CodexMDFallback([string]$ProjectDir) {
    $CodexMD = "$ProjectDir\CODEX.md"
    $SB = [System.Text.StringBuilder]::new()
    [void]$SB.AppendLine("# rein Project Instructions")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("## Task Progress")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("When working on a feature with ``docs/rein/changes/<name>/task.md``, after completing")
    [void]$SB.AppendLine("any task or sub-task, you MUST immediately mark it as done:")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("  rein task done <id>          # e.g., rein task done 1.2")
    [void]$SB.AppendLine("  rein task done <subtask-id>  # e.g., rein task done 1.2.0")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("Do NOT skip this step. Marking progress is mandatory, not optional.")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("## rein Workflow")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("This project uses rein for structured development. Key commands:")
    [void]$SB.AppendLine("- ``rein status`` - Check current workflow phase")
    [void]$SB.AppendLine("- ``rein task next`` - Show next unchecked task")
    [void]$SB.AppendLine("- ``rein task done <id>`` - Mark task complete")
    [void]$SB.AppendLine("- ``rein validate <feature>`` - Validate artifact completeness")

    # Skills
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("## Skills")
    Get-ChildItem "$WorkflowDir\skills" -Directory | ForEach-Object {
        $SkillFile = "$($_.FullName)\SKILL.md"
        if (Test-Path $SkillFile) {
            [void]$SB.AppendLine()
            [void]$SB.AppendLine("### $($_.Name)")
            [void]$SB.AppendLine()
            $content = Get-Content $SkillFile -Raw
            $body = $content -replace '(?s)^---\s*\n.*?\n---\s*\n?', ''
            [void]$SB.Append($body)
        }
    }

    # Commands
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("## Commands")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("Reference these commands by name when asking Codex to perform a workflow.")
    Get-ChildItem "$WorkflowDir\commands\*.md" | ForEach-Object {
        [void]$SB.AppendLine()
        [void]$SB.AppendLine("### $($_.BaseName)")
        [void]$SB.AppendLine()
        $content = Get-Content $_.FullName -Raw
        $body = $content -replace '(?s)^---\s*\n.*?\n---\s*\n?', ''
        [void]$SB.Append($body)
    }

    # Agents
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("## Agents")
    [void]$SB.AppendLine()
    [void]$SB.AppendLine("Reference these agents by name when asking Codex to adopt a perspective.")
    Get-ChildItem "$WorkflowDir\agents\*.md" | ForEach-Object {
        [void]$SB.AppendLine()
        [void]$SB.AppendLine("### $($_.BaseName)")
        [void]$SB.AppendLine()
        $content = Get-Content $_.FullName -Raw
        $body = $content -replace '(?s)^---\s*\n.*?\n---\s*\n?', ''
        [void]$SB.Append($body)
    }

    Set-Content -Path $CodexMD -Value $SB.ToString()
    Write-Host "  OK CODEX.md created"
}

function Generate-CodexConfigFallback([string]$CodexDir) {
    $ConfigFile = "$CodexDir\config.toml"
    $ReinCmd = if ($env:CLAUDE_CONFIG_DIR) { "$env:CLAUDE_CONFIG_DIR\bin\rein.exe" } else { "$env:USERPROFILE\.claude\bin\rein.exe" }

    $Content = @"
# rein Codex configuration
[features]
multi_agent = true

[[hooks.pre_command]]
command = "$ReinCmd hook guard"
description = "Block edits to rein-managed files"

[[hooks.pre_command]]
command = "$ReinCmd hook gate"
description = "Run tests before deploy commands"

[[hooks.post_command]]
command = "$ReinCmd hook format"
description = "Auto-format web files with prettier"
"@
    Set-Content -Path $ConfigFile -Value $Content
    Write-Host "  OK .codex/config.toml created"
}

# --- Shared helper: configure Codex config.toml ---
function Configure-CodexConfig([string]$ConfigFile, [string]$ReinCmd) {
    $CodexDir = Split-Path -Parent $ConfigFile
    New-Item -ItemType Directory -Path $CodexDir -Force | Out-Null

    $Content = @"
# rein Codex configuration
[features]
multi_agent = true

[[hooks.pre_command]]
command = "$ReinCmd hook guard"
description = "Block edits to rein-managed files"

[[hooks.pre_command]]
command = "$ReinCmd hook gate"
description = "Run tests before deploy commands"

[[hooks.post_command]]
command = "$ReinCmd hook format"
description = "Auto-format web files with prettier"
"@
    if (Test-Path $ConfigFile) {
        $existing = Get-Content $ConfigFile -Raw
        if ($existing -match "rein hook guard") {
            Write-Host "  OK config.toml already has rein hooks"
        } else {
            Add-Content -Path $ConfigFile -Value "`n$Content"
            Write-Host "  OK config.toml appended with rein hooks"
        }
    } else {
        Set-Content -Path $ConfigFile -Value $Content
        Write-Host "  OK config.toml created"
    }
}

# ============================================================
# Global Install
# ============================================================
if ($Global) {
    if ($Ide -eq "cursor") {
        $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }
        $BinDir = "$ConfigDir\bin"
        $CursorGlobalDir = "$env:USERPROFILE\.cursor"

        Write-Host "=== rein Global Installer (Cursor) ===" -ForegroundColor Cyan
        Write-Host "Target: $CursorGlobalDir"
        Write-Host ""

        # [1/6] Install binary
        Write-Host "[1/6] Installing rein CLI..." -ForegroundColor Yellow
        Install-Binary $BinDir

        # [2/6] Copy resources as .mdc rules
        Write-Host "[2/6] Installing Cursor rules..." -ForegroundColor Yellow
        Copy-Resources-Cursor $CursorGlobalDir

        # [3/6] Generate manifest
        Write-Host "[3/6] Generating protection manifest..." -ForegroundColor Yellow
        Generate-Manifest $CursorGlobalDir

        # [4/6] Configure hooks.json
        Write-Host "[4/6] Configuring hooks.json..." -ForegroundColor Yellow
        $HookCmd = "$BinDir\rein.exe"
        Configure-CursorHooks "$CursorGlobalDir\hooks.json" $HookCmd

        # [5/6] Add to PATH
        Write-Host "[5/6] Adding to PATH..." -ForegroundColor Yellow
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($currentPath -notlike "*$BinDir*") {
            [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$BinDir", "User")
            Write-Host "  OK Added $BinDir to User PATH"
            Write-Host "  INFO Restart terminal for PATH to take effect"
        } else {
            Write-Host "  OK PATH already configured"
        }

        # [6/6] Create artifact directories
        Write-Host "[6/6] Creating artifact directories..." -ForegroundColor Yellow
        $ProjectDir = Get-Location
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\changes" -Force | Out-Null
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
        Write-Host "  OK docs/rein/{changes,archive}"

        Write-Host ""
        Write-Host "=== Global Installation Complete (Cursor) ===" -ForegroundColor Green
        Write-Host ""
        Write-Host "Installed:"
        Write-Host "  Binary:    $BinDir\rein.exe"
        Write-Host "  Rules:     $CursorGlobalDir\rules\"
        Write-Host "  Hooks:     $CursorGlobalDir\hooks.json"
        Write-Host "  Artifacts: $ProjectDir\docs\rein\"
        Write-Host ""
        Write-Host "Reference rules in Cursor chat with @<rule-name>"

    } elseif ($Ide -eq "codex") {
        $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }
        $BinDir = "$ConfigDir\bin"
        $ProjectDir = Get-Location
        $CodexDir = "$ProjectDir\.codex"

        Write-Host "=== rein Global Installer (Codex) ===" -ForegroundColor Cyan
        Write-Host "Target: $ProjectDir"
        Write-Host ""

        # [1/6] Install binary
        Write-Host "[1/6] Installing rein CLI..." -ForegroundColor Yellow
        Install-Binary $BinDir

        # [2/6] Copy resources as CODEX.md
        Write-Host "[2/6] Generating CODEX.md and config..." -ForegroundColor Yellow
        Copy-Resources-Codex $ProjectDir

        # [3/6] Generate manifest
        Write-Host "[3/6] Generating protection manifest..." -ForegroundColor Yellow
        Generate-Manifest $CodexDir

        # [4/6] Configure config.toml
        Write-Host "[4/6] Configuring .codex/config.toml..." -ForegroundColor Yellow
        $HookCmd = "$BinDir\rein.exe"
        Configure-CodexConfig "$CodexDir\config.toml" $HookCmd

        # [5/6] Add to PATH
        Write-Host "[5/6] Adding to PATH..." -ForegroundColor Yellow
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($currentPath -notlike "*$BinDir*") {
            [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$BinDir", "User")
            Write-Host "  OK Added $BinDir to User PATH"
            Write-Host "  INFO Restart terminal for PATH to take effect"
        } else {
            Write-Host "  OK PATH already configured"
        }

        # [6/6] Create artifact directories
        Write-Host "[6/6] Creating artifact directories..." -ForegroundColor Yellow
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\changes" -Force | Out-Null
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
        Write-Host "  OK docs/rein/{changes,archive}"

        Write-Host ""
        Write-Host "=== Global Installation Complete (Codex) ===" -ForegroundColor Green
        Write-Host ""
        Write-Host "Installed:"
        Write-Host "  Binary:    $BinDir\rein.exe"
        Write-Host "  Rules:     $ProjectDir\CODEX.md"
        Write-Host "  Config:    $CodexDir\config.toml"
        Write-Host "  Artifacts: $ProjectDir\docs\rein\"
        Write-Host ""
        Write-Host "Codex will read CODEX.md automatically on session start"

    } else {
        $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }
        $BinDir = "$ConfigDir\bin"

        Write-Host "=== rein Global Installer ===" -ForegroundColor Cyan
        Write-Host "Target: $ConfigDir"
        Write-Host ""

        # Check existing installation
        if (Test-Path "$ConfigDir\.rein-manifest") {
            Write-Host "INFO Existing rein installation detected — upgrading" -ForegroundColor Yellow
        fi

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

        # [7/8] Inject task progress rule into CLAUDE.md
        Write-Host "[7/8] Injecting task-progress rule into CLAUDE.md..." -ForegroundColor Yellow
        Inject-ClaudeMd $ProjectDir

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
    }

# ============================================================
# Project Install
# ============================================================
} else {
    $ProjectDir = Get-Location

    if ($Ide -eq "cursor") {
        $CursorDir = "$ProjectDir\.cursor"
        $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }

        Write-Host "=== rein Project Installer (Cursor) ===" -ForegroundColor Cyan
        Write-Host "Target: $ProjectDir"
        Write-Host ""

        # Check existing installation
        if (Test-Path "$CursorDir\.rein-manifest") {
            Write-Host "INFO Existing rein installation detected — upgrading" -ForegroundColor Yellow
        }

        # [1/5] Install/upgrade rein CLI globally
        Write-Host "[1/5] Installing rein CLI..." -ForegroundColor Yellow
        $GlobalBinDir = "$ConfigDir\bin"
        Install-Binary $GlobalBinDir

        # [2/5] Copy resources as .mdc rules
        Write-Host "[2/5] Installing Cursor rules..." -ForegroundColor Yellow
        Copy-Resources-Cursor $CursorDir

        # [3/5] Generate manifest
        Write-Host "[3/5] Generating protection manifest..." -ForegroundColor Yellow
        Generate-Manifest $CursorDir

        # [4/5] Configure hooks.json
        Write-Host "[4/5] Configuring hooks.json..." -ForegroundColor Yellow
        $HookCmd = "$GlobalBinDir\rein.exe"
        Configure-CursorHooks "$CursorDir\hooks.json" $HookCmd

        # [5/5] Create artifact directories
        Write-Host "[5/5] Creating artifact directories..." -ForegroundColor Yellow
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\changes" -Force | Out-Null
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
        Write-Host "  OK docs/rein/{changes,archive}"

        Write-Host ""
        Write-Host "=== Project Installation Complete (Cursor) ===" -ForegroundColor Green
        Write-Host ""
        Write-Host "Installed:"
        Write-Host "  Rules:     $CursorDir\rules\"
        Write-Host "  Hooks:     $CursorDir\hooks.json"
        Write-Host "  Artifacts: $ProjectDir\docs\rein\"
        Write-Host ""
        Write-Host "Reference rules in Cursor chat with @<rule-name>"

    } elseif ($Ide -eq "codex") {
        $CodexDir = "$ProjectDir\.codex"
        $ConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { "$env:USERPROFILE\.claude" }

        Write-Host "=== rein Project Installer (Codex) ===" -ForegroundColor Cyan
        Write-Host "Target: $ProjectDir"
        Write-Host ""

        # Check existing installation
        if (Test-Path "$CodexDir\.rein-manifest") {
            Write-Host "INFO Existing rein installation detected — upgrading" -ForegroundColor Yellow
        }

        # [1/5] Install/upgrade rein CLI globally
        Write-Host "[1/5] Installing rein CLI..." -ForegroundColor Yellow
        $GlobalBinDir = "$ConfigDir\bin"
        Install-Binary $GlobalBinDir

        # [2/5] Copy resources as CODEX.md
        Write-Host "[2/5] Generating CODEX.md and config..." -ForegroundColor Yellow
        Copy-Resources-Codex $ProjectDir

        # [3/5] Generate manifest
        Write-Host "[3/5] Generating protection manifest..." -ForegroundColor Yellow
        Generate-Manifest $CodexDir

        # [4/5] Configure config.toml
        Write-Host "[4/5] Configuring .codex/config.toml..." -ForegroundColor Yellow
        $HookCmd = "$GlobalBinDir\rein.exe"
        Configure-CodexConfig "$CodexDir\config.toml" $HookCmd

        # [5/5] Create artifact directories
        Write-Host "[5/5] Creating artifact directories..." -ForegroundColor Yellow
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\changes" -Force | Out-Null
        New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
        Write-Host "  OK docs/rein/{changes,archive}"

        Write-Host ""
        Write-Host "=== Project Installation Complete (Codex) ===" -ForegroundColor Green
        Write-Host ""
        Write-Host "Installed:"
        Write-Host "  Rules:     $ProjectDir\CODEX.md"
        Write-Host "  Config:    $CodexDir\config.toml"
        Write-Host "  Artifacts: $ProjectDir\docs\rein\"
        Write-Host ""
        Write-Host "Codex will read CODEX.md automatically on session start"

    } else {
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

        # Inject task progress rule into CLAUDE.md
        Inject-ClaudeMd $ProjectDir

        # [6/6] Verification
        Write-Host "[6/6] Verifying installation..." -ForegroundColor Yellow
        if (Test-Path "$GlobalBinDir\rein.exe") {
            try {
                $ver = & "$GlobalBinDir\rein.exe" --version 2>$null
                Write-Host "  OK rein $ver available at $GlobalBinDir\" -ForegroundColor Green
            } catch {
                Write-Host "  OK rein CLI installed at $GlobalBinDir\ (version check blocked by policy)" -ForegroundColor Green
            }
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
}
