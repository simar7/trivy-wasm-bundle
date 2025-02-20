package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type BenchmarkResult struct {
	Name       string
	Iteration  int
	Time       int64
	Memory     int64
	Allocations int
}

func parseBenchmarkOutput(file string) ([]BenchmarkResult, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var results []BenchmarkResult

	benchRegex := regexp.MustCompile(`Benchmark(\w+)-(\d+)\s+(\d+)\s+(\d+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := benchRegex.FindStringSubmatch(line); matches != nil {
			timeNs, _ := strconv.ParseInt(matches[3], 10, 64)
			memBytes, _ := strconv.ParseInt(matches[4], 10, 64)
			allocs, _ := strconv.Atoi(matches[5])

			results = append(results, BenchmarkResult{
				Name:        matches[1],
				Iteration:   1, // Since each benchmark line represents a single test run
				Time:        timeNs,
				Memory:      memBytes,
				Allocations: allocs,
			})
		}
	}

	return results, scanner.Err()
}

func generateMarkdown(results []BenchmarkResult) string {
	grouped := make(map[string][]BenchmarkResult)
	for _, result := range results {
		grouped[result.Name] = append(grouped[result.Name], result)
	}

	var sb strings.Builder
	sb.WriteString("## Benchmark Results\n\n")

	for name, benchmarks := range grouped {
		sb.WriteString(fmt.Sprintf("### %s\n\n", name))
		sb.WriteString("| Iteration | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |\n")
		sb.WriteString("|-----------|-------------|---------------|--------------------------|\n")
		for i, bench := range benchmarks {
			sb.WriteString(fmt.Sprintf("| %d | %d | %d | %d |\n", i+1, bench.Time, bench.Memory, bench.Allocations))
		}
		sb.WriteString("\n---\n\n")
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <benchmark_output_file>")
		os.Exit(1)
	}

	file := os.Args[1]
	results, err := parseBenchmarkOutput(file)
	if err != nil {
		fmt.Println("Error reading benchmark output:", err)
		os.Exit(1)
	}

	markdown := generateMarkdown(results)
	fmt.Println(markdown)
}

