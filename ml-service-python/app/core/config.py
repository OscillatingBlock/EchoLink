from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    APP_NAME: str = "ai-internal-service"
    HOST: str = "0.0.0.0"
    PORT: int = 8000

    # JWT settings
    JWT_SECRET: str = "change-me-to-a-secure-random-string"  # override in env
    JWT_ALGORITHM: str = "HS256"
    SESSION_TOKEN_EXPIRES_SEC: int = 30  # short-lived token for WS connect

    # Internal AI service host (where AI service listens for ws)
    AI_SERVICE_HOST: str = "ai-service"   # used to construct ws url
    AI_SERVICE_PORT: int = 8000

    # Redis
    REDIS_URL: str | None = None  # e.g. "redis://localhost:6379/0"

    class Config:
        env_file = ".env"

settings = Settings()
