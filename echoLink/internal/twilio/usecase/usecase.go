package usecase

import (
	"context"
	botuc "echoLink/internal/bot/usecase" // ← alias for bot usecase
	twModel "echoLink/internal/twilio/model"
	"echoLink/internal/twilio/repository"
	"echoLink/pkg/ml"
	"echoLink/pkg/twiml" // ← import generator
	"fmt"
)

type TwilioUsecase interface {
	HandleVoice(ctx context.Context, botID, callSid string) (string, error)
	HandleVoiceResponse(ctx context.Context, botID, callSid, speech string) (string, error)
}

type twilioUsecase struct {
	botRepo  botuc.BotUsecase
	callRepo repository.CallRepository
	mlClient *ml.Client
}

func NewTwilioUsecase(botRepo botuc.BotUsecase, callRepo repository.CallRepository, mlURL string) TwilioUsecase {
	return &twilioUsecase{
		botRepo:  botRepo,
		callRepo: callRepo,
		mlClient: ml.NewClient(mlURL),
	}
}

/* -------------------  /voice  ------------------- */
func (u *twilioUsecase) HandleVoice(ctx context.Context, botID, callSid string) (string, error) {
	fmt.Println("Starting voice flow")
	// FIX: We must use a method that performs a public, ID-only lookup.
	// This assumes the botuc.BotUsecase interface has a GetBotPublic method.
	bot, err := u.botRepo.GetBotPublic(ctx, botID)
	if err != nil {
		// Log the error for debugging purposes internally, but return a generic TwiML error.

		return twiml.SayGather("Bot not found", ""), nil
	}
	fmt.Printf("bot %v", bot.ID)
	fmt.Printf("bot %v", bot.Goal)

	// initialise call state
	state := &twModel.CallState{
		CallSid: callSid,
		BotID:   botID,
		Context: bot.Prompt(),
		Step:    0,
	}
	// Note: We ignore the error here for MVP stability
	_ = u.callRepo.Upsert(ctx, state)

	actionURL := "/v1/voice-response?bot_id=" + botID
	return twiml.SayGather("Hi, how can I help you?", actionURL), nil
}

/* -------------------  /voice-response  ------------------- */
// func (u *twilioUsecase) HandleVoiceResponse(ctx context.Context, botID, callSid, speech string) (string, error) {
// 	// FIX: Use the public lookup method here as well.
// 	bot, err := u.botRepo.GetBotPublic(ctx, botID)
// 	if err != nil {
// 		return twiml.SayGather("Bot not found", ""), nil
// 	}

// 	// get (or create) call state
// 	state, _ := u.callRepo.Get(ctx, callSid)
// 	if state == nil {
// 		state = &twModel.CallState{CallSid: callSid, BotID: botID, Context: bot.Prompt()}
// 	}

// 	state.UserInput = speech
// 	state.Step++

// 	// ---- call Python ML ----
// 	mlRes, err := u.mlClient.Process(ml.Request{
// 		BotID:     botID,
// 		UserInput: speech,
// 		Goal:      bot.Goal,
// 		Context:   state.Context,
// 		Webhook:   bot.Webhook,
// 		/* History:   []ml.Message{}, // can be filled later */
// 	})
// 	if err != nil {
// 		return twiml.SayGather("Sorry, the AI is unavailable.", ""), nil
// 	}

// 	// update state
// 	state.AssistantResp = mlRes.Response
// 	state.Context += "\nUser: " + speech + "\nAssistant: " + mlRes.Response
// 	_ = u.callRepo.Upsert(ctx, state)

// 	if mlRes.EndCall {
// 		return twiml.SayAndHangup(mlRes.Response), nil
// 	}

// 	actionURL := "/v1/voice-response?bot_id=" + botID
// 	return twiml.SayGather("hello", actionURL), nil
// }
/* -------------------  /voice-response  ------------------- */
func (u *twilioUsecase) HandleVoiceResponse(ctx context.Context, botID, callSid, speech string) (string, error) {
	// FIX: Use the public lookup method here as well.
	bot, err := u.botRepo.GetBotPublic(ctx, botID)
	if err != nil {
		return twiml.SayGather("Bot not found", ""), nil
	}

	// get (or create) call state
	state, _ := u.callRepo.Get(ctx, callSid)
	if state == nil {
		state = &twModel.CallState{CallSid: callSid, BotID: botID, Context: bot.Prompt()}
	}

	state.UserInput = speech
	state.Step++

	// ---- call Python ML ----
	mlRes, err := u.mlClient.Process(ml.Request{
		BotID:     botID,
		UserInput: speech,
		Goal:      bot.Goal,
		Context:   state.Context,
		Webhook:   bot.Webhook,
		/* History: []ml.Message{}, // can be filled later */
	})
	if err != nil {
		return twiml.SayGather("Sorry, the AI is unavailable.", ""), nil
	}

	// update state
	state.AssistantResp = mlRes.Response
	state.Context += "\nUser: " + speech + "\nAssistant: " + mlRes.Response
	_ = u.callRepo.Upsert(ctx, state)

	if mlRes.EndCall {
		return twiml.SayAndHangup(mlRes.Response), nil
	}

	actionURL := "/v1/voice-response?bot_id=" + botID
	return twiml.SayGather(mlRes.Response, actionURL), nil
}
