// utils/utils.go
package utils

import (
	"fmt"
	"strings"
)

// SanitizeIPAddress takes an IP address as input and removes any "Ip:" prefix,
// replacing all dots with hyphens to sanitize the IP address format.
func SanitizeIPAddress(ip string) string {
	ip = strings.ReplaceAll(ip, ".", "-")
	return strings.TrimPrefix(ip, "Ip:")
}

// FormatBytes converts a byte size into a human-readable string (e.g., 1.23 GB)
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
