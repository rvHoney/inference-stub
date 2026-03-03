# inference-stub

Minimalist, deterministic OpenAI-compatible stub for LLM infrastructure testing

## Overview

**inference-stub** is a specialized tool designed to simulate LLM inference streams. By providing a predictable, programmable backend, it allows for isolated performance analysis of AI Gateways and Proxy layers.

Unlike a real LLM, this stub removes inference variability, making it possible to measure the precise overhead of the networking stack (TTFT/TPOT) in Cloud-Native environments.

## Current Focus
- **Gateway Benchmarking:** Isolate proxy latency by using deterministic TTFT/TPOT settings.
- **Protocol Validation:** Ensure Gateway-level filters (Rate Limiting, Usage Tracking) behave correctly against standard OpenAI-compatible SSE streams.
- **CI/CD Integration:** Provide a lightweight, zero-cost alternative to real LLMs for automated integration tests.

## Getting Started

```bash
# Build the binary
make build

# Run the stub with 100ms TTFT and 20ms TPOT
./bin/inference-stub --ttft 100ms --tpot 20ms --port 8080
```

## Roadmap
- **Error Injection:** Support for simulating 429 Too Many Requests and 503 Service Unavailable.
- **Usage Reporting:** Implementation of the usage field in the final stream chunk for quota-testing.
- **Kubernetes Integration:** Helm charts and specialized manifests for kind deployments.

---
*Developed for the GSoC 2026 - kgateway Performance Benchmarking project.*