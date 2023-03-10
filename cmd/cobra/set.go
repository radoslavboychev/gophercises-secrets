package cobra

import (
	"fmt"

	"github.com/radoslavboychev/gophercises-secret/secret"
	"github.com/spf13/cobra"
)

// define the 'set' command
// creates a new Vault,
// invokes the Set() method for the vault
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("Value set successfully!")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
