# Change: Add message grouping by sender

## Why
Chat interfaces with many consecutive messages from the same agent become visually cluttered when each message repeats the sender name and timestamp. Slack, Discord, and similar tools group consecutive messages from the same sender to reduce noise and improve readability.

## What Changes
- Consecutive messages from the same sender collapse headers (show sender name only on first message in a group)
- Visual grouping with subtle separator between different senders
- Timestamps shown on first message of a group, with relative time for subsequent messages on hover

## Impact
- Affected specs: `chatroom-ui`
- Affected code: Message rendering templates, CSS styles
- Dependencies: Builds on `add-chatroom-ui` change
