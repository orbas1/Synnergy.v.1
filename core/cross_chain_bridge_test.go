package core

import (
        "encoding/base64"
        "encoding/json"
        "testing"
        "time"
)

func TestBridgeManager(t *testing.T) {
	l := NewLedger()
	l.Credit("alice", 100)
	bm := NewBridgeManager(l)
	bridgeID := bm.RegisterBridge("chainA", "chainB", "relayer1")
	if bridgeID == 0 {
		t.Fatalf("expected bridge id")
	}
	if !bm.IsRelayerAuthorized(bridgeID, "relayer1") {
		t.Fatalf("relayer1 should be authorized")
	}
	if bm.IsRelayerAuthorized(bridgeID, "relayer2") {
		t.Fatalf("relayer2 should not be authorized")
	}

	transferID, err := bm.Deposit(bridgeID, "alice", "bob", 50, "token")
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	if l.GetBalance("alice") != 50 {
		t.Fatalf("expected alice balance 50")
	}
        proofPayload := bridgeClaimProof{
                TransferID: transferID,
                BridgeID:   bridgeID,
                Recipient:  "bob",
                Amount:     50,
                TokenID:    "token",
                SourceTx:   "tx-hash",
                Signers:    []string{"relayer1"},
                Timestamp:  time.Now().Unix(),
        }
        badPayload := proofPayload
        badPayload.Checksum = "invalid"
        badData, _ := json.Marshal(badPayload)
        if err := bm.Claim(transferID, "relayer1", base64.StdEncoding.EncodeToString(badData)); err == nil {
                t.Fatalf("expected checksum validation to fail")
        }

        proofPayload.Checksum = proofPayload.expectedChecksum(bm.transfers[transferID])
        proofData, _ := json.Marshal(proofPayload)
        proof := base64.StdEncoding.EncodeToString(proofData)

        if err := bm.Claim(transferID, "relayer2", proof); err == nil {
                t.Fatalf("expected relayer authorization error")
        }
        if err := bm.Claim(transferID, "relayer1", proof); err != nil {
                t.Fatalf("claim failed: %v", err)
        }
        if l.GetBalance("bob") != 50 {
                t.Fatalf("expected bob balance 50")
        }
        record, err := bm.GetTransfer(transferID)
        if err != nil {
                t.Fatalf("get transfer failed: %v", err)
        }
        if !record.Claimed || record.ProofChecksum == "" || record.RelayTx != proofPayload.SourceTx {
                t.Fatalf("expected transfer metadata to be recorded")
        }
        if len(record.Signers) != 1 || record.Signers[0] != "relayer1" {
                t.Fatalf("unexpected signers stored: %+v", record.Signers)
        }
        if len(bm.ListTransfers()) != 1 {
                t.Fatalf("unexpected transfer list length")
        }

	if err := bm.RemoveBridge(bridgeID); err != nil {
		t.Fatalf("remove bridge failed: %v", err)
	}
	if _, err := bm.GetBridge(bridgeID); err == nil {
		t.Fatalf("expected bridge not found after removal")
	}
	if _, err := bm.Deposit(bridgeID, "alice", "carol", 10, "token"); err == nil {
		t.Fatalf("expected deposit to fail for removed bridge")
	}
}

func TestBridgeTransferManagerClaimValidatesProof(t *testing.T) {
        manager := NewBridgeTransferManager()
        transfer, err := manager.Deposit("bridge-1", "alice", "bob", 10, "token")
        if err != nil {
                t.Fatalf("deposit failed: %v", err)
        }

        payload := bridgeTransferClaimProof{
                TransferID: transfer.ID,
                BridgeID:   transfer.BridgeID,
                Recipient:  transfer.To,
                Amount:     transfer.Amount,
                TokenID:    transfer.TokenID,
                SourceTx:   "external-tx",
                Signers:    []string{"observer"},
                Timestamp:  time.Now().Unix(),
        }

        // Expect validation failure when checksum missing
        badData, _ := json.Marshal(payload)
        if err := manager.Claim(transfer.ID, base64.StdEncoding.EncodeToString(badData)); err == nil {
                t.Fatalf("expected claim to fail for missing checksum")
        }

        payload.Checksum = payload.expectedChecksum(transfer)
        data, _ := json.Marshal(payload)
        if err := manager.Claim(transfer.ID, base64.StdEncoding.EncodeToString(data)); err != nil {
                t.Fatalf("claim should succeed: %v", err)
        }

        claimed, ok := manager.GetTransfer(transfer.ID)
        if !ok {
                t.Fatalf("transfer not found after claim")
        }
        if claimed.Status != "released" || claimed.ProofChecksum == "" || claimed.SourceTx != payload.SourceTx {
                t.Fatalf("expected transfer metadata to be recorded")
        }
        if len(claimed.Signers) != 1 || claimed.Signers[0] != "observer" {
                t.Fatalf("unexpected signers stored: %+v", claimed.Signers)
        }
        if claimed.ReleasedAt.IsZero() {
                t.Fatalf("expected release timestamp to be recorded")
        }
}
