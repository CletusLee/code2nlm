## Why

Large context windows in LLMs suffer from the "Lost in the Middle" phenomenon, and many RAG applications struggle with extracting structured insights from overly large codebases. There is an immediate need for a fast, fully local tool to chunk wide and deep directory trees evenly into manageable Markdown files, preserving structural awareness without relying on external AI interventions or API validations.

## What Changes

- Implement a purely local Go command-line tool `code2nlm` to parse code repositories.
- Add features to generate project indexes, calculate word-counts dynamically, and split files using Tree-sitter AST or word-count boundaries.
- Support strict validation of capacity constraints against output tiers (`--max-sources`).
- Introduce zero-AI processing, ensuring structural mapping with injected contextual headers for every chunked markdown file.

## Capabilities

### New Capabilities
- `cli-core`: The basic structural foundation of the tool, utilizing Cobra for argument parsing and afero for resilient memory/disk filesystem handling.
- `scanner-engine`: High concurrency directory scanning, `.gitignore` parsing, and overall workspace virtual tree generation along with boundary checks.
- `chunking-engine`: File chunking strategy relying on structural AST parsing via go-tree-sitter to perform clean syntactic splits or fallback word counts.
- `markdown-generator`: Output formatting logic, injecting contextual headers into individual chunks and creating the central `00_Project_Index.md`.

### Modified Capabilities
*(None)*

## Impact

- Creates a new, self-contained Go CLI application.
- Introduces robust requirements for unit tests (Table-Driven and In-Memory Filesystem tests) for scanner, granularity computations, and AST implementations.
- No existing systems or dependencies are altered since this represents a fresh code generation path.
