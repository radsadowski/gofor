package cmd

import (
	"fmt"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "icmp",
	Short: "Ping the subject",
	Long:  "Send ICMP packet to a IPv4, IPv6 or a hostname subject",
	Run: func(cmd *cobra.Command, args []string) {
		subject := args[0]
		pinger, err := ping.NewPinger(subject)
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		err = pinger.Run() // Blocks until finished.
		if err != nil {
			panic(err)
		}
		fmt.Println(pinger.Statistics())
	},
	Args: cobra.ExactArgs(1),
}
