package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"json2cpp/internal/nameutil"
	"json2cpp/internal/types"
)

type Parser struct {
	structCounter int
	legacyCpp     bool
	camelCase     bool
}

func NewParser(legacyCpp bool, camelCase bool) *Parser {
	return &Parser{
		legacyCpp: legacyCpp,
		camelCase: camelCase,
	}
}

func (p *Parser) ParseFile(filename string) ([]*types.Struct, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return p.ParseValue(v, "Root")
}

func (p *Parser) ParseValue(v interface{}, suggestedName string) ([]*types.Struct, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		return p.parseObject(val, suggestedName)
	case []interface{}:
		// 배열의 경우 배열 남용을 위해 첫 번째 요소 분석
		if len(val) > 0 {
			return p.ParseValue(val[0], suggestedName+"Item")
		}
		return []*types.Struct{}, nil
	default:
		// 기본 타입은 struct가 아님
		return []*types.Struct{}, nil
	}
}

func (p *Parser) parseObject(obj map[string]interface{}, structName string) ([]*types.Struct, error) {
	structs := make([]*types.Struct, 0)

	// 현재 struct 생성
	current := &types.Struct{
		Name:   p.generateStructName(structName),
		Fields: make([]*types.Field, 0),
	}

	for key, value := range obj {
		field := &types.Field{
			Name:     p.generateFieldName(key),
			JSONName: key,
		}

		switch val := value.(type) {
		case nil:
			field.Type = types.JSONNull
			field.IsOptional = true

		case bool:
			field.Type = types.JSONBool

		case float64:
			if isInteger(val) {
				field.Type = types.JSONInt
			} else {
				field.Type = types.JSONFloat
			}

		case string:
			field.Type = types.JSONString

		case []interface{}:
			field.Type = types.JSONArray
			if len(val) > 0 {
				// 배열 요소의 타입 분석
				elemType := p.inferArrayElementType(val)
				if elemType == types.JSONObject {
					// 객체 배열인 경우 nested struct 생성
					nestedName := p.generateStructName(key) + "Item"
					nestedStructs, err := p.parseArrayOfObjects(val, nestedName)
					if err != nil {
						return nil, err
					}
					structs = append(structs, nestedStructs...)
					if len(nestedStructs) > 0 {
						field.NestedType = nestedStructs[0]
					}
				} else {
					// primitive array element type
					field.ElemType = elemType
				}
			}

		case map[string]interface{}:
			field.Type = types.JSONObject
			nestedName := p.generateStructName(key)
			nestedStructs, err := p.parseObject(val, nestedName)
			if err != nil {
				return nil, err
			}
			structs = append(structs, nestedStructs...)
			if len(nestedStructs) > 0 {
				field.NestedType = nestedStructs[len(nestedStructs)-1]
			}

		default:
			field.Type = types.JSONString // 기본값
		}

		current.Fields = append(current.Fields, field)
	}

	structs = append(structs, current)
	return structs, nil
}

func (p *Parser) parseArrayOfObjects(arr []interface{}, structName string) ([]*types.Struct, error) {
	if len(arr) == 0 {
		return nil, nil
	}

	// 첫 번째 객체로 struct 생성
	first, ok := arr[0].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return p.parseObject(first, structName)
}

func (p *Parser) inferArrayElementType(arr []interface{}) types.JSONType {
	if len(arr) == 0 {
		return types.JSONNull
	}

	// 모든 요소의 타입을 분석
	typeCounts := make(map[types.JSONType]int)
	for _, elem := range arr {
		if elem == nil {
			typeCounts[types.JSONNull]++
		} else if _, ok := elem.(bool); ok {
			typeCounts[types.JSONBool]++
		} else if f, ok := elem.(float64); ok {
			if isInteger(f) {
				typeCounts[types.JSONInt]++
			} else {
				typeCounts[types.JSONFloat]++
			}
		} else if _, ok := elem.(string); ok {
			typeCounts[types.JSONString]++
		} else if _, ok := elem.([]interface{}); ok {
			typeCounts[types.JSONArray]++
		} else if _, ok := elem.(map[string]interface{}); ok {
			typeCounts[types.JSONObject]++
		}
	}

	// 가장 많은 타입 선택
	maxCount := 0
	var result types.JSONType
	for t, count := range typeCounts {
		if count > maxCount {
			maxCount = count
			result = t
		}
	}

	return result
}

func isInteger(f float64) bool {
	return f == float64(int64(f))
}

func (p *Parser) generateStructName(key string) string {
	p.structCounter++
	return types.GenerateStructName(key)
}

func (p *Parser) generateFieldName(key string) string {
	// Delegate sanitization to centralized helper.
	return nameutil.SanitizeToCppIdentifier(key, p.camelCase, false)
}

// removed local keyword check — centralized in nameutil
