package main

import (
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api"
	"log"
)

func main(){
	cfg, err := config.SetupEnv()
	if err != nil {
		log.Fatalln("Config file is not loaded properly", err)
	}
	api.StartServer(cfg)
}