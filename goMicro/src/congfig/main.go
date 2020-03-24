package main

import (
	"log"

	"github.com/micro/go-micro/v2/config"
	//"github.com/micro/go-micro/config/encoder/toml"
	//"github.com/micro/go-micro/config/source"
	//"github.com/micro/go-micro/v2/config/source/file"
)

func main() {
	conf, _ := config.NewConfig()
	log.Printf("%T", conf)
	err := config.LoadFile("conf.json")
	if err != nil {
		log.Printf("load config failed: %v.", err)
	}

	//enc := toml.NewEncoder()
	//config.Load(file.NewSource(file.WithPath("tomlconf"), source.WithEncoder(enc)))
	log.Printf("config map: %v.",config.Map())
	log.Printf("host: %v.",config.Map()["hosts"])
}
