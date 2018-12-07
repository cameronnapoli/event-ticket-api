// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
    "fmt"
    // "encoding/json"
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "github.com/go-redis/redis"
    "strconv"
    // "crypto/sha1"
)


const TICKET_GA = 0
const TICKET_VIP = 1
const TICKET_ONE_DAY = 2
const LOCK_TIME_MS = 1000 * 60 * 5 // 5 minute lock time
const INITIAL_TICKET_COUNT = 1000

const GLOBAL_DEBUG = true

var GlobalRedisClient *redis.Client

type TicketResponse struct {
    UserToken string `json:"user_token"`
    TicketType int `json:"ticket_type"`
    PaymentToken string `json:"payment_token"`
}

/*
Build a go application that can handle high volume festival ticket purchases
----------------------------------------------------------------------------
Features:
    - Lock ticket purchase in for user on IP (Max number of users per IP)
    - Process user token, ticket type, and "payment"
    - Display how many tickets remain (open, purchasing, purchased)
*/

func generateHash(ticketNum int) string {
    return "cf23df2207d99a74fbe169e3eba035e633b65d94"
}


func InitializeRedisClient() *redis.Client {
    return  redis.NewClient(&redis.Options{
        Addr: "redis:6379", Password: "", DB: 0,
    })
}


// Create JSON key in Redis DB with initialized ticket count
func InitializeTickets() {
    // Set num_tickets in Redis to INITIAL_TICKET_COUNT
    err := GlobalRedisClient.Set("num_tickets", INITIAL_TICKET_COUNT, 0).Err()
	if err != nil {
		panic(err)
	}
}


// Return the remaining ticket count in JSON form
func GetRemainingTickets(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* GetRemainingTickets")

    var val string

    if !GLOBAL_DEBUG {
        tmp_val, err := GlobalRedisClient.Get("num_tickets").Result()
    	if err != nil {
    		panic(err)
    	}
        val = tmp_val
    } else {
        val = "200"
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"num_tickets": %s}`, val)
}


// Create lock on ticket to allow purchasing period
func LockTicket(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* LockTicket")


    if !GLOBAL_DEBUG {
        // Redis to drop ticket count
        numTickets, err := GlobalRedisClient.Get("num_tickets").Result()
        if err != nil {
            panic(err)
        }
        intNumTickets, err2 := strconv.Atoi(numTickets)
        if err2 != nil {
            panic(err2)
        }

        // Add salt to create more secure hash
        hash := generateHash(intNumTickets)

        // Write hash to redis
        err = GlobalRedisClient.Set(hash, intNumTickets, 0).Err()
    	if err != nil {
    		panic(err)
    	}

        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"process": "success", "token": %s}`, hash)
    } else {

    }
}


func CompleteTicketPurchase(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* CompleteTicketPurchase")
    // If lock is still valid, finish ticket purchase

}


// Release lock on ticket
func ReleaseTicket() {
    fmt.Println("* ReleasingTicket")

}


func main() {
    if !GLOBAL_DEBUG {
        GlobalRedisClient = InitializeRedisClient()
        InitializeTickets()
    }

    router := mux.NewRouter()

    router.HandleFunc("/remaining_tickets", GetRemainingTickets).Methods("GET")
    router.HandleFunc("/buy_ticket", LockTicket).Methods("POST")
    router.HandleFunc("/buy_ticket/{token}", CompleteTicketPurchase).Methods("POST")

    fmt.Println("Listening on port 8000.")
    log.Fatal(http.ListenAndServe(":8000", router))


    // tr := &TicketResponse{UserToken: "0f8238n2fn803f2", TicketType: TICKET_GA, PaymentToken: "82748719712"}
    // b, err := json.Marshal(tr)
    // if err != nil {
    //     fmt.Println(err); return
    // } fmt.Println(string(b))
}
