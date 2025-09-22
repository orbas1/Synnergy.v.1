package nodes

import (
	"testing"
	"time"
)

func TestApplyOptions(t *testing.T) {
	base := GeospatialConfig{}
	cfg := ApplyOptions(base, WithMaxHistory(5), WithRetentionWindow(time.Minute))
	if cfg.MaxHistory != 5 {
		t.Fatalf("expected max history 5 got %d", cfg.MaxHistory)
	}
	if cfg.RetentionWindow != time.Minute {
		t.Fatalf("expected retention %v got %v", time.Minute, cfg.RetentionWindow)
	}
}

type noopOption struct{}

func (noopOption) applyGeospatialOption(cfg *GeospatialConfig) { cfg.MaxHistory = 42 }

func TestApplyOptionsIgnoresNil(t *testing.T) {
	cfg := ApplyOptions(GeospatialConfig{}, nil, noopOption{})
	if cfg.MaxHistory != 42 {
		t.Fatalf("unexpected config %+v", cfg)
	}
}
