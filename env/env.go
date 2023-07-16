package env

import (
	"fmt"
	"os"
)

func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func Get(key string) string {
	return os.Getenv(key)
}

func GetDefault(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return defaultValue
	}
	return v
}

func MustGet(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		Exitf("env: var %q not set or is empty", key)
	}
	return v
}

func Exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
