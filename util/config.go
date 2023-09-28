package util

import "github.com/spf13/viper"

//Config stores all configurations of the application
// The values are read by viper from a config file ar envieroment variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

//LoadConfig reads configuration from file or eviroment variables.
func LoadConfig(path string) (confing Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&confing)
	return
}
