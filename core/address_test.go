package core

import (
    "encoding/hex"
    "strings"
    "testing"
)

func TestStringToAddress(t *testing.T) {
    cases := []struct {
        name    string
        input   string
        want    Address
        wantErr bool
    }{
        {
            name:  "with0x",
            input: "0x0123456789abcdef0123456789abcdef01234567",
            want:  Address("0x0123456789abcdef0123456789abcdef01234567"),
        },
        {
            name:  "without0x",
            input: "0123456789abcdef0123456789abcdef01234567",
            want:  Address("0x0123456789abcdef0123456789abcdef01234567"),
        },
        {
            name:  "uppercase",
            input: "0xABCDEF0123456789ABCDEF0123456789ABCDEF01",
            want:  Address("0xabcdef0123456789abcdef0123456789abcdef01"),
        },
        {
            name:    "short",
            input:   "0x1234",
            wantErr: true,
        },
        {
            name:    "invalidChars",
            input:   "0x0123456789abcdef0123456789abcdef0123456g",
            wantErr: true,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := StringToAddress(tc.input)
            if tc.wantErr {
                if err == nil {
                    t.Fatalf("expected error for %s", tc.input)
                }
                return
            }
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tc.want {
                t.Fatalf("got %s want %s", got, tc.want)
            }
        })
    }
}

func TestAddressHexBytesAndString(t *testing.T) {
    addrStr := "0x0123456789abcdef0123456789abcdef01234567"
    addr, err := StringToAddress(addrStr)
    if err != nil {
        t.Fatalf("StringToAddress: %v", err)
    }
    if addr.Hex() != addrStr {
        t.Fatalf("Hex mismatch: %s", addr.Hex())
    }
    if addr.String() != addrStr {
        t.Fatalf("String mismatch: %s", addr.String())
    }
    b := addr.Bytes()
    if len(b) != 20 {
        t.Fatalf("expected 20 bytes, got %d", len(b))
    }
    if hex.EncodeToString(b) != strings.TrimPrefix(addrStr, "0x") {
        t.Fatalf("Bytes mismatch: %x", b)
    }
}

func TestAddressBytesInvalid(t *testing.T) {
    addr := Address("0xzz")
    if len(addr.Bytes()) != 0 {
        t.Fatalf("expected empty slice for invalid address")
    }
}

func TestAddressShort(t *testing.T) {
    addr, err := StringToAddress("0x0123456789abcdef0123456789abcdef01234567")
    if err != nil {
        t.Fatalf("StringToAddress: %v", err)
    }
    if s := addr.Short(); s != "0x0123...4567" {
        t.Fatalf("unexpected short form: %s", s)
    }

    shortAddr := Address("0x1234")
    if s := shortAddr.Short(); s != "0x1234" {
        t.Fatalf("short address should remain unchanged, got %s", s)
    }
}

