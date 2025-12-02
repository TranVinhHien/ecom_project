from __future__ import annotations

from typing import Any, Dict, List, Optional

import numpy as np
import redis
import torch
from redis.commands.search.field import TagField, TextField, VectorField
from redis.commands.search.indexDefinition import IndexDefinition, IndexType
from redis.commands.search.query import Query
from redis.exceptions import ResponseError
from sentence_transformers import SentenceTransformer
from .call_api import get_product_detail_for_search
from .config import settings


class EmbeddingService:
    """Loads the embedding model and produces GPU-backed sentence embeddings."""

    def __init__(self, model_name: str) -> None:
        self._device = self._resolve_device()
        self._model = SentenceTransformer(
            model_name,
            device=self._device,
            trust_remote_code=True,
        )
        self._model.eval()

    @staticmethod
    def _resolve_device() -> str:
        if not torch.cuda.is_available():
            raise RuntimeError("CUDA-enabled GPU is required but not detected.")
        return "cuda"

    def embed(self, text: str) -> np.ndarray:
        vector = self._model.encode(
            [text],
            device=self._device,
            convert_to_numpy=True,
            normalize_embeddings=True,
            batch_size=1,
            show_progress_bar=False,
        )[0]
        return np.asarray(vector, dtype=np.float32)


class VectorDBService:
    """Handles Redis vector index creation and CRUD/search operations."""

    def __init__(self) -> None:
        self._redis = redis.Redis(
            host=settings.REDIS_HOST,
            port=settings.REDIS_PORT,
            decode_responses=False,
        )
        self._index_name = settings.INDEX_PRODUCT
        self._index_policy = settings.INDEX_POLICY
        self._vector_dim = settings.VECTOR_DIM
        self._doc_prefix = settings.DOC_PREFIX
        self._policy_prefix = settings.POLICY_PREFIX

    def create_index_if_not_exists(self) -> None:
        index = self._redis.ft(self._index_name)
        index_policy = self._redis.ft(self._index_policy)
        # index_policy.
        try:
            info = index.info()
            info_policy = index_policy.info()
            print(f"[INFO] Index '{self._index_name}' already exists with {info.get('num_docs', 0)} documents.")
            print(f"[INFO] Index '{self._index_policy}' already exists with {info_policy.get('num_docs', 0)} documents.")
            return
        except ResponseError as exc:
            message = str(exc).lower()
            # Check if error is about missing index (both patterns)
            if "unknown index name" not in message and "no such index" not in message:
                print(f"[ERROR] Unexpected Redis error while checking index: {exc}")
                raise
            print(f"[INFO] Index '{self._index_name}' not found. Creating...")

        schema = [
            TagField("doc_id"),
            TextField("text_content"),
            VectorField(
                "vector",
                "HNSW",
                {
                    "TYPE": "FLOAT32",
                    "DIM": self._vector_dim,
                    "DISTANCE_METRIC": "COSINE",
                    "INITIAL_CAP": 2000,
                    "M": 16,
                    "EF_CONSTRUCTION": 200,
                },
            ),
        ]
        definition = IndexDefinition(prefix=[self._doc_prefix], index_type=IndexType.HASH)
        definition_policy = IndexDefinition(prefix=[self._policy_prefix], index_type=IndexType.HASH)
        try:
            index.create_index(fields=schema, definition=definition)
            index_policy.create_index(fields=schema, definition=definition_policy)
            print(f"[INFO] Index '{self._index_name}' created successfully with prefix '{self._doc_prefix}'.")
            print(f"[INFO] Index '{self._index_policy}' created successfully with prefix '{self._policy_prefix}'.")
        except Exception as e:
            print(f"[ERROR] Failed to create index: {e}")
            raise

    def add_document(self, doc_type: str,doc_id: str,  text_content: str, vector: np.ndarray) -> None:
        key = f"{doc_type}:{doc_id}"
        mapping = {
            "doc_id": doc_id,
            "text_content": text_content,
            "vector": vector.astype(np.float32, copy=False).tobytes(),
        }
        self._redis.hset(name=key, mapping=mapping)
        print(f"[INFO] Document added: key={key}, doc_id={doc_id}, text_length={len(text_content)}")

    def update_document(self,doc_type: str, doc_id: str, text_content: str, vector: np.ndarray) -> None:
        key = f"{doc_type}:{doc_id}"
        mapping = {
            "text_content": text_content,
            "vector": vector.astype(np.float32, copy=False).tobytes(),
        }
        self._redis.hset(name=key, mapping=mapping)

    def delete_document(self,doc_type: str,doc_id: str) -> None:
        key = f"{doc_type}:{doc_id}"
        self._redis.delete(key)

    def search_documents(self, query_vector: np.ndarray, top_k: int, doc_type : str ) -> Optional[Dict[str, Any]]:
        index = None
        if doc_type == "product":
            index = self._redis.ft(self._index_name)
        elif doc_type == "policy":
            index = self._redis.ft(self._index_policy)
        if index is None:
            return None
        # Ensure index exists before searching
        try:
            index.info()
        except ResponseError as exc:
            message = str(exc).lower()
            if "unknown index name" in message or "no such index" in message:
                print(f"[WARNING] Index '{self._index_name}' does not exist. Creating it now...")
                self.create_index_if_not_exists()
                print(f"[INFO] Index created, but no documents to search yet. Returning empty result.")
                return {}
            raise
        
        knn_clause = f"*=>[KNN {top_k} @vector $query_vec AS vector_score]"
        query = (
            Query(knn_clause)
            .return_fields("doc_id", "text_content", "vector_score")
            .sort_by("vector_score", asc=True)
            .paging(0, top_k)
            .dialect(2)
        )
        params = {"query_vec": query_vector.astype(np.float32, copy=False).tobytes()}
        print(f"[INFO] Executing search: index={doc_type}, top_k={top_k}")
        result = index.search(query, query_params=params)
        print(f"[INFO] Search returned {result.total} results")
        trave: List[str] = []
        filtered_docs_with_scores: Dict[str, float] = {}
        for doc in getattr(result, "docs", []):
            doc_id = self._ensure_str(getattr(doc, "doc_id", ""))
            text_content = self._ensure_str(getattr(doc, "text_content", ""))
            trave.append(text_content)
            score_raw = float(getattr(doc, "vector_score", 0.0))
            similarity = max(0.0, 1.0 - score_raw)
            if similarity > 0.45:
                filtered_docs_with_scores[doc_id] = similarity
        data = get_product_detail_for_search(filtered_docs_with_scores) if doc_type == "product" else trave
        return data

    def close(self) -> None:
        self._redis.close()
        self._redis.connection_pool.disconnect()

    @staticmethod
    def _ensure_str(value: Any) -> str:
        if value is None:
            return ""
        if isinstance(value, bytes):
            return value.decode("utf-8")
        return str(value)
