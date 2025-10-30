# core/session_store.py
import uuid
from typing import Dict, Optional

class SessionStore:
    def __init__(self):
        self._sessions: Dict[str, dict] = {}

    def create(self, language: str, model: str) -> str:
        session_id = f"vsn_{uuid.uuid4().hex[:8]}"
        self._sessions[session_id] = {
            "language": language,
            "model": model,
            "history": []
        }
        return session_id

    def get(self, session_id: str) -> Optional[dict]:
        return self._sessions.get(session_id)

    def get_history(self, session_id: str):
        return self._sessions.get(session_id, {}).get("history", [])

    def update_history(self, session_id: str, role: str, text: str):
        if session_id in self._sessions:
            self._sessions[session_id]["history"].append({
                "role": role,
                "text": text
            })

session_store = SessionStore()