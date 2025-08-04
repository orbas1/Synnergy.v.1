package core

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Address string

func (a Address) Hex() string { return string(a) }

func (a Address) String() string { return string(a) }

func (a Address) Bytes() []byte { return []byte(a) }

func (a Address) Short() string {
	s := string(a)
	if len(s) <= 10 {
		return s
	}
	return s[:6] + "..." + s[len(s)-4:]
}

func StringToAddress(s string) (Address, error) {
	s = strings.ToLower(s)
	if len(s) != 42 || !strings.HasPrefix(s, "0x") {
		return "", fmt.Errorf("invalid address")
	}
	if _, err := hex.DecodeString(s[2:]); err != nil {
		return "", fmt.Errorf("invalid address")
	}
	return Address(s), nil
}

type Hash [32]byte
