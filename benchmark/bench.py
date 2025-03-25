#!/usr/bin/env python3
import os
import sys
import re
import subprocess
import shutil
import tempfile
from typing import Dict, List, Tuple, NamedTuple
import colorama
from colorama import Fore, Style


class BenchmarkResult(NamedTuple):
    test_case: str
    time_ns: float
    bytes_per_op: int
    allocs_per_op: int


def setup_colorama():
    colorama.init()


def create_temp_dir() -> str:
    temp_dir = tempfile.mkdtemp(prefix="jsonlogic_benchmark_", dir=os.getcwd())
    print(f"{Fore.BLUE}Working in temporary directory: {temp_dir}{Style.RESET_ALL}")
    return temp_dir


def cleanup(temp_dir: str):
    """Remove temporary directory and files"""
    print(f"{Fore.BLUE}Cleaning up temporary files...{Style.RESET_ALL}")
    shutil.rmtree(temp_dir, ignore_errors=True)


def create_published_environment(temp_dir: str, version: str) -> str:
    published_dir = os.path.join(temp_dir, "published")
    os.makedirs(published_dir, exist_ok=True)

    shutil.copy("benchmark_test.go", published_dir)

    with open(os.path.join(published_dir, "go.mod"), "w") as f:
        f.write(f"""module benchmark

go 1.21

require github.com/diegoholiveira/jsonlogic/v3 {version}
""")

    print(f"{Fore.BLUE}Created environment for published version {version}{Style.RESET_ALL}")
    return published_dir


def create_current_environment(temp_dir: str, version: str) -> str:
    current_dir = os.path.join(temp_dir, "current")
    os.makedirs(current_dir, exist_ok=True)

    shutil.copy("benchmark_test.go", current_dir)

    with open(os.path.join(current_dir, "go.mod"), "w") as f:
        f.write(f"""module benchmark

go 1.21

require github.com/diegoholiveira/jsonlogic/v3 {version}

replace github.com/diegoholiveira/jsonlogic/v3 => ../../..
""")

    print(f"{Fore.BLUE}Created environment for current version with replace directive{Style.RESET_ALL}")
    return current_dir


def run_benchmark(directory: str, is_published: bool, version: str) -> Dict[str, BenchmarkResult]:
    version_name = "published" if is_published else "current"
    print(f"{Fore.YELLOW}Running go mod tidy for {version_name} version...{Style.RESET_ALL}")

    try:
        subprocess.run(["go", "mod", "tidy"], cwd=directory, check=True,
                      stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    except subprocess.CalledProcessError as e:
        print(f"{Fore.RED}Error running go mod tidy: {e}{Style.RESET_ALL}")
        print(f"{Fore.RED}Stdout: {e.stdout.decode()}{Style.RESET_ALL}")
        print(f"{Fore.RED}Stderr: {e.stderr.decode()}{Style.RESET_ALL}")
        return {}

    print(f"{Fore.YELLOW}Running benchmark for {version_name} version...{Style.RESET_ALL}")
    try:
        result = subprocess.run(["go", "test", "-bench=.", "-benchmem", "-count=5"],
                               cwd=directory, check=True, capture_output=True, text=True)
    except subprocess.CalledProcessError as e:
        print(f"{Fore.RED}Error running benchmark: {e}{Style.RESET_ALL}")
        print(f"{Fore.RED}Stdout: {e.stdout}{Style.RESET_ALL}")
        print(f"{Fore.RED}Stderr: {e.stderr}{Style.RESET_ALL}")
        return {}

    with open(os.path.join(directory, "bench_output.txt"), "w") as f:
        f.write(result.stdout)

    results = {}

    print(f"{Fore.BLUE}First few lines of benchmark output:{Style.RESET_ALL}")
    for i, line in enumerate(result.stdout.split('\n')[:10]):
        print(f"  {line}")

    patterns = [
        r'BenchmarkJSONLogic/([^\s]+)\s+\d+\s+([\d\.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op',
        r'BenchmarkJSONLogic/([^\s]+)-\d+\s+\d+\s+([\d\.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op',
        r'Benchmark[^/]+/([^\s]+)\s+\d+\s+([\d\.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op',
        r'Benchmark[^/]+/([^\s-]+)(?:-\d+)?\s+\d+\s+([\d\.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op'
    ]

    matched = False
    for pattern in patterns:
        matches = re.findall(pattern, result.stdout)
        if matches:
            matched = True
            for match in matches:
                test_case = match[0]
                time_ns = float(match[1])
                bytes_per_op = int(match[2])
                allocs_per_op = int(match[3])

                benchmark_result = BenchmarkResult(
                    test_case=test_case,
                    time_ns=time_ns,
                    bytes_per_op=bytes_per_op,
                    allocs_per_op=allocs_per_op
                )

                results[test_case] = benchmark_result
                print(f"{Fore.BLUE}Recorded {version_name} result for {test_case}: "
                      f"{time_ns}ns/op, {bytes_per_op}B/op, {allocs_per_op}allocs/op{Style.RESET_ALL}")
            break

    if not matched:
        print(f"{Fore.RED}Failed to parse benchmark results. Here's the full output:{Style.RESET_ALL}")
        print(result.stdout)

        print(f"{Fore.YELLOW}Attempting line-by-line matching...{Style.RESET_ALL}")
        for line in result.stdout.split('\n'):
            for pattern in patterns:
                match = re.search(pattern, line)
                if match:
                    test_case = match.group(1)
                    time_ns = float(match.group(2))
                    bytes_per_op = int(match.group(3))
                    allocs_per_op = int(match.group(4))

                    benchmark_result = BenchmarkResult(
                        test_case=test_case,
                        time_ns=time_ns,
                        bytes_per_op=bytes_per_op,
                        allocs_per_op=allocs_per_op
                    )

                    results[test_case] = benchmark_result
                    print(f"{Fore.BLUE}Recorded {version_name} result for {test_case} (line-by-line): "
                          f"{time_ns}ns/op, {bytes_per_op}B/op, {allocs_per_op}allocs/op{Style.RESET_ALL}")

    return results


def run_benchmarks(temp_dir: str, version: str) -> Tuple[Dict[str, BenchmarkResult], Dict[str, BenchmarkResult]]:
    published_dir = create_published_environment(temp_dir, version)
    published_results = run_benchmark(published_dir, is_published=True, version=version)

    current_dir = create_current_environment(temp_dir, version)
    current_results = run_benchmark(current_dir, is_published=False, version=version)

    if not published_results:
        print(f"{Fore.RED}No results for published version. Benchmark failed.{Style.RESET_ALL}")

    if not current_results:
        print(f"{Fore.RED}No results for current version. Benchmark failed.{Style.RESET_ALL}")

    return published_results, current_results


def display_results(published_results: Dict[str, BenchmarkResult], current_results: Dict[str, BenchmarkResult], version: str):
    if not published_results or not current_results:
        print(f"{Fore.RED}Cannot display results - missing data.{Style.RESET_ALL}")
        return

    print(f"{Fore.GREEN}====================================================")
    print(f"JSONLogic Benchmark Results")
    print(f"===================================================={Style.RESET_ALL}")
    print()

    print(f"{Fore.BLUE}Comparing current vs {version}:{Style.RESET_ALL}")
    print()

    print(f"{'Test Case':<20} {'Version':<15} {'Time (ns/op)':<15} {'Memory (B/op)':<15} "
          f"{'Allocs/op':<15} {'Time Diff':<15} {'Memory Diff':<15} {'Allocs Diff':<15}")
    print("-" * 120)

    test_cases = set(published_results.keys()) | set(current_results.keys())

    total_time_improvement = 0
    total_memory_improvement = 0
    total_allocs_improvement = 0
    test_count = 0

    for test in sorted(test_cases):
        if test not in published_results or test not in current_results:
            print(f"{Fore.YELLOW}Warning: Missing data for test case: {test}{Style.RESET_ALL}")
            continue

        pub_result = published_results[test]
        cur_result = current_results[test]

        time_improvement = (pub_result.time_ns - cur_result.time_ns) / pub_result.time_ns * 100
        memory_improvement = (pub_result.bytes_per_op - cur_result.bytes_per_op) / pub_result.bytes_per_op * 100 if pub_result.bytes_per_op else 0
        allocs_improvement = (pub_result.allocs_per_op - cur_result.allocs_per_op) / pub_result.allocs_per_op * 100 if pub_result.allocs_per_op else 0

        total_time_improvement += time_improvement
        total_memory_improvement += memory_improvement
        total_allocs_improvement += allocs_improvement
        test_count += 1

        time_diff = f"{Fore.GREEN}+{time_improvement:.2f}%{Style.RESET_ALL}" if time_improvement > 0 else f"{Fore.YELLOW}{time_improvement:.2f}%{Style.RESET_ALL}"
        memory_diff = f"{Fore.GREEN}+{memory_improvement:.2f}%{Style.RESET_ALL}" if memory_improvement > 0 else f"{Fore.YELLOW}{memory_improvement:.2f}%{Style.RESET_ALL}"
        allocs_diff = f"{Fore.GREEN}+{allocs_improvement:.2f}%{Style.RESET_ALL}" if allocs_improvement > 0 else f"{Fore.YELLOW}{allocs_improvement:.2f}%{Style.RESET_ALL}"

        print(f"{test:<20} {version:<15} {pub_result.time_ns:<15.2f} "
              f"{pub_result.bytes_per_op:<15} {pub_result.allocs_per_op:<15} {'':15} {'':15} {'':15}")

        print(f"{test:<20} {'current':<15} {cur_result.time_ns:<15.2f} "
              f"{cur_result.bytes_per_op:<15} {cur_result.allocs_per_op:<15} "
              f"{time_diff:<15} {memory_diff:<15} {allocs_diff:<15}")

        print("-" * 120)

    if test_count > 0:
        avg_time_improvement = total_time_improvement / test_count
        avg_memory_improvement = total_memory_improvement / test_count
        avg_allocs_improvement = total_allocs_improvement / test_count

        print()
        print(f"{Fore.GREEN}Average Improvements:{Style.RESET_ALL}")

        time_avg = f"{Fore.GREEN}+{avg_time_improvement:.2f}%{Style.RESET_ALL}" if avg_time_improvement > 0 else f"{Fore.YELLOW}{avg_time_improvement:.2f}%{Style.RESET_ALL}"
        memory_avg = f"{Fore.GREEN}+{avg_memory_improvement:.2f}%{Style.RESET_ALL}" if avg_memory_improvement > 0 else f"{Fore.YELLOW}{avg_memory_improvement:.2f}%{Style.RESET_ALL}"
        allocs_avg = f"{Fore.GREEN}+{avg_allocs_improvement:.2f}%{Style.RESET_ALL}" if avg_allocs_improvement > 0 else f"{Fore.YELLOW}{avg_allocs_improvement:.2f}%{Style.RESET_ALL}"

        print(f"- Time: {time_avg}")
        print(f"- Memory: {memory_avg}")
        print(f"- Allocations: {allocs_avg}")
    else:
        print(f"{Fore.YELLOW}No test cases had both published and current results for comparison.{Style.RESET_ALL}")


def main():
    if len(sys.argv) < 2:
        print(f"{Fore.RED}Error: Missing version argument{Style.RESET_ALL}")
        print(f"Usage: {sys.argv[0]} <version>")
        print(f"Example: {sys.argv[0]} v3.7.5")
        sys.exit(1)

    version = sys.argv[1]

    if not os.path.isfile("benchmark_test.go"):
        print(f"{Fore.RED}Error: benchmark_test.go not found in current directory.{Style.RESET_ALL}")
        sys.exit(1)

    setup_colorama()
    print(f"{Fore.GREEN}Starting JSONLogic benchmark comparing local code with {version}...{Style.RESET_ALL}")

    temp_dir = create_temp_dir()

    try:
        published_results, current_results = run_benchmarks(temp_dir, version)

        display_results(published_results, current_results, version)

        print(f"{Fore.GREEN}Benchmark completed!{Style.RESET_ALL}")
    finally:
        cleanup(temp_dir)


if __name__ == "__main__":
    main()
