## 1. markdown/generator.go — Add projectName to GenerateIndex & FormatContextualHeader

- [x] 1.1 Add `projectName string` parameter to `GenerateIndex`; compute prefixed index filename as `<projectName>_000_Project_Index.md` and use it when writing the index file.
- [x] 1.2 Add `indexFileName string` parameter to `FormatContextualHeader`; replace the hard-coded `` `000_Project_Index.md` `` reference in the `Global Context` line with the passed-in filename.

## 2. chunking/chunker.go — Prepend folder-name prefix to chunk filenames

- [x] 2.1 Add a `FolderPrefix string` field to the `Chunker` struct (or reuse `ProjectName`), representing the base name of the input directory.
- [x] 2.2 In `flushChunk`, change the filename format from `<prefix>_NNN.md` to `<folderPrefix>_<prefix>_NNN.md`.
- [x] 2.3 Update the call to `markdown.FormatContextualHeader` inside `flushChunk` to pass the computed prefixed index filename (e.g., `<folderPrefix>_000_Project_Index.md`).

## 3. cmd/root.go — Wire projectName into GenerateIndex

- [x] 3.1 Update the call to `markdown.GenerateIndex(FS, OutputPath, virtualTree)` to pass `ProjectName` as the new `projectName` argument.

## 4. Update tests

- [x] 4.1 Update `markdown/generator.go` tests (if any) to assert the output file is named `<projectName>_000_Project_Index.md`.
- [x] 4.2 Update `chunking/chunker.go` / `cmd/root.go` tests that assert on output chunk filenames to expect the new `<folderName>_<lca>_NNN.md` format.
- [x] 4.3 Update any test helper or E2E test that checks `FormatContextualHeader` output to verify the `Global Context` line references the prefixed index filename.

## 5. Verify

- [x] 5.1 Build the binary (`go build ./...`) and confirm no compile errors.
- [x] 5.2 Run all tests (`go test ./...`) and confirm all pass.
- [x] 5.3 Manually run the tool against a sample directory and confirm output files are named `<folder>_000_Project_Index.md`, `<folder>_<lca>_001.md`, etc., and that the contextual header inside each chunk correctly references the prefixed index file.
