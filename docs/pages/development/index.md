# Contributing

We welcome contributions!

- Fork the repo and create a feature branch.
- Add tests for new rules or features.
- Run `make test` before submitting a PR.
- Ensure your commits are signed.

Thank you for improving tfcoach!

## Semantic Commits

We are using [conventional commits](https://www.conventionalcommits.org/en/v1.0.0-beta.4/) to release this project.

To streamline the whole process, we have enabled **squash-commits** on merge. So you just need to name your PR
correctly.

## Pre-Commit

We also have defined some [pre-commit](https://pre-commit.com/) rules. You can install these hooks via
`make init-precommit`, provided that you have installed the `pre-commit` tool beforehand.
