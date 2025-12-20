// Basic test for json2cpp generated code
#include <iostream>
#include <cmath>
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"

// Include generated files from ../out
#include "../out/types.h"
#include "../out/serializer.h"
#include "../out/rapidjson_adapter.h"

int main() {
    std::cout << "========================================" << std::endl;
    std::cout << "  json2cpp Basic Test" << std::endl;
    std::cout << "========================================" << std::endl;

    // Create a simple JSON document matching test_simple.json structure
    const char* json = R"({
        "name": "test",
        "value": 123,
        "active": true
    })";

    // Parse JSON
    rapidjson::Document doc;
    doc.Parse(json);

    if (doc.HasParseError()) {
        std::cerr << "Parse error!" << std::endl;
        return 1;
    }
    std::cout << "[PASS] JSON parsing" << std::endl;

    // Try to deserialize
    json2cpp::RapidJsonReader reader(doc);
    Root root;
    try {
        DeserializeRoot(root, reader);
        std::cout << "[PASS] Deserialization" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "[FAIL] Deserialization: " << e.what() << std::endl;
        return 1;
    }

    // Serialize back
    rapidjson::Document outDoc;
    outDoc.SetObject();
    json2cpp::RapidJsonWriter writer(outDoc, outDoc.GetAllocator());
    try {
        SerializeRoot(root, writer);
        std::cout << "[PASS] Serialization" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "[FAIL] Serialization: " << e.what() << std::endl;
        return 1;
    }

    // Convert to string for display
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> rapidWriter(buffer);
    outDoc.Accept(rapidWriter);

    std::cout << "\nSerialized JSON:\n" << buffer.GetString() << std::endl;
    std::cout << "\n✓ All basic tests passed!" << std::endl;

    return 0;
}
