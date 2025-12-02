"""
ProductDetailAgent - Agent chuyên xử lý lấy chi tiết sản phẩm và phân tích đánh giá
"""
import json
import logging
from typing import Dict, Any, Optional

import httpx
from google.adk import Agent
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.adk.models.lite_llm import LiteLlm
from google.genai import types
from pydantic import BaseModel, Field
from typing import List
from google.adk.tools.tool_context import ToolContext

from .product_detail_promtp import root_instruction


class ProductDetailResponse(BaseModel):
    """
    Cấu trúc dữ liệu chuẩn hóa cho việc trả về phân tích sản phẩm
    """
    response_text: str = Field(description="Nội dung phân tích chi tiết về sản phẩm và đánh giá")
    product_key: str = Field(default="", description="Product key được phân tích")


class ProductDetailAgent:
    """Agent chuyên xử lý lấy chi tiết sản phẩm và phân tích comments"""
    
    BASE_PRODUCT_URL = "http://172.26.127.95:9001/v1"  # TODO: Cập nhật URL API thực tế
    BASE_COMMENT_URL = "http://172.26.127.95:9002/v1"  # TODO: Cập nhật URL API thực tế
    
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
        """Tạo Agent ProductDetail với cấu hình riêng"""
        return Agent(
            model=self._llm_model if self._llm_model.startswith("gemini") else LiteLlm(model=self._llm_model),
            name="agent_product_detail",
            instruction=root_instruction(),
            description="Agent chuyên lấy chi tiết sản phẩm và phân tích đánh giá",
            tools=[
                self.get_product_detail,
            ],
            include_contents="none",
            output_schema=ProductDetailResponse,
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
    
    async def process_product_query(
        self,
        query: str,
        user_id: str,
        token: str,
        session_id: str,
        product_key: str
    ) -> Dict[str, Any]:
        """
        Xử lý câu hỏi về chi tiết sản phẩm từ người dùng
        
        Args:
            query: Câu hỏi/tên sản phẩm/product_key của người dùng
            user_id: ID người dùng
            token: JWT token để xác thực
            session_id: ID phiên làm việc
            
        Returns:
            Dict chứa phản hồi phân tích sản phẩm
        """
        logging.info("ProductDetailAgent đang xử lý câu hỏi: %s", query)
        
        try:
            # Gọi Agent ProductDetail để xử lý
            content = types.Content(role="user", parts=[types.Part(text=query)])
            session_data = await self._ensure_session(user_id, session_id)
            
            # Import process_agent_response từ util
            from ..util import process_agent_response
            
            final_response = None
            async for event in self.runner.run_async(
                user_id=session_data.user_id,
                session_id=session_data.id,
                new_message=content,
                state_delta={"token": token, "product_key": product_key}
            ):
                response = await process_agent_response(event)
                if response:
                    parsed = json.loads(response.get("text", "{}"))
                    response_text = parsed.get("response_text", "")
                    product_key = parsed.get("product_key", "")
                    
                    final_response = {
                        "text": response_text,
                        "product_key": product_key
                    }
            
            return final_response or {"text": "Không nhận được phản hồi từ agent"}
            
        except Exception as e:
            logging.exception("Lỗi trong ProductDetailAgent.process_product_query")
            return {
                "error": f"Lỗi khi xử lý yêu cầu: {str(e)}",
                "text": "Tôi gặp lỗi khi lấy chi tiết sản phẩm, vui lòng thử lại."
            }
    
    async def get_product_detail(
        self,
        tool_context: ToolContext,
        product_key: str
    ) -> Dict[str, Any]:
        """
        Tool: Lấy chi tiết sản phẩm kèm comments/reviews
        
        Args:
            tool_context: Context từ ADK
            product_key: Key của sản phẩm (VD: "android-tivi-box-ram-2g-android-tv-10-dual-wifi-bluetooth-netflix-remote-tim-kiem-bang-g")
        
        Returns:
            Dict chứa thông tin sản phẩm đã lọc và comments được nhóm theo rating
        """
        product_key_context = tool_context.state.get("product_key", "unknown")
        if product_key_context != "unknown":
            product_key = product_key_context
            
        logging.info("get_product_detail được gọi: product_key=%s", product_key)
        
        token = tool_context.state.get("token", "unknown")
        if token == "unknown":
            logging.warning("Token không được cung cấp trong tool_context.state")
            return {"error": "Authorization token is missing"}
        
        async with httpx.AsyncClient() as client:
            try:
                # BƯỚC 1: Lấy chi tiết sản phẩm
                product_response = await client.get(
                    f"{self.BASE_PRODUCT_URL}/product/getdetail/{product_key}",
                    headers={"Authorization": f"Bearer {token}"},
                    timeout=20.0
                )
                product_response.raise_for_status()
                product_result = product_response.json()
                
                if product_result.get("code") != 200:
                    return {"error": f"API trả về lỗi: {product_result.get('message', 'Unknown error')}"}
                
                raw_data = product_result.get("result", {}).get("data", {})
                if not raw_data:
                    return {"error": f"Không tìm thấy sản phẩm với key: {product_key}"}
                
                # Lọc dữ liệu sản phẩm - chỉ lấy thông tin cần thiết
                filtered_product = {
                    "brand": raw_data.get("brand", {}).get("name", ""),
                    "category": raw_data.get("category", {}).get("name", ""),
                    "options": [
                        {
                            "option_name": opt.get("option_name", ""),
                            "values": [v.get("value", "") for v in opt.get("values", [])]
                        }
                        for opt in raw_data.get("option", [])
                    ],
                    "product": {
                        "id": raw_data.get("product", {}).get("id", ""),
                        "name": raw_data.get("product", {}).get("name", ""),
                        "key": raw_data.get("product", {}).get("key", ""),
                        "description": raw_data.get("product", {}).get("description", ""),
                        "short_description": raw_data.get("product", {}).get("short_description", ""),
                        "image": raw_data.get("product", {}).get("image", ""),
                        "min_price": raw_data.get("product", {}).get("min_price", 0),
                        "max_price": raw_data.get("product", {}).get("max_price", 0),
                    },
                    "sku": [
                        {
                            "sku_name": sku.get("sku_name", ""),
                            "price": sku.get("price", 0),
                            "quantity": sku.get("quantity", 0)
                        }
                        for sku in raw_data.get("sku", [])
                    ]
                }
                
                product_id = filtered_product["product"]["id"]
                if not product_id:
                    return {
                        "product": filtered_product,
                        "comments": {
                            "data": [],
                            "limit": 0,
                            "totalElements": 0
                        },
                        "message": "Sản phẩm không có ID, không thể lấy comments"
                    }
                
                # BƯỚC 2: Lấy comments của sản phẩm
                comment_response = await client.get(
                    f"{self.BASE_COMMENT_URL}/comments",
                    params={
                        "product_id": product_id,
                        "page": 1,
                        "page_size": 20
                    },
                    headers={"Authorization": f"Bearer {token}"},
                    timeout=15.0
                )
                comment_response.raise_for_status()
                comment_result = comment_response.json()
                
                if comment_result.get("code") != 200:
                    logging.warning("Không lấy được comments: %s", comment_result.get("message"))
                    comments_data = []
                    total_elements = 0
                    limit = 0
                else:
                    comments_data = comment_result.get("result", {}).get("data", [])
                    total_elements = comment_result.get("result", {}).get("totalElements", 0)
                    limit = comment_result.get("result", {}).get("limit", 0)
                
                # BƯỚC 3: Nhóm comments theo rating (1-5 sao)
                grouped_comments = {1: [], 2: [], 3: [], 4: [], 5: []}
                
                for comment in comments_data:
                    rating = comment.get("rating", 0)
                    content = comment.get("content", "").strip()
                    
                    if rating in grouped_comments and content:
                        grouped_comments[rating].append(content)
                
                # Tạo cấu trúc output cho comments
                comments_output = {
                    "data": [
                        {
                            "star": star,
                            "count": len(comments),
                            "comments": comments
                        }
                        for star, comments in sorted(grouped_comments.items())
                        if len(comments) > 0  # Chỉ hiển thị rating có comments
                    ],
                    "limit": limit,
                    "totalElements": total_elements
                }
                
                # Trả về kết quả cuối cùng
                return {
                    "product": filtered_product,
                    "comments": comments_output,
                    "total_comments": total_elements
                }
                
            except httpx.HTTPStatusError as e:
                logging.error("HTTP error khi lấy chi tiết sản phẩm: %s", e)
                if e.response.status_code == 404:
                    return {"error": f"Không tìm thấy sản phẩm với key: {product_key}"}
                return {"error": f"Lỗi API: {e.response.status_code} - {e.response.text}"}
            except httpx.TimeoutException:
                logging.error("Timeout khi gọi API")
                return {"error": "Timeout khi lấy thông tin sản phẩm"}
            except Exception as e:
                logging.exception("Exception khi lấy chi tiết sản phẩm")
                return {"error": f"Lỗi hệ thống: {str(e)}"}
