package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"serverstatus/utils"
	"strings"
)

// GetActiveSSHSessions returns a string containing an HTML table of all active SSH sessions
// including user, terminal, login time, and IP address. If there are no active sessions, it
// returns a string indicating that. If there is an error running the "who" command, it returns
// a string indicating that.
func GetActiveSSHSessions() string {
	cmd := exec.Command("who")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "<p>Unable to retrieve active SSH sessions.</p>"
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "<p>No active SSH sessions found.</p>"
	}

	var data []struct {
		User      string
		Terminal  string
		LoginTime string
		IPAddress string
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			data = append(data, struct {
				User      string
				Terminal  string
				LoginTime string
				IPAddress string
			}{
				User:      fields[0],
				Terminal:  fields[1],
				LoginTime: fields[2] + " " + fields[3],
				IPAddress: utils.SanitizeIPAddress(fields[4]),
			})
		}
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
	t.Execute(&htmlOut, data)
	return htmlOut.String()
}

// GetPreviousSSHSessions returns a string containing an HTML table of the last 10 SSH sessions,
// including user, terminal, IP address, login time, and duration. If there are no recent logins,
// or if an error occurs while executing the "last" command, it returns a string indicating that.
func GetPreviousSSHSessions() string {
	cmd := exec.Command("last", "-n", "10")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "<p>Unable to retrieve previous SSH sessions.</p>"
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 || strings.Contains(lines[0], "wtmp begins") {
		return "<p>No recent SSH logins found.</p>"
	}

	var data []struct {
		User      string
		Terminal  string
		IPAddress string
		LoginTime string
		Duration  string
	}

	for _, line := range lines {
		if strings.Contains(line, "wtmp begins") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 10 {
			data = append(data, struct {
				User      string
				Terminal  string
				IPAddress string
				LoginTime string
				Duration  string
			}{
				User:      fields[0],
				Terminal:  fields[1],
				IPAddress: fields[2],
				LoginTime: fields[3] + " " + fields[4],
				Duration:  fields[9],
			})
		}
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
	t.Execute(&htmlOut, data)
	return htmlOut.String()
}
