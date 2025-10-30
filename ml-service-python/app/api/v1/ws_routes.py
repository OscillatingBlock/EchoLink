# api/v1/ws_routes.py
from fastapi import APIRouter, WebSocket, WebSocketDisconnect
from app.core.session_store import session_store
from app.core.stt_service import stt_service
from app.core.llm_service import llm_service
import asyncio

router = APIRouter()

@router.websocket("/session/{session_id}/ws")
async def ai_session_ws(websocket: WebSocket, session_id: str):
    session = session_store.get(session_id)
    if not session:
        await websocket.accept()
        await websocket.send_json({"type": "error", "message": "Invalid session"})
        await websocket.close(code=4000)
        return

    await websocket.accept()
    await websocket.send_json({"type": "connected", "session_id": session_id})

    # Per-session state
    audio_buffer = []
    silence_timer = None
    vad_triggered = False

    async def reset_vad():
        nonlocal silence_timer, vad_triggered
        if silence_timer:
            silence_timer.cancel()
        silence_timer = None
        vad_triggered = False

    try:
        while True:
            data = await websocket.receive()
            text_data = data.get("text")
            bytes_data = data.get("bytes")

            # === TEXT MESSAGE (fallback) ===
            if text_data:
                session_store.update_history(session_id, "user", text_data)
                history = session_store.get_history(session_id)
                reply = llm_service.generate_reply(history)
                session_store.update_history(session_id, "assistant", reply)
                await websocket.send_json({"type": "ai_text", "text": reply})
                continue

            # === AUDIO CHUNK ===
            if bytes_data:
                # Add to buffer
                audio_buffer = stt_service.add_audio_chunk(bytes_data, audio_buffer)

                # Reset VAD timer
                await reset_vad()

                # Start silence detection
                silence_timer = asyncio.create_task(
                    asyncio.sleep(0.7)  # 700ms
                )
                try:
                    await silence_timer
                    if stt_service.detect_silence(audio_buffer):
                        vad_triggered = True
                        # Trigger STT
                        await websocket.send_json({
                            "type": "partial_transcript",
                            "text": "Processing..."
                        })
                        transcript = stt_service.transcribe_buffer(audio_buffer)
                        if transcript:
                            await websocket.send_json({
                                "type": "final_transcript",
                                "text": transcript
                            })
                            # → LLM
                            session_store.update_history(session_id, "user", transcript)
                            history = session_store.get_history(session_id)
                            reply = llm_service.generate_reply(history)
                            session_store.update_history(session_id, "assistant", reply)
                            await websocket.send_json({
                                "type": "ai_text",
                                "text": reply
                            })
                        # Reset
                        audio_buffer = []
                        await reset_vad()
                except asyncio.CancelledError:
                    pass  # Timer reset by new chunk

            # === END OF INPUT (manual) ===
            end_input = data.get("text") == "end_of_input"  # or {"type": "end_of_input"}
            if end_input and audio_buffer:
                transcript = stt_service.transcribe_buffer(audio_buffer)
                if transcript:
                    await websocket.send_json({"type": "final_transcript", "text": transcript})
                    session_store.update_history(session_id, "user", transcript)
                    history = session_store.get_history(session_id)
                    reply = llm_service.generate_reply(history)
                    session_store.update_history(session_id, "assistant", reply)
                    await websocket.send_json({"type": "ai_text", "text": reply})
                audio_buffer = []

    except WebSocketDisconnect:
        print(f"WS closed: {session_id}")
    except Exception as e:
        print(f"WS error: {e}")
        await websocket.send_json({"type": "error", "message": str(e)})