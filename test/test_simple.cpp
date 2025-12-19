// Simple test for JSON serialization/deserialization
#include <iostream>
#include <fstream>
#include <string>
#include <cmath>
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/filereadstream.h"

#include "order.h"

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

std::string valueToString(const rapidjson::Value& v) {
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    v.Accept(writer);
    return buffer.GetString();
}

int main() {
    std::cout << "========================================" << std::endl;
    std::cout << "  json2cpp Serialization Test" << std::endl;
    std::cout << "========================================" << std::endl;

    // Test Order struct
    {
        std::cout << "\n=== Testing Order ===" << std::endl;

        rapidjson::Document original_doc;
        if (!readJsonFile("../../examples/order.json", original_doc)) {
            log_test("Order - Read JSON", false);
            return 1;
        }
        log_test("Order - Read JSON", true);

        Order order;
        order.FromJson(original_doc);
        log_test("Order - Deserialize (FromJson)", true);

        // Check some values
        bool values_ok = (order.symbol == "AAPL" &&
                         order.quantity == 100 &&
                         std::abs(order.price - 150.75) < 0.01);
        log_test("Order - Data validation", values_ok);

        rapidjson::Document serialized_doc;
        serialized_doc.SetObject();
        rapidjson::Document::AllocatorType& allocator = serialized_doc.GetAllocator();

        rapidjson::Value serialized_value(rapidjson::kObjectType);
        order.ToJson(serialized_value, allocator);
        log_test("Order - Serialize (ToJson)", true);

        std::cout << "\nOriginal JSON:" << std::endl;
        std::cout << valueToString(original_doc) << std::endl;
        std::cout << "\nSerialized JSON:" << std::endl;
        std::cout << valueToString(serialized_value) << std::endl;
    }

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
