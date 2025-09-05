package synnergy

import (
	"encoding/json"
	"testing"
)

func TestAIServicePredictAndPublish(t *testing.T) {
	svc := NewAIService()
	txJSON, _ := json.Marshal(map[string]any{"from": "a", "to": "b"})
	score, err := svc.PredictFraud(txJSON)
	if err != nil {
		t.Fatalf("predict: %v", err)
	}
	if score < 0 || score > 1 {
		t.Fatalf("score out of range: %f", score)
	}

	hash, err := svc.PublishModel("cid", 100)
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if _, ok := svc.FetchModel(hash); !ok {
		t.Fatalf("model not found")
	}
}
