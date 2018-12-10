package main

import (
	"fmt"
	"time"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
)

var BaseTimestamp int64 = time.Now().UnixNano() / 1000 // Random seed value
var TokenTicker int64 = 0

func generateToken() string {
    genToken := BaseTimestamp + TokenTicker
    TokenTicker++
    // Generate hash from timestamp
    byteArray := make([]byte, 8)
    binary.LittleEndian.PutUint64(byteArray, uint64(genToken))
    sum := sha1.Sum(byteArray)
    return hex.EncodeToString(sum[:])
}

func main() {
	fmt.Println("Starting...\n")
	// fmt.Println(time.Now().UnixNano())
	// fmt.Println(time.Now().UnixNano())

    set := make(map[string]bool)

    iter := 1000000

	for i := 0; i < iter; i++ {
        token := generateToken()
        _, hasKey := set[token]

        if hasKey {
            fmt.Printf("*** Collision at iter=%d\n", i)
            fmt.Printf("    token: %s\n", token)
            fmt.Printf("    BaseTimestamp: %d\n", BaseTimestamp)
            panic(0)
        } else {
            set[token] = true
        }

        // fmt.Printf("%4d) ", i)
        // fmt.Print(token)
        // fmt.Println()
	}

    fmt.Printf("* Completed %d iterations.\n", iter)
}
