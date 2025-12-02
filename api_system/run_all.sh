#!/bin/bash

# Script để chạy toàn bộ hệ thống Microservices
# Tác giả: [Tên của bạn]

# Dừng script nếu có bất kỳ lỗi nào xảy ra
set -e

# Lấy thư mục hiện tại
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Màu sắc hiển thị
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}==============================================${NC}"
echo -e "${BLUE}   STARTING E-COMMERCE MICROSERVICES SYSTEM   ${NC}"
echo -e "${BLUE}==============================================${NC}\n"

# Danh sách các thư mục service (Dựa trên ảnh bạn cung cấp)
# Lưu ý: ecom_payment_service dùng dấu gạch ngang (-), các cái khác dùng gạch dưới (_)
SERVICES=(
    "ecom_analytics_service"
    "ecom_order_service"
    "ecom_payment_service"
    "ecom_product_service"
)

# Hàm để chạy service
run_service_script() {
    local service_dir=$1
    local full_path="${SCRIPT_DIR}/${service_dir}"

    echo -e "${YELLOW}>>> Processing: ${service_dir} ...${NC}"

    if [ -d "$full_path" ]; then
        # Kiểm tra xem file script con có tồn tại không
        if [ -f "$full_path/docker-run.sh" ]; then
            echo -e "Found docker-run.sh, executing..."
            
            # Cấp quyền thực thi cho file con (đề phòng chưa có quyền)
            chmod +x "$full_path/docker-run.sh"
            
            # Chạy script con
            # Dùng ( ) để tạo sub-shell, giúp lệnh cd không ảnh hưởng đến script chính
            (cd "$full_path" && ./docker-run.sh)
            
            echo -e "${GREEN}✓ Service ${service_dir} deployed successfully!${NC}\n"
            echo -e "--------------------------------------------------\n"
        else
            echo -e "${RED}✗ Error: File docker-run.sh not found in ${service_dir}${NC}\n"
            exit 1
        fi
    else
        echo -e "${RED}✗ Error: Directory ${service_dir} does not exist${NC}\n"
        exit 1
    fi
}

# Vòng lặp chạy qua từng service
for service in "${SERVICES[@]}"; do
    run_service_script "$service"
    # Nghỉ 2 giây giữa các lần build để dễ theo dõi log và tránh quá tải tức thời
    sleep 2
done

echo -e "${BLUE}==============================================${NC}"
echo -e "${GREEN}   ALL SERVICES HAVE BEEN DEPLOYED!           ${NC}"
echo -e "${BLUE}==============================================${NC}"

# Liệt kê tất cả container đang chạy thuộc dự án (lọc theo tên ecom)
echo -e "\n${YELLOW}Current Running Containers:${NC}"
sudo docker ps --filter "name=ecom" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"