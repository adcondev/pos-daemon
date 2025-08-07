# Conventional Commit Message Generation Guidelines

## Overview

You are tasked with generating Conventional Commit messages for code changes. These commit messages MUST be in English
and follow the specification below.

## Format

```
<type>[(optional scope)][!]: <description>

[optional body]

[optional footer(s)]
```

## Types

- `feat`: A new feature for the user/application (MUST be used when adding new functionality)
- `fix`: A bug fix (MUST be used when fixing bugs)
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit

## Optional Scope

Scope MUST be a noun describing a section of the codebase surrounded by parentheses. Common scopes include:

- `(encoding)`
- `(connector)`
- `(printer)`
- `(imaging)`
- `(profiles)`
- `(protocol)`
- `(escpos)`
- `(errors)`
- `(logs)`
- `(config)`
- `(api)`
- `(models)`
- `(service)`
- `(utils)`

## Description

- MUST be imperative, present tense: "change" not "changed" nor "changes"
- MUST be lowercase
- MUST NOT end with a period
- MUST be concise (under 50 characters)

## Breaking Changes

Breaking changes MUST be indicated in one of two ways:

1. Adding `!` before the colon: `feat(api)!: remove deprecated endpoints`
2. Adding `BREAKING CHANGE:` in the footer:
   `BREAKING CHANGE: environment variables now take precedence over config files`

## Body

- Use the body to explain WHAT and WHY (not HOW)
- Separate paragraphs with blank lines
- Use bullet points with hyphens (`-`)
- SHOULD wrap at 72 characters

## Footers

- Footers MUST be separated from the body by a blank line
- Each footer MUST consist of a token, followed by either `: ` or ` #`
- Common footers include:
    - `Fixes: #123`
    - `Reviewed-by: Person Name`
    - `Refs: #456`
    - `BREAKING CHANGE: description of breaking change`

## Examples

### Simple feature

```
feat(auth): add OAuth2 authentication support
```

### Bug fix with scope

```
fix(printer): correct paper feed calculation
```

### Feature with breaking change

```
feat(api)!: remove deprecated endpoints

All v1 API endpoints have been removed as previously announced.
Users should migrate to v2 API endpoints.

BREAKING CHANGE: v1 API is no longer available
```

### Documentation update with body

```
docs(readme): update installation instructions

- Add Docker setup
- Update minimum requirements
- Fix typos in configuration section
```

## Instructions for AI

1. Analyze the code changes to determine the correct type and scope
2. Write a clear, concise description in imperative mood
3. Include a body explaining WHY the change was made and WHAT it accomplishes
4. Add relevant footers (especially for issue references or breaking changes)
5. Ensure the message is in English
6. Keep the first line (type+scope+description) under 50 characters
7. Wrap the body text at 72 characters

When presented with code changes, first identify the primary purpose of the change to select the appropriate type and
scope, then follow this format to generate a conventional commit message.

```