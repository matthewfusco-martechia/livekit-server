package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// healthHandler responds with a simple status message.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "LiveKit is up and running!")
}

// processSpeechWithOpenAI sends the speech input to OpenAI's Chat Completions API and returns the generated response.
func processSpeechWithOpenAI(speechInput string) (string, error) {
	// OpenAI Chat Completions endpoint.
	openaiURL := "https://api.openai.com/v1/chat/completions"

	// Prepare the request body using the chat format.
	requestBody := map[string]interface{}{
		"model": "gpt-4o-mini", // Using your chosen chat model.
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": speechInput,
			},
		},
		"max_tokens":  50,
		"temperature": 0.7,
	}
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Debug: Print the HTTP status and raw response from OpenAI.
	fmt.Println("OpenAI API status:", resp.StatusCode)
	fmt.Println("OpenAI raw response:", string(responseBody))

	var openaiResponse map[string]interface{}
	if err := json.Unmarshal(responseBody, &openaiResponse); err != nil {
		return "", err
	}

	// Check if OpenAI returned an error message.
	if errVal, exists := openaiResponse["error"]; exists {
		errorJSON, _ := json.Marshal(errVal)
		return "", fmt.Errorf("OpenAI API returned error: %s", string(errorJSON))
	}

	// Extract generated text from the Chat Completions response.
	choices, ok := openaiResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in OpenAI response")
	}
	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}
	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format in response")
	}
	generatedText, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no generated text in response")
	}
	return generatedText, nil
}

// webhookHandler processes incoming events, forwards the speech input to OpenAI, and returns the result.
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

	// Extract the speech input from the payload.
	speechInput, ok := payload["speech_input"].(string)
	if !ok {
		http.Error(w, "Missing or invalid speech_input", http.StatusBadRequest)
		return
	}

	// Process the speech input with OpenAI.
	processedOutput, err := processSpeechWithOpenAI(speechInput)
	if err != nil {
		http.Error(w, "Error processing speech: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a JSON response including the processed output.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message":          "Webhook received",
		"processed_output": processedOutput,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Health check endpoints.
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/health", healthHandler)
	// Webhook endpoint.
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Starting LiveKit on port 7880")
	if err := http.ListenAndServe(":7880", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
