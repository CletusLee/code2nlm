## ADDED Requirements

### Requirement: Chunk Stream Partitioning
The system SHALL group scanned files and generate chunk boundaries dynamically once a given file cluster reaches `--max-words`.

#### Scenario: Multi-File Split Boundary
- **WHEN** current stream surpasses `--max-words`
- **THEN** it safely concludes the active `0N_xxx.md` file and creates the next sequentially valid markdown file.

### Requirement: Contextual AST Parsing
The system SHALL parse large single files using `smacker/go-tree-sitter` bindings to accurately split code precisely after function boundaries rather than completely haphazardly.

#### Scenario: Class Method Split
- **WHEN** a single `150,000` word file triggers split limits
- **THEN** AST correctly matches nearest ending bracket `}` boundary and gracefully separates content
