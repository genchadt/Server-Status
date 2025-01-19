// metrics/cpu_load.go
package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"regexp"
)

// CPULoad represents CPU load averages
type CPULoad struct {
	Load1  string
	Load5  string
	Load15 string
}

// GetCPULoad retrieves CPU load details
func GetCPULoad() (*CPULoad, error) {
	cmd := exec.Command("uptime")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`load average: ([\d\.]+), ([\d\.]+), ([\d\.]+)`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) != 4 {
		return nil, err
	}

	return &CPULoad{
		Load1:  matches[1],
		Load5:  matches[2],
		Load15: matches[3],
	}, nil
}

// FormatCPULoad formats CPU load details into an HTML table
func FormatCPULoad(cpuLoad *CPULoad) string {
	if cpuLoad == nil {
		return "<p>Unable to retrieve CPU load details.</p>"
	}
	tmpl := `<table border="1">
    <tr><th>1 Minute Load</th><th>5 Minute Load</th><th>15 Minute Load</th></tr>
    <tr>
        <td style="text-align:center">{{.Load1}}</td>
        <td style="text-align:center">{{.Load5}}</td>
        <td style="text-align:center">{{.Load15}}</td>
    </tr>
    </table>`

	t := template.Must(template.New("cpuLoad").Parse(tmpl))
	var htmlOut bytes.Buffer
	t.Execute(&htmlOut, cpuLoad)
	return htmlOut.String()
}
