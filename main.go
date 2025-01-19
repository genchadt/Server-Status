package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
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
		handleFatalError(log, "Failed to initialize logger: %v", err)
	}
	defer log.Close()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load .env file: %v", err)
		handleFatalError(log, "Failed to load .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration: %v", err)
		handleFatalError(log, "Failed to load configuration: %v", err)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	metricsData := MetricsData{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.PackageUpdates = collectMetric(ctx, log, "package updates", metrics.GetPackageUpdates, metrics.FormatPackageUpdates)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.DiskDetails = collectMetric(ctx, log, "disk details", metrics.GetDiskDetails, metrics.FormatDiskDetails)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.CPULoadDetails = collectMetric(ctx, log, "CPU load details", metrics.GetCPULoad, metrics.FormatCPULoad)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.MemoryDetails = collectMetric(ctx, log, "memory details", metrics.GetMemoryDetails, metrics.FormatMemoryDetails)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.ActiveSSH = collectMetric(ctx, log, "active SSH sessions", metrics.GetActiveSSHSessions, metrics.FormatActiveSSHSessions)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.PreviousSSH = collectMetric(ctx, log, "previous SSH sessions", metrics.GetPreviousSSHSessions, metrics.FormatPreviousSSHSessions)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.NetworkDetails = collectMetric(ctx, log, "network details", metrics.GetNetworkDetails, metrics.FormatNetworkDetails)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.CrowdSecAlerts = collectMetric(ctx, log, "CrowdSec alerts", metrics.GetCrowdSecAlerts, metrics.FormatCrowdSecAlerts)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsData.CrowdSecDecisions = collectMetric(ctx, log, "CrowdSec decisions", metrics.GetCrowdSecDecisions, metrics.FormatCrowdSecDecisions)
	}()

	wg.Wait()

	// Construct Email Body
	emailData := email.EmailData{
		ServerHostname:    serverHostname,
		ServerTime:        time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
		ServerUptime:      metrics.GetUptime(),
		LastRebootTime:    metrics.GetLastRebootTime(),
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
		handleFatalError(log, "Failed to construct email body: %v", err)
	}

	// Send Email
	if err := sendReport(emailBody, log, cfg); err != nil {
		log.Error("Failed to send email: %v", err)
		handleFatalError(log, "Failed to send email: %v", err)
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

// collectMetric is a generic helper function to collect and format metrics with context
func collectMetric[T any](ctx context.Context, logger *logger.FileLogger, metricName string, getFunc func() (T, error), formatFunc func(T) (string, error)) string {
	var data T
	var err error

	done := make(chan struct{})
	go func() {
		data, err = getFunc()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Error("Timeout while collecting %s: %v", metricName, ctx.Err())
		return fmt.Sprintf("<p>Timeout while collecting %s.</p>", metricName)
	case <-done:
		if err != nil {
			logger.Error("Failed to get %s: %v", metricName, err)
			return fmt.Sprintf("<p>Unable to get %s.</p>", metricName)
		}
		formatted, err := formatFunc(data)
		if err != nil {
			logger.Error("Failed to format %s: %v", metricName, err)
			return fmt.Sprintf("<p>Unable to format %s.</p>", metricName)
		}
		return formatted
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

func handleFatalError(log *logger.FileLogger, message string, err error) {
	log.Error(message, err)
	log.Close()
	os.Exit(1)
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
