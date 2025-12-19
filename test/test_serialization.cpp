// Test serialization and deserialization of generated C++ classes
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <cmath>
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/filereadstream.h"
#include "rapidjson/filewritestream.h"
#include "rapidjson/prettywriter.h"

// Include all generated headers
#include "order.h"      // Order struct
#include "user.h"       // User struct
#include "product.h"    // Product struct
#include "edgecases.h"  // EdgeCases struct
#include "apiresponse.h" // ApiResponse struct
#include "config.h"     // Config struct

int g_passed = 0;
int g_failed = 0;

void log_test(const std::string& test_name, bool passed) {
    if (passed) {
        std::cout << "[PASS] " << test_name << std::endl;
        g_passed++;
    } else {
        std::cout << "[FAIL] " << test_name << std::endl;
        g_failed++;
    }
}

// Helper function to read JSON file
bool readJsonFile(const std::string& filename, rapidjson::Document& doc) {
    FILE* fp = fopen(filename.c_str(), "rb");
    if (!fp) {
        std::cerr << "Failed to open file: " << filename << std::endl;
        return false;
    }

    char readBuffer[65536];
    rapidjson::FileReadStream is(fp, readBuffer, sizeof(readBuffer));
    doc.ParseStream(is);
    fclose(fp);

    if (doc.HasParseError()) {
        std::cerr << "JSON parse error in file: " << filename << std::endl;
        return false;
    }

    return true;
}

// Helper function to convert rapidjson::Value to string
std::string valueToString(const rapidjson::Value& v) {
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    v.Accept(writer);
    return buffer.GetString();
}

// Helper function to compare two JSON values (relaxed comparison)
bool compareJsonValues(const rapidjson::Value& v1, const rapidjson::Value& v2, bool strict = false) {
    if (v1.GetType() != v2.GetType()) {
        // Allow int/int64 mismatch
        if ((v1.IsInt() || v1.IsInt64()) && (v2.IsInt() || v2.IsInt64())) {
            int64_t val1 = v1.IsInt64() ? v1.GetInt64() : v1.GetInt();
            int64_t val2 = v2.IsInt64() ? v2.GetInt64() : v2.GetInt();
            return val1 == val2;
        }
        // Allow int/double mismatch for numbers
        if ((v1.IsInt() || v1.IsInt64() || v1.IsDouble()) &&
            (v2.IsInt() || v2.IsInt64() || v2.IsDouble())) {
            double val1 = v1.IsDouble() ? v1.GetDouble() : (v1.IsInt64() ? (double)v1.GetInt64() : v1.GetInt());
            double val2 = v2.IsDouble() ? v2.GetDouble() : (v2.IsInt64() ? (double)v2.GetInt64() : v2.GetInt());
            return std::abs(val1 - val2) < 1e-9;
        }
        return false;
    }

    if (v1.IsNull()) {
        return v2.IsNull();
    } else if (v1.IsBool()) {
        return v1.GetBool() == v2.GetBool();
    } else if (v1.IsInt()) {
        return v1.GetInt() == v2.GetInt();
    } else if (v1.IsInt64()) {
        return v1.GetInt64() == v2.GetInt64();
    } else if (v1.IsDouble()) {
        // Use epsilon comparison for floating point
        return std::abs(v1.GetDouble() - v2.GetDouble()) < 1e-6;
    } else if (v1.IsString()) {
        return std::string(v1.GetString()) == std::string(v2.GetString());
    } else if (v1.IsArray()) {
        if (v1.Size() != v2.Size()) {
            return false;
        }
        for (rapidjson::SizeType i = 0; i < v1.Size(); ++i) {
            if (!compareJsonValues(v1[i], v2[i], strict)) {
                return false;
            }
        }
        return true;
    } else if (v1.IsObject()) {
        // Relaxed object comparison - allow extra fields in v2
        for (rapidjson::Value::ConstMemberIterator it = v1.MemberBegin();
             it != v1.MemberEnd(); ++it) {
            if (!v2.HasMember(it->name)) {
                if (strict) {
                    std::cout << "  Missing field in serialized: " << it->name.GetString() << std::endl;
                    return false;
                }
                // Skip null/optional fields
                if (!it->value.IsNull()) {
                    std::cout << "  Missing non-null field: " << it->name.GetString() << std::endl;
                    return false;
                }
            } else {
                if (!compareJsonValues(it->value, v2[it->name], strict)) {
                    std::cout << "  Field mismatch: " << it->name.GetString() << std::endl;
                    return false;
                }
            }
        }
        return true;
    }

    return false;
}

// Test template for any struct that has FromJson and ToJson methods
template<typename T>
bool testSerializationRoundTrip(const std::string& test_name,
                                 const std::string& json_file) {
    std::cout << "\n=== Testing " << test_name << " ===" << std::endl;

    // Read original JSON
    rapidjson::Document original_doc;
    if (!readJsonFile(json_file, original_doc)) {
        log_test(test_name + " - Read JSON", false);
        return false;
    }
    log_test(test_name + " - Read JSON", true);

    // Deserialize from JSON to C++ object
    T obj;
    try {
        obj.FromJson(original_doc);
        log_test(test_name + " - Deserialize (FromJson)", true);
    } catch (const std::exception& e) {
        std::cerr << "Exception during FromJson: " << e.what() << std::endl;
        log_test(test_name + " - Deserialize (FromJson)", false);
        return false;
    }

    // Serialize back to JSON
    rapidjson::Document serialized_doc;
    serialized_doc.SetObject();
    rapidjson::Document::AllocatorType& allocator = serialized_doc.GetAllocator();

    rapidjson::Value serialized_value(rapidjson::kObjectType);
    try {
        obj.ToJson(serialized_value, allocator);
        log_test(test_name + " - Serialize (ToJson)", true);
    } catch (const std::exception& e) {
        std::cerr << "Exception during ToJson: " << e.what() << std::endl;
        log_test(test_name + " - Serialize (ToJson)", false);
        return false;
    }

    // Compare original and serialized JSON
    bool same = compareJsonValues(original_doc, serialized_value, false);
    log_test(test_name + " - Round-trip comparison", same);

    if (!same) {
        std::cout << "Original JSON:" << std::endl;
        std::cout << valueToString(original_doc) << std::endl;
        std::cout << "\nSerialized JSON:" << std::endl;
        std::cout << valueToString(serialized_value) << std::endl;
    }

    return same;
}

// Test basic types
void testBasicTypes() {
    std::cout << "\n### Testing Basic Types ###" << std::endl;

    // Test integer
    {
        rapidjson::Document doc;
        doc.SetObject();
        doc.AddMember("value", 42, doc.GetAllocator());

        int value = 0;
        if (doc.HasMember("value") && doc["value"].IsInt()) {
            value = doc["value"].GetInt();
        }

        rapidjson::Value result(rapidjson::kObjectType);
        result.AddMember("value", value, doc.GetAllocator());

        bool passed = (result["value"].GetInt() == 42);
        log_test("Basic int serialization", passed);
    }

    // Test double
    {
        rapidjson::Document doc;
        doc.SetObject();
        doc.AddMember("value", 3.14159, doc.GetAllocator());

        double value = 0.0;
        if (doc.HasMember("value") && doc["value"].IsDouble()) {
            value = doc["value"].GetDouble();
        }

        rapidjson::Value result(rapidjson::kObjectType);
        result.AddMember("value", value, doc.GetAllocator());

        bool passed = std::abs(result["value"].GetDouble() - 3.14159) < 1e-9;
        log_test("Basic double serialization", passed);
    }

    // Test string
    {
        rapidjson::Document doc;
        doc.SetObject();
        rapidjson::Value str_val;
        str_val.SetString("Hello World", doc.GetAllocator());
        doc.AddMember("value", str_val, doc.GetAllocator());

        std::string value;
        if (doc.HasMember("value") && doc["value"].IsString()) {
            value = doc["value"].GetString();
        }

        rapidjson::Value result(rapidjson::kObjectType);
        rapidjson::Value result_str;
        result_str.SetString(value.c_str(), doc.GetAllocator());
        result.AddMember("value", result_str, doc.GetAllocator());

        bool passed = (std::string(result["value"].GetString()) == "Hello World");
        log_test("Basic string serialization", passed);
    }

    // Test bool
    {
        rapidjson::Document doc;
        doc.SetObject();
        doc.AddMember("value", true, doc.GetAllocator());

        bool value = false;
        if (doc.HasMember("value") && doc["value"].IsBool()) {
            value = doc["value"].GetBool();
        }

        rapidjson::Value result(rapidjson::kObjectType);
        result.AddMember("value", value, doc.GetAllocator());

        bool passed = (result["value"].GetBool() == true);
        log_test("Basic bool serialization", passed);
    }

    // Test array
    {
        rapidjson::Document doc;
        doc.SetObject();
        rapidjson::Value arr(rapidjson::kArrayType);
        arr.PushBack(1, doc.GetAllocator());
        arr.PushBack(2, doc.GetAllocator());
        arr.PushBack(3, doc.GetAllocator());
        doc.AddMember("values", arr, doc.GetAllocator());

        std::vector<int> values;
        if (doc.HasMember("values") && doc["values"].IsArray()) {
            const rapidjson::Value& arr = doc["values"];
            for (rapidjson::SizeType i = 0; i < arr.Size(); ++i) {
                if (arr[i].IsInt()) {
                    values.push_back(arr[i].GetInt());
                }
            }
        }

        rapidjson::Value result(rapidjson::kObjectType);
        rapidjson::Value result_arr(rapidjson::kArrayType);
        for (size_t i = 0; i < values.size(); ++i) {
            result_arr.PushBack(values[i], doc.GetAllocator());
        }
        result.AddMember("values", result_arr, doc.GetAllocator());

        bool passed = (result["values"].Size() == 3 &&
                      result["values"][0].GetInt() == 1 &&
                      result["values"][1].GetInt() == 2 &&
                      result["values"][2].GetInt() == 3);
        log_test("Basic array serialization", passed);
    }
}

int main() {
    std::cout << "========================================" << std::endl;
    std::cout << "  json2cpp Serialization Test Suite" << std::endl;
    std::cout << "========================================" << std::endl;

    // Test basic types first
    testBasicTypes();

    // Test generated structs from JSON examples
    testSerializationRoundTrip<Order>("Order", "../examples/order.json");
    testSerializationRoundTrip<User>("User", "../examples/user.json");
    testSerializationRoundTrip<Product>("Product", "../examples/product.json");
    testSerializationRoundTrip<EdgeCases>("EdgeCases", "../examples/edge-cases.json");
    testSerializationRoundTrip<ApiResponse>("ApiResponse", "../examples/api-response.json");
    testSerializationRoundTrip<Config>("Config", "../examples/config.json");

    std::cout << "\n========================================" << std::endl;
    std::cout << "  Test Summary" << std::endl;
    std::cout << "========================================" << std::endl;
    std::cout << "Passed: " << g_passed << std::endl;
    std::cout << "Failed: " << g_failed << std::endl;
    std::cout << "Total:  " << (g_passed + g_failed) << std::endl;

    if (g_failed == 0) {
        std::cout << "\n✓ All tests passed!" << std::endl;
        return 0;
    } else {
        std::cout << "\n✗ Some tests failed!" << std::endl;
        return 1;
    }
}
