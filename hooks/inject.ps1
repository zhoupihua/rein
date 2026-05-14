# inject Hook (UserPromptExpansion → /code-review)
# Inject team review checklist into /code-review command
$ChecklistPath = Join-Path $env:CLAUDE_PROJECT_DIR ".rein\checklists\review.md"
if (Test-Path $ChecklistPath) {
    Get-Content $ChecklistPath -Raw
}
