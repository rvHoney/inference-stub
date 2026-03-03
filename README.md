# inference-stub

> Minimalist, deterministic OpenAI-compatible stub for LLM infrastructure testing.

[![Go Version](https://img.shields.io/badge/go-1.22.0-blue.svg)](https://golang.org)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-yellow.svg)](https://opensource.org/licenses/Apache-2.0)
[![CI](https://github.com/rvHoney/inference-stub/actions/workflows/release.yml/badge.svg)](https://github.com/rvHoney/inference-stub/actions/workflows/release.yml)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://github.com/rvHoney/inference-stub/pkgs/container/inference-stub)

## Overview

**inference-stub** is a specialized tool designed to simulate LLM inference streams. By providing a predictable, programmable backend, it allows for isolated performance analysis of AI Gateways and Proxy layers.

Unlike a real LLM, this stub removes inference variability, making it possible to measure the precise overhead of the networking stack (TTFT/TPOT) in Cloud-Native environments. It supports both stream and non-stream requests, returning dynamically generated Lorem Ipsum text based on configurable parameters.

## Current Focus
- **Gateway Benchmarking:** Isolate proxy latency by using deterministic TTFT/TPOT settings.
- **Protocol Validation:** Ensure Gateway-level filters (Rate Limiting, Usage Tracking) behave correctly against standard OpenAI-compatible JSON responses and SSE streams.
- **CI/CD Integration:** Provide a lightweight, zero-cost alternative to real LLMs for automated integration tests.

## Getting Started

```bash
# Build the binary
make build

# Run the stub with 100ms TTFT, 20ms TPOT, and fixed payload length of 15 words
./bin/inference-stub --ttft 100ms --tpot 20ms --length 15 --port 8080
```

### Configuration Flags
- `--port` (default `8080`): The port to listen on.
- `--ttft` (default `100ms`): Time to first token. Simulates the initial processing delay.
- `--tpot` (default `20ms`): Time per output token. Simulates the delay between generation steps.
- `--length` (default `50`): The exact number of Lorem Ipsum words to generate in the mock response.
- `--timeout` (default `1m0s`): Timeout for requests.
- `--debug` (default `false`): Enable debug logging.

## Roadmap
- **Error Injection:** Support for simulating 429 Too Many Requests and 503 Service Unavailable.
- **Usage Reporting:** Implementation of the usage field in the final stream chunk for quota-testing.
- **Kubernetes Integration:** Helm charts and specialized manifests for kind deployments.

---
*Developed for the GSoC 2026 - kgateway Performance Benchmarking project.*