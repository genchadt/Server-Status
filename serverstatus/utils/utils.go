// utils/utils.go
package utils

import "strings"

// Dots-to-dashes: Avoid accidental navigation to malicious hosts
func SanitizeIPAddress(ip string) string {
	// Raw CrowdSec decisions include IPs with "Ip:" prefix
	ip = strings.TrimPrefix(ip, "Ip:")

	return strings.ReplaceAll(ip, ".", "-")
}
