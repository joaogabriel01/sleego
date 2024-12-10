package sleego

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

type ShutdownPolicyImpl struct {
	shutdown     func() error
	c            chan string
	timesToAlert []int
}

func NewShutdownPolicyImpl(c chan string, timesToAlert []int) ShutdownPolicy {
	return &ShutdownPolicyImpl{
		shutdown:     shutdown,
		c:            c,
		timesToAlert: timesToAlert,
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
		log.Println("Shutting down now")
		return s.shutdown()
	}

	log.Printf("Shutting down scheduled in %v", duration)

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
					log.Println(msg)
					s.c <- msg
				}
			}()
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		log.Println("Shutting down now")
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
