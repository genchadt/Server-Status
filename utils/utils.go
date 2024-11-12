// utils/utils.go
package utils

import "strings"

// SanitizeIPAddress takes an IP address as input and removes any "Ip:" prefix,
// replacing all dots with hyphens to sanitize the IP address format.
func SanitizeIPAddress(ip string) string {
	// Raw CrowdSec decisions include IPs with "Ip:" prefix
	ip = strings.TrimPrefix(ip, "Ip:")

	return strings.ReplaceAll(ip, ".", "-")
}
