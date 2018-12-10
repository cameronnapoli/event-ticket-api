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


//=============================================
//================ GLOBAL VARS ================
//=============================================
const TICKET_GA = 0
const TICKET_VIP = 1
const TICKET_ONE_DAY = 2
const LOCK_TIME = time.Second * 5 // 5 minute lock time in seconds
const INITIAL_TICKET_COUNT = 100000

const GLOBAL_DEBUG = true

var GlobalRedisClient *redis.Client = nil

var BaseTimestamp int64 = time.Now().UnixNano() / 1000 // Random seed value
var TokenTicker int64 = 0

type TicketPaymentPayload struct {
    UserToken    string `json:"user_token"`
    TicketType   int    `json:"ticket_type"`
    PaymentToken string `json:"payment_token"`
}


//=========================================
//================ HELPERS ================
//=========================================
func payloadToJson(tp *TicketPaymentPayload) string {
    s, err := json.Marshal(tp)
    if err != nil {
        panic(err)
    }
    return string(s)
}


func generateToken() string {
    genToken := BaseTimestamp + TokenTicker
    TokenTicker++
    // Generate hash from timestamp
    byteArray := make([]byte, 8)
    binary.LittleEndian.PutUint64(byteArray, uint64(genToken))
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


//=========================================
//================ ROUTING ================
//=========================================
func WriteErrorResponse(w *http.ResponseWriter, err string) {
    (*w).WriteHeader(403)
    fmt.Fprintf(*w, `{"success": false, "errorMessage": "%s"}`, err)
}


func BasicSuccessResponse(w *http.ResponseWriter) {
    (*w).WriteHeader(200)
    fmt.Fprintf(*w, `{"success": true}`)
}


//=======================================
//================ REDIS ================
//=======================================
func GetRedisClient() *redis.Client {
    if GlobalRedisClient == nil {
        GlobalRedisClient = redis.NewClient(&redis.Options{
            Addr: "redis:6379", Password: "", DB: 0,
        })
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
