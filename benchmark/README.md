# JSONLogic Benchmark

This benchmark suite allows you to compare the performance of different versions of the JSONLogic Go implementation.

## Methodology

The benchmark compares two versions of the JSONLogic Go implementation:

1. **Published Version**: The specified release version from GitHub
2. **Current Version**: Your local code changes

The benchmark:

- Creates temporary Go modules for both versions
- Runs the Go benchmark tests with the same set of test cases for both versions
- Compares the results across multiple metrics:
  - Execution time (ns/op)
  - Memory usage (B/op)
  - Memory allocations (allocs/op)
- Calculates improvement/degradation percentages for each metric
- Provides an overall summary showing average improvements

## Test Cases

The benchmark uses a variety of JSONLogic operations to thoroughly test performance:

- Simple equality checks
- Complex conditions with logical operators
- Variable access with nested paths
- Array operations
- Conditional logic

## Requirements

- Go (1.21+)
- Python 3.x
- Python packages: colorama (installed automatically with uv)
- uv (Python package manager)

## Usage

Run the benchmark by executing the `bench` script with a published version to compare against:

```bash
./bench v3.7.5
```

Where `v3.7.5` should be replaced with the specific version tag you want to compare against.

## Understanding the Results

The benchmark results display:

- A side-by-side comparison of each test case for both versions
- Percentage differences for time, memory, and allocations
- Positive percentages (green) indicate improvements in your current code
- Negative percentages (yellow) indicate performance regressions
- Average improvements across all test cases
