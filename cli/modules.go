package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	synn "synnergy"
)

// ModuleDefinition captures an enterprise CLI surface and the opcodes required
// to operate it. The metadata feeds the CLI, JavaScript control plane and the
// detailed whitepaper architecture guide.
type ModuleDefinition struct {
	Name        string
	Command     string
	Category    string
	Description string
	Opcodes     []string
}

// OpcodeStatus reports gas cost coverage for an opcode that underpins a CLI
// module. Missing opcodes are surfaced through the CLI and UI so operators can
// spot documentation drift before deploying upgrades.
type OpcodeStatus struct {
	Name string `json:"name"`
	Cost uint64 `json:"cost"`
}

// ModuleStatus is exported through the modules CLI command and the web control
// plane for Stage 81. It aligns CLI, VM, consensus, wallet and node workflows
// with the runtime gas schedule to guarantee fault-tolerant orchestration.
type ModuleStatus struct {
	Name           string         `json:"name"`
	Command        string         `json:"command"`
	Category       string         `json:"category"`
	Description    string         `json:"description"`
	Opcodes        []OpcodeStatus `json:"opcodes"`
	MissingOpcodes []string       `json:"missingOpcodes,omitempty"`
}

var moduleDefinitions = []ModuleDefinition{
	{
		Name:        "Consensus Control",
		Command:     "consensus",
		Category:    "consensus",
		Description: "Adaptive PoW/PoS orchestration with dynamic thresholds and validator policy enforcement.",
		Opcodes: []string{
			"MineBlock",
			"AdjustWeights",
			"TransitionThreshold",
			"SetAvailability",
			"SetPoWRewards",
		},
	},
	{
		Name:        "Virtual Machine",
		Command:     "simplevm",
		Category:    "execution",
		Description: "Lifecycle management for the deterministic Synnergy VM including execution timeouts and gas metering.",
		Opcodes: []string{
			"VMCreate",
			"VMStart",
			"VMStop",
			"VMStatus",
			"VMExec",
		},
	},
	{
		Name:        "Wallet & Identity",
		Command:     "wallet",
		Category:    "wallet",
		Description: "Secure wallet generation, encryption and identity flows with audited gas accounting.",
		Opcodes: []string{
			"WalletNew",
			"VerifySignature",
			"IDWalletRegister",
		},
	},
	{
		Name:        "Validator Nodes",
		Command:     "node",
		Category:    "node",
		Description: "Stake assignment, slashing, rehabilitation and block production for validator fleets.",
		Opcodes: []string{
			"NodeInfo",
			"NodeStake",
			"NodeSlash",
			"NodeRehab",
			"NodeAddTx",
			"NodeMempool",
			"NodeMine",
		},
	},
	{
		Name:        "Authority Governance",
		Command:     "authority",
		Category:    "governance",
		Description: "Authority node registration, voting and term renewal backed by digital signatures and quorum rules.",
		Opcodes: []string{
			"AuthorityApplyVote",
			"RenewAuthorityTerm",
			"UpdateMemberRole",
		},
	},
	{
		Name:        "Module Catalogue",
		Command:     "modules",
		Category:    "cli",
		Description: "Stage 81 module inventory with gas catalogue inspection for CLI, VM and web function web integrations.",
		Opcodes: []string{
			"ModuleCatalogueList",
			"ModuleCatalogueInspect",
		},
	},
}

// ModuleCatalogue returns a defensive copy of the module definitions. Tests and
// the JavaScript UI consume this to guarantee deterministic ordering.
func ModuleCatalogue() []ModuleDefinition {
	out := make([]ModuleDefinition, len(moduleDefinitions))
	copy(out, moduleDefinitions)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ModuleStatuses resolves opcode coverage for each module definition.
func ModuleStatuses() []ModuleStatus {
	defs := ModuleCatalogue()
	statuses := make([]ModuleStatus, 0, len(defs))
	for _, def := range defs {
		statuses = append(statuses, resolveModule(def))
	}
	return statuses
}

func resolveModule(def ModuleDefinition) ModuleStatus {
	status := ModuleStatus{
		Name:        def.Name,
		Command:     def.Command,
		Category:    def.Category,
		Description: def.Description,
		Opcodes:     make([]OpcodeStatus, 0, len(def.Opcodes)),
	}
	missing := make([]string, 0)
	for _, name := range def.Opcodes {
		cost := synn.GasCost(name)
		status.Opcodes = append(status.Opcodes, OpcodeStatus{Name: name, Cost: cost})
		if !synn.HasOpcode(name) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		status.MissingOpcodes = missing
	}
	return status
}

func newModulesCommand() *cobra.Command {
	modulesCmd := &cobra.Command{
		Use:   "modules",
		Short: "Discover enterprise CLI modules and their gas coverage",
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List modules, CLI entry points and opcode coverage",
		Run: func(cmd *cobra.Command, args []string) {
			statuses := ModuleStatuses()
			if jsonOutput {
				printOutput(statuses)
				return
			}
			printOutput(renderModuleTable(statuses))
		},
	}
	modulesCmd.AddCommand(listCmd)
	return modulesCmd
}

func init() {
	rootCmd.AddCommand(newModulesCommand())
}

func renderModuleTable(statuses []ModuleStatus) string {
	if len(statuses) == 0 {
		return "no modules registered"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%-24s %-12s %-12s %-10s %-s\n", "MODULE", "COMMAND", "CATEGORY", "OPCODES", "GAS COSTS")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 80))
	for _, status := range statuses {
		gasValues := make([]string, len(status.Opcodes))
		for i, op := range status.Opcodes {
			gasValues[i] = fmt.Sprintf("%s=%d", op.Name, op.Cost)
		}
		fmt.Fprintf(&b, "%-24s %-12s %-12s %-10d %s\n", status.Name, status.Command, status.Category, len(status.Opcodes), strings.Join(gasValues, ", "))
		if len(status.MissingOpcodes) > 0 {
			fmt.Fprintf(&b, "  ! missing documentation: %s\n", strings.Join(status.MissingOpcodes, ", "))
		}
	}
	return b.String()
}
