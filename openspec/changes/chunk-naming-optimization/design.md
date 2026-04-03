## Context

Currently, the `code2nlm` CLI chunker groups processed source files into generically named output files (e.g., `01_chunk.md`, `02_chunk.md`). While efficient at packing tokens up to the strict `--max-words` limit, it loses vital spatial context at the file system level for both humans and Retrieval-Augmented Generation (RAG) models. Users want output chunks to be thoughtfully named after the overarching directories they contain (e.g., `src_components_001.md`), without sacrificing the token-packing mechanism that keeps total file counts low.

## Goals / Non-Goals

**Goals:**
- Dynamically generate output filenames based on the Lowest Common Ancestor (LCA) directory of the files within a chunk.
- Format the chosen LCA path with underscores `_` to denote directory depth (e.g., `src_utils_formatting`).
- Append sequential numeric suffixes per path-prefix to guarantee uniqueness. Use three-digit zero-padding (e.g., `src_001.md`, `src_002.md`) to prevent lexicographical sorting issues where `src_10.md` might appear before `src_2.md`.
- Force the global project index file to be named `000_Project_Index.md`, leveraging numeric priority to ensure it always stays at the absolute top of the LLM Context Window during RAG retrieval.
- Ensure no impact to the current word-bound chunking logic (token packing density).

**Non-Goals:**
- We are *not* changing the word-counting algorithms or Tree-sitter split integrations.
- We are *not* changing how files are topologically sorted; `afero.Walk`'s natural alphabetical order is already optimal.

## Decisions

- **Lowest Common Ancestor Logic**: We will implement a lightweight utility that takes the slice `currentPaths` immediately prior to `flushChunk` and computes their common directory ancestor.
- **Prefix Collision Tracking**: The `Chunker` struct will utilize a state map `map[string]int` to track how many times a given prefix (e.g., `src_utils`) has been used. During string formatting, a `%03d` specifier MUST be used to ensure fixed-length suffixes like `_001.md`.
- **Root Fallback**: If files from completely disjoint roots are bunched into one chunk (e.g., `cmd/main.go` and `scanner/engine.go`), their LCA reduces to the root. In such cases, the fallback prefix shall be `root_001.md` or `global_001.md`.

## Risks / Trade-offs

- **Risk**: Test fixtures depending heavily on the legacy `01_chunk.md` format will break.
  - **Mitigation**: Update assertions in `e2e_test.go` and `cases_test.go` to validate against the new logical dynamic names.
