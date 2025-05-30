package config

import (
	"github.com/spf13/viper"
)

// Configuration struct to hold environment variables
type Config struct {
	DB       DB
	Token    Token
	App      App
	External External
}

type App struct {
	Name string
	Env  string
	Port string
}

type External struct {
	WarehouseServiceURL         string
	WarehouseServiceStaticToken string
	ProductServiceURL           string
}

type DB struct {
	DSN string
}

type Token struct {
	JWTSecret []byte
	JWTStatic string
}

func LoadConfig() (*Config, error) {

	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, err
	}

	c := &Config{
		App: App{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
		},
		DB: DB{
			DSN: viper.GetString("DB_DSN"),
		},
		Token: Token{
			JWTSecret: []byte(viper.GetString("JWT_SECRET")),
			JWTStatic: viper.GetString("JWT_STATIC"),
		},
		External: External{
			WarehouseServiceURL:         viper.GetString("WAREHOUSE_SERVICE_URL"),
			WarehouseServiceStaticToken: viper.GetString("WAREHOUSE_SERVICE_STATIC_TOKEN"),
			ProductServiceURL:           viper.GetString("PRODUCT_SERVICE_URL"),
		},
	}

	return c, nil
}
