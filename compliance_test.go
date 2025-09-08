package synnergy

import "testing"

func TestComplianceServiceKYCAndRisk(t *testing.T) {
	svc := NewComplianceService()
	if _, err := svc.ValidateKYC("addr1", []byte("kycdata")); err != nil {
		t.Fatalf("validate kyc: %v", err)
	}
	if err := svc.RecordFraud("addr1", 5); err != nil {
		t.Fatalf("record fraud: %v", err)
	}
	if score := svc.RiskScore("addr1"); score != 5 {
		t.Fatalf("expected risk score 5 got %d", score)
	}
	if err := svc.EraseKYC("addr1"); err != nil {
		t.Fatalf("erase kyc: %v", err)
	}
	if len(svc.AuditTrail("addr1")) == 0 {
		t.Fatalf("expected audit trail")
	}

	if _, err := svc.ValidateKYC("", []byte("kycdata")); err == nil {
		t.Fatalf("expected error for empty address")
	}
	if err := svc.RecordFraud("", 1); err == nil {
		t.Fatalf("expected error for empty fraud address")
	}
	if err := svc.RecordFraud("addr1", 0); err == nil {
		t.Fatalf("expected error for non-positive severity")
	}
	if err := svc.EraseKYC("missing"); err == nil {
		t.Fatalf("expected error for missing kyc record")
	}
}

func TestComplianceServiceMonitorTransaction(t *testing.T) {
	svc := NewComplianceService()
	tx := ComplianceTransaction{ID: "tx1", From: "a", Amount: 100}
	if !svc.MonitorTransaction(tx, 50) {
		t.Fatalf("expected anomaly detection")
	}
}
