package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sample test structure for clarity
type SampleStruct struct {
	Field int
}

func TestArrayMapRetainIfPresent_Basic(t *testing.T) {
	// Test case 1: Simple mapping, no retained value
	array := []int{1, 2, 3}
	mapper := func(n int) (int, *string) {
		return n * 2, nil
	}

	mapped, retained := ArrayMapKeepSecondIfPresent(array, mapper)

	assert.Equal(t, []int{2, 4, 6}, mapped)
	assert.Nil(t, retained)
}

func TestArrayMapRetainIfPresent_RetainedValue(t *testing.T) {
	// Test case 2: Mapping with retained value
	array := []int{1, 2, 3}
	retainedValue := "retain-me"
	mapper := func(n int) (int, *string) {
		if n == 2 {
			return n * 2, &retainedValue
		}
		return n * 2, nil
	}

	mapped, retained := ArrayMapKeepSecondIfPresent(array, mapper)

	assert.Equal(t, []int{2, 4, 6}, mapped)
	assert.Equal(t, &retainedValue, retained)
}

func TestArrayMapRetainIfPresent_RetainLastNonNil(t *testing.T) {
	// Test case 3: Multiple retained values, last one is retained
	array := []int{1, 2, 3}
	firstRetained := "first"
	secondRetained := "second"
	mapper := func(n int) (int, *string) {
		switch n {
		case 1:
			return n * 2, &firstRetained
		case 3:
			return n * 2, &secondRetained
		default:
			return n * 2, nil
		}
	}

	mapped, retained := ArrayMapKeepSecondIfPresent(array, mapper)

	assert.Equal(t, []int{2, 4, 6}, mapped)
	assert.Equal(t, &secondRetained, retained)
}

func TestArrayMapRetainIfPresent_ErrorHandling(t *testing.T) {
	// Test case 4: Retain an error value
	array := []int{1, 2, 3}
	mapper := func(n int) (int, *error) {
		if n == 2 {
			err := errors.New("error encountered")
			return n * 2, &err
		}
		return n * 2, nil
	}

	mapped, retained := ArrayMapKeepSecondIfPresent(array, mapper)

	assert.Equal(t, []int{2, 4, 6}, mapped)
	assert.NotNil(t, retained)
	assert.EqualError(t, *retained, "error encountered")
}
