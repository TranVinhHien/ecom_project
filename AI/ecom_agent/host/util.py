from google.genai import types
from google.adk.runners import Runner
from google.adk.sessions import BaseSessionService 
from google.adk.events import Event
from datetime import datetime
import tempfile, os, mimetypes
import base64
import uuid
import json
from dotenv import load_dotenv
load_dotenv()
import os
import jwt
secret_key = os.getenv("SECRET_KEY","")
# ANSI color codes for terminal output
class Colors:
    RESET = "\033[0m"
    BOLD = "\033[1m"
    UNDERLINE = "\033[4m"

    # Foreground colors
    BLACK = "\033[30m"
    RED = "\033[31m"
    GREEN = "\033[32m"
    YELLOW = "\033[33m"
    BLUE = "\033[34m"
    MAGENTA = "\033[35m"
    CYAN = "\033[36m"
    WHITE = "\033[37m"

    # Background colors
    BG_BLACK = "\033[40m"
    BG_RED = "\033[41m"
    BG_GREEN = "\033[42m"
    BG_YELLOW = "\033[43m"
    BG_BLUE = "\033[44m"
    BG_MAGENTA = "\033[45m"
    BG_CYAN = "\033[46m"
    BG_WHITE = "\033[47m"

 
async def display_state(
    session_service: BaseSessionService, app_name: str, user_id: str, session_id: str, label: str = "Current State"
):
    """Display the current session state in a formatted way."""
    try:
        session = await session_service.get_session(
            app_name=app_name, user_id=user_id, session_id=session_id
        )

        # Format the output with clear sections
        print(f"\n{'-' * 10} {label} {'-' * 10}")

        # Handle the user name
        # user_name = await session.state.get("user_name", "Unknown")
        print(f"👤 User: {session}")



        print("-" * (22 + len(label)))
    except Exception as e:
        print(f"Error displaying state: {e}")

async def process_agent_response(event: Event):
    """Process and display agent response events."""
    # Check for specific parts first
    has_specific_part = False
    if event.content and event.content.parts:
        for part in event.content.parts:
            if hasattr(part, "executable_code") and part.executable_code:
                # Access the actual code string via .code
                print(
                    f"  Debug: Agent generated code:\n```python\n{part.executable_code.code}\n```"
                )
                has_specific_part = True
            elif hasattr(part, "code_execution_result") and part.code_execution_result:
                # Access outcome and output correctly
                print(
                    f"  Debug: Code Execution Result: {part.code_execution_result.outcome} - Output:\n{part.code_execution_result.output}"
                )
                has_specific_part = True
            elif hasattr(part, "tool_response") and part.tool_response:
                # Print tool response information
                print(f"  Tool Response: {part.tool_response.output}")
                has_specific_part = True
            # Also print any text parts found in any event for debugging
            elif hasattr(part, "text") and part.text and not part.text.isspace():
                print(f"  Text: '{part.text.strip()}'")
    print(
                f"\n{Colors.BG_YELLOW}{Colors.WHITE}{Colors.BOLD}╔══ Event final ═════════════════════════════════════════{Colors.RESET}"
            )
    print(f"{Colors.CYAN}{Colors.BOLD}{event}{Colors.RESET}")
    print(
                f"{Colors.BG_YELLOW}{Colors.WHITE}{Colors.BOLD}╚═════════════════════════════════════════════════════════════{Colors.RESET}\n"
            )
    # Check for final response after specific parts
    final_response = None
    if event.is_final_response():

        # return file
        if part.function_response and part.function_response.response:
            if type(part.function_response.response.get("text")) == str:
                final_response =part.function_response.response
                # {
                #     "text": part.function_response.response.get("text","Lỗi không lấy được câu trả lời từ Agent")
                # }
                return final_response
        
            if part.function_response.response.get("result")[0].get("kind") == "text":
                final_response = {
                    "text": part.function_response.response.get("result")[0].get("text","Lỗi không lấy được câu trả lời từ Agent")
                }
                return final_response
            if part.function_response.response.get("result")[0].get("kind") == "file":
                final_response = {
                    "result":part.function_response.response.get("result")
                }
                return final_response
            if part.function_response.response.get("result")[0].get("kind") == "data":
                return {
                    "result":part.function_response.response.get("result")
                }
        # return text agent
        if (
            event.content
            and event.content.parts
            and hasattr(event.content.parts[0], "text")
            and event.content.parts[0].text
        ):
            final_response = {
                "text":event.content.parts[0].text.strip()
            }
            # Use colors and formatting to make the final response stand out
            
            # print(
            #     f"\n{Colors.BG_BLUE}{Colors.WHITE}{Colors.BOLD}╔══ AGENT RESPONSE ═════════════════════════════════════════{Colors.RESET}"
            # )
            # print(f"{Colors.CYAN}{Colors.BOLD}{final_response}{Colors.RESET}")
            # print(
            #     f"{Colors.BG_BLUE}{Colors.WHITE}{Colors.BOLD}╚═════════════════════════════════════════════════════════════{Colors.RESET}\n"
            # )
        else:
            print(
                f"\n{Colors.BG_RED}{Colors.WHITE}{Colors.BOLD}==> Final Agent Response: [No text content in final event]{Colors.RESET}\n"
            )
            return {"text":"Không nhận được kết quả trả về."}

    return final_response

LOG_DIR = os.getenv("LOG_DIR", "logs")
os.makedirs(LOG_DIR, exist_ok=True)
import csv

async def call_agent_async(runner: Runner, user_id: str, session_id: str, query: str, token: str):
    """Call the agent asynchronously with the user's query and log all processing information."""
    # Khởi tạo thông tin log
    start_time = datetime.now()
    log_entry = {
        "timestamp": start_time.isoformat(),
        "user_id": user_id,
        "session_id": session_id,
        "query": query,
        "events": [],
        "final_response": None,
        "processing_time_ms": None,
        "success": False,
        "error": None
    }
    
    content = types.Content(role="user", parts=[types.Part(text=query)])
    print(
        f"\n{Colors.BG_GREEN}{Colors.BLACK}{Colors.BOLD}--- Running Query: {query} ---{Colors.RESET}"
    )
    
    final_response_text = None
    state_delta: dict[str, str] = {
        "token": token,
        "user_id": user_id,
        "session_id": session_id
    }
    
    try:
        event_count = 0
        exec_time = datetime.now()
        async for event in runner.run_async(
            user_id=user_id, session_id=session_id, new_message=content, state_delta=state_delta
        ):
            event_count += 1
            
            # Log thông tin chi tiết của event
            event_log = {
                "event_number": event_count,
                "event_type": type(event).__name__,
                "timestamp": datetime.now().isoformat(),
                "time_since_last_event_ms": (datetime.now() - exec_time).total_seconds() * 1000,
                "is_final_response": event.is_final_response(),
                "content_parts_count": len(event.content.parts) if event.content and event.content.parts else 0,
                "event_details": await _extract_event_details(event)
            }
            exec_time = datetime.now()
            log_entry["events"].append(event_log)
            
            # Process each event and get the final response if available
            response = await process_agent_response(event)
            if response:
                final_response_text = response
                log_entry["final_response"] = response
        
        # Đánh dấu thành công
        log_entry["success"] = True
        
    except Exception as e:
        print(f"Error during agent call: {e}")
        log_entry["error"] = str(e)
        final_response_text = {"text": e.__str__()}
        log_entry["final_response"] = final_response_text
    
    # Tính thời gian xử lý
    end_time = datetime.now()
    processing_time = (end_time - start_time).total_seconds() * 1000  # milliseconds
    log_entry["processing_time_ms"] = processing_time
    
    # Ghi log vào file
    await _write_logs(log_entry)
    
    print("Agent call completed.")
    return final_response_text

async def _extract_event_details(event: Event) -> dict:
    """Trích xuất thông tin chi tiết từ event để log."""
    details = {
        "has_content": bool(event.content),
        "parts_info": []
    }
    
    if event.content and event.content.parts:
        for i, part in enumerate(event.content.parts):
            part_info = {
                "part_index": i,
                "part_type": type(part).__name__,
                "id": event.invocation_id
            }
            
            # Kiểm tra các loại part khác nhau
            if hasattr(part, "text") and part.text:
                part_info["text_length"] = len(part.text)
                part_info["text_preview"] = part.text[:100] + "..." if len(part.text) > 100 else part.text

            if hasattr(part, "function_call") and part.function_call:
                part_info["has_function_call"] = True
                part_info["function_call_name"] = part.function_call.name
                part_info["function_call_args"] = part.function_call.args
                part_info["candidates_token_count"] = event.usage_metadata.candidates_token_count if event.usage_metadata.candidates_token_count else None
                part_info["cached_content_token_count"] = event.usage_metadata.cached_content_token_count if event.usage_metadata.cached_content_token_count else None
                part_info["prompt_token_count"] = event.usage_metadata.prompt_token_count if event.usage_metadata.prompt_token_count else None
                part_info["total_token_count"] = event.usage_metadata.total_token_count if event.usage_metadata.total_token_count else None

            if hasattr(part, "executable_code") and part.executable_code:
                part_info["has_executable_code"] = True
                part_info["code_length"] = len(part.executable_code.code) if part.executable_code.code else 0
            
            if hasattr(part, "code_execution_result") and part.code_execution_result:
                part_info["has_execution_result"] = True
                part_info["execution_outcome"] = str(part.code_execution_result.outcome)
            
            if hasattr(part, "tool_response") and part.tool_response:
                part_info["has_tool_response"] = True
                part_info["tool_output_length"] = len(str(part.tool_response.output)) if part.tool_response.output else 0
            
            if hasattr(part, "function_response") and part.function_response:
                part_info["has_function_response"] = True
                if part.function_response.response:
                    part_info["function_response_keys"] = list(part.function_response.response.keys())
            
            details["parts_info"].append(part_info)
    
    return details

async def _write_logs(log_entry: dict):
    """Ghi log vào file JSON và CSV."""
    try:
        # Tạo tên file với timestamp
        date_str = datetime.now().strftime("%Y%m%d")
        json_filename = f"{LOG_DIR}/agent_calls_{date_str}.json"
        csv_filename = f"{LOG_DIR}/agent_calls_{date_str}.csv"
        
        # Ghi vào file JSON
        await _write_json_log(json_filename, log_entry)
        
        # Ghi vào file CSV
        await _write_csv_log(csv_filename, log_entry)
        
    except Exception as e:
        print(f"Error writing logs: {e}")

async def _write_json_log(filename: str, log_entry: dict):
    """Ghi log vào file JSON."""
    try:
        # Đọc file hiện tại nếu có
        logs = []
        if os.path.exists(filename):
            with open(filename, 'r', encoding='utf-8') as f:
                try:
                    logs = json.load(f)
                except json.JSONDecodeError:
                    logs = []
        
        # Thêm log mới
        logs.append(log_entry)
        
        # Ghi lại file
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(logs, f, indent=2, ensure_ascii=False)
            
    except Exception as e:
        print(f"Error writing JSON log: {e}")

async def _write_csv_log(filename: str, log_entry: dict):
    """Ghi log vào file CSV."""
    try:
        # Chuẩn bị dữ liệu CSV (flatten các thông tin phức tạp)
        csv_row = {
            "timestamp": log_entry["timestamp"],
            "user_id": log_entry["user_id"],
            "session_id": log_entry["session_id"],
            "query": log_entry["query"][:200] + "..." if len(log_entry["query"]) > 200 else log_entry["query"],
            "events_count": len(log_entry["events"]),
            "final_response_type": _get_response_type(log_entry["final_response"]),
            "final_response_text": _get_response_text_preview(log_entry["final_response"]),
            "processing_time_ms": log_entry["processing_time_ms"],
            "success": log_entry["success"],
            "error": log_entry["error"] or ""
        }
        
        # Kiểm tra xem file đã tồn tại chưa
        file_exists = os.path.exists(filename)
        
        # Ghi vào file CSV
        with open(filename, 'a', newline='', encoding='utf-8') as f:
            writer = csv.DictWriter(f, fieldnames=csv_row.keys())
            
            # Ghi header nếu file chưa tồn tại
            if not file_exists:
                writer.writeheader()
            
            writer.writerow(csv_row)
            
    except Exception as e:
        print(f"Error writing CSV log: {e}")

def _get_response_type(response) -> str:
    """Lấy loại response để ghi vào CSV."""
    if not response:
        return "None"
    
    if isinstance(response, dict):
        if "text" in response:
            return "text"
        elif "result" in response:
            return "result"
        else:
            return "dict"
    
    return str(type(response).__name__)

def _get_response_text_preview(response) -> str:
    """Lấy preview của response text để ghi vào CSV."""
    if not response:
        return ""
    
    if isinstance(response, dict):
        if "text" in response:
            text = response["text"]
            return text[:100] + "..." if len(text) > 100 else text
        elif "result" in response:
            return str(response["result"])[:100] + "..."
    
    return str(response)[:100] + "..."



def check_token(token: str) -> str|None:
    """Check if the provided token is valid."""
    
    # Implement your token validation logic here
    # For example, check against a database or an external service
    try:
        # decoded = jwt.decode(token, secret_key, algorithms=["HS256"])
        return None
    except jwt.ExpiredSignatureError as e:
        return f"Token đã hết hạn: {e}"
    except jwt.InvalidTokenError as e:
        return "Token không hợp lệ"
