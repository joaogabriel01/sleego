package sleego

import "testing"

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     FileConfig
		wantErr bool
	}{
		{
			name: "valid config with shutdown",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: "09:00", AllowedTo: "18:00"},
				},
				Shutdown: "23:59",
			},
		},
		{
			name: "shutdown disabled",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: "09:00", AllowedTo: "18:00"},
				},
				Shutdown: "",
			},
		},
		{
			name: "no app policy",
			cfg: FileConfig{
				Shutdown: "23:59",
			},
		},
		{
			name: "invalid shutdown",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: "09:00", AllowedTo: "18:00"},
				},
				Shutdown: "24:00",
			},
			wantErr: true,
		},
		{
			name: "invalid allowed from",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: "25:00", AllowedTo: "18:00"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid allowed to",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: "09:00", AllowedTo: "18:60"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty app name",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: " ", AllowedFrom: "09:00", AllowedTo: "18:00"},
				},
			},
			wantErr: true,
		},
		{
			name: "single digit hour",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: "9:00", AllowedTo: "18:00"},
				},
			},
			wantErr: true,
		},
		{
			name: "time with spaces",
			cfg: FileConfig{
				Apps: []AppConfig{
					{Name: "code", AllowedFrom: " 09:00", AllowedTo: "18:00"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
