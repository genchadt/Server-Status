package metrics

import (
    "bytes"
    "html/template"
    "os/exec"
    "regexp"
    // "strings"
)

func GetCPULoadDetails() string {
    cmd := exec.Command("uptime")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "<p>Unable to retrieve CPU load details.</p>"
    }

    re := regexp.MustCompile(`load average: ([\d\.]+), ([\d\.]+), ([\d\.]+)`)
    matches := re.FindStringSubmatch(out.String())
    if len(matches) != 4 {
        return "<p>Unable to parse CPU load details.</p>"
    }

    data := struct {
        Load1  string
        Load5  string
        Load15 string
    }{
        Load1:  matches[1],
        Load5:  matches[2],
        Load15: matches[3],
    }

    tmpl := `<table border="1">
    <tr><th>1 Minute Load</th><th>5 Minute Load</th><th>15 Minute Load</th></tr>
    <tr>
        <td style="text-align:center">{{.Load1}} %</td>
        <td style="text-align:center">{{.Load5}} %</td>
        <td style="text-align:center">{{.Load15}} %</td>
    </tr>
    </table>`

    t := template.Must(template.New("cpuLoad").Parse(tmpl))
    var htmlOut bytes.Buffer
    t.Execute(&htmlOut, data)
    return htmlOut.String()
}
