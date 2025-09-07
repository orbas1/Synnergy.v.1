package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// generateTables parses the list of functions and writes opcode and gas tables.
func generateTables(functionsPath, opcodePath, gasPath string, base uint64) error {
	f, err := os.Open(functionsPath)
	if err != nil {
		return err
	}
	defer f.Close()

	re := regexp.MustCompile(`^([^:]+):.*func\s+(?:\([^)]*\)\s+)?([A-Za-z0-9_]+)\s*\(`)
	uniq := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			path, fn := m[1], m[2]
			if strings.Contains(path, "_test.go") || strings.HasPrefix(path, "cli/") || strings.HasPrefix(path, "cmd/") {
				continue
			}
			uniq[fn] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	funcs := make([]string, 0, len(uniq))
	for k := range uniq {
		funcs = append(funcs, k)
	}
	sort.Strings(funcs)

	opFile, err := os.Create(opcodePath)
	if err != nil {
		return err
	}
	defer opFile.Close()
	fmt.Fprintln(opFile, "| Function | Opcode |")
	fmt.Fprintln(opFile, "|---|---|")

	gasFile, err := os.Create(gasPath)
	if err != nil {
		return err
	}
	defer gasFile.Close()
	fmt.Fprintln(gasFile, "| Function | Gas Cost |")
	fmt.Fprintln(gasFile, "|---|---|")

	for i, fn := range funcs {
		opcode := base + uint64(i) + 1
		fmt.Fprintf(opFile, "| `%s` | `0x%06X` |\n", fn, opcode)
		fmt.Fprintf(gasFile, "| `%s` | `1` |\n", fn)
	}
	return nil
}

func main() {
	var (
		functions = flag.String("functions", "functions_list.txt", "path to functions list")
		opcodeOut = flag.String("opcode-out", "opcodes_list.md", "output file for opcode table")
		gasOut    = flag.String("gas-out", "gas_table_list.md", "output file for gas table")
		base      = flag.Uint64("base", 0x100000, "starting opcode value")
	)
	flag.Parse()

	if err := generateTables(*functions, *opcodeOut, *gasOut, *base); err != nil {
		log.Fatalf("opcode generation failed: %v", err)
	}
	absOp, _ := filepath.Abs(*opcodeOut)
	absGas, _ := filepath.Abs(*gasOut)
	log.Printf("wrote %s and %s", absOp, absGas)
}
