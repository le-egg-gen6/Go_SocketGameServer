package config

import "github.com/spf13/viper"

type Config struct {

	//Database
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
	DB_HOST     string
	DB_PORT     int

	//Http Server
	HTTP_PORT  int
	SECRET_KEY string

	//Socket Server
	SOCKET_PORT int
}

func LoadConfig() *Config {
	viper.SetConfigName("server_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading server config " + err.Error())
	}

	return &Config{
		// Database
		DB_USERNAME: viper.GetString("database.username"),
		DB_PASSWORD: viper.GetString("database.password"),
		DB_NAME:     viper.GetString("database.name"),
		DB_HOST:     viper.GetString("database.host"),
		DB_PORT:     viper.GetInt("database.port"),

		// HTTP Server
		HTTP_PORT:  viper.GetInt("http.port"),
		SECRET_KEY: viper.GetString("http.secret_key"),

		// Socket Server
		SOCKET_PORT: viper.GetInt("socket.port"),
	}
}
