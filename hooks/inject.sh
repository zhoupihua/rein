# Context Injection Hook (UserPromptExpansion → /review)
# Inject team review checklist into /review command
if [ -f "${CLAUDE_PROJECT_DIR}/.claude/checklists/review.md" ]; then
    cat "${CLAUDE_PROJECT_DIR}/.claude/checklists/review.md"
fi
