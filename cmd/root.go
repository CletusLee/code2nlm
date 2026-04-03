package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"code2nlm/chunking"
	"code2nlm/markdown"
	"code2nlm/scanner"
)

var (
	InputPath  string
	OutputPath string
	MaxSources int
	MaxWords   int
	IgnoreFile string
	Strategy   string
	FS         afero.Fs
)

var rootCmd = &cobra.Command{
	Use:   "code2nlm",
	Short: "Chunk codebases into LLM-friendly Markdown files",
	Long: `code2nlm is a fast, local-only CLI tool that transforms entire code repositories 
into structured Markdown notebooks optimized for LLMs like NotebookLM.

It supports:
- Automatic path normalization for cross-platform consistency
- Smart AST-aware chunking for Go (or fallback text-splitting)
- Gitignore-respecting concurrent directory scanning
- Dynamic project indexing with spatial context mapping
- Built-in Denoising: Automatically strips massive inline Source Maps (*.map URIs) and Base64 Data URIs to preserve valuable LLM token context.`,
	Example: `  # Process current directory into default nlm_output
  code2nlm

  # Process specific directory with custom word limit
  code2nlm -i ./src -o ./md_out -w 50000

  # Allow for a larger number of output source files
  code2nlm -m 300

  # Manual ignore file override
  code2nlm --ignore-file .nlignore`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no flags were changed and no arguments provided, show help
		if !cmd.Flags().Changed("input") &&
			!cmd.Flags().Changed("output") &&
			!cmd.Flags().Changed("max-sources") &&
			!cmd.Flags().Changed("max-words") &&
			!cmd.Flags().Changed("ignore-file") &&
			!cmd.Flags().Changed("strategy") &&
			len(args) == 0 {
			return cmd.Help()
		}
		return runChunking()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	FS = afero.NewOsFs()
	rootCmd.Flags().StringVarP(&InputPath, "input", "i", "./", "Input directory path")
	rootCmd.Flags().StringVarP(&OutputPath, "output", "o", "./nlm_output", "Output directory")
	rootCmd.Flags().IntVarP(&MaxSources, "max-sources", "m", 50, "Max number of output files")
	rootCmd.Flags().IntVarP(&MaxWords, "max-words", "w", 100000, "Max word count per file")
	rootCmd.Flags().StringVar(&IgnoreFile, "ignore-file", ".gitignore", "Path to ignore list")
	rootCmd.Flags().StringVarP(&Strategy, "strategy", "s", "ast", "Chunking strategy (dir or ast)")
}

func runChunking() error {
	// Sanitize paths to fix Windows PowerShell escaping issues where trailing slashes escape closing quotes.
	InputPath = filepath.Clean(strings.Trim(InputPath, "\" '"))
	OutputPath = filepath.Clean(strings.Trim(OutputPath, "\" '"))

	totalBytes, virtualTree, err := scanner.ScanDirectory(FS, InputPath, IgnoreFile)
	if err != nil {
		return err
	}

	estimatedWords := float64(totalBytes) / 5.0
	requiredFiles := int(math.Ceil(estimatedWords / float64(MaxWords)))

	if requiredFiles > MaxSources {
		msg := fmt.Sprintf("FATAL: Project is too large for the current configuration. It requires at least %d files, but --max-sources is set to %d. Please either increase --max-sources (-m %d), increase --max-words, or exclude more directories.", requiredFiles, MaxSources, requiredFiles)
		return fmt.Errorf(msg)
	}

	err = markdown.GenerateIndex(FS, OutputPath, virtualTree)
	if err != nil {
		return err
	}

	c := &chunking.Chunker{
		FS:          FS,
		MaxWords:    MaxWords,
		InputPath:   InputPath,
		OutputPath:  OutputPath,
		ProjectName: filepath.Base(filepath.Clean(InputPath)),
	}

	return c.Process(virtualTree)
}
