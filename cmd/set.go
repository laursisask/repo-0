package cmd

import (
	"fmt"
	"os"

	"github.com/Keeper-Security/linux-keyring-utility/pkg/keyring"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Args:  cobra.MinimumNArgs(2),
	Short: "Set a secret in the keyring.",
	Long:  `Add a KSM Config to the keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		secret := args[1]
		provider := keyring.SecretProvider{}

		err := provider.Set(appName, secret)
		if err != nil {
			fmt.Println("Unable to set secret:", err)
			os.Exit(1)
		}
	},
}
