## ADDED Requirements

### Requirement: CI Workflow on Pull Requests
The project SHALL have a GitHub Actions CI workflow that runs automatically on pull request events.

#### Scenario: PR triggers CI
- **WHEN** a pull request is opened or updated
- **THEN** the CI workflow runs automatically

### Requirement: Go Test Execution
The CI workflow SHALL execute Go tests to verify code correctness.

#### Scenario: Tests pass
- **WHEN** the CI workflow runs
- **THEN** `go test ./...` executes successfully

#### Scenario: Tests fail
- **WHEN** a test fails
- **THEN** the CI workflow fails and reports the error

### Requirement: Go Build Verification
The CI workflow SHALL verify that the Go code compiles successfully.

#### Scenario: Build succeeds
- **WHEN** the CI workflow runs
- **THEN** `go build ./...` completes without errors

### Requirement: Code Quality Linting
The CI workflow SHALL run golangci-lint to enforce code quality standards.

#### Scenario: Lint passes
- **WHEN** the CI workflow runs
- **THEN** golangci-lint checks pass

#### Scenario: Lint fails
- **WHEN** linting issues are found
- **THEN** the CI workflow fails and reports the issues

### Requirement: Go Version Consistency
The CI workflow SHALL use the Go version specified in go.mod.

#### Scenario: Go version matches go.mod
- **WHEN** the CI workflow runs
- **THEN** Go 1.24.1 is used (matching go.mod)
