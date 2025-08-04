package core

// CharityPool – 5% cut from every gas fee routed to on-chain philanthropy.
//
// Mechanics
// ---------
// * **Income**: Tx execution calls `charityPool.Deposit(fee)` with exactly 5 % of gas
//   burned. Funds accumulate in `CharityPoolAccount`.
// * **Daily distribution**: 24 h cron (triggered via `Tick`) pays out 50 % of the
//   contract balance – 50 % split equally across the **30 winning charities** of
//   the current 90-day cycle (max 5 per category) and 50 % sent to an internal
//   charity wallet (`InternalCharityAccount`).
// * **Cycle**: 90-day windows numbered from genesis timestamp; charities must
//   **register 30 d before cycle end**; community (ID-token holders) vote in the
//   last 15 d; top 5 per category win.
//
// Dependencies: ledger (state+balance), authority (ID token check), security
// (sig verify for charity wallet keys). All times in UTC.

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

//---------------------------------------------------------------------
// Categories
//---------------------------------------------------------------------

type CharityCategory uint8

const (
	HungerRelief CharityCategory = iota + 1
	ChildrenHelp
	WildlifeHelp
	SeaSupport
	DisasterSupport
	WarSupport
)

func (c CharityCategory) String() string {
	switch c {
	case HungerRelief:
		return "HungerRelief"
	case ChildrenHelp:
		return "ChildrenHelp"
	case WildlifeHelp:
		return "WildlifeHelp"
	case SeaSupport:
		return "SeaSupport"
	case DisasterSupport:
		return "DisasterSupport"
	case WarSupport:
		return "WarSupport"
	default:
		return "Unknown"
	}
}

//---------------------------------------------------------------------
// Pool configuration
//---------------------------------------------------------------------

const (
	cycleDuration      = 90 * 24 * time.Hour
	registrationCutoff = 30 * 24 * time.Hour
	votingWindow       = 15 * 24 * time.Hour
	dailyPayout        = 24 * time.Hour
)

var (
	CharityPoolAccount     Address
	InternalCharityAccount Address
)

func init() {
	var err error

	CharityPoolAccount, err = StringToAddress("0x43686172697479426F7800000000000000000000")
	if err != nil {
		panic("invalid CharityPoolAccount: " + err.Error())
	}

	InternalCharityAccount, err = StringToAddress("0x496E7443686172697479426F7800000000000000")
	if err != nil {
		panic("invalid InternalCharityAccount: " + err.Error())
	}
}

//---------------------------------------------------------------------
// Expected external APIs
//---------------------------------------------------------------------

type electorate interface {
	IsIDTokenHolder(addr Address) bool
}

//---------------------------------------------------------------------
// CharityPool struct
//---------------------------------------------------------------------

func NewCharityPool(lg *logrus.Logger, led StateRW, el electorate, genesis time.Time) *CharityPool {
	return &CharityPool{logger: lg, led: led, vote: el, genesis: genesis}
}

//---------------------------------------------------------------------
// Deposit – called by VM interpreter when charging gas fees (5% of total fee).
//---------------------------------------------------------------------

func (cp *CharityPool) Deposit(from Address, amount uint64) error {
	return cp.led.Transfer(from, CharityPoolAccount, amount)
}

//---------------------------------------------------------------------
// Register – charity registers for next cycle (open >=30d before cycle end).
//---------------------------------------------------------------------

func (cp *CharityPool) Register(addr Address, name string, cat CharityCategory) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	now := time.Now().UTC()
	cycle := cp.currentCycle(now.Add(registrationCutoff)) // registering for *next* cycle if within cutoff
	if len(name) == 0 || len(name) > 64 {
		return errors.New("invalid name")
	}
	if cat < HungerRelief || cat > WarSupport {
		return errors.New("invalid category")
	}
	key := regKey(cycle, addr)
	if exists, _ := cp.led.HasState(key); exists {
		return errors.New("already registered")
	}
	// Ensure category cap <=5
	count := cp.countCategoryRegistrations(cycle, cat)
	if count >= 5 {
		return errors.New("category full for cycle")
	}
	r := CharityRegistration{Addr: addr, Name: name, Category: cat, Cycle: cycle}
	cp.led.SetState(key, mustJSON(r))
	cp.logger.Printf("charity %s registered for cycle %d cat %s", addr.Short(), cycle, cat)
	return nil
}

//---------------------------------------------------------------------
// Vote – ID-token holders vote for a charity during last 15d of cycle.
//---------------------------------------------------------------------

func (cp *CharityPool) Vote(voter, charity Address) error {
	if !cp.vote.IsIDTokenHolder(voter) {
		return errors.New("not verified voter")
	}
	now := time.Now().UTC()
	cycle := cp.currentCycle(now)
	cycleEnd := cp.cycleEnd(cycle)
	if cycleEnd.Sub(now) > votingWindow {
		return errors.New("voting not open yet")
	}
	regRaw, err := cp.led.GetState(regKey(cycle, charity))
	if err != nil || len(regRaw) == 0 {
		return errors.New("charity not registered this cycle")
	}
	// no double vote per cycle and voter
	vkey := voteKey(cycleHash(cycle), voter)
	if ok, _ := cp.led.HasState(vkey); ok {
		return errors.New("already voted this cycle")
	}
	cp.led.SetState(vkey, charity.Bytes())
	var reg CharityRegistration
	_ = json.Unmarshal(regRaw, &reg)
	reg.VoteCount++
	cp.led.SetState(regKey(cycle, charity), mustJSON(reg))
	return nil
}

type Voter interface {
	IsIDTokenHolder(Address) bool
}

func cycleHash(cycle uint64) Hash {
	var id Hash
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], cycle)
	h := sha256.Sum256(buf[:])
	copy(id[:], h[:])
	return id
}

//---------------------------------------------------------------------
// Tick – called each block; handles daily payouts and cycle transitions.
//---------------------------------------------------------------------

func (cp *CharityPool) Tick(ts time.Time) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Daily payout?
	if ts.Unix()-cp.lastDaily >= int64(dailyPayout.Seconds()) {
		cp.distributeDaily()
		cp.lastDaily = ts.Unix()
	}

	// At cycle end +1 block, compute winners & persist for next 90d payouts.
	if cp.isCycleBoundary(ts) {
		cp.finaliseCycle(cp.currentCycle(ts))
	}
}

//---------------------------------------------------------------------
// Internal helpers
//---------------------------------------------------------------------

func (cp *CharityPool) distributeDaily() {
	bal := cp.led.BalanceOf(CharityPoolAccount)
	if bal == 0 {
		return
	}
	half := bal / 2

	cycle := cp.currentCycle(time.Now().UTC())
	winners := cp.winnerList(cycle)
	per := uint64(0)
	if len(winners) > 0 {
		per = half / uint64(len(winners))
		for _, w := range winners {
			_ = cp.led.Transfer(CharityPoolAccount, w, per)
		}
	}
	_ = cp.led.Transfer(CharityPoolAccount, InternalCharityAccount, half)
	cp.logger.Printf("charity daily payout half=%d perCharity=%d to %d winners", half, per, len(winners))
}

func (cp *CharityPool) finaliseCycle(cycle uint64) {
	// Determine top 5 per category by votes.
	cats := make(map[CharityCategory][]CharityRegistration)
	iter := cp.led.PrefixIterator([]byte(fmt.Sprintf("charity:reg:%d:", cycle)))
	for iter.Next() {
		var r CharityRegistration
		_ = json.Unmarshal(iter.Value(), &r)
		cats[r.Category] = append(cats[r.Category], r)
	}
	var winners []Address
	for cat, list := range cats {
		sort.Slice(list, func(i, j int) bool { return list[i].VoteCount > list[j].VoteCount })
		n := 5
		if len(list) < 5 {
			n = len(list)
		}
		for i := 0; i < n; i++ {
			winners = append(winners, list[i].Addr)
		}
		cp.logger.Printf("cycle %d cat %s winners: %d", cycle, cat, n)
	}
	// Persist winners set for daily payouts (next cycle period)
	cp.led.SetState(winKey(cycle), mustJSON(winners))
}

func (cp *CharityPool) winnerList(cycle uint64) []Address {
	raw, _ := cp.led.GetState(winKey(cycle))
	if len(raw) == 0 {
		return nil
	}
	var out []Address
	_ = json.Unmarshal(raw, &out)
	return out
}

//---------------------------------------------------------------------
// Winners – exported accessor returning winners for a given cycle
//---------------------------------------------------------------------

func (cp *CharityPool) Winners(cycle uint64) ([]Address, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	list := cp.winnerList(cycle)
	if list == nil {
		return nil, fmt.Errorf("no winners recorded for cycle %d", cycle)
	}
	out := make([]Address, len(list))
	copy(out, list)
	return out, nil
}

//---------------------------------------------------------------------
// GetRegistration – return registration info if present
//---------------------------------------------------------------------

func (cp *CharityPool) GetRegistration(cycle uint64, addr Address) (CharityRegistration, bool, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	var reg CharityRegistration
	raw, err := cp.led.GetState(regKey(cycle, addr))
	if err != nil {
		return reg, false, err
	}
	if len(raw) == 0 {
		return reg, false, nil
	}
	if err := json.Unmarshal(raw, &reg); err != nil {
		return reg, false, err
	}
	return reg, true, nil
}

//---------------------------------------------------------------------
// Cycle maths
//---------------------------------------------------------------------

func (cp *CharityPool) currentCycle(t time.Time) uint64 {
	if t.Before(cp.genesis) {
		return 0
	}
	return uint64(t.Sub(cp.genesis) / cycleDuration)
}

func (cp *CharityPool) cycleEnd(cycle uint64) time.Time {
	return cp.genesis.Add(time.Duration(cycle+1) * cycleDuration)
}
func (cp *CharityPool) isCycleBoundary(ts time.Time) bool {
	return ts.UTC().Truncate(cycleDuration) == ts.UTC()
}

//---------------------------------------------------------------------
// Ledger key helpers
//---------------------------------------------------------------------

func regKey(cycle uint64, addr Address) []byte {
	return []byte(fmt.Sprintf("charity:reg:%d:%s", cycle, addr.Hex()))
}

func winKey(cycle uint64) []byte { return []byte(fmt.Sprintf("charity:winners:%d", cycle)) }

func (cp *CharityPool) countCategoryRegistrations(cycle uint64, cat CharityCategory) int {
	iter := cp.led.PrefixIterator([]byte(fmt.Sprintf("charity:reg:%d:", cycle)))
	count := 0
	for iter.Next() {
		var r CharityRegistration
		_ = json.Unmarshal(iter.Value(), &r)
		if r.Category == cat {
			count++
		}
	}
	return count
}

//---------------------------------------------------------------------
// END charity_pool.go
//---------------------------------------------------------------------
