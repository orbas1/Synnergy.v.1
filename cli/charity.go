package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"synnergy/core"
)

type dummyElectorate struct{}

func (dummyElectorate) IsIDTokenHolder(addr core.Address) bool { return true }

type memState struct {
	balances map[core.Address]uint64
	store    map[string][]byte
}

func newMemState() *memState {
	return &memState{balances: map[core.Address]uint64{}, store: map[string][]byte{}}
}

func (m *memState) Transfer(from, to core.Address, amount uint64) error {
	m.balances[from] -= amount
	m.balances[to] += amount
	return nil
}

func (m *memState) SetState(key, value []byte)          { m.store[string(key)] = value }
func (m *memState) GetState(key []byte) ([]byte, error) { return m.store[string(key)], nil }
func (m *memState) HasState(key []byte) (bool, error)   { _, ok := m.store[string(key)]; return ok, nil }
func (m *memState) BalanceOf(addr core.Address) uint64  { return m.balances[addr] }

type memIter struct {
	keys  []string
	idx   int
	store map[string][]byte
}

func (m *memState) PrefixIterator(prefix []byte) core.StateIterator {
	p := string(prefix)
	keys := []string{}
	for k := range m.store {
		if strings.HasPrefix(k, p) {
			keys = append(keys, k)
		}
	}
	return &memIter{keys: keys, idx: -1, store: m.store}
}

func (it *memIter) Next() bool {
	it.idx++
	return it.idx < len(it.keys)
}

func (it *memIter) Value() []byte { return it.store[it.keys[it.idx]] }

var charityState = newMemState()
var charityPool = core.NewCharityPool(logrus.New(), charityState, dummyElectorate{}, time.Now())

func init() {
	poolCmd := &cobra.Command{Use: "charity_pool", Short: "Charity pool operations"}

	registerCmd := &cobra.Command{
		Use:   "register [addr] [category] [name]",
		Args:  cobra.ExactArgs(3),
		Short: "Register a charity with the pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := core.Address(args[0])
			catVal, _ := strconv.Atoi(args[1])
			return charityPool.Register(addr, args[2], core.CharityCategory(catVal))
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [voterAddr] [charityAddr]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for a charity during the cycle.",
		RunE: func(cmd *cobra.Command, args []string) error {
			voter := core.Address(args[0])
			charity := core.Address(args[1])
			return charityPool.Vote(voter, charity)
		},
	}

	tickCmd := &cobra.Command{
		Use:   "tick [timestamp]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Manually trigger pool cron tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			ts := time.Now()
			if len(args) == 1 {
				if v, err := strconv.ParseInt(args[0], 10, 64); err == nil {
					ts = time.Unix(v, 0)
				}
			}
			charityPool.Tick(ts)
		},
	}

	registrationCmd := &cobra.Command{
		Use:   "registration [addr] [cycle]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Show registration info for a charity.",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := core.Address(args[0])
			var cycle uint64
			if len(args) == 2 {
				c, _ := strconv.ParseUint(args[1], 10, 64)
				cycle = c
			}
			reg, ok, err := charityPool.GetRegistration(cycle, addr)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.Marshal(reg)
			fmt.Println(string(b))
			return nil
		},
	}

	winnersCmd := &cobra.Command{
		Use:   "winners [cycle]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "List winning charities for a cycle.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cycle uint64
			if len(args) == 1 {
				c, _ := strconv.ParseUint(args[0], 10, 64)
				cycle = c
			}
			ws, err := charityPool.Winners(cycle)
			if err != nil {
				return err
			}
			for _, a := range ws {
				fmt.Println(a)
			}
			return nil
		},
	}

	poolCmd.AddCommand(registerCmd, voteCmd, tickCmd, registrationCmd, winnersCmd)

	mgmtCmd := &cobra.Command{Use: "charity_mgmt", Short: "Charity pool management"}

	donateCmd := &cobra.Command{
		Use:   "donate [from] [amt]",
		Args:  cobra.ExactArgs(2),
		Short: "Donate tokens to the charity pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			from := core.Address(args[0])
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			charityState.balances[from] += amt
			return charityPool.Deposit(from, amt)
		},
	}

	withdrawCmd := &cobra.Command{
		Use:   "withdraw [to] [amt]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw internal charity funds.",
		RunE: func(cmd *cobra.Command, args []string) error {
			to := core.Address(args[0])
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			return charityState.Transfer(core.InternalCharityAccount, to, amt)
		},
	}

	balancesCmd := &cobra.Command{
		Use:   "balances",
		Short: "Show pool and internal balances.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("pool: %d internal: %d\n", charityState.BalanceOf(core.CharityPoolAccount), charityState.BalanceOf(core.InternalCharityAccount))
		},
	}

	mgmtCmd.AddCommand(donateCmd, withdrawCmd, balancesCmd)

	rootCmd.AddCommand(poolCmd, mgmtCmd)
}
