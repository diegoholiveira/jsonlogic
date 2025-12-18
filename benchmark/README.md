# JSONLogic Benchmark

Benchmark suite to compare performance between different versions of the JSONLogic library.

## Prerequisites

- Go 1.21+
- `benchstat` (install with: `go install golang.org/x/perf/cmd/benchstat@latest`)

## Usage

Compare your current code against a published version:

```bash
./bench v3.7.5
```

The script will:
1. Create an isolated git worktree for the target version
2. **Copy current benchmark code** to target version (ensures fair comparison)
3. Run **comprehensive benchmarks** for both versions (10 iterations each)
4. Display statistical comparison using `benchstat`
5. Clean up automatically

### Benchmark Suites

By default, the script runs the **comprehensive suite** (8 complex benchmarks, ~6-8s total):
```bash
./bench v3.7.5  # Fast, realistic comparison
```

To run the **detailed suite** (all 65 benchmarks, ~60s total):
```bash
./bench v3.7.5 BenchmarkDetailed
```

To run specific benchmark categories:
```bash
./bench v3.7.5 BenchmarkMathOperations
./bench v3.7.5 BenchmarkArrayOperationsScaling
```

## Understanding the Output

`benchstat` shows the performance comparison:

```
name                old time/op    new time/op    delta
JSONLogic/simple-8    900ns ± 2%     850ns ± 3%   -5.56%  (p=0.000 n=10+10)
```

- **old time/op**: Target version performance
- **new time/op**: Current code performance
- **delta**: Percentage change (negative = improvement, positive = regression)
- **±**: Variation/noise in measurements
- **p-value**: Statistical significance (p < 0.05 means the difference is real)

## Benchmark Suites

### Comprehensive Suite (Default)

The comprehensive suite contains 8 complex, realistic benchmarks that exercise multiple operators:

1. **user_validation** - Complex user validation with `and`, `or`, `>=`, `in`, `==`, `+`, `var`
2. **data_pipeline** - Filter + reduce chain for data aggregation
3. **business_rules** - Nested if/else with conditional pricing logic
4. **array_validation** - Combines `all`, `some`, `none` for array validation
5. **string_processing** - String operations with `in`, `substr`, and comparisons
6. **custom_operators** - Tests `contains_all`, `contains_any`, `contains_none`
7. **complex_data_transform** - Chained `filter` + `map` with complex conditions
8. **nested_conditions** - Multi-level conditionals with `missing`, `missing_some`, `cat`

Each benchmark represents real-world usage patterns and exercises 5-7 operators per test.

### Detailed Suite

The detailed suite contains all individual operator benchmarks organized by category:

### Core Operations
- **baseline_noop**: Minimal baseline benchmark (just `true`)
- **simple_equal**: Basic equality check
- **complex_condition**: Nested logical operators with `and`
- **nested_var**: Deep variable path access with defaults
- **complex_logic**: Conditional if/else logic

### Array Operations
- **array_operations**: Array map operations
- **reduce_sum**: Reduce operation with sum accumulation
- **filter_even_numbers**: Filter with modulo operation
- **all_validation**: All operator for validation patterns
- **some_validation**: Some operator for validation patterns
- **merge_arrays**: Merge multiple arrays

### Custom Operators
- **contains_all**: Tests if all elements exist in array
- **contains_any**: Tests if any element exists in array
- **contains_none**: Tests if no elements exist in array

### String Operations
- **string_concatenation**: String concatenation with `cat`
- **substring_extraction**: Substring extraction with `substr`

### Math Operations
- **max_operation**: Find maximum value
- **min_operation**: Find minimum value
- **modulo_operation**: Modulo operator

### Logic Operations
- **or_operation**: Logical OR operator

### Field Validation
- **missing_fields**: Missing field detection

### Complex Scenarios
- **deeply_nested_operations**: Nested filter and map operations

## Benchmark Categories

The benchmark suite is organized into multiple categories for targeted testing:

### Run All Benchmarks
```bash
go test -bench=. ./benchmark/
```

### Run Specific Categories

**Main benchmark suite** (all test cases):
```bash
go test -bench=BenchmarkJSONLogic$ ./benchmark/
```

**Parallel benchmarks** (concurrent performance testing):
```bash
go test -bench=BenchmarkJSONLogicParallel ./benchmark/
```

**Scaling benchmarks** (tests with 10, 100, 1000 element arrays):
```bash
go test -bench=BenchmarkArrayOperationsScaling ./benchmark/
```

**Math operations only**:
```bash
go test -bench=BenchmarkMathOperations ./benchmark/
```

**String operations only**:
```bash
go test -bench=BenchmarkStringOperations ./benchmark/
```

**Logic operations only**:
```bash
go test -bench=BenchmarkLogicOperations ./benchmark/
```

**Custom operators only**:
```bash
go test -bench=BenchmarkCustomOperators ./benchmark/
```

## Benchmark Types

### 1. Standard Benchmarks (`BenchmarkJSONLogic`)
Tests all 22 core operations with realistic data sizes.

### 2. Parallel Benchmarks (`BenchmarkJSONLogicParallel`)
Tests concurrent usage with `RunParallel` for:
- simple_equal
- map
- reduce
- filter

### 3. Scaling Benchmarks (`BenchmarkArrayOperationsScaling`)
Tests performance with different array sizes (10, 100, 1000 elements) for:
- map
- filter
- reduce
- all
- some
- none

### 4. Categorical Benchmarks
Organized by operation type for focused testing:
- Math: +, -, *, /, %, max, min, abs
- String: cat, substr, in
- Logic: and, or, !, if
- Custom: contains_all, contains_any, contains_none

## Advanced Profiling

### Memory profiling:
```bash
cd benchmark
go test -bench=. -benchmem -memprofile=mem.prof
go tool pprof -http=:8080 mem.prof
```

### CPU profiling:
```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof -http=:8080 cpu.prof
```

### Compare scaling performance:
```bash
go test -bench=BenchmarkArrayOperationsScaling/map -benchmem
```

This will show how map performance scales from 10 to 1000 elements.
