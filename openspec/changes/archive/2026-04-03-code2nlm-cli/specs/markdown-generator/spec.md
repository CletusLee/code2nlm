## ADDED Requirements

### Requirement: Index Initialization
The system SHALL automatically generate `00_Project_Index.md` mapping the directory tree cleanly before populating chunks.

#### Scenario: Standard Project Execution
- **WHEN** generating output
- **THEN** the initial file emitted is always `00_Project_Index.md`

### Requirement: Context Headers
The system SHALL enforce exact metadata contextual headers inside every created chunk, stating the project module mapping clearly to assist RAG processing.

#### Scenario: Markdown Chunk Initialization
- **WHEN** the engine opens a new `01_xxx.md` mapping file
- **THEN** the absolute first lines generated must include `# Module: [Domain...` and contextual references prior to code lines.
