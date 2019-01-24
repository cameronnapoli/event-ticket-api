package main

import (
    "fmt"
    "net/http"
)

func WriteSuccessResponse(w *http.ResponseWriter) {
    (*w).WriteHeader(http.StatusOK)
    fmt.Fprintf(*w, `{"success": true}`)
}

func WriteNumTicketsResponse(w *http.ResponseWriter, numTickets string) {
    (*w).WriteHeader(http.StatusOK)
    fmt.Fprintf(*w, `{"num_tickets": %s}`, numTickets)
}

func WriteTokenResponse(w *http.ResponseWriter, token string) {
    (*w).WriteHeader(http.StatusOK)
    fmt.Fprintf(*w, `{"success": true, "token": %s}`, token)
}

func WriteErrorResponse(w *http.ResponseWriter, err string) {
    (*w).WriteHeader(403)
    fmt.Fprintf(*w, `{"success": false, "errorMessage": "%s"}`, err)
}
