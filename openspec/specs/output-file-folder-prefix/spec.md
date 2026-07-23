# output-file-folder-prefix

## ADDED Requirements

### Requirement: Input Folder Name Prefix on All Output Files
Every file written to the output directory SHALL be prefixed with the base name of the `--input` directory (the "folder name"), followed by an underscore, before any other name segment. This applies to the project index file and every chunk file.

#### Scenario: Index file gets folder prefix
- **WHEN** the tool is run with `--input ./my-project`
- **THEN** the generated index file is named `my-project_000_Project_Index.md` (not `000_Project_Index.md`)

#### Scenario: Chunk file gets folder prefix
- **WHEN** the tool is run with `--input ./my-project` and a chunk's LCA resolves to `src_components`
- **THEN** the chunk file is named `my-project_src_components_001.md` (not `src_components_001.md`)

#### Scenario: Global/root chunk gets folder prefix
- **WHEN** files from disjoint paths produce a `global` LCA
- **THEN** the chunk file is named `<folder-name>_global_001.md`

#### Scenario: Folder name derived from Linux absolute input path
- **WHEN** the user supplies `--input /test/product` (Linux absolute path)
- **THEN** the prefix used is `product` (the base name only, no leading slash or parent path components)

#### Scenario: Folder name derived from Windows absolute input path
- **WHEN** the user supplies `--input D:/test/product` or `--input D:\test\product` (Windows absolute path)
- **THEN** the prefix used is `product` (the base name only, no drive letter or parent path components)

#### Scenario: Trailing path separator is ignored
- **WHEN** the user supplies `--input /test/product/` or `--input D:/test/product/` (trailing separator)
- **THEN** the prefix used is still `product` (the trailing separator is stripped before extracting the base name)
