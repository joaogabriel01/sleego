package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/joaogabriel01/sleego"

	_ "embed"
)

type AppEntry struct {
	nameEntry        *widget.Entry
	allowedFromEntry *widget.Entry
	allowedToEntry   *widget.Entry
	container        *fyne.Container
}

var (
	a                 fyne.App
	w                 fyne.Window
	ctx               context.Context
	cancel            context.CancelFunc
	loader            sleego.Loader
	configPath        string
	fileConfig        sleego.FileConfig
	appConfigs        []sleego.AppConfig
	entries           []AppEntry
	appsContainer     *fyne.Container
	shutdownTimeEntry *widget.Entry

	//go:embed assets/sleego_icon.png
	iconData []byte
)

func main() {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	parseFlags()

	initializeApp()
	setupTrayIcon()
	setupWindow()
	loadConfigurations()
	createUI()
	w.ShowAndRun()
}

func parseFlags() {
	configUser, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}
	configPath = *flag.String("config", configUser+"/sleego/config.json", "Path to config file")
	flag.Parse()
}

func initializeApp() {
	a = app.NewWithID("sleego.gui")
	icon := fyne.NewStaticResource("icon.png", iconData)
	a.SetIcon(icon)
	w = a.NewWindow("Configuration")
}

func setupTrayIcon() {
	icon := w.Icon()
	trayMenu := fyne.NewMenu("",
		fyne.NewMenuItem("Open", func() {
			w.Show()
		}),
		fyne.NewMenuItem("Quit", func() {
			a.Quit()
		}),
	)
	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayMenu(trayMenu)
		desk.SetSystemTrayIcon(icon)
	}
}

func setupWindow() {
	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("Quit", func() {
			a.Quit()
		})),
	))
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	screenSize := w.Canvas().Size()
	if screenSize.Width == 0 || screenSize.Height == 0 {
		screenSize = fyne.NewSize(800, 600)
	}
	w.Resize(screenSize)
}

func loadConfigurations() {
	var err error
	fileConfig, err = loader.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}
	appConfigs = fileConfig.Apps
}

func createUI() {
	appsContainer = container.NewVBox()
	for _, config := range appConfigs {
		addApp(config)
	}

	addButton := widget.NewButton("Add", func() {
		addApp(sleego.AppConfig{})
	})

	shutdownTimeEntry = widget.NewEntry()
	shutdownTimeEntry.SetPlaceHolder("Enter shutdown time (HH:MM)")
	if fileConfig.Shutdown != "" {
		shutdownTimeEntry.SetText(fileConfig.Shutdown)
	}

	saveButton := widget.NewButton("Save", saveConfigurations)
	runButton := widget.NewButton("Run", runPolicies)

	mainContainer := container.NewVBox(
		appsContainer,
		addButton,
		saveButton,
		shutdownTimeEntry,
		runButton,
	)

	w.SetContent(container.NewScroll(mainContainer))
}

func addApp(config sleego.AppConfig) {
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
	appsContainer.Add(widget.NewLabel(""))

	entries = append(entries, AppEntry{
		nameEntry:        nameEntry,
		allowedFromEntry: allowedFromEntry,
		allowedToEntry:   allowedToEntry,
		container:        form,
	})
}

func removeEntry(index int) {
	entry := entries[index]
	appsContainer.Remove(entry.container)
	entries = append(entries[:index], entries[index+1:]...)
}

func saveConfigurations() {
	updatedAppConfigs := make([]sleego.AppConfig, len(entries))
	for i, entry := range entries {
		updatedAppConfigs[i] = sleego.AppConfig{
			Name:        entry.nameEntry.Text,
			AllowedFrom: entry.allowedFromEntry.Text,
			AllowedTo:   entry.allowedToEntry.Text,
		}
	}
	shutdownConfig := shutdownTimeEntry.Text

	updatedConfigs := sleego.FileConfig{
		Apps:     updatedAppConfigs,
		Shutdown: shutdownConfig,
	}

	if err := loader.Save(configPath, updatedConfigs); err != nil {
		log.Printf("Error saving the configuration file: %v", err)
		return
	}
	fileConfig = updatedConfigs
	dialog.ShowInformation("Success", "Configurations saved successfully!", w)
	log.Println("Configurations saved successfully!")
}

func runPolicies() {
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	monitor := &sleego.ProcessorMonitorImpl{}
	processChan := make(chan string)
	policy := sleego.NewProcessPolicyImpl(monitor, nil, processChan)
	log.Printf("Starting process policy with config: %+v of path: %s", fileConfig.Apps, configPath)
	dialog.ShowInformation("Running", "Applying the policy...", w)

	go func() {
		if err := policy.Apply(ctx, fileConfig.Apps); err != nil {
			log.Printf("Error applying the policy: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-processChan:
				a.SendNotification(fyne.NewNotification("Alert", msg))
			case <-ctx.Done():
				return
			}
		}
	}()

	shutdownTimeStr := shutdownTimeEntry.Text
	shutdownTime, err := time.Parse("15:04", shutdownTimeStr)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid shutdown time format"), w)
		return
	}
	log.Printf("Starting shutdown policy with time: %s", shutdownTimeStr)
	channelToAlert := make(chan string)
	alerts := []int{10, 3, 1}
	shutdownPolicy := sleego.NewShutdownPolicyImpl(channelToAlert, alerts)

	go func() {
		if err := shutdownPolicy.Apply(ctx, shutdownTime); err != nil {
			log.Printf("Error applying the shutdown policy: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-channelToAlert:
				a.SendNotification(fyne.NewNotification("Alert", msg))
			case <-ctx.Done():
				return
			}
		}
	}()
}
