// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "strconv"
    "time"
)

/*
Go application that can handle high volume festival ticket purchases
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
        log.Fatal(err)
    }
}

// Return the remaining ticket count in JSON form
func GetRemainingTickets(w http.ResponseWriter, r *http.Request) {
    client := GetRedisClient()

    val, err := client.Get("num_tickets").Result()
    if err != nil {
        log.Fatal(err)
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"num_tickets": %s}`, val)
}

// Create lock on ticket to allow purchasing period
func LockTicket(w http.ResponseWriter, r *http.Request) {
    client := GetRedisClient()
    params := mux.Vars(r) // map[string]string

    argsErr := CheckArgsInParams(params, "type")
    if argsErr != nil {
        WriteErrorResponse(&w, "parameters not all found")
        return
    }

    ticketType, _ := params["type"]
    tVal := CheckTicketType(ticketType)
    if tVal < 0 {
        WriteErrorResponse(&w, "ticket type is invalid.")
        return
    }

    // Redis query to get ticket count
    numTickets, err := client.Get("num_tickets").Result()
    if err != nil {
        log.Fatal(err)
    }

    var intNumTickets int
    intNumTickets, err = strconv.Atoi(numTickets)
    if err != nil {
        log.Fatal(err)
    }

    if intNumTickets == 0 {
        WriteErrorResponse(&w, "No tickets remaining.")
        return
    } else if intNumTickets < 0 {
        log.Fatal("number of tickets in redis < 0. Something went very wrong.")
    }

    // Store token in Redis
    token := generateToken()
    err = client.Set(token, tVal, 0).Err()
    if err != nil {
        log.Fatal(err)
    }

    // Decrement number of tickets
    _, err = client.Decr("num_tickets").Result()
    if err != nil {
        log.Fatal(err)
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"process": "success", "token": %s}`, token)

    doneChan := make(chan bool)
    go func() {
        time.Sleep(LOCK_TIME)
        ReleaseTicket(token, true)
        doneChan <- true
    }()
}

// Endpoint function to finalize ticket purchase
func CompleteTicketPurchase(w http.ResponseWriter, r *http.Request) {
    // NOTE: in a real platform, there would be a step here to process a payment
    //   authentication for the type of ticket the user purchased.
    client := GetRedisClient()

    params := mux.Vars(r)                         // map[string]string
    argsErr := CheckArgsInParams(params, "token") //, "paymentToken")
    if argsErr != nil {
        WriteErrorResponse(&w, "parameters not all found")
        return
    }

    token, _ := params["token"]
    // paymentToken, _ := params["paymentToken"]

    // see if r.token exists in Redis
    ticketType, err := client.Get(token).Result()
    if err != nil {
        WriteErrorResponse(&w, concatStrings("token '", token, "' not found in RedisDB."))
        return
    }

    var ticketTypeVal int
    ticketTypeVal, err = strconv.Atoi(ticketType)
    if err != nil {
        log.Fatal(err)
    }

    tp := &TicketPaymentPayload{UserToken: token,
        TicketType:   ticketTypeVal,
        PaymentToken: ""} // PaymentToken: paymentToken}
    // if it does add to second completed purchases table with token, ticketType
    err = client.SAdd("purchases", payloadToJson(tp), 0).Err()
    if err != nil {
        log.Fatal(err)
    }

    ReleaseTicket(token, false)
    BasicSuccessResponse(&w)
}

// Release lock on ticket
func ReleaseTicket(token string, incr bool) {
    client := GetRedisClient()

    // see if token exists in Redis
    _, err := client.Get(token).Result()
    if err != nil {
        fmt.Printf("Token '%s' doesn't exist in Redis.\n", token)
        return
    }

    // Remove token from purchase pool if it does
    _, err = client.Del(token).Result()
    if err != nil {
        fmt.Printf("Token delete failed for '%s'.", token)
        return
    }

    if incr {
        _, err = client.Incr("num_tickets").Result()
        if err != nil {
            log.Fatal(err)
        }
    }
}

func main() {
    ResetDB()
    InitializeTickets()

    router := mux.NewRouter()
    router.HandleFunc("/remaining_tickets", GetRemainingTickets).Methods("GET")
    router.HandleFunc("/buy_ticket/{type}", LockTicket).Methods("POST")
    router.HandleFunc("/complete_purchase/{token}", CompleteTicketPurchase).Methods("POST")

    fmt.Println("Listening on port 8000.")
    log.Fatal(http.ListenAndServe(":8000", router))
}
