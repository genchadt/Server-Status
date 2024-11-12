package metrics

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetUptime retrieves the system uptime and returns it as a string in the format
// "X days, HH:MM"
func GetUptime() string {
	cmd := exec.Command("uptime", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unable to retrieve uptime"
	}
	return strings.TrimSpace(out.String())
}

// GetLastRebootTime retrieves the time of the last reboot and returns it as a string in
// the format "YYYY-MM-DD HH:MM:SS".
func GetLastRebootTime() string {
	cmd := exec.Command("uptime", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unable to retrieve last reboot time"
	}
	return strings.TrimSpace(out.String())
}
