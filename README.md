<div align="center">

# go-argus-benchmark

**Performance comparison between [go-argus] and [go-playground/validator/v10]**

[go-argus]: https://github.com/kamalyes/go-argus
[go-playground/validator/v10]: https://github.com/go-playground/validator

[![Benchmark](https://github.com/kamalyes/go-argus-benchmark/actions/workflows/benchmark.yml/badge.svg)](https://github.com/kamalyes/go-argus-benchmark/actions/workflows/benchmark.yml)

[中文文档](./README_zh.md)

</div>

---

## 📊 Benchmark Charts

### Latency — Sequential

![Latency Sequential](./benchmarks/latency.svg)

### Latency — Parallel

![Latency Parallel](./benchmarks/latency_parallel.svg)

### Memory Allocation — Sequential

![Allocs Sequential](./benchmarks/allocs.svg)

> Charts are auto-generated on every push via GitHub Actions.

---

## 🚀 Quick Start

```bash
# Run all benchmarks
go test -run='^$' -bench=. -benchmem -count=3 -timeout=30m ./... | tee benchmark_output.txt

# Generate charts + BENCHMARKS.md
go run ./bootstrap/report

# Parse real benchmark output
go test -run='^$' -bench=. -benchmem -count=1 -timeout=10m ./... > benchmark_output.txt
go run ./bootstrap/report -parse benchmark_output.txt
```

---

## 📁 Project Structure

```bash
go-argus-benchmark/
├── benchmark_test.go              # Main benchmark test cases
├── benchmark_tags_test.go         # Tag-level benchmarks
├── models.go                      # Test struct models
├── bootstrap/
│   └── report/
│       └── main.go                # Report + SVG chart generator
├── benchmarks/
│   ├── latency.svg                # Sequential latency chart
│   ├── latency_parallel.svg       # Parallel latency chart
│   ├── allocs.svg                 # Memory allocation chart
│   └── *.json                     # Raw benchmark data
├── .github/workflows/
│   └── benchmark.yml              # CI auto-benchmark workflow
├── BENCHMARKS.md                  # Detailed data tables
└── go.mod
```

---

## ⚙️ CI Automation

The [benchmark workflow](./.github/workflows/benchmark.yml) runs automatically on every push to `main`/`master`:

1. Runs `go test -bench` with 5 iterations
2. Parses results and generates SVG charts
3. Commits updated charts + `BENCHMARKS.md` back to the repo

---

## 📋 Full Data

See [BENCHMARKS.md](./BENCHMARKS.md) for detailed comparison tables with all scenarios.

---

## 📝 License

This project is licensed under the MIT License.
