// metrics/crowdsec.go
package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"os/exec"
	"serverstatus/utils"
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

// GetCrowdSecAlerts retrieves CrowdSec alerts
func GetCrowdSecAlerts() ([]Alert, error) {
	cmd := exec.Command("cscli", "alerts", "list", "-o", "raw")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

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
			Value:     utils.SanitizeIPAddress(record[2]),
			Reason:    record[3],
			Country:   record[4],
			AS:        record[5],
			Decisions: record[6],
			CreatedAt: record[7],
		})
	}

	return alerts, nil
}

func FormatCrowdSecAlerts(alerts []Alert) (string, error) {
	if len(alerts) == 0 {
		return "<p>No alerts available.</p>", nil
	}

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
	err := t.Execute(&htmlOut, alerts)
	if err != nil {
		return "<p>Error generating CrowdSec alerts table.</p>", fmt.Errorf("error executing template: %v", err)
	}

	return htmlOut.String(), nil
}

// GetCrowdSecDecisions retrieves CrowdSec decisions
func GetCrowdSecDecisions() ([]Decision, error) {
	cmd := exec.Command("cscli", "decisions", "list", "-o", "raw")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

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
			IP:          utils.SanitizeIPAddress(record[2]),
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

	return decisions, nil
}

// FormatCrowdSecDecisions formats CrowdSec decisions into an HTML table
func FormatCrowdSecDecisions(decisions []Decision) (string, error) {
	if len(decisions) == 0 {
		return "<p>No decisions available.</p>", nil
	}

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
	err := t.Execute(&htmlOut, decisions)
	if err != nil {
		return "<p>Error generating CrowdSec decisions table.</p>", fmt.Errorf("error executing template: %v", err)
	}

	return htmlOut.String(), nil
}
