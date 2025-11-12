"""
VoucherAgent - Agent chuyên xử lý tra cứu và tư vấn voucher/mã giảm giá
"""
import json
import logging
from typing import Dict, Any, List, Optional, Literal

import httpx
from google.adk import Agent
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.adk.models.lite_llm import LiteLlm
from google.genai import types
from pydantic import BaseModel, Field
from google.adk.tools.tool_context import ToolContext

from .voucher_promtp import root_instruction


class VoucherResponse(BaseModel):
    """
    Cấu trúc dữ liệu chuẩn hóa cho việc trả về thông tin voucher
    """
    response_text: str = Field(description="Nội dung phản hồi chi tiết về voucher cho người dùng")
    voucher_codes: List[str] = Field(
        default_factory=list,
        description="Danh sách mã voucher được đề cập trong phản hồi"
    )


class VoucherAgent:
    """Agent chuyên xử lý tra cứu và tư vấn voucher"""
    
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
        """Tạo Agent Voucher với cấu hình riêng"""
        return Agent(
            model=self._llm_model if self._llm_model.startswith("gemini") else LiteLlm(model=self._llm_model),
            name="agent_voucher",
            instruction=root_instruction(),
            description="Agent chuyên về tư vấn voucher và mã giảm giá cho khách hàng",
            tools=[
                self.get_vouchers,
            ],
            include_contents="none",
            output_schema=VoucherResponse,
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
    
    async def process_voucher_query(
        self,
        query: str,
        user_id: str,
        token: str,
        session_id: str
    ) -> Dict[str, Any]:
        """
        Xử lý câu hỏi về voucher từ người dùng
        
        Args:
            query: Câu hỏi của người dùng về voucher
            user_id: ID người dùng
            token: JWT token để xác thực
            session_id: ID phiên làm việc
            
        Returns:
            Dict chứa phản hồi về voucher
        """
        logging.info("VoucherAgent đang xử lý câu hỏi: %s", query)
        
        try:
            # Gọi Agent Voucher để xử lý
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
                    voucher_codes = parsed.get("voucher_codes", [])
                    
                    final_response = {
                        "text": response_text,
                        "voucher_codes": voucher_codes
                    }
            
            return final_response or {"text": "Không nhận được phản hồi từ agent"}
            
        except Exception as e:
            logging.exception("Lỗi trong VoucherAgent.process_voucher_query")
            return {
                "error": f"Lỗi khi xử lý yêu cầu: {str(e)}",
                "text": "Tôi gặp lỗi khi tra cứu voucher, vui lòng thử lại."
            }
    
    async def get_vouchers(
        self,
        tool_context: ToolContext,
        owner_type: Optional[Literal["PLATFORM", "SHOP"]] = None,
        shop_id: Optional[str] = None,
        applies_to_type: Optional[Literal["ORDER_TOTAL", "SHIPPING_FEE"]] = None,
        sort_by: Optional[Literal["discount_asc", "discount_desc", "created_at"]] = "discount_desc"
    ) -> Dict[str, Any]:
        """
        Tool: Lấy danh sách voucher với các bộ lọc tùy chọn
        
        Args:
            tool_context: Context từ ADK
            owner_type: Lọc theo chủ sở hữu (PLATFORM: voucher sàn, SHOP: voucher shop)
            shop_id: UUID của shop cụ thể (chỉ dùng khi owner_type=SHOP)
            applies_to_type: Lọc theo loại áp dụng (ORDER_TOTAL: giảm tổng đơn, SHIPPING_FEE: giảm phí ship)
            sort_by: Sắp xếp kết quả (discount_asc: giảm ít->nhiều, discount_desc: giảm nhiều->ít, created_at: mới nhất)
            
        Returns:
            Dict chứa danh sách voucher đã lọc và sắp xếp
        """
        logging.info(
            "get_vouchers được gọi với: owner_type=%s, shop_id=%s, applies_to_type=%s, sort_by=%s",
            owner_type, shop_id, applies_to_type, sort_by
        )
        
        token = tool_context.state.get("token", "unknown")
        if token == "unknown":
            logging.warning("Token không được cung cấp trong tool_context.state")
            return {"error": "Authorization token is missing"}
        
        # Xây dựng query params
        params = {}
        if owner_type:
            params["owner_type"] = owner_type
        if shop_id and owner_type == "SHOP":
            params["shop_id"] = shop_id
        if applies_to_type:
            params["applies_to_type"] = applies_to_type
        if sort_by:
            params["sort_by"] = sort_by
        
        async with httpx.AsyncClient() as client:
            try:
                response = await client.get(
                    f"{self.BASE_URL}/vouchers",
                    headers={"Authorization": f"Bearer {token}"},
                    params=params,
                    timeout=15.0
                )
                response.raise_for_status()
                result = response.json()
                
                if result.get("code") != 200:
                    return {"error": f"API trả về lỗi: {result.get('message', 'Unknown error')}"}
                
                # Lấy danh sách voucher
                vouchers = result.get("result", {}).get("data", [])
                
                if not vouchers:
                    # Tạo thông báo chi tiết dựa trên bộ lọc
                    filter_desc = []
                    if owner_type == "PLATFORM":
                        filter_desc.append("voucher sàn")
                    elif owner_type == "SHOP":
                        filter_desc.append(f"voucher shop {shop_id if shop_id else ''}")
                    if applies_to_type == "ORDER_TOTAL":
                        filter_desc.append("giảm tổng đơn")
                    elif applies_to_type == "SHIPPING_FEE":
                        filter_desc.append("giảm phí ship")
                    
                    filter_text = " ".join(filter_desc) if filter_desc else "voucher"
                    return {
                        "message": f"Không tìm thấy {filter_text} nào.",
                        "vouchers": [],
                        "total": 0,
                        "filters_applied": {
                            "owner_type": owner_type,
                            "shop_id": shop_id,
                            "applies_to_type": applies_to_type,
                            "sort_by": sort_by
                        }
                    }
                
                return {
                    "vouchers": vouchers,
                    "total": len(vouchers),
                    "filters_applied": {
                        "owner_type": owner_type,
                        "shop_id": shop_id,
                        "applies_to_type": applies_to_type,
                        "sort_by": sort_by
                    }
                }
                
            except httpx.HTTPStatusError as e:
                logging.error("HTTP error khi lấy danh sách voucher: %s", e)
                return {"error": f"Lỗi API: {e.response.status_code}"}
            except Exception as e:
                logging.exception("Exception khi lấy danh sách voucher")
                return {"error": f"Lỗi hệ thống: {str(e)}"}
