package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"strings"
)

// GetPackageUpdates retrieves a list of upgradable packages and returns it as an HTML table.
// If there are no upgradable packages, it returns a message indicating all packages are up to date.
// In case of an error executing the command, it returns an error message.
func GetPackageUpdates() string {
	cmd := exec.Command("apt", "list", "--upgradable")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "<p>Unable to retrieve package updates.</p>"
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) <= 1 {
		return "<p>All packages are up to date.</p>"
	}

	var data []struct {
		Package        string
		CurrentVersion string
		NewVersion     string
	}

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pkgInfo := strings.Split(fields[0], "/")
			if len(pkgInfo) > 0 {
				data = append(data, struct {
					Package        string
					CurrentVersion string
					NewVersion     string
				}{
					Package:        pkgInfo[0],
					CurrentVersion: fields[1],
					NewVersion:     fields[1],
				})
			}
		}
	}

	tmpl := `<table border="1">
    <tr><th>Package</th><th>Current Version</th><th>New Version</th></tr>
    {{range .}}
    <tr><td>{{.Package}}</td><td>{{.CurrentVersion}}</td><td>{{.NewVersion}}</td></tr>
    {{end}}
    </table>`
	t := template.Must(template.New("packageUpdates").Parse(tmpl))
	var htmlOut bytes.Buffer
	t.Execute(&htmlOut, data)
	return htmlOut.String()
}
