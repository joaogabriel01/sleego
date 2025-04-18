package sleego

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/joaogabriel01/sleego/internal/logger"
)

type MockShutdown struct {
	called bool
}

func (m *MockShutdown) Shutdown() error {
	m.called = true
	return nil
}

var ctxOk = context.Background()

func TestShutdownPolicyImpl_Apply_ShutdownCalled(t *testing.T) {
	mockShutdown := &MockShutdown{}
	alertCh := make(chan string, 1)
	timesToAlert := []int{1}

	policy := &ShutdownPolicyImpl{
		shutdown: func() error {
			return mockShutdown.Shutdown()
		},
		c:            alertCh,
		timesToAlert: timesToAlert,
		logger:       logger.NewLoggerMock(),
	}

	endTime := time.Now().Add(3 * time.Second)
	endTimeStr := endTime.Format("15:04:05")
	t.Log(endTimeStr)
	endTimeParsed, _ := time.Parse("15:04:05", endTimeStr)
	go func() {
		err := policy.Apply(ctxOk, endTimeParsed)
		if err != nil {
			t.Errorf("Apply returned error: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	if !mockShutdown.called {
		t.Errorf("Expected shutdown to be called, but it was not")
	}
}

func TestShutdownPolicyImpl_Apply_AlertSent(t *testing.T) {
	mockShutdown := &MockShutdown{}
	alertCh := make(chan string, 1)
	timesToAlert := []int{1} // 1 minute before shutdown

	policy := &ShutdownPolicyImpl{
		shutdown: func() error {
			return mockShutdown.Shutdown()
		},
		c:            alertCh,
		timesToAlert: timesToAlert,
		logger:       logger.NewLoggerMock(),
	}

	// Set endTime to 2 minutes from now
	endTime := time.Now().Add(2 * time.Minute)
	go func() {
		err := policy.Apply(ctxOk, endTime)
		if err != nil {
			t.Errorf("Apply returned error: %v", err)
		}
	}()

	// Wait for alert
	select {
	case msg := <-alertCh:
		expectedMsg := "Shutting down in 1 minutes"
		if msg != expectedMsg {
			t.Errorf("Expected alert message '%s', got '%s'", expectedMsg, msg)
		}
	case <-time.After(3 * time.Minute):
		t.Errorf("Did not receive expected alert message")
	}
}

func TestShutdownPolicyImpl_Apply_ShutdownError(t *testing.T) {
	mockShutdown := &MockShutdown{}
	alertCh := make(chan string, 1)
	timesToAlert := []int{1}

	policy := &ShutdownPolicyImpl{
		shutdown: func() error {
			return errors.New("shutdown failed")
		},
		c:            alertCh,
		timesToAlert: timesToAlert,
		logger:       logger.NewLoggerMock(),
	}

	endTime := time.Now().Add(2 * time.Second)
	go func() {
		err := policy.Apply(ctxOk, endTime)
		if err == nil {
			t.Errorf("Expected error from shutdown, but got none")
		}
	}()

	time.Sleep(3 * time.Second)

	if mockShutdown.called {
		t.Errorf("Shutdown should not have been called due to error")
	}
}

func TestShutdownPolicyImpl_ContextCancelled(t *testing.T) {
	mockShutdown := &MockShutdown{}
	alertCh := make(chan string, 1)
	timesToAlert := []int{1}

	policy := &ShutdownPolicyImpl{
		shutdown: func() error {
			return mockShutdown.Shutdown()
		},
		c:            alertCh,
		timesToAlert: timesToAlert,
		logger:       logger.NewLoggerMock(),
	}
	ctx, cancel := context.WithCancel(ctxOk)
	cancel()
	endTime := time.Now().Add(2 * time.Second)
	go func() {
		err := policy.Apply(ctx, endTime)
		if err == nil {
			t.Errorf("Expected error from context cancellation, but got none")
		}
	}()

	// Cancel context
	time.Sleep(1 * time.Second)
}
