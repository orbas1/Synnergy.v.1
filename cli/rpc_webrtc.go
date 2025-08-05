package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	webrtcRPC   = core.NewWebRTCRPC()
	webrtcPeers = make(map[string]<-chan []byte)
)

func init() {
	rpcCmd := &cobra.Command{
		Use:   "rpc_webrtc",
		Short: "Simulate WebRTC-style RPC",
	}

	connectCmd := &cobra.Command{
		Use:   "connect <peerID>",
		Args:  cobra.ExactArgs(1),
		Short: "Connect a peer",
		Run: func(cmd *cobra.Command, args []string) {
			ch := webrtcRPC.Connect(args[0])
			webrtcPeers[args[0]] = ch
			fmt.Println("connected")
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send <peerID> <msg>",
		Args:  cobra.ExactArgs(2),
		Short: "Send a message to a peer",
		Run: func(cmd *cobra.Command, args []string) {
			if !webrtcRPC.Send(args[0], []byte(args[1])) {
				fmt.Println("peer not found")
			}
		},
	}

	recvCmd := &cobra.Command{
		Use:   "recv <peerID>",
		Args:  cobra.ExactArgs(1),
		Short: "Receive a message from a peer",
		Run: func(cmd *cobra.Command, args []string) {
			ch, ok := webrtcPeers[args[0]]
			if !ok {
				fmt.Println("peer not found")
				return
			}
			select {
			case msg := <-ch:
				fmt.Println(string(msg))
			default:
				fmt.Println("no message")
			}
		},
	}

	disconnectCmd := &cobra.Command{
		Use:   "disconnect <peerID>",
		Args:  cobra.ExactArgs(1),
		Short: "Disconnect a peer",
		Run: func(cmd *cobra.Command, args []string) {
			webrtcRPC.Disconnect(args[0])
			delete(webrtcPeers, args[0])
		},
	}

	rpcCmd.AddCommand(connectCmd, sendCmd, recvCmd, disconnectCmd)
	rootCmd.AddCommand(rpcCmd)
}
