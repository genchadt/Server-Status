package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"serverstatus/email"
	"serverstatus/logger"
	"serverstatus/metrics"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logPath := filepath.Join("logs", "serverstatus.log")
	logger, err := logger.NewFileLogger(logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Error("Failed to load .env file: %v", err)
		os.Exit(1)
	}

	// Log start of execution
	logger.Info("Starting server status check")

	// Server details here
	serverHostname := os.Getenv("SERVER_HOSTNAME")
	if serverHostname == "" {
		serverHostname = "Lightsail Web"
	}

	serverTime := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
	logger.Info("Collecting server metrics")

	// Metrics collection with error logging
	serverUptime := metrics.GetUptime()
	lastRebootTime := metrics.GetLastRebootTime()

	// Collect all metrics with logging
	metrics := collectMetrics(logger)

	// Construct Email Body
	emailBody := email.ConstructEmailBody(email.EmailData{
		ServerHostname:    serverHostname,
		ServerTime:        serverTime,
		ServerUptime:      serverUptime,
		LastRebootTime:    lastRebootTime,
		PackageUpdates:    metrics.packageUpdates,
		DiskDetails:       metrics.diskDetails,
		MemoryDetails:     metrics.memoryDetails,
		CPULoadDetails:    metrics.cpuLoadDetails,
		ActiveSSH:         metrics.activeSSH,
		PreviousSSH:       metrics.previousSSH,
		NetworkDetails:    metrics.networkDetails,
		CrowdSecAlerts:    metrics.crowdSecAlerts,
		CrowdSecDecisions: metrics.crowdSecDecisions,
	})

	// Send Email
	today := time.Now().Format("2006-01-02")
	subject := fmt.Sprintf("Daily System Report, %s", today)
	recipient := os.Getenv("TO_EMAIL")
	if recipient == "" {
		recipient = "webmaster@timothywb.com"
	}

	if err := email.SendEmail(subject, emailBody, recipient); err != nil {
		logger.LogEmailError(err)
		os.Exit(1)
	}

	logger.LogEmailSuccess(recipient)
	logger.Info("Server status check completed successfully")
}

// MetricsData holds all collected metrics
type MetricsData struct {
	packageUpdates    string
	diskDetails       string
	cpuLoadDetails    string
	memoryDetails     string
	activeSSH         string
	previousSSH       string
	networkDetails    string
	crowdSecAlerts    string
	crowdSecDecisions string
}

// collectMetrics gathers all system metrics with logging
func collectMetrics(logger *logger.FileLogger) MetricsData {
	logger.Info("Collecting package updates")
	packageUpdates := metrics.GetPackageUpdates()

	logger.Info("Collecting disk details")
	diskDetails := metrics.GetDiskDetails()

	logger.Info("Collecting CPU load details")
	cpuLoadDetails := metrics.GetCPULoadDetails()

	logger.Info("Collecting memory details")
	memoryDetails := metrics.GetMemoryDetails()

	logger.Info("Collecting SSH session information")
	activeSSH := metrics.GetActiveSSHSessions()
	previousSSH := metrics.GetPreviousSSHSessions()

	logger.Info("Collecting network details")
	networkDetails := metrics.GetNetworkDetails()

	logger.Info("Collecting CrowdSec information")
	crowdSecAlerts := metrics.GetCrowdSecAlerts()
	crowdSecDecisions := metrics.GetCrowdSecDecisions()

	return MetricsData{
		packageUpdates:    packageUpdates,
		diskDetails:       diskDetails,
		cpuLoadDetails:    cpuLoadDetails,
		memoryDetails:     memoryDetails,
		activeSSH:         activeSSH,
		previousSSH:       previousSSH,
		networkDetails:    networkDetails,
		crowdSecAlerts:    crowdSecAlerts,
		crowdSecDecisions: crowdSecDecisions,
	}
}
