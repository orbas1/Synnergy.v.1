package synnergy

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// TrainingJob represents a model training task.
type TrainingJob struct {
	ID          string
	DatasetCID  string
	ModelCID    string
	Status      string
	StartedAt   time.Time
	CompletedAt time.Time
}

// TrainingManager orchestrates training jobs.
type TrainingManager struct {
	mu     sync.RWMutex
	jobs   map[string]TrainingJob
	nextID uint64
}

// NewTrainingManager creates a new TrainingManager.
func NewTrainingManager() *TrainingManager {
	return &TrainingManager{jobs: make(map[string]TrainingJob)}
}

// Start begins a new training job and returns its ID.
// Both datasetCID and modelCID must be provided or an error is returned.
func (m *TrainingManager) Start(datasetCID, modelCID string) (string, error) {
	if datasetCID == "" || modelCID == "" {
		return "", errors.New("datasetCID and modelCID required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := fmt.Sprintf("job-%d", m.nextID)
	m.jobs[id] = TrainingJob{
		ID:         id,
		DatasetCID: datasetCID,
		ModelCID:   modelCID,
		Status:     "running",
		StartedAt:  time.Now().UTC(),
	}
	return id, nil
}

// Status returns a training job by ID.
func (m *TrainingManager) Status(id string) (TrainingJob, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[id]
	return job, ok
}

// List returns all jobs.
func (m *TrainingManager) List() []TrainingJob {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]TrainingJob, 0, len(m.jobs))
	for _, j := range m.jobs {
		out = append(out, j)
	}
	return out
}

// Cancel marks a running job as cancelled.
func (m *TrainingManager) Cancel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	job, ok := m.jobs[id]
	if !ok {
		return errors.New("job not found")
	}
	if job.Status != "running" {
		return errors.New("job not running")
	}
	job.Status = "cancelled"
	job.CompletedAt = time.Now().UTC()
	m.jobs[id] = job
	return nil
}

// Complete marks a running job as finished successfully.
func (m *TrainingManager) Complete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	job, ok := m.jobs[id]
	if !ok {
		return errors.New("job not found")
	}
	if job.Status != "running" {
		return errors.New("job not running")
	}
	job.Status = "completed"
	job.CompletedAt = time.Now().UTC()
	m.jobs[id] = job
	return nil
}
