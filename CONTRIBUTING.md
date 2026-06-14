# Contributing to Still

Thank you for your interest in Still. This document explains how to contribute to
the open-source Community Edition.

## Before you contribute

By submitting a pull request, you agree that your contributions will be licensed
under the Apache License 2.0 and you grant the project maintainers the rights
described in our [Contributor License Agreement](CLA.md). We may ask you to sign
a digital CLA before a large PR is merged.

## Getting started

1. Fork the repository.
2. Clone your fork and create a feature branch from `main`.
3. Run `make install` to install dependencies.
4. Run `make env` to generate local environment files.
5. Run `make infra` to start PostgreSQL and MinIO.
6. Run `make backend` and `make mobile` in separate terminals.

See [README.md](README.md) for more details.

## Code style

- **Go**: standard formatting (`gofmt`), run `go vet ./...`.
- **TypeScript**: strict mode, run `yarn lint`.
- **Protocol Buffers**: define APIs in `proto/still/v1/` first, then run `make proto`.

## Pull request process

1. Open an issue or discussion before large changes.
2. Keep PRs focused on a single concern.
3. Add tests when applicable.
4. Ensure `make lint` and `make test` pass.
5. Update relevant documentation (`docs/` or `README.md`).
6. Fill out the PR template if one is provided.

## What belongs in this repo

This is the Community Edition. Contributions should benefit self-hosters and the
open-source community. Proprietary or SaaS-only features belong in the separate
commercial repository.

## Community

- GitHub Discussions: open a thread for questions or ideas.
- Issues: report bugs or request features.

## Code of conduct

Be respectful, constructive, and inclusive. Harassment or discriminatory behavior
will not be tolerated.
