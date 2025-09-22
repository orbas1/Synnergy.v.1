package core

import (
        "context"
        "fmt"
        "math"
        "testing"
        "time"
)

func TestThreshold(t *testing.T) {
	sc := NewSynnergyConsensus()
	if sc.Threshold(2, 3) != sc.Alpha*2+sc.Beta*3 {
		t.Fatalf("threshold calculation incorrect")
	}
}

func TestAdjustWeightsAndAvailability(t *testing.T) {
	sc := NewSynnergyConsensus()
	sc.SetAvailability(true, false, true)
	sc.AdjustWeights(0.5, 0.5)
	weights := sc.WeightsSnapshot()
	if weights.PoS != 0 {
		t.Fatalf("PoS weight should be zero when unavailable")
	}
	total := weights.PoW + weights.PoS + weights.PoH
	if math.Abs(total-1) > 1e-9 {
		t.Fatalf("weights not normalized")
	}
}

func TestTransitionThreshold(t *testing.T) {
	sc := NewSynnergyConsensus()
	tt := sc.TransitionThreshold(1, 2, 3)
	expected := sc.Tload(1) + sc.Tsecurity(2) + sc.Tstake(3)
	if tt != expected {
		t.Fatalf("transition threshold mismatch")
	}
}

func TestDifficultyAdjust(t *testing.T) {
	sc := NewSynnergyConsensus()
	if sc.DifficultyAdjust(1, 20, 10) != 2 {
		t.Fatalf("difficulty adjust incorrect")
	}
}

func TestSelectValidator(t *testing.T) {
	sc := NewSynnergyConsensus()
	stakes := map[string]uint64{"a": 1, "b": 1}
	addr := sc.SelectValidator("seed", stakes)
	if _, ok := stakes[addr]; !ok {
		t.Fatalf("unexpected validator: %s", addr)
	}
	if sc.SelectValidator("seed", map[string]uint64{}) != "" {
		t.Fatalf("expected empty string when no stakes")
	}
}

func TestSelectValidatorMajorityStake(t *testing.T) {
	sc := NewSynnergyConsensus()
	stakes := map[string]uint64{"a": 60, "b": 40}
	wins := map[string]int{"a": 0, "b": 0}
	for i := 0; i < 200; i++ {
		seed := fmt.Sprintf("seed-%d", i)
		addr := sc.SelectValidator(seed, stakes)
		if addr == "" {
			t.Fatalf("expected validator selection for seed %s", seed)
		}
		if _, ok := wins[addr]; !ok {
			t.Fatalf("unexpected validator: %s", addr)
		}
		wins[addr]++
	}
	if wins["a"] <= wins["b"] {
		t.Fatalf("expected majority stakeholder to win more often: %+v", wins)
	}
}

func TestFinalizeBlockRewards(t *testing.T) {
        sc := NewSynnergyConsensus()
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "v1")
	b := NewBlock([]*SubBlock{sb}, "")
	vm := NewValidatorManager(1)
	_ = vm.Add(context.Background(), "v1", 5)
	votes := map[string]bool{"v1": true, "v2": true, "v3": false}
	if !sc.FinalizeBlock(b, votes, vm, 2) {
		t.Fatalf("expected block to finalize")
	}
	if !b.Finalized {
		t.Fatalf("block not marked finalized")
	}
	if vm.Stake("v1") != 7 {
		t.Fatalf("reward not applied")
	}
}

func TestChooseChainPrefersWellFinalizedChain(t *testing.T) {
        sc := NewSynnergyConsensus()
        now := time.Now().Unix()

        chainA := []*Block{
                testConsensusBlock(nil, now, true, "valA1", 5),
        }
        chainA = append(chainA, testConsensusBlock(chainA[len(chainA)-1], now+12, true, "valA2", 6))
        chainA = append(chainA, testConsensusBlock(chainA[len(chainA)-1], now+24, true, "valA3", 6))

        chainB := []*Block{
                testConsensusBlock(nil, now-600, false, "valB", 2),
        }
        chainB = append(chainB, testConsensusBlock(chainB[len(chainB)-1], now-420, false, "valB", 2))
        chainB = append(chainB, testConsensusBlock(chainB[len(chainB)-1], now-240, false, "valB", 2))
        chainB = append(chainB, testConsensusBlock(chainB[len(chainB)-1], now-60, false, "valB", 1))

        selected := sc.ChooseChain([][]*Block{chainB, chainA})
        if selected == nil || selected[0] != chainA[0] {
                t.Fatalf("expected chainA to be selected")
        }
}

func TestChooseChainSkipsInvalidCandidates(t *testing.T) {
        sc := NewSynnergyConsensus()
        now := time.Now().Unix()

        valid := []*Block{
                testConsensusBlock(nil, now, true, "val1", 3),
                testConsensusBlock(nil, now+20, true, "val2", 3),
        }
        // fix prev hash linkage for valid chain
        valid[1].PrevHash = valid[0].Hash
        valid[1].Hash = valid[1].HeaderHash(valid[1].Nonce)

        invalid := []*Block{
                testConsensusBlock(nil, now+100, false, "valX", 1),
                testConsensusBlock(nil, now+90, false, "valX", 1), // timestamp regression
        }

        selected := sc.ChooseChain([][]*Block{invalid, valid})
        if selected == nil || selected[0] != valid[0] {
                t.Fatalf("expected valid chain to be selected when competitor invalid")
        }
}

func testConsensusBlock(prev *Block, timestamp int64, finalized bool, validator string, txCount int) *Block {
        txs := make([]*Transaction, txCount)
        for i := range txs {
                txs[i] = &Transaction{
                        ID:     fmt.Sprintf("%s-tx-%d-%d", validator, timestamp, i),
                        From:   validator,
                        To:     fmt.Sprintf("recipient-%d", i),
                        Amount: 1,
                }
        }
        sb := NewSubBlock(txs, validator)
        if timestamp > 0 {
                sb.Timestamp = timestamp - 1
                sb.PohHash = sb.Hash()
                sb.Signature = signSubBlock(sb.Validator, sb.PohHash)
        }
        blk := NewBlock([]*SubBlock{sb}, "")
        if prev != nil {
                blk.PrevHash = prev.Hash
        }
        if timestamp > 0 {
                blk.Timestamp = timestamp
        }
        blk.Nonce = uint64(txCount + 1)
        blk.Hash = blk.HeaderHash(blk.Nonce)
        blk.Finalized = finalized
        return blk
}

func TestValidateSubBlock(t *testing.T) {
	sc := NewSynnergyConsensus()
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "val")
	if !sc.ValidateSubBlock(sb) {
		t.Fatalf("expected valid sub-block")
	}
	sb.Signature = "bad"
	if sc.ValidateSubBlock(sb) {
		t.Fatalf("expected invalid sub-block")
	}
}

func TestValidateBlock(t *testing.T) {
	sc := NewSynnergyConsensus()
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "val")
	block := NewBlock([]*SubBlock{sb}, "")
	sc.MineBlock(block, 1)
	if !sc.ValidateBlock(block) {
		t.Fatalf("expected valid block")
	}
	block.SubBlocks[0].Transactions = nil
	if sc.ValidateBlock(block) {
		t.Fatalf("expected invalid block")
	}
}

func TestValidateSubBlockRegulatory(t *testing.T) {
	sc := NewSynnergyConsensus()
	mgr := NewRegulatoryManager()
	mgr.AddRegulation(Regulation{ID: "r1", MaxAmount: 10})
	rn := NewRegulatoryNode("rn", mgr)
	sc.SetRegulatoryNode(rn)
	w, err := NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	rn.RegisterWallet(w)

	tx := NewTransaction(w.Address, "bob", 5, 0, 0)
	if _, err := w.Sign(tx); err != nil {
		t.Fatalf("sign: %v", err)
	}
	sb := NewSubBlock([]*Transaction{tx}, "val")
	if !sc.ValidateSubBlock(sb) {
		t.Fatalf("expected valid sub-block")
	}

	tx2 := NewTransaction(w.Address, "bob", 20, 0, 0)
	if _, err := w.Sign(tx2); err != nil {
		t.Fatalf("sign2: %v", err)
	}
	sb2 := NewSubBlock([]*Transaction{tx2}, "val")
	if sc.ValidateSubBlock(sb2) {
		t.Fatalf("expected regulatory rejection")
	}
}

func TestValidateSubBlockWithoutRegNode(t *testing.T) {
	sc := NewSynnergyConsensus()
	mgr := NewRegulatoryManager()
	mgr.AddRegulation(Regulation{ID: "r1", MaxAmount: 10})
	rn := NewRegulatoryNode("rn", mgr)
	sc.SetRegulatoryNode(rn)

	tx := NewTransaction("alice", "bob", 20, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "val")
	if sc.ValidateSubBlock(sb) {
		t.Fatalf("expected rejection with regulatory node")
	}
	sc.SetRegulatoryNode(nil)
	if !sc.ValidateSubBlock(sb) {
		t.Fatalf("expected approval without regulatory node")
	}
}
