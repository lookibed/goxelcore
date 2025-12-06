package coders

import (
	"encoding/json"
	"fmt"
	"io/ioutil" // For ReadFile
	"strings"

	"goxelcore/data" // For dv.Value
	"goxelcore/io"   // For io.Path
)

// ParseJSON parses a JSON string into a dv.Value.
// Corresponds to C++ json::parse(std::string_view filename, std::string_view source)
// and json::parse(std::string_view source).
func ParseJSON(filename string, source string) (data.Value, error) {
	var raw interface{}
	decoder := json.NewDecoder(strings.NewReader(source))
	// Use Number instead of float64 for numbers
	decoder.UseNumber() 

	if err := decoder.Decode(&raw); err != nil {
		return data.Value{}, &ParsingError{
			Filename: filename,
			Message:  fmt.Sprintf("JSON parsing failed: %v", err),
		}
	}
	return convertToDvValue(raw), nil
}

// ReadJSONFile reads a JSON file and parses it into a dv.Value.
// Not directly in C++ json::, but useful for AssetsLoader.
func ReadJSONFile(p io.Path) (data.Value, error) {
	content, err := io.ReadString(p)
	if err != nil {
		return data.Value{}, fmt.Errorf("failed to read JSON file %s: %w", p.String(), err)
	}
	return ParseJSON(p.String(), content), nil
}

// convertToDvValue recursively converts an interface{} (from json.Unmarshal) to dv.Value.
func convertToDvValue(raw interface{}) data.Value {
	switch v := raw.(type) {
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return data.NewValue(i)
		}
		if f, err := v.Float64(); err == nil {
			return data.NewValue(f)
		}
		return data.NewValue(v.String()) // Fallback to string if number parsing fails
	case bool:
		return data.NewValue(v)
	case string:
		return data.NewValue(v)
	case []interface{}:
		list := make([]data.Value, len(v))
		for i, item := range v {
			list[i] = convertToDvValue(item)
		}
		return data.NewValue(list)
	case map[string]interface{}:
		obj := make(map[string]data.Value)
		for key, item := range v {
			obj[key] = convertToDvValue(item)
		}
		return data.NewValue(obj)
	case nil:
		return data.Value{} // Represents ValueTypeNone
	default:
		// Should not happen with typical JSON types
		log.Printf("Warning: Unhandled type during JSON to dv.Value conversion: %T\n", v)
		return data.Value{}
	}
}

// Stringify converts a dv.Value to a JSON string.
// Corresponds to C++ json::stringify(const dv::value& value, bool nice, const std::string& indent, bool escapeUtf8).
func Stringify(val data.Value, nice bool, indent string, escapeUtf8 bool) (string, error) {
	// Go's json marshaler handles UTF-8 by default and doesn't have direct control over escaping
	// Indent is handled by json.MarshalIndent
	
	raw := convertDvValueToRaw(val)
	var bytes []byte
	var err error
	if nice {
		bytes, err = json.MarshalIndent(raw, "", indent)
	} else {
		bytes, err = json.Marshal(raw)
	}
	if err != nil {
		return "", fmt.Errorf("failed to stringify dv.Value to JSON: %w", err)
	}
	return string(bytes), nil
}

// convertDvValueToRaw recursively converts a dv.Value to a raw interface{} for json.Marshal.
func convertDvValueToRaw(val data.Value) interface{} {
	switch val.Type {
	case data.ValueTypeNone:
		return nil
	case data.ValueTypeNumber:
		return val.Value.(float64)
	case data.ValueTypeInteger:
		return val.Value.(int64)
	case data.ValueTypeBoolean:
		return val.Value.(bool)
	case data.ValueTypeString:
		return val.Value.(string)
	case data.ValueTypeBytes:
		return val.Value.([]byte)
	case data.ValueTypeList:
		list := val.Value.([]data.Value)
		rawList := make([]interface{}, len(list))
		for i, item := range list {
			rawList[i] = convertDvValueToRaw(item)
		}
		return rawList
	case data.ValueTypeObject:
		obj := val.Value.(map[string]data.Value)
		rawMap := make(map[string]interface{})
		for key, item := range obj {
			rawMap[key] = convertDvValueToRaw(item)
		}
		return rawMap
	default:
		return nil
	}
}