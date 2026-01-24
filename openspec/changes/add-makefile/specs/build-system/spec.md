## ADDED Requirements

### Requirement: Makefile Build Targets
The project SHALL provide a Makefile with standard targets for common development operations.

#### Scenario: Display available targets
- **WHEN** user runs `make` or `make help`
- **THEN** all available targets with descriptions are displayed

#### Scenario: Build the project
- **WHEN** user runs `make build`
- **THEN** the Go binary is compiled to `./bin/agent-chat`

#### Scenario: Run tests
- **WHEN** user runs `make test`
- **THEN** all Go tests are executed with verbose output

#### Scenario: Run linter
- **WHEN** user runs `make lint`
- **THEN** code is checked with `go vet` and staticcheck (if available)

#### Scenario: Clean build artifacts
- **WHEN** user runs `make clean`
- **THEN** the `./bin/` directory is removed

#### Scenario: Install binary
- **WHEN** user runs `make install`
- **THEN** the binary is installed to `$GOPATH/bin` or `$HOME/go/bin`

### Requirement: Cross-Platform Build Support
The Makefile SHALL support building for Linux and macOS platforms.

#### Scenario: Build for current platform
- **WHEN** user runs `make build`
- **THEN** binary is built for the current OS and architecture

#### Scenario: Build for Linux
- **WHEN** user runs `make build-linux`
- **THEN** binary is cross-compiled for Linux amd64

#### Scenario: Build for macOS
- **WHEN** user runs `make build-darwin`
- **THEN** binary is cross-compiled for macOS (darwin) amd64
