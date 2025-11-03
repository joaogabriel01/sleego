package sleego

import (
	"github.com/shirou/gopsutil/v4/process"
)

// Process defines the behavior of a process
type Process interface {
	GetInfo() (ProcessInfo, error)
	Kill() error
}

// ProcessorMonitor will be used to interact with the system processes
type ProcessorMonitor interface {
	GetRunningProcesses() ([]Process, error)
}

// ProcessInfo contains the information of a process
type ProcessInfo struct {
	Name     string
	Pid      int
	Category []string
}

// This is the adapter to the Process interface from the gopsutil library
type ProcessImpl struct {
	proc             *process.Process
	categoryOperator CategoryOperator
}

func (p *ProcessImpl) GetInfo() (ProcessInfo, error) {
	name, err := p.proc.Name()
	if err != nil {
		return ProcessInfo{}, err
	}
	pid := int(p.proc.Pid)
	if p.categoryOperator != nil {
		categories := p.categoryOperator.GetCategoriesOf(name)
		if len(categories) != 0 {
			return ProcessInfo{Name: name, Pid: pid, Category: categories}, nil
		}
	}

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
		processes = append(processes, newProcessWithCategoryOperator(proc))
	}
	return processes, nil
}

func newProcessWithCategoryOperator(proc *process.Process) *ProcessImpl {
	return &ProcessImpl{
		proc:             proc,
		categoryOperator: GetCategoryOperator(),
	}
}

var _ Process = &ProcessImpl{}
var _ ProcessorMonitor = &ProcessorMonitorImpl{}
