package synnergy

import "strings"

// SupportedContractLanguages enumerates the smart contract languages that the
// Synnergy virtual machine recognises. Contracts authored in these languages
// can be compiled to bytecode compatible with the VM's opcode set.
var SupportedContractLanguages = []string{
	"wasm",
	"golang",
	"javascript",
	"solidity",
	"rust",
	"python",
	"yul",
}

// IsLanguageSupported reports whether the supplied language is recognised as
// compatible with the VM. The check is case-insensitive.
func IsLanguageSupported(lang string) bool {
	lang = strings.ToLower(lang)
	for _, l := range SupportedContractLanguages {
		if l == lang {
			return true
		}
	}
	return false
}
