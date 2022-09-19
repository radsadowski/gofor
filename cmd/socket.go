package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var socketCommand = &cobra.Command{
	Use:   "socket subject:port",
	Short: "Check if socket is alive",
	Long:  "Send TCP or UDP segment to a socket subject",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("to do...")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateArguments(args)
	},
}

func init() {
	socketCommand.Flags().StringP("protocol", "p", "tcp", "specify protocol to use")
	rootCmd.AddCommand(socketCommand)
}
