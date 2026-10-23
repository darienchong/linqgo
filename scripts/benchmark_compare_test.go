// Package scripts contains developer tooling that is executed via `go test`.
//
// benchmark_compare_test.go runs all Functioner benchmarks, snapshots the raw
// results into /data as YY_MM_DD_HHmm.txt, and prints a tabular A/B comparison
// of the previous run against the current run using benchstat. The comparison
// reports percentage differences in ns/op (shown as sec/op), B/op and
// allocs/op.
//
// Usage:
//
//	RUN_BENCHMARKS=1 go test ./scripts/ -run TestBenchmarkCompare -v
//
// Optional environment overrides (forwarded to `go test`):
//
//	BENCHTIME=2s   -> passed as -benchtime
//	BENCHCOUNT=6   -> passed as -count (>= 6 lets benchstat show confidence)
//
// benchstat must be on
// PATH or installed under GOBIN/GOPATH/bin; install it with:
//
//	go install golang.org/x/perf/cmd/benchstat@latest
package scripts

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

const (
	// User-modifiable configuration variables
	benchTime  = "1s"
	benchCount = "6"

	// skipPattern is passed to `go test -skip` so the between-run comparison
	// omits the "baseline" arms of each benchmark. Those arms only exist to
	// compare a hand-written baseline against the LINQ library within a single
	// run; across runs they add noise to the benchstat table. Set to "" to
	// include them again. The leading slash targets the sub-benchmark segment
	// (e.g. "baseline @ 1000000") without matching top-level benchmark names.
	skipPattern = "/baseline"
)

// DO NOT EDIT
const (
	// benchPackage is the package that is compiled and benchmarked.
	benchPackage = "./benchmarks/"

	// dataDirRel is where raw benchmark snapshots are stored.
	dataDirRel = "data"

	// timestampLayout renders filenames as YY_MM_DD_HHmm.txt.
	timestampLayout = "06_01_02_1504"
)

//goland:noinspection ALL
func TestBenchmarkCompare(t *testing.T) {
	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("locate repo root: %v", err)
	}

	benchstatBin, err := findBenchstat()
	if err != nil {
		t.Fatalf("%v", err)
	}

	dataDir := filepath.Join(root, dataDirRel)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatalf("create data dir %q: %v", dataDir, err)
	}

	// Capture the previous snapshot BEFORE writing the new one so we never
	// compare a run against itself (even on a same-minute filename collision).
	prevFile, err := latestDataFile(dataDir)
	if err != nil {
		t.Fatalf("find previous data file: %v", err)
	}

	// Run every benchmark in the package, streaming output to the terminal for
	// live progress.
	curText, err := runBenchmarks(root)
	if err != nil {
		t.Fatalf("run benchmarks: %v\n%s", err, curText)
	}

	// Persist the raw results as YY_MM_DD_HHmm.txt.
	curFile := filepath.Join(dataDir, time.Now().Format(timestampLayout)+".txt")
	if err := os.WriteFile(curFile, []byte(curText), 0o644); err != nil {
		t.Fatalf("write results %q: %v", curFile, err)
	}
	t.Logf("wrote raw results to %s", curFile)

	// Render the comparison table via benchstat (previous = base, current = new).
	fmt.Fprintf(os.Stdout, "\n=== Benchmark comparison (benchstat) ===\n")
	fmt.Fprintf(os.Stdout, "current  run: %s\n", filepath.Base(curFile))
	if prevFile == "" || prevFile == curFile {
		fmt.Fprintf(os.Stdout, "previous run: (none) - nothing to compare against yet\n\n")
		return
	}
	fmt.Fprintf(os.Stdout, "previous run: %s\n\n", filepath.Base(prevFile))

	table, err := runBenchstat(benchstatBin, prevFile, curFile)
	fmt.Fprint(os.Stdout, table)
	if err != nil {
		t.Fatalf("run benchstat: %v", err)
	}
	fmt.Fprintln(os.Stdout, "\n(vs base %: positive = regression, negative = improvement; ~ = within noise)")
}

// findRepoRoot walks upwards from the working directory until it finds go.mod.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %q", dir)
		}
		dir = parent
	}
}

// findBenchstat resolves the benchstat binary from PATH, then GOBIN, then
// GOPATH/bin, then the default $HOME/go/bin. It returns an actionable error
// with install instructions when benchstat cannot be found.
func findBenchstat() (string, error) {
	if p, err := exec.LookPath("benchstat"); err == nil {
		return p, nil
	}

	var candidates []string
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		candidates = append(candidates, filepath.Join(gobin, "benchstat"))
	}
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		gopath := strings.TrimSpace(string(out))
		if gopath != "" {
			candidates = append(candidates, filepath.Join(gopath, "bin", "benchstat"))
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "go", "bin", "benchstat"))
	}

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("benchstat not found on PATH or in %v; install with: go install golang.org/x/perf/cmd/benchstat@latest", candidates)
}

// runBenchmarks executes `go test` for the configured package, running every
// benchmark it contains, and returns the combined stdout+stderr output.
func runBenchmarks(root string) (string, error) {
	args := []string{
		"test",
		"-run", "^$", // skip unit tests; benchmarks only
		"-bench", ".", // run all benchmarks in the package
		"-benchmem",
	}
	if bt := os.Getenv("BENCHTIME"); bt != "" {
		args = append(args, "-benchtime", bt)
	} else {
		args = append(args, "-benchtime", benchTime)
	}
	if count := os.Getenv("BENCHCOUNT"); count != "" {
		args = append(args, "-count", count)
	} else {
		args = append(args, "-count", benchCount)
	}
	if //goland:noinspection GoBoolExpressions
	skipPattern != "" {
		args = append(args, "-skip", skipPattern)
	}
	args = append(args, benchPackage)

	var buf bytes.Buffer
	sink := io.MultiWriter(&buf, os.Stdout)

	cmd := exec.Command("go", args...)
	cmd.Dir = root
	cmd.Stdout = sink
	cmd.Stderr = sink

	err := cmd.Run()
	return buf.String(), err
}

// runBenchstat invokes benchstat comparing base (previous) against new
// (current), returning the rendered comparison tables.
func runBenchstat(bin, baseFile, newFile string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command(bin, baseFile, newFile)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// latestDataFile returns the most recent snapshot in dataDir. Filenames follow
// the zero-padded YY_MM_DD_HHmm layout, so lexical ordering is chronological.
func latestDataFile(dataDir string) (string, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".txt") {
			files = append(files, e.Name())
		}
	}
	if len(files) == 0 {
		return "", nil
	}
	sort.Strings(files)
	return filepath.Join(dataDir, files[len(files)-1]), nil
}
