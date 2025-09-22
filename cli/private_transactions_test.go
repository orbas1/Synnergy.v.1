package cli

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"synnergy/core"
)

func TestPrivateTransactionsCommands(t *testing.T) {
	privTxMgr = core.NewPrivateTxManager()
	t.Cleanup(func() { privTxMgr = core.NewPrivateTxManager() })

	key := "0123456789abcdef0123456789abcdef"
	plaintext := "secret"
	out, err := executeCLICommand(t, "private-tx", "encrypt", key, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	cipher, err := hex.DecodeString(out)
	if err != nil {
		t.Fatalf("decode cipher: %v", err)
	}
	if len(cipher) <= len(plaintext) {
		t.Fatalf("ciphertext too short: %d", len(cipher))
	}

	generated, err := core.Encrypt([]byte(key), []byte(plaintext))
	if err != nil {
		t.Fatalf("encrypt fixture: %v", err)
	}
	decoded, err := executeCLICommand(t, "private-tx", "decrypt", key, hex.EncodeToString(generated))
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if decoded != plaintext {
		t.Fatalf("expected plaintext %q, got %q", plaintext, decoded)
	}

	payload := hex.EncodeToString([]byte{0x01, 0x02})
	nonce := hex.EncodeToString([]byte{0x03, 0x04})
	body, err := json.Marshal(map[string]string{"payload": payload, "nonce": nonce})
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "tx.json")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := executeCLICommand(t, "private-tx", "send", path); err != nil {
		t.Fatalf("send: %v", err)
	}
	txs := privTxMgr.List()
	if len(txs) != 1 {
		t.Fatalf("expected 1 tx, got %d", len(txs))
	}
	if hex.EncodeToString(txs[0].Payload) != payload || hex.EncodeToString(txs[0].Nonce) != nonce {
		t.Fatalf("unexpected stored tx: %+v", txs[0])
	}
}
