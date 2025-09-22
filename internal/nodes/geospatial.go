package nodes

import (
	"errors"
	"math"
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
	// HistoryWithin constrains returned records to a time window.
	HistoryWithin(subject string, query HistoryQuery) []GeoRecord
	// Latest returns the most recent record for the subject.
	Latest(subject string) (GeoRecord, bool)
	// Summary returns aggregated metrics for the subject.
	Summary(subject string) (GeoSummary, bool)
}

// HistoryQuery filters returned geospatial history.
type HistoryQuery struct {
	Since time.Time
	Until time.Time
}

// GeoSummary captures aggregate statistics for a subject.
type GeoSummary struct {
	Subject string
	Count   int
	Bounds  GeoBounds
	Last    GeoRecord
}

// GeoBounds describes the min/max coordinates recorded for a subject.
type GeoBounds struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

type geospatialConfig struct {
	maxHistory      int
	retentionWindow time.Duration
}

// GeospatialOption configures optional behaviour.
type GeospatialOption func(*geospatialConfig)

// WithGeospatialMaxHistory limits the number of history entries retained per
// subject. Non-positive values retain the default unlimited behaviour.
func WithGeospatialMaxHistory(limit int) GeospatialOption {
	return func(cfg *geospatialConfig) {
		if limit > 0 {
			cfg.maxHistory = limit
		}
	}
}

// WithGeospatialRetention prunes history entries older than the supplied
// duration during subsequent writes. A zero duration disables pruning.
func WithGeospatialRetention(window time.Duration) GeospatialOption {
	return func(cfg *geospatialConfig) {
		if window > 0 {
			cfg.retentionWindow = window
		}
	}
}

// GeospatialNode provides an in-memory geospatial tracker backed by a
// BasicNode. Records are stored in-memory and suitable for testing and
// development scenarios.
type GeospatialNode struct {
	*BasicNode
	mu     sync.RWMutex
	points map[string][]GeoRecord
	cfg    geospatialConfig
}

// NewGeospatialNode creates a new geospatial node.
func NewGeospatialNode(id Address, opts ...GeospatialOption) *GeospatialNode {
	cfg := geospatialConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &GeospatialNode{
		BasicNode: NewBasicNode(id),
		points:    make(map[string][]GeoRecord),
		cfg:       cfg,
	}
}

// Record captures a latitude and longitude for the given subject.
func (n *GeospatialNode) Record(subject string, lat, lon float64) error {
	if subject == "" {
		return errors.New("subject is required")
	}
	if math.Abs(lat) > 90 {
		return errors.New("latitude out of range")
	}
	if math.Abs(lon) > 180 {
		return errors.New("longitude out of range")
	}

	rec := GeoRecord{Subject: subject, Latitude: lat, Longitude: lon, Timestamp: time.Now().UTC()}

	n.mu.Lock()
	entries := append(n.points[subject], rec)
	if n.cfg.retentionWindow > 0 {
		cutoff := rec.Timestamp.Add(-n.cfg.retentionWindow)
		filtered := entries[:0]
		for _, existing := range entries {
			if existing.Timestamp.After(cutoff) || existing.Timestamp.Equal(cutoff) {
				filtered = append(filtered, existing)
			}
		}
		entries = append([]GeoRecord(nil), filtered...)
	}
	if n.cfg.maxHistory > 0 && len(entries) > n.cfg.maxHistory {
		entries = append([]GeoRecord(nil), entries[len(entries)-n.cfg.maxHistory:]...)
	}
	n.points[subject] = entries
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

// HistoryWithin returns history constrained to the provided window.
func (n *GeospatialNode) HistoryWithin(subject string, query HistoryQuery) []GeoRecord {
	history := n.History(subject)
	if len(history) == 0 {
		return history
	}
	out := history[:0]
	for _, rec := range history {
		if !query.Since.IsZero() && rec.Timestamp.Before(query.Since) {
			continue
		}
		if !query.Until.IsZero() && rec.Timestamp.After(query.Until) {
			continue
		}
		out = append(out, rec)
	}
	return append([]GeoRecord(nil), out...)
}

// Latest returns the most recent record for the subject.
func (n *GeospatialNode) Latest(subject string) (GeoRecord, bool) {
	n.mu.RLock()
	records := n.points[subject]
	n.mu.RUnlock()
	if len(records) == 0 {
		return GeoRecord{}, false
	}
	latest := records[0]
	for _, rec := range records[1:] {
		if rec.Timestamp.After(latest.Timestamp) {
			latest = rec
		}
	}
	return latest, true
}

// Summary returns aggregate metrics for the subject.
func (n *GeospatialNode) Summary(subject string) (GeoSummary, bool) {
	records := n.History(subject)
	if len(records) == 0 {
		return GeoSummary{}, false
	}
	bounds := GeoBounds{
		MinLat: records[0].Latitude,
		MaxLat: records[0].Latitude,
		MinLon: records[0].Longitude,
		MaxLon: records[0].Longitude,
	}
	for _, rec := range records[1:] {
		if rec.Latitude < bounds.MinLat {
			bounds.MinLat = rec.Latitude
		}
		if rec.Latitude > bounds.MaxLat {
			bounds.MaxLat = rec.Latitude
		}
		if rec.Longitude < bounds.MinLon {
			bounds.MinLon = rec.Longitude
		}
		if rec.Longitude > bounds.MaxLon {
			bounds.MaxLon = rec.Longitude
		}
	}
	last := records[len(records)-1]
	return GeoSummary{Subject: subject, Count: len(records), Bounds: bounds, Last: last}, true
}

// Ensure GeospatialNode implements GeospatialNodeInterface.
var _ GeospatialNodeInterface = (*GeospatialNode)(nil)
