package util

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Custom type to handle PostgreSQL arrays of integers
type PGIntArray []int

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
