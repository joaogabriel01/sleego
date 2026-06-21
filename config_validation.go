package sleego

import (
	"fmt"
	"strings"
	"time"
)

const configTimeLayout = "15:04"

// ValidateConfig rejects invalid configuration before any policy starts.
func ValidateConfig(cfg FileConfig) error {
	if cfg.Shutdown != "" {
		if err := validateConfigTime("shutdown", cfg.Shutdown); err != nil {
			return err
		}
	}

	for i, app := range cfg.Apps {
		if strings.TrimSpace(app.Name) == "" {
			return fmt.Errorf("apps[%d].name is required", i)
		}
		if err := validateConfigTime(fmt.Sprintf("apps[%d].allowed_from", i), app.AllowedFrom); err != nil {
			return err
		}
		if err := validateConfigTime(fmt.Sprintf("apps[%d].allowed_to", i), app.AllowedTo); err != nil {
			return err
		}
	}

	return nil
}

func validateConfigTime(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if strings.TrimSpace(value) != value {
		return fmt.Errorf("%s must use HH:MM format", field)
	}
	if len(value) != 5 || value[2] != ':' || !isDigit(value[0]) || !isDigit(value[1]) || !isDigit(value[3]) || !isDigit(value[4]) {
		return fmt.Errorf("%s must use HH:MM format", field)
	}
	if _, err := time.Parse(configTimeLayout, value); err != nil {
		return fmt.Errorf("%s must use HH:MM format: %w", field, err)
	}
	return nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
