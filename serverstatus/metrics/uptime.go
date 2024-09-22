package metrics

import (
    "bytes"
    "os/exec"
    "strings"
)

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
