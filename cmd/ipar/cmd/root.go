package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ipar",
	Short: "The InterPlanetary Archive CLI",
	Long:  `...`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var (
	Verbose bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initRootConfig)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.AddCommand(iaCmd)
	rootCmd.AddCommand(addCmd)
}

func initRootConfig() {
	if Verbose {
		SlogSetup(slog.LevelDebug)
	} else {
		SlogSetup(slog.LevelInfo)
	}
}

func printAndDie(args ...any) {
	fmt.Println(args...)
	os.Exit(1)
}

func SlogSetup(lvl slog.Level) {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	)
}
