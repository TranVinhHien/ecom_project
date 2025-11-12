
import asyncio
import json
import logging
import os
import aiohttp
from typing import Dict, Any, Optional

from dotenv import load_dotenv
from datetime import datetime
from mcp.server.fastmcp import FastMCP
from google.adk.tools.tool_context import ToolContext

from mcp.server.stdio import stdio_server
import sys
load_dotenv()

# --- Logging Setup ---

# --- Logging Setup ---
LOG_FILE_PATH = os.path.join(os.path.dirname(__file__), "school_mcp_server_activity.log")

# SỬA LẠI ĐOẠN NÀY
logging.basicConfig(
    level=logging.DEBUG,
    format="%(asctime)s - %(levelname)s - [%(filename)s:%(lineno)d] - %(message)s",
    handlers=[logging.StreamHandler(sys.stderr)]
)


# Lấy logger gốc
logger = logging.getLogger() 

# Tạo một FileHandler với encoding UTF-8
# Đây là dòng quan trọng nhất
file_handler = logging.FileHandler(LOG_FILE_PATH, mode="w", encoding="utf-8")

# Tạo một formatter và gán nó cho handler
formatter = logging.Formatter("%(asctime)s - %(levelname)s - [%(filename)s:%(lineno)d] - %(message)s")
file_handler.setFormatter(formatter)

# Thêm handler đã cấu hình vào logger gốc
# Nếu có handler cũ, có thể cần xóa đi trước logger.handlers.clear()
logger.addHandler(file_handler)

# API Configuration
API_BASE_URL = "https://ai-api.bitech.vn/api"
async def make_api_request(
    method: str, 
    endpoint: str, 
    *,  # Dấu sao quan trọng ở đây
    data: Dict = None, 
    auth_required: bool = True,
    token: Optional[str] = None
) -> Dict[str, Any]:
    """Thực hiện HTTP request đến API, sử dụng global ACCESS_TOKEN."""
    url = f"{API_BASE_URL}{endpoint}"
    headers = {"Content-Type": "application/json"}
    # print(f"Token trong get_profile: {token}")
    # logging.debug(f"Thực hiện request: {method.upper()} {token} với auth_required={auth_required}")
    token = token.strip('"')
    token=token.replace('"','')
    if isinstance(token, bytes):
            token = token.decode("utf-8", errors="ignore")
    access_token = token
    if auth_required:
        if not token:
            logging.warning(f"Yêu cầu xác thực cho endpoint '{endpoint}' nhưng không tìm thấy token.")
            return {"success": False, "message": "Lỗi xác thực: Bạn chưa đăng nhập hoặc phiên đã hết hạn. Vui lòng sử dụng tool 'login'."}
        

        headers["Authorization"] = f"Bearer {token}"
    
    logging.debug(f"Thực hiện request: {method.upper()} {url} với auth_required={auth_required}")
    
    async with aiohttp.ClientSession() as session:
        try:
            # (Phần còn lại của hàm giữ nguyên)
            if method.upper() == "GET":
                async with session.get(url, headers=headers) as response:
                    return await response.json()
            elif method.upper() == "POST":
                async with session.post(url, headers=headers, json=data) as response:
                    return await response.json()
            elif method.upper() == "PUT":
                async with session.put(url, headers=headers, json=data) as response:
                    return await response.json()
        except Exception as e:
            logging.error(f"Lỗi khi gọi API tới {url}: {e}", exc_info=True)
            return {"success": False, "message": f"Lỗi kết nối API: {str(e)}"}

# --- MCP Server Setup ---
logging.info("Tạo MCP Server cho hệ thống quản lý trường học...")
mcp = FastMCP("school-management-mcp-server")

# --- Authentication Functions ---


@mcp.tool()
async def get_order(accessToken:str,status:str) -> str:
    """
    Xem đơn hàng.

    - Dùng khi học sinh cần xem thông tin đơn hàng.

    """

    
    result = await make_api_request("GET", "/student/schedule", auth_required=True,token=accessToken)
    if not result or not result.get("success", False):
        return json.dumps({
            "success": False,
            "message": result.get("message", "Không lấy được dữ liệu lịch học."),
            "result": result,
            "type":"error",
        }, ensure_ascii=False)
    return json.dumps({
        "type": "data",
        "data_type":"schedule",
        "result": result.get("data").get("student_schedule"),
        "message": "Lấy dữ liệu lịch học thành công."},
        ensure_ascii=False)


# --- MCP Server Runner ---
async def run_mcp_stdio_server():
    """Chạy MCP server, lắng nghe kết nối qua standard input/output."""
    async with stdio_server() as (read_stream, write_stream):
        logging.info("MCP Stdio Server: Bắt đầu handshake với client...")
        await mcp._mcp_server.run(
            read_stream,
            write_stream,
            mcp._mcp_server.create_initialization_options()
        )
        logging.info("MCP Stdio Server: Kết thúc hoặc client đã ngắt kết nối.")

if __name__ == "__main__":
    logging.info("Khởi động School Management MCP Server qua stdio...")
    try:
        asyncio.run(run_mcp_stdio_server())
    except KeyboardInterrupt:
        logging.info("\nMCP Server (stdio) đã dừng bởi người dùng.")
    except Exception as e: 
        logging.critical(
            f"MCP Server (stdio) gặp lỗi không xử lý được: {e}", exc_info=True
        )
    finally:
        logging.info("MCP Server (stdio) đã thoát.")
        
        
        
        
        
@mcp.tool()
async def create_class(accessToken:str,course_id: int = 0, semester: str = "", academic_year: str = "", 
                      max_capacity: int = 0, start_date: str = "", end_date: str = "") -> str:

    """
    Tạo mới một lớp học trong hệ thống (dành cho quản lý).
    """

    # Lấy danh sách khóa học từ API
    courses_resp = await make_api_request("GET", "/manager/courses", auth_required=True, token=accessToken)
    courses = []
    if courses_resp and "data" in courses_resp and "courses" in courses_resp["data"]:
        courses = [
            {"label": f"{c['course_name']} ({c['course_code']})", "value": str(c["course_id"])}
            for c in courses_resp["data"]["courses"]
        ]
    class_form = {
    "fields": [
        {
            "field": "course_id",
            "type": "combobox",
            "label": "Chọn khóa học",
            "options": courses,
            "value": str(course_id) if course_id else (courses[0]["value"] if courses else None),
        },
        {
            "field": "semester",
            "type": "combobox",
            "label": "Học kỳ",
            "options": ["Học kỳ 1", "Học kỳ 2", "Học kỳ hè"],
            "value": semester or "Học kỳ 1",
        },
        {
            "field": "academic_year",
            "type": "text",
            "label": "Năm học (current year - next year)",
            "value": academic_year if academic_year else f"{datetime.now().year}-{datetime.now().year + 1}",
        },
        {
            "field": "max_capacity",
            "type": "number",
            "label": "Sức chứa tối đa",
            "value": int(max_capacity) if max_capacity > 0 else 30,
        },
        {
            "field": "start_date",
            "type": "date",
            "label": "Ngày bắt đầu",
            "value": (start_date or datetime.now()).strftime("%Y-%m-%d"),
        },
        {
            "field": "end_date",
            "type": "date",
            "label": "Ngày kết thúc",
            "value": (end_date or datetime.now()).strftime("%Y-%m-%d"),
        },
    ],
    "title": "Tạo lớp học mới",
    "description": "Vui lòng điền thông tin lớp học mới",
    "submit_label": "Tạo lớp học",
    "submit_endpoint": f"{API_BASE_URL}/manager/create-class",
}

    
    
    # result = await make_api_request("POST", "/manager/create-class" , data=class_data, auth_required=True,token=accessToken)
    return json.dumps({"type": "data","data_type":"form","message":"tạo lớp học","result":class_form}, ensure_ascii=False)

@mcp.tool()
async def update_class( accessToken:str,class_id: int=0,semester: str = "",course_id: int=0, academic_year: str = "",
                      max_capacity: int = 0, start_date: str = "", 
                      end_date: str = "", status: str = "",) -> str:
    """
    Chuẩn bị form cập nhật thông tin lớp học.
    - Nếu caller không truyền tham số (semester, academic_year, max_capacity,
      start_date, end_date, status) thì hàm sẽ gọi API GET /manager/get-class/{class_id}
      để lấy giá trị hiện tại và dùng những giá trị đó.
    - Trả về: JSON chứa message + form đã điền sẵn.
    """

    # Gọi API để lấy thông tin lớp học hiện tại
    current_data = await make_api_request(
        "GET", f"/manager/get-class/{class_id}", 
        auth_required=True, token=accessToken
    )

    # Nếu caller không truyền giá trị thì dùng từ current_data
    semester = semester or current_data.get("semester", "")
    academic_year = academic_year or current_data.get("academic_year", "")
    max_capacity = max_capacity or current_data.get("max_capacity", 0)
    start_date = start_date or current_data.get("start_date", "")
    end_date = end_date or current_data.get("end_date", "")
    status = status or current_data.get("status", "")

    status_options = [
        {"label": "Mở đăng ký", "value": "OPEN"},
        {"label": "Đang học", "value": "IN_PROGRESS"},
        {"label": "Hoàn thành", "value": "COMPLETED"},
        {"label": "Hủy bỏ", "value": "CLOSED"},
    ]

    update_form = {
        "fields": [
            {"field": "class_id", "type": "number", "label": "ID lớp học", "value": class_id, "readonly": True},
  {
            "field": "semester",
            "type": "combobox",
            "label": "Học kỳ",
            "options": ["Học kỳ 1","Học kỳ 2","Học kỳ hè"],
            "value": semester or "Học kỳ 1",
        },
            {"field": "academic_year", "type": "text", "label": "Năm học(current_year-next_year)", "value": academic_year or f'{datetime.now().year}-{datetime.now().year + 1}'},
            {"field": "max_capacity", "type": "number", "label": "Sức chứa tối đa", "value": int(max_capacity)},
            {"field": "start_date", "type": "date", "label": "Ngày bắt đầu", "value": start_date or datetime.now()},
            {"field": "end_date", "type": "date", "label": "Ngày kết thúc", "value": end_date or datetime.now()},
            {"field": "status", "type": "combobox", "label": "Trạng thái", 
             "options": status_options, "value": status or status_options[0]["value"]},
        ],
        "title": "Cập nhật lớp học",
        "description": f"Vui lòng chỉnh sửa thông tin cho lớp ID: {class_id}",
        "submit_label": "Cập nhật lớp học",
        "submit_endpoint": f"{API_BASE_URL}/manager/update-class/{class_id}"
    }

    return json.dumps({
        "type": "data",
        "data_type": "form",
        "message": f"Form cập nhật lớp học ID {class_id} đã được chuẩn bị.",
        "result": update_form
    }, ensure_ascii=False)
