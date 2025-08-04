package nodes

import "time"

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
