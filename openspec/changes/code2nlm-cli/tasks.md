## 1. Project Initialization & Architecture

- [x] 1.1 Scaffold Cobra CLI structure and setup `root.go` with strict configuration flags (`--input`, `--output`, `--max-sources`, `--max-words`, `--ignore-file`, `--strategy`).
- [x] 1.2 Implement the `spf13/afero` interface dependency injection for all filesystem IO operations to ensure testability.
- [x] 1.3 Add standard Go testing table-tests layout to verify argument parsing validity correctly.

## 2. Scanner & Filtering Infrastructure

- [x] 2.1 Integrate `sabhiram/go-gitignore` and logic to locate/parse standard `.gitignore` or specific `--ignore-file` patterns.
- [x] 2.2 Implement a highly efficient sequential directory traversal algorithm using `filepath.WalkDir`. Since we only read metadata and avoid opening files, single-threaded scanning is extremely fast; omitting `sync.WaitGroup` prevents unnecessary context switch overhead.
- [x] 2.3 Develop the Virtual Tree tabulation logic utilizing a "Double-Pass I/O" avoidance strategy. Directly read `FileInfo.Size()` (Bytes) during the walk to compute the total size without opening the file content.

## 3. Strict Capacity Validation

- [x] 3.1 Implement the dynamic file counting formula `ceil(Estimated Words / --max-words)` by applying an estimation formula (e.g., Estimated Words = Bytes / 5) to the pre-calculated sizes.
- [x] 3.2 Add strict abort condition: emit clean `stderr` contextual validation error and non-zero exit if required files exceed `--max-sources`.
- [x] 3.3 Create corresponding `TestGranularity` testing scenarios evaluating edge bounds and strict abort triggers via mock FS.

## 4. Markdown Generation & Indexing

- [x] 4.1 Create functionality to synthesize the overall structured `00_Project_Index.md` generated out of the valid Virtual Tree.
- [x] 4.2 Build markdown generation utilities specifically injecting exact `# Module: [Domain Name]` contextual metadata headers formatting at the top of every chunk.
- [x] 4.3 Implement Path Normalization logic: ensure all paths written to Markdown (including Index and Contextual Headers) use Go's `filepath.ToSlash()`. This prevents Windows backslashes (`\`) from causing directory structure hallucinations in LLMs.

## 5. File Chunking Engine & AST Support

- [x] 5.1 Implement the sequential processing loop streaming code strings into `0N_xxx.md` output blobs wrapping safely when nearing file word thresholds.
- [x] 5.2 Integrate `smacker/go-tree-sitter` bindings and logic to successfully detect nearest structural node closures (`}`) recursively when chunking large codebase targets.
- [x] 5.3 Write functional logic for the fallback mechanism, allowing word-splitting directly if the parser encounters unstructured text beyond parsing comprehension.

## 6. End-to-End Orchestration & E2E Validation

- [x] 6.1 Connect the CLI Cobra handler to Scanner, Validator, and Chunk Engines effectively driving workflows within the `RunE` setup.
- [x] 6.2 Write definitive `TestE2E_CLI` end-to-end tests validating index mapping consistency, distinct contextual header placements, and multi-file generated chunks precisely dynamically against a fixture mock directory.

## 7. CI/CD & Cross-Platform Compilation

- [x] 7.1 Setup GitHub Actions workflow to automate the build, test, and release pipelines.
- [x] 7.2 Configure cross-platform CGO compilation toolchains (using `zig cc` or `xgo`) within the pipeline to build fully static executable binaries for Windows, macOS, and Linux targets to resolve `go-tree-sitter` dependencies automatically.
