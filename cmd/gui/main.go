package main

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/joaogabriel01/sleego"
)

type AppEntry struct {
	nameEntry        *widget.Entry
	allowedFromEntry *widget.Entry
	allowedToEntry   *widget.Entry
	container        *fyne.Container
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := app.New()
	w := a.NewWindow("Configuration")
	screenSize := w.Canvas().Size()
	if screenSize.Width == 0 || screenSize.Height == 0 {
		screenSize = fyne.NewSize(800, 600)
	}
	w.Resize(screenSize)

	configPath := "./config.json"
	loader := sleego.Loader{}

	appConfigs, err := loader.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	var entries []AppEntry
	appsContainer := container.NewVBox()

	removeEntry := func(index int) {
		entry := entries[index]
		appsContainer.Remove(entry.container)
		entries = append(entries[:index], entries[index+1:]...)
	}

	addApp := func(config sleego.AppConfig) {
		nameEntry := widget.NewEntry()
		nameEntry.SetText(config.Name)

		allowedFromEntry := widget.NewEntry()
		allowedFromEntry.SetText(config.AllowedFrom)

		allowedToEntry := widget.NewEntry()
		allowedToEntry.SetText(config.AllowedTo)

		index := len(entries)
		removeBtn := widget.NewButton("Remove", func() {
			removeEntry(index)
		})

		form := container.NewVBox(
			widget.NewLabel("Name:"),
			nameEntry,
			widget.NewLabel("Allowed From:"),
			allowedFromEntry,
			widget.NewLabel("Allowed To:"),
			allowedToEntry,
			removeBtn,
			widget.NewSeparator(),
		)

		appsContainer.Add(form)
		spacer := widget.NewLabel("")
		spacer.Resize(fyne.NewSize(0, 10))
		appsContainer.Add(spacer)

		entries = append(entries, AppEntry{
			nameEntry:        nameEntry,
			allowedFromEntry: allowedFromEntry,
			allowedToEntry:   allowedToEntry,
			container:        form,
		})
	}

	for _, config := range appConfigs {
		addApp(config)
	}

	addButton := widget.NewButton("Add", func() {
		addApp(sleego.AppConfig{})
	})

	saveButton := widget.NewButton("Save", func() {
		updatedConfigs := make([]sleego.AppConfig, len(entries))
		for i, entry := range entries {
			updatedConfigs[i] = sleego.AppConfig{
				Name:        entry.nameEntry.Text,
				AllowedFrom: entry.allowedFromEntry.Text,
				AllowedTo:   entry.allowedToEntry.Text,
			}
		}

		if err := loader.Save(configPath, updatedConfigs); err != nil {
			log.Printf("Error saving the configuration file: %v", err)
			return
		}
		appConfigs = updatedConfigs
		dialog.ShowInformation("Success", "Configurations saved successfully!", w)
		log.Println("Configurations saved successfully!")
	})

	runButton := widget.NewButton("Run", func() {
		cancel()
		ctx, cancel = context.WithCancel(context.Background())
		monitor := &sleego.ProcessorMonitorImpl{}
		policy := sleego.NewProcessPolicyImpl(monitor, nil)
		log.Printf("Starting process policy with config: %+v of path: %s", appConfigs, configPath)
		dialog.ShowInformation("Running", "Applying the policy...", w)

		go policy.Apply(ctx, appConfigs)
	})

	mainContainer := container.NewVBox(
		appsContainer,
		addButton,
		saveButton,
		runButton,
	)

	w.SetContent(container.NewScroll(mainContainer))
	w.ShowAndRun()
}
