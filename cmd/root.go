package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"json2cpp/internal/codegen"
	"json2cpp/internal/parser"
	"json2cpp/internal/types"

	"github.com/spf13/cobra"
)

const version = "1.1.0"

var (
	inputFile     string
	outputDir     string
	parserBackend string
	legacyCpp     bool
	namespace     string
	camelCase     bool
	optionalNull  bool
	merge         bool
	stringRef     bool
	overwrite     bool
	showVersion   bool
)

var rootCmd = &cobra.Command{
	Use:   "json2cpp",
	Short: "JSON to C++ code generator with parser-agnostic adapter pattern",
	Long: `JSON to C++ code generator that creates parser-independent structs and serialization code.
	Uses adapter pattern to support multiple JSON parsers without code changes:
	- RapidJSON
	- nlohmann/json
	- JsonCpp

	Generates separate files: types.h, serializer.h/cpp, and adapter files.
	Requires C++11 or later.`,
	RunE: run,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input JSON file (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./generated", "Output directory for generated files")
	rootCmd.Flags().BoolVar(&legacyCpp, "legacy-cpp", false, "Generate C++03 compatible code")
	rootCmd.Flags().StringVar(&namespace, "namespace", "", "C++ namespace for generated types")
	rootCmd.Flags().BoolVar(&camelCase, "camelcase", false, "Use camelCase for field names (default: snake_case)")
	rootCmd.Flags().BoolVar(&optionalNull, "optional-null", false, "Generate Optional<T> for nullable fields")
	rootCmd.Flags().BoolVar(&merge, "merge", false, "Merge multiple JSON files (supports wildcards)")

	rootCmd.MarkFlagRequired("input")

	// Deprecated flags (kept for compatibility but ignored)
	rootCmd.Flags().StringVarP(&parserBackend, "parser", "p", "", "[DEPRECATED] Parser flag is ignored - adapter pattern supports all parsers")
	rootCmd.Flags().BoolVar(&stringRef, "string-ref", false, "[DEPRECATED] String-ref flag is ignored")
	rootCmd.Flags().BoolVar(&overwrite, "overwrite", false, "[DEPRECATED] Output directory is always overwritten")
	rootCmd.Flags().MarkHidden("parser")
	rootCmd.Flags().MarkHidden("string-ref")
	rootCmd.Flags().MarkHidden("overwrite")
}

func run(cmd *cobra.Command, args []string) error {
	// Check input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Create JSON parser
	p := parser.NewParser(legacyCpp, camelCase)

	// Collect type information
	var allStructs []*types.Struct

	if merge {
		// Process multiple JSON files (supports wildcards)
		files, err := filepath.Glob(inputFile)
		if err != nil {
			return fmt.Errorf("failed to glob input files: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no files matched pattern: %s", inputFile)
		}

		fmt.Printf("Merging %d files...\n", len(files))
		for _, file := range files {
			structs, err := p.ParseFile(file)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", file, err)
			}
			allStructs = types.MergeTypes(allStructs, structs)
		}
	} else {
		// Process single file
		structs, err := p.ParseFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to parse input file: %w", err)
		}
		allStructs = structs
	}

	if len(allStructs) == 0 {
		return fmt.Errorf("no structs generated from input")
	}

	// Configure adapter-based code generator
	cfg := codegen.Config{
		LegacyCPP:    legacyCpp,
		Namespace:    namespace,
		CamelCase:    camelCase,
		OptionalNull: optionalNull,
	}

	// Create adapter generator
	gen := codegen.NewAdapterGenerator(cfg, outputDir)

	// Type information
	typeInfo := &types.TypeInfo{
		Structs: allStructs,
	}

	// Generate all files (types.h, serializer.h/cpp, adapter files)
	fmt.Printf("Generating parser-agnostic code...\n")
	if err := gen.GenerateFiles(typeInfo); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Print summary
	fmt.Printf("\n✓ Generated successfully in: %s\n", outputDir)
	fmt.Printf("  - types.h (data structures)\n")
	fmt.Printf("  - serializer.h/cpp (serialization functions)\n")
	fmt.Printf("  - json_adapter.h (base interface)\n")
	fmt.Printf("  - rapidjson_adapter.h/cpp\n")
	fmt.Printf("  - nlohmann_adapter.h/cpp\n")
	fmt.Printf("  - jsoncpp_adapter.h/cpp\n")
	fmt.Printf("\nStructs: %d\n", len(allStructs))
	if legacyCpp {
		fmt.Printf("Mode: C++03 compatible (deprecated, use C++11+)\n")
	} else {
		fmt.Printf("Mode: Modern C++ (C++11+)\n")
	}

	return nil
}
