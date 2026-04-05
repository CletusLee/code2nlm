## ADDED Requirements

### Requirement: Index Initialization
The system SHALL automatically generate `000_Project_Index.md` mapping the directory tree cleanly before populating chunks.

#### Scenario: Standard Project Execution
- **WHEN** generating output
- **THEN** the initial file emitted is always `000_Project_Index.md`

### Requirement: Context Headers
The system SHALL enforce exact metadata contextual headers inside every created chunk, stating the project module mapping clearly to assist RAG processing.

#### Scenario: Markdown Chunk Initialization
- **WHEN** the engine opens a new dynamically named chunk mapping file (e.g. `src_components_001.md`)
- **THEN** the absolute first lines generated must include `# Module: [Lowest Common Ancestor]` and correctly reflect its logical context.
