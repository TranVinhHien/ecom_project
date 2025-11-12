# PowerShell Script để cài đặt uv, tạo môi trường và chạy ứng dụng FastAPI

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Bắt đầu thiết lập môi trường..." -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. Kiểm tra và cài đặt uv nếu chưa có
try {
    $uvVersion = uv --version
    Write-Host "✓ uv đã được cài đặt: $uvVersion" -ForegroundColor Green
} catch {
    Write-Host "uv chưa được cài đặt. Đang cài đặt uv..." -ForegroundColor Yellow
    
    # Cài đặt uv trên Windows
    irm https://astral.sh/uv/install.ps1 | iex
    
    # Refresh PATH
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    Write-Host "✓ Đã cài đặt uv thành công!" -ForegroundColor Green
}

# 2. Tạo môi trường ảo với uv (nếu chưa có)
if (-Not (Test-Path ".venv")) {
    Write-Host "Đang tạo môi trường ảo với uv..." -ForegroundColor Yellow
    uv venv
    Write-Host "✓ Đã tạo môi trường ảo" -ForegroundColor Green
} else {
    Write-Host "✓ Môi trường ảo đã tồn tại" -ForegroundColor Green
}

# 3. Cài đặt các thư viện từ pyproject.toml hoặc requirements.txt
Write-Host "Đang cài đặt các thư viện..." -ForegroundColor Yellow

if (Test-Path "pyproject.toml") {
    Write-Host "Sử dụng pyproject.toml..." -ForegroundColor Cyan
    uv pip install -e .
} elseif (Test-Path "requirements.txt") {
    Write-Host "Sử dụng requirements.txt..." -ForegroundColor Cyan
    uv pip install -r requirements.txt
} else {
    Write-Host "⚠ Không tìm thấy pyproject.toml hoặc requirements.txt" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Đã cài đặt thư viện thành công!" -ForegroundColor Green

# 4. Chạy ứng dụng với uvicorn
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Đang khởi động ứng dụng FastAPI..." -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# Sử dụng uvicorn để chạy ứng dụng
# Dựa vào run_server.py, module chính là host:app
uv run uvicorn host:app --host 0.0.0.0 --port 9102 --reload

# Hoặc nếu bạn muốn chạy file run_server.py:
# uv run python run_server.py
