# 📞 EchoLink: AI Phone Agents (for hackMujX)

-----

## 💡 Project Overview

**EchoLink** is a fully functional AI Phone Agent platform built for the **hackMujX** hackathon. It allows users to connect their own Twilio phone numbers and deploy conversational AI bots using Google's Gemini models.

Crucially, **EchoLink** operates entirely using **RESTful webhooks and TwiML**, eliminating the need for complex, real-time WebSocket media streams. This approach makes deployment simple, robust, and highly scalable.

### Problem Solved

Traditional voice AI often requires managing real-time audio streams, STT/TTS services, and complex state management across network latency. EchoLink simplifies this by leveraging Twilio's built-in speech recognition (`<Gather>`) to convert speech to text, passing only **text-in/text-out** to the Gemini LLM.

## ✨ Core Features

  * **REST + TwiML Architecture:** Uses standard HTTP POST webhooks for call management, making it highly reliable.
  * **Gemini Integration:** Connects to the Gemini API via a microservice for high-quality, conversational AI processing.
  * **Twilio Integration:** Connects directly to the user's Twilio account to manage phone numbers.
  * **Stateless Conversation Management:** Maintains conversation history and context in the PostgreSQL database (`CallState` model) between Twilio webhooks.
  * **Protected API:** All user and bot management endpoints are secured using **JSON Web Tokens (JWT)**.

## 📐 System Architecture and Call Flow

The application is split into two primary services: the **Main Backend (Go/Echo)** and the **ML Microservice (Go/Gemini)**.

The call flow for every user response is a full loop:

```
sequenceDiagram
    participant C as Customer
    participant T as Twilio
    participant E as EchoLink Backend (Go:8080)
    participant M as ML Service (Go:5000)
    participant P as PostgreSQL (Bun)

    C->>T: Dials Phone Number
    T->>E: POST /v1/voice?bot_id=...
    E->>P: Retrieve Bot Config
    E->>P: Save Initial Call State
    E-->>T: TwiML: <Say> + <Gather action=/v1/voice-response>
    Note over C,T: Customer speaks ("I want to book...")
    T->>E: POST /v1/voice-response?SpeechResult=...
    E->>P: Retrieve Call History (Context)
    E->>M: POST /ml/process {input, context}
    M->>M: Call Gemini API (with System Prompt)
    M-->>E: Response Text
    E->>P: Save Updated Call State
    E-->>T: TwiML: <Say> + <Gather> (Next Turn)
    T-->>C: AI speaks the response

```

-----

## 🚀 Setup and Installation

### Prerequisites

1.  **Go:** Go environment installed.
2.  **PostgreSQL:** Database instance running and accessible.
3.  **Twilio:** Account SID, Auth Token, and a registered Phone Number SID.
4.  **Gemini API Key:** A valid API key from Google AI Studio.
5.  **ngrok:** (or similar) to expose your localhost to the internet for Twilio webhooks.

### Step 1: Configuration

Update your `config.yaml` file to point to your services and set secrets.

  * `postgres.url`: Your PostgreSQL connection string.
  * `jwt.secret`: A strong, random secret key for token signing.
  * `ml_service.url`: **Must be set to `http://localhost:5000`** during development.

### Step 2: Running the ML Microservice

This service handles the heavy lifting of the Gemini API calls.

1.  Set the API Key environment variable:
    ```bash
    export GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
    ```
2.  Navigate to the directory containing `ml_service.go` and run:
    ```bash
    go run ml_service.go 
    # Service will start on :5000
    ```

### Step 3: Running the Main Backend

1.  Ensure all dependencies are installed and the database is initialized.
2.  Run the main application:
    ```bash
    go run cmd/api/main.go 
    # Service will start on :8080
    ```

### Step 4: Testing the Live Call Flow

1.  **Start ngrok** to expose your main backend:

    ```bash
    ngrok http 8080
    ```

    (Note the public HTTPS URL provided by ngrok.)

2.  **Run the Registration/Auth Call** (Must use real credentials to fetch phone number data):

    ```bash
    curl -X POST "http://localhost:8080/v1/connect-twilio" \
    -H "Content-Type: application/json" \
    -d '{
      "first_name": "Dev", "last_name": "Hack",
      "email": "dev@hack.com",
      "account_sid": "AC...", 
      "auth_token": "xxx",
      "phone_number_sid": "PN..."
    }'
    ```

    Save the `access_token` and the `bot_id` from this response.

3.  **Configure Twilio Webhook:** In the Twilio Console, set the Voice URL for your registered phone number to:

    **`[YOUR_NGROK_URL]/v1/voice?bot_id=[YOUR_NEW_BOT_ID]`**

4.  **Call your Twilio number\!** The AI will answer, and the conversation will flow between the two microservices.
