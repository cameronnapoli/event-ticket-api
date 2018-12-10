// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
    "fmt"
    "log"
    "time"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    "strings"
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
    client := GetRedisClient()
    err := client.Set("num_tickets", INITIAL_TICKET_COUNT, 0).Err()
    if err != nil {
        panic(err)
    }
}


// Return the remaining ticket count in JSON form
func GetRemainingTickets(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* GetRemainingTickets")
    client := GetRedisClient()
    var val string

    if !GLOBAL_DEBUG {
        tmp_val, err := client.Get("num_tickets").Result()
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
    client := GetRedisClient()
    params := mux.Vars(r) // map[string]string

    argsErr := CheckArgsInParams(params, "ticketType")
    if argsErr != nil {
        WriteErrorResponse(&w, "parameters not all found")
    }

    ticketType, _ := params["ticketType"]

    if !GLOBAL_DEBUG {
        // Redis to drop ticket count
        numTickets, err := client.Get("num_tickets").Result()
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
        token := generateToken()

        // Write ticket type to Redis as well
        err4 := client.Set(token, ticketType, 0).Err()
        if err4 != nil {
            panic(err4)
        }

        // Decrement number of tickets
        _, err5 := client.Decr("num_tickets").Result()
        if err5 != nil {
            panic(err5)
        }

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
        fmt.Fprintf(w, `{"success": true, "token": %s}`, generateToken())
    }
}


// Endpoint function to finalize ticket purchase
func CompleteTicketPurchase(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* CompleteTicketPurchase")
    client := GetRedisClient()

    params := mux.Vars(r) // map[string]string
    argsErr := CheckArgsInParams(params, "token", "paymentToken")
    if argsErr != nil {
        WriteErrorResponse(&w, "parameters not all found")
    }

    token, _ := params["token"]
    paymentToken, _ := params["paymentToken"]

    // see if r.token exists in Redis
    ticketType, err := client.Get(token).Result()
    if err != nil {
        WriteErrorResponse(&w, strings.Join([]string{"token '", token, "' not found in RedisDB."}, ""))
        return
    }

    ticketTypeVal, err2 := strconv.Atoi(ticketType)
    if err2 != nil {
        panic(err2)
    }

    tp := &TicketPaymentPayload{UserToken: token,
                                TicketType: ticketTypeVal,
                                PaymentToken: paymentToken}
    // if it does add to second completed purchases table with token, ticketType
    err3 := client.SAdd("purchases", payloadToJson(tp), 0).Err()
    if err3 != nil {
        panic(err3)
    }

    ReleaseTicket(token)
    BasicSuccessResponse(&w)
}


// Release lock on ticket
func ReleaseTicket(token string) {
    fmt.Println("* ReleasingTicket")
    client := GetRedisClient()

    // see if token exists in Redis
    _, err := client.Get(token).Result()
    if err != nil {
        fmt.Printf("Token '%s' doesn't exist in RedisDB.", token)
        return;
    }

    // Remove token from purchase pool if it does
    _, err2 := client.Del(token).Result()
    if err2 != nil {
        fmt.Printf("Token delete failed for '%s'.", token)
        return;
    }
}


func main() {
    if !GLOBAL_DEBUG {
        ResetDB()
        InitializeTickets()
    }

    router := mux.NewRouter()

    router.HandleFunc("/remaining_tickets", GetRemainingTickets).Methods("GET")
    router.HandleFunc("/buy_ticket", LockTicket).Methods("POST")
    router.HandleFunc("/buy_ticket/{token}", CompleteTicketPurchase).Methods("POST")

    fmt.Println("Listening on port 8000.")
    log.Fatal(http.ListenAndServe(":8000", router))
}
