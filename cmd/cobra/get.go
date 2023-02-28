package cobra

import (
	"fmt"

	"github.com/radoslavboychev/gophercises-secret/secret"
	"github.com/spf13/cobra"
)

// define the 'get' console command
// invokes the Get() command for the vault
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value set")
			return
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
