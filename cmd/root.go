package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "orion-cli",
	Short:         "AI destekli görevleri GitHub Issue'larına dönüştürür",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() error {
	return rootCmd.Execute()
}
