/*
Copyright © 2023 Netmaker Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gravitl/devops/netmaker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var (
	cfgFile  string
	debug    bool
	config   netmaker.Config
	logger   *slog.Logger
	logLevel *slog.LevelVar
	replace  func(groups []string, a slog.Attr) slog.Attr
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "test command",
	Args:  cobra.ExactArgs(1),
	Short: "run tests",
	Long: `run an integration test on netmaker network
	.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "use debugg logging")
	rootCmd.PersistentFlags().BoolP("cleanup", "c", false, "use to remove netmaker interface and config from extclient")
	rootCmd.PersistentFlags().StringP("network", "n", "devops", "network name")
	rootCmd.PersistentFlags().String("digitalocean_token", "", "digitalocean token")
	rootCmd.PersistentFlags().StringP("api", "a", "https://api.clustercat.com", "api endpoint")
	rootCmd.PersistentFlags().StringP("tag", "t", "devops", "digital ocean droplet tag")
	rootCmd.PersistentFlags().StringP("masterkey", "m", "secretkey", "netmaker masterkey")
	rootCmd.PersistentFlags().Int("timeout", 1, "cleaup timeout in minutes ")
	rootCmd.PersistentFlags().StringP("key", "k", "", "ssh private key")
	rootCmd.PersistentFlags().StringArray("ranges", []string{"10.0.3.0/24"}, "egress gateway range(s)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".test" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".test")
	}

	viper.BindPFlags(rootCmd.Flags())
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	viper.Unmarshal(&config)
	if config.Key == "" {
		key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_ed25519")
		if err != nil {
			cobra.CheckErr("invalid configuration, ssh key not set" + err.Error())
		}
		config.Key = string(key)
	}
	netmaker.SetCxt(config.Api, config.Masterkey)
}

func setupLoging(name string) {
	// setup logging
	f, err := os.Create(os.TempDir() + "/" + name + ".log")
	cobra.CheckErr(err)
	//defer f.Close() -- don't close file here
	logLevel = &slog.LevelVar{}
	replace = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			a.Value = slog.StringValue(filepath.Base(a.Value.String()))
		}
		return a
	}
	logger = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stderr, f), &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace, Level: logLevel}))
	slog.SetDefault(logger)
	slog.With("TEST", name)
	if debug {
		logLevel.Set(slog.LevelDebug)
	}
}
