package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const mlServicePort = ":5555"
const geminiApiUrl = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-preview-09-2025:generateContent?key="
const maxRetries = 3

// --- Request and Response Structures based on Backend Docs ---

type MLRequest struct {
	BotID     string `json:"bot_id"`
	UserInput string `json:"user_input"`
	Goal      string `json:"goal"`
	Context   string `json:"context"`
	Webhook   string `json:"webhook"`
	// conversation_history and context are implicitly combined for the prompt
}

type MLResponse struct {
	Response    string `json:"response"`
	Action      string `json:"action"`
	WebhookData any    `json:"webhook_data"`
	EndCall     bool   `json:"end_call"`
}

// --- Gemini API Structures ---

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPayload struct {
	Contents          []GeminiContent `json:"contents"`
	SystemInstruction GeminiContent   `json:"systemInstruction"`
	Tools             []any           `json:"tools,omitempty"`
}

type GeminiResult struct {
	Candidates []struct {
		Content GeminiContent `json:"content"`
	} `json:"candidates"`
}

// --- Handlers ---

func processHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting processHandler")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Error: GEMINI_API_KEY not set")
		http.Error(w, "API Key missing", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req MLRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// 1. Construct the System Prompt
	// This defines the AI's role and rules (Polite, focused, uses provided Goal/Context)
	systemPrompt := fmt.Sprintf(`You are an AI phone agent designed to conduct conversations based on a specific goal and context.
	Your Goal: %s
	Current Context: %s
	You must keep responses concise and conversational, suitable for a voice call.
	If the conversation is complete and all necessary information is gathered, your final response must be the final confirmation.`,
		req.Goal, req.Context)

	// 2. Construct the User Query (Current input)
	userQuery := req.UserInput

	// 3. Prepare the Gemini Payload
	payload := GeminiPayload{
		Contents: []GeminiContent{
			{
				Role:  "user",
				Parts: []GeminiPart{{Text: userQuery}},
			},
		},
		SystemInstruction: GeminiContent{
			Parts: []GeminiPart{{Text: systemPrompt}},
		},
		Tools: []any{
			map[string]any{"google_search": map[string]any{}}, // Enable grounding
		},
	}

	geminiResponse, err := callGeminiAPI(r.Context(), payload, apiKey)
	if err != nil {
		log.Printf("Gemini API call failed: %v", err)
		http.Error(w, "AI service failed to generate content", http.StatusInternalServerError)
		return
	}
	log.Printf("Gemini Response Text: %s", geminiResponse) // Log what Gemini actually returned
	// 4. Determine ML Response (Simplified for MVP: always ask for more info unless the prompt is simple)
	mlRes := MLResponse{
		Response: geminiResponse,
		Action:   "continue", // Assume conversation continues
		EndCall:  false,
	}

	// Basic check to see if the conversation can end (e.g., if the user is confirming)
	if strings.Contains(strings.ToLower(userQuery), "thank you") || strings.Contains(strings.ToLower(userQuery), "goodbye") {
		mlRes.Response = "Thank you for calling! Goodbye."
		mlRes.EndCall = true
	}

	// 5. Send structured response back to main backend
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mlRes); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// callGeminiAPI handles the external API call with exponential backoff
func callGeminiAPI(ctx context.Context, payload GeminiPayload, apiKey string) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{}
	var result GeminiResult
	var lastError error

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, "POST", geminiApiUrl+apiKey, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastError = fmt.Errorf("HTTP request failed: %w", err)
			time.Sleep(time.Duration(1<<i) * time.Second) // Exponential backoff
			continue
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			lastError = fmt.Errorf("Gemini API returned status %d: %s", resp.StatusCode, string(respBody))
			time.Sleep(time.Duration(1<<i) * time.Second)
			continue
		}

		if err := json.Unmarshal(respBody, &result); err != nil {
			return "", fmt.Errorf("failed to unmarshal Gemini response: %w", err)
		}

		if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
			return result.Candidates[0].Content.Parts[0].Text, nil
		}

		return "", fmt.Errorf("gemini response was empty")
	}

	return "", lastError
}

func main() {
	log.Printf("ML Microservice starting on %s...", mlServicePort)

	http.HandleFunc("/ml/process", processHandler)

	if err := http.ListenAndServe(mlServicePort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
