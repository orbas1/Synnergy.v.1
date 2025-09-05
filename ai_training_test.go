package synnergy

import "testing"

func TestTrainingLifecycle(t *testing.T) {
	m := NewTrainingManager()
	if _, err := m.Start("", "model"); err == nil {
		t.Fatalf("expected error for empty dataset")
	}
	id, err := m.Start("dataset", "model")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if job, ok := m.Status(id); !ok || job.Status != "running" {
		t.Fatalf("unexpected job status: %#v %v", job, ok)
	}
	if err := m.Complete(id); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if job, ok := m.Status(id); !ok || job.Status != "completed" {
		t.Fatalf("expected completed job: %#v %v", job, ok)
	}
	if err := m.Cancel(id); err == nil {
		t.Fatalf("cancel should fail on completed job")
	}
	if len(m.List()) != 1 {
		t.Fatalf("expected one job in list")
	}
}
