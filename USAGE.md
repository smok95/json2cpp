# json2cpp 사용 가이드

## 빠른 시작

### 1. 빌드

```bash
# Go 1.21 이상 필요
go mod download
go build -o json2cpp
```

### 2. 기본 사용법

```bash
# 기본 (RapidJSON)
./json2cpp -i examples/order.json -o out/

# nlohmann/json 사용
./json2cpp -i examples/order.json -o out/ --parser nlohmann

# JsonCpp 사용
./json2cpp -i examples/order.json -o out/ --parser jsoncpp

# pre-C++11 호환 코드 생성
./json2cpp -i examples/order.json -o out/ --legacy-cpp --parser rapidjson

# 네임스페이스 지정
./json2cpp -i examples/order.json -o out/ --namespace trading
```

### 3. 고급 사용법

#### 여러 JSON 파일 병합

```bash
# data 폴더의 모든 JSON 파일을 병합
./json2cpp -i "data/*.json" -o out/ --merge --legacy-cpp
```

#### Optional 필드 처리

```bash
# null 값을 Optional<T>로 처리
./json2cpp -i input.json -o out/ --optional-null --legacy-cpp
```

#### 필드 명명 규칙

```bash
# camelCase 사용
./json2cpp -i input.json -o out/ --camelcase

# snake_case 사용 (기본)
./json2cpp -i input.json -o out/
```

#### 문자열 최적화

```bash
# const std::string& 사용
./json2cpp -i market_data.json -o include/ --legacy-cpp --string-ref
```

## JSON 파서 선택 가이드

### RapidJSON (기본값)

**언제 사용하나요?**
- 최고의 성능이 필요할 때 (HFT, 고빈도 거래)
- pre-C++11 환경 (MSVC v120, GCC 4.8 등)
- 메모리 효율성이 중요할 때

**생성된 코드:**

```cpp
struct Order {
    int64_t id;
    std::string symbol;
    double price;

    void FromJson(const rapidjson::Value& v) {
        if (v.HasMember("id") && v["id"].IsInt64()) {
            id = v["id"].GetInt64();
        }
        // ...
    }

    void ToJson(rapidjson::Value& v, rapidjson::Document::AllocatorType& a) const {
        v.AddMember("id", rapidjson::Value().SetInt64(id), a);
        // ...
    }
};
```

**사용 예제:**

```cpp
#include "generated.h"
#include "rapidjson/document.h"

std::string json_str = R"({"id": 123, "symbol": "AAPL", "price": 150.5})";
rapidjson::Document doc;
doc.Parse(json_str.c_str());

Order order;
order.FromJson(doc);
```

### nlohmann/json

**언제 사용하나요?**
- 개발 편의성과 가독성이 우선일 때
- C++11 이상 환경
- 프로토타이핑 및 빠른 개발

**생성된 코드:**

```cpp
struct Order {
    int64_t id;
    std::string symbol;
    double price;

    void from_json(const nlohmann::json& j) {
        if (j.contains("id")) {
            id = j["id"].get<int64_t>();
        }
        // ...
    }

    nlohmann::json to_json() const {
        nlohmann::json j;
        j["id"] = id;
        j["symbol"] = symbol;
        j["price"] = price;
        return j;
    }
};

// ADL 헬퍼
inline void from_json(const nlohmann::json& j, Order& obj) {
    obj.from_json(j);
}

inline void to_json(nlohmann::json& j, const Order& obj) {
    j = obj.to_json();
}
```

**사용 예제:**

```cpp
#include "generated.h"
#include <nlohmann/json.hpp>

std::string json_str = R"({"id": 123, "symbol": "AAPL", "price": 150.5})";
nlohmann::json j = nlohmann::json::parse(json_str);

Order order;
order.from_json(j);

// 또는 ADL 사용
Order order2 = j.get<Order>();
```

### JsonCpp

**언제 사용하나요?**
- 안정성과 검증된 라이브러리가 필요할 때
- pre-C++11 환경에서 RapidJSON 대안이 필요할 때
- 기존 JsonCpp 프로젝트와 통합

**생성된 코드:**

```cpp
struct Order {
    int64_t id;
    std::string symbol;
    double price;

    void FromJson(const Json::Value& v) {
        if (v.isMember("id") && v["id"].isInt64()) {
            id = v["id"].asInt64();
        }
        // ...
    }

    void ToJson(Json::Value& v) const {
        v["id"] = Json::Value::Int64(id);
        v["symbol"] = symbol;
        v["price"] = price;
    }
};
```

**사용 예제:**

```cpp
#include "generated.h"
#include <json/json.h>

std::string json_str = R"({"id": 123, "symbol": "AAPL", "price": 150.5})";
Json::Value root;
Json::Reader reader;
reader.parse(json_str, root);

Order order;
order.FromJson(root);
```

## 타입 추론 규칙

### 1. 단일 파일 처리

입력 JSON의 각 필드는 다음 규칙에 따라 C++ 타입으로 변환됩니다:

```json
{
  "int_field": 123,           // int64_t
  "float_field": 123.45,      // double
  "string_field": "text",     // std::string
  "bool_field": true,         // bool
  "null_field": null,         // Optional<T> (--optional-null 옵션)
  "array_field": [1, 2, 3],   // std::vector<int64_t>
  "object_field": {}          // struct
}
```

### 2. 다중 파일 병합

여러 JSON 파일을 처리할 때:

- **같은 이름의 필드**: 타입이 다르다면 상위 타입으로 승격
- **존재하지 않는 필드**: Optional로 처리
- **배열 요소**: 모든 파일에서 추론된 공통 타입 사용

예:

```json
// file1.json
{"price": 100}

// file2.json
{"price": 99.9}

// 결과: price는 double로 추론됨
```

## 생성된 코드 구조 비교

| 기능 | RapidJSON | nlohmann/json | JsonCpp |
|------|-----------|---------------|---------|
| **역직렬화 메서드** | `FromJson(Value&)` | `from_json(json&)` | `FromJson(Value&)` |
| **직렬화 메서드** | `ToJson(Value&, Alloc&)` | `to_json()` → `json` | `ToJson(Value&)` |
| **ADL 지원** | ❌ | ✅ | ❌ |
| **pre-C++11** | ✅ | ❌ (C++11+) | ✅ |
| **성능** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **편의성** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

## Optional 타입

`--optional-null` 옵션 사용 시:

### pre-C++11 모드

```cpp
template<typename T>
struct Optional {
    bool has;
    T value;

    Optional() : has(false) {}
    Optional(const T& v) : has(true), value(v) {}

    bool IsValid() const { return has; }
    const T& Get() const { return value; }
    void Set(const T& v) { has = true; value = v; }
    void Clear() { has = false; }
};
```

### C++17 모드 (nlohmann만 해당)

```cpp
#include <optional>

template<typename T>
using Optional = std::optional<T>;
```

## 통합 예제

### 1. HFT 거래 시스템 (RapidJSON)

```bash
./json2cpp -i orders.json -o include/ --legacy-cpp --namespace hft --parser rapidjson
```

```cpp
#include "generated.h"

namespace hft {
    Order order;
    rapidjson::Document doc;
    doc.Parse(json_data);
    order.FromJson(doc);
}
```

### 2. 마이크로서비스 (nlohmann/json)

```bash
./json2cpp -i api_schema.json -o src/ --parser nlohmann --namespace api
```

```cpp
#include "generated.h"

namespace api {
    auto j = nlohmann::json::parse(request_body);
    Request req = j.get<Request>();
}
```

### 3. 설정 파일 처리 (JsonCpp)

```bash
./json2cpp -i "config/*.json" -o src/ --merge --parser jsoncpp
```

```cpp
#include "generated.h"

Json::Value root;
Json::Reader reader;
reader.parse(config_file, root);

Config config;
config.FromJson(root);
```

## 성능 최적화 팁

### RapidJSON

```cpp
// 메모리 풀 재사용
rapidjson::Document doc;
doc.Parse(json_string);

rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();
order.ToJson(doc, allocator);
```

### nlohmann/json

```cpp
// 복사 방지
const nlohmann::json& j = nlohmann::json::parse(str);
Order order;
order.from_json(j);  // 복사 없이 참조 전달
```

### JsonCpp

```cpp
// FastWriter 사용
Json::FastWriter writer;
Json::Value root;
order.ToJson(root);
std::string output = writer.write(root);
```

## 문제 해결

### 빌드 오류

```bash
# 의존성 확인
go mod tidy

# Go 버전 확인
go version  # 1.21+ 필요

# 빌드
go build -o json2cpp
```

### 생성된 코드 컴파일 오류

#### RapidJSON 헤더 누락

```cpp
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
```

#### nlohmann/json 설치

```bash
# vcpkg
vcpkg install nlohmann-json

# CMake
find_package(nlohmann_json 3.11.0 REQUIRED)
target_link_libraries(myapp nlohmann_json::nlohmann_json)
```

#### JsonCpp 설치

```bash
# Ubuntu/Debian
sudo apt-get install libjsoncpp-dev

# CMake
find_package(jsoncpp REQUIRED)
target_link_libraries(myapp jsoncpp_lib)
```

### C++ 표준 설정

```cmake
# pre-C++11 (RapidJSON, JsonCpp)
set(CMAKE_CXX_STANDARD 98)

# C++11+ (nlohmann/json)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
```

### 런타임 오류

1. **JSON 키 매칭 실패**
   - JSON 키 이름이 C++ 필드 이름과 일치하는지 확인
   - `--camelcase` 옵션 사용 여부 확인

2. **타입 변환 실패**
   - JSON 데이터 타입이 예상과 일치하는지 확인
   - Optional 필드 처리 확인

3. **파서별 API 차이**
   - RapidJSON: `HasMember()`, `IsInt64()`, `GetInt64()`
   - nlohmann/json: `contains()`, `is_number_integer()`, `get<T>()`
   - JsonCpp: `isMember()`, `isInt64()`, `asInt64()`

## CLI 옵션 전체 목록

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-i, --input` | 입력 JSON 파일 (필수) | - |
| `-o, --output` | 출력 디렉토리 | `./out` |
| `-p, --parser` | JSON 파서 (rapidjson, nlohmann, jsoncpp) | `rapidjson` |
| `--legacy-cpp` | pre-C++11 호환 코드 생성 | `false` |
| `--namespace` | C++ 네임스페이스 | 없음 |
| `--camelcase` | camelCase 필드 이름 사용 | `false` |
| `--optional-null` | Optional<T> 타입 생성 | `false` |
| `--merge` | 여러 JSON 파일 병합 | `false` |
| `--string-ref` | const std::string& 사용 | `false` |
| `--overwrite` | 기존 파일 덮어쓰기 | `false` |

## 지원

- GitHub Issues: 프로젝트 이슈 트래커
- 문서: README.md 참조
- 예제: examples/ 폴더 확인
