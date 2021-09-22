package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type MessageQueueEngine string

const (
	MessageQueueEngineRedis MessageQueueEngine = "redis"
	MessageQueueEngineNats  MessageQueueEngine = "nats"
	MessageQueueEngineNSQ   MessageQueueEngine = "nsq"
)

type DatabaseEngine string

const (
	DatabaseEngineElasticsearch DatabaseEngine = "elasticsearch"
)

type Config struct {
	Enabled        bool   `mapstructure:"enabled"`
	Port           int    `mapstructure:"port"`
	Retry          uint   `mapstructure:"retry" validate:"lte=50"`
	RetryMechanism string `mapstructure:"retry_mechanism"`
	NoOfWorker     int    `mapstructure:"no_of_worker" validate:"required,lte=100"`
	Monitor        struct {
		Enabled bool `mapstructure:"enabled"`
		Port    int  `mapstructure:"port"`
	} `mapstructure:"monitor"`
	DB struct {
		Engine DatabaseEngine `mapstructure:"engine" validate:"oneof=elasticsearch"`
		Host   string         `mapstructure:"host" validate:""`
	} `mapstructure:"db"`
	GRPC struct {
		Enabled bool   `mapstructure:"enabled"`
		ApiKey  string `mapstructure:"api_key"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"grpc"`
	MessageQueue struct {
		Engine     MessageQueueEngine `mapstructure:"engine" validate:"oneof=redis nats nsq"`
		Topic      string             `mapstructure:"topic" validate:"alphanum"`
		QueueGroup string             `mapstructure:"queue_group" validate:"alphanum"`
		Redis      struct {
			Cluster  bool   `mapstructure:"cluster"`
			Addr     string `mapstructure:"addr"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			DB       int    `mapstructure:"db"`
		} `mapstructure:"redis"`
		NATS struct {
			JetStream bool `mapstructure:"js"`
		} `mapstructure:"nats"`
		NSQ struct {
		} `mapstructure:"nsq"`
	} `mapstructure:"message_queue"`
}

func (c *Config) SetDefault() {
	c.NoOfWorker = runtime.NumCPU()
	c.Enabled = true
	c.Port = 3000
	c.Monitor.Port = 3222
	c.GRPC.Port = 5222
	c.MessageQueue.Redis.Addr = "localhost:6379"
	c.MessageQueue.NATS.JetStream = true
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

	// log.Println(cfgFile)
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
