// metrics/disk_details.go
package metrics

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
	"time"
)

// DiskUsage represents disk usage information
type DiskUsage struct {
	Filesystem string
	Size       string
	Used       string
	Available  string
	UsePercent string
	MountedOn  string
}

// GetDiskDetails retrieves disk usage information
func GetDiskDetails() ([]DiskUsage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "df", "-h")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve disk details: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("no disk information available")
	}

	var diskUsages []DiskUsage

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		diskUsages = append(diskUsages, DiskUsage{
			Filesystem: fields[0],
			Size:       fields[1],
			Used:       fields[2],
			Available:  fields[3],
			UsePercent: fields[4],
			MountedOn:  fields[5],
		})
	}

	return diskUsages, nil
}

// FormatDiskDetails takes the disk usage data and formats it into an HTML table string.
func FormatDiskDetails(diskUsages []DiskUsage) (string, error) {
	if len(diskUsages) == 0 {
		return "<p>No disk information available.</p>", nil
	}

	tmpl := `<table border="1">
    <tr>
        <th>Filesystem</th>
        <th>Size</th>
        <th>Used</th>
        <th>Available</th>
        <th>Use%</th>
        <th>Mounted On</th>
    </tr>
    {{range .}}
    <tr>
        <td>{{.Filesystem}}</td>
        <td>{{.Size}}</td>
        <td>{{.Used}}</td>
        <td>{{.Available}}</td>
        <td>{{.UsePercent}}</td>
        <td>{{.MountedOn}}</td>
    </tr>
    {{end}}
    </table>`

	t := template.Must(template.New("diskDetails").Parse(tmpl))
	var htmlOut bytes.Buffer
	err := t.Execute(&htmlOut, diskUsages)
	if err != nil {
		return "<p>Error formatting disk details</p>", fmt.Errorf("error executing template: %v", err)
	}
	return htmlOut.String(), nil
}
