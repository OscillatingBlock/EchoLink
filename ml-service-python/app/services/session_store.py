import json
import uuid
from typing import Optional
from app.core.config import settings

# Try to import redis; fallback to in-memory dict
try:
    import redis
    _has_redis = True
except Exception:
    _has_redis = False

class InMemoryStore:
    def __init__(self):
        self._store = {}

    def create_session(self, session_id: str, payload: dict, ttl: int = 1800):
        self._store[session_id] = {"payload": payload, "ttl": ttl}
    def get(self, session_id: str):
        return self._store.get(session_id)
    def delete(self, session_id: str):
        self._store.pop(session_id, None)

class RedisStore:
    def __init__(self, redis_url: str):
        self._client = redis.from_url(redis_url, decode_responses=True)

    def create_session(self, session_id: str, payload: dict, ttl: int = 1800):
        key = f"ai:session:{session_id}"
        self._client.set(key, json.dumps(payload), ex=ttl)

    def get(self, session_id: str) -> Optional[dict]:
        key = f"ai:session:{session_id}"
        data = self._client.get(key)
        return json.loads(data) if data else None

    def delete(self, session_id: str):
        key = f"ai:session:{session_id}"
        self._client.delete(key)

# Initialize store instance
if settings.REDIS_URL and _has_redis:
    store = RedisStore(settings.REDIS_URL)
else:
    store = InMemoryStore()

def new_session_id(prefix: str = "vsn") -> str:
    return f"{prefix}_{uuid.uuid4().hex[:12]}"
