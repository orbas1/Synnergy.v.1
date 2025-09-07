package main

import (
    "encoding/hex"
    "flag"
    "fmt"
    "io"
    "os"
    "strings"

    "synnergy/internal/p2p"
)

func runWithManager(mgr *p2p.Manager, args []string, out io.Writer) int {
    if len(args) == 0 {
        fmt.Fprintln(out, "expected subcommand add-peer|list-peers")
        return 1
    }
    switch args[0] {
    case "add-peer":
        fs := flag.NewFlagSet("add-peer", flag.ContinueOnError)
        id := fs.String("id", "", "peer id")
        addr := fs.String("addr", "", "peer address")
        pub := fs.String("pubkey", "", "peer public key (hex)")
        fs.SetOutput(out)
        if err := fs.Parse(args[1:]); err != nil {
            return 1
        }
        if *id == "" || *addr == "" {
            fmt.Fprintln(out, "id and addr required")
            return 1
        }
        peer := p2p.Peer{ID: *id, Address: *addr}
        if *pub != "" {
            b, err := hex.DecodeString(strings.TrimSpace(*pub))
            if err != nil {
                fmt.Fprintln(out, "invalid pubkey")
                return 1
            }
            peer.PubKey = b
        }
        mgr.AddPeer(peer)
        fmt.Fprintln(out, "peer added")
        return 0
    case "list-peers":
        peers := mgr.ListPeers()
        for _, p := range peers {
            fmt.Fprintf(out, "%s %s\n", p.ID, p.Address)
        }
        return 0
    default:
        fmt.Fprintln(out, "unknown subcommand")
        return 1
    }
}

var defaultManager = p2p.NewManager()

func run(args []string) int {
    return runWithManager(defaultManager, args, os.Stdout)
}

func main() {
    os.Exit(run(os.Args[1:]))
}
