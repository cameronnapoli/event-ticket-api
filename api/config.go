package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		Database int    `json:"database"`
	} `json:"redis"`
	Mysql struct {
		Address  string `json:"address"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"mysql"`
}

func GetConfigObject() Configuration {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Can't open config file: ", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	Config := Configuration{}

	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("Can't decode config JSON: ", err)
	}

	return Config
}

func GetSQLConfig() (address, username, password, database string) {
	Config := GetConfigObject()
    return Config.Mysql.Address,
		   Config.Mysql.Username,
		   Config.Mysql.Password,
		   Config.Mysql.Database
}

func GetRedisConfig() (Addr, Password string, DB int) {
	Config := GetConfigObject()
	return Config.Redis.Address,
		   Config.Redis.Password,
		   Config.Redis.Database
}