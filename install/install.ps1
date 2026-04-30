# rein install script (Windows)
# Run from your project root: powershell -ExecutionPolicy Bypass -File \path\to\rein\install\install.ps1

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WorkflowDir = Split-Path -Parent $ScriptDir
$ProjectDir = Get-Location

Write-Host "=== rein Installer ===" -ForegroundColor Cyan
Write-Host "Workflow source: $WorkflowDir"
Write-Host "Target project:  $ProjectDir"
Write-Host ""

# 1. Create artifact directories
Write-Host "[1/9] Creating artifact directories..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\specs" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\plans" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\tasks" -Force | Out-Null
New-Item -ItemType Directory -Path "$ProjectDir\docs\rein\archive" -Force | Out-Null
Write-Host "  OK docs/rein/specs/, docs/rein/plans/, docs/rein/tasks/, docs/rein/archive/"

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
Copy-Item "$WorkflowDir\hooks\gate.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\gate.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\leak-guard.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\leak-guard.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\inject.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\inject.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\guard.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\guard.ps1" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\guard-bash.sh" "$ProjectDir\.claude\hooks\"
Copy-Item "$WorkflowDir\hooks\guard-bash.ps1" "$ProjectDir\.claude\hooks\"
Write-Host "  OK All hooks installed"

# 6. Copy checklists
Write-Host "[6/9] Installing checklists..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$ProjectDir\.claude\checklists" -Force | Out-Null
if (Test-Path "$WorkflowDir\templates\checklists\review.md") {
    Copy-Item "$WorkflowDir\templates\checklists\review.md" "$ProjectDir\.claude\checklists\" -Force
    Write-Host "  OK review.md checklist installed"
} else {
    Write-Host "  SKIP No review checklist template found"
}

# 7. Generate manifest
Write-Host "[7/9] Generating protection manifest..." -ForegroundColor Yellow
$ManifestFile = "$ProjectDir\.claude\.rein-manifest"
$ManifestLines = @(
    "# rein Managed Files - DO NOT EDIT",
    "# These files are protected from modification by the guard hook.",
    "# To allow edits to a specific file, remove its line from this manifest.",
    ""
)
# Enumerate installed files relative to project root
$Dirs = @(
    "$ProjectDir\.claude\hooks",
    "$ProjectDir\.claude\commands",
    "$ProjectDir\.claude\agents",
    "$ProjectDir\.claude\checklists"
)
foreach ($Dir in $Dirs) {
    if (Test-Path $Dir) {
        Get-ChildItem $Dir -File | ForEach-Object {
            $RelPath = $_.FullName.Substring($ProjectDir.Length + 1).Replace('\', '/')
            $ManifestLines += $RelPath
        }
    }
}
# Skills are directories with subdirectories
if (Test-Path "$ProjectDir\.claude\skills") {
    Get-ChildItem "$ProjectDir\.claude\skills" -Directory | ForEach-Object {
        $SkillDir = $_.FullName.Substring($ProjectDir.Length + 1).Replace('\', '/') + "/"
        $ManifestLines += $SkillDir
    }
}
Set-Content -Path $ManifestFile -Value $ManifestLines
$ManifestCount = ($ManifestLines | Where-Object { $_ -notmatch '^\s*#' -and $_ -ne '' }).Count
Write-Host "  OK $ManifestCount entries in .rein-manifest"

# 8. Configure settings.json
Write-Host "[8/9] Configuring hooks in settings.json..." -ForegroundColor Yellow
$SettingsFile = "$ProjectDir\.claude\settings.json"
if (Test-Path $SettingsFile) {
    Write-Host "  INFO settings.json exists - merge hooks manually if needed"
    Write-Host "  See hooks/hooks.json for the full configuration template"
} else {
    New-Item -ItemType Directory -Path "$ProjectDir\.claude" -Force | Out-Null
    $Settings = @'
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\session-start.ps1\""
          }
        ]
      }
    ],
    "PreToolUse": [
      {
        "matcher": "Edit|Write|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\guard.ps1\""
          }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\guard-bash.ps1\""
          },
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\gate.ps1\""
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\format.ps1\""
          }
        ]
      },
      {
        "matcher": "Read|Bash",
        "hooks": [
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\leak-guard.ps1\""
          }
        ]
      }
    ],
    "UserPromptExpansion": [
      {
        "matcher": "code-review",
        "hooks": [
          {
            "type": "command",
            "command": "powershell -ExecutionPolicy Bypass -File \"${CLAUDE_PROJECT_DIR}\\.claude\\hooks\\inject.ps1\""
          }
        ]
      }
    ]
  }
}
'@
    Set-Content -Path $SettingsFile -Value $Settings
    Write-Host "  OK settings.json created with all hooks"
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
Write-Host "  1. session-start   - Inject using-rein skill (SessionStart)"
Write-Host "  2. guard           - Block edits to rein-managed files (PreToolUse: Edit|Write|MultiEdit)"
Write-Host "  3. guard-bash      - Block destructive cmds on rein files (PreToolUse: Bash)"
Write-Host "  4. gate            - Run tests before deploy (PreToolUse: Bash)"
Write-Host "  5. format          - Auto-format with Prettier (PostToolUse: Write|Edit|MultiEdit)"
Write-Host "  6. leak-guard      - Block secrets in output (PostToolUse: Read|Bash)"
Write-Host "  7. inject          - Inject review checklist (UserPromptExpansion: /code-review)"
Write-Host ""
Write-Host "Protection:" -ForegroundColor Cyan
Write-Host "  rein-managed files are listed in .claude/.rein-manifest"
Write-Host "  Edit/Write on these files will be blocked automatically"
Write-Host "  To allow edits, remove the file's entry from the manifest"
Write-Host ""
Write-Host "Verification steps:"
Write-Host "1. Start a new Claude Code session"
Write-Host "2. The using-rein skill should be auto-injected"
Write-Host "3. Try /triage to test the workflow"
Write-Host "4. Try /code-review to test checklist injection"
