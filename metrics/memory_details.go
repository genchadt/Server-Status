// metrics/memory_details.go
package metrics

import (
	"bytes"
	"fmt"
	"html/template"
	"serverstatus/utils"

	"github.com/shirou/gopsutil/mem"
)

// MemoryData represents memory usage information
type MemoryData struct {
	Total       string
	Used        string
	Free        string
	UsedPercent float64
}

// GetMemoryDetails retrieves memory usage details
func GetMemoryDetails() (*MemoryData, error) {
	// Use gopsutil to get memory information
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve memory details: %v", err)
	}

	// Convert memory values to human-readable format
	total := utils.FormatBytes(memoryInfo.Total)
	used := utils.FormatBytes(memoryInfo.Used)
	free := utils.FormatBytes(memoryInfo.Free)

	return &MemoryData{
		Total:       total,
		Used:        used,
		Free:        free,
		UsedPercent: memoryInfo.UsedPercent,
	}, nil
}

// FormatMemoryDetails formats memory details into an HTML table
func FormatMemoryDetails(memoryData *MemoryData) (string, error) {
	if memoryData == nil {
		return "<p>Unable to retrieve memory details.</p>", fmt.Errorf("nil MemoryData provided")
	}

	tmpl := `<table border="1">
    <tr>
        <th>Total</th>
        <th>Used</th>
        <th>Free</th>
        <th>Used%</th>
    </tr>
    <tr>
        <td>{{.Total}}</td>
        <td>{{.Used}}</td>
        <td>{{.Free}}</td>
        <td>{{.UsedPercent}}%</td>
    </tr>
    </table>`

	t := template.Must(template.New("memoryDetails").Parse(tmpl))
	var htmlOut bytes.Buffer
	err := t.Execute(&htmlOut, memoryData)
	if err != nil {
		return "<p>Error formatting memory details</p>", fmt.Errorf("error executing template: %v", err)
	}
	return htmlOut.String(), nil
}
