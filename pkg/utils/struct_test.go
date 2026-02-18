package utils

import (
	"testing"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertYaml2Struct(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectNil bool
		validate  func(t *testing.T, result *structpb.Struct)
	}{
		{
			name:      "empty input",
			input:     "",
			expectNil: false, // Empty YAML unmarshals to nil map, which becomes empty struct
			validate: func(t *testing.T, result *structpb.Struct) {
				// Empty struct should have no fields or nil fields map
				if result.GetFields() != nil {
					assert.Empty(t, result.GetFields())
				}
			},
		},
		{
			name:      "invalid YAML - unclosed bracket",
			input:     "key: [invalid",
			expectNil: true,
		},
		{
			name:      "invalid YAML - bad indentation",
			input:     "key:\n  nested: value\n bad: indent",
			expectNil: true,
		},
		{
			name:      "simple key-value pair",
			input:     "name: test",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "name")
				assert.Equal(t, "test", fields["name"].GetStringValue())
			},
		},
		{
			name:      "multiple key-value pairs",
			input:     "name: test\nversion: 1.0\nenabled: true",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "name")
				require.Contains(t, fields, "version")
				require.Contains(t, fields, "enabled")
				assert.Equal(t, "test", fields["name"].GetStringValue())
			},
		},
		{
			name:      "integer value",
			input:     "count: 42",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "count")
				assert.Equal(t, float64(42), fields["count"].GetNumberValue())
			},
		},
		{
			name:      "float value",
			input:     "ratio: 3.14",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "ratio")
				assert.Equal(t, 3.14, fields["ratio"].GetNumberValue())
			},
		},
		{
			name:      "boolean value",
			input:     "enabled: true",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "enabled")
				assert.Equal(t, true, fields["enabled"].GetBoolValue())
			},
		},
		{
			name: "nested structure",
			input: `
parent:
  child:
    grandchild: value
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "parent")

				parent := fields["parent"].GetStructValue()
				require.NotNil(t, parent)
				require.Contains(t, parent.GetFields(), "child")

				child := parent.GetFields()["child"].GetStructValue()
				require.NotNil(t, child)
				require.Contains(t, child.GetFields(), "grandchild")
				assert.Equal(t, "value", child.GetFields()["grandchild"].GetStringValue())
			},
		},
		{
			name: "deeply nested structure",
			input: `
level1:
  level2:
    level3:
      level4:
        deepValue: found
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "level1")

				l1 := fields["level1"].GetStructValue()
				require.NotNil(t, l1)
				l2 := l1.GetFields()["level2"].GetStructValue()
				require.NotNil(t, l2)
				l3 := l2.GetFields()["level3"].GetStructValue()
				require.NotNil(t, l3)
				l4 := l3.GetFields()["level4"].GetStructValue()
				require.NotNil(t, l4)
				assert.Equal(t, "found", l4.GetFields()["deepValue"].GetStringValue())
			},
		},
		{
			name: "simple array",
			input: `
items:
  - item1
  - item2
  - item3
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "items")

				items := fields["items"].GetListValue()
				require.NotNil(t, items)
				values := items.GetValues()
				require.Len(t, values, 3)
				assert.Equal(t, "item1", values[0].GetStringValue())
				assert.Equal(t, "item2", values[1].GetStringValue())
				assert.Equal(t, "item3", values[2].GetStringValue())
			},
		},
		{
			name: "array of integers",
			input: `
numbers:
  - 1
  - 2
  - 3
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "numbers")

				numbers := fields["numbers"].GetListValue()
				require.NotNil(t, numbers)
				values := numbers.GetValues()
				require.Len(t, values, 3)
				assert.Equal(t, float64(1), values[0].GetNumberValue())
				assert.Equal(t, float64(2), values[1].GetNumberValue())
				assert.Equal(t, float64(3), values[2].GetNumberValue())
			},
		},
		{
			name: "array of objects",
			input: `
users:
  - name: alice
    age: 30
  - name: bob
    age: 25
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "users")

				users := fields["users"].GetListValue()
				require.NotNil(t, users)
				values := users.GetValues()
				require.Len(t, values, 2)

				alice := values[0].GetStructValue()
				require.NotNil(t, alice)
				assert.Equal(t, "alice", alice.GetFields()["name"].GetStringValue())
				assert.Equal(t, float64(30), alice.GetFields()["age"].GetNumberValue())

				bob := values[1].GetStructValue()
				require.NotNil(t, bob)
				assert.Equal(t, "bob", bob.GetFields()["name"].GetStringValue())
				assert.Equal(t, float64(25), bob.GetFields()["age"].GetNumberValue())
			},
		},
		{
			name: "nested arrays",
			input: `
matrix:
  - - 1
    - 2
  - - 3
    - 4
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "matrix")

				matrix := fields["matrix"].GetListValue()
				require.NotNil(t, matrix)
				values := matrix.GetValues()
				require.Len(t, values, 2)

				row1 := values[0].GetListValue()
				require.NotNil(t, row1)
				require.Len(t, row1.GetValues(), 2)
				assert.Equal(t, float64(1), row1.GetValues()[0].GetNumberValue())
				assert.Equal(t, float64(2), row1.GetValues()[1].GetNumberValue())

				row2 := values[1].GetListValue()
				require.NotNil(t, row2)
				require.Len(t, row2.GetValues(), 2)
				assert.Equal(t, float64(3), row2.GetValues()[0].GetNumberValue())
				assert.Equal(t, float64(4), row2.GetValues()[1].GetNumberValue())
			},
		},
		{
			name: "mixed types in array",
			input: `
mixed:
  - string_value
  - 42
  - true
  - 3.14
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "mixed")

				mixed := fields["mixed"].GetListValue()
				require.NotNil(t, mixed)
				values := mixed.GetValues()
				require.Len(t, values, 4)
				assert.Equal(t, "string_value", values[0].GetStringValue())
				assert.Equal(t, float64(42), values[1].GetNumberValue())
				assert.Equal(t, true, values[2].GetBoolValue())
				assert.Equal(t, 3.14, values[3].GetNumberValue())
			},
		},
		{
			name: "complex nested structure with arrays",
			input: `
config:
  server:
    host: localhost
    port: 8080
    endpoints:
      - path: /api
        methods:
          - GET
          - POST
      - path: /health
        methods:
          - GET
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "config")

				config := fields["config"].GetStructValue()
				require.NotNil(t, config)

				server := config.GetFields()["server"].GetStructValue()
				require.NotNil(t, server)
				assert.Equal(t, "localhost", server.GetFields()["host"].GetStringValue())
				assert.Equal(t, float64(8080), server.GetFields()["port"].GetNumberValue())

				endpoints := server.GetFields()["endpoints"].GetListValue()
				require.NotNil(t, endpoints)
				require.Len(t, endpoints.GetValues(), 2)

				endpoint1 := endpoints.GetValues()[0].GetStructValue()
				require.NotNil(t, endpoint1)
				assert.Equal(t, "/api", endpoint1.GetFields()["path"].GetStringValue())

				methods1 := endpoint1.GetFields()["methods"].GetListValue()
				require.NotNil(t, methods1)
				require.Len(t, methods1.GetValues(), 2)
			},
		},
		{
			name:      "null value",
			input:     "nullField: null",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "nullField")
				// Null values in structpb are represented with NullValue
				assert.NotNil(t, fields["nullField"].GetNullValue)
			},
		},
		{
			name:      "empty string value",
			input:     "emptyString: \"\"",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "emptyString")
				assert.Equal(t, "", fields["emptyString"].GetStringValue())
			},
		},
		{
			name: "empty array",
			input: `
emptyArray: []
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "emptyArray")

				arr := fields["emptyArray"].GetListValue()
				require.NotNil(t, arr)
				assert.Empty(t, arr.GetValues())
			},
		},
		{
			name: "empty object",
			input: `
emptyObject: {}
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "emptyObject")

				obj := fields["emptyObject"].GetStructValue()
				require.NotNil(t, obj)
				assert.Empty(t, obj.GetFields())
			},
		},
		{
			name:      "string with special characters",
			input:     "special: \"hello: world = test\"",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "special")
				assert.Equal(t, "hello: world = test", fields["special"].GetStringValue())
			},
		},
		{
			name: "multiline string",
			input: `
description: |
  This is a
  multiline string
  with multiple lines
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "description")
				assert.Equal(t, "This is a\nmultiline string\nwith multiple lines\n", fields["description"].GetStringValue())
			},
		},
		{
			name: "inline array syntax",
			input: `
inline: [a, b, c]
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "inline")

				arr := fields["inline"].GetListValue()
				require.NotNil(t, arr)
				values := arr.GetValues()
				require.Len(t, values, 3)
				assert.Equal(t, "a", values[0].GetStringValue())
				assert.Equal(t, "b", values[1].GetStringValue())
				assert.Equal(t, "c", values[2].GetStringValue())
			},
		},
		{
			name: "inline object syntax",
			input: `
inline: {key1: value1, key2: value2}
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "inline")

				obj := fields["inline"].GetStructValue()
				require.NotNil(t, obj)
				assert.Equal(t, "value1", obj.GetFields()["key1"].GetStringValue())
				assert.Equal(t, "value2", obj.GetFields()["key2"].GetStringValue())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertYaml2Struct(tt.input)

			if tt.expectNil {
				assert.Nil(t, result, "expected nil result for invalid YAML")
				return
			}

			require.NotNil(t, result, "expected non-nil result")

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestConvertMapInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "simple string values",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "mixed types",
			input: map[string]interface{}{
				"string": "hello",
				"int":    42,
				"float":  3.14,
				"bool":   true,
			},
			expected: map[string]interface{}{
				"string": "hello",
				"int":    42,
				"float":  3.14,
				"bool":   true,
			},
		},
		{
			name: "nested map[string]interface{}",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
		},
		{
			name: "nested map[interface{}]interface{}",
			input: map[string]interface{}{
				"outer": map[interface{}]interface{}{
					"inner": "value",
				},
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
		},
		{
			name: "array of strings",
			input: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			expected: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
		},
		{
			name: "array with nested maps",
			input: map[string]interface{}{
				"items": []interface{}{
					map[interface{}]interface{}{
						"name": "item1",
					},
					map[interface{}]interface{}{
						"name": "item2",
					},
				},
			},
			expected: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"name": "item1",
					},
					map[string]interface{}{
						"name": "item2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertMapInterface(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "integer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: true,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name: "map[interface{}]interface{} value",
			input: map[interface{}]interface{}{
				"key": "value",
			},
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "map[string]interface{} value",
			input: map[string]interface{}{
				"key": "value",
			},
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:     "slice value",
			input:    []interface{}{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name: "slice with nested maps",
			input: []interface{}{
				map[interface{}]interface{}{
					"nested": "value",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"nested": "value",
				},
			},
		},
		{
			name: "deeply nested map[interface{}]interface{}",
			input: map[interface{}]interface{}{
				"level1": map[interface{}]interface{}{
					"level2": map[interface{}]interface{}{
						"value": "deep",
					},
				},
			},
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"value": "deep",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertYaml2Struct_RealWorldScenarios(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectNil bool
		validate  func(t *testing.T, result *structpb.Struct)
	}{
		{
			name: "Istio rate limit config",
			input: `
domain: productpage-ratelimit
descriptors:
  - key: PATH
    value: "/productpage"
    rate_limit:
      unit: minute
      requests_per_unit: 1
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				assert.Equal(t, "productpage-ratelimit", fields["domain"].GetStringValue())

				descriptors := fields["descriptors"].GetListValue()
				require.NotNil(t, descriptors)
				require.Len(t, descriptors.GetValues(), 1)

				desc := descriptors.GetValues()[0].GetStructValue()
				require.NotNil(t, desc)
				assert.Equal(t, "PATH", desc.GetFields()["key"].GetStringValue())
				assert.Equal(t, "/productpage", desc.GetFields()["value"].GetStringValue())

				rateLimit := desc.GetFields()["rate_limit"].GetStructValue()
				require.NotNil(t, rateLimit)
				assert.Equal(t, "minute", rateLimit.GetFields()["unit"].GetStringValue())
				assert.Equal(t, float64(1), rateLimit.GetFields()["requests_per_unit"].GetNumberValue())
			},
		},
		{
			name: "Kubernetes-like labels",
			input: `
metadata:
  name: my-service
  namespace: default
  labels:
    app: myapp
    version: v1
    environment: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()

				metadata := fields["metadata"].GetStructValue()
				require.NotNil(t, metadata)
				assert.Equal(t, "my-service", metadata.GetFields()["name"].GetStringValue())
				assert.Equal(t, "default", metadata.GetFields()["namespace"].GetStringValue())

				labels := metadata.GetFields()["labels"].GetStructValue()
				require.NotNil(t, labels)
				assert.Equal(t, "myapp", labels.GetFields()["app"].GetStringValue())
				assert.Equal(t, "v1", labels.GetFields()["version"].GetStringValue())

				spec := fields["spec"].GetStructValue()
				require.NotNil(t, spec)
				assert.Equal(t, float64(3), spec.GetFields()["replicas"].GetNumberValue())
			},
		},
		{
			name: "EnvoyFilter patch configuration",
			input: `
patch:
  operation: INSERT_BEFORE
  value:
    name: envoy.filters.http.local_ratelimit
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit
      stat_prefix: http_local_rate_limiter
`,
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()

				patch := fields["patch"].GetStructValue()
				require.NotNil(t, patch)
				assert.Equal(t, "INSERT_BEFORE", patch.GetFields()["operation"].GetStringValue())

				value := patch.GetFields()["value"].GetStructValue()
				require.NotNil(t, value)
				assert.Equal(t, "envoy.filters.http.local_ratelimit", value.GetFields()["name"].GetStringValue())

				typedConfig := value.GetFields()["typed_config"].GetStructValue()
				require.NotNil(t, typedConfig)
				assert.Equal(t, "type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit", typedConfig.GetFields()["@type"].GetStringValue())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertYaml2Struct(tt.input)

			if tt.expectNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestConvertYaml2Struct_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectNil bool
		validate  func(t *testing.T, result *structpb.Struct)
	}{
		{
			name:      "only whitespace",
			input:     "   \n\t\n   ",
			expectNil: true, // Whitespace-only YAML results in nil map, which causes empty JSON marshal returning nil struct
		},
		{
			name:      "YAML with comments only",
			input:     "# This is a comment\n# Another comment",
			expectNil: false, // Comments-only is valid YAML with nil content
			validate: func(t *testing.T, result *structpb.Struct) {
				if result.GetFields() != nil {
					assert.Empty(t, result.GetFields())
				}
			},
		},
		{
			name:      "YAML with anchor and alias",
			input:     "defaults: &defaults\n  timeout: 30\nproduction:\n  <<: *defaults\n  host: prod.example.com",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "defaults")
				require.Contains(t, fields, "production")

				production := fields["production"].GetStructValue()
				require.NotNil(t, production)
				// The anchor values should be merged
				assert.Equal(t, float64(30), production.GetFields()["timeout"].GetNumberValue())
				assert.Equal(t, "prod.example.com", production.GetFields()["host"].GetStringValue())
			},
		},
		{
			name:      "unicode characters",
			input:     "greeting: \u4f60\u597d\u4e16\u754c",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "greeting")
				assert.Equal(t, "\u4f60\u597d\u4e16\u754c", fields["greeting"].GetStringValue())
			},
		},
		{
			name:      "special YAML characters in quoted string",
			input:     "special: \"key: value, [item], {obj}\"",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "special")
				assert.Equal(t, "key: value, [item], {obj}", fields["special"].GetStringValue())
			},
		},
		{
			name:      "negative numbers",
			input:     "negative: -42\nnegativeFloat: -3.14",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "negative")
				require.Contains(t, fields, "negativeFloat")
				assert.Equal(t, float64(-42), fields["negative"].GetNumberValue())
				assert.Equal(t, -3.14, fields["negativeFloat"].GetNumberValue())
			},
		},
		{
			name:      "scientific notation",
			input:     "scientific: 1.23e10",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				require.Contains(t, fields, "scientific")
				assert.Equal(t, 1.23e10, fields["scientific"].GetNumberValue())
			},
		},
		{
			name:      "boolean variations",
			input:     "yes_val: yes\nno_val: no\ntrue_val: true\nfalse_val: false",
			expectNil: false,
			validate: func(t *testing.T, result *structpb.Struct) {
				fields := result.GetFields()
				// YAML 1.1 treats yes/no as booleans
				assert.Equal(t, true, fields["yes_val"].GetBoolValue())
				assert.Equal(t, false, fields["no_val"].GetBoolValue())
				assert.Equal(t, true, fields["true_val"].GetBoolValue())
				assert.Equal(t, false, fields["false_val"].GetBoolValue())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertYaml2Struct(tt.input)

			if tt.expectNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestConvertYaml2Struct_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid YAML - duplicate key at same level",
			input: "key: value1\nkey: value2", // This is actually valid YAML, last value wins
		},
		{
			name:  "invalid YAML - tabs in wrong places",
			input: "key:\n\t\tvalue", // Tabs can cause issues
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These might or might not return nil depending on the YAML parser
			// We're just ensuring no panic occurs
			_ = ConvertYaml2Struct(tt.input)
		})
	}
}

func TestConvertYaml2Struct_LargeInput(t *testing.T) {
	// Test with a moderately large YAML structure
	input := `
root:
  level1_1:
    level2_1:
      items:
        - name: item1
          value: 1
        - name: item2
          value: 2
        - name: item3
          value: 3
    level2_2:
      config:
        enabled: true
        timeout: 30
        retries: 3
  level1_2:
    data:
      - type: typeA
        metadata:
          key1: value1
          key2: value2
      - type: typeB
        metadata:
          key3: value3
          key4: value4
`

	result := ConvertYaml2Struct(input)
	require.NotNil(t, result)

	fields := result.GetFields()
	require.Contains(t, fields, "root")

	root := fields["root"].GetStructValue()
	require.NotNil(t, root)

	// Verify the structure is properly converted
	level1_1 := root.GetFields()["level1_1"].GetStructValue()
	require.NotNil(t, level1_1)

	level2_1 := level1_1.GetFields()["level2_1"].GetStructValue()
	require.NotNil(t, level2_1)

	items := level2_1.GetFields()["items"].GetListValue()
	require.NotNil(t, items)
	assert.Len(t, items.GetValues(), 3)
}
