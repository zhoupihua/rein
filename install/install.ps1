# Alloy install script (Windows)
# Run from your project root: powershell -ExecutionPolicy Bypass -File \path\to\Alloy\install\install.ps1

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WorkflowDir = Split-Path -Parent $ScriptDir
$ProjectDir = Get-Location

Write-Host "=== Alloy Installer ===" -ForegroundColor Cyan
Write-Host "Workflow source: $WorkflowDir"
Write-Host "Target project:  $ProjectDir"
Write-Host ""

# 1. Create artifact directories
Write-Host "[1/8] Creating artifact directories..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\specs" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\changes" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\archive" -Force | Out-Null
Write-Host "  OK specs/, changes/, archive/"

# 2. Copy commands
Write-Host "[2/8] Installing commands..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\commands" -Force | Out-Null
Copy-Item "$WorkflowDir\commands\*.md" "$ProjectDir\.claude\commands\"
$CmdCount = (Get-ChildItem "$ProjectDir\.claude\commands\*.md").Count
Write-Host "  OK $CmdCount commands installed"

# 3. Copy skills
Write-Host "[3/8] Installing skills..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\skills" -Force | Out-Null
Copy-Item -Path "$WorkflowDir\skills\*" -Destination "$ProjectDir\.claude\skills\" -Recurse -Force
$SkillCount = (Get-ChildItem "$ProjectDir\.claude\skills" -Directory).Count
Write-Host "  OK $SkillCount skills installed"

# 4. Copy agents
Write-Host "[4/8] Installing agents..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\agents" -Force | Out-Null
Copy-Item "$WorkflowDir\agents\*.md" "$ProjectDir\.claude\agents\"
$AgentCount = (Get-ChildItem "$ProjectDir\.claude\agents\*.md").Count
Write-Host "  OK $AgentCount agents installed"

# 5. Copy hooks
Write-Host "[5/8] Installing hooks..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\hooks" -Force | Out-Null
Copy-Item "$WorkflowDir\hooks\session-start.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\session-start.ps1" "$ProjectDir\.claude\hooks\"
Write-Host "  OK Hooks installed"

# 6. Configure settings.json
Write-Host "[6/8] Configuring hooks in settings.json..." -ForegroundColor Yellow
$SettingsFile = "$ProjectDir\.claude\settings.json"
if (Test-Path $SettingsFile) {
    Write-Host "  INFO settings.json exists - merge hooks manually if needed"
} else {
    New-Item -ItemType Directory -Path "$ProjectDir\.claude" -Force | Out-Null
    $Settings = @{
        hooks = @{
            SessionStart = @(
                @{
                    matcher = ""
                    hooks = @(
                        @{
                            type = "command"
                            command = 'powershell -ExecutionPolicy Bypass -File "${CLAUDE_PROJECT_DIR}\.claude\hooks\session-start.ps1"'
                        }
                    )
                }
            )
        }
    } | ConvertTo-Json -Depth 5
    Set-Content -Path $SettingsFile -Value $Settings
    Write-Host "  OK settings.json created with session-start hook"
}

# 7. Append workflow instructions to CLAUDE.md
Write-Host "[7/8] Updating CLAUDE.md..." -ForegroundColor Yellow
$ClaudeMd = "$ProjectDir\CLAUDE.md"
$WorkflowBlock = @"

## Alloy

This project uses Alloy for structured AI-assisted development.

### Commands
- ``/triage`` - Classify a change as L1/L2/L3
- ``/quick`` - L1: <=5 lines, no logic impact
- ``/fix`` - L2: 1-3 files, clear requirements
- ``/feature`` - L3: Full 8-step workflow
- ``/spec`` - Generate change artifacts
- ``/plan`` - Task breakdown
- ``/build`` - Execute tasks from tasks.md
- ``/test`` - TDD workflow
- ``/review`` - Five-axis code review
- ``/ship`` - Fan-out review + GO/NO-GO
- ``/simplify`` - Code simplification
- ``/resume`` - Resume from breakpoint

### Artifact Directories
- ``specs/`` - Published specs (long-lived)
- ``changes/`` - Active changes (short-lived)
- ``archive/`` - Archived changes
"@

if (Test-Path $ClaudeMd) {
    $Content = Get-Content $ClaudeMd -Raw
    if ($Content -notmatch "Alloy") {
        Add-Content -Path $ClaudeMd -Value $WorkflowBlock
        Write-Host "  OK Workflow instructions appended to CLAUDE.md"
    } else {
        Write-Host "  INFO CLAUDE.md already contains Alloy section"
    }
} else {
    Set-Content -Path $ClaudeMd -Value "# CLAUDE.md$WorkflowBlock"
    Write-Host "  OK CLAUDE.md created with workflow instructions"
}

# 8. Handle AGENTS.md (Codex CLI compatibility)
Write-Host "[8/8] Checking for Codex CLI..." -ForegroundColor Yellow
$AgentsMd = "$ProjectDir\AGENTS.md"
if (Test-Path $AgentsMd) {
    Write-Host "  INFO AGENTS.md found - Codex CLI detected"
} else {
    Write-Host "  INFO No AGENTS.md found - skipping Codex CLI setup"
}

Write-Host ""
Write-Host "=== Installation Complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "Verification steps:"
Write-Host "1. Start a new Claude Code session"
Write-Host "2. The using-workflow skill should be auto-injected"
Write-Host "3. Try /triage to test the workflow"
Write-Host "4. Try /spec test-feature to test artifact generation"