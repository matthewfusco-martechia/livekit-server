package main

import (
  "fmt"
  "net/http"
)

// healthHandler responds with a simple status message.
func healthHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "LiveKit is up and running!")
}

func main() {
  // Set up the health endpoint.
  http.HandleFunc("/health", healthHandler)
  fmt.Println("Starting LiveKit on port 7880")
  // Start the server on port 7880.
  err := http.ListenAndServe(":7880", nil)
  if err != nil {
    fmt.Println("Error starting server:", err)
  }
}
