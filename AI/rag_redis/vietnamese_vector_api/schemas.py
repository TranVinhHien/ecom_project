from typing import Any, Dict, List, Optional


from pydantic import BaseModel, Field


class DocumentIn(BaseModel):
    doc_id: str = Field(..., min_length=1)
    doc_type: str = Field(..., min_length=1, description="product or policy")
    text_content: str = Field(..., min_length=1)


class DocumentUpdate(BaseModel):
    text_content: str = Field(..., min_length=1)
    doc_type: str = Field(..., min_length=1, description="product or policy")



class DocumentResponse(BaseModel):
    status: str
    doc_id: str


class SearchQuery(BaseModel):
    query_text: str = Field(..., min_length=1)
    top_k: int = Field(5, ge=1, le=50)
    doc_type: str = Field("", description="product or polices")


# class SearchResultItem(BaseModel):
#     doc_id: str
#     text_content: str
#     score: float


class SearchResponse(BaseModel):
    results:  Dict[str, Any] | None
