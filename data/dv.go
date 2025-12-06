package data

import (
	"fmt"
	"log"
	"math"
	"reflect"

	"goxelcore" // For size_t (ubyte, integer_t, number_t)
)

// value_type corresponds to C++ dv::value_type enum.
type ValueType uint8

const (
	ValueTypeNone   ValueType = iota
	ValueTypeNumber           // float64
	ValueTypeBoolean          // bool
	ValueTypeInteger          // int64
	ValueTypeObject           // map[string]Value
	ValueTypeList             // []Value
	ValueTypeBytes            // []byte
	ValueTypeString           // string
)

// TypeName returns the string representation of a ValueType.
// Corresponds to C++ dv::type_name(value_type type)
func (vt ValueType) TypeName() string {
	switch vt {
	case ValueTypeNone:
		return "none"
	case ValueTypeNumber:
		return "number"
	case ValueTypeBoolean:
		return "boolean"
	case ValueTypeInteger:
		return "integer"
	case ValueTypeObject:
		return "object"
	case ValueTypeList:
		return "list"
	case ValueTypeBytes:
		return "bytes"
	case ValueTypeString:
		return "string"
	default:
		return "unknown"
	}
}

// Value corresponds to C++ dv::value class.
// It's a dynamic value type capable of holding various data types.
type Value struct {
	Type  ValueType
	Value interface{} // Holds the actual data
}

// NewValue creates a new Value from a given Go type.
// This acts as a constructor.
func NewValue(v interface{}) Value {
	val := Value{}
	val.Set(v)
	return val
}

// Set sets the Value's content based on the Go type of 'v'.
// Corresponds to C++ value::operator= overloads.
func (v *Value) Set(data interface{}) {
	if data == nil {
		v.Type = ValueTypeNone
		v.Value = nil
		return
	}

	switch val := data.(type) {
	case int, int8, int16, int32, int64:
		v.Type = ValueTypeInteger
		v.Value = reflect.ValueOf(val).Int()
	case uint, uint8, uint16, uint32, uint64: // Convert unsigned to signed int64 if possible
		if reflect.ValueOf(val).Uint() > math.MaxInt64 { 
			v.Type = ValueTypeNumber
			v.Value = float64(reflect.ValueOf(val).Uint())
		} else {
			v.Type = ValueTypeInteger
			v.Value = int64(reflect.ValueOf(val).Uint())
		}
	case float32, float64:
		v.Type = ValueTypeNumber
		v.Value = reflect.ValueOf(val).Float()
	case bool:
		v.Type = ValueTypeBoolean
		v.Value = val
	case string:
		v.Type = ValueTypeString
		v.Value = val
	case []byte:
		v.Type = ValueTypeBytes
		v.Value = val
	case []Value:
		v.Type = ValueTypeList
		v.Value = val
	case map[string]Value:
		v.Type = ValueTypeObject
		v.Value = val
	case Value: // If setting from another Value, copy it.
		*v = val
	default:
		// Fallback for unhandled types, or error.
		log.Printf("Warning: Unhandled type for dv.Value.Set: %T\n", data)
		v.Type = ValueTypeNone
		v.Value = nil
	}
}

// AsString returns the value as a string. Panics if type is not string.
// Corresponds to C++ value::asString()
func (v Value) AsString() string {
	if v.Type != ValueTypeString {
		throwTypeError(v.Type, ValueTypeString)
	}
	return v.Value.(string)
}

// AsInteger returns the value as an int64. Panics if type is not integer or number.
// Corresponds to C++ value::asInteger()
func (v Value) AsInteger() int64 {
	switch v.Type {
	case ValueTypeInteger:
		return v.Value.(int64)
	case ValueTypeNumber:
		return int64(v.Value.(float64))
	default:
		throwTypeError(v.Type, ValueTypeInteger)
	}
	return 0
}

// AsNumber returns the value as a float64. Panics if type is not number or integer.
// Corresponds to C++ value::asNumber()
func (v Value) AsNumber() float64 {
	switch v.Type {
	case ValueTypeNumber:
		return v.Value.(float64)
	case ValueTypeInteger:
		return float64(v.Value.(int64))
	default:
		throwTypeError(v.Type, ValueTypeNumber)
	}
	return 0.0
}

// AsBoolean returns the value as a bool. Panics if type is not boolean.
// Corresponds to C++ value::asBoolean()
func (v Value) AsBoolean() bool {
	if v.Type != ValueTypeBoolean {
		throwTypeError(v.Type, ValueTypeBoolean)
	}
	return v.Value.(bool)
}

// AsBytes returns the value as a []byte. Panics if type is not bytes.
// Corresponds to C++ value::asBytes()
func (v Value) AsBytes() []byte {
	if v.Type != ValueTypeBytes {
		throwTypeError(v.Type, ValueTypeBytes)
	}
	return v.Value.([]byte)
}

// AsObject returns the value as a map[string]Value. Panics if type is not object.
// Corresponds to C++ value::asObject()
func (v Value) AsObject() map[string]Value {
	if v.Type != ValueTypeObject {
		throwTypeError(v.Type, ValueTypeObject)
	}
	return v.Value.(map[string]Value)
}

// AsList returns the value as a []Value. Panics if type is not list.
// Corresponds to C++ value::asList()
func (v Value) AsList() []Value {
	if v.Type != ValueTypeList {
		throwTypeError(v.Type, ValueTypeList)
	}
	return v.Value.([]Value)
}

// Add adds a value to a list or object. Panics if not list/object.
// Corresponds to C++ value::add()
func (v *Value) Add(val Value) {
	switch v.Type {
	case ValueTypeList:
		v.Value = append(v.Value.([]Value), val)
	case ValueTypeObject:
		// Not directly supported in C++ add(value v), usually add(key, value)
		// For Go, this method is typically for lists.
		panic(fmt.Sprintf("cannot add Value to object without key (type %s)", v.Type.TypeName()))
	default:
		panic(fmt.Sprintf("cannot add to non-list/object type (type %s)", v.Type.TypeName()))
	}
}

// AddToKey adds a value to an object with a key. (Go specific helper)
func (v *Value) AddToKey(key string, val Value) {
	if v.Type != ValueTypeObject {
		panic(fmt.Sprintf("cannot add to non-object type (type %s)", v.Type.TypeName()))
	}
	v.Value.(map[string]Value)[key] = val
}

// GetAtKey returns the value at a given key for an object.
// Corresponds to C++ value::operator[](const key_t& key)
func (v Value) GetAtKey(key string) Value {
	if v.Type != ValueTypeObject {
		throwTypeError(v.Type, ValueTypeObject)
	}
	if obj, ok := v.Value.(map[string]Value); ok {
		if val, found := obj[key]; found {
			return val
		}
	}
	return Value{} // Return none type if not found
}

// GetAtIndex returns the value at a given index for a list.
// Corresponds to C++ value::operator[](size_t index)
func (v Value) GetAtIndex(index int) Value {
	if v.Type != ValueTypeList {
		throwTypeError(v.Type, ValueTypeList)
	}
	if list, ok := v.Value.([]Value); ok {
		if index >= 0 && index < len(list) {
			return list[index]
		}
	}
	return Value{} // Return none type if out of bounds
}

// Has checks if an object has a specific key.
// Corresponds to C++ value::has()
func (v Value) Has(key string) bool {
	if v.Type != ValueTypeObject {
		return false
	}
	if obj, ok := v.Value.(map[string]Value); ok {
		_, found := obj[key]
		return found
	}
	return false
}

// Size returns the size of a list or object.
// Corresponds to C++ value::size()
func (v Value) Size() goxelcore.size_t {
	switch v.Type {
	case ValueTypeList:
		return goxelcore.size_t(len(v.Value.([]Value)))
	case ValueTypeObject:
		return goxelcore.size_t(len(v.Value.(map[string]Value)))
	case ValueTypeString:
		return goxelcore.size_t(len(v.Value.(string)))
	case ValueTypeBytes:
		return goxelcore.size_t(len(v.Value.([]byte)))
	default:
		return 0
	}
}

// IsEmpty checks if the value is empty.
// Corresponds to C++ value::empty()
func (v Value) IsEmpty() bool {
	return v.Size() == 0 && v.Type != ValueTypeNone // C++ empty checks size
}

// IsNone checks if the value is of type none.
// Corresponds to C++ value::operator==(nullptr)
func (v Value) IsNone() bool {
	return v.Type == ValueTypeNone
}

// IsString checks if the value is of type string.
func (v Value) IsString() bool {
	return v.Type == ValueTypeString
}

// IsObject checks if the value is of type object.
func (v Value) IsObject() bool {
	return v.Type == ValueTypeObject
}

// IsList checks if the value is of type list.
func (v Value) IsList() bool {
	return v.Type == ValueTypeList
}

// IsInteger checks if the value is of type integer.
func (v Value) IsInteger() bool {
	return v.Type == ValueTypeInteger
}

// IsNumber checks if the value is of type number.
func (v Value) IsNumber() bool {
	return v.Type == ValueTypeNumber
}

// IsBoolean checks if the value is of type boolean.
func (v Value) IsBoolean() bool {
	return v.Type == ValueTypeBoolean
}

// throwTypeError is a helper to panic on type mismatches.
// Corresponds to C++ throw_type_error.
func throwTypeError(got, expected ValueType) {
	panic(fmt.Sprintf("type error: expected %s, got %s", expected.TypeName(), got.TypeName()))
}

// IsNumeric checks if the value is integer or number.
// Corresponds to C++ is_numeric().
func IsNumeric(val Value) bool {
	return val.IsInteger() || val.IsNumber()
}

// Object creates a new object (map) Value.
// Corresponds to C++ dv::object().
func Object() Value {
	return Value{Type: ValueTypeObject, Value: make(map[string]Value)}
}

// List creates a new list Value.
// Corresponds to C++ dv::list().
func List() Value {
	return Value{Type: ValueTypeList, Value: make([]Value, 0)}
}

// OptionalValue represents a nullable value reference.
// Corresponds to C++ dv::optionalvalue struct.
type OptionalValue struct {
	Ptr *Value
}

// Get returns the underlying Value, or a zero Value if nil.
func (ov OptionalValue) Get() Value {
	if ov.Ptr == nil {
		return Value{} // Returns ValueTypeNone
	}
	return *ov.Ptr
}

// Exists checks if the optional value points to a valid value.
func (ov OptionalValue) Exists() bool {
	return ov.Ptr != nil && ov.Ptr.Type != ValueTypeNone
}

// GetAsString attempts to retrieve the value as a string.
// Corresponds to C++ optionalvalue::get(std::string& dst)
func (ov OptionalValue) GetAsString() (string, bool) {
	if !ov.Exists() || ov.Ptr.Type != ValueTypeString {
		return "", false
	}
	return ov.Ptr.Value.(string), true
}

// GetAsBoolean attempts to retrieve the value as a bool.
func (ov OptionalValue) GetAsBoolean() (bool, bool) {
	if !ov.Exists() || ov.Ptr.Type != ValueTypeBoolean {
		return false, false
	}
	return ov.Ptr.Value.(bool), true
}

// GetAsInteger attempts to retrieve the value as an int64.
func (ov OptionalValue) GetAsInteger() (int64, bool) {
	if !ov.Exists() {
		return 0, false
	}
	switch ov.Ptr.Type {
	case ValueTypeInteger:
		return ov.Ptr.Value.(int64), true
	case ValueTypeNumber:
		return int64(ov.Ptr.Value.(float64)), true
	default:
		return 0, false
	}
}

// GetAsNumber attempts to retrieve the value as a float64.
func (ov OptionalValue) GetAsNumber() (float64, bool) {
	if !ov.Exists() {
		return 0.0, false
	}
	switch ov.Ptr.Type {
	case ValueTypeNumber:
		return ov.Ptr.Value.(float64), true
	case ValueTypeInteger:
		return float64(ov.Ptr.Value.(int64)), true
	default:
		return 0.0, false
	}
}