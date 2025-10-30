from datetime import datetime, timedelta
from jose import jwt
from app.core.config import settings

def create_session_token(session_id: str, expires_in: int | None = None) -> str:
    if expires_in is None:
        expires_in = settings.SESSION_TOKEN_EXPIRES_SEC
    now = datetime.utcnow()
    payload = {
        "sub": "ai_session",
        "session_id": session_id,
        "iat": int(now.timestamp()),
        "exp": int((now + timedelta(seconds=expires_in)).timestamp()),
        "aud": "ai-service"
    }
    token = jwt.encode(payload, settings.JWT_SECRET, algorithm=settings.JWT_ALGORITHM)
    return token
