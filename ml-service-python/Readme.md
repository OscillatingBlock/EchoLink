

## File Tree

fastapi-app/
├─ app/
│  ├─ __init__.py
│  ├─ main.py
│  ├─ core/
│  │  ├─ __init__.py
│  │  └─ config.py
│  ├─ api/
│  │  ├─ __init__.py
│  │  └─ v1/
│  │     ├─ __init__.py
│  │     └─ internal_routes.py
│  ├─ services/
│  │  ├─ __init__.py
│  │  └─ session_store.py
│  └─ utils/
│     ├─ __init__.py
│     └─ auth.py
├─ requirements.txt
└─ tests/
   └─ test_internal_start.py


## TODOs

| Phase | Component                    | Tech                           | Status            |
| ----- | ---------------------------- | ------------------------------ | ----------------- |
| 1     | WebSocket base communication | FastAPI                        | ✅ Done            |
| 2     | Audio Streaming & STT        | Whisper (local or lightweight) | 🔥 Next           |
| 3     | Triage Orchestration         | Rule-based + NER to start      | ⏳ Pending         |
| 4     | LLM Response Generation      | Mistral (local later)          | ⏳ Pending         |
| 5     | TTS                          | Coqui TTS or Edge TTS          | ⏳ Optional next   |
| 6     | Deploy GPU Model             | Cloud server                   | ✅ Supported later |
