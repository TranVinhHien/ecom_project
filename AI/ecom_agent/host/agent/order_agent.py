"""
OrderAgent - Agent chuyên xử lý tra cứu và quản lý đơn hàng
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

from .order_promtp import root_intruction


class OrderResponse(BaseModel):
    """
    Cấu trúc dữ liệu chuẩn hóa cho việc trả về thông tin đơn hàng
    """
    response_text: str = Field(description="Nội dung phản hồi chi tiết về đơn hàng cho người dùng")
    order_ids: List[str] = Field(
        default_factory=list,
        description="Danh sách ID các đơn hàng được đề cập trong phản hồi"
    )


class OrderAgent:
    """Agent chuyên xử lý tra cứu và quản lý thông tin đơn hàng"""
    
    BASE_URL = "http://172.26.127.95:9002/v1"
    
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
        """Tạo Agent Order với cấu hình riêng"""
        return Agent(
            model=self._llm_model if self._llm_model.startswith("gemini") else LiteLlm(model=self._llm_model),
            name="agent_order",
            instruction=root_intruction(),
            description="Agent chuyên về tra cứu và quản lý đơn hàng cho khách hàng",
            tools=[
                self.search_orders_detail,
            ],
            include_contents="none",
            output_schema=OrderResponse,
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
    
    async def process_order_query(
        self,
        query: str,
        user_id: str,
        token: str,
        session_id: str
    ) -> Dict[str, Any]:
        """
        Xử lý câu hỏi về đơn hàng từ người dùng
        
        Args:
            query: Câu hỏi của người dùng về đơn hàng
            user_id: ID người dùng
            session_id: ID phiên làm việc
            
        Returns:
            Dict chứa phản hồi về đơn hàng
        """
        logging.info("OrderAgent đang xử lý câu hỏi: %s", query)
        
        try:
            # Gọi Agent Order để xử lý
            content = types.Content(role="user", parts=[types.Part(text=query)])
            session_data = await self._ensure_session(user_id, session_id)
            
            # Import process_agent_response từ util
            from ..util import process_agent_response
            
            final_response = None
            async for event in self.runner.run_async(
                user_id=session_data.user_id,
                session_id=session_data.id,
                new_message=content,
                state_delta={"token": token}
            ):
                response = await process_agent_response(event)
                if response:
                    parsed = json.loads(response.get("text", "{}"))
                    response_text = parsed.get("response_text", "")
                    order_ids = parsed.get("order_ids", [])
                    
                    final_response = {
                        "text": response_text,
                        "order_ids": order_ids
                    }
            
            return final_response or {"text": "Không nhận được phản hồi từ agent"}
            
        except Exception as e:
            logging.exception("Lỗi trong OrderAgent.process_order_query")
            return {
                "error": f"Lỗi khi xử lý yêu cầu: {str(e)}",
                "text": "Tôi gặp lỗi khi tra cứu đơn hàng, vui lòng thử lại."
            }
    
    async def search_orders_detail(
        self,
        tool_context: ToolContext,
        status: Optional[str] = None,
        shop_id: Optional[str] = None,
        min_amount: Optional[float] = None,
        max_amount: Optional[float] = None,
        created_from: Optional[str] = None,
        created_to: Optional[str] = None,
        paid_from: Optional[str] = None,
        paid_to: Optional[str] = None,
        processing_from: Optional[str] = None,
        processing_to: Optional[str] = None,
        shipped_from: Optional[str] = None,
        shipped_to: Optional[str] = None,
        completed_from: Optional[str] = None,
        completed_to: Optional[str] = None,
        cancelled_from: Optional[str] = None,
        cancelled_to: Optional[str] = None,
        page: int = 1,
        limit: int = 20,
    ) -> Dict[str, Any]:
        """
        Tool: Tìm kiếm đơn hàng chi tiết với nhiều bộ lọc linh hoạt
        
        Args:
            status: Trạng thái đơn hàng (AWAITING_PAYMENT, PROCESSING, SHIPPED, COMPLETED, CANCELED, REFUNDED)
            shop_id: ID của shop (VD: "shop015", "shop022")
            min_amount: Giá trị tối thiểu của đơn hàng (VD: 100000)
            max_amount: Giá trị tối đa của đơn hàng (VD: 500000)
            created_from: Đơn tạo từ ngày (format: YYYY-MM-DD)
            created_to: Đơn tạo đến ngày (format: YYYY-MM-DD)
            paid_from: Đơn thanh toán từ ngày (format: YYYY-MM-DD)
            paid_to: Đơn thanh toán đến ngày (format: YYYY-MM-DD)
            processing_from: Đơn bắt đầu xử lý từ ngày (format: YYYY-MM-DD)
            processing_to: Đơn bắt đầu xử lý đến ngày (format: YYYY-MM-DD)
            shipped_from: Đơn giao vận chuyển từ ngày (format: YYYY-MM-DD)
            shipped_to: Đơn giao vận chuyển đến ngày (format: YYYY-MM-DD)
            completed_from: Đơn hoàn thành từ ngày (format: YYYY-MM-DD)
            completed_to: Đơn hoàn thành đến ngày (format: YYYY-MM-DD)
            cancelled_from: Đơn hủy từ ngày (format: YYYY-MM-DD)
            cancelled_to: Đơn hủy đến ngày (format: YYYY-MM-DD)
            page: Số trang (mặc định 1)
            limit: Số đơn hàng mỗi trang (mặc định 20, tối đa 100)
        
        Returns:
            Dict chứa mảng kết quả [{order, order_shop}, ...]
        """
        logging.info("search_orders_detail được gọi với các tham số: status=%s, shop_id=%s, min_amount=%s, max_amount=%s",
                     status, shop_id, min_amount, max_amount)
        token = tool_context.state.get("token", "unknown")
        if token == "unknown":
            logging.warning("Token không được cung cấp trong tool_context.state")
            return {"error": "Authorization token is missing"}
        
        async with httpx.AsyncClient() as client:
            try:
                # Xây dựng query params động
                params = {}
                
                if status:
                    valid_statuses = ["AWAITING_PAYMENT", "PROCESSING", "SHIPPED", "COMPLETED", "CANCELED", "REFUNDED"]
                    if status.upper() not in valid_statuses:
                        return {"error": f"Trạng thái '{status}' không hợp lệ. Các trạng thái hợp lệ: {', '.join(valid_statuses)}"}
                    params["status"] = status.upper()
                
                if shop_id:
                    params["shop_id"] = shop_id
                
                if min_amount is not None:
                    params["min_amount"] = min_amount
                
                if max_amount is not None:
                    params["max_amount"] = max_amount
                
                # Các tham số thời gian
                date_params = {
                    "created_from": created_from,
                    "created_to": created_to,
                    "paid_from": paid_from,
                    "paid_to": paid_to,
                    "processing_from": processing_from,
                    "processing_to": processing_to,
                    "shipped_from": shipped_from,
                    "shipped_to": shipped_to,
                    "completed_from": completed_from,
                    "completed_to": completed_to,
                    "cancelled_from": cancelled_from,
                    "cancelled_to": cancelled_to,
                    "page": page,
                    "page_size": limit,
                }
                
                for key, value in date_params.items():
                    if value:
                        params[key] = value
                
                response = await client.get(
                    f"{self.BASE_URL}/orders/search/detail",
                    params=params,
                    headers={"Authorization": f"Bearer {token}"},
                    timeout=20.0
                )
                response.raise_for_status()
                result = response.json()
                
                if result.get("code") != 200:
                    return {"error": f"API trả về lỗi: {result.get('message', 'Unknown error')}"}
                
                # Trả về mảng kết quả
                orders = result.get("result", [])
                
                if not orders:
                    return {"message": "Không tìm thấy đơn hàng nào phù hợp với điều kiện tìm kiếm."}
                
                return {"orders": orders, "total": len(orders)}
                
            except httpx.HTTPStatusError as e:
                logging.error("HTTP error khi tìm kiếm đơn hàng: %s", e)
                return {"error": f"Lỗi API: {e.response.status_code}"}
            except Exception as e:
                logging.exception("Exception khi tìm kiếm đơn hàng")
                return {"error": f"Lỗi hệ thống: {str(e)}"}
    
    async def get_order_detail(self, shop_order_id: str) -> Dict[str, Any]:
        """
        Tool: Lấy thông tin chi tiết của một đơn hàng cụ thể
        
        Args:
            shop_order_id: ID của đơn hàng cần xem chi tiết
        
        Returns:
            Dict chứa thông tin chi tiết đơn hàng (order + order_shop)
        """
        logging.info("get_order_detail được gọi: shop_order_id=%s", shop_order_id)
        
        async with httpx.AsyncClient() as client:
            try:
                response = await client.get(
                    f"{self.BASE_URL}/orders/{shop_order_id}",
                    timeout=15.0
                )
                response.raise_for_status()
                result = response.json()
                
                if result.get("code") != 200:
                    return {"error": f"API trả về lỗi: {result.get('message', 'Unknown error')}"}
                
                return result.get("result", {})
                
            except httpx.HTTPStatusError as e:
                logging.error("HTTP error khi lấy chi tiết đơn hàng: %s", e)
                return {"error": f"Lỗi API: {e.response.status_code}"}
            except Exception as e:
                logging.exception("Exception khi lấy chi tiết đơn hàng")
                return {"error": f"Lỗi hệ thống: {str(e)}"}

