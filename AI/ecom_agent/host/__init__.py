from fastapi import FastAPI, HTTPException,Request
from pydantic import BaseModel
from typing import Optional, Dict, Any
from  .agent_main import HostAgent, session_service
import uuid
from .util import call_agent_async,check_token
from contextlib import asynccontextmanager
from  .call_api import get_user_info
from dotenv import load_dotenv
import jwt
from fastapi.middleware.cors import CORSMiddleware
import logging

load_dotenv()
AGENT_NAME = "Host_Agent"

print("initializing host agent")
host = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    global host
    print("Initializing host agent...")
    host = await HostAgent.create(name=AGENT_NAME)
    print("HostAgent initialized successfully")
    
    yield
    print("Shutting down host agent...")



print("HostAgent initialized")

app = FastAPI(
    title="Agent Host API",
    description="API for managing agent sessions and messages",
    version="1.0.0",
    lifespan=lifespan

)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

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

# Pydantic models for request/response
class ErrorResponse(BaseModel):
    success: bool
    error: str

class CreateSessionRequest(BaseModel):
    user_id: Optional[str] = None
    state: Optional[Dict[str, Any]] = {}

class CreateSessionResponse(BaseModel):
    success: bool
    session_id: str
    message: str

class SendMessageRequest(BaseModel):
    message: str
    user_id: Optional[str] = None
    session_id: Optional[str] = None

class SendMessageResponse(BaseModel):
    success: bool
    response: Optional[dict[str, Any]] = None
    session_id: str
    error: Optional[str] = None

class HealthResponse(BaseModel):
    status: str
    agent_name: str

class ErrorResponse(BaseModel):
    success: bool
    error: str

@app.post("/api/session", response_model=CreateSessionResponse)
async def create_session(request: CreateSessionRequest,raw_request: Request):
    """API để tạo session mới"""
    try:
        # Tạo session_id mới nếu không được cung cấp
        session_id =  str(uuid.uuid4())
        app_name = AGENT_NAME
        lang = request.state.get("lang","VN") 
        # state = request.state or {}
        headers = raw_request.headers
        token = headers.get("Authorization", "").replace("Bearer ", "")
        # check token
        if not token:
            return HTTPException(
                status_code=401,
                detail="Authorization token is missing"
            )
        user_info = get_user_info(token)
        user_id = user_info.get("userId")
        if not user_id or not user_info:
            return HTTPException(
                status_code=400,
                detail="user_id and user_info are required"
            )

        # Tạo session
        new_session = await session_service.create_session(
            app_name=app_name,
            user_id=user_id,
            state={
                "user_info":user_info,
                "lang":lang
                },
            session_id=session_id,
        )
        return CreateSessionResponse(
            success=True,
            session_id=session_id,
            message="Session created successfully",
            session_data=new_session
        )
        
    except Exception as e:
        return HTTPException(
            status_code=500,
            detail=f"Failed to create session: {str(e)}"
        )

@app.post("/api/message", response_model=SendMessageResponse)
async def send_message(request: SendMessageRequest,raw_request: Request):
    """API để gửi tin nhắn cho agent"""
    headers = raw_request.headers
    token = headers.get("Authorization", "").replace("Bearer ", "")
    session_id = request.session_id or ""
    if check := check_token(token):
        return SendMessageResponse(
            success=False,
            error=check,
            session_id=session_id
        )
    if not request.message:
        return SendMessageResponse(
            success=False,
            error="Message is required",
            session_id=session_id
        )

    try:
        user_id = jwt.decode(token, options={"verify_signature": False}).get("userId")
        if not session_id:
            return SendMessageResponse(
                success=False,
                error="session_id is required",
                session_id=session_id
            )
        if not user_id:
            return SendMessageResponse(
                success=False,
                error="user_id is invalid in token, please login again",
                session_id=session_id
            )

        response = await call_agent_async(host.runner, user_id, session_id, request.message, token)
        if not response:
            return SendMessageResponse(
                success=False,
                error="Lỗi hệ thống. Vui lòng thử lại",
                session_id=session_id
            )

        return SendMessageResponse(
            success=True,
            response=response,
            session_id=session_id
        )

    except Exception as e:
        logging.exception("Failed to send message", exc_info=e)
        return SendMessageResponse(
            success=False,
            error=f"Lỗi hệ thống. Vui lòng thử lại: {str(e)}",
            session_id=session_id
        )

@app.get("/api/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        agent_name=AGENT_NAME
    )

@app.get("/api/session/{session_id}")
async def get_session(session_id: str, raw_request: Request):
    """Lấy thông tin session"""
    headers = raw_request.headers
    token = headers.get("Authorization", "").replace("Bearer ", "")
    if check := check_token(token):
        return {
            "success": False,
            "status_code": 401,
            "detail": check
        }
    user_id = jwt.decode(token, options={"verify_signature": False}).get("userId")
    if not user_id:
        return {
            "success": False,
            "status_code": 400,
            "detail": "user_id is error in token, please login again"
        }
    try:
        session =await session_service.get_session(session_id=session_id,app_name=AGENT_NAME,user_id=user_id)
        if not session:
            return {
                "success": False,
                "status_code": 404,
                "detail": "Session not found"
            }

        return {
            "success": True,
            "session_id": session_id,
            "session_data": session
        }
        
    except Exception as e:
        return {
            "success": False,
            "status_code": 500,
            "detail": f"Failed to get session: {str(e)}"
        }
@app.get("/api/list_sessions")
async def list_sessions(raw_request: Request):
    """Lấy danh sách các session của người dùng"""
    headers = raw_request.headers
    token = headers.get("Authorization", "").replace("Bearer ", "")
    if check := check_token(token):
        return ErrorResponse(
            success=False,
            error=check
        )
    user_id = jwt.decode(token, options={"verify_signature": False}).get("userId")
    if not user_id:
        return ErrorResponse(
            success=False,
            error="user_id is error in token, please login again"
        )

    try:
        sessions = await session_service.list_sessions(user_id=user_id, app_name=AGENT_NAME)
        list_sessions = [s.id for s in sessions.sessions]
        return {
            "success": True,
            "sessions": list_sessions,
            "message": f"Found {len(list_sessions)} sessions for user {user_id}"
        }
    except Exception as e:
        return {
            "success": False,
            "status_code": 500,
            "detail": f"Failed to list sessions: {str(e)}"
        }

@app.delete("/api/session/{session_id}")
async def delete_session( session_id: str, raw_request: Request):
    """Xóa session"""
    headers = raw_request.headers
    token = headers.get("Authorization", "").replace("Bearer ", "")
    if check := check_token(token):
        return {
            "success": False,
            "status_code": 401,
            "detail": check
        }
    user_id = jwt.decode(token, options={"verify_signature": False}).get("userId")
    if not user_id:
        return {
            "success": False,
            "status_code": 400,
            "detail": "user_id is error in token, please login again"
        }
    try:
        result = await session_service.delete_session(app_name=AGENT_NAME,session_id=session_id,user_id=user_id)

        return {
            "success": True,
            "message": f"Session {session_id} deleted successfully",
            "result": result
        }
        
    except Exception as e:
        return {
            "success": False,
            "status_code": 500,
            "detail": f"Failed to delete session: {str(e)}"
        }
