package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"strings"
)

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
				IPAddress: fields[4],
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
