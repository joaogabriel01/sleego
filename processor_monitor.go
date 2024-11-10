package sleego

import (
	"github.com/shirou/gopsutil/process"
)

// This is the adapter to the Process interface from the gopsutil library
type ProcessImpl struct {
	proc *process.Process
}

func (p *ProcessImpl) GetInfo() (ProcessInfo, error) {
	name, err := p.proc.Name()
	if err != nil {
		return ProcessInfo{}, err
	}
	pid := int(p.proc.Pid)
	return ProcessInfo{Name: name, Pid: pid}, nil
}

func (p *ProcessImpl) Kill() error {
	return p.proc.Kill()
}

// This is the adapter to the ProcessorMonitor interface from the gopsutil library
type ProcessorMonitorImpl struct {
}

func (p *ProcessorMonitorImpl) GetRunningProcesses() ([]Process, error) {
	procs, err := process.Processes()
	if err != nil {
		return []Process{}, err
	}
	processes := make([]Process, 0, len(procs))
	for _, proc := range procs {
		processes = append(processes, &ProcessImpl{proc: proc})
	}
	return processes, nil
}

var _ Process = &ProcessImpl{}
var _ ProcessorMonitor = &ProcessorMonitorImpl{}
