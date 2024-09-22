package metrics

import (
    "bytes"
    "html/template"
    "os/exec"
    "strings"
)

func GetNetworkDetails() string {
    cmd := exec.Command("ip", "-brief", "addr", "show")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "<p>Unable to retrieve network details.</p>"
    }

    lines := strings.Split(strings.TrimSpace(out.String()), "\n")
    if len(lines) == 0 {
        return "<p>No network information available.</p>"
    }

    var data []struct {
        Interface string
        State     string
        IPAddress string
    }

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) >= 3 {
            data = append(data, struct {
                Interface string
                State     string
                IPAddress string
            }{
                Interface: fields[0],
                State:     fields[1],
                IPAddress: fields[2],
            })
        }
    }

    tmpl := `<table border="1">
    <tr><th>Interface</th><th>State</th><th>IP Address</th></tr>
    {{range .}}
    <tr>
        <td>{{.Interface}}</td>
        <td>{{.State}}</td>
        <td>{{.IPAddress}}</td>
    </tr>
    {{end}}
    </table>`

    t := template.Must(template.New("networkDetails").Parse(tmpl))
    var htmlOut bytes.Buffer
    t.Execute(&htmlOut, data)
    return htmlOut.String()
}
