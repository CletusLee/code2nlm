## Context

`code2nlm` transforms a codebase directory into a set of Markdown chunk files suitable for LLM ingestion (e.g., NotebookLM). Currently the output filenames are derived purely from the internal Lowest Common Ancestor (LCA) path of the source files within the scanned directory, with no reference to the root input folder name. When users aggregate outputs from multiple projects into one destination, or compare outputs across runs with different `--input` paths, the filenames give no indication of their origin.

The `ProjectName` value (derived as `filepath.Base(filepath.Clean(InputPath))`) already exists in `cmd/root.go` and is passed to the `Chunker`. It just isn't used as part of file naming yet.

## Goals / Non-Goals

**Goals:**
- Prepend the input folder's base name to every output file (index + all chunks).
- Keep cross-reference strings inside generated Markdown consistent with new names.
- Require zero new flags or config — the folder name is derived automatically from `--input`.

**Non-Goals:**
- Sanitizing or slugifying the folder name (we preserve it as-is; users control the directory name).
- Supporting a user-supplied custom prefix (out of scope for this change).
- Changing the `000_` numeric ordering of the index file (the prefix sits before it).

## Decisions

### Decision 1: Where to derive the prefix

**Choice**: Use the existing `ProjectName` field (`filepath.Base(filepath.Clean(InputPath))`) already computed in `cmd/root.go` and passed into `Chunker`.

**Rationale**: The value is already computed and semantically correct. No new state is introduced. The same value is used for the in-document `**Project**:` metadata line.

**Alternative considered**: Adding a dedicated `--prefix` flag. Rejected — the folder name is the natural, zero-config source of truth for the prefix.

### Decision 2: Prefix format for output files

**Choice**: `<folderName>_<rest>` — a single underscore separating the folder prefix from the LCA segment or the index's numeric prefix.

| Before | After (`--input ./my-project`) |
|---|---|
| `000_Project_Index.md` | `my-project_000_Project_Index.md` |
| `src_components_001.md` | `my-project_src_components_001.md` |
| `global_001.md` | `my-project_global_001.md` |

**Rationale**: Consistent with the existing underscore-delimited naming convention already used for LCA path segments.

### Decision 3: Propagating `projectName` to `GenerateIndex`

**Choice**: Add a `projectName string` parameter to `markdown.GenerateIndex(...)` so it can construct the prefixed index filename.

**Rationale**: `GenerateIndex` currently has no knowledge of the project name. Adding it as an explicit parameter keeps the function pure and testable, without introducing global state.

**Alternative considered**: Embedding the project name in the `Chunker` and having `Chunker.Process` call a new index-aware method. Rejected — `GenerateIndex` is called independently in `cmd/root.go` before the chunker; coupling them would complicate flow.

### Decision 4: Updating the cross-reference string in `FormatContextualHeader`

**Choice**: Change the hard-coded `000_Project_Index.md` reference in `FormatContextualHeader` to accept `indexFileName string` as a parameter, so callers pass the computed prefixed name.

**Rationale**: The contextual header that `Chunker` injects into every chunk currently hard-codes the index filename. If we only rename the file without updating the header string, the reference breaks. Passing the actual index filename keeps the contract intact.

## Risks / Trade-offs

- **Existing test fixtures use old filenames** → All test assertions on output filenames (`000_Project_Index.md`, bare LCA names) must be updated. Mitigation: tests are clearly scoped and easy to locate.
- **Folder names with special characters** (spaces, Unicode) could produce unusual filenames → Acceptable for now; `Non-Goals` explicitly excludes sanitization. Users are responsible for folder naming.

## Migration Plan

This is a purely additive naming change. No existing stored files are migrated; each run always produces fresh output. No rollback needed.
