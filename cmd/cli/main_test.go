package main

import (
	"reflect"
	"testing"

	"github.com/joaogabriel01/sleego"
)

type fakeConfigLoader struct {
	config sleego.FileConfig
	err    error
}

func (f fakeConfigLoader) Load(_ string) (sleego.FileConfig, error) {
	return f.config, f.err
}

type recordingCategoryOperator struct {
	categories map[string][]string
}

func (r *recordingCategoryOperator) GetCategoriesOf(_ string) []string {
	return nil
}

func (r *recordingCategoryOperator) SetProcessByCategories(categories map[string][]string) {
	r.categories = categories
}

func TestLoadConfigSetsCategoriesOnOperator(t *testing.T) {
	categories := map[string][]string{
		"games": {"steam.exe"},
	}
	loader := fakeConfigLoader{
		config: sleego.FileConfig{
			Apps: []sleego.AppConfig{
				{Name: "games", AllowedFrom: "09:00", AllowedTo: "17:00"},
			},
			Categories: categories,
		},
	}
	categoryOp := &recordingCategoryOperator{}

	config, err := loadConfig("config.json", loader, categoryOp)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	if !reflect.DeepEqual(config.Categories, categories) {
		t.Fatalf("loadConfig() categories = %v, want %v", config.Categories, categories)
	}
	if !reflect.DeepEqual(categoryOp.categories, categories) {
		t.Fatalf("SetProcessByCategories() categories = %v, want %v", categoryOp.categories, categories)
	}
}
