package metrics

import (
    "bytes"
    "encoding/csv"
    "html/template"
    "os/exec"
    "strings"
)

// Alert represents a CrowdSec alert
type Alert struct {
    ID        string
    Scope     string
    Value     string
    Reason    string
    Country   string
    AS        string
    Decisions string
    CreatedAt string
}

// Decision represents a CrowdSec decision
type Decision struct {
    ID          string
    Source      string
    IP          string
    Reason      string
    Action      string
    Country     string
    AS          string
    EventsCount string
    Expiration  string
    Simulated   string
    AlertID     string
}

func GetCrowdSecAlerts() string {
    cmd := exec.Command("cscli", "alerts", "list", "-o", "raw")
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
        return "<p>Unable to retrieve CrowdSec alerts.</p>"
    }

    output := strings.TrimSpace(out.String())
    lines := strings.Split(output, "\n")
    if len(lines) <= 1 {
        return "<p>No alerts available.</p>"
    }

    reader := csv.NewReader(strings.NewReader(output))
    records, err := reader.ReadAll()
    if err != nil {
        return "<p>Error parsing CrowdSec alerts data.</p>"
    }

    // First line is header, rest are data
    var alerts []Alert
    for i, record := range records {
        if i == 0 {
            continue // skip header
        }
        if len(record) < 8 {
            continue // skip incomplete records
        }
        alerts = append(alerts, Alert{
            ID:        record[0],
            Scope:     record[1],
            Value:     record[2],
            Reason:    record[3],
            Country:   record[4],
            AS:        record[5],
            Decisions: record[6],
            CreatedAt: record[7],
        })
    }

    // Generate HTML table using template
    tmpl := `<table border="1">
    <tr>
        <th>ID</th>
        <th>Scope</th>
        <th>Value</th>
        <th>Reason</th>
        <th>Country</th>
        <th>AS</th>
        <th>Decisions</th>
        <th>Created At</th>
    </tr>
    {{range .}}
    <tr>
        <td>{{.ID}}</td>
        <td>{{.Scope}}</td>
        <td>{{.Value}}</td>
        <td>{{.Reason}}</td>
        <td>{{.Country}}</td>
        <td>{{.AS}}</td>
        <td>{{.Decisions}}</td>
        <td>{{.CreatedAt}}</td>
    </tr>
    {{end}}
    </table>`

    t := template.Must(template.New("crowdSecAlerts").Parse(tmpl))
    var htmlOut bytes.Buffer
    err = t.Execute(&htmlOut, alerts)
    if err != nil {
        return "<p>Error generating CrowdSec alerts table.</p>"
    }

    return htmlOut.String()
}

func GetCrowdSecDecisions() string {
    cmd := exec.Command("cscli", "decisions", "list", "-o", "raw")
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
        return "<p>Unable to retrieve CrowdSec decisions.</p>"
    }

    output := strings.TrimSpace(out.String())
    lines := strings.Split(output, "\n")
    if len(lines) <= 1 {
        return "<p>No decisions available.</p>"
    }

    reader := csv.NewReader(strings.NewReader(output))
    records, err := reader.ReadAll()
    if err != nil {
        return "<p>Error parsing CrowdSec decisions data.</p>"
    }

    // First line is header, rest are data
    var decisions []Decision
    for i, record := range records {
        if i == 0 {
            continue // skip header
        }
        if len(record) < 11 {
            continue // skip incomplete records
        }
        decisions = append(decisions, Decision{
            ID:          record[0],
            Source:      record[1],
            IP:          record[2],
            Reason:      record[3],
            Action:      record[4],
            Country:     record[5],
            AS:          record[6],
            EventsCount: record[7],
            Expiration:  record[8],
            Simulated:   record[9],
            AlertID:     record[10],
        })
    }

    // Generate HTML table using template
    tmpl := `<table border="1">
    <tr>
        <th>ID</th>
        <th>Source</th>
        <th>IP</th>
        <th>Reason</th>
        <th>Action</th>
        <th>Country</th>
        <th>AS</th>
        <th>Events Count</th>
        <th>Expiration</th>
        <th>Simulated</th>
        <th>Alert ID</th>
    </tr>
    {{range .}}
    <tr>
        <td>{{.ID}}</td>
        <td>{{.Source}}</td>
        <td>{{.IP}}</td>
        <td>{{.Reason}}</td>
        <td>{{.Action}}</td>
        <td>{{.Country}}</td>
        <td>{{.AS}}</td>
        <td>{{.EventsCount}}</td>
        <td>{{.Expiration}}</td>
        <td>{{.Simulated}}</td>
        <td>{{.AlertID}}</td>
    </tr>
    {{end}}
    </table>`

    t := template.Must(template.New("crowdSecDecisions").Parse(tmpl))
    var htmlOut bytes.Buffer
    err = t.Execute(&htmlOut, decisions)
    if err != nil {
        return "<p>Error generating CrowdSec decisions table.</p>"
    }

    return htmlOut.String()
}
