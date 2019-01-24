package main

import (
    "fmt"
    "github.com/go-redis/redis"
    "github.com/micro/go-config"
    "log"
    "strconv"
)

var GlobalRedisClient *redis.Client = nil

func GetRedisConfig() (Addr, Password string, DB int) {
    config.LoadFile("config.json")
    return config.Get("redis", "Addr").String("localhost:6379"), // localhost:6379 default value
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

func GetNumTickets(client *redis.Client) int {
    numTickets, err := client.Get("num_tickets").Result()
    if err != nil {
        log.Fatal(err)
    }

    intNumTickets, err := strconv.Atoi(numTickets)
    if err != nil {
        log.Fatal(err)
    }

    return intNumTickets
}

func ResetDB()  {
    client := GetRedisClient()
    _, err := client.FlushAll().Result()
    if err != nil {
        fmt.Println("* Failed to flush Redis DB")
        panic(err)
    }
}
