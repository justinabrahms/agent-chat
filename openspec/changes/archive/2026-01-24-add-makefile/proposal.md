# Change: Add Makefile with Standard Build Targets

## Why
The project lacks a consistent build system interface. A Makefile provides a standard, discoverable way to build, test, lint, and clean the project across different development environments and CI systems.

## What Changes
- Add Makefile with targets: build, test, lint, clean, install, help
- Support cross-platform builds for Linux and macOS (darwin)
- Auto-generate help from target comments

## Impact
- Affected specs: build-system (new capability)
- Affected code: Root directory (new Makefile)
- No breaking changes to existing functionality
