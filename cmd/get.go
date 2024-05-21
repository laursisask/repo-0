package cmd

import (
	"fmt"
	"os"

	"github.com/Keeper-Security/linux-keyring-utility/pkg/keyring"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get a secret from the keyring.",
	Long:  `Set a KSM Config from the keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		provider := keyring.SecretProvider{}

		secret, err := provider.Get(appName)
		if err != nil {
			fmt.Println("Unable to get secret:", err)
			os.Exit(1)
		}

		fmt.Println(secret)
	},
}
