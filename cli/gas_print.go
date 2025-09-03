package cli

import (
	"fmt"
	synn "synnergy"
)

// gasPrint outputs the configured gas price for the given opcode name.
// Stage 23 exposes gas awareness across consensus and governance CLI
// commands so operators can estimate costs before execution.
func gasPrint(name string) {
	fmt.Printf("gas cost: %d\n", synn.GasCost(name))
}
