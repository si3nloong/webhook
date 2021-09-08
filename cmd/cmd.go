package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Enabled        bool   `mapstructure:"enabled"`
	Port           string `mapstructure:"port" validate:"numeric"`
	Retry          uint   `mapstructure:"retry" validate:"lte=50"`
	RetryMechanism string `mapstructure:"retry_mechanism" validate:""`
	GRPC           struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    string `mapstructure:"port" validate:"numeric"`
	} `mapstructure:"grpc"`
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "SianLoong", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "MIT")
	rootCmd.PersistentFlags().StringP("version", "v", "", "1.0.0-alpha0")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	// log.Println("here 2")
	// viper.SetDefault("config", "config.yaml")
	// viper.SetDefault("CURLHOOK_REDIS_PORT", "6379")
	// viper.SetDefault("CURLHOOK_REDIS_CLUSTER", false)
	// viper.SetDefault("grpc.")

	log.Println(cfgFile)
	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
