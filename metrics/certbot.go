package metrics

import (
	"bytes"
	"html/template"
	"os/exec"
	"regexp"
	"strings"
)

func GetCertbotCerts() string {
	cmd := exec.Command("certbot", "certificates")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "<p>Unable to retrieve Certbot certificates</p>"
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return "<p>No Certbot certificates found.</p>"
	}

	type Certificate struct {
		CertificateName string
		Domains         string
		ExpiryDate      string
		Validity        string
		CertPath        string
		KeyPath         string
	}

	var certs []Certificate
	var cert Certificate
	certFound := false

	reCertName := regexp.MustCompile(`^\s*Certificate Name:\s*(.*)`)
	reDomains := regexp.MustCompile(`^\s*Domains:\s*(.*)`)
	reExpiry := regexp.MustCompile(`^\s*Expiry Date:\s*(.*)\s*\(VALID:\s*(.*)\)`)
	reCertPath := regexp.MustCompile(`^\s*Certificate Path:\s*(.*)`)
	reKeyPath := regexp.MustCompile(`^\s*Private Key Path:\s*(.*)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if matches := reCertName.FindStringSubmatch(line); matches != nil {
			if certFound {
				certs = append(certs, cert)
			}
			cert = Certificate{}
			cert.CertificateName = matches[1]
			certFound = true
		} else if matches := reDomains.FindStringSubmatch(line); matches != nil {
			cert.Domains = matches[1]
		} else if matches := reExpiry.FindStringSubmatch(line); matches != nil {
			cert.ExpiryDate = matches[1]
			cert.Validity = matches[2]
		} else if matches := reCertPath.FindStringSubmatch(line); matches != nil {
			cert.CertPath = matches[1]
		} else if matches := reKeyPath.FindStringSubmatch(line); matches != nil {
			cert.KeyPath = matches[1]
		}
	}

	if certFound {
		certs = append(certs, cert)
	}

	if len(certs) == 0 {
		return "<p>No Certbot certificates found.</p>"
	}

	tmpl := `<table border="1">
    <tr>
        <th>Certificate Name</th>
        <th>Domains</th>
        <th>Expiry Date</th>
        <th>Validity</th>
        <th>Certificate Path</th>
        <th>Private Key Path</th>
    </tr>
    {{range .}}
    <tr>
        <td>{{.CertificateName}}</td>
        <td>{{.Domains}}</td>
        <td>{{.ExpiryDate}}</td>
        <td>{{.Validity}}</td>
        <td>{{.CertPath}}</td>
        <td>{{.KeyPath}}</td>
    </tr>
    {{end}}
    </table>`

	t := template.Must(template.New("certbotCertificates").Parse(tmpl))
	var htmlOut bytes.Buffer
	err = t.Execute(&htmlOut, certs)
	if err != nil {
		return "<p>Unable to render Certbot certificates</p>"
	}

	return htmlOut.String()
}
