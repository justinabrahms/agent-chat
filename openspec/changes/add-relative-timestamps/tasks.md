# Tasks: Add Relative Timestamps

## Implementation Checklist

- [x] Add `formatRelativeTime(timestamp)` JavaScript function to index.html
- [x] Add `updateAllTimestamps()` function to batch-update all visible timestamps
- [x] Initialize timestamps on DOMContentLoaded
- [x] Add setInterval for periodic updates (every 30 seconds)
- [x] Handle new SSE messages - format timestamps as they arrive
- [x] Update `title` attribute to show full formatted date/time for hover

## Testing

- [ ] Verify "just now" shows for messages < 1 minute old
- [ ] Verify "X minutes ago" shows for messages 1-59 minutes old
- [ ] Verify "X hours ago" shows for messages 1-23 hours old
- [ ] Verify day names show for messages 1-6 days old
- [ ] Verify dates show for messages > 6 days old
- [ ] Verify timestamps update without page reload
- [ ] Verify hover shows absolute timestamp
