package cmd

import (
	"fmt"
	iping "goforit/internal/ping"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

type InvalidSubjectError struct {
	detail string
}

func (invalidSubjectError *InvalidSubjectError) Error() string {
	return fmt.Sprintf("Invalid subject specified. Must be an IPv4, IPv6 or hostname. '%s' given.", invalidSubjectError.detail)
}

type NoArgumentError struct {
}

func (noArgumentError *NoArgumentError) Error() string {
	return "Requires at least one argument."
}

func validateArguments(args []string) error {
	if len(args) < 1 {
		return &NoArgumentError{}
	}
	subject := args[0]
	if !iping.IsValidIPAddress(subject) && !iping.IsValidHostname(subject) {
		return &InvalidSubjectError{subject}
	}

	return nil
}

func runPingCommand(cmd *cobra.Command, subject string) {
	timeout, _ := cmd.Flags().GetDuration("timeout")
	interval, _ := cmd.Flags().GetDuration("interval")
	// try to ping the subject until the timeout is reached
	// wait the interval value between next ping
	// c := make(chan bool)

	// go func() {
	fmt.Println("connecting to peer...")
	pinger, _ := ping.NewPinger(subject)
	pinger.Count = 1
	pinger.OnFinish = func(stats *ping.Statistics) {
		if stats.PacketsRecv == 0 {
			fmt.Println("failed to connect.")
			return
		}
		fmt.Println("finished")
	}
	// pinger.OnRecv = func(pkt *ping.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	// }
	pinger.Timeout = timeout
	pinger.Interval = interval
	pinger.Debug = true
	err := pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
}

var pingCommand = &cobra.Command{
	Use:   "ping",
	Short: "Ping the subject",
	Long:  "Send ICMP packet to a IPv4, IPv6 or a hostname subject",
	Run: func(cmd *cobra.Command, args []string) {
		runPingCommand(cmd, args[0])
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateArguments(args)
	},
}

func init() {
	rootCmd.AddCommand(pingCommand)
}
