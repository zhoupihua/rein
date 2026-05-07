# Using rein with GitHub Copilot

## Setup

### Copilot Instructions

Copilot supports creating agent skills using a `.github/skills` directory in your repository.

```bash
mkdir -p .github/skills

# Copy essential rein skills
cp -r .claude/skills/define .github/skills/
cp -r .claude/skills/tdd .github/skills/
cp -r .claude/skills/code-review .github/skills/
```

For more details, refer [Creating agent skills for GitHub Copilot](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-skills).

### Agent Personas

Copilot supports specialized agent personas:

```bash
mkdir -p .github/agents
cp .claude/agents/code-reviewer.md .github/agents/code-reviewer.md
cp .claude/agents/test-engineer.md .github/agents/test-engineer.md
cp .claude/agents/security-auditor.md .github/agents/security-auditor.md
```

Invoke agents in Copilot Chat:
- `@code-reviewer Review this PR`
- `@test-engineer Analyze test coverage for this module`
- `@security-auditor Check this endpoint for vulnerabilities`

### Custom Instructions (User Level)

For skills you want across all repositories:

1. Open VS Code → Settings → GitHub Copilot → Custom Instructions
2. Add your most-used skill summaries

## Recommended Configuration

### .github/copilot-instructions.md

```markdown
# Project Coding Standards

## Testing
- Write tests before code (TDD)
- For bugs: write a failing test first, then fix (Prove-It pattern)
- Test hierarchy: unit > integration > e2e
- Run `go test ./...` after every change

## Code Quality
- Review across five axes: correctness, readability, architecture, security, performance
- Every PR must pass: lint, type check, tests, build
- No secrets in code or version control

## Implementation
- Build in small, verifiable increments
- Each increment: implement → test → verify → commit
- Never mix formatting changes with behavior changes

## Boundaries
- Always: Run tests before commits, validate user input
- Ask first: Database schema changes, new dependencies
- Never: Commit secrets, remove failing tests, skip verification
```

## Usage Tips

1. **Keep instructions concise** — Copilot instructions work best when focused
2. **Use agents for review** — The code-reviewer, test-engineer, and security-auditor agents are designed for Copilot's agent model
3. **Reference in chat** — When working on a specific phase, paste the relevant skill content into Copilot Chat
4. **Combine with PR reviews** — Set up Copilot to review PRs using the code-reviewer agent persona
