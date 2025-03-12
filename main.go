package main

import (
    "fmt"
    "net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "LiveKit is up and running!")
}

func main() {
    // Respond on both "/" and "/health"
    http.HandleFunc("/", healthHandler)
    http.HandleFunc("/health", healthHandler)
    fmt.Println("Starting LiveKit on port 7880")
    if err := http.ListenAndServe(":7880", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
