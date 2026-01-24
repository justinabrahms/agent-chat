## 1. Implementation

- [x] 1.1 Create `internal/config/config.go` with Config struct and Load function
- [x] 1.2 Add YAML parsing with `gopkg.in/yaml.v3`
- [x] 1.3 Implement default config path resolution (`~/.config/agent-chat/config.yaml`)
- [x] 1.4 Add helpful error messages for missing/invalid config files
- [x] 1.5 Add `--config` flag to main.go for alternate config path
- [x] 1.6 Integrate config loading into main.go with proper precedence (flags > env > config > defaults)
- [x] 1.7 Add tests for config loading (valid, invalid, missing scenarios)
