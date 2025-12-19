package types

import (
	"json2cpp/internal/nameutil"
)

type JSONType int

const (
	JSONNull JSONType = iota
	JSONBool
	JSONInt
	JSONFloat
	JSONString
	JSONArray
	JSONObject
)

func (t JSONType) String() string {
	switch t {
	case JSONNull:
		return "null"
	case JSONBool:
		return "bool"
	case JSONInt:
		return "int"
	case JSONFloat:
		return "float"
	case JSONString:
		return "string"
	case JSONArray:
		return "array"
	case JSONObject:
		return "object"
	default:
		return "unknown"
	}
}

func (t JSONType) ToCppType() string {
	switch t {
	case JSONNull:
		return "Optional<null_t>"
	case JSONBool:
		return "bool"
	case JSONInt:
		return "int64_t"
	case JSONFloat:
		return "double"
	case JSONString:
		return "std::string"
	case JSONArray:
		return "std::vector"
	case JSONObject:
		return "struct"
	default:
		return "unknown"
	}
}

type Field struct {
	Name       string
	JSONName   string
	Type       JSONType
	NestedType *Struct  // for object/array (for JSONObject or array of objects)
	ElemType   JSONType // for arrays: element type when not an object
	IsOptional bool
}

type Struct struct {
	Name   string
	Fields []*Field
}

type TypeInfo struct {
	Structs []*Struct
}

func (s *Struct) GetField(name string) *Field {
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (s *Struct) HasField(name string) bool {
	return s.GetField(name) != nil
}

func MergeTypes(types1, types2 []*Struct) []*Struct {
	result := make([]*Struct, 0)
	structMap := make(map[string]*Struct)

	// 첫 번째 타입 집합 추가
	for _, s := range types1 {
		structMap[s.Name] = s
		result = append(result, &Struct{
			Name:   s.Name,
			Fields: append([]*Field{}, s.Fields...),
		})
	}

	// 두 번째 타입 집합 병합
	for _, s2 := range types2 {
		if s1, exists := structMap[s2.Name]; exists {
			// 같은 이름의 struct가 있으면 필드 병합
			mergeStructFields(s1, s2)
		} else {
			// 새로운 struct 추가
			result = append(result, &Struct{
				Name:   s2.Name,
				Fields: append([]*Field{}, s2.Fields...),
			})
		}
	}

	return result
}

func mergeStructFields(s1, s2 *Struct) {
	fieldMap := make(map[string]*Field)
	for _, f := range s1.Fields {
		fieldMap[f.JSONName] = f
	}

	for _, f2 := range s2.Fields {
		if f1, exists := fieldMap[f2.JSONName]; exists {
			// 같은 필드가 있으면 타입 승격
			f1.Type = promoteType(f1.Type, f2.Type)
			if f1.IsOptional || f2.IsOptional {
				f1.IsOptional = true
			}
		} else {
			// 새로운 필드는 optional로 추가
			f2.IsOptional = true
			s1.Fields = append(s1.Fields, f2)
		}
	}
}

func promoteType(t1, t2 JSONType) JSONType {
	// 타입 우선순위: object > array > string > float > int > bool > null
	types := []JSONType{t1, t2}
	for _, t := range []JSONType{JSONObject, JSONArray, JSONString, JSONFloat, JSONInt, JSONBool, JSONNull} {
		for _, tt := range types {
			if t == tt {
				return t
			}
		}
	}
	return t1
}

func GenerateStructName(key string) string {
	// Use centralized sanitizer to produce PascalCase struct/type names.
	return nameutil.SanitizeToCppIdentifier(key, false, true)
}

func SortStructs(structs []*Struct) {
	// 의존성 순서로 정렬 (nested struct가 먼저 오도록)
	graph := make(map[string][]string)
	for _, s := range structs {
		deps := []string{}
		for _, f := range s.Fields {
			if f.Type == JSONObject && f.NestedType != nil {
				deps = append(deps, f.NestedType.Name)
			}
		}
		graph[s.Name] = deps
	}

	// 위상 정렬
	visited := make(map[string]bool)
	result := make([]*Struct, 0)

	var visit func(string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true
		for _, dep := range graph[name] {
			visit(dep)
		}
		// struct 찾아서 추가
		for _, s := range structs {
			if s.Name == name {
				result = append(result, s)
				break
			}
		}
	}

	for _, s := range structs {
		visit(s.Name)
	}

	// 결과를 뒤집어서 의존성이 낮은 것부터
	for i := 0; i < len(result)/2; i++ {
		j := len(result) - 1 - i
		result[i], result[j] = result[j], result[i]
	}

	// structs 슬라이스 업데이트
	for i, s := range result {
		structs[i] = s
	}
}
