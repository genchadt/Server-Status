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

type GetPackageUpdatesFunc func() ([]metrics.PackageUpdate, error)
type FormatPackageUpdatesFunc func([]metrics.PackageUpdate) (string, error)

type GetDiskDetailsFunc func() ([]metrics.DiskUsage, error)
type FormatDiskDetailsFunc func([]metrics.DiskUsage) (string, error)

type GetCPULoadFunc func() (*metrics.CPULoad, error)
type FormatCPULoadFunc func(*metrics.CPULoad) (string, error)

type GetMemoryDetailsFunc func() (*metrics.MemoryData, error)
type FormatMemoryDetailsFunc func(*metrics.MemoryData) (string, error)

type GetActiveSSHSessionsFunc func() ([]metrics.ActiveSession, error)
type FormatActiveSSHSessionsFunc func([]metrics.ActiveSession) (string, error)

type GetPreviousSSHSessionsFunc func() ([]metrics.PreviousSession, error)
type FormatPreviousSSHSessionsFunc func([]metrics.PreviousSession) (string, error)

type GetNetworkDetailsFunc func() ([]metrics.NetworkInterface, error)
type FormatNetworkDetailsFunc func([]metrics.NetworkInterface) (string, error)

type GetCrowdSecAlertsFunc func() ([]metrics.Alert, error)
type FormatCrowdSecAlertsFunc func([]metrics.Alert) (string, error)

type GetCrowdSecDecisionsFunc func() ([]metrics.Decision, error)
type FormatCrowdSecDecisionsFunc func([]metrics.Decision) (string, error)

// MetricConfig defines the configuration for each metric
type MetricConfig struct {
	Name       string
	Timeout    time.Duration
	GetFunc    interface{} // Keep as interface{} to hold different function types
	FormatFunc interface{} // Keep as interface{} to hold different function types
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

			// Call collectMetric with explicit type arguments
			switch cfg.Name {
			case "package updates":
				metricsData.PackageUpdates = collectMetric[[]metrics.PackageUpdate](ctx, log, cfg.Name, cfg.GetFunc.(GetPackageUpdatesFunc), cfg.FormatFunc.(FormatPackageUpdatesFunc))
			case "disk details":
				metricsData.DiskDetails = collectMetric[[]metrics.DiskUsage](ctx, log, cfg.Name, cfg.GetFunc.(GetDiskDetailsFunc), cfg.FormatFunc.(FormatDiskDetailsFunc))
			case "CPU load details":
				metricsData.CPULoadDetails = collectMetric[*metrics.CPULoad](ctx, log, cfg.Name, cfg.GetFunc.(GetCPULoadFunc), cfg.FormatFunc.(FormatCPULoadFunc))
			case "memory details":
				metricsData.MemoryDetails = collectMetric[*metrics.MemoryData](ctx, log, cfg.Name, cfg.GetFunc.(GetMemoryDetailsFunc), cfg.FormatFunc.(FormatMemoryDetailsFunc))
			case "active SSH sessions":
				metricsData.ActiveSSH = collectMetric[[]metrics.ActiveSession](ctx, log, cfg.Name, cfg.GetFunc.(GetActiveSSHSessionsFunc), cfg.FormatFunc.(FormatActiveSSHSessionsFunc))
			case "previous SSH sessions":
				metricsData.PreviousSSH = collectMetric[[]metrics.PreviousSession](ctx, log, cfg.Name, cfg.GetFunc.(GetPreviousSSHSessionsFunc), cfg.FormatFunc.(FormatPreviousSSHSessionsFunc))
			case "network details":
				metricsData.NetworkDetails = collectMetric[[]metrics.NetworkInterface](ctx, log, cfg.Name, cfg.GetFunc.(GetNetworkDetailsFunc), cfg.FormatFunc.(FormatNetworkDetailsFunc))
			case "CrowdSec alerts":
				metricsData.CrowdSecAlerts = collectMetric[[]metrics.Alert](ctx, log, cfg.Name, cfg.GetFunc.(GetCrowdSecAlertsFunc), cfg.FormatFunc.(FormatCrowdSecAlertsFunc))
			case "CrowdSec decisions":
				metricsData.CrowdSecDecisions = collectMetric[[]metrics.Decision](ctx, log, cfg.Name, cfg.GetFunc.(GetCrowdSecDecisionsFunc), cfg.FormatFunc.(FormatCrowdSecDecisionsFunc))
			default:
				log.Error("Unknown metric type: %s", cfg.Name)
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
func collectMetric[T any](ctx context.Context, logger *logger.FileLogger, metricName string, getFunc interface{}, formatFunc interface{}) string {
	var data T
	var err error

	done := make(chan struct{})
	go func() {
		defer close(done)

		switch f := getFunc.(type) {
		case GetPackageUpdatesFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetDiskDetailsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetCPULoadFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetMemoryDetailsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetActiveSSHSessionsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetPreviousSSHSessionsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetNetworkDetailsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetCrowdSecAlertsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		case GetCrowdSecDecisionsFunc:
			result, getErr := f()
			if getErr == nil {
				data = any(result).(T)
			}
			err = getErr
		default:
			logger.Error("Unknown GetFunc type for metric: %s", metricName)
			return
		}
	}()

	select {
	case <-ctx.Done():
		logger.Error("Timeout while collecting %s: %v", metricName, ctx.Err())
		return fmt.Sprintf("<p>Timeout while collecting %s: %v</p>", metricName, ctx.Err())
	case <-done:
		if err != nil {
			logger.Error("Failed to get %s: %v", metricName, err)
			return fmt.Sprintf("<p>Unable to get %s: %v</p>", metricName, err)
		}

		var formatted string
		switch f := formatFunc.(type) {
		case FormatPackageUpdatesFunc:
			if typedData, ok := any(data).([]metrics.PackageUpdate); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.PackageUpdate for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatDiskDetailsFunc:
			if typedData, ok := any(data).([]metrics.DiskUsage); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.DiskUsage for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatCPULoadFunc:
			if typedData, ok := any(data).(*metrics.CPULoad); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to *metrics.CPULoad for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatMemoryDetailsFunc:
			if typedData, ok := any(data).(*metrics.MemoryData); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to *metrics.MemoryData for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatActiveSSHSessionsFunc:
			if typedData, ok := any(data).([]metrics.ActiveSession); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.ActiveSession for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatPreviousSSHSessionsFunc:
			if typedData, ok := any(data).([]metrics.PreviousSession); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.PreviousSession for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatNetworkDetailsFunc:
			if typedData, ok := any(data).([]metrics.NetworkInterface); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.NetworkInterface for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatCrowdSecAlertsFunc:
			if typedData, ok := any(data).([]metrics.Alert); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.Alert for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		case FormatCrowdSecDecisionsFunc:
			if typedData, ok := any(data).([]metrics.Decision); ok {
				formatted, err = f(typedData)
			} else {
				logger.Error("Failed to cast data to []metrics.Decision for metric: %s", metricName)
				return fmt.Sprintf("<p>Unable to format %s: invalid data type</p>", metricName)
			}
		default:
			logger.Error("Unknown FormatFunc type for metric: %s", metricName)
			return fmt.Sprintf("<p>Unable to format %s: unknown format function type</p>", metricName)
		}

		if err != nil {
			logger.Error("Failed to format %s: %v", metricName, err)
			return fmt.Sprintf("<p>Unable to format %s: %v</p>", metricName, err)
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
