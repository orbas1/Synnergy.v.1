package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

type stateStore struct {
	balances map[core.Address]uint64
	kv       map[string][]byte
}

type stateIter struct {
	keys []string
	idx  int
	kv   map[string][]byte
}

func (m *stateIter) Next() bool {
	if m.idx >= len(m.keys) {
		return false
	}
	m.idx++
	return true
}

func (m *stateIter) Value() []byte {
	return m.kv[m.keys[m.idx-1]]
}

func newStateStore() *stateStore {
	return &stateStore{balances: make(map[core.Address]uint64), kv: make(map[string][]byte)}
}

func (m *stateStore) Transfer(from, to core.Address, amount uint64) error {
	if m.balances[from] < amount {
		return errors.New("insufficient balance")
	}
	m.balances[from] -= amount
	m.balances[to] += amount
	return nil
}

func (m *stateStore) SetState(key, value []byte) {
	m.kv[string(key)] = value
}

func (m *stateStore) GetState(key []byte) ([]byte, error) {
	v, ok := m.kv[string(key)]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *stateStore) HasState(key []byte) (bool, error) {
	_, ok := m.kv[string(key)]
	return ok, nil
}

func (m *stateStore) PrefixIterator(prefix []byte) core.StateIterator {
	keys := make([]string, 0)
	for k := range m.kv {
		if strings.HasPrefix(k, string(prefix)) {
			keys = append(keys, k)
		}
	}
	return &stateIter{keys: keys, kv: m.kv}
}

func (m *stateStore) BalanceOf(addr core.Address) uint64 {
	return m.balances[addr]
}

var state = newStateStore()

func init() {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "In-memory StateRW utilities",
	}

	setCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Args:  cobra.ExactArgs(2),
		Short: "Set a key/value pair",
		Run: func(cmd *cobra.Command, args []string) {
			state.SetState([]byte(args[0]), []byte(args[1]))
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <key>",
		Args:  cobra.ExactArgs(1),
		Short: "Get value for a key",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := state.GetState([]byte(args[0]))
			if err != nil {
				fmt.Println("not found")
				return
			}
			fmt.Println(string(v))
		},
	}

	hasCmd := &cobra.Command{
		Use:   "has <key>",
		Args:  cobra.ExactArgs(1),
		Short: "Check if key exists",
		Run: func(cmd *cobra.Command, args []string) {
			ok, _ := state.HasState([]byte(args[0]))
			fmt.Println(ok)
		},
	}

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Transfer balance",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := state.Transfer(core.Address(args[0]), core.Address(args[1]), amt); err != nil {
				fmt.Println(err)
			}
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show balance",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(state.BalanceOf(core.Address(args[0])))
		},
	}

	iterCmd := &cobra.Command{
		Use:   "iterate <prefix>",
		Args:  cobra.ExactArgs(1),
		Short: "List key/values with prefix",
		Run: func(cmd *cobra.Command, args []string) {
			it := state.PrefixIterator([]byte(args[0]))
			for it.Next() {
				fmt.Println(string(it.Value()))
			}
		},
	}

	cmd.AddCommand(setCmd, getCmd, hasCmd, transferCmd, balanceCmd, iterCmd)
	rootCmd.AddCommand(cmd)
}
