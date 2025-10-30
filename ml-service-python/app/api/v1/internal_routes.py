# api/v1/internal_routes.py
from fastapi import APIRouter
from app.core.session_store import session_store
from app.core.config import settings

router = APIRouter()

@router.post("/session/start")
async def start_session(language: str = "en", model: str = "default"):
    session_id = session_store.create(language, model)
    ws_url = f"ws://{settings.AI_SERVICE_HOST}:{settings.AI_SERVICE_PORT}/session/{session_id}/ws"
    return {
        "session_id": session_id,
        "ws_url": ws_url
    }