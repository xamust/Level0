package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"service/internal/app/service"
)

var configsPath string

func init() {
	flag.StringVar(&configsPath, "configs-path", "configs/config.toml", "Path to configs...")
}

func main() {
	flag.Parse()
	configs := service.NewConfig()
	meta, err := toml.DecodeFile(configsPath, configs)
	if err != nil {
		log.Fatalln(err)
	}

	if len(meta.Undecoded()) != 0 {
		log.Fatal("Undecoded configs param: ", meta.Undecoded())
	}

	//start service...
	service := service.NewService(configs)
	if err = service.Start(); err != nil {
		log.Fatalln("Error on start:", err)
	}
}
