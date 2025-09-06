package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

// execNFTCLI runs the CLI with args capturing stdout.
func execNFTCLI(args ...string) (string, error) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetErr(new(bytes.Buffer))
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	w.Close()
	os.Stdout = stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

// TestNFTMarketplaceCLI exercises minting, listing and buying NFTs.
func TestNFTMarketplaceCLI(t *testing.T) {
	nftMarket = core.NewNFTMarketplace()

	out, err := execNFTCLI("nft", "mint", "id1", "alice", "meta", "100")
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if !strings.Contains(out, "minted") {
		t.Fatalf("unexpected output: %s", out)
	}

	out, err = execNFTCLI("--json", "nft", "list", "id1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "id1") {
		t.Fatalf("expected id in output: %s", out)
	}

	out, err = execNFTCLI("nft", "buy", "id1", "bob")
	if err != nil {
		t.Fatalf("buy: %v", err)
	}
	if !strings.Contains(out, "transferred") {
		t.Fatalf("unexpected output: %s", out)
	}
}

// TestNFTMarketplaceInvalidPrice ensures zero price is rejected.
func TestNFTMarketplaceInvalidPrice(t *testing.T) {
	nftMarket = core.NewNFTMarketplace()
	if _, err := execNFTCLI("nft", "mint", "id2", "alice", "meta", "0"); err == nil {
		t.Fatalf("expected error for invalid price")
	}
}
