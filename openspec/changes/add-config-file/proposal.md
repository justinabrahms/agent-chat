# Change: Add config file support

## Why

Currently, agent-chat requires users to pass command-line flags or set environment variables for configuration. A config file at `~/.config/agent-chat/config.yaml` would provide a more convenient way to persist settings across sessions.

## What Changes

- Add YAML config file loading from `~/.config/agent-chat/config.yaml`
- Support `port`, `gastown-dir`, and `multiclaude-dir` settings
- Add `--config` flag to specify alternate config file path
- Provide helpful error messages for invalid config files
- Establish precedence: flags > env vars > config file > defaults

## Impact

- Affected specs: config (new capability)
- Affected code: `cmd/agent-chat/main.go`, new `internal/config/` package
