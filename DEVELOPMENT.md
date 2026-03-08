# Development Setup

## Git Hooks

We use Git pre-commit hooks to catch issues locally before pushing to CI/CD.

### Installing Hooks

Run the setup script from the repository root:

```bash
./scripts/setup-hooks.sh
```

Or manually:

```bash
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### What the Pre-Commit Hook Checks

The pre-commit hook runs the same checks as our CI/CD pipeline:

1. **go mod tidy** - Ensures go.mod and go.sum are up to date
2. **go vet** - Reports suspicious constructs
3. **golangci-lint** - Runs comprehensive linting (if installed)
4. **go build** - Verifies the code compiles

### Bypassing Hooks

If you need to commit without running hooks (not recommended):

```bash
git commit --no-verify
```

### Installing golangci-lint

The pre-commit hook will skip golangci-lint if it's not installed. To install it:

**macOS:**
```bash
brew install golangci-lint
```

**Linux:**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

**Windows:**
```powershell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Why Use Pre-Commit Hooks?

- **Catch issues early** - Find problems before CI/CD
- **Save time** - No waiting for CI to fail on simple issues
- **Better code quality** - Consistent checks for all commits
- **Learn best practices** - Immediate feedback on code issues

## CI/CD Pipeline

Our GitHub Actions workflow (`.github/workflows/backend.yml`) runs:

- All the same checks as pre-commit hooks
- Unit tests
- Integration tests (if configured)

Having pre-commit hooks means you'll catch most issues locally before pushing.
