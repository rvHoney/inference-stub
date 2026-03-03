package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Config
		wantErr bool
	}{
		{
			name: "default values",
			args: []string{},
			want: &Config{
				Port:    8080,
				TTFT:    100 * time.Millisecond,
				TPOT:    20 * time.Millisecond,
				Timeout: 1 * time.Minute,
				Length:  50,
			},
			wantErr: false,
		},
		{
			name: "custom values",
			args: []string{
				"-port", "9090",
				"-ttft", "50ms",
				"-tpot", "10ms",
				"-length", "15",
			},
			want: &Config{
				Port:    9090,
				TTFT:    50 * time.Millisecond,
				TPOT:    10 * time.Millisecond,
				Timeout: 1 * time.Minute,
				Length:  15,
			},
			wantErr: false,
		},
		{
			name:    "invalid duration format",
			args:    []string{"-ttft", "invalid"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid port format",
			args:    []string{"-port", "notaport"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse_EnvVars(t *testing.T) {
	os.Setenv("PORT", "7070")
	os.Setenv("TIMEOUT", "2m")
	os.Setenv("TTFT", "300ms")
	os.Setenv("TPOT", "15ms")
	os.Setenv("DEBUG", "true")
	os.Setenv("LENGTH", "100")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("TIMEOUT")
		os.Unsetenv("TTFT")
		os.Unsetenv("TPOT")
		os.Unsetenv("DEBUG")
		os.Unsetenv("LENGTH")
	}()

	// Call Parse with empty args to rely solely on env vars
	got, err := Parse([]string{})
	if err != nil {
		t.Fatalf("Parse() failed with env vars set: %v", err)
	}

	want := &Config{
		Port:    7070,
		Timeout: 2 * time.Minute,
		TTFT:    300 * time.Millisecond,
		TPOT:    15 * time.Millisecond,
		Debug:   true,
		Length:  100,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %v, want %v", got, want)
	}
}
