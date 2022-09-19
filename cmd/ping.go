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

func runPingCommand(cmd *cobra.Command, subject string) bool {
	timeout, _ := cmd.Flags().GetDuration("timeout")
	interval, _ := cmd.Flags().GetDuration("interval")
	var status bool

	pinger, _ := ping.NewPinger(subject)
	pinger.Timeout = timeout
	pinger.Interval = interval
	pinger.Debug = true
	pinger.OnFinish = func(stats *ping.Statistics) {
		if stats.PacketsRecv == 0 {
			status = false
		}
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		pinger.Stop()
		status = true
	}

	pinger.Run()
	return status
}

var pingCommand = &cobra.Command{
	Use:   "ping [subject]",
	Short: "Ping the subject",
	Long:  "Send ICMP packet to a IPv4, IPv6 or a hostname subject",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connecting to the subject...")
		isAvailable := runPingCommand(cmd, args[0])
		if !isAvailable {
			fmt.Println("timeout - subject is not available.")
			return
		}
		fmt.Println("subject available.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateArguments(args)
	},
}

func init() {
	rootCmd.AddCommand(pingCommand)
}
