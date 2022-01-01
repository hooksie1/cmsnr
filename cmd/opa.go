package cmd

import (
	"github.com/spf13/cobra"
)

// opaCmd represents the opa command
var opaCmd = &cobra.Command{
	Use:   "opa",
	Short: "OPA Controls the OPA commands for cmsnr",
	Long:  "opa allows cmsnr to interact with the OPA CRDs in the cluster",
	//Run: func(cmd *cobra.Command, args []string) {
	//	cmd.Help()
	//	os.Exit(1)
	//},
}

func init() {
	rootCmd.AddCommand(opaCmd)
}
