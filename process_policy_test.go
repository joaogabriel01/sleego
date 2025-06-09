package sleego

import (
	"context"
	"testing"
	"time"
)

// *************** MOCKS *************** //

// MockProcess is a mock implementation of Process
type MockProcess struct {
	info   ProcessInfo
	killed bool
}

func (p *MockProcess) GetInfo() (ProcessInfo, error) {
	return p.info, nil
}

func (p *MockProcess) Kill() error {
	p.killed = true
	return nil
}

// MockProcessorMonitor is a mock implementation of ProcessorMonitor
type MockProcessorMonitor struct {
	processes []Process
}

func (m *MockProcessorMonitor) GetRunningProcesses() ([]Process, error) {
	return m.processes, nil
}

type MockCategoryOperator struct{}

func (m *MockCategoryOperator) GetCategoriesOf(processName string) []string {
	return []string{"MockCategory"}
}

func (m *MockCategoryOperator) SetProcessByCategories(categoriesByProcess map[string][]string) {
	// Mock implementation does nothing
}

var mockCategoryOperator = &MockCategoryOperator{}

// *************** TESTS *************** //

// TestEnforceProcessPolicy tests the Apply method
func TestEnforceProcessPolicy_KillProcessAndVerifyChannel(t *testing.T) {
	mockProcess := &MockProcess{
		info: ProcessInfo{
			Name: "Notepad",
			Pid:  1234,
		},
	}

	mockMonitor := &MockProcessorMonitor{
		processes: []Process{mockProcess},
	}

	appsConfig := []AppConfig{
		{
			Name:        "Notepad",
			AllowedFrom: "09:00",
			AllowedTo:   "17:00",
		},
	}

	mockNow := func() time.Time {
		return time.Date(2023, 10, 10, 18, 0, 0, 0, time.UTC) // 18:00 UTC on October 10, 2023
	}

	ch := make(chan string, 1)
	policy := NewProcessPolicyImpl(mockMonitor, mockCategoryOperator, mockNow, ch)
	policy.enforceProcessPolicy(appsConfig)

	select {
	case alert := <-ch:
		if alert != "Killing process: Notepad, PID: 1234" {
			t.Errorf("Expected alert to be 'Killing process: Notepad', got %s", alert)
		}
	default:
		t.Errorf("Expected alert to be sent, but it was not")
	}

	if !mockProcess.killed {
		t.Errorf("Expected process to be killed, but it was not")
	}
}

func TestEnforceProcessPolicy_DoNotKillProcess(t *testing.T) {
	mockProcess := &MockProcess{
		info: ProcessInfo{
			Name: "Calculator",
			Pid:  5678,
		},
	}

	mockMonitor := &MockProcessorMonitor{
		processes: []Process{mockProcess},
	}

	appsConfig := []AppConfig{
		{
			Name:        "Calculator",
			AllowedFrom: "00:00",
			AllowedTo:   "23:59",
		},
	}

	mockNow := func() time.Time {
		return time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC) // 12:00 UTC on October 10, 2023
	}

	policy := NewProcessPolicyImpl(mockMonitor, mockCategoryOperator, mockNow, nil)
	policy.enforceProcessPolicy(appsConfig)

	if mockProcess.killed {
		t.Errorf("Expected process not to be killed, but it was")
	}
}

func TestEnforceProcessPolicy_ProcessNotInConfig(t *testing.T) {
	mockProcess := &MockProcess{
		info: ProcessInfo{
			Name: "UnknownApp",
			Pid:  9999,
		},
	}

	mockMonitor := &MockProcessorMonitor{
		processes: []Process{mockProcess},
	}

	appsConfig := []AppConfig{
		{
			Name:        "KnownApp",
			AllowedFrom: "09:00",
			AllowedTo:   "17:00",
		},
	}

	mockNow := func() time.Time {
		return time.Date(2023, 10, 10, 15, 0, 0, 0, time.UTC)
	}

	policy := NewProcessPolicyImpl(mockMonitor, categoryOperator, mockNow, nil)
	policy.enforceProcessPolicy(appsConfig)

	if mockProcess.killed {
		t.Errorf("Process not in config should not be killed")
	}
}

func TestIsAllowedToRun_InvalidTimeFormat(t *testing.T) {
	appConfig := AppConfig{
		AllowedFrom: "invalid",
		AllowedTo:   "invalid",
	}

	policy := NewProcessPolicyImpl(nil, nil, nil, nil)
	result := policy.isAllowedToRun(appConfig)

	if result {
		t.Errorf("Expected isAllowedToRun to return false for invalid time format")
	}
}

func TestIsAllowedToRun_EmptyTimes(t *testing.T) {
	appConfig := AppConfig{
		AllowedFrom: "",
		AllowedTo:   "",
	}

	policy := NewProcessPolicyImpl(nil, nil, nil, nil)
	result := policy.isAllowedToRun(appConfig)

	if result {
		t.Errorf("Expected isAllowedToRun to return false for empty times")
	}
}

func TestApply(t *testing.T) {
	mockProcess := &MockProcess{
		info: ProcessInfo{
			Name: "BlockedApp",
			Pid:  4321,
		},
	}

	mockMonitor := &MockProcessorMonitor{
		processes: []Process{mockProcess},
	}

	appsConfig := []AppConfig{
		{
			Name:        "BlockedApp",
			AllowedFrom: "09:00",
			AllowedTo:   "17:00",
		},
	}

	mockNow := func() time.Time {
		return time.Date(2023, 10, 10, 18, 0, 0, 0, time.UTC)
	}

	policy := NewProcessPolicyImpl(mockMonitor, nil, mockNow, nil)
	ctx := context.Background()
	go policy.Apply(ctx, appsConfig)

	time.Sleep(100 * time.Millisecond) // Wait for enforcement

	if !mockProcess.killed {
		t.Errorf("Apply should have killed the process")
	}
}

func TestApply_ProcessAllowed(t *testing.T) {
	mockProcess := &MockProcess{
		info: ProcessInfo{
			Name: "AllowedApp",
			Pid:  8765,
		},
	}

	mockMonitor := &MockProcessorMonitor{
		processes: []Process{mockProcess},
	}

	appsConfig := []AppConfig{
		{
			Name:        "AllowedApp",
			AllowedFrom: "00:00",
			AllowedTo:   "23:59",
		},
	}

	mockNow := func() time.Time {
		return time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
	}

	policy := NewProcessPolicyImpl(mockMonitor, nil, mockNow, nil)
	ctx := context.Background()
	go policy.Apply(ctx, appsConfig)

	time.Sleep(100 * time.Millisecond) // Wait for enforcement

	if mockProcess.killed {
		t.Errorf("Apply should not have killed the process")
	}
}
func TestIsAllowedToRun(t *testing.T) {
	tests := []struct {
		name      string
		appConfig AppConfig
		mockNow   time.Time
		expected  bool
	}{
		{
			name: "Within allowed time",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC), // 12:00 UTC on October 10, 2023
			expected: true,
		},
		{
			name: "Before allowed time",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 8, 59, 59, 0, time.UTC), // 08:59:59 UTC on October 10, 2023
			expected: false,
		},
		{
			name: "After allowed time",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 17, 0, 1, 0, time.UTC), // 17:00:01 UTC on October 10, 2023
			expected: false,
		},
		{
			name: "Invalid AllowedFrom format",
			appConfig: AppConfig{
				AllowedFrom: "invalid",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC), // 10:00 UTC on October 10, 2023
			expected: false,
		},
		{
			name: "Invalid AllowedTo format",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "invalid",
			},
			mockNow:  time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC), // 10:00 UTC on October 10, 2023
			expected: false,
		},
		{
			name: "Exactly at AllowedFrom",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 9, 0, 0, 0, time.UTC), // 09:00 UTC on October 10, 2023
			expected: true,
		},
		{
			name: "Exactly at AllowedTo",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 17, 0, 0, 0, time.UTC), // 17:00 UTC on October 10, 2023
			expected: true,
		},
		{
			name: "AllowedFrom equals AllowedTo",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "09:00",
			},
			mockNow:  time.Date(2023, 10, 10, 9, 0, 0, 0, time.UTC), // 09:00 UTC on October 10, 2023
			expected: true,
		},
		{
			name: "AllowedFrom later than AllowedTo (overnight)",
			appConfig: AppConfig{
				AllowedFrom: "22:00",
				AllowedTo:   "06:00",
			},
			mockNow:  time.Date(2023, 10, 10, 23, 0, 0, 0, time.UTC), // 23:00 UTC on October 10, 2023
			expected: true,
		},
		{
			name: "AllowedFrom later than AllowedTo (overnight) outside allowed time",
			appConfig: AppConfig{
				AllowedFrom: "22:00",
				AllowedTo:   "06:00",
			},
			mockNow:  time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC), // 12:00 UTC on October 10, 2023
			expected: false,
		},
		{
			name: "Empty AllowedFrom",
			appConfig: AppConfig{
				AllowedFrom: "",
				AllowedTo:   "17:00",
			},
			mockNow:  time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC), // 10:00 UTC on October 10, 2023
			expected: false,
		},
		{
			name: "Empty AllowedTo",
			appConfig: AppConfig{
				AllowedFrom: "09:00",
				AllowedTo:   "",
			},
			mockNow:  time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC), // 10:00 UTC on October 10, 2023
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := NewProcessPolicyImpl(nil, nil, func() time.Time {
				return tt.mockNow
			}, nil)
			result := policy.isAllowedToRun(tt.appConfig)
			if result != tt.expected {
				t.Errorf("isAllowedToRun() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestNewProcessPolicyImpl tests the NewProcessPolicyImpl constructor
func TestNewProcessPolicyImpl(t *testing.T) {
	mockMonitor := &MockProcessorMonitor{}
	policy := NewProcessPolicyImpl(mockMonitor, nil, nil, nil)

	if policy.monitor != mockMonitor {
		t.Errorf("Expected monitor to be %v, got %v", mockMonitor, policy.monitor)
	}

	if policy.now == nil {
		t.Errorf("Expected now function to be set, got nil")
	}
}
