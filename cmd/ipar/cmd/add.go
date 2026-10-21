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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a path to your IPAR",
	Long:  `...`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var config struct {
			ArchivalDir *pathlib.Path   `mapstructure:"archival_dir"`
			Paths       []*pathlib.Path `mapstructure:"paths"`
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
			return pathlib.NewPath(userHomeDir).Join("ipar", "self").String()
		}())
		vi.Set("paths", args)
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

		for _, p := range config.Paths {
			pathArchivalLink := config.ArchivalDir.Join(p.Name())
			if exists, err := pathArchivalLink.Exists(); !exists {
				if err != nil {
					printAndDie(err)
				}
				if err := pathArchivalLink.Symlink(p); err != nil {
					printAndDie(err)
				}
			} else {
				if isSymlink, err := pathArchivalLink.IsSymlink(); !isSymlink {
					if err != nil {
						printAndDie(err)
					}
				}
				// TODO: Check if symlink already points to <p>. If not,
				// panic.
			}

			coreAPI, err := rpc.NewLocalApi()
			if err != nil {
				printAndDie(err)
			}
			ipfs := ipar.NewIpfsAdder(config.ArchivalDir, coreAPI)

			path, err := ipfs.Add(ctx, pathArchivalLink, pathArchivalLink.Name())
			if err != nil {
				printAndDie(err)
			}
			cidArchivalLink := config.ArchivalDir.Join(path.RootCid().String())
			if err := cidArchivalLink.Symlink(p); err != nil {
				printAndDie(err)
			}

			slog.Debug("added path", "dir", pathArchivalLink.String(), "cid", path.RootCid().String())
		}
	},
}
