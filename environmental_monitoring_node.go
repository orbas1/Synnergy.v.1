package synnergy

import (
	"strconv"
	"strings"
	"sync"
)

// EnvCondition evaluates sensor bytes and returns true when the action should trigger.
type EnvCondition struct {
	SensorID  string
	Operator  string // one of ">", "<", ">=", "<=", "=="
	Threshold float64
}

// Evaluate returns true if the sensor reading encoded in data satisfies the condition.
func (c EnvCondition) Evaluate(data []byte) bool {
	v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return false
	}
	switch c.Operator {
	case ">":
		return v > c.Threshold
	case "<":
		return v < c.Threshold
	case ">=":
		return v >= c.Threshold
	case "<=":
		return v <= c.Threshold
	case "==":
		return v == c.Threshold
	default:
		return false
	}
}

// EnvironmentalMonitoringNode aggregates external sensor data and evaluates registered conditions.
type EnvironmentalMonitoringNode struct {
	mu         sync.RWMutex
	conditions map[string]EnvCondition
}

// NewEnvironmentalMonitoringNode constructs a new monitoring node instance.
func NewEnvironmentalMonitoringNode() *EnvironmentalMonitoringNode {
	return &EnvironmentalMonitoringNode{conditions: make(map[string]EnvCondition)}
}

// SetCondition registers or updates a condition for a sensor.
func (n *EnvironmentalMonitoringNode) SetCondition(c EnvCondition) {
	n.mu.Lock()
	n.conditions[c.SensorID] = c
	n.mu.Unlock()
}

// Trigger evaluates incoming sensor data against a stored condition and returns true if it matches.
func (n *EnvironmentalMonitoringNode) Trigger(sensorID string, data []byte) bool {
	n.mu.RLock()
	cond, ok := n.conditions[sensorID]
	n.mu.RUnlock()
	if !ok {
		return false
	}
	return cond.Evaluate(data)
}
