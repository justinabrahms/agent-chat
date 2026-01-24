# Change: Add Light Theme with System Auto-Detection

## Why
The chat UI currently only supports dark theme. Users working in bright environments or preferring light interfaces have no option to switch. Adding theme support with system auto-detection provides a better user experience across different lighting conditions and user preferences.

## What Changes
- Add light theme CSS variables alongside existing dark theme
- Add system theme auto-detection using `prefers-color-scheme` media query
- Add UI toggle in header for switching between light/dark/system modes
- Persist theme preference in localStorage

## Impact
- Affected specs: `chatroom-ui` (modified)
- Affected code: `internal/server/static/style.css`, `internal/server/templates/index.html`
- External dependencies: None (uses standard CSS and JS APIs)
