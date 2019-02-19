package main

import (
    "github.com/go-redis/redis"
    "log"
    "strconv"
)

var GlobalRedisClient *redis.Client = nil

func GetRedisClient() *redis.Client {
    Addr, Password, DB := GetRedisConfig()
    if GlobalRedisClient == nil {
        GlobalRedisClient = redis.NewClient(&redis.Options{Addr: Addr, Password: Password, DB: DB})
    }
    _, err := GlobalRedisClient.Ping().Result()
    if err != nil {
        log.Printf("Cannot connect to Redis server at address: '%s'", Addr)
        log.Fatal(err)
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
        log.Print("Failed to flush Redis DB")
        log.Fatal(err)
    }
}
