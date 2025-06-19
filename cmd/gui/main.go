package main

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/joaogabriel01/sleego"
	"github.com/joaogabriel01/sleego/internal/logger"
)

type AppEntry struct {
	nameEntry        *widget.Entry
	allowedFromEntry *widget.Entry
	allowedToEntry   *widget.Entry
	container        *fyne.Container
}

type CategoryEntry struct {
	nameEntry *widget.Entry
	appsEntry *widget.Entry
	container *fyne.Container
}

var (
	a                   fyne.App
	w                   fyne.Window
	ctx                 context.Context
	cancel              context.CancelFunc
	loader              sleego.Loader
	configPath          string
	logLevel            string
	loggerInstance      logger.Logger
	fileConfig          sleego.FileConfig
	appConfigs          []sleego.AppConfig
	entries             map[string]AppEntry
	categoryEntries     map[string]CategoryEntry
	appsContainer       *fyne.Container
	categoriesContainer *fyne.Container
	shutdownTimeEntry   *widget.Entry

	screenSize       fyne.Size
	applicationsSize fyne.Size
	categoriesSize   fyne.Size

	categoriesOperator sleego.CategoryOperator

	//go:embed assets/sleego.ico
	iconData []byte
)

func main() {
	entries = make(map[string]AppEntry)
	categoryEntries = make(map[string]CategoryEntry)
	categoriesOperator = sleego.GetCategoryOperator()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	parseFlags()

	var err error
	loggerInstance, err = logger.Get(logLevel)
	if err != nil {
		log.Fatalf("Error getting logger instance: %v", err)
	}

	initializeApp()
	setupTrayIcon()
	setupWindow()
	loadConfigurations()
	root := container.NewStack()
	root.Objects = []fyne.CanvasObject{createUI()}

	w.SetContent(root)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			newSize := w.Canvas().Size()
			if screenSize != newSize {
				loggerInstance.Debug(fmt.Sprintf("Window size changed: %v", newSize))
				screenSize = newSize
				root.Objects = []fyne.CanvasObject{createUI()}
				root.Refresh()
			}
		}
	}()

	w.ShowAndRun()
}

func parseFlags() {
	configUser, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}
	configPathP := flag.String("config", configUser+"/sleego/config.json", "Path to config file")

	logLevelP := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")

	flag.Parse()

	configPath = *configPathP
	logLevel = *logLevelP
	if logLevel != "debug" && logLevel != "info" && logLevel != "warn" && logLevel != "error" {
		logLevel = "info"
	}
}

func initializeApp() {
	a = app.NewWithID("sleego.gui")
	icon := fyne.NewStaticResource("icon.png", iconData)
	a.SetIcon(icon)
	w = a.NewWindow("Sleego")
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
	screenSize = w.Canvas().Size()
	if screenSize.Width == 0 || screenSize.Height == 0 {
		screenSize = fyne.NewSize(1200, 700)
	}
	w.Resize(screenSize)
}

func loadConfigurations() {
	var err error
	fileConfig, err = loader.Load(configPath)
	if err != nil {
		loggerInstance.Error("Error loading configurations: " + err.Error())
		dialog.ShowError(fmt.Errorf("error loading configurations: %v", err), w)
	}
	appConfigs = fileConfig.Apps
}

func createUI() *fyne.Container {

	applicationsSize = fyne.NewSize(screenSize.Width*0.7, screenSize.Height*0.6)
	categoriesSize = fyne.NewSize(screenSize.Width*0.7, screenSize.Height*0.6)

	mainContent := buildAppContent()
	categoryContent := buildCategoryContent()
	tabs := setupTabs(mainContent, categoryContent)

	shutdownTimeEntry = widget.NewEntry()
	shutdownTimeEntry.SetPlaceHolder("Enter shutdown time (HH:MM)")
	if fileConfig.Shutdown != "" {
		shutdownTimeEntry.SetText(fileConfig.Shutdown)
	}

	saveButton := widget.NewButton("Save", saveConfigurations)
	runButton := widget.NewButton("Run", runPolicies)

	runningProcess := container.NewVBox(
		widget.NewLabel("Running Process"),
		// will be updated with the running process
	)

	content := container.NewVBox(
		container.NewHBox(
			tabs,
			runningProcess,
		),
		shutdownTimeEntry,
		container.NewHBox(
			saveButton,
			runButton,
		),
	)
	return content

}

func newBox(content fyne.CanvasObject) *fyne.Container {
	square := canvas.NewRectangle(color.Black)
	context := container.NewStack(
		square,
		content,
	)
	return context
}

func setupTabs(mainContent fyne.CanvasObject, categoryContent fyne.CanvasObject) *container.AppTabs {
	mainBox := newBox(mainContent)
	catBox := newBox(categoryContent)

	mainScroll := container.NewScroll(mainBox)
	mainScroll.SetMinSize(applicationsSize)

	categoryScroll := container.NewScroll(catBox)
	categoryScroll.SetMinSize(categoriesSize)

	loggerInstance.Debug(fmt.Sprintf("Main scroll size: %v %v", applicationsSize.Width, applicationsSize.Height))
	loggerInstance.Debug(fmt.Sprintf("Category scroll size: %v %v", categoriesSize.Width, categoriesSize.Height))

	tabs := container.NewAppTabs(
		container.NewTabItem("Applications", mainScroll),
		container.NewTabItem("Categories", categoryScroll),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}

func saveConfigurations() {
	updatedAppConfigs := make([]sleego.AppConfig, len(entries))
	for _, entry := range entries {
		updatedAppConfigs = append(updatedAppConfigs, sleego.AppConfig{
			Name:        entry.nameEntry.Text,
			AllowedFrom: entry.allowedFromEntry.Text,
			AllowedTo:   entry.allowedToEntry.Text,
		})
	}
	categoriesConfig := make(map[string][]string)
	for _, entry := range categoryEntries {
		categoryName := entry.nameEntry.Text
		appsText := entry.appsEntry.Text
		if categoryName == "" || appsText == "" {
			continue
		}
		categoriesConfig[categoryName] = strings.Split(appsText, ";")
		loggerInstance.Debug(fmt.Sprintf("Category: %s, Apps: %v", categoryName, categoriesConfig[categoryName]))
	}
	shutdownConfig := shutdownTimeEntry.Text

	updatedConfigs := sleego.FileConfig{
		Apps:       updatedAppConfigs,
		Categories: categoriesConfig,
		Shutdown:   shutdownConfig,
	}

	if err := loader.Save(configPath, updatedConfigs); err != nil {
		loggerInstance.Error("Error saving the configuration file: " + err.Error())
		dialog.ShowError(fmt.Errorf("error saving the configuration file: %v", err), w)
		return
	}
	fileConfig = updatedConfigs
	categoriesOperator.SetProcessByCategories(fileConfig.Categories)
	dialog.ShowInformation("Success", "Configurations saved successfully!", w)
	loggerInstance.Info("Configurations saved successfully!")
}

func runPolicies() {
	cancel()
	ctx, cancel = context.WithCancel(context.Background())

	monitor := &sleego.ProcessorMonitorImpl{}
	processChan := make(chan string)
	policy := sleego.NewProcessPolicyImpl(monitor, categoriesOperator, nil, processChan)
	loggerInstance.Debug(fmt.Sprintf("Starting process policy with config: %+v of path: %s", fileConfig, configPath))
	dialog.ShowInformation("Running", "Applying the policy...", w)

	go func() {
		if err := policy.Apply(ctx, fileConfig.Apps); err != nil {
			loggerInstance.Error("Error applying the process policy: " + err.Error())
			dialog.ShowError(fmt.Errorf("error applying the process policy: %v", err), w)
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

	loggerInstance.Info("Starting shutdown policy with time:" + shutdownTimeStr)
	channelToAlert := make(chan string)
	alerts := []int{10, 3, 1}
	shutdownPolicy := sleego.NewShutdownPolicyImpl(channelToAlert, alerts)

	go func() {
		if err := shutdownPolicy.Apply(ctx, shutdownTime); err != nil {
			loggerInstance.Error("Error applying the shutdown policy: " + err.Error())
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
