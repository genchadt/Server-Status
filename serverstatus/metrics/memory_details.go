package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"strings"
)

// GetMemoryDetails retrieves the current memory and swap usage and returns it as an HTML
// table. If there is an error while executing the "free -h" command, or if there is no
// memory information available, it returns a string indicating that. The table will have
// the following columns: total, used, free, shared, buff/cache, and available.
func GetMemoryDetails() string {
	cmd := exec.Command("free", "-h")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "<p>Unable to retrieve memory details.</p>"
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) < 3 {
		return "<p>Memory information unavailable.</p>"
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
	t.Execute(&htmlOut, map[string]interface{}{
		"Headers": headers,
		"Data":    data,
	})
	return htmlOut.String()
}
