package email

type EmailData struct {
	ServerHostname    string
	ServerTime        string
	ServerUptime      string
	LastRebootTime    string
	PackageUpdates    string
	DiskDetails       string
	MemoryDetails     string
	CPULoadDetails    string
	ActiveSSH         string
	PreviousSSH       string
	NetworkDetails    string
	CrowdSecAlerts    string
	CrowdSecDecisions string
}
