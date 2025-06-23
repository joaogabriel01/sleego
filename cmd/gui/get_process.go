package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joaogabriel01/sleego"
)

func buildGetProcessContent(monitor sleego.ProcessorMonitor) fyne.CanvasObject {

	listContainer := UpdateProcessList(monitor)

	scroll := container.NewVScroll(listContainer)

	scroll.SetMinSize(fyne.NewSize(300, 200))

	refreshButton := widget.NewButton("Refresh", func() {

		newList := UpdateProcessList(monitor)
		scroll.Content = newList
		scroll.Refresh()
	})

	return container.NewBorder(
		widget.NewLabelWithStyle("Process List", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		refreshButton,
		nil, nil,
		scroll,
	)
}

func UpdateProcessList(monitor sleego.ProcessorMonitor) *fyne.Container {
	processes, err := monitor.GetRunningProcesses()
	if err != nil {
		return container.NewVBox(
			widget.NewLabel("Error fetching processes: " + err.Error()),
		)
	}

	var processList []fyne.CanvasObject
	for _, process := range processes {
		processInfo, err := process.GetInfo()
		if err != nil {
			continue
		}

		processList = append(processList, widget.NewLabel(processInfo.Name+" (PID: "+fmt.Sprintf("%d", processInfo.Pid)+")"))
	}

	return container.NewVBox(processList...)

}
