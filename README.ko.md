# json2cpp - JSON to C++ Code Generator

*[English](README.md)*

JSON 데이터를 입력받아 C++ struct 및 JSON 직렬화/역직렬화 코드를 자동 생성하는 CLI 도구입니다. pre-C++11 (legacy toolchains, 예: MSVC v120)에서도 사용할 수 있도록 설계되었으며, HFT(고빈도 거래) 및 실서비스 환경에서 사용될 고성능 코드를 생성합니다.

## ✨ 주요 특징

- **Go 기반 CLI 도구** - 단일 실행 파일로 배포
- **파서 독립적 설계** - RapidJSON, nlohmann/json, JsonCpp를 어댑터 패턴으로 지원
- **Pre-C++11 지원** - C++03 호환 코드로 레거시 컴파일러에서 동작
- **중첩 객체/배열 지원** - 복잡한 JSON 구조 처리
- **타입 추론** - 여러 JSON 파일에서 공통 타입 생성
- **Optional 필드 처리** - null 값을 위한 안전한 타입 시스템
- **임베디드 템플릿** - 외부 의존성 없이 단독 실행

## 🚀 설치

### 릴리스에서 다운로드 (권장)

[Releases](https://github.com/smok95/json2cpp/releases)에서 최신 바이너리를 다운로드하세요:
- `json2cpp-linux-amd64` - Linux 64비트
- `json2cpp-windows-amd64.exe` - Windows 64비트
- `json2cpp-darwin-amd64` - macOS Intel
- `json2cpp-darwin-arm64` - macOS Apple Silicon

### 소스에서 빌드

```bash
# 저장소 복제
git clone https://github.com/smok95/json2cpp.git
cd json2cpp

# 빌드
go build -o json2cpp

# 또는 설치
go install
```

## 📖 사용법

### 기본 사용

```bash
# 기본 설정으로 생성
json2cpp -i input.json -o output/

# Pre-C++11 호환 코드 생성
json2cpp -i input.json -o output/ --legacy-cpp

# 네임스페이스 지정
json2cpp -i input.json -o output/ --namespace myapp

# 여러 JSON 파일 병합
json2cpp -i "data/*.json" -o output/ --merge
```

### 생성되는 파일

어댑터 패턴을 사용하여 파서 독립적인 코드를 생성합니다:

```
output/
├── types.h                  # 순수 데이터 구조 (파서 독립적)
├── serializer.h             # 직렬화 함수 선언
├── serializer.cpp           # 직렬화 구현
├── json_ptr.h               # C++03 호환 스마트 포인터
├── json_adapter.h           # 어댑터 기본 인터페이스
├── rapidjson_adapter.h/cpp  # RapidJSON 구현
├── nlohmann_adapter.h/cpp   # nlohmann/json 구현
└── jsoncpp_adapter.h/cpp    # JsonCpp 구현
```

### 생성된 코드 사용

```cpp
#include "types.h"
#include "serializer.h"
#include "rapidjson_adapter.h"  // 또는 nlohmann_adapter.h, jsoncpp_adapter.h

// RapidJSON을 사용한 역직렬화 예제
rapidjson::Document doc;
doc.Parse(jsonString.c_str());
json2cpp::RapidJsonReader reader(doc);

MyStruct obj;
DeserializeMyStruct(obj, reader);

// 직렬화 예제
rapidjson::Document outDoc;
json2cpp::RapidJsonWriter writer(outDoc, outDoc.GetAllocator());
SerializeMyStruct(obj, writer);
```

## CLI 옵션

| 옵션 | 설명 |
|------|------|
| `-i, --input` | 입력 JSON 파일 (필수) |
| `-o, --output` | 출력 디렉토리 (기본값: `./generated`) |
| `--legacy-cpp` | C++03 호환 코드 생성 |
| `--namespace` | C++ 네임스페이스 지정 |
| `--camelcase` | 필드명을 camelCase로 생성 (기본값: snake_case) |
| `--optional-null` | null 가능 필드에 Optional&lt;T&gt; 생성 |
| `--merge` | 여러 JSON 파일 병합 (와일드카드 지원) |
| `--overwrite` | 기존 파일 덮어쓰기 |

## 필드 이름 변환 규칙

기본적으로 JSON 키는 `snake_case`로 변환됩니다:

| JSON 키 | C++ 필드 (snake_case) | C++ 필드 (--camelcase) |
|---------|------------------------|-------------------------|
| `"userName"` | `user_name` | `userName` |
| `"isCompress"` | `is_compress` | `isCompress` |
| `"HTTPStatus"` | `http_status` | `httpStatus` |
| `"getHTTPResponse"` | `get_http_response` | `getHttpResponse` |
| `"base64Encode"` | `base_64_encode` | `base64Encode` |

## 타입 매핑

| JSON 타입 | C++ 타입 |
|-----------|----------|
| 정수 | `int64_t` |
| 실수 | `double` |
| 문자열 | `std::string` |
| 불리언 | `bool` |
| null | `Optional<T>` (`--optional-null` 사용 시) |
| 객체 | `struct` |
| 배열 | `std::vector<T>` |

## JSON 파서 비교

### RapidJSON (기본)
- **최고의 성능** - HFT 및 고성능 시스템에 최적화
- **Pre-C++11 완벽 지원** - 레거시 컴파일러 호환
- 헤더 온리, 최소한의 의존성

### nlohmann/json
- **최고의 편의성** - 직관적이고 사용하기 쉬운 API
- **C++11 이상 필요** - 모던 C++ 기능 활용
- C++에서 JSON과 유사한 문법

### JsonCpp
- **안정성과 호환성** - 널리 사용되는 검증된 라이브러리
- **Pre-C++11 지원** - 레거시 환경에서 사용 가능
- 신뢰할 수 있고 안정적

## Pre-C++11 제약사항

- `std::optional`, `std::string_view` 미사용
- `auto`, range-based for 루프 미사용
- 커스텀 `UniquePtr` 구현 (C++03용 safe bool idiom)
- nlohmann/json은 C++11 이상 필요

## 예제

`examples/` 디렉토리에서 샘플 JSON 파일과 생성된 코드를 확인하세요.

## 프로젝트 구조

```
json2cpp/
├── cmd/                    # CLI 명령 구현
├── internal/
│   ├── codegen/           # 코드 생성 (어댑터 패턴)
│   ├── parser/            # JSON 파싱
│   ├── nameutil/          # 명명 규칙
│   └── types/             # 타입 시스템
├── templates/
│   └── adapter/           # 임베디드 어댑터 템플릿
├── examples/              # 예제 JSON 파일
├── test/                  # C++ 통합 테스트
├── main.go
└── README.md
```

## 테스트

```bash
# Go 테스트 실행
go test ./...

# C++ 통합 테스트 실행
cd test
cmake -B build -S .
cmake --build build
./build/test_basic
```

## 기여

1. 저장소 Fork
2. 기능 브랜치 생성 (`git checkout -b feature/amazing-feature`)
3. 변경사항 커밋 (`git commit -m 'Add some amazing feature'`)
4. 브랜치에 Push (`git push origin feature/amazing-feature`)
5. Pull Request 생성

## 라이선스

MIT License - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요
