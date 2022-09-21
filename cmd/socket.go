package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

func runSocketCommand(cmd *cobra.Command, host string) {
	var successChannel = make(chan bool)
	timeout, _ := cmd.Flags().GetDuration("timeout")
	interval, _ := cmd.Flags().GetDuration("interval")

	go func() {
		fmt.Printf("Trying to connect to %s\n", host)
		waited := interval / time.Second
		for {
			conn, err := net.Dial("tcp", host)
			if err == nil {
				conn.Close()
				successChannel <- true
			}
			time.Sleep(interval)
			fmt.Printf("\rWaited: %s", waited*time.Second)
			waited += interval / time.Second
		}
	}()

	select {
	case <-successChannel:
		fmt.Println("Host is available.")
	case <-time.After(timeout):
		fmt.Println("\nTimeout - host is not available.")
	}
}

var socketCommand = &cobra.Command{
	Use:   "socket host:port",
	Short: "Connect to the host",
	Long:  "Send TCP or UDP segment to a host",
	Run: func(cmd *cobra.Command, args []string) {
		runSocketCommand(cmd, args[0])
	},
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	return validateArguments(args)
	// },
}

func init() {
	socketCommand.Flags().StringP("protocol", "p", "tcp", "specify protocol to use")
	rootCmd.AddCommand(socketCommand)
}
