package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initHfConfig)

	hfCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	hfCmd.AddCommand(hfAddCmd)
	hfCmd.AddCommand(hfRemoveCmd)
}

var hfCmd = &cobra.Command{
	Use:   "hf",
	Short: "HuggingFace sub-commands",
	Long:  `...`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func initHfConfig() {
	if Verbose {
		SlogSetup(slog.LevelDebug)
	} else {
		SlogSetup(slog.LevelInfo)
	}
}
