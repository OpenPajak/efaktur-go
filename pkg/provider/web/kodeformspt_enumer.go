// Code generated by "enumer -type=KodeFormSpt -trimprefix=KodeFormSpt_ -json -text -yaml -sql"; DO NOT EDIT.

package web

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _KodeFormSptName = "A1A2B1B2B3"

var _KodeFormSptIndex = [...]uint8{0, 2, 4, 6, 8, 10}

const _KodeFormSptLowerName = "a1a2b1b2b3"

func (i KodeFormSpt) String() string {
	i -= 1
	if i < 0 || i >= KodeFormSpt(len(_KodeFormSptIndex)-1) {
		return fmt.Sprintf("KodeFormSpt(%d)", i+1)
	}
	return _KodeFormSptName[_KodeFormSptIndex[i]:_KodeFormSptIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _KodeFormSptNoOp() {
	var x [1]struct{}
	_ = x[KodeFormSpt_A1-(1)]
	_ = x[KodeFormSpt_A2-(2)]
	_ = x[KodeFormSpt_B1-(3)]
	_ = x[KodeFormSpt_B2-(4)]
	_ = x[KodeFormSpt_B3-(5)]
}

var _KodeFormSptValues = []KodeFormSpt{KodeFormSpt_A1, KodeFormSpt_A2, KodeFormSpt_B1, KodeFormSpt_B2, KodeFormSpt_B3}

var _KodeFormSptNameToValueMap = map[string]KodeFormSpt{
	_KodeFormSptName[0:2]:       KodeFormSpt_A1,
	_KodeFormSptLowerName[0:2]:  KodeFormSpt_A1,
	_KodeFormSptName[2:4]:       KodeFormSpt_A2,
	_KodeFormSptLowerName[2:4]:  KodeFormSpt_A2,
	_KodeFormSptName[4:6]:       KodeFormSpt_B1,
	_KodeFormSptLowerName[4:6]:  KodeFormSpt_B1,
	_KodeFormSptName[6:8]:       KodeFormSpt_B2,
	_KodeFormSptLowerName[6:8]:  KodeFormSpt_B2,
	_KodeFormSptName[8:10]:      KodeFormSpt_B3,
	_KodeFormSptLowerName[8:10]: KodeFormSpt_B3,
}

var _KodeFormSptNames = []string{
	_KodeFormSptName[0:2],
	_KodeFormSptName[2:4],
	_KodeFormSptName[4:6],
	_KodeFormSptName[6:8],
	_KodeFormSptName[8:10],
}

// KodeFormSptString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func KodeFormSptString(s string) (KodeFormSpt, error) {
	if val, ok := _KodeFormSptNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _KodeFormSptNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to KodeFormSpt values", s)
}

// KodeFormSptValues returns all values of the enum
func KodeFormSptValues() []KodeFormSpt {
	return _KodeFormSptValues
}

// KodeFormSptStrings returns a slice of all String values of the enum
func KodeFormSptStrings() []string {
	strs := make([]string, len(_KodeFormSptNames))
	copy(strs, _KodeFormSptNames)
	return strs
}

// IsAKodeFormSpt returns "true" if the value is listed in the enum definition. "false" otherwise
func (i KodeFormSpt) IsAKodeFormSpt() bool {
	for _, v := range _KodeFormSptValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for KodeFormSpt
func (i KodeFormSpt) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for KodeFormSpt
func (i *KodeFormSpt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("KodeFormSpt should be a string, got %s", data)
	}

	var err error
	*i, err = KodeFormSptString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for KodeFormSpt
func (i KodeFormSpt) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for KodeFormSpt
func (i *KodeFormSpt) UnmarshalText(text []byte) error {
	var err error
	*i, err = KodeFormSptString(string(text))
	return err
}

// MarshalYAML implements a YAML Marshaler for KodeFormSpt
func (i KodeFormSpt) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for KodeFormSpt
func (i *KodeFormSpt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = KodeFormSptString(s)
	return err
}

func (i KodeFormSpt) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *KodeFormSpt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of KodeFormSpt: %[1]T(%[1]v)", value)
	}

	val, err := KodeFormSptString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}