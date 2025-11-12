"""
RAGAgent - Agent chuyên xử lý RAG (Retrieval-Augmented Generation)
"""
import json
import logging
from typing import Dict, Any

import httpx
from google.adk import Agent
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.adk.models.lite_llm import LiteLlm
from google.genai import types
from pydantic import BaseModel, Field
from typing import List

from host.agent.rag_promtp import AGENT_RAG_INSTRUCTION


class ProductAnalysisResponse(BaseModel):
    """
    Cấu trúc dữ liệu chuẩn hóa cho việc phân tích sản phẩm.
    """
    analysis_text: str = Field(description="Đoạn văn bản chi tiết phân tích và tư vấn sản phẩm cho người dùng.")
    product_ids: List[str] = Field(
        default_factory=list,
        description="Danh sách (mảng) các ID của sản phẩm đã được sử dụng trong phân tích."
    )


class RAGAgent:
    """Agent chuyên xử lý RAG (Retrieval-Augmented Generation) cho sản phẩm và chính sách"""
    
    def __init__(self, llm_model: str, session_service: InMemorySessionService):
        self._llm_model = llm_model
        self._session_service = session_service
        self._agent = self._create_agent()
        self.runner = Runner(
            app_name=self._agent.name,
            agent=self._agent,
            session_service=self._session_service,
        )
    
    def _create_agent(self) -> Agent:
        """Tạo Agent RAG với cấu hình riêng"""
        return Agent(
            model=self._llm_model if self._llm_model.startswith("gemini") else LiteLlm(model=self._llm_model),
            name="agent_rag",
            instruction=AGENT_RAG_INSTRUCTION,
            description="Bạn là 1 người tư vấn viên tư vấn sản phẩm cho khách hàng.",
            include_contents="none",
            output_schema=ProductAnalysisResponse,
            disallow_transfer_to_parent=True,
            disallow_transfer_to_peers=True
        )
    
    async def _ensure_session(self, user_id: str, session_id: str):
        """Đảm bảo session tồn tại, tạo mới nếu chưa có"""
        session_data = await self._session_service.get_session(
            user_id=user_id,
            session_id=session_id,
            app_name=self._agent.name,
        )
        if not session_data:
            session_data = await self._session_service.create_session(
                app_name=self._agent.name,
                user_id=user_id,
                session_id=session_id,
                state={}
            )
        return session_data
    
    async def search_and_analyze(
        self, 
        query: str, 
        intent: str,
        user_id: str,
        session_id: str
    ) -> Dict[str, Any]:
        """
        Tìm kiếm và phân tích dữ liệu RAG
        
        Args:
            query: Câu hỏi của người dùng
            intent: "product" hoặc "policy"
            user_id: ID người dùng
            session_id: ID phiên làm việc
            
        Returns:
            Dict chứa kết quả phân tích và sản phẩm khớp
        """
        logging.info("RAGAgent đang xử lý với intent: %s", intent)
        
        try:
            # 1. Lấy dữ liệu từ RAG
            if intent == "product":
                raw_data = await self._search_products_rag(query)
            elif intent == "policy":
                raw_data = await self._search_policies_rag(query)
            else:
                return {"error": f"Intent '{intent}' không hợp lệ. Chỉ chấp nhận 'product' hoặc 'policy'."}  
            if not raw_data or "error" in raw_data:
                return raw_data
            
            # 2. Chuẩn bị context cho LLM phân tích
            try:
                context_json = json.dumps(raw_data, ensure_ascii=False, indent=2)
            except TypeError:
                if hasattr(raw_data, 'model_dump_json'):
                    context_json = raw_data.model_dump_json(indent=2)
                elif hasattr(raw_data, 'dict'):
                    context_json = json.dumps(raw_data.dict(), ensure_ascii=False, indent=2)
                else:
                    context_json = str(raw_data)
            
            final_prompt = f"""
[BỐI CẢNH SẢN PHẨM]
{context_json}

[CÂU HỎI CỦA NGƯỜI DÙNG] 
{query}

[PHÂN TÍCH VÀ TRẢ LỜI]
"""
            
            # 3. Gọi Agent RAG phân tích
            content = types.Content(role="user", parts=[types.Part(text=final_prompt)])
            session_data = await self._ensure_session(user_id, session_id)
            
            # Import process_agent_response từ util
            from ..util import process_agent_response
            
            final_response = None
            async for event in self.runner.run_async(
                user_id=session_data.user_id,
                session_id=session_data.id,
                new_message=content
            ):
                response = await process_agent_response(event)
                if response:
                    parsed = json.loads(response.get("text", "{}"))
                    analysis_text = parsed.get("analysis_text", "")
                    product_ids = parsed.get("product_ids", [])
                    
                    # Lọc sản phẩm khớp
                    matched_products = [
                        p for p in raw_data.get("products", [])
                        if isinstance(p, dict) and (p.get("product") or {}).get("id") in product_ids
                    ]
                    
                    final_response = {
                        "text": analysis_text,
                        "matched_products": matched_products
                    }
            
            return final_response or {"text": "Không nhận được phản hồi từ agent phân tích"}
            
        except Exception as e:
            logging.exception("Lỗi trong RAGAgent.search_and_analyze")
            return {
                "error": f"Lỗi khi xử lý RAG: {str(e)}",
                "text": "Tôi gặp lỗi khi phân tích dữ liệu, vui lòng thử lại."
            }
    
    async def _search_products_rag(self, query: str) -> Dict[str, Any]:
        """Tìm kiếm sản phẩm qua API RAG"""
        async with httpx.AsyncClient() as client:
            try:
                response = await client.post(
                    "http://localhost:8000/search",
                    json={"query_text": query, "top_k": 5, "doc_type": "product"},
                    timeout=10.0
                )
                response.raise_for_status()
                results = response.json().get("results", [])
                
                if not results:
                    return {"message": "Không tìm thấy sản phẩm nào phù hợp."}
                
                return {"products": results}
                
            except httpx.HTTPStatusError as e:
                logging.error("HTTP error khi tìm sản phẩm: %s", e)
                return {"error": f"Lỗi API: {e.response.status_code}"}
            except Exception as e:
                logging.exception("Exception khi tìm sản phẩm")
                return {"error": f"Lỗi hệ thống: {str(e)}"}
    
    async def _search_policies_rag(self, query: str) -> Dict[str, Any]:
        """Tìm kiếm chính sách qua API RAG"""
        async with httpx.AsyncClient() as client:
            try:
                response = await client.post(
                    "http://localhost:8000/search",
                    json={"query_text": query, "top_k": 5, "doc_type": "policy"},
                    timeout=10.0
                )
                response.raise_for_status()
                results = response.json().get("results", [])
                
                if not results:
                    return {"message": "Không tìm thấy chính sách nào phù hợp."}
                
                return {"policies": results}
                
            except httpx.HTTPStatusError as e:
                logging.error("HTTP error khi tìm chính sách: %s", e)
                return {"error": f"Lỗi API: {e.response.status_code}"}
            except Exception as e:
                logging.exception("Exception khi tìm chính sách")
                return {"error": f"Lỗi hệ thống: {str(e)}"}
