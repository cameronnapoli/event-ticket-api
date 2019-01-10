// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
    "bytes"
    "crypto/sha1"
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    "errors"
    "time"
)

//================ GLOBAL VARS ================
const TICKET_GA = "GA"
const TICKET_VIP = "VIP"
const TICKET_ONE_DAY = "ONE_DAY"
const LOCK_TIME = time.Second * 30 * 1 // 30 second lock time
const INITIAL_TICKET_COUNT = 50000

var BaseTimestamp int64 = time.Now().UnixNano() / 100000 // Random seed value
var TokenTicker int64 = 0

type TicketPaymentPayload struct {
    UserToken    string `json:"user_token"`
    TicketType   int    `json:"ticket_type"`
    PaymentToken string `json:"payment_token"`
}

//================ HELPERS ================
func payloadToJson(tp *TicketPaymentPayload) string {
    s, err := json.Marshal(tp)
    if err != nil {
        panic(err)
    }
    return string(s)
}

func concatStrings(strs... string) string {
    var b bytes.Buffer
    for _, s := range strs {
        b.WriteString(s)
    }
    return b.String()
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

func CheckTicketType(ticketType string) int {
    switch ticketType {
        case TICKET_GA:
            return 0
        case TICKET_VIP:
            return 1
        case TICKET_ONE_DAY:
            return 2
        default:
            return -1
    }
}
