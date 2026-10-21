package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initIaConfig)

	iaCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	iaCmd.AddCommand(iaAddCmd)
	iaCmd.AddCommand(iaRemoveCmd)
}

var iaCmd = &cobra.Command{
	Use:   "ia",
	Short: "Internet Archive sub-commands",
	Long:  `...`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func initIaConfig() {
	if Verbose {
		SlogSetup(slog.LevelDebug)
	} else {
		SlogSetup(slog.LevelInfo)
	}
}
