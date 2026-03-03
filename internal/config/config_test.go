package config

import (
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
