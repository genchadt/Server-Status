// metrics/memory_details.go
package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"strings"
)

// MemoryData represents memory usage information
type MemoryData struct {
	Headers []string
	Data    []map[string]string
}

// GetMemoryDetails retrieves memory usage details
func GetMemoryDetails() (*MemoryData, error) {
	cmd := exec.Command("free", "-h")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) < 3 {
		return nil, err
	}

	headers := strings.Fields(lines[0])
	var data []map[string]string

	for _, line := range lines[1:3] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		row := make(map[string]string)
		for i, header := range headers {
			if i < len(fields) {
				row[header] = fields[i]
			}
		}
		data = append(data, row)
	}

	return &MemoryData{
		Headers: headers,
		Data:    data,
	}, nil
}

// FormatMemoryDetails formats memory details into an HTML table
func FormatMemoryDetails(memoryData *MemoryData) string {
	if memoryData == nil {
		return "<p>Unable to retrieve memory details.</p>"
	}
	tmpl := `<table border="1">
    <tr>
    {{range $key := .Headers}}<th>{{$key}}</th>{{end}}
    </tr>
    {{range .Data}}
    <tr>
        {{range $key := $.Headers}}<td>{{index . $key}}</td>{{end}}
    </tr>
    {{end}}
    </table>`

	t := template.Must(template.New("memoryDetails").Parse(tmpl))
	var htmlOut bytes.Buffer
	t.Execute(&htmlOut, memoryData)
	return htmlOut.String()
}
