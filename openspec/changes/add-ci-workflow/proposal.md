# Change: Add GitHub Actions CI Workflow

## Why
The project currently has no continuous integration. PRs can be merged without verifying that tests pass, the code compiles, or linting checks succeed. This increases the risk of introducing regressions and code quality issues.

## What Changes
- Add `.github/workflows/ci.yml` with GitHub Actions workflow
- Run Go tests on PR events
- Run golangci-lint for code quality checks
- Verify build succeeds
- Use Go 1.24.1 (from go.mod)

## Impact
- Affected specs: ci (new capability)
- Affected code: `.github/workflows/ci.yml` (new file)
- No breaking changes to existing functionality
