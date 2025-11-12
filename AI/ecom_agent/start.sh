#!/bin/bash

# Script để cài đặt uv, tạo môi trường và chạy ứng dụng FastAPI

set -e  # Dừng script nếu có lỗi

echo "=========================================="
echo "Bắt đầu thiết lập môi trường..."
echo "=========================================="

# 1. Kiểm tra và cài đặt uv nếu chưa có
if ! command -v uv &> /dev/null; then
    echo "uv chưa được cài đặt. Đang cài đặt uv..."
    curl -LsSf https://astral.sh/uv/install.sh | sh
    
    # Thêm uv vào PATH cho session hiện tại
    export PATH="$HOME/.cargo/bin:$PATH"
    
    echo "✓ Đã cài đặt uv thành công!"
else
    echo "✓ uv đã được cài đặt"
fi

# 2. Tạo môi trường ảo với uv (nếu chưa có)
if [ ! -d ".venv" ]; then
    echo "Đang tạo môi trường ảo với uv..."
    uv venv
    echo "✓ Đã tạo môi trường ảo"
else
    echo "✓ Môi trường ảo đã tồn tại"
fi

# 3. Kích hoạt môi trường ảo
echo "Đang kích hoạt môi trường ảo..."
source .venv/bin/activate

# 4. Cài đặt các thư viện từ pyproject.toml hoặc requirements.txt
echo "Đang cài đặt các thư viện..."
if [ -f "pyproject.toml" ]; then
    echo "Sử dụng pyproject.toml..."
    uv pip install -e .
elif [ -f "requirements.txt" ]; then
    echo "Sử dụng requirements.txt..."
    uv pip install -r requirements.txt
else
    echo "⚠ Không tìm thấy pyproject.toml hoặc requirements.txt"
    exit 1
fi

echo "✓ Đã cài đặt thư viện thành công!"

# 5. Chạy ứng dụng với uvicorn
echo "=========================================="
echo "Đang khởi động ứng dụng FastAPI..."
echo "=========================================="

# Sử dụng uvicorn để chạy ứng dụng
# Dựa vào run_server.py, module chính là host:app
uv run uvicorn host:app --host 0.0.0.0 --port 9102 --reload

# Hoặc nếu bạn muốn chạy file run_server.py:
# uv run python run_server.py
