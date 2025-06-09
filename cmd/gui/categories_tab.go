package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildCategoryContent() fyne.CanvasObject {
	categoriesContainer = container.NewVBox()
	for category, apps := range fileConfig.Categories {
		loggerInstance.Debug(fmt.Sprintf("Category: %s, Apps: %v", category, apps))
		addCategory(category, apps)
	}

	addButton := widget.NewButton("Add", func() { addCategory("", []string{}) })

	return container.NewVBox(
		categoriesContainer,
		addButton,

		widget.NewSeparator(),
	)
}

func addCategory(category string, apps []string) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(category)

	appsEntry := widget.NewEntry()
	appsEntry.SetText(strings.Join(apps, ";"))

	removeBtn := widget.NewButton("Remove", func() {
		removeCategory(nameEntry.Text)
	})
	form := container.NewVBox(
		widget.NewLabel("Category:"),
		nameEntry,
		widget.NewLabel("Apps:"),
		appsEntry,
		removeBtn,
		widget.NewSeparator(),
	)
	categoriesContainer.Add(form)
	categoryEntries[category] = CategoryEntry{
		nameEntry: nameEntry,
		appsEntry: appsEntry,
		container: form,
	}

	categoryEntries[category].nameEntry.OnChanged = func(newText string) {
		if _, exists := categoryEntries[newText]; exists {
			delete(categoryEntries, category)
			loggerInstance.Debug(fmt.Sprintf("Updated category name from %s to %s", category, newText))
		}
		categoryEntries[newText] = CategoryEntry{
			nameEntry: nameEntry,
			appsEntry: appsEntry,
			container: form,
		}
	}

}

func removeCategory(categoryName string) {
	var categoryKey string
	for key, entry := range categoryEntries {
		if entry.nameEntry.Text == categoryName {
			categoryKey = key
			break
		}
	}

	if categoryKey != "" {
		entry := categoryEntries[categoryKey]
		categoriesContainer.Remove(entry.container)
		delete(categoryEntries, categoryKey)
		loggerInstance.Debug(fmt.Sprintf("Removed category: %s", categoryName))
	} else {
		loggerInstance.Error(fmt.Sprintf("Category %s does not exist", categoryName))
	}
}
