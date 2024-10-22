package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// JSONMap is a custom type to handle JSON objects
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for JSONMap
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface for JSONMap
func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &m)
}

// Custom type to handle PostgreSQL arrays
type PGIntArray []int
type PGStringArray []string

// Implement the sql.Scanner interface
func (a *PGIntArray) Scan(src interface{}) error {
	if src == nil {
		*a = []int{}
		return nil
	}

	// Convert the src (which will be a byte array) to a string
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot convert %v to string", src)
	}

	// Remove the curly braces from the array string and split by commas
	str = strings.Trim(str, "{}")
	if str == "" {
		*a = []int{}
		return nil
	}

	// Split the string by commas and convert to an integer slice
	elements := strings.Split(str, ",")
	result := make([]int, len(elements))
	for i, elem := range elements {
		var value int
		_, err := fmt.Sscanf(elem, "%d", &value)
		if err != nil {
			return fmt.Errorf("failed to parse element '%s': %w", elem, err)
		}
		result[i] = value
	}

	*a = result
	return nil
}

// Implement the driver.Valuer interface
func (a PGIntArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Convert the slice to a string in the format "{1,2,3}"
	elements := make([]string, len(a))
	for i, v := range a {
		elements[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("{%s}", strings.Join(elements, ",")), nil
}

// Implement the sql.Scanner interface
func (a *PGStringArray) Scan(src interface{}) error {
	if src == nil {
		*a = []string{}
		return nil
	}

	// Convert the src (which is a byte array) to a string
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot convert %v to string", src)
	}

	// Remove the curly braces from the array string and split by commas
	str = strings.Trim(str, "{}")
	if str == "" {
		*a = []string{}
		return nil
	}

	// Split the string by commas and convert to a string slice
	elements := strings.Split(str, ",")
	for i, elem := range elements {
		elements[i] = strings.Trim(elem, `"`) // Remove any double quotes around the strings
	}

	*a = elements
	return nil
}

// Implement the driver.Valuer interface
func (a PGStringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Convert the slice to a string in the format '{"foo","bar"}'
	elements := make([]string, len(a))
	for i, v := range a {
		elements[i] = `"` + v + `"` // Add quotes around each string element
	}
	return "{" + strings.Join(elements, ",") + "}", nil
}
