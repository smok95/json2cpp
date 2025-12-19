# json2cpp - JSON to C++ Code Generator

JSON 데이터를 입력받아 C++ struct 및 rapidjson 기반 직렬화/역직렬화 코드를 자동 생성하는 CLI 도구입니다. pre-C++11 (legacy toolchains, 예: MSVC v120)에서도 사용할 수 있도록 설계되었으며, HFT(고빈도 거래) 및 실서비스 환경에서 사용될 고성능 코드를 생성합니다.

## ✨ 주요 특징

- **Go 기반 CLI 도구** - 단일 실행 파일로 배포
- **pre-C++11 지원** - C++11 이전 환경(legacy toolchains)에서도 동작하도록 제한된 기능 사용
- **rapidjson 기반** - 고성능 JSON 직렬화/역직렬화
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

### 기본 사용

```bash
json2cpp -i input.json -o out/
```

### pre-C++11 호환 코드 생성

```bash
json2cpp -i input.json -o out/ --legacy-cpp --namespace myapp
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
| `--legacy-cpp` | pre-C++11 (legacy toolchains) 호환 코드 생성 |
| `--namespace` | C++ 네임스페이스 지정 |
| `--camelcase` | camelCase 필드 이름 사용 |
| `--optional-null` | Optional<T> 타입 생성 |
| `--merge` | 여러 JSON 파일 병합 |
# json2cpp

JSON 입력으로 C++ `struct`와 rapidjson 기반 직렬화/역직렬화를 생성하는 CLI 도구입니다.

## 주요 기능

- Go 기반 CLI
- pre-C++11 코드 생성 옵션
- rapidjson 사용
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
json2cpp -i input.json -o out/
json2cpp -i input.json -o out/ --legacy-cpp --namespace myapp
```

## 옵션 (주요)

- `-i, --input`: 입력 JSON 파일 (필수)
- `-o, --output`: 출력 디렉토리 (기본: ./out)
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

## 제약 (pre-C++11)

- `std::optional`, `std::string_view` 미사용
- 일부 최신 문법 제한

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
