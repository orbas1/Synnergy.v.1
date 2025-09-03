package nodes

import (
	"sort"
	"sync"
	"time"
)

// GeoRecord represents a single geospatial data point captured by the network.
type GeoRecord struct {
	Subject   string    // entity the record relates to
	Latitude  float64   // decimal degrees
	Longitude float64   // decimal degrees
	Timestamp time.Time // time the reading was taken
}

// GeospatialNodeInterface extends NodeInterface with geospatial capabilities.
type GeospatialNodeInterface interface {
	NodeInterface
	// Record stores a latitude/longitude pair for the given subject.
	Record(subject string, lat, lon float64) error
	// History returns the recorded locations for a subject ordered by time.
	History(subject string) []GeoRecord
}

// GeospatialNode provides an in-memory geospatial tracker backed by a
// BasicNode. Records are stored in-memory and suitable for testing and
// development scenarios.
type GeospatialNode struct {
	*BasicNode
	mu     sync.RWMutex
	points map[string][]GeoRecord
}

// NewGeospatialNode creates a new geospatial node.
func NewGeospatialNode(id Address) *GeospatialNode {
	return &GeospatialNode{BasicNode: NewBasicNode(id), points: make(map[string][]GeoRecord)}
}

// Record captures a latitude and longitude for the given subject.
func (n *GeospatialNode) Record(subject string, lat, lon float64) error {
	n.mu.Lock()
	rec := GeoRecord{Subject: subject, Latitude: lat, Longitude: lon, Timestamp: time.Now()}
	n.points[subject] = append(n.points[subject], rec)
	n.mu.Unlock()
	return nil
}

// History returns all recorded locations for a subject ordered by time.
func (n *GeospatialNode) History(subject string) []GeoRecord {
	n.mu.RLock()
	recs := n.points[subject]
	out := make([]GeoRecord, len(recs))
	copy(out, recs)
	n.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.Before(out[j].Timestamp) })
	return out
}

// Ensure GeospatialNode implements GeospatialNodeInterface.
var _ GeospatialNodeInterface = (*GeospatialNode)(nil)
