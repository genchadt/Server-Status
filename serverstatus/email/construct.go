package email

import (
    "bytes"
    "html/template"
)

func ConstructEmailBody(data EmailData) string {
    tmpl := `
    <html>
        <head>
            <style>
                body {
                    font-family: Arial, sans-serif;
                }
                h1 {
                    font-size: 24px;
                    color: #333333;
                }
                h2 {
                    font-size: 20px;
                    color: #555555;
                }
                table {
                    width: 100%;
                    border-collapse: collapse;
                }
                table, th, td {
                    border: 1px solid #dddddd;
                }
                th, td {
                    padding: 8px;
                    text-align: left;
                }
                th {
                    background-color: #f2f2f2;
                    color: #333333;
                }
                tr:nth-child(even) {
                    background-color: #f9f9f9;
                }
            </style>
        </head>
        <body>
            <h1>Server Status Report: {{.ServerHostname}}</h1>
            <h2>Server Time:</h2>
                <p>{{.ServerTime}}</p>
            <h2>Uptime:</h2>
                <p>{{.ServerUptime}}</p>
            <h2>Available Package Updates:</h2>
                <p>{{.PackageUpdates}}</p>
            <h2>Disk Information:</h2>
                <p>{{.DiskDetails}}</p>
            <h2>Memory Usage:</h2>
                <p>{{.MemoryDetails}}</p>
            <h2>CPU Load:</h2>
                <p>{{.CPULoadDetails}}</p>
            <h2>Active SSH Sessions:</h2>
                <p>{{.ActiveSSH}}</p>
            <h2>Previous SSH Sessions:</h2>
                <p>{{.PreviousSSH}}</p>
            <h2>Network Information:</h2>
                <p>{{.NetworkDetails}}</p>
            <h2>CrowdSec Alerts:</h2>
                <p>{{.CrowdSecAlerts}}</p>
            <h2>CrowdSec Decisions:</h2>
                <p>{{.CrowdSecDecisions}}</p>
        </body>
    </html>
    `
    t := template.Must(template.New("emailBody").Parse(tmpl))
    var htmlOut bytes.Buffer
    t.Execute(&htmlOut, data)
    return htmlOut.String()
}
