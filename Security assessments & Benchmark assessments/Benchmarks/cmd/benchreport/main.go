package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type FuncInfo struct {
	PackagePath string
	Name        string
}

type BenchResult struct {
	Package string
	Name    string
	NsPerOp float64
	BPerOp  float64
	Allocs  float64
}

func main() {
	funcs := collectFunctions(".")
	fmt.Printf("Total functions: %d\n", len(funcs))
	funcsByPkg := map[string]int{}
	for _, f := range funcs {
		funcsByPkg[f.PackagePath]++
	}
	fmt.Println("Function counts by package:")
	for pkg, count := range funcsByPkg {
		fmt.Printf("  %s: %d\n", pkg, count)
	}

	pkgs := listPackages()
	benchResults := runBenchmarks(pkgs)
	sort.Slice(benchResults, func(i, j int) bool { return benchResults[i].NsPerOp > benchResults[j].NsPerOp })
	fmt.Println("\nSlowest benchmarks:")
	for i, b := range benchResults {
		if i >= 10 {
			break
		}
		fmt.Printf("%s %s %.0f ns/op %.0f B/op %.0f allocs/op\n", b.Package, b.Name, b.NsPerOp, b.BPerOp, b.Allocs)
	}
}

func collectFunctions(root string) []FuncInfo {
	var funcs []FuncInfo
	fset := token.NewFileSet()
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}
		pkg := packagePath(path)
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				funcs = append(funcs, FuncInfo{PackagePath: pkg, Name: fn.Name.Name})
			}
		}
		return nil
	})
	return funcs
}

func packagePath(path string) string {
	dir := filepath.Dir(path)
	mod, err := modulePath()
	if err != nil {
		return dir
	}
	rel, err := filepath.Rel(".", dir)
	if err != nil {
		return dir
	}
	return filepath.Join(mod, filepath.ToSlash(rel))
}

func modulePath() (string, error) {
	out, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func listPackages() []string {
	out, err := exec.Command("go", "list", "./...").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var pkgs []string
	for _, l := range lines {
		if strings.Contains(l, "cmd/synnergy") { // skip known failing package
			continue
		}
		pkgs = append(pkgs, l)
	}
	return pkgs
}

func runBenchmarks(pkgs []string) []BenchResult {
	var results []BenchResult
	benchRe := regexp.MustCompile(`^Benchmark(\S+)\s+\d+\s+([0-9.]+) ns/op(?:\s+([0-9.]+) B/op\s+([0-9.]+) allocs/op)?`)
	for _, pkg := range pkgs {
		cmd := exec.Command("go", "test", pkg, "-run=^$", "-bench", ".", "-benchmem", "-count=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "benchmark failed for %s: %v\n%s\n", pkg, err, string(out))
			continue
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			line := scanner.Text()
			if m := benchRe.FindStringSubmatch(line); m != nil {
				ns, _ := strconv.ParseFloat(m[2], 64)
				var b, a float64
				if len(m) > 3 {
					b, _ = strconv.ParseFloat(m[3], 64)
					a, _ = strconv.ParseFloat(m[4], 64)
				}
				results = append(results, BenchResult{
					Package: pkg,
					Name:    m[1],
					NsPerOp: ns,
					BPerOp:  b,
					Allocs:  a,
				})
			}
		}
	}
	return results
}
