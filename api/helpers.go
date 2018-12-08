// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
    "github.com/go-redis/redis"
    "time"
    "fmt"
    "net/http"
    "encoding/json"
    "errors"
    "crypto/sha1"
    "encoding/binary"
    "encoding/hex"
)

const TICKET_GA = 0
const TICKET_VIP = 1
const TICKET_ONE_DAY = 2
const LOCK_TIME = time.Second * 5 // 5 minute lock time in seconds
const INITIAL_TICKET_COUNT = 1000

const GLOBAL_DEBUG = true

var GlobalRedisClient *redis.Client


type TicketPaymentPayload struct {
    UserToken    string `json:"user_token"`
    TicketType   int    `json:"ticket_type"`
    PaymentToken string `json:"payment_token"`
}


func payloadToJson(tp *TicketPaymentPayload) string {
    s, err := json.Marshal(tp)
    if err != nil {
        panic(err)
    }
    return string(s)
}


func generateToken(ticketNum int) string {
    now := time.Now().UnixNano() + ticketNum
    byteArray := make([]byte, 8)
    binary.LittleEndian.PutUint64(byteArray, uint64(now))
    sum := sha1.Sum(byteArray)
    return hex.EncodeToString(sum[:])
}


func CheckArgsInParams(params map[string]string, reqArgs... string) error {
    for _, reqArg := range reqArgs {
        if _, ok := params[reqArg]; !ok {
            return errors.New("Argument missing from request.")
        }
    }
    return nil
}


func WriteErrorResponse(w *http.ResponseWriter, err string) {
    (*w).WriteHeader(403)
    fmt.Fprintf(*w, `{"success": false, "errorMessage": "%s"}`, err)
}


func BasicSuccessResponse(w *http.ResponseWriter) {
    (*w).WriteHeader(200)
    fmt.Fprintf(*w, `{"success": true}`)
}


func InitializeRedisClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: "redis:6379", Password: "", DB: 0,
    })
}
