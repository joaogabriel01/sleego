package sleego

import (
	"context"
	"fmt"
	"time"

	"github.com/joaogabriel01/sleego/internal/logger"
)

// Time to sleep between checks
const sleepTime = 5 * time.Second

// ProcessPolicyImpl is the implementation of the ProcessPolicy interface
type ProcessPolicyImpl struct {
	monitor          ProcessorMonitor
	categoryOperator CategoryOperator
	now              func() time.Time
	alertCh          chan string
	logger           logger.Logger
}

// NewProcessPolicyImpl creates a new ProcessPolicyImpl
func NewProcessPolicyImpl(monitor ProcessorMonitor, categoryOperator CategoryOperator, now func() time.Time, alert chan string) *ProcessPolicyImpl {
	if now == nil {
		now = time.Now
	}
	logger, err := logger.Get()
	if err != nil {
		panic(fmt.Sprintf("failed to get logger: %v", err))
	}
	return &ProcessPolicyImpl{monitor: monitor, categoryOperator: categoryOperator, now: now, alertCh: alert, logger: logger}
}

// Apply will check the running processes and kill the ones that are not allowed to run
func (p *ProcessPolicyImpl) Apply(ctx context.Context, appsConfig []AppConfig) error {
	for {
		if ctx.Err() != nil {
			p.logger.Debug("Context cancelled, stopping process policy")
			return nil
		}
		p.enforceProcessPolicy(appsConfig)
		time.Sleep(sleepTime)
	}
}

func (p *ProcessPolicyImpl) enforceProcessPolicy(appsConfig []AppConfig) {
	processes, err := p.monitor.GetRunningProcesses()
	if err != nil {
		p.logger.Error(fmt.Sprintf("Error getting running processes: %v", err))
		return
	}

	for _, process := range processes {
		info, err := process.GetInfo()
		if err != nil {
			p.logger.Error(fmt.Sprintf("Error getting process info: %v", err))
			continue
		}

		p.logger.Debug(fmt.Sprintf("Checking process: %s, PID: %d", info.Name, info.Pid))
		for _, appConfig := range appsConfig {
			if info.Name == appConfig.Name || (p.categoryOperator != nil && existElementInSlice(p.categoryOperator.GetCategoriesOf(info.Name), appConfig.Name)) {

				// Check if the process is running outside the allowed hours
				if !p.isAllowedToRun(appConfig) {
					msg := fmt.Sprintf("Killing process: %s, PID: %d", info.Name, info.Pid)
					if p.alertCh != nil {
						p.alertCh <- msg
					}
					p.logger.Info(msg)
					err := process.Kill()
					if err != nil {
						p.logger.Error(fmt.Sprintf("Error killing process: %v", err))
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
		p.logger.Error("Error parsing allowedFrom: " + err.Error())
		return false
	}
	allowedTo, err := time.Parse("15:04", appConfig.AllowedTo)
	if err != nil {
		p.logger.Error("Error parsing allowedTo: " + err.Error())
		return false
	}
	allowedFrom = time.Date(now.Year(), now.Month(), now.Day(), allowedFrom.Hour(), allowedFrom.Minute(), 0, 0, now.Location())
	allowedTo = time.Date(now.Year(), now.Month(), now.Day(), allowedTo.Hour(), allowedTo.Minute(), 0, 0, now.Location())

	p.logger.Debug(fmt.Sprintf("AllowedFrom: %s, AllowedTo: %s, Now: %s", allowedFrom.Format("15:04"), allowedTo.Format("15:04"), now.Format("15:04")))

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

func existElementInSlice(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

var _ ProcessPolicy = &ProcessPolicyImpl{}
