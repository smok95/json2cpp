# json2cpp - JSON to C++ Code Generator

JSON 데이터를 입력받아 C++ struct 및 JSON 직렬화/역직렬화 코드를 자동 생성하는 CLI 도구입니다. pre-C++11 (legacy toolchains, 예: MSVC v120)에서도 사용할 수 있도록 설계되었으며, HFT(고빈도 거래) 및 실서비스 환경에서 사용될 고성능 코드를 생성합니다.

## ✨ 주요 특징

- **Go 기반 CLI 도구** - 단일 실행 파일로 배포
- **pre-C++11 지원** - C++11 이전 환경(legacy toolchains)에서도 동작하도록 제한된 기능 사용
- **다중 파서 지원** - RapidJSON, nlohmann/json, JsonCpp 중 선택 가능
- **중첩 객체/배열 지원** - 복잡한 JSON 구조 처리
- **타입 병합 추론** - 여러 JSON 파일에서 공통 타입 생성
- **Optional 필드 처리** - null 값을 위한 안전한 타입 시스템

## 🚀 설치

```bash
# Go가 설치되어 있다면
go install github.com/yourusername/json2cpp@latest

# 또는 소스에서 직접 빌드
git clone https://github.com/yourusername/json2cpp.git
cd json2cpp
go build -o json2cpp
```

## 📖 사용법

### 기본 사용 (RapidJSON - 기본값)

```bash
json2cpp -i input.json -o out/
```

### nlohmann/json 사용

```bash
json2cpp -i input.json -o out/ --parser nlohmann
```

### JsonCpp 사용

```bash
json2cpp -i input.json -o out/ --parser jsoncpp
```

### pre-C++11 호환 코드 생성

```bash
json2cpp -i input.json -o out/ --legacy-cpp --namespace myapp --parser rapidjson
```

### 여러 JSON 파일 병합

```bash
json2cpp -i "data/*.json" -o out/ --merge --legacy-cpp
```

### CLI 옵션

| 옵션 | 설명 |
|------|------|
| `-i, --input` | 입력 JSON 파일 (필수) |
| `-o, --output` | 출력 디렉토리 (기본값: ./out) |
| `-p, --parser` | JSON 파서 선택 (rapidjson, nlohmann, jsoncpp) (기본값: rapidjson) |
| `--legacy-cpp` | pre-C++11 (legacy toolchains) 호환 코드 생성 |
| `--namespace` | C++ 네임스페이스 지정 |
| `--camelcase` | camelCase 필드 이름 사용 |
| `--optional-null` | Optional<T> 타입 생성 |
| `--merge` | 여러 JSON 파일 병합 |
| `--string-ref` | const std::string& 사용 |
| `--overwrite` | 기존 파일 덮어쓰기 |
# json2cpp

JSON 입력으로 C++ `struct`와 JSON 직렬화/역직렬화를 생성하는 CLI 도구입니다.

## 주요 기능

- Go 기반 CLI
- pre-C++11 코드 생성 옵션
- **다중 JSON 파서 지원**: RapidJSON, nlohmann/json, JsonCpp
- 중첩 객체 및 배열 지원
- 여러 파일 병합 시 타입 병합

## 설치

```bash
go install github.com/yourusername/json2cpp@latest
# or build from source
git clone https://github.com/yourusername/json2cpp.git
cd json2cpp
go build -o json2cpp
```

## 사용법

```bash
# 기본 (RapidJSON)
json2cpp -i input.json -o out/

# nlohmann/json 사용
json2cpp -i input.json -o out/ --parser nlohmann

# JsonCpp 사용
json2cpp -i input.json -o out/ --parser jsoncpp

# Legacy C++ 모드
json2cpp -i input.json -o out/ --legacy-cpp --namespace myapp
```

## 옵션 (주요)

- `-i, --input`: 입력 JSON 파일 (필수)
- `-o, --output`: 출력 디렉토리 (기본: ./out)
- `-p, --parser`: JSON 파서 (rapidjson, nlohmann, jsoncpp) (기본: rapidjson)
- `--legacy-cpp`: pre-C++11 호환 코드 생성
- `--namespace`: C++ 네임스페이스 지정
- `--camelcase`: camelCase 필드 이름 사용
- `--optional-null`: Optional<T> 타입 생성
- `--merge`: 여러 JSON 파일 병합
- `--string-ref`: const std::string& 사용
- `--overwrite`: 기존 파일 덮어쓰기

## 타입 매핑

- JSON 정수 → `int64_t`
- JSON 실수 → `double`
- JSON 문자열 → `std::string`
- JSON 불리언 → `bool`
- JSON null → `Optional<T>` (옵션)
- JSON 객체 → `struct`
- JSON 배열 → `std::vector<T>`

## JSON 파서별 특징

### RapidJSON (기본)
- **최고의 성능** - HFT 및 고성능 시스템에 최적화
- **pre-C++11 완벽 지원** - 레거시 컴파일러 호환
- 메서드: `FromJson(rapidjson::Value&)`, `ToJson(rapidjson::Value&, AllocatorType&)`

### nlohmann/json
- **최고의 편의성** - 직관적이고 사용하기 쉬운 API
- **C++11 이상 필요** - 모던 C++ 기능 활용
- 메서드: `from_json(nlohmann::json&)`, `to_json()`
- ADL (Argument-Dependent Lookup) 헬퍼 함수 자동 생성

### JsonCpp
- **안정성과 호환성** - 널리 사용되는 검증된 라이브러리
- **pre-C++11 지원** - 레거시 환경에서 사용 가능
- 메서드: `FromJson(Json::Value&)`, `ToJson(Json::Value&)`

## 제약 (pre-C++11 모드)

- `std::optional`, `std::string_view` 미사용
- 일부 최신 문법 제한
- nlohmann/json은 C++11 이상 필요

## 예제

입력 예제 및 생성된 구조체는 `examples/`를 참고하세요.

## 프로젝트 구조

```
json2cpp/
├── cmd/
├── internal/
│   ├── codegen/
│   ├── parser/
│   └── types/
├── examples/
├── main.go
└── README.md
```

## 테스트

```bash
go test ./...
```

## 기여

Fork → 브랜치 생성 → 커밋 → PR 요청

## 라이선스

MIT
