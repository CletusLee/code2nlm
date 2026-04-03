## Context

Processing complex, large-scale code repositories into a Retrieval-Augmented Generation (RAG) system, such as NotebookLM, often risks context fragmentation, causing a "Lost in the Middle" phenomenon. `code2nlm` addresses this by evenly partitioning code off-line via deterministic, structural mechanisms into standard `.md` files without requiring live API interactions.

## Goals / Non-Goals

**Goals:**
- Build a fast and reliable entirely local Go CLI tool (`code2nlm`).
- Parse large directory trees efficiently using highly concurrent strategies.
- Maintain syntactic fidelity on file splits using Tree-sitter's AST parsing functionality.
- Provide strict fail-fast validation regarding chunk maximums (`--max-sources`, `--max-words`).
- Output contextually-rich Markdown chunks detailing file structures and included paths.

**Non-Goals:**
- Automatic uploads to LLMs or NotebookLM deployments (the tool relies 100% on local disk processing).
- Deep linting or codebase semantic restructuring beyond plain chunking boundaries.

## Decisions

- **CLI Foundation:** `spf13/cobra` provides standardized argument parsing, while `spf13/afero` replaces direct `os` system calls. This makes unit testing incredibly fast using an in-memory mapped filesystem rather than disk reads.
- **Parse & Ignore Rules:** Standard `sabhiram/go-gitignore` perfectly models Git's blocklist behavior, while standard `filepath.WalkDir` allows for swift filesystem traversal to calculate upfront totals.
- **AST Chunking:** `smacker/go-tree-sitter` offers CGO interoperability to robustly find class, struct, and function boundaries `}` to enact a clean split just after the bracket. For unsupported extensions, it gracefully degrades to word-count boundaries.
- **Predictive Strict Validation:** Before executing the disk-intensive split, total words are estimated to determine if `ceil(Total Words / max-words)` breaches `--max-sources`. If it does, the operation FATALs out immediately to avoid incomplete states. To prevent a "Double-Pass I/O" performance penalty where files must be fully read just to determine limits, we use a fast evaluation strategy during `filepath.WalkDir`. By reading `FileInfo.Size()` to compute Bytes, we roughly estimate word counts (e.g., 1 Word ≈ 5 Bytes), completing the validation in milliseconds without opening file handles.
- **Path Normalization:** To prevent AI confusion caused by Windows backslashes (`\`) when interpreting directory trees, the system enforces strict POSIX-style paths. All file paths injected into the Markdown outputs (including `00_Project_Index.md` and Contextual Headers in each chunk) are normalized using Go's `filepath.ToSlash()`.

## Risks / Trade-offs

- **Risk:** CGO Dependency overhead tied to `go-tree-sitter` can hamper static compilation on specific environments. 
  **Mitigation:** We will manage this explicitly via CI/CD Pipelines (e.g., GitHub Actions) paired with cross-platform CGO compilation toolchains like `zig cc` or `xgo`. This ensures we distribute completely static binaries for Windows, macOS, and Linux, saving end-users from the burden of building C/C++ environments themselves.
- **Risk:** Heavily minimized/obfuscated files devoid of clean structural boundaries.
  **Mitigation:** The system will gracefully fall back to a safe text-split directly once the word count watermark is fully reached and no boundary is identified.
