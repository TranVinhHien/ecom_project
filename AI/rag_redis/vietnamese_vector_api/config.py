from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    REDIS_HOST: str = "172.26.127.95"
    REDIS_PORT: int = 6379
    EMBEDDING_MODEL_NAME: str = "dangvantuan/vietnamese-document-embedding"
    VECTOR_DIM: int = 768
    INDEX_PRODUCT: str = "document_index"
    INDEX_POLICY: str = "policy_index"
    DOC_PREFIX: str = "product:"
    POLICY_PREFIX: str = "policy:"
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = False


settings = Settings()
