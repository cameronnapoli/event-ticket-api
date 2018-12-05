package main

import (
    "fmt"
    "encoding/json"
    // "github.com/gorilla/mux"
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

func


func main() {
    fmt.Println("Starting application...")

    tr := &TicketResponse{UserToken: "0f8238n2fn803f2",
                          TicketType: TICKET_GA,
                          PaymentToken: "82748719712"}

    b, err := json.Marshal(tr)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(b))
}
