package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/chigopher/pathlib"
	"github.com/go-viper/mapstructure/v2"
	"github.com/iparchive/ipar-cli"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hfAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Huggingface model to your IPAR",
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

		hfRepo := ipar.NewHuggingfaceRepositoryHttpAPI(os.Getenv("HF_TOKEN"))

		for _, identifier := range config.Identifiers {
			modelArchivalDir := config.ArchivalDir.Join(identifier)

			slog.Debug("getting model files...", "id", identifier, "modelArchivalDir", modelArchivalDir.String())
			model, err := hfRepo.DownloadModel(identifier, modelArchivalDir)
			if err != nil {
				printAndDie(err)
			}

			coreAPI, err := rpc.NewLocalApi()
			if err != nil {
				printAndDie(err)
			}

			ipfs := ipar.NewIpfsAdder(config.ArchivalDir, coreAPI)

			path, err := ipfs.Add(ctx, model.ModelArchivalDir, modelArchivalDir.Name())
			if err != nil {
				printAndDie(err)
			}
			slog.Debug("got model", "id", identifier, "path", path)
		}
	},
}
