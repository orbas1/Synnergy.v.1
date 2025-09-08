package main

import (
	"fmt"
	"os"

	"synnergy/internal/security"
)

var sm = security.NewSecretsManager()

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: secrets-manager [set|get] ...")
		return
	}
	switch os.Args[1] {
	case "set":
		if len(os.Args) != 4 {
			fmt.Println("usage: secrets-manager set <key> <value>")
			return
		}
		if err := sm.Store(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println("ok")
	case "get":
		if len(os.Args) != 3 {
			fmt.Println("usage: secrets-manager get <key>")
			return
		}
		v, err := sm.Retrieve(os.Args[2])
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(v)
	default:
		fmt.Println("unknown command")
	}
}
