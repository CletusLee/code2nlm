package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

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
	RunE: func(cmd *cobra.Command, args []string) error {
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
		OutputPath:  OutputPath,
		ProjectName: filepath.Base(filepath.Clean(InputPath)),
	}

	return c.Process(virtualTree)
}
