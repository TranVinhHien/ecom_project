import asyncio
import json
import uuid
from datetime import datetime
from typing import Any, AsyncIterable, List

import httpx
import nest_asyncio

from dotenv import load_dotenv
from google.adk import Agent
from google.adk.agents.readonly_context import ReadonlyContext
from google.adk.artifacts import InMemoryArtifactService
from google.adk.memory.in_memory_memory_service import InMemoryMemoryService
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.adk.tools.tool_context import ToolContext
from google.genai import types
from typing import Any, AsyncIterable, List, Optional  # Đảm bảo import MappingProxyType nếu bạn dùng Python 3.9 trở lên
import os
from google.adk.models import LlmResponse, LlmRequest
from google.adk.agents.callback_context import CallbackContext
import base64
import litellm 
from google.adk.models.lite_llm import LiteLlm  # For multi-model support
import logging

from .agent import RAGAgent, OrderAgent, VoucherAgent, ProductDetailAgent
from .util import process_agent_response
from google.adk.tools.base_tool import BaseTool
from typing import Dict, Any, Optional
from google.adk.agents.invocation_context import InvocationContext
from google.adk.tools.mcp_tool.mcp_toolset import MCPToolset, StdioServerParameters, StdioConnectionParams
import sys
from google.adk.sessions import DatabaseSessionService 

from pathlib import Path
PATH_TO_SCHOOL_MCP_SERVER = str((Path(__file__).parent / "school_mcp_server.py").resolve())
#######################################################################################
#############                                                             #############
#############                                                             #############
#############                                                             #############
############# uv run uvicorn host:app --host 0.0.0.0 --port 9102 --reload #############
#############                                                             #############
#############                                                             #############
#######################################################################################from google.adk.tools.base_tool import BaseTool

load_dotenv()
nest_asyncio.apply()
llm_model = os.getenv("LLM_MODEL")
db_url = os.getenv("DB_URL")

session_service = DatabaseSessionService(db_url=db_url) 
session_service_memory = InMemorySessionService()

def safe_parse_json(text_content):
    if not text_content or not text_content.strip():
        return None  # hoặc raise Exception("Empty response from tool")
    try:
        return json.loads(text_content)
    except json.JSONDecodeError:
        raise Exception(f"Invalid JSON: {text_content}")

def bmc_trim_llm_request(
    callback_context: CallbackContext, llm_request: LlmRequest
) -> Optional[LlmResponse]:

    max_prev_user_interactions = int(os.environ.get("HISTORY_LENGTH","-1"))
    max_prev_user_interactions = 15
    temp_processed_list = []
    
    if max_prev_user_interactions == -1:
        return None 
    else:
        user_message_count = 0
        # Iterate in reverse order
        for i in range(len(llm_request.contents) - 1, -1, -1):
            item = llm_request.contents[i]
            if item.parts[0] and item.parts[0].function_response and item.parts[0].function_response.response and item.parts[0].function_response.response.get("result") and item.parts[0].function_response.response.get("result")[0].get("kind") in ["data", "image"]:
                item.parts[0].function_response.response["result"][0]["data"] = {
                    "message": item.parts[0].function_response.response["result"][0]["data"]["message"]
                }
            # Check if the item is a user message and has text content and is not a transfer to agent content
            if item.role == "user" and item.parts[0] and item.parts[0].text and item.parts[0].text != "For context:":
                user_message_count += 1

            if user_message_count > max_prev_user_interactions:
                temp_processed_list.append(item) # make sure we add this user message.
                break
            temp_processed_list.append(item)

        final_list = temp_processed_list[::-1]

        if user_message_count < max_prev_user_interactions:
            logging.info("User message count did not reach the allowed limit. List remains unchanged.")
        else:
            logging.info(f"User message count reached {max_prev_user_interactions}. List truncated.")
            llm_request.contents = final_list


    return None 

async def  after_tool_call(
    tool: BaseTool, args: Dict[str, Any], tool_context: ToolContext, tool_response: Dict
) ->Dict | str:
    """Callback sau khi gọi tool - lưu response vào file"""
    
    if hasattr(tool_response, "model_dump"):
        raw_response = tool_response.model_dump()
    elif hasattr(tool_response, "content"):

        raw_response = tool_response.content
    elif isinstance(tool_response, dict):

        raw_response = tool_response
    else:
        print("[Callback] Không biết unwrap kiểu gì")
        return None
    if raw_response.get("isError"):
        return None
    text_content = raw_response.get("content", [])[0].get("text", "")
    parsed_json = safe_parse_json(text_content)
    if parsed_json.get("type") == "text":
        return None
    tool_context.actions.skip_summarization=True
    tool_context.actions.escalate=True
    return parsed_json # Return the modified dictionary

class HostAgent:
    """Agent chính điều phối các tác vụ"""

    def __init__(
        self,
        name_agent:str
    ):
        # Khởi tạo RAGAgent như một dependency
        self._rag_agent = RAGAgent(
            llm_model=llm_model,
            session_service=session_service_memory
        )
        
        # Khởi tạo OrderAgent
        self._order_agent = OrderAgent(
            llm_model=llm_model,
            session_service=session_service_memory
        )
        
        # Khởi tạo VoucherAgent
        self._voucher_agent = VoucherAgent(
            llm_model=llm_model,
            session_service=session_service_memory
        )
        
        # Khởi tạo ProductDetailAgent
        self._product_detail_agent = ProductDetailAgent(
            llm_model=llm_model,
            session_service=session_service_memory
        )
        
        self._agent = self.create_agent(name_agent)
        self._user_id = "host_agent"
        self.runner = Runner(
            app_name=self._agent.name,
            agent=self._agent,
            session_service=session_service,
        )

    @classmethod
    async def create(
        cls,
        name: str
    ):
        instance =  HostAgent(name)
        return instance
    
    def create_agent(self,name) -> Agent:
        return Agent(
            model=llm_model if llm_model.startswith("gemini") else LiteLlm(model=llm_model),
            name=name,
            instruction=self.root_instruction,
            description="Bạn là 1 người tư vấn viên chăm sóc khách hàng trong bán sản phẩm ở sàn thương mại điện tử.",
            tools=[
                self.agent_rag_tool,
                self.agent_order_tool,
                self.agent_voucher_tool,
                self.agent_product_detail_tool,
                self.remote_complaint_page,
            ],
            before_model_callback=bmc_trim_llm_request,            
        )
    # def root_instruction(self, context: ReadonlyContext) -> str:
    #     lang = context.state.get("lang")
    #     user_info = context.state.get("user_info")
        
    #     # Sử dụng f-string với ba dấu ngoặc nhọn để format sau
    #     # (Hoặc giữ nguyên nếu bạn format ngay lập tức như code gốc)
    #     text = f"""
    #     # BỐI CẢNH VÀ VAI TRÒ
    #     Bạn là một trợ lý AI chuyên nghiệp, đóng vai trò là nhân viên tư vấn và chăm sóc khách hàng cho một sàn thương mại điện tử.
    #     Mục tiêu chính của bạn là hỗ trợ khách hàng giải quyết các yêu cầu của họ một cách chính xác và hiệu quả.

    #     # QUY TẮC VÀNG (TUYỆT ĐỐI TUÂN THỦ)
    #     1.  **ƯU TIÊN YÊU CẦU HIỆN TẠI:** Bạn phải luôn xử lý yêu cầu *hiện tại* của người dùng. Nếu người dùng đưa ra một yêu cầu (ví dụ: "kiểm tra đơn hàng của tôi"), bạn phải thực hiện nó bằng cách sử dụng công cụ thích hợp.
    #     2.  **KHÔNG TỪ CHỐI DỰA TRÊN LỊCH SỬ THẤT BẠI:** Nếu một hành động trước đó không thành công (ví dụ: `agent_order_tool` báo lỗi hoặc không tìm thấy đơn hàng) và người dùng *vẫn lặp lại yêu cầu*, bạn **TUYỆT ĐỐI** không được trả lời rằng bạn đã thử và thất bại. Bạn **PHẢI** thực thi lại yêu cầu đó bằng công cụ một lần nữa. Luôn giả định rằng yêu cầu mới là một lệnh mới cần được thực thi.
    #     3.  **THÁI ĐỘ:** Luôn giữ thái độ lịch sự, thân thiện và chuyên nghiệp.
    #     4.  **HẠN CHẾ HỎI LẠI:** Các công cụ được thiết kế để tự động xử lý. Chỉ hỏi lại thông tin khi *thực sự* không thể suy ra từ yêu cầu của người dùng.
        
    #     # QUY TRÌNH XỬ LÝ
    #     1.  **Phân tích yêu cầu:** Đọc kỹ yêu cầu của khách hàng để hiểu rõ họ muốn gì.
    #     2.  **Lựa chọn công cụ (Tool):** Dựa trên phân tích, chọn *chính xác* một trong các công cụ dưới đây.
    #     3.  **Thực thi:** Gọi công cụ đã chọn với các tham số chính xác.
        
    #     # DANH SÁCH CÔNG CỤ (TOOLS)

    #     ## 1. Tool: agent_rag_tool
    #     * **Mục đích:** Tìm kiếm thông tin chung về sản phẩm, danh mục sản phẩm, và các chính sách (bảo hành, đổi trả, vận chuyển).
    #     * **Khi nào sửs dụng:** Khi khách hàng hỏi về:
    #         * Thông tin sản phẩm: "Laptop A có tốt không?", "So sánh điện thoại B và C."
    #         * Tìm kiếm sản phẩm: "Tìm cho tôi laptop dưới 20 triệu."
    #         * Chính sách: "Chính sách đổi trả thế nào?", "Phí vận chuyển ra sao?"
    #     * **Tham số:**
    #         * `query` (str): Tóm tắt ngắn gọn, đầy đủ ngữ cảnh yêu cầu tìm kiếm.
    #         * `intent` (str): "product" (cho sản phẩm) hoặc "policy" (cho chính sách).
    #     * **Ví dụ gọi:**
    #         * Hỏi: "So sánh iPhone 14 và 15" -> `agent_rag_tool(query="So sánh iPhone 14 và 15", intent="product")`
    #         * Hỏi: "Chính sách bảo hành" -> `agent_rag_tool(query="Chính sách bảo hành", intent="policy")`

    #     ## 2. Tool: agent_order_tool
    #     * **Mục đích:** Tra cứu thông tin *cụ thể* về đơn hàng của khách hàng (dựa trên `user_info`).
    #     * **Khi nào sử dụng:** Khi khách hàng hỏi về các đơn hàng *của riêng họ*:
    #         * Trạng thái đơn hàng: "Đơn hàng của tôi đâu rồi?"
    #         * Lịch sử mua hàng: "Xem các đơn tôi mua tháng 10."
    #         * Chi tiết đơn hàng: "Đơn hàng 123ABC có những gì?"
    #         * Các yêu cầu lọc: "Đơn đã hủy giá trên 100k", "Đơn ở shop015", "Đơn thanh toán hôm qua"
    #     * **Tham số:**
    #         * `query` (str): Nguyên văn hoặc tóm tắt câu hỏi của khách hàng về đơn hàng.
    #     * **Ví dụ gọi:**
    #         * Hỏi: "Xem đơn hàng của tôi" -> `agent_order_tool(query="Xem đơn hàng của tôi")`
    #         * Hỏi: "Đơn ở shop015" -> `agent_order_tool(query="Đơn ở shop015")`

    #     ## 3. Tool: remote_complaint_page
    #     * **Mục đích:** Tra cứu và tư vấn về voucher/mã giảm giá cho khách hàng.
    #     * **Khi nào sử dụng:** Khi khách hàng hỏi về:
    #         * Danh sách voucher: "Có voucher nào không?", "Voucher nào tốt nhất?"
    #         * Tư vấn voucher: "Đơn 200k dùng voucher gì?", "Voucher giảm nhiều nhất"
    #         * Loại voucher: "Voucher giảm ship", "Voucher shop ABC"
    #     * **Tham số:**   ?:
    #         * `category` (str): Loại phàn nàn gồm có 
    #             BUG: người dùng báo cáo lỗi hệ thống. 
    #             COMPLAINT: Phàn nàn về các chất lượng dịch vụ hoặc các vấn đề liên quan đến mua sắm khác..
    #             SUGGESTION: Gợi ý cải thiện dịch vụ.
    #             OTHER: Các loại khác.
    #         * `content` (str): Toàn bộ nội nội dung khiếu nại của khách hàng , bạn có thể tóm tắt nêu các ý khiếu nại của người dùng ở trong đây và trả về/.
    #     * **Ví dụ gọi:**
    #         * Hỏi: "Hệ thống đang rất ít voucher khuyến mãi cho người dùng mới" -> `remote_complaint_page(category="SUGGESTION",content="Hệ thống đang rất ít voucher khuyến mãi cho người dùng mới cần phải cải thiện thêm")`
            
            
    #     ## 4. Tool: agent_voucher_tool
    #     * **Mục đích:** Chuyển hướng người dùng tới trang khiếu nại và tự điền khiếu nại giúp cho người dùng.
    #     * **Khi nào sử dụng:** Khi khách hàng có phàn nàn và muốn khiếu nại về sản phẩm/dịch vụ.:
    #     * **Tham số:**
    #         * `query` (str): Nguyên văn câu hỏi của khách hàng về voucher.
    #     * **Ví dụ gọi:**
    #         * Hỏi: "Có voucher nào không?" -> `agent_voucher_tool(query="Có voucher nào không?")`
    #         * Hỏi: "Đơn 200k dùng voucher gì?" -> `agent_voucher_tool(query="Đơn 200k dùng voucher gì?")`

    #     # PHÂN BIỆT RÕ RÀNG
    #     * `agent_rag_tool`: Dùng cho thông tin chung, catalog sản phẩm, chính sách (cho mọi người).
    #     * `agent_order_tool`: Dùng cho thông tin cá nhân, lịch sử mua hàng (chỉ cho người dùng hiện tại).
    #     * `agent_voucher_tool`: Dùng cho tra cứu và tư vấn voucher/mã giảm giá.
        
    #     # DỮ LIỆU BỐI CẢNH
    #     * **Ngôn ngữ giao tiếp:** {lang}
    #     * **Thông tin người dùng:** {user_info}
    #     Ngày tháng năm hiện tại (YYYY-MM-DD): {datetime.now().strftime("%Y-%m-%d")}
    #     """
    #     return text
    def root_instruction(self, context: ReadonlyContext) -> str:
        lang = context.state.get("lang")
        user_info = context.state.get("user_info")
        
        text = f"""
# VAI TRÒ
Trợ lý AI chăm sóc khách hàng sàn thương mại điện tử.

# QUY TẮC
1. Ưu tiên yêu cầu hiện tại
2. Retry nếu tool lỗi và user lặp lại yêu cầu
3. Thái độ lịch sự, thân thiện

# TOOL MAPPING
- **agent_rag_tool**: Sản phẩm, chính sách chung (intent: "product"|"policy")
- **agent_order_tool**: Đơn hàng cá nhân (lọc: status, shop, giá, ngày)
- **agent_voucher_tool**: Voucher/mã giảm giá
- **agent_product_detail_tool**: Chi tiết sản phẩm + phân tích đánh giá (input: tên/key sản phẩm) nếu trong câu truy vấn của người dùng có product_key, lập tức sử dụng tool này.
- **remote_complaint_page**: Khiếu nại (category: BUG|COMPLAINT|SUGGESTION|OTHER)

# CONTEXT
Lang: {lang} | User: {user_info} | Date: {datetime.now().strftime("%Y-%m-%d")}
        """
        return text
    async def agent_rag_tool(self, query: str, intent: str, tool_context: ToolContext):
        """
        Tìm kiếm sản phẩm/chính sách qua RAG
        
        Args:
            query: Câu hỏi người dùng liên quan đến sản phẩm/chính sách
            intent: "product" | "policy"
        """
        logging.info("HostAgent.agent_rag_tool được gọi với intent: %s", intent)
        
        user_id = tool_context.state.get("user_id", "unknown")
        session_id = tool_context.state.get("session_id", str(uuid.uuid4()))
        
        # Gọi RAGAgent để xử lý
        result = await self._rag_agent.search_and_analyze(
            query=query,
            intent=intent,
            user_id=user_id,
            session_id=session_id
        )
        
        # Báo cho ADK không tóm tắt lại và escalate kết quả
        tool_context.actions.skip_summarization = True
        tool_context.actions.escalate = True
        
        return result
    
    async def remote_complaint_page(self, category: str, content: str, tool_context: ToolContext):
        """
        Chuyển hướng trang khiếu nại
        
        Args:
            category: BUG | COMPLAINT | SUGGESTION | OTHER
            content: Nội dung khiếu nại
        """
        
      
        tool_context.actions.skip_summarization = True
        tool_context.actions.escalate = True
        return {
            "category": category,
            "content": content,
            "text":"Cảm ơn bạn đã góp ý, hệ thống đã chuyển hướng bạn tới trang góp ý, bạn có thể gửi góp ý tới sàn để phía sàn sau này cải thiện tốt hơn!!!"
        }
    
    
    async def agent_order_tool(self, query: str, tool_context: ToolContext):
        """
        Tra cứu đơn hàng cá nhân
        
        Args:
            query: Câu hỏi của người dùng về đơn hàng        
        Examples:
            - "Xem đơn hàng của tôi"
            - "Đơn tháng 10"
            - "Đơn đã hủy giá trên 100k"
            - "Đơn thanh toán hôm qua ở shop015"
        """
        logging.info("HostAgent.agent_order_tool được gọi với query: %s", query)
        
        token = tool_context.state.get("token", "unknown")
        user_id = tool_context.state.get("user_id", "unknown")
        if token == "unknown":
            
            logging.warning("Token không được cung cấp trong tool_context.state")
            return {"error": "Authorization token is missing"}
        session_id = tool_context.state.get("session_id", str(uuid.uuid4()))
        
        # Gọi OrderAgent để xử lý
        result = await self._order_agent.process_order_query(
            query=query,
            user_id=user_id,
            token=token,
            session_id=session_id,
        )
        
        # Báo cho ADK không tóm tắt lại và escalate kết quả
        tool_context.actions.skip_summarization = True
        tool_context.actions.escalate = True
        
        return result
    
    async def agent_voucher_tool(self, query: str, tool_context: ToolContext):
        """
        Tra cứu và tư vấn voucher
        
        Args:
            query: Câu hỏi của người dùng về voucher
        Examples:
            - "Có voucher nào không?"
            - "Voucher giảm nhiều nhất"
            - "Đơn 200k dùng voucher gì?"
            - "Voucher giảm ship"
        """
        logging.info("HostAgent.agent_voucher_tool được gọi với query: %s", query)
        
        token = tool_context.state.get("token", "unknown")
        user_id = tool_context.state.get("user_id", "unknown")
        if token == "unknown":
            logging.warning("Token không được cung cấp trong tool_context.state")
            return {"error": "Authorization token is missing"}
        session_id = tool_context.state.get("session_id", str(uuid.uuid4()))
        
        # Gọi VoucherAgent để xử lý
        result = await self._voucher_agent.process_voucher_query(
            query=query,
            user_id=user_id,
            token=token,
            session_id=session_id,
        )
        
        # Báo cho ADK không tóm tắt lại và escalate kết quả
        tool_context.actions.skip_summarization = True
        tool_context.actions.escalate = True
        
        return result
    
    async def agent_product_detail_tool(self, query: str, tool_context: ToolContext):
        """
        Phân tích tóm tắt từ đánh giá chi tiết sản phẩm
        Khi yêu cầu người dùng bắt đầu bằng từ "Gợi ý sản phẩm" sẽ thực hiện gọi tới tool này.
        Args:
            query: Yêu cầu cua người dùng về chi tiết sản phẩm nào
        """
        logging.info("HostAgent.agent_product_detail_tool được gọi với query: %s", query)
        
        token = tool_context.state.get("token", "unknown")
        user_id = tool_context.state.get("user_id", "unknown")
        product_key = tool_context.state.get("product_key", "unknown")
        if token == "unknown":
            logging.warning("Token không được cung cấp trong tool_context.state")
            return {"error": "Authorization token is missing"}
        if product_key == "unknown":
            logging.warning("Product key không được cung cấp trong tool_context.state")
            return {"error": "Không tìm thấy sản phẩm để phân tích"}
        session_id = tool_context.state.get("session_id", str(uuid.uuid4()))
        
        # Gọi ProductDetailAgent để xử lý
        result = await self._product_detail_agent.process_product_query(
            query=query,
            user_id=user_id,
            token=token,
            session_id=session_id,
            product_key=product_key
        )
        
        # Báo cho ADK không tóm tắt lại và escalate kết quả
        tool_context.actions.skip_summarization = True
        tool_context.actions.escalate = True
        
        return result