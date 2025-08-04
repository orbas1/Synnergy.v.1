package synnergy

import (
    "bytes"
    "testing"
)

func TestAIServiceEndToEnd(t *testing.T) {
    s := NewAIService()
    // PredictFraud
    score, err := s.PredictFraud([]byte(`{"foo":1}`))
    if err != nil || score < 0 || score > 1 {
        t.Fatalf("predict: %v %f", err, score)
    }
    // OptimiseBaseFee
    fee, err := s.OptimiseBaseFee([]byte(`{"avgGasPrice":100}`))
    if err != nil || fee != 110 {
        t.Fatalf("opt fee: %v %d", err, fee)
    }
    // ForecastVolume
    vol, err := s.ForecastVolume([]byte(`{"recentTxs":1000}`))
    if err != nil || vol != 1050 {
        t.Fatalf("forecast: %v %d", err, vol)
    }
    // Publish and fetch model
    hash, err := s.PublishModel("cid1", 100)
    if err != nil {
        t.Fatalf("publish: %v", err)
    }
    meta, ok := s.FetchModel(hash)
    if !ok || meta.CID != "cid1" {
        t.Fatalf("fetch mismatch")
    }
    // List model and buy
    listingID := s.ListModel(hash, "cid1", "seller", 500)
    escrowID, err := s.BuyModel(listingID, "buyer", 500)
    if err != nil {
        t.Fatalf("buy model: %v", err)
    }
    // Rent model
    rentID, err := s.RentModel(listingID, "renter", 1, 500)
    if err != nil || rentID == "" {
        t.Fatalf("rent model: %v", err)
    }
    // Release escrow
    if err := s.ReleaseEscrow(escrowID); err != nil {
        t.Fatalf("release escrow: %v", err)
    }
    s.mu.RLock()
    if !s.escrows[escrowID].Released {
        t.Fatalf("escrow not released")
    }
    s.mu.RUnlock()
}

func TestModelMarketplace(t *testing.T) {
    m := NewModelMarketplace()
    id := m.AddListing("hash", "cid", "seller", 100)
    l, ok := m.Get(id)
    if !ok || l.Price != 100 {
        t.Fatalf("get listing")
    }
    if err := m.Update(id, 200); err != nil {
        t.Fatalf("update: %v", err)
    }
    if err := m.Remove(id, "seller"); err != nil {
        t.Fatalf("remove: %v", err)
    }
}

func TestTrainingManager(t *testing.T) {
    tm := NewTrainingManager()
    id := tm.Start("data", "model")
    if _, ok := tm.Status(id); !ok {
        t.Fatalf("status not found")
    }
    if len(tm.List()) != 1 {
        t.Fatalf("list size")
    }
    if err := tm.Cancel(id); err != nil {
        t.Fatalf("cancel: %v", err)
    }
    job, _ := tm.Status(id)
    if job.Status != "cancelled" {
        t.Fatalf("expected cancelled")
    }
}

func TestInferenceEngine(t *testing.T) {
    e := NewInferenceEngine()
    e.LoadModel("m1", []byte("model"))
    out, err := e.Run("m1", []byte("input"))
    if err != nil || len(out) == 0 {
        t.Fatalf("run: %v %d", err, len(out))
    }
    res := e.Analyse([]string{"tx1", "tx2"})
    if len(res) != 2 {
        t.Fatalf("analyse len")
    }
    if res[0].Score < 0 || res[0].Score > 1 {
        t.Fatalf("score range")
    }
}

func TestDriftMonitor(t *testing.T) {
    d := NewDriftMonitor()
    d.UpdateBaseline("m1", 0.5)
    if !d.HasDrift("m1", 0.8, 0.2) {
        t.Fatalf("expected drift")
    }
    if d.HasDrift("m1", 0.55, 0.2) {
        t.Fatalf("unexpected drift")
    }
}

func TestSecureStorage(t *testing.T) {
    s := NewSecureStorage()
    key := bytes.Repeat([]byte{1}, 32)
    if err := s.Store("hash", []byte("data"), key); err != nil {
        t.Fatalf("store: %v", err)
    }
    data, err := s.Retrieve("hash", key)
    if err != nil || !bytes.Equal(data, []byte("data")) {
        t.Fatalf("retrieve: %v", err)
    }
}

func TestAnomalyDetector(t *testing.T) {
    a := NewAnomalyDetector(2)
    a.Update(10)
    a.Update(12)
    a.Update(11)
    if !a.IsAnomalous(20) {
        t.Fatalf("expected anomaly")
    }
    if a.IsAnomalous(12) {
        t.Fatalf("unexpected anomaly")
    }
}

