// Test with RapidJSON
#include "test_output/types.h"
#include "test_output/serializer.h"
#include "test_output/rapidjson_adapter.h"
#include <iostream>

int main() {
    // Test JSON string
    const char* json = R"({
        "username": "testuser",
        "age": 25,
        "email": "test@example.com",
        "balance": 100.5,
        "is_active": true,
        "roles": ["admin", "user"],
        "scores": [95, 87, 92],
        "id": 12345,
        "profile": {
            "first_name": "Test",
            "last_name": "User",
            "bio": "Test bio",
            "website": "https://example.com"
        },
        "settings": {
            "theme": "dark",
            "notifications": {
                "email": true,
                "push": false,
                "sms": true
            }
        }
    })";

    // Parse with RapidJSON
    rapidjson::Document doc;
    doc.Parse(json);

    if (doc.HasParseError()) {
        std::cerr << "Parse error!" << std::endl;
        return 1;
    }

    // Deserialize using adapter
    json2cpp::RapidJsonReader reader(doc);
    Root root;
    DeserializeRoot(root, reader);

    // Verify
    std::cout << "Username: " << root.username << std::endl;
    std::cout << "Age: " << root.age << std::endl;
    std::cout << "Email: " << root.email << std::endl;
    std::cout << "Balance: " << root.balance << std::endl;
    std::cout << "Active: " << (root.is_active ? "true" : "false") << std::endl;
    std::cout << "Roles: " << root.roles.size() << std::endl;
    std::cout << "Scores: " << root.scores.size() << std::endl;
    std::cout << "Profile name: " << root.profile.first_name << " " << root.profile.last_name << std::endl;
    std::cout << "Theme: " << root.settings.theme << std::endl;

    // Serialize back
    rapidjson::Document outDoc;
    outDoc.SetObject();
    json2cpp::RapidJsonWriter writer(outDoc, outDoc.GetAllocator());
    SerializeRoot(root, writer);

    // Convert to string
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> rapidWriter(buffer);
    outDoc.Accept(rapidWriter);

    std::cout << "\nSerialized JSON:\n" << buffer.GetString() << std::endl;

    std::cout << "\n✓ RapidJSON test passed!" << std::endl;
    return 0;
}
