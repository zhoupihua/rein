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
Write-Host "[1/9] Creating artifact directories..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\specs" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\changes" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\archive" -Force | Out-Null
Write-Host "  OK specs/, changes/, archive/"

# 2. Copy commands
Write-Host "[2/9] Installing commands..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\commands" -Force | Out-Null
Copy-Item "$WorkflowDir\commands\*.md" "$ProjectDir\.claude\commands\"
$CmdCount = (Get-ChildItem "$ProjectDir\.claude\commands\*.md").Count
Write-Host "  OK $CmdCount commands installed"

# 3. Copy skills
Write-Host "[3/9] Installing skills..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\skills" -Force | Out-Null
Copy-Item -Path "$WorkflowDir\skills\*" -Destination "$ProjectDir\.claude\skills\" -Recurse -Force
$SkillCount = (Get-ChildItem "$ProjectDir\.claude\skills" -Directory).Count
Write-Host "  OK $SkillCount skills installed"

# 4. Copy agents
Write-Host "[4/9] Installing agents..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\agents" -Force | Out-Null
Copy-Item "$WorkflowDir\agents\*.md" "$ProjectDir\.claude\agents\"
$AgentCount = (Get-ChildItem "$ProjectDir\.claude\agents\*.md").Count
Write-Host "  OK $AgentCount agents installed"

# 5. Copy hooks
Write-Host "[5/9] Installing hooks..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\hooks" -Force | Out-Null
Copy-Item "$WorkflowDir\hooks\session-start.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\session-start.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\format.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\format.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\test-gateway.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\test-gateway.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\secret-scan.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\secret-scan.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\context-inject.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\context-inject.ps1" "$ProjectDir\.claude\hooks\"
Write-Host "  OK All hooks installed (session-start, format, test-gateway, secret-scan, context-inject)"

# 6. Copy checklists
Write-Host "[6/9] Installing checklists..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\checklists" -Force | Out-Null
if (Test-Path "$WorkflowDir\templates\checklists\review.md") {
    Copy-Item "$WorkflowDir\templates\checklists\review.md" "$ProjectDir\.claude\checklists\" -Force
    Write-Host "  OK review.md checklist installed"
} else {
    Write-Host "  SKIP No review checklist template found"
}

# 7. Configure settings.json
Write-Host "[7/9] Configuring hooks in settings.json..." -ForegroundColor Yellow
$SettingsFile = "$ProjectDir\.claude\settings.json"
$HookBase = 'powershell -ExecutionPolicy Bypass -File "${CLAUDE_PROJECT_DIR}\.claude\hooks'
if (Test-Path $SettingsFile) {
    Write-Host "  INFO settings.json exists - merge hooks manually if needed"
    Write-Host "  See hooks/hooks.json for the full configuration template"
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
                            command = "$HookBase\session-start.ps1""
                        }
                    )
                }
            )
            PostToolUse = @(
                @{
                    matcher = "Write|Edit|MultiEdit"
                    hooks = @(
                        @{
                            type = "command"
                            command = "$HookBase\format.ps1""
                        }
                    )
                }
                @{
                    matcher = "Read|Bash"
                    hooks = @(
                        @{
                            type = "command"
                            command = "$HookBase\secret-scan.ps1""
                        }
                    )
                }
            )
            PreToolUse = @(
                @{
                    matcher = "Bash"
                    hooks = @(
                        @{
                            type = "command"
                            command = "$HookBase\test-gateway.ps1""
                        }
                    )
                }
            )
            UserPromptExpansion = @(
                @{
                    matcher = "review"
                    hooks = @(
                        @{
                            type = "command"
                            command = "$HookBase\context-inject.ps1""
                        }
                    )
                }
            )
        }
    } | ConvertTo-Json -Depth 10
    Set-Content -Path $SettingsFile -Value $Settings
    Write-Host "  OK settings.json created with all hooks"
}

# 8. Append workflow instructions to CLAUDE.md
Write-Host "[8/9] Updating CLAUDE.md..." -ForegroundColor Yellow
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

# 9. Handle AGENTS.md (Codex CLI compatibility)
Write-Host "[9/9] Checking for Codex CLI..." -ForegroundColor Yellow
$AgentsMd = "$ProjectDir\AGENTS.md"
if (Test-Path $AgentsMd) {
    Write-Host "  INFO AGENTS.md found - Codex CLI detected"
} else {
    Write-Host "  INFO No AGENTS.md found - skipping Codex CLI setup"
}

Write-Host ""
Write-Host "=== Installation Complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "Installed hooks:" -ForegroundColor Cyan
Write-Host "  1. SessionStart  - Inject using-workflow skill"
Write-Host "  2. Format        - Auto-format with Prettier (PostToolUse: Write|Edit|MultiEdit)"
Write-Host "  3. Test Gateway  - Run tests before deploy (PreToolUse: Bash)"
Write-Host "  4. Secret Scan   - Block secrets in output (PostToolUse: Read|Bash)"
Write-Host "  5. Context Inject- Inject review checklist (UserPromptExpansion: /review)"
Write-Host ""
Write-Host "Verification steps:"
Write-Host "1. Start a new Claude Code session"
Write-Host "2. The using-workflow skill should be auto-injected"
Write-Host "3. Try /triage to test the workflow"
Write-Host "4. Try /review to test checklist injection"
