package core

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestComplianceServiceKYCAndRisk(t *testing.T) {
	svc := NewComplianceService()
	if _, err := svc.ValidateKYC("addr1", []byte("kycdata")); err != nil {
		t.Fatalf("validate kyc: %v", err)
	}
	svc.RecordFraud("addr1", 5)
	if score := svc.RiskScore("addr1"); score != 5 {
		t.Fatalf("expected risk score 5 got %d", score)
	}
	svc.EraseKYC("addr1")
	if len(svc.AuditTrail("addr1")) == 0 {
		t.Fatalf("expected audit trail")
	}
}

func TestComplianceServiceMonitorTransaction(t *testing.T) {
	svc := NewComplianceService()
	tx := ComplianceTransaction{ID: "tx1", From: "a", Amount: 100}
	if !svc.MonitorTransaction(tx, 50) {
		t.Fatalf("expected anomaly detection")
	}
}

func TestComplianceServiceVerifyZKP(t *testing.T) {
	svc := NewComplianceService()
	blob := []byte("data")
	h := sha256.Sum256(blob)
	if !svc.VerifyZKP(blob, hex.EncodeToString(h[:]), "unused") {
		t.Fatalf("expected verification success")
	}
}
