# code2nlm

> **Transform entire code repositories into structured Markdown notebooks — optimized for NotebookLM, RAG pipelines, and any LLM that needs to reason over a full codebase.**

`code2nlm` is a fast, local-only CLI tool that scans a source directory, intelligently chunks the code into token-bounded Markdown files, and produces a `000_Project_Index.md` as the master manifest. The result is a ready-to-upload notebook set that lets tools like Google NotebookLM "read" a project as if it were a collection of documents.

---

## Table of Contents

- [Why code2nlm?](#why-code2nlm)
- [How It Works](#how-it-works)
- [Installation](#installation)
  - [Download a Pre-built Binary](#download-a-pre-built-binary)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
  - [Basic Example](#basic-example)
  - [All Flags](#all-flags)
  - [Common Recipes](#common-recipes)
- [Using the Output with NotebookLM](#using-the-output-with-notebooklm)
- [Using the Output with Other LLMs / RAG Pipelines](#using-the-output-with-other-llms--rag-pipelines)
- [Output File Structure](#output-file-structure)
- [Advanced: Custom Ignore Rules](#advanced-custom-ignore-rules)
- [Building from Source](#building-from-source)
- [Contributing](#contributing)
- [License](#license)

---

## Why code2nlm?

Large language models like Gemini or GPT-4 have finite context windows. You cannot paste an entire codebase into a single prompt. Existing workarounds (e.g., manual copy-paste, custom scripts) are brittle and produce noisy output.

`code2nlm` solves this by:

| Problem | How code2nlm Solves It |
|---|---|
| Context window limits | Splits code into configurable token-bounded chunks |
| Noisy source maps / Base64 blobs | Automatically strips `*.map` URIs and Base64 Data URIs |
| Hidden files and build artifacts | Respects `.gitignore` (or any custom ignore file) |
| RAG retrieval ambiguity | Names each chunk after its dominant directory (e.g., `src_utils_001.md`) |
| Lost navigational context | `000_Project_Index.md` lists every source file with its chunk location |
| Cross-platform path confusion | Normalizes all paths to forward-slash format |

---

## How It Works

```
Your Repo/                         nlm_output/
├── src/                  ──►      ├── 000_Project_Index.md   ← master file map
│   ├── utils/                     ├── src_utils_001.md       ← chunked by directory
│   ├── components/                ├── src_components_001.md
│   └── ...                        ├── src_components_002.md
├── tests/                         ├── tests_001.md
└── ...                            └── global_001.md          ← cross-dir files
```

1. **Scan** — Walks the input directory respecting `.gitignore` rules; collects file paths and byte sizes. Skips hidden files and binary/media extensions automatically.
2. **Index** — Writes `000_Project_Index.md` listing every source file found.
3. **Chunk** — Groups files into Markdown chunks within the `--max-words` budget. Each chunk is named after the **Lowest Common Ancestor (LCA)** directory of the files it contains, then suffixed with a zero-padded counter (`_001`, `_002` …) for correct lexicographical ordering.
4. **Denoise** — Before writing, strips Base64 Data URIs and Source Map inline blobs to reclaim token budget.

---

## Installation

### Download a Pre-built Binary

Go to the [**Releases page**](https://github.com/YOUR_ORG/code2nlm/releases) and download the binary for your platform:

| Platform | File |
|---|---|
| Windows (64-bit) | `code2nlm-windows-amd64.exe` |
| Linux (64-bit) | `code2nlm-linux-amd64` |
| Linux (ARM64) | `code2nlm-linux-arm64` |

**Windows:**
```powershell
# Rename after download for convenience
Rename-Item code2nlm-windows-amd64.exe code2nlm.exe
# Run directly (no install required)
.\code2nlm.exe --help
```

**Linux:**
```bash
chmod +x code2nlm-linux-amd64
sudo mv code2nlm-linux-amd64 /usr/local/bin/code2nlm
code2nlm --help
```

### Build from Source

Requirements: **Go 1.20+**

```bash
git clone https://github.com/YOUR_ORG/code2nlm.git
cd code2nlm
go build -o code2nlm .
```

> **Note on CGO / Tree-sitter**: The tool includes optional AST-aware chunking for Go source files. This feature requires CGO. If `CGO_ENABLED=0`, it gracefully falls back to line-boundary splitting for all file types.

---

## Usage

### Basic Example

```bash
# Process the current directory, output to ./nlm_output
code2nlm -i ./

# Process a specific repository
code2nlm -i ../my-project

# Specify a custom output directory
code2nlm -i ../my-project -o ./notebooklm-docs
```

### All Flags

```
Usage:
  code2nlm [flags]

Flags:
  -i, --input string         Input directory path (default "./")
  -o, --output string        Output directory (default "./nlm_output")
  -m, --max-sources int      Max number of output .md files (default 50)
  -w, --max-words int        Max word count per output file (default 100000)
      --ignore-file string   Path to ignore rules file (default ".gitignore")
  -s, --strategy string      Chunking strategy: "ast" or "dir" (default "ast")
  -h, --help                 help for code2nlm
```

### Common Recipes

```bash
# Large monorepo — increase file limit and word budget
code2nlm -i ./monorepo -m 300 -w 150000

# Small library — tighter chunks for better RAG retrieval precision  
code2nlm -i ./my-lib -m 20 -w 50000

# Use a custom ignore file instead of .gitignore
code2nlm -i ./project --ignore-file .nlmignore

# Target only the src/ subdirectory
code2nlm -i ./project/src -o ./project/nlm_out
```

> **Too many files error?**  
> If you see `FATAL: Project is too large`, increase `-m` (e.g., `-m 200`) to allow more output files.

---

## Using the Output with NotebookLM

[Google NotebookLM](https://notebooklm.google.com/) lets you upload documents and then ask questions about them. The chunked Markdown files produced by `code2nlm` are perfect for this workflow.

### Step-by-Step

**Step 1: Generate the Markdown files**

```bash
code2nlm -i ./my-project -o ./nlm_output -m 100
```

This creates a folder like:
```
nlm_output/
├── 000_Project_Index.md
├── src_001.md
├── src_components_001.md
├── src_utils_001.md
└── ...
```

**Step 2: Open NotebookLM**

Navigate to [notebooklm.google.com](https://notebooklm.google.com) and create a **New Notebook**.

**Step 3: Upload the Markdown Files**

Click **"+ Add Source"** → **"Upload file"**. 

Select **all** `.md` files from your `nlm_output/` folder.

> **Tip:** Always include `000_Project_Index.md` first — it contains the full file map and helps NotebookLM orient itself within the codebase.

**Step 4: Ask Questions**

Once uploaded, you can ask NotebookLM questions such as:

- *"Explain the overall architecture of this project."*
- *"How does the authentication flow work? Which files are involved?"*
- *"What does the `QueryEngine` class do and where is it defined?"*
- *"Find all places where database connections are opened."*
- *"Summarize the purpose of the `src/utils/` directory."*
- *"How does error handling work in the API layer?"*

### Tips for Best Results

| Tip | Why It Helps |
|---|---|
| Upload **all** chunks including the Index | The Index helps the model navigate, just like a table of contents |
| Use smaller `--max-words` (e.g. 50000) for precision | Smaller chunks = better retrieval for specific questions |
| Use larger `--max-words` for broad overviews | Larger chunks = more holistic context for architecture questions |
| Ask about **specific files or directories** | NotebookLM can cite its sources — mention file paths for targeted answers |
| Ask for **diagrams or summaries** | NotebookLM can generate audio overviews and study guides from code |

---

## Using the Output with Other LLMs / RAG Pipelines

The output format is generic Markdown and works with any LLM or RAG stack.

### Direct LLM Prompt Injection

For small projects, you can concatenate all chunks directly into a prompt:

```bash
cat nlm_output/*.md | clip  # Copies to clipboard (Windows CLI)
# or ... | pbcopy on some Linux environments
```

### RAG Pipeline Integration

Each chunk file has a **structured header** that makes it ideal for vector embedding:

```markdown
# Module: src_utils
**Project**: my-project
**Global Context**: Please refer to `000_Project_Index.md` for the complete directory structure.

## Included Paths in this Chunk
* `src/utils/auth.ts`
* `src/utils/crypto.ts`

---
### File: `src/utils/auth.ts`
...
```

- **`# Module:`** → Use as document title for metadata
- **`## Included Paths`** → Use as keyword tags for filtering  
- **`000_Project_Index.md`** → Load first as the master context document

### With Claude / GPT / Gemini Chat

1. Upload the `000_Project_Index.md` to give the model the full file map
2. Paste or upload specific chunk files when asking about that module
3. Reference the index when asking cross-cutting questions

---

## Output File Structure

### `000_Project_Index.md` — The Master Index

Always generated first. Contains a flat list of every source file discovered:

```markdown
# Project Index

## Directory Structure

- `src/main.ts`
- `src/utils/auth.ts`
- `src/components/App.tsx`
...
```

This file is pinned at the top lexicographically (`000_` prefix) so RAG systems always retrieve it with highest priority.

### Chunk Files — e.g., `src_utils_001.md`

Named after the **Lowest Common Ancestor** directory of contained files:

| Chunk Name | What It Contains |
|---|---|
| `global_001.md` | Files from diverse/root-level directories |
| `src_001.md` | Top-level `src/` files |
| `src_utils_001.md` | Files from `src/utils/` |
| `src_components_001.md` | Files from `src/components/` |
| `src_components_002.md` | Overflow from `src/components/` (next chunk) |

The three-digit zero-padded counter (`_001`, `_002`) ensures correct alphabetical ordering in all file browsers and RAG systems.

---

## Advanced: Custom Ignore Rules

By default, `code2nlm` uses your project's `.gitignore`. You can specify a different file:

```bash
code2nlm -i ./project --ignore-file .nlmignore
```

The ignore file uses standard `.gitignore` syntax:

```gitignore
# .nlmignore example
node_modules/
dist/
build/
*.min.js
*.bundle.js
coverage/
.env*
__pycache__/
*.pyc
```

**Files always skipped** (hardcoded):
- Hidden files and directories (names starting with `.`)
- Binary and media files (`.exe`, `.png`, `.jpg`, `.mp4`, `.zip`, etc.)
- Already processed output (`nlm_output/`)

---

## Building from Source

```bash
# Clone repository
git clone https://github.com/YOUR_ORG/code2nlm.git
cd code2nlm

# Run tests
go test -v ./...

# Build for current platform
go build -o code2nlm .

# Cross-compile for Windows from Linux (requires Zig)
GOOS=windows GOARCH=amd64 \
  CC="zig cc -target x86_64-windows-gnu" \
  CXX="zig c++ -target x86_64-windows-gnu" \
  CGO_ENABLED=1 go build -trimpath -o code2nlm-windows-amd64.exe .
```

The CI/CD pipeline (`.github/workflows/release.yml`) automatically builds binaries for all supported platforms when a version tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you'd like to change.

```bash
git clone https://github.com/YOUR_ORG/code2nlm.git
cd code2nlm
go test ./...  # must pass before submitting PR
```

---

## License

MIT © code2nlm contributors
