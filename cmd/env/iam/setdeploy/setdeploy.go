package setdeploy

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "setdeploy",
	Short: "Set Apigee Deployer role for a SA on an environment",
	Long:  "Set Apigee Deployer role for a SA on an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return shared.SetIAMServiceAccount(serviceAccountName, "deploy")
	},
}

var serviceAccountName string

func init() {

	Cmd.Flags().StringVarP(&serviceAccountName, "name", "n",
		"", "Service Account Name")

	_ = Cmd.MarkFlagRequired("name")
}
