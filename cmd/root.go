package cmd

import (
	"log"
	"os"

	"github.com/Lofter1/githook-manager/cmd/prepush"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigDir  = "./.git"
	defaultConfigName = ".githook-manager"
	defaultConfigType = "yaml"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "githook-manager",
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

	rootCmd.AddCommand(prepush.PrepushCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.githook-manager.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(defaultConfigDir)
		viper.SetConfigType(defaultConfigType)
		viper.SetConfigName(defaultConfigName)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
}
