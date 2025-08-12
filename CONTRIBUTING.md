# Development Workflow Guide

## Overview

This document describes the development workflow for the pos-daemon project. We follow a pull request-based workflow
with automated CI/CD, conventional commits, and semantic versioning.

## Workflow Principles

1. **All changes go through pull requests** - Direct pushes to main/master are disabled
2. **Conventional commits are mandatory** - Ensures automatic versioning and changelog generation
3. **CI must pass** - All tests, linting, and security checks must pass before merging
4. **Code review is required** - At least one approval needed (except for dependabot patches)

## Development Process

### 1. Setting Up Your Environment

```bash
# Clone the repository
git clone https://github.com/AdConDev/pos-daemon.git
cd pos-daemon

# Install Go (1.23 or later)
# See: https://golang.org/doc/install

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
npm install -g @commitlint/cli @commitlint/config-conventional

# Verify setup
go version
golangci-lint version
commitlint --version
```

### 2. Creating a Feature Branch

```bash
# Always branch from the latest main/master
git checkout main
git pull origin main

# Create a feature branch
# Use prefixes: feat/, fix/, docs/, refactor/, test/, ci/
git checkout -b feat/add-new-printer-support

# For bugs
git checkout -b fix/connection-timeout-issue
```

### 3. Making Changes

Follow these guidelines:

1. **Write tests first** (TDD approach)
2. **Keep commits atomic** - One logical change per commit
3. **Run tests locally** before pushing:
   ```bash
   go test -race -cover ./...
   golangci-lint run
   ```

### 4. Committing Changes

We use [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Feature
git commit -m "feat(printer): add support for Epson TM-T88VII"

# Bug fix
git commit -m "fix(connector): resolve timeout on slow networks"

# Breaking change (triggers major version)
git commit -m "feat(api)!: redesign printer configuration interface"

# With scope
git commit -m "docs(readme): update installation instructions"

# With body and footer
git commit -m "fix(encoding): handle UTF-8 characters properly

This fixes an issue where special characters were not being
encoded correctly for certain printer models.

Fixes #123"
```

### 5. Creating a Pull Request

1. Push your branch:
   ```bash
   git push origin feat/your-feature
   ```

2. Go to GitHub and create a PR

3. **PR Title MUST follow conventional commits format**:
    - ✅ `feat: add support for new printer model`
    - ✅ `fix: resolve connection timeout issue`
    - ❌ `Added new feature` (wrong format)
    - ❌ `Fix bug` (too vague)

4. Fill out the PR template completely

5. Wait for CI checks to pass

### 6. Code Review Process

- **Authors**: Respond to feedback promptly
- **Reviewers**:
    - Check for bugs, performance issues, and maintainability
    - Ensure tests are adequate
    - Verify conventional commit format
    - Be constructive and specific

### 7. Merging

Once approved and CI passes:

- The PR will be automatically merged using squash merge
- The squash commit message will use the PR title (must be conventional format)
- A release will be automatically created if the commit type triggers one

## Commit Types and Versioning

| Type       | Description             | Version Bump  | Changelog |
|------------|-------------------------|---------------|-----------|
| `feat`     | New feature             | Minor (0.X.0) | ✅         |
| `fix`      | Bug fix                 | Patch (0.0.X) | ✅         |
| `docs`     | Documentation only      | None          | ❌         |
| `style`    | Code style (formatting) | None          | ❌         |
| `refactor` | Code refactoring        | None          | ❌         |
| `perf`     | Performance improvement | Patch (0.0.X) | ✅         |
| `test`     | Add/update tests        | None          | ✅         |
| `build`    | Build system changes    | None          | ✅         |
| `ci`       | CI/CD changes           | None          | ✅         |
| `chore`    | Other changes           | None          | ❌         |
| `revert`   | Revert previous commit  | Varies        | ✅         |
| `deps`     | Dependency updates      | Patch (0.0.X) | ✅         |

**Breaking Changes**: Add `!` after type or include `BREAKING CHANGE:` in commit body → Major (X.0.0)

## CI/CD Pipeline

Our CI pipeline runs on every PR and includes:

1. **Validation**
    - PR title format
    - Commit message format

2. **Testing**
    - Unit tests on multiple OS/Go versions
    - Race condition detection
    - Coverage reporting

3. **Code Quality**
    - golangci-lint checks
    - go mod tidy verification

4. **Security**
    - Gosec security scanning
    - Trivy vulnerability scanning
    - Dependency review

5. **Automation**
    - Auto-labeling based on files changed
    - PR size labeling
    - Conflict detection

## Release Process

Releases are **fully automated**:

1. When a PR is merged to main/master
2. The system analyzes commits since last release
3. If releasable changes exist (feat, fix, perf):
    - Version is bumped according to commit types
    - CHANGELOG.md is updated
    - Git tag is created
    - GitHub Release is published
    - Go module proxy is notified

## Working Alone vs Team

### Solo Development

- You can work directly on main/master (if branch protection allows)
- Still use conventional commits for automatic releases
- Consider using PRs anyway for CI validation

### Team Development (2+ people)

- **Always use pull requests**
- Require code reviews
- Use branch protection rules
- Communicate in PR comments
- Assign PRs to reviewers

## Common Scenarios

### Hotfix Process

```bash
# Branch from main
git checkout -b fix/critical-bug main

# Make fix with proper commit
git commit -m "fix: prevent data corruption on disconnect"

# Push and create PR marked as urgent
git push origin fix/critical-bug
```

### Feature Development

```bash
# Long-running feature branch
git checkout -b feat/major-feature

# Regular commits as you work
git commit -m "feat(parser): add basic parsing logic"
git commit -m "test(parser): add unit tests"
git commit -m "docs(parser): add API documentation"

# Keep branch updated
git fetch origin
git rebase origin/main
```

### Dependency Updates

- Dependabot creates PRs automatically
- Patches are auto-merged if tests pass
- Minor/major updates require manual review

## Best Practices

1. **Write meaningful commit messages** - They become the changelog
2. **Keep PRs small** - Easier to review, less risk
3. **Update tests** - Maintain or increase coverage
4. **Document public APIs** - Use godoc conventions
5. **Respond to CI failures quickly** - Don't leave PRs in failed state
6. **Use draft PRs** - For work in progress
7. **Link issues** - Use "Fixes #123" in PR description

## Troubleshooting

### CI Failures

1. **Linting errors**: Run `golangci-lint run` locally
2. **Test failures**: Check OS-specific issues, race conditions
3. **Commit validation**: Ensure conventional format
4. **Security issues**: Review Gosec/Trivy reports

### Release Issues

- Check workflow logs in Actions tab
- Ensure commits follow conventional format
- Verify no `[skip release]` in commit messages
- Check branch protection settings

## Questions?

- Check existing issues/discussions
- Read the codebase and tests
- Ask in PR comments
- Create a discussion for broader topics
