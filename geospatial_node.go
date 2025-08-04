package synnergy

import (
	"sync"
	"time"
)

// GeoRecord stores a geospatial data point.
type GeoRecord struct {
	Subject   string
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

// GeospatialNode collects and retrieves geospatial information for subjects.
type GeospatialNode struct {
	mu      sync.RWMutex
	records map[string][]GeoRecord
}

// NewGeospatialNode constructs a new GeospatialNode instance.
func NewGeospatialNode() *GeospatialNode {
	return &GeospatialNode{records: make(map[string][]GeoRecord)}
}

// Record stores a geospatial data point for the specified subject.
func (n *GeospatialNode) Record(subject string, lat, lon float64) {
	n.mu.Lock()
	n.records[subject] = append(n.records[subject], GeoRecord{
		Subject:   subject,
		Latitude:  lat,
		Longitude: lon,
		Timestamp: time.Now().UTC(),
	})
	n.mu.Unlock()
}

// History returns a copy of all recorded locations for a subject.
func (n *GeospatialNode) History(subject string) []GeoRecord {
	n.mu.RLock()
	defer n.mu.RUnlock()
	recs := n.records[subject]
	out := make([]GeoRecord, len(recs))
	copy(out, recs)
	return out
}
