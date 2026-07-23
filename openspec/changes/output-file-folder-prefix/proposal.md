## Why

Currently, when `code2nlm` is run against a directory (e.g., `--input ./my-project`), output files are named based on the internal Lowest Common Ancestor (LCA) path of the source files, but carry no information about which *root input folder* they originated from. When users process multiple projects or run the tool against named sub-directories, all outputs land in the same flat folder with ambiguous names like `src_001.md` or `global_001.md`, making it impossible to tell at a glance which source project or folder each file belongs to.

## What Changes

- The basename of the `--input` directory (the "folder name") will be prepended as a prefix to **every** generated output file, including the index file.
- The index file will be renamed from `000_Project_Index.md` to `<folder-name>_000_Project_Index.md`.
- Chunk files will be renamed from `<lca>_001.md` to `<folder-name>_<lca>_001.md`.
- All internal cross-references (e.g., the `Global Context` header line referencing `000_Project_Index.md`) will be updated to use the new prefixed filename.

## Capabilities

### New Capabilities
- `output-file-folder-prefix`: Prefix every generated output file (index and chunks) with the base name of the `--input` folder, making output files unambiguously attributable to their source directory.

### Modified Capabilities
- `chunking-engine`: Chunk filename format changes — the folder-name prefix is prepended to the LCA-based filename segment (e.g., `src_001.md` → `my-project_src_001.md`).
- `markdown-generator`: Index file naming changes — the folder-name prefix is prepended (e.g., `000_Project_Index.md` → `my-project_000_Project_Index.md`). The contextual header's cross-reference link must also reflect the new index filename.

## Impact

- `chunking/chunker.go` — filename generation logic in `flushChunk`.
- `markdown/generator.go` — index filename constant and the `Global Context` header string in `FormatContextualHeader`.
- `cmd/root.go` — `ProjectName` (already derived from `filepath.Base(InputPath)`) can serve as the folder prefix; the same value must be passed through to `GenerateIndex` so it can construct the prefixed index filename.
- All existing tests that assert on specific output file names (`000_Project_Index.md`, bare LCA names) must be updated to expect the new prefix format.
