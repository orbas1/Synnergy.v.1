package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"synnergy/pkg/version"
)

var (
	manifestFormat string
)

// manifestResponse captures the CLI tree in a serialisable structure so the web
// interface can offer rich command discovery and validation without shelling
// out to multiple help invocations.
type manifestResponse struct {
	Version     string            `json:"version"`
	GeneratedAt string            `json:"generated_at"`
	CommandSet  []manifestCommand `json:"commands"`
}

type manifestCommand struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Segments     []string          `json:"segments"`
	UseLine      string            `json:"use"`
	Short        string            `json:"short"`
	Long         string            `json:"long,omitempty"`
	Example      string            `json:"example,omitempty"`
	Aliases      []string          `json:"aliases,omitempty"`
	Hidden       bool              `json:"hidden"`
	Category     string            `json:"category"`
	Flags        []manifestFlag    `json:"flags"`
	Requirements manifestExtraInfo `json:"requirements"`
}

type manifestExtraInfo struct {
	Persistent bool `json:"persistent"`
}

type manifestFlag struct {
	Name       string `json:"name"`
	Shorthand  string `json:"shorthand,omitempty"`
	Usage      string `json:"usage"`
	Default    string `json:"default,omitempty"`
	ValueType  string `json:"type"`
	Persistent bool   `json:"persistent"`
	Required   bool   `json:"required"`
	Hidden     bool   `json:"hidden"`
}

var guiManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Emit a machine-readable CLI manifest for web and desktop clients",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifest := manifestResponse{
			Version:     version.Version,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			CommandSet:  collectManifest(RootCmd()),
		}

		var out []byte
		switch strings.ToLower(manifestFormat) {
		case "", "json":
			data, err := json.MarshalIndent(manifest, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal manifest: %w", err)
			}
			out = append(data, '\n')
		default:
			return fmt.Errorf("unsupported manifest format %q", manifestFormat)
		}

		if _, err := cmd.OutOrStdout().Write(out); err != nil {
			return fmt.Errorf("write manifest: %w", err)
		}
		return nil
	},
}

func init() {
	guiManifestCmd.Flags().StringVarP(&manifestFormat, "format", "f", "json", "Output format (json)")
	guiCmd.AddCommand(guiManifestCmd)
}

func collectManifest(root *cobra.Command) []manifestCommand {
	var cmds []manifestCommand
	if root == nil {
		return cmds
	}

	stack := []*cobra.Command{root}
	seen := make(map[*cobra.Command]struct{})

	for len(stack) > 0 {
		n := len(stack) - 1
		cmd := stack[n]
		stack = stack[:n]
		if cmd == nil {
			continue
		}
		if _, ok := seen[cmd]; ok {
			continue
		}
		seen[cmd] = struct{}{}

		if cmd != root {
			if cmd.Hidden {
				// Skip hidden helpers to keep the manifest focused on
				// supported UX flows.
				continue
			}

			segs := commandSegments(root, cmd)
			flags := collectFlags(cmd)
			category := "root"
			if len(segs) > 0 {
				category = segs[0]
			}

			entry := manifestCommand{
				Name:     cmd.Name(),
				Path:     strings.Join(segs, " "),
				Segments: segs,
				UseLine:  strings.TrimSpace(cmd.UseLine()),
				Short:    strings.TrimSpace(cmd.Short),
				Long:     strings.TrimSpace(cmd.Long),
				Example:  strings.TrimSpace(cmd.Example),
				Aliases:  append([]string(nil), cmd.Aliases...),
				Hidden:   cmd.Hidden,
				Category: category,
				Flags:    flags,
				Requirements: manifestExtraInfo{
					Persistent: cmd.PersistentFlags().HasFlags(),
				},
			}
			cmds = append(cmds, entry)
		}

		for _, child := range cmd.Commands() {
			stack = append(stack, child)
		}
	}

	sort.Slice(cmds, func(i, j int) bool {
		return strings.Join(cmds[i].Segments, " ") < strings.Join(cmds[j].Segments, " ")
	})
	return cmds
}

func collectFlags(cmd *cobra.Command) []manifestFlag {
	var flags []manifestFlag

	gather := func(fs *pflag.FlagSet, persistent bool) {
		if fs == nil {
			return
		}
		fs.VisitAll(func(f *pflag.Flag) {
			mf := manifestFlag{
				Name:       f.Name,
				Shorthand:  f.Shorthand,
				Usage:      strings.TrimSpace(f.Usage),
				Default:    f.DefValue,
				ValueType:  f.Value.Type(),
				Persistent: persistent,
				Required:   flagRequired(f),
				Hidden:     f.Hidden,
			}
			flags = append(flags, mf)
		})
	}

	gather(cmd.PersistentFlags(), true)
	gather(cmd.LocalFlags(), false)

	sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
	return flags
}

func flagRequired(f *pflag.Flag) bool {
	if f == nil {
		return false
	}
	const requiredKey = "cobra_annotation_bash_completion_one_required_flag"
	_, ok := f.Annotations[requiredKey]
	return ok
}

func commandSegments(root, cmd *cobra.Command) []string {
	if cmd == nil {
		return nil
	}
	raw := strings.Fields(cmd.CommandPath())
	if len(raw) > 0 && raw[0] == root.Name() {
		raw = raw[1:]
	}
	return raw
}
