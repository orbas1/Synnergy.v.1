package synnergy

import (
	"math"
	"testing"
)

func TestMovingAverageModel(t *testing.T) {
	model := MovingAverageModel{Window: 3}
	data := []float64{10, 12, 14, 16, 18}
	preds, err := model.Forecast(data, 2)
	if err != nil {
		t.Fatalf("forecast error: %v", err)
	}
	if math.Abs(preds[0]-16) > 1e-6 {
		t.Errorf("expected first prediction 16, got %v", preds[0])
	}
	if math.Abs(preds[1]-16.6666667) > 1e-6 {
		t.Errorf("expected second prediction ~16.67, got %v", preds[1])
	}
}

func TestLinearRegressionModel(t *testing.T) {
	model := LinearRegressionModel{}
	data := []float64{1, 2, 3, 4}
	preds, err := model.Forecast(data, 2)
	if err != nil {
		t.Fatalf("forecast error: %v", err)
	}
	if len(preds) != 2 {
		t.Fatalf("expected two predictions, got %d", len(preds))
	}
	if math.Abs(preds[0]-5) > 1e-6 {
		t.Errorf("expected first prediction 5, got %v", preds[0])
	}
	if math.Abs(preds[1]-6) > 1e-6 {
		t.Errorf("expected second prediction 6, got %v", preds[1])
	}
}

func TestAR1Model(t *testing.T) {
	model := AR1Model{}
	data := []float64{2, 4, 8}
	preds, err := model.Forecast(data, 2)
	if err != nil {
		t.Fatalf("forecast error: %v", err)
	}
	if len(preds) != 2 {
		t.Fatalf("expected two predictions, got %d", len(preds))
	}
	if math.Abs(preds[0]-16) > 1e-6 {
		t.Errorf("expected first prediction 16, got %v", preds[0])
	}
}

func TestForecastSeriesDefault(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	preds, err := ForecastSeries(nil, data, 1)
	if err != nil {
		t.Fatalf("forecast error: %v", err)
	}
	if len(preds) != 1 {
		t.Fatalf("expected one prediction, got %d", len(preds))
	}
}
