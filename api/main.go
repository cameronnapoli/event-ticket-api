// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
    "fmt"
    "github.com/go-redis/redis"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "strconv"
    // "encoding/json"
    "time"
)

/*
Build a go application that can handle high volume festival ticket purchases
----------------------------------------------------------------------------
Features:
    - Lock ticket purchase in for user on IP (Max number of users per IP)
    - Process user token, ticket type, and "payment"
    - Display how many tickets remain (open, purchasing, purchased)
*/

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

    params := mux.Vars(r) // map[string]string

    var ticketType int

    if val, ok := params["ticket_type"]; ok {
        iVal, err := strconv.Atoi(val)
        if err != nil {
            WriteErrorResponse(&w, err.Error())
        }
        ticketType = iVal

    } else {
        WriteErrorResponse(&w, "ticket_type not found in parameters")
    }

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

        if intNumTickets == 0 {
            WriteErrorResponse(&w, "No tickets remaining.")
        } else if intNumTickets < 0 {
            panic("number of tickets in redis < 0. Something went very wrong.")
        }

        // Add salt to create more secure hash
        token := generateToken(intNumTickets)

        // Write hash to redis with value being the lock
        err3 := GlobalRedisClient.SAdd("token_set", token, 0).Err()
        if err3 != nil {
            panic(err3)
        }

        // Write ticket type to Redis as well
        err4 := GlobalRedisClient.Set(token, ticketType, 0).Err()
        if err4 != nil {
            panic(err4)
        }

        // Decrement number of tickets
        result, err5 := GlobalRedisClient.Decr("num_tickets").Result()
        if err5 != nil {
            panic(err5)
        }
        fmt.Println(result)

        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"process": "success", "token": %s}`, token)

        // Create callback to ReleaseTicket(token)
        // TODO: make sure this is a nonblocking action
        doneChan := make(chan bool)
        go func() {
            time.Sleep(LOCK_TIME)
            ReleaseTicket(token)
            doneChan <- true
        }()
    } else {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"process": "success", "token": %s}`, generateToken(0))
    }
}

func CompleteTicketPurchase(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* CompleteTicketPurchase")

    // If lock is still valid, finish ticket purchase

    // get token from r

    // see if r.token exists in Redis

    // if it does add to second completed purchases table with token, ticketType

}

// Release lock on ticket
func ReleaseTicket(token string) {
    fmt.Println("* ReleasingTicket")

    // see if token exists in Redis set

    // Remove token from purchase pool if it does

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
