# Contributing to Goticks

Thanks for your interest in contributing to Goticks! We welcome issues, pull requests, and discussion.

## How to Contribute

1. Fork the repository.
2. Create a new branch for your change: `git checkout -b feature/my-feature` or `git checkout -b fix/bug`.
3. Make your changes, ensuring code is properly formatted (`go fmt ./...`) and linted (`golangci-lint run`).
4. Add tests for any new functionality or bug fixes (`go test ./...`).
5. Commit your changes and push to your fork.
6. Open a pull request against the `main` branch.

## Style Guide

- Follow the existing code style and conventions.
- Keep functions small and focused.
- Write clear, concise documentation for exported types and functions.

## Running Tests

```bash
go test ./... -v
```  

## Code Reviews

- All changes are reviewed by at least one maintainer.
- Address review comments promptly.

Thank you for helping improve Goticks!
