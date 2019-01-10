package main

import (
    "fmt"
    "github.com/go-redis/redis"
    "github.com/micro/go-config"
)

var GlobalRedisClient *redis.Client = nil

func GetRedisConfig() (Addr, Password string, DB int) {
    config.LoadFile("config.json")
    return config.Get("redis", "Addr").String("localhost:6379"), // default value
           config.Get("redis", "Password").String(""),
           config.Get("redis", "Addr").Int(0)
}

func GetRedisClient() *redis.Client {
    Addr, Password, DB := GetRedisConfig()
    if GlobalRedisClient == nil {
        GlobalRedisClient = redis.NewClient(&redis.Options{Addr: Addr, Password: Password, DB: DB})
    }
    return GlobalRedisClient
}

func ResetDB()  {
    client := GetRedisClient()
    _, err := client.FlushAll().Result()
    if err != nil {
        fmt.Println("* Failed to flush Redis DB")
        panic(err)
    }
}
