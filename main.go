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

// MetricConfig defines the configuration for each metric
type MetricConfig struct {
	Name       string
	Timeout    time.Duration
	GetFunc    interface{}
	FormatFunc interface{}
}

// metricConfigs is a slice of MetricConfig, defining each metric to be collected
var metricConfigs = []MetricConfig{
	{
		Name:       "package updates",
		Timeout:    15 * time.Second,
		GetFunc:    metrics.GetPackageUpdates,
		FormatFunc: metrics.FormatPackageUpdates,
	},
	{
		Name:       "disk details",
		Timeout:    10 * time.Second,
		GetFunc:    metrics.GetDiskDetails,
		FormatFunc: metrics.FormatDiskDetails,
	},
	{
		Name:       "CPU load details",
		Timeout:    5 * time.Second,
		GetFunc:    metrics.GetCPULoad,
		FormatFunc: metrics.FormatCPULoad,
	},
	{
		Name:       "memory details",
		Timeout:    5 * time.Second,
		GetFunc:    metrics.GetMemoryDetails,
		FormatFunc: metrics.FormatMemoryDetails,
	},
	{
		Name:       "active SSH sessions",
		Timeout:    5 * time.Second,
		GetFunc:    metrics.GetActiveSSHSessions,
		FormatFunc: metrics.FormatActiveSSHSessions,
	},
	{
		Name:       "previous SSH sessions",
		Timeout:    5 * time.Second,
		GetFunc:    metrics.GetPreviousSSHSessions,
		FormatFunc: metrics.FormatPreviousSSHSessions,
	},
	{
		Name:       "network details",
		Timeout:    10 * time.Second,
		GetFunc:    metrics.GetNetworkDetails,
		FormatFunc: metrics.FormatNetworkDetails,
	},
	{
		Name:       "CrowdSec alerts",
		Timeout:    10 * time.Second,
		GetFunc:    metrics.GetCrowdSecAlerts,
		FormatFunc: metrics.FormatCrowdSecAlerts,
	},
	{
		Name:       "CrowdSec decisions",
		Timeout:    10 * time.Second,
		GetFunc:    metrics.GetCrowdSecDecisions,
		FormatFunc: metrics.FormatCrowdSecDecisions,
	},
}

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
		handleFatalError(log, "Failed to load .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
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

	// Metrics collection with per-metric timeouts
	var wg sync.WaitGroup
	metricsData := MetricsData{}

	for _, config := range metricConfigs {
		wg.Add(1)
		go func(cfg MetricConfig) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()

			// Use type assertions to call the correct functions
			switch getFunc := cfg.GetFunc.(type) {
			case func() ([]metrics.PackageUpdate, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.PackageUpdate) (string, error))
				metricsData.PackageUpdates = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() ([]metrics.DiskUsage, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.DiskUsage) (string, error))
				metricsData.DiskDetails = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() (*metrics.CPULoad, error):
				formatFunc := cfg.FormatFunc.(func(*metrics.CPULoad) (string, error))
				metricsData.CPULoadDetails = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() (*metrics.MemoryData, error):
				formatFunc := cfg.FormatFunc.(func(*metrics.MemoryData) (string, error))
				metricsData.MemoryDetails = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() ([]metrics.ActiveSession, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.ActiveSession) (string, error))
				metricsData.ActiveSSH = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() ([]metrics.PreviousSession, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.PreviousSession) (string, error))
				metricsData.PreviousSSH = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() ([]metrics.NetworkInterface, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.NetworkInterface) (string, error))
				metricsData.NetworkDetails = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() ([]metrics.Alert, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.Alert) (string, error))
				metricsData.CrowdSecAlerts = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			case func() ([]metrics.Decision, error):
				formatFunc := cfg.FormatFunc.(func([]metrics.Decision) (string, error))
				metricsData.CrowdSecDecisions = collectMetric(ctx, log, cfg.Name, getFunc, formatFunc)
			}
		}(config)
	}

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
		handleFatalError(log, "Failed to construct email body: %v", err)
	}

	// Send Email with retry
	if err := sendReportWithRetry(emailBody, log, cfg, 3); err != nil {
		log.Error("Failed to send email after retries: %v", err)
		// Consider saving the report to a file for later delivery
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
		return fmt.Sprintf("<p>Timeout while collecting %s: %v</p>", metricName, ctx.Err()) // Include error in output
	case <-done:
		if err != nil {
			logger.Error("Failed to get %s: %v", metricName, err)
			return fmt.Sprintf("<p>Unable to get %s: %v</p>", metricName, err) // Include error in output
		}
		formatted, err := formatFunc(data)
		if err != nil {
			logger.Error("Failed to format %s: %v", metricName, err)
			return fmt.Sprintf("<p>Unable to format %s: %v</p>", metricName, err) // Include error in output
		}
		return formatted
	}
}

// sendReportWithRetry sends the email report with a specified number of retries
func sendReportWithRetry(emailBody string, logger *logger.FileLogger, cfg *config.Config, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		err = sendReport(emailBody, logger, cfg)
		if err == nil {
			return nil // Success
		}
		logger.Error("Attempt %d: Failed to send email: %v", i+1, err)
		if i < retries-1 {
			time.Sleep(5 * time.Second) // Wait before retrying
		}
	}
	return err // Return the last error after all retries failed
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

// handleFatalError logs the fatal error and exits
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
