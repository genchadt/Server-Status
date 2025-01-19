// metrics/ssh_sessions.go
package metrics

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"serverstatus/utils"
	"strings"
)

// ActiveSession represents an active SSH session
type ActiveSession struct {
	User      string
	Terminal  string
	LoginTime string
	IPAddress string
}

// PreviousSession represents a previous SSH session
type PreviousSession struct {
	User      string
	Terminal  string
	IPAddress string
	LoginTime string
	Duration  string
}

// GetActiveSSHSessions retrieves active SSH sessions
func GetActiveSSHSessions() ([]ActiveSession, error) {
	cmd := exec.Command("who")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil, nil
	}

	var data []ActiveSession

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			data = append(data, ActiveSession{
				User:      fields[0],
				Terminal:  fields[1],
				LoginTime: fields[2] + " " + fields[3],
				IPAddress: utils.SanitizeIPAddress(fields[4]),
			})
		}
	}

	return data, nil
}

// FormatActiveSSHSessions formats active SSH sessions into an HTML table
func FormatActiveSSHSessions(data []ActiveSession) (string, error) {
	if len(data) == 0 {
		return "<p>No active SSH sessions found.</p>", nil
	}

	tmpl := `<table border="1">
    <tr><th>User</th><th>Terminal</th><th>Login Time</th><th>IP Address</th></tr>
    {{range .}}
    <tr>
        <td>{{.User}}</td>
        <td>{{.Terminal}}</td>
        <td>{{.LoginTime}}</td>
        <td>{{.IPAddress}}</td>
    </tr>
    {{end}}
    </table>`

	t := template.Must(template.New("activeSSH").Parse(tmpl))
	var htmlOut bytes.Buffer
	err := t.Execute(&htmlOut, data)
	if err != nil {
		return "<p>Error formatting active SSH sessions</p>", fmt.Errorf("error executing template: %v", err)
	}
	return htmlOut.String(), nil
}

// GetPreviousSSHSessions retrieves previous SSH sessions
func GetPreviousSSHSessions() ([]PreviousSession, error) {
	cmd := exec.Command("last", "-n", "10")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 || strings.Contains(lines[0], "wtmp begins") {
		return nil, nil
	}

	var data []PreviousSession

	for _, line := range lines {
		if strings.Contains(line, "wtmp begins") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 10 {
			data = append(data, PreviousSession{
				User:      fields[0],
				Terminal:  fields[1],
				IPAddress: fields[2],
				LoginTime: fields[3] + " " + fields[4],
				Duration:  fields[9],
			})
		}
	}

	return data, nil
}

// FormatPreviousSSHSessions formats previous SSH sessions into an HTML table
func FormatPreviousSSHSessions(data []PreviousSession) (string, error) {
	if len(data) == 0 {
		return "<p>No recent SSH logins found.</p>", nil
	}

	tmpl := `<table border="1">
    <tr><th>User</th><th>Terminal</th><th>IP Address</th><th>Login Time</th><th>Duration</th></tr>
    {{range .}}
    <tr>
        <td>{{.User}}</td>
        <td>{{.Terminal}}</td>
        <td>{{.IPAddress}}</td>
        <td>{{.LoginTime}}</td>
        <td>{{.Duration}}</td>
    </tr>
    {{end}}
    </table>`

	t := template.Must(template.New("previousSSH").Parse(tmpl))
	var htmlOut bytes.Buffer
	err := t.Execute(&htmlOut, data)
	if err != nil {
		return "<p>Error formatting previous SSH sessions</p>", fmt.Errorf("error executing template: %v", err)
	}
	return htmlOut.String(), nil
}
