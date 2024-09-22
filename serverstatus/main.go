package main

import (
    "fmt"
    "os"
    "time"

    "serverstatus/email"
    "serverstatus/metrics"
)

func main() {
    // Server details here
    serverHostname := "Lightsail Web"
    serverTime := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
    serverUptime := metrics.GetUptime()

    // Metrics collection
    packageUpdates := metrics.GetPackageUpdates()
    diskDetails := metrics.GetDiskDetails()
    cpuLoadDetails := metrics.GetCPULoadDetails()
    memoryDetails := metrics.GetMemoryDetails()
    activeSSH := metrics.GetActiveSSHSessions()
    previousSSH := metrics.GetPreviousSSHSessions()
    networkDetails := metrics.GetNetworkDetails()
    crowdSecAlerts := metrics.GetCrowdSecAlerts()
    crowdSecDecisions := metrics.GetCrowdSecDecisions()

    // Construct Email Body
    emailBody := email.ConstructEmailBody(email.EmailData{
        ServerHostname:  serverHostname,
        ServerTime:      serverTime,
        ServerUptime:    serverUptime,
        PackageUpdates:  packageUpdates,
        DiskDetails:     diskDetails,
        MemoryDetails:   memoryDetails,
        CPULoadDetails:  cpuLoadDetails,
        ActiveSSH:       activeSSH,
        PreviousSSH:     previousSSH,
        NetworkDetails:  networkDetails,
        CrowdSecAlerts:  crowdSecAlerts,
        CrowdSecDecisions: crowdSecDecisions,
    })

    // Send Email
    today := time.Now().Format("2006-01-02")
    subject := fmt.Sprintf("Daily System Report, %s", today)
    err := email.SendEmail(subject, emailBody, "webmaster@timothywb.com")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to send email: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Email sent successfully.")
}
