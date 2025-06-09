package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joaogabriel01/sleego"
)

func buildAppContent() fyne.CanvasObject {
	appsContainer = container.NewVBox()
	for _, cfg := range appConfigs {
		addApp(cfg)
	}

	addButton := widget.NewButton("Add New App", func() { addApp(sleego.AppConfig{}) })

	return container.NewVBox(
		appsContainer,
		addButton,
		widget.NewSeparator(),
	)
}

func createResizableContainer(object fyne.CanvasObject, width, height float32) *fyne.Container {
	object.Resize(fyne.NewSize(width, height))
	return container.NewWithoutLayout(object)
}

func addApp(config sleego.AppConfig) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(config.Name)

	allowedFromEntry := widget.NewEntry()
	allowedFromEntry.SetText(config.AllowedFrom)

	allowedToEntry := widget.NewEntry()
	allowedToEntry.SetText(config.AllowedTo)

	removeBtn := widget.NewButton("Remove", func() {
		if nameEntry.Text == "" {
			loggerInstance.Error("Cannot remove entry with empty name")
			return
		}
		removeEntry(nameEntry.Text)
		loggerInstance.Debug(fmt.Sprintf("Removed entry with name: %s", nameEntry.Text))
	})

	form := container.NewVBox(
		widget.NewLabel("Name:"),
		createResizableContainer(nameEntry, applicationsSize.Width*0.95, 36),
		widget.NewLabel("Allowed From:"),
		createResizableContainer(allowedFromEntry, applicationsSize.Width*0.95, 36),
		widget.NewLabel("Allowed To:"),
		createResizableContainer(allowedToEntry, applicationsSize.Width*0.95, 36),
		createResizableContainer(removeBtn, applicationsSize.Width*0.95, 40),
		widget.NewSeparator(),
	)

	appsContainer.Add(form)
	appsContainer.Add(widget.NewLabel(""))

	entries[nameEntry.Text] = AppEntry{
		nameEntry:        nameEntry,
		allowedFromEntry: allowedFromEntry,
		allowedToEntry:   allowedToEntry,
		container:        form,
	}

	nameEntry.OnChanged = func(newName string) {
		if _, exists := entries[config.Name]; exists {
			delete(entries, config.Name)
			loggerInstance.Debug(fmt.Sprintf("Updated app name from %s to %s", config.Name, newName))
		}
		config.Name = newName
		entries[newName] = AppEntry{
			nameEntry:        nameEntry,
			allowedFromEntry: allowedFromEntry,
			allowedToEntry:   allowedToEntry,
			container:        form,
		}
	}
}

func removeEntry(name string) {
	if entry, exists := entries[name]; exists {
		appsContainer.Remove(entry.container)
		delete(entries, name)
	} else {
		loggerInstance.Error(fmt.Sprintf("Entry with name %s does not exist", name))
	}
}
