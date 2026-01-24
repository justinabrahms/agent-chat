## ADDED Requirements

### Requirement: Config File Loading

The system SHALL load configuration from a YAML file at `~/.config/agent-chat/config.yaml` by default.

#### Scenario: Default config path exists
- **WHEN** a config file exists at `~/.config/agent-chat/config.yaml`
- **THEN** the system loads configuration values from that file

#### Scenario: Default config path does not exist
- **WHEN** no config file exists at the default path
- **THEN** the system starts normally using defaults and other configuration sources
- **AND** no error is displayed

### Requirement: Config Flag Override

The system SHALL support a `--config` flag to specify an alternate config file path.

#### Scenario: Custom config path specified
- **WHEN** the user specifies `--config /path/to/config.yaml`
- **THEN** the system loads configuration from that path instead of the default

#### Scenario: Custom config path does not exist
- **WHEN** the user specifies `--config /nonexistent/path.yaml`
- **AND** the file does not exist
- **THEN** the system exits with a helpful error message indicating the file was not found

### Requirement: Configuration Settings

The config file SHALL support the following settings:
- `port`: HTTP server port (integer)
- `gastown-dir`: Path to Gas Town .beads directory (string)
- `multiclaude-dir`: Path to multiclaude directory (string)

#### Scenario: Valid config file
- **WHEN** the config file contains valid YAML with supported settings
- **THEN** those settings are applied to the application

#### Scenario: Invalid YAML syntax
- **WHEN** the config file contains invalid YAML syntax
- **THEN** the system exits with an error message describing the YAML parse error and line number

#### Scenario: Unknown config keys
- **WHEN** the config file contains unknown keys
- **THEN** the system ignores them and logs a warning

### Requirement: Configuration Precedence

Configuration values SHALL be resolved in this order (highest to lowest priority):
1. Command-line flags
2. Environment variables
3. Config file values
4. Default values

#### Scenario: Flag overrides config file
- **WHEN** `--port 9000` is specified on command line
- **AND** config file contains `port: 8000`
- **THEN** port 9000 is used

#### Scenario: Env var overrides config file
- **WHEN** `PORT=9000` environment variable is set
- **AND** config file contains `port: 8000`
- **AND** no `--port` flag is specified
- **THEN** port 9000 is used

#### Scenario: Config file overrides default
- **WHEN** config file contains `port: 8000`
- **AND** no `--port` flag or `PORT` env var is set
- **THEN** port 8000 is used

### Requirement: Helpful Error Messages

The system SHALL provide user-friendly error messages for config file issues.

#### Scenario: YAML parse error
- **WHEN** config file has invalid YAML
- **THEN** error message includes the filename, line number, and description of the syntax error

#### Scenario: Permission denied
- **WHEN** config file exists but is not readable
- **THEN** error message indicates permission was denied and suggests checking file permissions
