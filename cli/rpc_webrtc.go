package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	webrtcRPC = core.NewWebRTCRPC()
	peerChans = map[string]<-chan []byte{}
)

func init() {
	cmd := &cobra.Command{
		Use:   "rpcwebrtc",
		Short: "Simulate WebRTC RPC messaging",
	}

	connectCmd := &cobra.Command{
		Use:   "connect [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Connect a peer",
		Run: func(cmd *cobra.Command, args []string) {
			ch := webrtcRPC.Connect(args[0])
			peerChans[args[0]] = ch
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send [id] [msg]",
		Args:  cobra.ExactArgs(2),
		Short: "Send message to peer",
		Run: func(cmd *cobra.Command, args []string) {
			if webrtcRPC.Send(args[0], []byte(args[1])) {
				fmt.Println("sent")
			} else {
				fmt.Println("failed")
			}
		},
	}

	recvCmd := &cobra.Command{
		Use:   "recv [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Receive a message for peer",
		Run: func(cmd *cobra.Command, args []string) {
			ch, ok := peerChans[args[0]]
			if !ok {
				fmt.Println("not connected")
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
		Use:   "disconnect [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Disconnect a peer",
		Run: func(cmd *cobra.Command, args []string) {
			webrtcRPC.Disconnect(args[0])
			delete(peerChans, args[0])
		},
	}

	cmd.AddCommand(connectCmd, sendCmd, recvCmd, disconnectCmd)
	rootCmd.AddCommand(cmd)
}
