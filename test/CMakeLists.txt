cmake_minimum_required(VERSION 3.10)
project(json2cpp_test)

# C++ 표준 설정
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# rapidjson 헤더 경로 설정
# rapidjson은 헤더 온리 라이브러리이므로 헤더 경로만 필요
include(FetchContent)
FetchContent_Declare(
    rapidjson
    GIT_REPOSITORY https://github.com/Tencent/rapidjson.git
    GIT_TAG v1.1.0
)
FetchContent_MakeAvailable(rapidjson)

# 생성된 헤더 파일 경로
set(GENERATED_DIR ${CMAKE_SOURCE_DIR}/../out)

# 기본 테스트 (가장 간단한 테스트 - CI에서 사용)
add_executable(test_basic test_basic.cpp ${GENERATED_DIR}/serializer_rapidjson.cpp)
target_include_directories(test_basic PRIVATE
    ${rapidjson_SOURCE_DIR}/include
    ${GENERATED_DIR}
)

# 테스트 실행 파일 (예제 파일 필요)
add_executable(test_serialization test_serialization.cpp)
target_include_directories(test_serialization PRIVATE
    ${rapidjson_SOURCE_DIR}/include
    ${GENERATED_DIR}
)

# 간단한 테스트 (단일 JSON 파일)
add_executable(test_simple test_simple.cpp)
target_include_directories(test_simple PRIVATE
    ${rapidjson_SOURCE_DIR}/include
    ${GENERATED_DIR}
)

# 컴파일 옵션
if(MSVC)
    target_compile_options(test_basic PRIVATE /W4)
    target_compile_options(test_serialization PRIVATE /W4)
    target_compile_options(test_simple PRIVATE /W4)
else()
    target_compile_options(test_basic PRIVATE -Wall -Wextra -pedantic)
    target_compile_options(test_serialization PRIVATE -Wall -Wextra -pedantic)
    target_compile_options(test_simple PRIVATE -Wall -Wextra -pedantic)
endif()
