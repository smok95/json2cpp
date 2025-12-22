# json2cpp - JSON to C++ Code Generator

*[한국어](README.ko.md)*

A CLI tool that generates C++ structs and JSON serialization/deserialization code from JSON data. Designed to work with pre-C++11 environments (legacy toolchains like MSVC v120) and generates high-performance code suitable for HFT (High-Frequency Trading) and production environments.

## ✨ Key Features

- **Go-based CLI Tool** - Distributed as a single executable
- **Parser-Agnostic Design** - Supports RapidJSON, nlohmann/json, and JsonCpp via adapter pattern
- **Pre-C++11 Support** - Works with legacy toolchains using C++03 compatible code
- **Nested Objects/Arrays** - Handles complex JSON structures
- **Type Inference** - Merges types from multiple JSON files
- **Optional Field Handling** - Safe type system for null values
- **Embedded Templates** - No external dependencies required

## 🚀 Installation

### From Release (Recommended)

Download the latest binary from [Releases](https://github.com/smok95/json2cpp/releases):
- `json2cpp-linux-amd64` - Linux 64-bit
- `json2cpp-windows-amd64.exe` - Windows 64-bit
- `json2cpp-darwin-amd64` - macOS Intel
- `json2cpp-darwin-arm64` - macOS Apple Silicon

### Build from Source

```bash
# Clone the repository
git clone https://github.com/smok95/json2cpp.git
cd json2cpp

# Build
go build -o json2cpp

# Or install
go install
```

## 📖 Usage

### Basic Usage

```bash
# Generate with default settings
json2cpp -i input.json -o output/

# Pre-C++11 compatible code
json2cpp -i input.json -o output/ --legacy-cpp

# With namespace
json2cpp -i input.json -o output/ --namespace myapp

# Merge multiple JSON files
json2cpp -i "data/*.json" -o output/ --merge
```

### Generated Files

The tool generates parser-agnostic code using the adapter pattern:

```
output/
├── types.h                  # Pure data structures (parser-independent)
├── serializer.h             # Serialization function declarations
├── serializer.cpp           # Serialization implementations
├── json_ptr.h               # C++03-compatible smart pointer
├── json_adapter.h           # Base adapter interface
├── rapidjson_adapter.h/cpp  # RapidJSON implementation
├── nlohmann_adapter.h/cpp   # nlohmann/json implementation
└── jsoncpp_adapter.h/cpp    # JsonCpp implementation
```

### Using Generated Code

```cpp
#include "types.h"
#include "serializer.h"
#include "rapidjson_adapter.h"  // or nlohmann_adapter.h, jsoncpp_adapter.h

// Deserialization example with RapidJSON
rapidjson::Document doc;
doc.Parse(jsonString.c_str());
json2cpp::RapidJsonReader reader(doc);

MyStruct obj;
DeserializeMyStruct(obj, reader);

// Serialization example
rapidjson::Document outDoc;
json2cpp::RapidJsonWriter writer(outDoc, outDoc.GetAllocator());
SerializeMyStruct(obj, writer);
```

## CLI Options

| Option | Description |
|--------|-------------|
| `-i, --input` | Input JSON file (required) |
| `-o, --output` | Output directory (default: `./generated`) |
| `--legacy-cpp` | Generate C++03 compatible code |
| `--namespace` | C++ namespace for generated types |
| `--camelcase` | Use camelCase for field names (default: snake_case) |
| `--optional-null` | Generate Optional&lt;T&gt; for nullable fields |
| `--merge` | Merge multiple JSON files (supports wildcards) |
| `--overwrite` | Overwrite existing files |

## Field Naming Convention

By default, JSON keys are converted to `snake_case`:

| JSON Key | C++ Field (snake_case) | C++ Field (--camelcase) |
|----------|------------------------|-------------------------|
| `"userName"` | `user_name` | `userName` |
| `"isCompress"` | `is_compress` | `isCompress` |
| `"HTTPStatus"` | `http_status` | `httpStatus` |
| `"getHTTPResponse"` | `get_http_response` | `getHttpResponse` |
| `"base64Encode"` | `base_64_encode` | `base64Encode` |

## Type Mapping

| JSON Type | C++ Type |
|-----------|----------|
| Integer | `int64_t` |
| Float | `double` |
| String | `std::string` |
| Boolean | `bool` |
| Null | `Optional<T>` (with `--optional-null`) |
| Object | `struct` |
| Array | `std::vector<T>` |

## JSON Parser Comparison

### RapidJSON (Default)
- **Best Performance** - Optimized for HFT and high-performance systems
- **Full Pre-C++11 Support** - Compatible with legacy compilers
- Header-only, minimal dependencies

### nlohmann/json
- **Best Convenience** - Intuitive and easy-to-use API
- **Requires C++11+** - Uses modern C++ features
- JSON-like syntax in C++

### JsonCpp
- **Stability & Compatibility** - Widely used and well-tested
- **Pre-C++11 Support** - Works in legacy environments
- Reliable and stable

## Pre-C++11 Constraints

- No `std::optional`, `std::string_view`
- No `auto`, range-based for loops
- Custom `UniquePtr` implementation (safe bool idiom for C++03)
- nlohmann/json requires C++11+

## Examples

See `examples/` directory for sample JSON files and generated code.

## Project Structure

```
json2cpp/
├── cmd/                    # CLI command implementation
├── internal/
│   ├── codegen/           # Code generation (adapter pattern)
│   ├── parser/            # JSON parsing
│   ├── nameutil/          # Naming conventions
│   └── types/             # Type system
├── templates/
│   └── adapter/           # Embedded adapter templates
├── examples/              # Example JSON files
├── test/                  # C++ integration tests
├── main.go
└── README.md
```

## Testing

```bash
# Run Go tests
go test ./...

# Run C++ integration tests
cd test
cmake -B build -S .
cmake --build build
./build/test_basic
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details
