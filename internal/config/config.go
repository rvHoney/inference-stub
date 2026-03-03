// Package config handles command-line flags parsing.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config stores inference-stub server config such as port and TTFT and TPOT durations
type Config struct {
	Port    int
	TTFT    time.Duration
	TPOT    time.Duration
	Timeout time.Duration
	Debug   bool
	Length  int
}

// getEnv* reads an environment variable and falls back to a default value

func getEnvStr(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	strValue := getEnvStr(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	strValue := getEnvStr(key, "")
	if value, err := time.ParseDuration(strValue); err == nil {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	strValue := getEnvStr(key, "")
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value
	}
	return fallback
}

// Parse initializes a Config struct according to startup flags.
func Parse(args []string) (*Config, error) {
	cfg := &Config{}

	fs := flag.NewFlagSet("inference-stub", flag.ContinueOnError)

	fs.IntVar(&cfg.Port, "port", getEnvInt("PORT", 8080), "The port to listen on")
	fs.DurationVar(&cfg.Timeout, "timeout", getEnvDuration("TIMEOUT", 1*time.Minute), "Timeout for the request")
	fs.DurationVar(&cfg.TTFT, "ttft", getEnvDuration("TTFT", 100*time.Millisecond), "Time to first token")
	fs.DurationVar(&cfg.TPOT, "tpot", getEnvDuration("TPOT", 20*time.Millisecond), "Time per output token")
	fs.BoolVar(&cfg.Debug, "debug", getEnvBool("DEBUG", false), "Enable debug logging")
	fs.IntVar(&cfg.Length, "length", getEnvInt("LENGTH", 50), "Number of generated lorem ipsum words")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of inference-stub:\n")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
