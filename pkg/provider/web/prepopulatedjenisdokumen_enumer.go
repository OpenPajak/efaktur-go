// Code generated by "enumer -type=PrepopulatedJenisDokumen -trimprefix=PrepopulatedJenisDokumen_ -sql"; DO NOT EDIT.

package web

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const _PrepopulatedJenisDokumenName = "UNSELECTEDFPMPIBPEBCUKAIBC40BC25BC27BC41"

var _PrepopulatedJenisDokumenIndex = [...]uint8{0, 10, 13, 16, 19, 24, 28, 32, 36, 40}

const _PrepopulatedJenisDokumenLowerName = "unselectedfpmpibpebcukaibc40bc25bc27bc41"

func (i PrepopulatedJenisDokumen) String() string {
	if i < 0 || i >= PrepopulatedJenisDokumen(len(_PrepopulatedJenisDokumenIndex)-1) {
		return fmt.Sprintf("PrepopulatedJenisDokumen(%d)", i)
	}
	return _PrepopulatedJenisDokumenName[_PrepopulatedJenisDokumenIndex[i]:_PrepopulatedJenisDokumenIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _PrepopulatedJenisDokumenNoOp() {
	var x [1]struct{}
	_ = x[PrepopulatedJenisDokumen_UNSELECTED-(0)]
	_ = x[PrepopulatedJenisDokumen_FPM-(1)]
	_ = x[PrepopulatedJenisDokumen_PIB-(2)]
	_ = x[PrepopulatedJenisDokumen_PEB-(3)]
	_ = x[PrepopulatedJenisDokumen_CUKAI-(4)]
	_ = x[PrepopulatedJenisDokumen_BC40-(5)]
	_ = x[PrepopulatedJenisDokumen_BC25-(6)]
	_ = x[PrepopulatedJenisDokumen_BC27-(7)]
	_ = x[PrepopulatedJenisDokumen_BC41-(8)]
}

var _PrepopulatedJenisDokumenValues = []PrepopulatedJenisDokumen{PrepopulatedJenisDokumen_UNSELECTED, PrepopulatedJenisDokumen_FPM, PrepopulatedJenisDokumen_PIB, PrepopulatedJenisDokumen_PEB, PrepopulatedJenisDokumen_CUKAI, PrepopulatedJenisDokumen_BC40, PrepopulatedJenisDokumen_BC25, PrepopulatedJenisDokumen_BC27, PrepopulatedJenisDokumen_BC41}

var _PrepopulatedJenisDokumenNameToValueMap = map[string]PrepopulatedJenisDokumen{
	_PrepopulatedJenisDokumenName[0:10]:       PrepopulatedJenisDokumen_UNSELECTED,
	_PrepopulatedJenisDokumenLowerName[0:10]:  PrepopulatedJenisDokumen_UNSELECTED,
	_PrepopulatedJenisDokumenName[10:13]:      PrepopulatedJenisDokumen_FPM,
	_PrepopulatedJenisDokumenLowerName[10:13]: PrepopulatedJenisDokumen_FPM,
	_PrepopulatedJenisDokumenName[13:16]:      PrepopulatedJenisDokumen_PIB,
	_PrepopulatedJenisDokumenLowerName[13:16]: PrepopulatedJenisDokumen_PIB,
	_PrepopulatedJenisDokumenName[16:19]:      PrepopulatedJenisDokumen_PEB,
	_PrepopulatedJenisDokumenLowerName[16:19]: PrepopulatedJenisDokumen_PEB,
	_PrepopulatedJenisDokumenName[19:24]:      PrepopulatedJenisDokumen_CUKAI,
	_PrepopulatedJenisDokumenLowerName[19:24]: PrepopulatedJenisDokumen_CUKAI,
	_PrepopulatedJenisDokumenName[24:28]:      PrepopulatedJenisDokumen_BC40,
	_PrepopulatedJenisDokumenLowerName[24:28]: PrepopulatedJenisDokumen_BC40,
	_PrepopulatedJenisDokumenName[28:32]:      PrepopulatedJenisDokumen_BC25,
	_PrepopulatedJenisDokumenLowerName[28:32]: PrepopulatedJenisDokumen_BC25,
	_PrepopulatedJenisDokumenName[32:36]:      PrepopulatedJenisDokumen_BC27,
	_PrepopulatedJenisDokumenLowerName[32:36]: PrepopulatedJenisDokumen_BC27,
	_PrepopulatedJenisDokumenName[36:40]:      PrepopulatedJenisDokumen_BC41,
	_PrepopulatedJenisDokumenLowerName[36:40]: PrepopulatedJenisDokumen_BC41,
}

var _PrepopulatedJenisDokumenNames = []string{
	_PrepopulatedJenisDokumenName[0:10],
	_PrepopulatedJenisDokumenName[10:13],
	_PrepopulatedJenisDokumenName[13:16],
	_PrepopulatedJenisDokumenName[16:19],
	_PrepopulatedJenisDokumenName[19:24],
	_PrepopulatedJenisDokumenName[24:28],
	_PrepopulatedJenisDokumenName[28:32],
	_PrepopulatedJenisDokumenName[32:36],
	_PrepopulatedJenisDokumenName[36:40],
}

// PrepopulatedJenisDokumenString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PrepopulatedJenisDokumenString(s string) (PrepopulatedJenisDokumen, error) {
	if val, ok := _PrepopulatedJenisDokumenNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _PrepopulatedJenisDokumenNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to PrepopulatedJenisDokumen values", s)
}

// PrepopulatedJenisDokumenValues returns all values of the enum
func PrepopulatedJenisDokumenValues() []PrepopulatedJenisDokumen {
	return _PrepopulatedJenisDokumenValues
}

// PrepopulatedJenisDokumenStrings returns a slice of all String values of the enum
func PrepopulatedJenisDokumenStrings() []string {
	strs := make([]string, len(_PrepopulatedJenisDokumenNames))
	copy(strs, _PrepopulatedJenisDokumenNames)
	return strs
}

// IsAPrepopulatedJenisDokumen returns "true" if the value is listed in the enum definition. "false" otherwise
func (i PrepopulatedJenisDokumen) IsAPrepopulatedJenisDokumen() bool {
	for _, v := range _PrepopulatedJenisDokumenValues {
		if i == v {
			return true
		}
	}
	return false
}

func (i PrepopulatedJenisDokumen) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *PrepopulatedJenisDokumen) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of PrepopulatedJenisDokumen: %[1]T(%[1]v)", value)
	}

	val, err := PrepopulatedJenisDokumenString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
