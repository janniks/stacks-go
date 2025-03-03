# Contributing to stacks-go

Thank you for your interest in contributing to the stacks-go project! This document provides guidelines to make the contribution process smooth and effective.

## Conventional Commits

This project uses the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages. This leads to more readable messages that are easy to follow when looking through the project history and helps with automatic version bumping and changelog generation.

### Commit Message Format

Each commit message consists of a **header**, an optional **body**, and an optional **footer**:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

#### Type

Must be one of the following:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries

#### Scope

The scope is optional and should be a noun describing a section of the codebase (e.g., "parser", "crypto", "serialization").

#### Subject

The subject contains a succinct description of the change:

- Use the imperative, present tense: "change" not "changed" nor "changes"
- Don't capitalize the first letter
- No period (.) at the end

#### Body

The body should include the motivation for the change and contrast this with previous behavior.

#### Footer

The footer should contain any information about **Breaking Changes** and is also the place to reference GitHub issues.

### Examples

```
feat(parser): add ability to parse binary stacks transactions

This introduces a transaction parser that can decode binary-encoded Stacks transactions.

Closes #123
```

```
fix(crypto): correct signature verification for multisig transactions

The verification logic was incorrectly checking the threshold condition.

Fixes #456
```

## Pull Request Process

1. Ensure your PR includes tests for any new functionality.
2. Update the README.md or documentation with details of changes if applicable.
3. The PR should work with the CI checks passing.
4. Your PR will be reviewed by maintainers, who might request changes.

Thank you for contributing to stacks-go!
