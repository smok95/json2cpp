package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"json2cpp/internal/codegen"
	"json2cpp/internal/parser"
	"json2cpp/internal/types"

	"github.com/spf13/cobra"
)

var (
	inputFile    string
	outputDir    string
	legacyCpp    bool
	namespace    string
	camelCase    bool
	optionalNull bool
	merge        bool
	stringRef    bool
	overwrite    bool
)

var rootCmd = &cobra.Command{
	Use:   "json2cpp",
	Short: "JSON to C++ code generator with rapidjson",
	Long: `JSON to C++ code generator that creates structs and serialization code
	for pre-C++11 compatible C++ (legacy toolchains) with rapidjson library.`,
	RunE: run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input JSON file (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./out", "Output directory")
	rootCmd.Flags().BoolVar(&legacyCpp, "legacy-cpp", false, "Generate legacy C++ (pre-C++11) compatible code")
	rootCmd.Flags().StringVar(&namespace, "namespace", "", "C++ namespace")
	rootCmd.Flags().BoolVar(&camelCase, "camelcase", false, "Use camelCase for field names")
	rootCmd.Flags().BoolVar(&optionalNull, "optional-null", false, "Generate Optional<T> for null fields")
	rootCmd.Flags().BoolVar(&merge, "merge", false, "Merge multiple JSON files")
	rootCmd.Flags().BoolVar(&stringRef, "string-ref", false, "Use const std::string& for strings")
	rootCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")

	rootCmd.MarkFlagRequired("input")
}

func run(cmd *cobra.Command, args []string) error {
	// 입력 파일 확인
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// 출력 디렉토리 생성
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 파서 생성
	p := parser.NewParser(legacyCpp, camelCase)

	// 타입 정보 수집
	var allStructs []*types.Struct

	if merge {
		// 여러 JSON 파일 처리 (와일드카드 지원)
		files, err := filepath.Glob(inputFile)
		if err != nil {
			return fmt.Errorf("failed to glob input files: %w", err)
		}

		for _, file := range files {
			structs, err := p.ParseFile(file)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", file, err)
			}
			allStructs = types.MergeTypes(allStructs, structs)
		}
	} else {
		// 단일 파일 처리
		structs, err := p.ParseFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to parse input file: %w", err)
		}
		allStructs = structs
	}

	if len(allStructs) == 0 {
		return fmt.Errorf("no structs generated from input")
	}

	// 코드 생성기 설정
	cfg := codegen.Config{
		LegacyCPP:    legacyCpp,
		Namespace:    namespace,
		CamelCase:    camelCase,
		OptionalNull: optionalNull,
		StringRef:    stringRef,
	}

	gen := codegen.NewGenerator(cfg)

	// 타입 정보
	typeInfo := &types.TypeInfo{
		Structs: allStructs,
	}

	// C++ 코드 생성
	code, err := gen.Generate(typeInfo)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// 출력 파일 결정
	outputFile := filepath.Join(outputDir, "generated.h")
	if !overwrite {
		// 파일이 이미 존재하면 새로운 이름 생성
		counter := 1
		baseName := "generated"
		for {
			outputFile = filepath.Join(outputDir, fmt.Sprintf("%s_%d.h", baseName, counter))
			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				break
			}
			counter++
		}
	}

	// 파일 쓰기
	if err := ioutil.WriteFile(outputFile, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", outputFile)
	fmt.Printf("Structs: %d\n", len(allStructs))

	return nil
}
