package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"serverstatus/config"
	"serverstatus/email"
	"serverstatus/logger"
	"serverstatus/metrics"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logPath := filepath.Join("logs", "serverstatus.log")
	log, err := logger.NewFileLogger(logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load .env file: %v", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Log start of execution
	log.Info("Starting server status check")

	// Set server hostname
	serverHostname := cfg.ServerHostname
	if serverHostname == "" {
		log.Info("SERVER_HOSTNAME environment variable not set. Using hostname from system.")
		serverHostname, _ = os.Hostname()
	}

	// Handle graceful shutdown
	setupGracefulShutdown(log)

	// Metrics collection with error logging and timeout
	serverUptime := metrics.GetUptime()
	lastRebootTime := metrics.GetLastRebootTime()

	// Collect all metrics with logging
	metricsData := collectMetrics(log)

	// Construct Email Body
	emailData := email.EmailData{
		ServerHostname:    serverHostname,
		ServerTime:        time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
		ServerUptime:      serverUptime,
		LastRebootTime:    lastRebootTime,
		PackageUpdates:    metricsData.PackageUpdates,
		DiskDetails:       metricsData.DiskDetails,
		MemoryDetails:     metricsData.MemoryDetails,
		CPULoadDetails:    metricsData.CPULoadDetails,
		ActiveSSH:         metricsData.ActiveSSH,
		PreviousSSH:       metricsData.PreviousSSH,
		NetworkDetails:    metricsData.NetworkDetails,
		CrowdSecAlerts:    metricsData.CrowdSecAlerts,
		CrowdSecDecisions: metricsData.CrowdSecDecisions,
	}
	emailBody, err := email.ConstructEmailBody(emailData)
	if err != nil {
		log.Error("Failed to construct email body: %v", err)
		os.Exit(1)
	}

	// Send Email
	if err := sendReport(emailBody, log, cfg); err != nil {
		log.Error("Failed to send email: %v", err)
		os.Exit(1)
	}

	log.Info("Server status check completed successfully")
}

// MetricsData holds all collected metrics
type MetricsData struct {
	PackageUpdates    string
	DiskDetails       string
	CPULoadDetails    string
	MemoryDetails     string
	ActiveSSH         string
	PreviousSSH       string
	NetworkDetails    string
	CrowdSecAlerts    string
	CrowdSecDecisions string
}

// collectMetrics gathers all system metrics with logging
func collectMetrics(logger *logger.FileLogger) MetricsData {
	logger.Info("Collecting package updates...")
	packageUpdatesData, err := metrics.GetPackageUpdates()
	if err != nil {
		logger.Error("Failed to get package updates: %v", err)
	}
	packageUpdates := metrics.FormatPackageUpdates(packageUpdatesData)

	logger.Info("Collecting disk details...")
	diskDetails, err := metrics.GetDiskDetails()
	if err != nil {
		logger.Error("Failed to get disk details: %v", err)
	}
	diskDetailsHTML := metrics.FormatDiskDetails(diskDetails)

	logger.Info("Collecting CPU load details...")
	cpuLoad, err := metrics.GetCPULoad()
	var cpuLoadDetails string
	if err != nil {
		logger.Error("Failed to get CPU load details: %v", err)
		cpuLoadDetails = "<p>Unable to retrieve CPU load details.</p>"
	} else {
		cpuLoadDetails, err = metrics.FormatCPULoad(cpuLoad) // Now expects string, error
		if err != nil {
			logger.Error("Failed to format CPU load details: %v", err)
			cpuLoadDetails = "<p>Unable to format CPU load details.</p>"
		}
	}

	logger.Info("Collecting memory details...")
	memoryData, err := metrics.GetMemoryDetails()
	if err != nil {
		logger.Error("Failed to get memory details: %v", err)
	}
	memoryDetails := metrics.FormatMemoryDetails(memoryData)

	logger.Info("Collecting SSH session information...")
	activeSSHData, err := metrics.GetActiveSSHSessions()
	if err != nil {
		logger.Error("Failed to get active SSH sessions: %v", err)
	}
	activeSSH := metrics.FormatActiveSSHSessions(activeSSHData)

	previousSSHData, err := metrics.GetPreviousSSHSessions()
	if err != nil {
		logger.Error("Failed to get previous SSH sessions: %v", err)
	}
	previousSSH := metrics.FormatPreviousSSHSessions(previousSSHData)

	logger.Info("Collecting network details...")
	networkData, err := metrics.GetNetworkDetails()
	if err != nil {
		logger.Error("Failed to get network details: %v", err)
	}
	networkDetails := metrics.FormatNetworkDetails(networkData)

	logger.Info("Collecting CrowdSec information...")
	crowdSecAlertsData, err := metrics.GetCrowdSecAlerts()
	if err != nil {
		logger.Error("Failed to get CrowdSec alerts: %v", err)
	}
	crowdSecAlerts := metrics.FormatCrowdSecAlerts(crowdSecAlertsData)

	crowdSecDecisionsData, err := metrics.GetCrowdSecDecisions()
	if err != nil {
		logger.Error("Failed to get CrowdSec decisions: %v", err)
	}
	crowdSecDecisions := metrics.FormatCrowdSecDecisions(crowdSecDecisionsData)

	return MetricsData{
		PackageUpdates:    packageUpdates,
		DiskDetails:       diskDetailsHTML,
		CPULoadDetails:    cpuLoadDetails,
		MemoryDetails:     memoryDetails,
		ActiveSSH:         activeSSH,
		PreviousSSH:       previousSSH,
		NetworkDetails:    networkDetails,
		CrowdSecAlerts:    crowdSecAlerts,
		CrowdSecDecisions: crowdSecDecisions,
	}
}

// sendReport sends the email report
func sendReport(emailBody string, logger *logger.FileLogger, cfg *config.Config) error {
	today := time.Now().Format("2006-01-02")
	subject := fmt.Sprintf("Daily System Report, %s", today)
	recipient := cfg.ToEmail
	if recipient == "" {
		return fmt.Errorf("TO_EMAIL environment variable not set")
	}

	if err := email.SendEmail(subject, emailBody, recipient); err != nil {
		logger.LogEmailError(err)
		return err
	}

	logger.LogEmailSuccess(recipient)
	return nil
}

// setupGracefulShutdown handles graceful shutdown of the application
func setupGracefulShutdown(log *logger.FileLogger) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChannel
		log.Info("Received signal: %v", sig)
		log.Info("Shutting down gracefully...")

		// TODO: Add any cleanup tasks here if needed

		os.Exit(0)
	}()
}
