package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	f, err := os.Open("functions_list.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	re := regexp.MustCompile(`func\s+([A-Za-z0-9_]+)\s*\(`)
	uniq := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			uniq[m[1]] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	funcs := make([]string, 0, len(uniq))
	for k := range uniq {
		funcs = append(funcs, k)
	}
	sort.Strings(funcs)

	opFile, err := os.Create("opcodes_list.md")
	if err != nil {
		panic(err)
	}
	defer opFile.Close()
	fmt.Fprintln(opFile, "| Function | Opcode |")
	fmt.Fprintln(opFile, "|---|---|")

	gasFile, err := os.Create("gas_table_list.md")
	if err != nil {
		panic(err)
	}
	defer gasFile.Close()
	fmt.Fprintln(gasFile, "| Function | Gas Cost |")
	fmt.Fprintln(gasFile, "|---|---|")

	base := 0x100000
	for i, fn := range funcs {
		opcode := base + i + 1
		fmt.Fprintf(opFile, "| `%s` | `0x%06X` |\n", fn, opcode)
		fmt.Fprintf(gasFile, "| `%s` | `1` |\n", fn)
	}
}
