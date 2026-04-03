## Why

Currently, `code2nlm` generates chunk files using a simple numeric sequence (e.g., `01_chunk.md`, `02_chunk.md`). While this prevents naming collisions, it makes it difficult for LLMs (and humans) to understand the general context of the chunk from its filename alone. Furthermore, we need a smarter strategy to name chunks based on the directories of the files they contain (e.g., `src_components_1.md`), while continuing to pack files tightly to minimize the total number of chunks.

## What Changes

- Modify the chunk generation logic to dynamically determine a descriptive filename based on the directory structure of the files contained within the chunk.
- Consolidate paths by finding the lowest common directory among the files in a chunk to serve as the chunk's prefix.
- Delineate multistory folder paths with underscores `_` (e.g., `src_utils`).
- Append enumerations if multiple chunks share the same folder context (e.g., `src_utils_1.md`, `src_utils_2.md`) to resolve collisions.
- Preserve the existing tight-packing logic (up to `--max-words`) to ensure minimized total output files.

## Capabilities

### New Capabilities
None.

### Modified Capabilities

- `chunking-engine`: Update chunking engine to generate contextual file names derived from directory hierarchies rather than sequential integers, while preserving token-packing efficiency.
- `markdown-generator`: Modify or verify that index generation maps the newly contextualized chunk file names.

## Impact

- `chunking/chunker.go`: The `flushChunk` local function will require logic to compute the lowest common ancestor directory path from `currentPaths`.
- Output file names will change completely, impacting any consumer relying on the `00_chunk.md` format.
- `cmd/cases_test.go` and `cmd/e2e_test.go`: File name assertions will need to be updated.
