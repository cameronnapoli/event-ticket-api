// Written by: Cameron Napoli
// Date Created: Dec. 5, 2018

package main

import (
	"github.com/go-redis/redis"
	"time"
	// "crypto/sha1"
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

func generateToken(ticketNum int) string {
	// TOOD: implement
	return "cf23df2207d99a74fbe169e3eba035e633b65d94"
}

func WriteErrorResponse(w *http.ResponseWriter, err string) {
	(*w).WriteHeader(403)
	fmt.Fprintf(*w, `{"success": false, "errorMessage": "%s"}`, err)
}

func InitializeRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "redis:6379", Password: "", DB: 0,
	})
}
