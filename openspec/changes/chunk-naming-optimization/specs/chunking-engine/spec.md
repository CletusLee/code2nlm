## MODIFIED Requirements

### Requirement: Chunk Stream Partitioning
The system SHALL group scanned files and generate chunk boundaries dynamically once a given file cluster reaches `--max-words`. Crucially, the filename of the output chunk SHALL NOT be a simple sequential integer (like `01_chunk.md`), but MUST reflect the lowest common ancestor directory path of all internal files, normalized with underscores and appended with a **zero-padded** sequence number for collisions (e.g. `src_utils_math_001.md`).

#### Scenario: Single Directory Chunk Naming
- **WHEN** all files within a chunk share the `src/components/button` directory path
- **THEN** the chunk safely concludes upon hitting `--max-words` and is explicitly named `src_components_button_001.md`

#### Scenario: Same Directory Collision Name Iteration
- **WHEN** a second chunk is generated from files exclusively within `src/components/button`
- **THEN** the system increments the tracker to emit `src_components_button_002.md` (formatted securely via `%03d`)

#### Scenario: Distinct Directory Fallback
- **WHEN** files from completely disjoint paths (e.g., `cmd/main.go` and `scanner/engine.go`) are packed into a single final leftover chunk
- **THEN** their lowest common ancestor mathematically reduces to the root directory, yielding names like `root_001.md` or `global_001.md`