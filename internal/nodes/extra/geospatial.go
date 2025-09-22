package nodes

import "time"

// GeoRecord represents a single geospatial data point captured by the network.
type GeoRecord struct {
	Subject   string    // entity the record relates to
	Latitude  float64   // decimal degrees
	Longitude float64   // decimal degrees
	Timestamp time.Time // time the reading was taken
}

// HistoryQuery filters returned geospatial history.
type HistoryQuery struct {
	Since time.Time
	Until time.Time
}

// GeoBounds describes the min/max coordinates recorded for a subject.
type GeoBounds struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

// GeoSummary captures aggregate statistics for a subject.
type GeoSummary struct {
	Subject string
	Count   int
	Bounds  GeoBounds
	Last    GeoRecord
}

// GeospatialOption configures behaviour for specialised implementations.
type GeospatialOption interface {
	applyGeospatialOption(*GeospatialConfig)
}

// GeospatialConfig mirrors the configuration used by runtime nodes. It is kept
// separate to avoid import cycles when the extra package is consumed by CLI or
// tooling packages.
type GeospatialConfig struct {
	MaxHistory      int
	RetentionWindow time.Duration
}

type geospatialOptionFunc func(*GeospatialConfig)

func (f geospatialOptionFunc) applyGeospatialOption(cfg *GeospatialConfig) { f(cfg) }

// WithMaxHistory configures the maximum number of entries retained per subject.
func WithMaxHistory(limit int) GeospatialOption {
	return geospatialOptionFunc(func(cfg *GeospatialConfig) {
		if limit > 0 {
			cfg.MaxHistory = limit
		}
	})
}

// WithRetentionWindow configures the rolling retention window for entries.
func WithRetentionWindow(window time.Duration) GeospatialOption {
	return geospatialOptionFunc(func(cfg *GeospatialConfig) {
		if window > 0 {
			cfg.RetentionWindow = window
		}
	})
}

// ApplyOptions helper is exported for tooling packages that need to build
// configuration structs without depending on the runtime implementation.
func ApplyOptions(base GeospatialConfig, opts ...GeospatialOption) GeospatialConfig {
	cfg := base
	for _, opt := range opts {
		if opt != nil {
			opt.applyGeospatialOption(&cfg)
		}
	}
	return cfg
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
