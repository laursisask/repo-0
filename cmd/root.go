package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)
  

var rootCmd = &cobra.Command{
	Use:   "lku",
	Short: "KSM Keyring Utility",
	Long: `Keyring Utility for for saving and retrieving KSM configs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}
  
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
}

