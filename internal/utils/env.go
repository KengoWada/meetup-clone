// Package utils contains common functionality and helper functions that are
// used throughout the application. These utilities are designed to be reused
// across various modules, ensuring consistent behavior and reducing code duplication.
package utils

import (
	"fmt"
	"os"
	"strconv"
)

// EnvGetString retrieves the value of the specified environment variable as a string.
// If the environment variable is not set, it returns the provided fallback value.
// If the fallback value is an empty string, the function will panic.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//   - fallback: The value to return if the environment variable is not found.
//
// Returns:
//   - A string representing the value of the environment variable or the fallback value.
func EnvGetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		if fallback == "" {
			panic(fmt.Sprintf("%s can not be empty", key))
		}

		return fallback
	}

	return value
}

// EnvGetInt retrieves the value of the specified environment variable as an integer.
// If the environment variable is not set or cannot be converted to an integer,
// it returns the provided fallback value.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//   - fallback: The value to return if the environment variable is not found
//     or cannot be parsed as an integer.
//
// Returns:
//   - An integer representing the value of the environment variable or the fallback value.
func EnvGetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return valueInt
}

// EnvGetBool retrieves the value of the specified environment variable as a boolean.
// If the environment variable is not set or cannot be converted to a boolean,
// it returns the provided fallback value.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//   - fallback: The value to return if the environment variable is not found
//     or cannot be parsed as a boolean.
//
// Returns:
//   - A boolean representing the value of the environment variable or the fallback value.
func EnvGetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	switch value {
	case "false":
		return false
	case "true":
		return true
	default:
		return fallback
	}
}
