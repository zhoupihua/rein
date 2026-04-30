# inject Hook (UserPromptExpansion → /review)
# Inject team review checklist into /review command
$ChecklistPath = Join-Path $env:CLAUDE_PROJECT_DIR ".claude\checklists\review.md"
if (Test-Path $ChecklistPath) {
    Get-Content $ChecklistPath -Raw
}
