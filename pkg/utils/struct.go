package utils

import (
	"encoding/json"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"gopkg.in/yaml.v2"
)

func ConvertYaml2Struct(str string) *structpb.Struct {
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(str), &data); err != nil {
		return nil
	}

	// Convert to JSON then to structpb.Struct
	jsonBytes, err := json.Marshal(convertMapInterface(data))
	if err != nil {
		return nil
	}

	result := &structpb.Struct{}
	if err := json.Unmarshal(jsonBytes, result); err != nil {
		return nil
	}

	return result
}

// convertMapInterface converts map[interface{}]interface{} to map[string]interface{}
// which is needed for JSON marshaling
func convertMapInterface(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = convertValue(v)
	}
	return result
}

func convertValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			keyStr, ok := k.(string)
			if !ok {
				continue
			}
			m[keyStr] = convertValue(v)
		}
		return m
	case map[string]interface{}:
		return convertMapInterface(val)
	case []interface{}:
		arr := make([]interface{}, len(val))
		for i, item := range val {
			arr[i] = convertValue(item)
		}
		return arr
	default:
		return v
	}
}
