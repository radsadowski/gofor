package cmd

import (
	"fmt"
	"goforit/internal/validation"

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
	if !validation.IsValidIPAddress(subject) && !validation.IsValidHostname(subject) {
		return &InvalidSubjectError{subject}
	}

	return nil
}

func runPingCommand(cmd *cobra.Command, host string) bool {
	timeout, _ := cmd.Flags().GetDuration("timeout")
	interval, _ := cmd.Flags().GetDuration("interval")
	var status bool

	pinger, _ := ping.NewPinger(host)
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
	Use:   "ping host",
	Short: "Ping the host",
	Long:  "Send ICMP packet to a IPv4, IPv6 or a hostname subject",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Connecting to the subject...")
		isAvailable := runPingCommand(cmd, args[0])
		if !isAvailable {
			fmt.Println("Timeout - host is not available.")
			return
		}
		fmt.Println("host is available.")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateArguments(args)
	},
	Example: "- gofor ping 192.168.1.1 --timeout=10s --interval=2s\n- gofor ping google.com -t 5s -i 2s",
}

func init() {
	rootCmd.AddCommand(pingCommand)
}
