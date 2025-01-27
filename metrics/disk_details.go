package metrics

import (
	"bytes"
	"fmt"
	"html/template"

	"serverstatus/logger"

	"github.com/shirou/gopsutil/disk"
)

// Initialize a logger instance for the metrics package
var log logger.Logger

func init() {
	logger, err := logger.NewFileLogger("logs/metrics.log")
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	log = logger
}

// DiskUsage represents disk usage information
type DiskUsage struct {
	Device      string
	Mountpoint  string
	Fstype      string
	Total       string
	Free        string
	Used        string
	UsedPercent float64
}

// GetDiskDetails retrieves disk usage information
func GetDiskDetails() ([]DiskUsage, error) {
	parts, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve disk partitions: %v", err)
	}

	var diskUsages []DiskUsage

	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			// Log the error using the logger
			log.Error("failed to retrieve disk usage for %s: %v", part.Mountpoint, err)
			continue
		}

		diskUsages = append(diskUsages, DiskUsage{
			Device:      part.Device,
			Mountpoint:  part.Mountpoint,
			Fstype:      part.Fstype,
			Total:       fmt.Sprintf("%.2f GB", float64(usage.Total)/1024/1024/1024),
			Free:        fmt.Sprintf("%.2f GB", float64(usage.Free)/1024/1024/1024),
			Used:        fmt.Sprintf("%.2f GB", float64(usage.Used)/1024/1024/1024),
			UsedPercent: usage.UsedPercent,
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
        <th>Device</th>
        <th>Mountpoint</th>
        <th>Fstype</th>
        <th>Total</th>
        <th>Free</th>
        <th>Used</th>
        <th>Used%</th>
    </tr>
    {{range .}}
    <tr>
        <td>{{.Device}}</td>
        <td>{{.Mountpoint}}</td>
        <td>{{.Fstype}}</td>
        <td>{{.Total}}</td>
        <td>{{.Free}}</td>
        <td>{{.Used}}</td>
        <td>{{.UsedPercent}}%</td>
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
