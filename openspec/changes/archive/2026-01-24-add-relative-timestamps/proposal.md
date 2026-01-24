# Add Relative Timestamps

## Summary

Replace static HH:MM timestamps with dynamic relative timestamps ("2 minutes ago" style) that update periodically without page reload. Absolute timestamps remain visible on hover.

## Motivation

- Relative timestamps provide better context for message recency
- Users can quickly understand "how long ago" without mental date math
- Common UX pattern in chat applications (Slack, Discord, etc.)

## Scope

### In Scope
- JavaScript utility function to format relative timestamps
- Periodic update of all visible timestamps (every 30 seconds)
- Keep absolute timestamp on hover via title attribute
- Support for: "just now", "X minutes ago", "X hours ago", "yesterday", "Mon", date

### Out of Scope
- Timezone configuration (uses browser timezone)
- Localization/i18n of relative time strings

## Technical Approach

1. **Template changes**: Keep `data-timestamp` attribute with Unix milliseconds, update `title` attribute to show formatted absolute time
2. **JavaScript**: Add `formatRelativeTime()` function to index.html
3. **Initial render**: Format timestamps on page load via JS (not server-side)
4. **Periodic updates**: setInterval every 30 seconds to refresh all timestamps
5. **SSE handling**: Format new messages as they arrive

## Risks

- Low: Minor JS overhead for periodic updates
- Mitigation: Batch DOM updates, skip if no visible timestamps

## Dependencies

None - pure frontend change

## Status

Draft - Pending Review
