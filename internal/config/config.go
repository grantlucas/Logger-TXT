// Package config resolves the effective configuration from flags, env vars, and defaults.
package config

import "os"

// DefaultFilePath is the default log file location.
const DefaultFilePath = "./log.txt"

// ResolveFilePath returns the effective log file path.
// Precedence: flag value > LOGGERTXT_PATH env var > default (./log.txt).
func ResolveFilePath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envPath := os.Getenv("LOGGERTXT_PATH"); envPath != "" {
		return envPath
	}
	return DefaultFilePath
}
