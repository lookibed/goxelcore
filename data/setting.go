package data

import (
	"fmt"
	"strconv"
)

// SettingFormat represents the formatting style for a setting.
// Corresponds to C++ setting_format enum in voxelcore/src/data/setting.hpp
type SettingFormat int

const (
	SettingFormatSimple  SettingFormat = iota // simple
	SettingFormatPercent                   // percent
)

// BaseSetting provides common fields for all setting types.
type BaseSetting struct {
	Format SettingFormat
	// Note: C++ has an ObservableSetting base class with observer logic.
	// For now, we will simplify and not implement the observer pattern
	// directly in these base settings. If needed, this can be added later
	// via callbacks or channels.
}

// Setting defines the common behavior for all setting types.
type Setting interface {
	ResetToDefault()
	GetFormat() SettingFormat
	String() string // Go's Stringer interface for toString()
	// GetValue() interface{} // Depending on how we want to retrieve values generically
}

// NumberSetting represents a float64 setting with min/max bounds.
// Corresponds to C++ NumberSetting in voxelcore/src/data/setting.hpp
type NumberSetting struct {
	BaseSetting
	Value        float64
	DefaultValue float64
	Min          float64
	Max          float64
}

// NewNumberSetting creates a new NumberSetting.
func NewNumberSetting(value, min, max float64, format SettingFormat) *NumberSetting {
	return &NumberSetting{
		BaseSetting:  BaseSetting{Format: format},
		Value:        value,
		DefaultValue: value, // Initial value is also the default
		Min:          min,
		Max:          max,
	}
}

// ResetToDefault resets the setting to its default value.
func (s *NumberSetting) ResetToDefault() {
	s.Value = s.DefaultValue
}

// GetFormat returns the format of the setting.
func (s *NumberSetting) GetFormat() SettingFormat {
	return s.Format
}

// String returns the string representation of the setting's value.
func (s *NumberSetting) String() string {
	if s.Format == SettingFormatPercent {
		return fmt.Sprintf("%.0f%%", s.Value*100)
	}
	return fmt.Sprintf("%v", s.Value)
}

// IntegerSetting represents an int setting with min/max bounds.
// Corresponds to C++ IntegerSetting in voxelcore/src/data/setting.hpp
type IntegerSetting struct {
	BaseSetting
	Value        int
	DefaultValue int
	Min          int
	Max          int
}

// NewIntegerSetting creates a new IntegerSetting.
func NewIntegerSetting(value, min, max int, format SettingFormat) *IntegerSetting {
	return &IntegerSetting{
		BaseSetting:  BaseSetting{Format: format},
		Value:        value,
		DefaultValue: value, // Initial value is also the default
		Min:          min,
		Max:          max,
	}
}

// ResetToDefault resets the setting to its default value.
func (s *IntegerSetting) ResetToDefault() {
	s.Value = s.DefaultValue
}

// GetFormat returns the format of the setting.
func (s *IntegerSetting) GetFormat() SettingFormat {
	return s.Format
}

// String returns the string representation of the setting's value.
func (s *IntegerSetting) String() string {
	return strconv.Itoa(s.Value)
}

// FlagSetting represents a boolean setting.
// Corresponds to C++ FlagSetting in voxelcore/src/data/setting.hpp
type FlagSetting struct {
	BaseSetting
	Value        bool
	DefaultValue bool
}

// NewFlagSetting creates a new FlagSetting.
func NewFlagSetting(value bool, format SettingFormat) *FlagSetting {
	return &FlagSetting{
		BaseSetting:  BaseSetting{Format: format},
		Value:        value,
		DefaultValue: value, // Initial value is also the default
	}
}

// ResetToDefault resets the setting to its default value.
func (s *FlagSetting) ResetToDefault() {
	s.Value = s.DefaultValue
}

// GetFormat returns the format of the setting.
func (s *FlagSetting) GetFormat() SettingFormat {
	return s.Format
}

// String returns the string representation of the setting's value.
func (s *FlagSetting) String() string {
	return strconv.FormatBool(s.Value)
}

// Toggle flips the boolean value of the flag setting.
func (s *FlagSetting) Toggle() {
	s.Value = !s.Value
}

// StringSetting represents a string setting.
// Corresponds to C++ StringSetting in voxelcore/src/data/setting.hpp
type StringSetting struct {
	BaseSetting
	Value        string
	DefaultValue string
}

// NewStringSetting creates a new StringSetting.
func NewStringSetting(value string, format SettingFormat) *StringSetting {
	return &StringSetting{
		BaseSetting:  BaseSetting{Format: format},
		Value:        value,
		DefaultValue: value, // Initial value is also the default
	}
}

// ResetToDefault resets the setting to its default value.
func (s *StringSetting) ResetToDefault() {
	s.Value = s.DefaultValue
}

// GetFormat returns the format of the setting.
func (s *StringSetting) GetFormat() SettingFormat {
	return s.Format
}

// String returns the string representation of the setting's value.
func (s *StringSetting) String() string {
	return s.Value
}

// CreatePercent creates a NumberSetting with percentage format.
// Corresponds to C++ NumberSetting::createPercent
func CreatePercent(def float64) *NumberSetting {
	return NewNumberSetting(def, 0.0, 1.0, SettingFormatPercent)
}