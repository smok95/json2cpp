# json2cpp 사용 가이드

## 빠른 시작

### 1. 빌드

```bash
# Go 1.21 이상 필요
go mod download
go build -o json2cpp
```

./json2cpp -i examples/order.json -o out/ --legacy-cpp

```bash
# 단일 JSON 파일 처리
./json2cpp -i examples/order.json -o out/

# VS2013 호환 코드 생성
./json2cpp -i examples/order.json -o out/ --vs2013

./json2cpp -i "data/*.json" -o out/ --merge --legacy-cpp
./json2cpp -i examples/order.json -o out/ --namespace trading
```

### 3. 고급 사용법

#### 여러 JSON 파일 병합
./json2cpp -i input.json -o out/ --optional-null --legacy-cpp
```bash
# data 폴더의 모든 JSON 파일을 병합
./json2cpp -i "data/*.json" -o out/ --merge --vs2013
```

// pre-C++11 환경에서는 명시적 타입 체크

```bash
# null 값을 Optional<T>로 처리
./json2cpp -i input.json -o out/ --optional-null --vs2013
```

#### 필드 명명 규칙

./json2cpp -i orders.json -o include/ --legacy-cpp --namespace hft
# camelCase 사용
./json2cpp -i input.json -o out/ --camelcase

# snake_case 사용 (기본)
./json2cpp -i input.json -o out/
```
./json2cpp -i market_data.json -o include/ --legacy-cpp --string-ref
## 타입 추론 규칙

### 1. 단일 파일 처리

입력 JSON의 각 필드는 다음 규칙에 따라 C++ 타입으로 변환됩니다:

```json
{
  "int_field": 123,           // int64_t
  "float_field": 123.45,      // double
  "string_field": "text",     // std::string
  "bool_field": true,         // bool
  "null_field": null,         // Optional<T> (옵션)
  "array_field": [1, 2, 3],   // std::vector<int64_t>
  "object_field": {}          // struct
}
```

### 2. 다중 파일 병합

여러 JSON 파일을 처리할 때:

- **같은 이름의 필드**: 타입이 다륈다면 상위 타입으로 승격
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

## 생성된 코드 구조

### 기본 구조

```cpp
struct GeneratedStruct {
    // 멤버 변수
    int64_t id;
    std::string name;
    std::vector<Item> items;
    
    // 역직렬화
    void FromJson(const rapidjson::Value& v);
    
    // 직렬화
    void ToJson(rapidjson::Value& v, 
                rapidjson::Document::AllocatorType& a) const;
};
```

### VS2013 특화 기능

#### Optional 타입

```cpp
template<typename T>
struct Optional {
    bool has;
    T value;
    
    bool IsValid() const { return has; }
    const T& Get() const { return value; }
    void Set(const T& v) { has = true; value = v; }
    void Clear() { has = false; }
};
```

#### 안전한 타입 체크

```cpp
// VS2013에서는 명시적 타입 체크
if (v.HasMember("id") && v["id"].IsInt64()) {
    id = v["id"].GetInt64();
} else if (v.HasMember("id") && v["id"].IsInt()) {
    id = static_cast<int64_t>(v["id"].GetInt());
}
```

## 통합 예제

### 1. HFT 거래 시스템

```bash
# 거래 명령 JSON 생성
./json2cpp -i orders.json -o include/ --vs2013 --namespace hft
```

### 2. 설정 파일 처리

```bash
# 설정 파일 스키마 생성
./json2cpp -i "config/*.json" -o src/ --merge --camelcase
```

### 3. 마켓 데이터

```bash
# 실시간 마켓 데이터 구조
./json2cpp -i market_data.json -o include/ --vs2013 --string-ref
```

## 성능 최적화 팁

### 1. rapidjson 설정

```cpp
// Document에 allocator 설정
rapidjson::Document doc;
doc.Parse(json_string);

// 재사용을 위해 allocator 참조 전달
rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();
order.ToJson(doc, allocator);
```

### 2. String 처리

```bash
# 대량의 문자열 처리 시
./json2cpp -i data.json -o out/ --string-ref
```

### 3. 메모리 풀 사용

```cpp
// rapidjson 메모리 풀 사용
rapidjson::MemoryPoolAllocator<> allocator;
rapidjson::Document doc(&allocator);
```

## 문제 해결

### 빌드 오류

```bash
# 의존성 확인
go mod tidy

# Go 버전 확인
go version  # 1.21+ 필요
```

### 생성된 코드 컴파일 오류

1. **rapidjson 헤더 누락**
   ```cpp
   #include "rapidjson/document.h"
   #include "rapidjson/writer.h"
   #include "rapidjson/stringbuffer.h"
   ```

2. **C++11 설정 확인**
   ```cmake
   set(CMAKE_CXX_STANDARD 11)
   set(CMAKE_CXX_STANDARD_REQUIRED ON)
   ```

3. **VS2013 툴셋 확인**
   - 프로젝트 속성 → 플랫폼 도구집합 → v120

### 런타임 오류

1. **JSON 키 매칭 실패**
   - JSON 키 이름이 C++ 필드 이름과 일치하는지 확인
   - `--camelcase` 옵션 사용 여부 확인

2. **타입 변환 실패**
   - JSON 데이터 타입이 예상과 일치하는지 확인
   - Optional 필드 처리 확인

## 지원

- GitHub Issues: 프로젝트 이슈 트래커
- 문서: README.md 참조
- 예제: examples/ 폴터 확인