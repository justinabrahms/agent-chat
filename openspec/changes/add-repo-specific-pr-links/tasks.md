## 1. Implementation
- [x] 1.1 Add `RepoURLs` map to Server struct for workspace-to-URL mapping
- [x] 1.2 Implement `loadRepoURLs()` function to read repo info from state.json
- [x] 1.3 Modify `linkifyIssueRefs()` to accept optional repo URL parameter
- [x] 1.4 Update `renderMarkdown()` to pass repo URL context
- [x] 1.5 Update message templates to pass workspace context to markdown function
- [x] 1.6 Handle `.git` suffix in repo URLs correctly

## 2. Testing
- [x] 2.1 Add test case for PR links with repo URL
- [x] 2.2 Add test case for PR links without repo URL (fallback)
- [x] 2.3 Add test case for handling `.git` suffix in repo URLs
- [x] 2.4 Verify all existing tests pass
