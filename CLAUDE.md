<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

## CI Requirements

**Spec Requirement:** Substantial Go code changes require specification documents. CI will fail if you change >50 lines of Go code or >3 Go files without corresponding changes in `openspec/changes/` or `openspec/specs/`. See `openspec/AGENTS.md` for details on creating proposals.