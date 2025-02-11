// metrics/network_details.go
package metrics

import (
	"bytes"
	"fmt"
	"html/template"
	"net"

	"github.com/vishvananda/netlink"
)

// NetworkInterface represents network interface information
type NetworkInterface struct {
	Interface string
	State     string
	IPAddress net.IP
}

// GetNetworkDetails retrieves network interface details
func GetNetworkDetails() ([]NetworkInterface, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	var data []NetworkInterface

	for _, link := range links {
		attrs := link.Attrs()
		if attrs == nil {
			continue
		}

		addrs, err := netlink.AddrList(link, 0)
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			data = append(data, NetworkInterface{
				Interface: attrs.Name,
				State:     linkStateToString(attrs.OperState),
				IPAddress: addr.IP,
			})
		}
	}

	return data, nil
}

// FormatNetworkDetails formats network details into an HTML table
func FormatNetworkDetails(data []NetworkInterface) (string, error) {
	if len(data) == 0 {
		return "<p>No network information available.</p>", nil
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
	err := t.Execute(&htmlOut, data)
	if err != nil {
		return "<p>Error formatting network details.</p>", fmt.Errorf("error executing template: %v", err)
	}
	return htmlOut.String(), nil
}

// Helper function to convert netlink.LinkOperState to string
func linkStateToString(state netlink.LinkOperState) string {
	switch state {
	case netlink.OperUp:
		return "UP"
	case netlink.OperDown:
		return "DOWN"
	case netlink.OperUnknown:
		return "UNKNOWN"
	default:
		return "OTHER"
	}
}
