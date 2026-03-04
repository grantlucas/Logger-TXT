// Package config resolves the effective configuration from flags, env vars, and defaults.
package config

import "os"

// DefaultFilePath is the default log file location.
const DefaultFilePath = "./log.txt"

// ResolveFilePath returns the effective log file path.
// Precedence: flag value > local ./log.txt (if it exists) > LOGGERTXT_PATH env var > ./log.txt default.
func ResolveFilePath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if _, err := os.Stat(DefaultFilePath); err == nil {
		return DefaultFilePath
	}
	if envPath := os.Getenv("LOGGERTXT_PATH"); envPath != "" {
		return envPath
	}
	return DefaultFilePath
}
