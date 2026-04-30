# inject Hook (UserPromptExpansion → /code-review)
# Inject team review checklist into /code-review command
if [ -f "${CLAUDE_PROJECT_DIR}/.claude/checklists/review.md" ]; then
    cat "${CLAUDE_PROJECT_DIR}/.claude/checklists/review.md"
fi
