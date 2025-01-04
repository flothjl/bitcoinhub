package main

import (
	"os"
)

func GetEnvWithDefault(key, defaultVal string) string {
	r := os.Getenv(key)
	if r == "" {
		return defaultVal
	}
	return r
}
