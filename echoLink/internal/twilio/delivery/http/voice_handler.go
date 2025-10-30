// internal/twilio/delivery/http/voice_handler.go
package twiliohttp

import (
	"net/http"

	"echoLink/internal/twilio/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase usecase.TwilioUsecase
}

func NewHandler(usecase usecase.TwilioUsecase) *Handler {
	return &Handler{usecase: usecase}
}

// POST /v1/voice
func (h *Handler) HandleVoice(c echo.Context) error {
	botID := c.QueryParam("bot_id")
	callSid := c.FormValue("CallSid")
	if botID == "" || callSid == "" {
		return c.XML(http.StatusBadRequest, `<Response><Say>Invalid request</Say></Response>`)
	}

	twiml, err := h.usecase.HandleVoice(c.Request().Context(), botID, callSid)
	if err != nil {
		return c.XML(http.StatusInternalServerError, `<Response><Say>Server error</Say></Response>`)
	}

	c.Response().Header().Set("Content-Type", "application/xml")
	return c.Blob(http.StatusOK, "application/xml", []byte(twiml))
}

// POST /v1/voice-response
func (h *Handler) HandleVoiceResponse(c echo.Context) error {
	botID := c.QueryParam("bot_id")
	callSid := c.FormValue("CallSid")
	speech := c.FormValue("SpeechResult")

	if botID == "" || callSid == "" || speech == "" {
		return c.XML(http.StatusBadRequest, `<Response><Say>Invalid input</Say></Response>`)
	}

	twiml, err := h.usecase.HandleVoiceResponse(c.Request().Context(), botID, callSid, speech)
	if err != nil {
		return c.XML(http.StatusInternalServerError, `<Response><Say>Server error</Say></Response>`)
	}

	c.Response().Header().Set("Content-Type", "application/xml")
	return c.Blob(http.StatusOK, "application/xml", []byte(twiml))
}
