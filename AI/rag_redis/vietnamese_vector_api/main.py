from contextlib import asynccontextmanager
from fastapi import Depends, FastAPI, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.concurrency import run_in_threadpool
from typing import Any, Dict, List, Optional

from .config import settings
from .schemas import (
    DocumentIn,
    DocumentResponse,
    DocumentUpdate,
    SearchQuery,
    SearchResponse,
)
from .services import EmbeddingService, VectorDBService


@asynccontextmanager
async def lifespan(app: FastAPI):  # pragma: no cover - executed on app startup
    embedding_service = EmbeddingService(settings.EMBEDDING_MODEL_NAME)
    vector_service = VectorDBService()
    vector_service.create_index_if_not_exists()

    app.state.embedding_service = embedding_service
    app.state.vector_service = vector_service

    try:
        yield
    finally:
        vector_service.close()
        app.state.embedding_service = None
        app.state.vector_service = None


app = FastAPI(lifespan=lifespan, title="Vietnamese Vector Search API")

# CORS configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Cho phép tất cả origins. Trong production nên chỉ định cụ thể
    allow_credentials=True,
    allow_methods=["*"],  # Cho phép tất cả HTTP methods
    allow_headers=["*"],  # Cho phép tất cả headers
)


def get_embedding_service(request: Request) -> EmbeddingService:
    service = getattr(request.app.state, "embedding_service", None)
    if not isinstance(service, EmbeddingService):
        HTTPException(status_code=500, detail="embedding_service unavailable")
        return
    return service


def get_vector_db_service(request: Request) -> VectorDBService:
    service = getattr(request.app.state, "vector_service", None)
    if not isinstance(service, VectorDBService):
        HTTPException(status_code=500, detail="vector_service unavailable")
        return
    return service


@app.post("/documents", response_model=DocumentResponse)
async def add_document(
    payload: DocumentIn,
    embed_service: EmbeddingService = Depends(get_embedding_service),
    vdb_service: VectorDBService = Depends(get_vector_db_service),
) -> DocumentResponse:
    if payload.doc_type not in ["product", "policy"]:
        HTTPException(status_code=400, detail="doc_type must be 'product' or 'policy'")
        return
    vector = await run_in_threadpool(embed_service.embed, payload.text_content)
    await run_in_threadpool(
        vdb_service.add_document,
        payload.doc_type,
        payload.doc_id,
        payload.text_content,
        vector,
    )
    return DocumentResponse(status="success", doc_id=payload.doc_id)


@app.put("/documents/{doc_id}", response_model=DocumentResponse)
async def update_document(
    doc_id: str,
    payload: DocumentUpdate,
    embed_service: EmbeddingService = Depends(get_embedding_service),
    vdb_service: VectorDBService = Depends(get_vector_db_service),
) -> DocumentResponse:
    if payload.doc_type not in ["product", "policy"]:
        HTTPException(status_code=400, detail="doc_type must be 'product' or 'policy'")
        return

    vector = await run_in_threadpool(embed_service.embed, payload.text_content)
    
    await run_in_threadpool(
        vdb_service.update_document,
        payload.doc_type,
        doc_id,
        payload.text_content,
        vector,
    )
    return DocumentResponse(status="updated", doc_id=doc_id)


@app.delete("/documents/{doc_id}/{doc_type}", response_model=DocumentResponse)
async def delete_document(
    doc_id: str,
    doc_type: str,
    vdb_service: VectorDBService = Depends(get_vector_db_service),
) -> DocumentResponse:
    if doc_type not in ["product", "policy"]:
        HTTPException(status_code=400, detail="doc_type must be 'product' or 'policy'")
        return

    await run_in_threadpool(vdb_service.delete_document,doc_type, doc_id)
    return DocumentResponse(status="deleted", doc_id=doc_id)


@app.post("/search", response_model=Dict[str, Any] | None)
async def search_documents(
    query: SearchQuery,
    embed_service: EmbeddingService = Depends(get_embedding_service),
    vdb_service: VectorDBService = Depends(get_vector_db_service),
) -> Dict[str, Any] | None:
    if query.doc_type not in ["product", "policy"]:
        raise HTTPException(status_code=400, detail="doc_type must be 'product', 'policy' or empty")
    query_vector = await run_in_threadpool(embed_service.embed, query.query_text)
    results = await run_in_threadpool(vdb_service.search_documents, query_vector, query.top_k, query.doc_type)
    return {"results": results}
