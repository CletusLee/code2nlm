## 1. Directory Tree Analytics

- [ ] 1.1 Create `LCA` (Lowest Common Ancestor) utility function that takes a slice of string paths and computes their deepest shared directory prefix.
- [ ] 1.2 Implement logic to normalize this LCA string (replace slashes with underscores) so it natively functions as a clean file prefix naming scheme.

## 2. Chunking Engine Modification

- [ ] 2.1 Add a `prefixCounts map[string]int` state tracker to `Chunker` struct or `Process` function to keep track of iterative frequencies per prefix.
- [ ] 2.2 Update `flushChunk` inside `Process` to determine the LCA dynamically from the `currentPaths` array prior to file generation.
- [ ] 2.3 Modify the `fmt.Sprintf` filename generator dynamically within `Chunker` to use the schema `[prefix]_[%03d].md` (three-digit zero-padding) rather than iterating statically on `chunkIndex`.
- [ ] 2.4 If LCA resolves to an empty root (i.e. cross-directory fallback cluster), default the naming prefix string to `"global"` or `"root"` (e.g., `global_001.md`).
- [ ] 2.5 Rename the project index generation to `000_Project_Index.md` in the `markdown` package.

## 3. Validation & Testing

- [ ] 3.1 Update the `cmd/cases_test.go` multi-file chunk tests to dynamically resolve the new `[prefix]_[number].md` filenames.
- [ ] 3.2 Update `cmd/e2e_test.go` to assert against the newly generated logical paths (e.g. `000_Project_Index.md` and `src_001.md`) instead of `01_chunk.md`.
- [ ] 3.3 Add an isolated Unit Test covering purely the multi-path string Lowest Common Ancestor logic algorithm to guarantee safety against edge-cases.
