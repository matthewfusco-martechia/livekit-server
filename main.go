package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// healthHandler responds with a simple status message.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "LiveKit is up and running!")
}

// webhookHandler processes incoming realtime integration events.
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	fmt.Println("Webhook received:", payload)

	// TODO: Here you can integrate with the OpenAI API as per the realtime integration docs.
	// For now, simply respond with a success message.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received")
}

func main() {
	// Respond on both "/" and "/health" for general health checks.
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/health", healthHandler)
	// Set up the webhook endpoint.
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Starting LiveKit on port 7880")
	if err := http.ListenAndServe(":7880", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
