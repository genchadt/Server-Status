// metrics/package_updates.go
package metrics

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
)

// PackageUpdate represents a package update
type PackageUpdate struct {
	Package        string
	CurrentVersion string
	NewVersion     string
}

// GetPackageUpdates retrieves a list of upgradable packages
func GetPackageUpdates() ([]PackageUpdate, error) {
	cmd := exec.Command("apt", "list", "--upgradable")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	var data []PackageUpdate

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pkgInfo := strings.Split(fields[0], "/")
			if len(pkgInfo) > 0 {
				data = append(data, PackageUpdate{
					Package:        pkgInfo[0],
					CurrentVersion: fields[1],
					NewVersion:     fields[1],
				})
			}
		}
	}

	return data, nil
}

// FormatPackageUpdates formats package updates into an HTML table
func FormatPackageUpdates(data []PackageUpdate) (string, error) {
	if len(data) == 0 {
		return "<p>All packages are up to date.</p>", nil
	}

	tmpl := `<table border="1">
    <tr><th>Package</th><th>Current Version</th><th>New Version</th></tr>
    {{range .}}
    <tr><td>{{.Package}}</td><td>{{.CurrentVersion}}</td><td>{{.NewVersion}}</td></tr>
    {{end}}
    </table>`
	t := template.Must(template.New("packageUpdates").Parse(tmpl))
	var htmlOut bytes.Buffer
	err := t.Execute(&htmlOut, data)
	if err != nil {
		return "<p>Error formatting package updates</p>", fmt.Errorf("error executing template: %v", err)
	}
	return htmlOut.String(), nil
}
