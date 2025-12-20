cmake_minimum_required(VERSION 3.10)
project(json2cpp_test)

# C++ 표준 설정
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# rapidjson 헤더 경로 설정
# 시스템에 설치된 RapidJSON 사용 시도, 없으면 FetchContent 사용
find_path(RAPIDJSON_INCLUDE_DIR
    NAMES rapidjson/rapidjson.h
    PATHS
        /usr/include
        /usr/local/include
        /opt/homebrew/include
        $ENV{RAPIDJSON_INCLUDE_DIR}
)

if(RAPIDJSON_INCLUDE_DIR)
    # 시스템에 설치된 RapidJSON 발견
    message(STATUS "Using system RapidJSON at ${RAPIDJSON_INCLUDE_DIR}")
    set(RAPIDJSON_INCLUDE_PATH ${RAPIDJSON_INCLUDE_DIR})
else()
    # 시스템에 없으면 FetchContent로 다운로드
    message(STATUS "RapidJSON not found in system, using FetchContent")
    include(FetchContent)
    FetchContent_Declare(
        rapidjson
        GIT_REPOSITORY https://github.com/Tencent/rapidjson.git
        GIT_TAG master
    )
    FetchContent_GetProperties(rapidjson)
    if(NOT rapidjson_POPULATED)
        FetchContent_Populate(rapidjson)
    endif()
    set(RAPIDJSON_INCLUDE_PATH ${rapidjson_SOURCE_DIR}/include)
endif()

# 생성된 헤더 파일 경로
set(GENERATED_DIR ${CMAKE_SOURCE_DIR}/../out)

# 기본 테스트 (가장 간단한 테스트 - CI에서 사용)
add_executable(test_basic test_basic.cpp ${GENERATED_DIR}/serializer.cpp ${GENERATED_DIR}/rapidjson_adapter.cpp)
target_include_directories(test_basic PRIVATE
    ${RAPIDJSON_INCLUDE_PATH}
    ${CMAKE_SOURCE_DIR}/..
)

# 테스트 실행 파일 (예제 파일 필요)
add_executable(test_serialization test_serialization.cpp)
target_include_directories(test_serialization PRIVATE
    ${RAPIDJSON_INCLUDE_PATH}
    ${GENERATED_DIR}
)

# 간단한 테스트 (단일 JSON 파일)
add_executable(test_simple test_simple.cpp)
target_include_directories(test_simple PRIVATE
    ${RAPIDJSON_INCLUDE_PATH}
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
