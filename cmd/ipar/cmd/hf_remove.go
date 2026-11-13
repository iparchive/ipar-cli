package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chigopher/pathlib"
	"github.com/go-viper/mapstructure/v2"
	"github.com/iparchive/ipar-cli"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hfRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a Huggingface model from your IPAR",
	Long:  `...`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var config struct {
			ArchivalDir *pathlib.Path                     `mapstructure:"archival_dir"`
			Identifiers []ipar.HuggingfaceModelIdentifier `mapstructure:"model_ids"`
		}

		vi := viper.NewWithOptions(
			viper.WithDecodeHook(
				mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					ipar.StringToPath())))
		vi.SetEnvPrefix("ipar")
		vi.AutomaticEnv()
		vi.SetDefault("archival_dir", func() string {
			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				printAndDie(err)
			}
			return pathlib.NewPath(userHomeDir).Join("ipar", "huggingface").String()
		}())
		vi.SetDefault("model_ids", []string{})
		vi.Set("model_ids", args)
		if err := vi.Unmarshal(&config); err != nil {
			printAndDie(err)
		}

		ok, err := config.ArchivalDir.Exists()
		if err != nil {
			printAndDie(err)
		}
		if !ok {
			if err := config.ArchivalDir.MkdirAll(); err != nil {
				printAndDie(err)
			}
		}

		for _, identifier := range config.Identifiers {
			modelArchivalDir := config.ArchivalDir.Join(identifier)

			coreAPI, err := rpc.NewLocalApi()
			if err != nil {
				printAndDie(err)
			}

			ipfs := ipar.NewIpfsAdder(config.ArchivalDir, coreAPI)

			err = ipfs.Remove(ctx, modelArchivalDir, modelArchivalDir.Name())
			fmt.Println(modelArchivalDir, modelArchivalDir.Name())
			if err != nil {
				printAndDie(err)
			}
			slog.Debug("removed model", "id", identifier)
		}
	},
}
