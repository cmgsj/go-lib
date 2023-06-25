// Package env provides a simple interface to environment variables.
package env

import (
	"fmt"
	"os"
)

// Get returns the value of the environment variable named by the key.
// If the variable is not set, it returns an empty string.
func Get(key string) string {
	return os.Getenv(key)
}

// GetDefault returns the value of the environment variable named by the key.
// If the variable is not set, it returns the default value.
func GetDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// MustGet returns the value of the environment variable named by the key.
// If the variable is not set, it exits the program.
func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Fprintf(os.Stderr, "error: must set the %q environment variable\n", key)
		os.Exit(1)
	}
	return val
}
