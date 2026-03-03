// Package config handles command-line flags parsing.
package config

import (
	"flag"
	"fmt"
	"os"
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

// Parse initializes a Config struct according to startup flags.
func Parse(args []string) (*Config, error) {
	cfg := &Config{}

	fs := flag.NewFlagSet("inference-stub", flag.ContinueOnError)

	fs.IntVar(&cfg.Port, "port", 8080, "The port to listen on")
	fs.DurationVar(&cfg.Timeout, "timeout", 1*time.Minute, "Timeout for the request")
	fs.DurationVar(&cfg.TTFT, "ttft", 100*time.Millisecond, "Time to first token")
	fs.DurationVar(&cfg.TPOT, "tpot", 20*time.Millisecond, "Time per output token")
	fs.BoolVar(&cfg.Debug, "debug", false, "Enable debug logging")
	fs.IntVar(&cfg.Length, "length", 50, "Number of generated lorem ipsum words")

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
