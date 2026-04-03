## ADDED Requirements

### Requirement: File System Walk Traversal
The system SHALL concurrently gather filesystem data to calculate the exact word count size before starting file chunk generation.

#### Scenario: Scan Calculation Match
- **WHEN** system evaluates `filepath.WalkDir` over standard directories
- **THEN** it sums and caches word counts of matching files locally.

### Requirement: Strict Ignore Filtering
The system SHALL parse `.gitignore` (or specified override file via `--ignore-file`) and apply matching patterns to block files or directories correctly from the scanner queue.

#### Scenario: Exclusion Hit
- **WHEN** `node_modules` or `*.exe` exists in blocklist
- **THEN** scanner correctly ignores matching paths natively using go-gitignore matches.
