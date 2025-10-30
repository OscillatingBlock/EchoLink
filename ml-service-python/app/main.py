from fastapi import FastAPI
from app.api.v1.internal_routes import router as internal_router
from app.core.config import settings
from app.api.v1.ws_routes import router as ws_router

app = FastAPI(title=settings.APP_NAME)

app.include_router(internal_router, prefix="/internal/ai")
app.include_router(ws_router, prefix="/internal/ai")