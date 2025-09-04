package synnergy

import (
	"errors"
)

// PriceModel defines interface for forecasting future prices.
type PriceModel interface {
	Forecast(data []float64, months int) ([]float64, error)
}

// MovingAverageModel predicts future prices using simple moving average.
type MovingAverageModel struct {
	Window int
}

// Forecast generates predictions for the next `months` months based on a moving average.
func (m MovingAverageModel) Forecast(data []float64, months int) ([]float64, error) {
	if m.Window <= 0 {
		return nil, errors.New("window must be positive")
	}
	if len(data) < m.Window {
		return nil, errors.New("not enough data for the specified window")
	}
	result := make([]float64, months)
	series := append([]float64{}, data...)
	for i := 0; i < months; i++ {
		sum := 0.0
		for j := len(series) - m.Window; j < len(series); j++ {
			sum += series[j]
		}
		next := sum / float64(m.Window)
		result[i] = next
		series = append(series, next)
	}
	return result, nil
}

// LinearRegressionModel predicts prices by fitting a line to historical data
// and extrapolating into the future. It uses ordinary least squares.
type LinearRegressionModel struct{}

// Forecast returns linear regression predictions for the next `months` months.
func (LinearRegressionModel) Forecast(data []float64, months int) ([]float64, error) {
	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}
	n := float64(len(data))
	var sumX, sumY, sumXY, sumX2 float64
	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return nil, errors.New("cannot compute regression")
	}
	b := (n*sumXY - sumX*sumY) / denom
	a := (sumY - b*sumX) / n
	result := make([]float64, months)
	for i := 0; i < months; i++ {
		x := float64(len(data) + i)
		result[i] = a + b*x
	}
	return result, nil
}

// AR1Model performs a simple autoregressive forecasting of order 1.
// It estimates the coefficient from data and projects future values.
type AR1Model struct{}

// Forecast predicts future prices using AR(1) model.
func (AR1Model) Forecast(data []float64, months int) ([]float64, error) {
	if len(data) < 2 {
		return nil, errors.New("need at least two data points")
	}
	var sumXX, sumXY float64
	for i := 1; i < len(data); i++ {
		x := data[i-1]
		y := data[i]
		sumXX += x * x
		sumXY += x * y
	}
	if sumXX == 0 {
		return nil, errors.New("zero variance in data")
	}
	phi := sumXY / sumXX
	last := data[len(data)-1]
	result := make([]float64, months)
	for i := 0; i < months; i++ {
		next := phi * last
		result[i] = next
		last = next
	}
	return result, nil
}

// ForecastSeries selects a model and returns predictions for `months` months.
// Passing a nil model defaults to a MovingAverageModel with window 3.
func ForecastSeries(model PriceModel, data []float64, months int) ([]float64, error) {
	if model == nil {
		model = MovingAverageModel{Window: 3}
	}
	return model.Forecast(data, months)
}
