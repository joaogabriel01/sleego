package sleego

import (
	"log"
	"time"
)

// Time to sleep between checks
const sleepTime = 5 * time.Second

// ProcessPolicyImpl is the implementation of the ProcessPolicy interface
type ProcessPolicyImpl struct {
	monitor ProcessorMonitor
	now     func() time.Time
}

// NewProcessPolicyImpl creates a new ProcessPolicyImpl
func NewProcessPolicyImpl(monitor ProcessorMonitor, now func() time.Time) *ProcessPolicyImpl {
	if now == nil {
		now = time.Now
	}
	return &ProcessPolicyImpl{monitor: monitor, now: now}
}

// Apply will check the running processes and kill the ones that are not allowed to run
func (p *ProcessPolicyImpl) Apply(appsConfig []AppConfig) error {
	for {
		p.enforceProcessPolicy(appsConfig)
		time.Sleep(sleepTime)
	}
}

func (p *ProcessPolicyImpl) enforceProcessPolicy(appsConfig []AppConfig) {
	processes, err := p.monitor.GetRunningProcesses()
	if err != nil {
		log.Println("Error getting running processes:", err)
		return
	}

	for _, process := range processes {
		info, err := process.GetInfo()
		if err != nil {
			log.Println("Error getting process info:", err)
			continue
		}
		for _, appConfig := range appsConfig {
			if info.Name == appConfig.Name {
				// Check if the process is running outside the allowed hours
				if !p.isAllowedToRun(appConfig) {
					log.Println("Killing process:", info.Name, "PID:", info.Pid)
					err := process.Kill()
					if err != nil {
						log.Printf("Error killing process %s: %v", info.Name, err)
						continue
					}
				}
			}
		}
	}
}

func (p *ProcessPolicyImpl) isAllowedToRun(appConfig AppConfig) bool {
	now := p.now()
	allowedFrom, err := time.Parse("15:04", appConfig.AllowedFrom)
	if err != nil {
		log.Println("Error parsing allowedFrom:", err)
		return false
	}
	allowedTo, err := time.Parse("15:04", appConfig.AllowedTo)
	if err != nil {
		log.Println("Error parsing allowedTo:", err)
		return false
	}
	allowedFrom = time.Date(now.Year(), now.Month(), now.Day(), allowedFrom.Hour(), allowedFrom.Minute(), 0, 0, now.Location())
	allowedTo = time.Date(now.Year(), now.Month(), now.Day(), allowedTo.Hour(), allowedTo.Minute(), 0, 0, now.Location())
	log.Println("AllowedFrom:", allowedFrom, "AllowedTo:", allowedTo, "Now:", now)
	if allowedFrom.After(allowedTo) {
		// If AllowedFrom is later than AllowedTo, it means the app is allowed to run overnight
		// So we need to check if the current time is outside the allowed time
		if now.Before(allowedFrom) && now.After(allowedTo) {
			return false
		}
	} else {
		// If AllowedFrom is earlier than AllowedTo, it means the app is allowed to run during the day
		// So we need to check if the current time is outside the allowed time
		if now.Before(allowedFrom) || now.After(allowedTo) {
			return false
		}
	}
	return true
}
