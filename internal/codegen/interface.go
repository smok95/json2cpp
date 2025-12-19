package codegen

import "json2cpp/internal/types"

// ParserType represents the JSON parser backend to use
type ParserType string

const (
	ParserRapidJSON ParserType = "rapidjson"
	ParserNlohmann  ParserType = "nlohmann"
	ParserJsonCpp   ParserType = "jsoncpp"
)

// Config holds configuration for code generation
type Config struct {
	Parser       ParserType
	LegacyCPP    bool
	Namespace    string
	CamelCase    bool
	OptionalNull bool
	StringRef    bool
}

// CodeGenerator is the interface that all parser-specific generators must implement
type CodeGenerator interface {
	// Generate generates C++ code from TypeInfo
	Generate(info *types.TypeInfo) (string, error)
}

// NewGenerator creates a new code generator based on the parser type
func NewGenerator(cfg Config) CodeGenerator {
	// Default to rapidjson if not specified
	if cfg.Parser == "" {
		cfg.Parser = ParserRapidJSON
	}

	switch cfg.Parser {
	case ParserRapidJSON:
		return NewRapidJSONGenerator(cfg)
	case ParserNlohmann:
		return NewNlohmannGenerator(cfg)
	case ParserJsonCpp:
		return NewJsonCppGenerator(cfg)
	default:
		// Fallback to rapidjson
		return NewRapidJSONGenerator(cfg)
	}
}
