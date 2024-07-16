package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort             string
	Dsn                    string // Data Source Name or DB_URL
	AppSecret              string
	TwillioAccountSid      string
	TwillioAuthToken       string
	TwillioFromPhoneNumber string
}

func SetupEnv() (cfg AppConfig, err error) {

	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	dsn := os.Getenv("DSN")
	if len(dsn) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	twillioAccountSid := os.Getenv("TWILLIO_ACCOUNT_SID")
	if len(twillioAccountSid) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	twillioAuthToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	if len(twillioAuthToken) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	twillioFromPhoneNumber := os.Getenv("TWILLIO_FROM_PHONE_NUMBER")
	if len(twillioFromPhoneNumber) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	return AppConfig{
		ServerPort:             httpPort,
		Dsn:                    dsn,
		AppSecret:              appSecret,
		TwillioAccountSid:      twillioAccountSid,
		TwillioAuthToken:       twillioAuthToken,
		TwillioFromPhoneNumber: twillioFromPhoneNumber,
	}, nil
}
