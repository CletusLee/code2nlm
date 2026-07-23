## MODIFIED Requirements

### Requirement: Chunk Stream Partitioning
The system SHALL group scanned files and generate chunk boundaries dynamically once a given file cluster reaches `--max-words`. The filename of the output chunk MUST be prefixed with the base name of the `--input` folder (the "folder name"), followed by an underscore, then the lowest common ancestor directory path of all internal files normalized with underscores, and finally a **zero-padded** sequence number for collisions. Format: `<folder-name>_<lca>_NNN.md` (e.g. `my-project_src_utils_math_001.md`).

#### Scenario: Multi-File Split Boundary
- **WHEN** current stream surpasses `--max-words`
- **THEN** it safely concludes the active chunk file and creates the next sequentially valid markdown file, both carrying the folder-name prefix.

#### Scenario: Single Directory Chunk Naming
- **WHEN** all files within a chunk share the `src/components/button` directory path and the input folder is `my-project`
- **THEN** the chunk is explicitly named `my-project_src_components_button_001.md`

#### Scenario: Same Directory Collision Name Iteration
- **WHEN** a second chunk is generated from files exclusively within `src/components/button`
- **THEN** the system increments the tracker to emit `my-project_src_components_button_002.md` (formatted via `%03d`)

#### Scenario: Distinct Directory Fallback
- **WHEN** files from completely disjoint paths are packed into a single final leftover chunk and the input folder is `my-project`
- **THEN** their lowest common ancestor reduces to root/global, yielding `my-project_global_001.md`

### Requirement: Contextual AST Parsing
The system SHALL parse large single files using `smacker/go-tree-sitter` bindings to accurately split code precisely after function boundaries rather than completely haphazardly.

#### Scenario: Class Method Split
- **WHEN** a single `150,000` word file triggers split limits
- **THEN** AST correctly matches nearest ending bracket `}` boundary and gracefully separates content
