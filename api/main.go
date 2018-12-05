package main

// https://github.com/danhenriquesc/go-poll

import (
    "fmt"
    // "encoding/json"
    "net/http"
    "log"
    "github.com/gorilla/mux"
    // "github.com/go-redis/redis"
)


const TICKET_GA = 0
const TICKET_VIP = 1
const TICKET_ONE_DAY = 2



type TicketResponse struct {
    UserToken string `json:"user_token"`
    TicketType int `json:"ticket_type"`
    PaymentToken string `json:"payment_token"`
}


// Build a go application that can handle high volume festival ticket purchases

/*
Features:

Lock ticket purchase in for user on IP (Max number of users per IP)

Process user token, ticket type, and "payment"

Display how many tickets remain (open, purchasing, purchased)

*/

func GetRemainingTickets(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* GetRemainingTickets")

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "<p>Html content.</p>")
}

func LockTicket(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* LockTicket")
}

func BuyTicket(w http.ResponseWriter, r *http.Request) {
    fmt.Println("* BuyTicket")
}


func main() {
    router := mux.NewRouter()

    router.HandleFunc("/remaining_tickets", GetRemainingTickets).Methods("GET")
    router.HandleFunc("/buy_ticket", LockTicket).Methods("POST")
    router.HandleFunc("/buy_ticket/{token}", BuyTicket).Methods("POST")

    fmt.Println("Listening on port 8000.")
    log.Fatal(http.ListenAndServe(":8000", router))


    // tr := &TicketResponse{UserToken: "0f8238n2fn803f2",
    //                       TicketType: TICKET_GA,
    //                       PaymentToken: "82748719712"}
    //
    // b, err := json.Marshal(tr)
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // fmt.Println(string(b))
}
