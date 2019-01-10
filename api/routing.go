package main

import (
    "fmt"
    "net/http"
)

func WriteErrorResponse(w *http.ResponseWriter, err string) {
    (*w).WriteHeader(403)
    fmt.Fprintf(*w, `{"success": false, "errorMessage": "%s"}`, err)
}

func BasicSuccessResponse(w *http.ResponseWriter) {
    (*w).WriteHeader(200)
    fmt.Fprintf(*w, `{"success": true}`)
}
