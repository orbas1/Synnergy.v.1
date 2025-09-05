package synnergy

import "testing"

func TestInferenceEngineRun(t *testing.T) {
	e := NewInferenceEngine()
	e.LoadModel("m1", []byte("model"))
	out, err := e.Run("m1", []byte("input"))
	if err != nil || len(out) == 0 {
		t.Fatalf("run: %v", err)
	}
	res := e.Analyse([]string{"tx1", "tx2"})
	if len(res) != 2 || res[0].TxID != "tx1" {
		t.Fatalf("analyse results unexpected: %+v", res)
	}
	if res[0].Score < 0 || res[0].Score > 1 {
		t.Fatalf("score out of range")
	}
}
