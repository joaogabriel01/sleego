package sleego

import (
	"os"
	"testing"

	"github.com/shirou/gopsutil/v4/process"
)

func TestProcessImpl_GetInfo(t *testing.T) {
	// Get the current process
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatalf("Failed to create process: %v", err)
	}

	p := &ProcessImpl{proc: proc}

	info, err := p.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() returned error: %v", err)
	}

	if info.Pid != os.Getpid() {
		t.Errorf("Expected Pid %d, got %d", os.Getpid(), info.Pid)
	}

	name, err := proc.Name()
	if err != nil {
		t.Fatalf("proc.Name() returned error: %v", err)
	}

	if info.Name != name {
		t.Errorf("Expected Name %s, got %s", name, info.Name)
	}
}

func TestProcessImpl_Kill(t *testing.T) {
	// Note: Killing processes during tests is not recommended.
	// This test is being skipped for safety.
	t.Skip("TestProcessImpl_Kill is being skipped to avoid killing processes during tests")
}

func TestProcessorMonitorImpl_GetRunningProcesses(t *testing.T) {
	monitor := &ProcessorMonitorImpl{}

	processes, err := monitor.GetRunningProcesses()
	if err != nil {
		t.Fatalf("GetRunningProcesses() returned error: %v", err)
	}

	if len(processes) == 0 {
		t.Errorf("Expected some processes, but none were found")
	}

	for _, p := range processes {
		info, err := p.GetInfo()
		t.Log(info)
		if err != nil {
			t.Errorf("GetInfo() returned error: %v", err)
			continue
		}
	}
}
