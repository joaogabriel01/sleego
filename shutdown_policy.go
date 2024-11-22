package sleego

import (
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

func NewShutdownPolicy(c chan string, timesToAlert []int) ShutdownPolicy {
	return &ShutdownPolicyImpl{
		shutdown:     shutdown,
		c:            c,
		timesToAlert: timesToAlert,
	}
}

// Apply will shutdown the computer at the specified time.
func (s *ShutdownPolicyImpl) Apply(endTime time.Time) error {

	now := time.Now()

	endDateTime := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, now.Location())

	if endDateTime.Before(now) {
		endDateTime = endDateTime.Add(24 * time.Hour)
	}

	duration := endDateTime.Sub(now)

	for _, timeToAlert := range s.timesToAlert {
		durationToAlert := duration - time.Duration(timeToAlert)*time.Minute
		if durationToAlert > 0 {
			go func(t int) {
				time.Sleep(durationToAlert)
				msg := fmt.Sprintf("Shutting down in %v minutes", t)
				log.Println(msg)
				s.c <- msg
			}(timeToAlert)
		}
	}

	if duration > 0 {
		log.Printf("Shutting down in %v", duration)
		time.Sleep(duration)
	}

	log.Println("Shutting down now")
	return s.shutdown()
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
