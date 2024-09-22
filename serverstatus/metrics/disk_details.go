package metrics

import (
    "bytes"
    "html/template"
    "os/exec"
    "strings"
)

func GetDiskDetails() string {
    cmd := exec.Command("df", "-h")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "<p>Unable to retrieve disk details.</p>"
    }

    lines := strings.Split(out.String(), "\n")
    if len(lines) < 2 {
        return "<p>No disk information available.</p>"
    }

    headers := strings.Fields(lines[0])
    var data []map[string]string

    for _, line := range lines[1:] {
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

    t := template.Must(template.New("diskDetails").Parse(tmpl))
    var htmlOut bytes.Buffer
    t.Execute(&htmlOut, map[string]interface{}{
        "Headers": headers,
        "Data":    data,
    })
    return htmlOut.String()
}
