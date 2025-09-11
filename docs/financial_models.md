# Financial Modelling for Synthron Coin

This document outlines common algorithms for forecasting the price of the Synthron coin and provides guidance on implementing predictions within the Synnergy blockchain.

## Algorithms

- **Moving Average**: Smooths short‑term fluctuations and projects future prices using the average of the last *N* observations.
- **Linear Regression**: Fits a straight line to historical price data and extrapolates into the future.
- **Autoregressive (AR) Models**: Predict future values based on previous values in the series. The provided code includes an AR(1) example.
- **ARIMA**: Combines autoregression, differencing, and moving averages. Suitable for non‑stationary series.
- **LSTM Neural Networks**: Recurrent neural networks capable of capturing complex temporal dynamics.
- **Monte Carlo Simulation**: Uses random sampling to simulate possible future price paths and derive probabilistic forecasts.

## Running Forecasts

The package `financial_prediction.go` provides simple implementations of moving average, linear regression, and AR(1) models. Example usage:

```go
prices := []float64{10, 12, 14, 16, 18}
model := MovingAverageModel{Window: 3}
forecast, _ := model.Forecast(prices, 6) // 6‑month forecast
```

To integrate forecasts on-chain, store historical price data in a smart contract or oracle. Nodes can run the forecasting models off-chain and submit results back to the blockchain for transparency and further processing.

## Limitations

Forecast accuracy depends on data quality and market behaviour; backtest models regularly and avoid relying on a single technique for critical decisions.

