package sleego

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/joaogabriel01/sleego/internal/logger"
)

// ShutdownPolicy defines the behavior for shutting down the system
type ShutdownPolicy interface {
	Apply(ctx context.Context, endTime time.Time) error
}

type ShutdownPolicyImpl struct {
	shutdown     func() error
	c            chan string
	timesToAlert []int
	logger       logger.Logger
}

func NewShutdownPolicyImpl(c chan string, timesToAlert []int) ShutdownPolicy {
	logger, err := logger.Get()
	if err != nil {
		panic(fmt.Sprintf("failed to get logger: %v", err))
	}

	return &ShutdownPolicyImpl{
		shutdown:     shutdown,
		c:            c,
		timesToAlert: timesToAlert,
		logger:       logger,
	}
}

// Apply schedules a shutdown at the specified time.
func (s *ShutdownPolicyImpl) Apply(ctx context.Context, endTime time.Time) error {
	now := time.Now()
	shutdownTime := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, now.Location())

	if shutdownTime.Before(now) {
		shutdownTime = shutdownTime.Add(24 * time.Hour)
	}

	duration := time.Until(shutdownTime)
	if duration <= 0 {
		s.logger.Debug("Shutting down now")
		return s.shutdown()
	}

	s.logger.Info(fmt.Sprintf("Shutting down scheduled in %v", duration))

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for _, timeToAlert := range s.timesToAlert {
		alertDuration := duration - time.Duration(timeToAlert)*time.Minute
		if alertDuration > 0 {
			timeToAlert := timeToAlert
			go func() {
				select {
				case <-ctx.Done():
				case <-time.After(alertDuration):
					msg := fmt.Sprintf("Shutting down in %d minutes", timeToAlert)
					s.logger.Debug(msg)
					s.c <- msg
				}
			}()
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		s.logger.Debug("Shutting down now")
		return s.shutdown()
	}
}

var _ ShutdownPolicy = &ShutdownPolicyImpl{}

func shutdown() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("shutdown", "/s", "/f", "/t", "0")
	case "linux":
		cmd = exec.Command("shutdown", "-h", "now")
	case "darwin":
		cmd = exec.Command("sudo", "shutdown", "-h", "now")
	default:
		return errors.New("unsupported operating system")
	}
	return cmd.Run()
}
