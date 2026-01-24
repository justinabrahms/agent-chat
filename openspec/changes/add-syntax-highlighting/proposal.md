# Change: Add syntax highlighting for code blocks

## Why
Code blocks in chat messages are currently displayed as plain text without any visual differentiation for programming languages. Syntax highlighting makes code more readable and easier to understand at a glance.

## What Changes
- Add highlight.js library via CDN
- Parse code blocks with language hints (e.g., ```go, ```python)
- Apply syntax highlighting to detected code blocks
- Support both dark and light themes
- Auto-detect language when no hint is provided

## Impact
- Affected specs: message-display (new capability)
- Affected code: `internal/server/templates/index.html`, `internal/server/static/style.css`
