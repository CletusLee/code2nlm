## MODIFIED Requirements

### Requirement: Index Initialization
The system SHALL automatically generate `<folder-name>_000_Project_Index.md` (where `<folder-name>` is the base name of the `--input` directory) mapping the directory tree cleanly before populating chunks. The `projectName` parameter MUST be accepted by `GenerateIndex` and used to construct the prefixed filename.

#### Scenario: Standard Project Execution
- **WHEN** generating output with `--input ./my-project`
- **THEN** the initial file emitted is `my-project_000_Project_Index.md`

### Requirement: Context Headers
The system SHALL enforce exact metadata contextual headers inside every created chunk, stating the project module mapping clearly to assist RAG processing. The `Global Context` line in the header MUST reference the correct prefixed index filename (e.g., `` `my-project_000_Project_Index.md` ``), not the bare `000_Project_Index.md`. The `FormatContextualHeader` function SHALL accept the actual index filename as a parameter rather than hard-coding it.

#### Scenario: Markdown Chunk Initialization
- **WHEN** the engine opens a new dynamically named chunk mapping file (e.g. `my-project_src_components_001.md`)
- **THEN** the absolute first lines generated must include `# Module: [Lowest Common Ancestor]` and the `**Global Context**` line MUST reference `` `my-project_000_Project_Index.md` `` (the prefixed index file).
