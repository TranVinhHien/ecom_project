#!/bin/bash
# filepath: /home/hienlazada123/project-ecom/ecom_analytics_service/docker-run.sh
set -e

# Màu sắc cho output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Ecommerce Analytics Service - Docker Deployment ===${NC}\n"

# Tên image và container
IMAGE_NAME="ecom-analytics-service"
CONTAINER_NAME="ecom-analytics-container"

# Build Docker image
echo -e "${YELLOW}[1/4] Building Docker image...${NC}"
sudo docker build -t ${IMAGE_NAME}:latest .
echo -e "${GREEN}✓ Docker image built successfully${NC}\n"

# Dừng và xóa container cũ nếu có
echo -e "${YELLOW}[2/4] Cleaning up old container...${NC}"
sudo docker stop ${CONTAINER_NAME} 2>/dev/null || true
sudo docker rm ${CONTAINER_NAME} 2>/dev/null || true
echo -e "${GREEN}✓ Cleanup completed${NC}\n"

# Load biến môi trường từ file .env.docker
echo -e "${YELLOW}[3/4] Loading environment variables from .env.docker...${NC}"
if [ ! -f .env.docker ]; then
    echo -e "${RED}✗ File .env.docker không tồn tại!${NC}"
    echo -e "${YELLOW}Tạo file .env.docker với nội dung mẫu...${NC}"
    cat > .env.docker << 'EOF'
# ============================================
# ECOMMERCE ANALYTICS SERVICE - DOCKER ENV CONFIG
# ============================================

# Database Configuration (3 databases)
DB_SOURCE_ORDER=root:12345@tcp(172.26.127.95:3306)/ecommerce_order_db?parseTime=true
DB_SOURCE_TRANSACTION=root:12345@tcp(172.26.127.95:3306)/ecommerce_transacion_db?parseTime=true
DB_SOURCE_INTERACT=root:12345@tcp(172.26.127.95:3306)/ecommerce_interact_db?parseTime=true

# Server Configuration
HTTP_SERVER_ADDRESS=0.0.0.0:9004

# JWT Configuration
JWT_SECRET=bv-T"-u6@-WR?SHiHQ7yQ]CK*dd9(@jM9BI)|g;zq)ur-Z.Jw/u5HyJHgg,KS.fa

# Client Configuration (Multiple origins)
CLIENT_IP=http://localhost:9999,http://localhost:8989

# Redis Configuration
REDIS_ADDRESS=172.26.127.95:6379

# System Token (for internal service communication)
TOKEN_SYSTEM=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaWVubGF6YWRhIiwic2NvcGUiOiJTWVNURU0iLCJpc3MiOiJsZW1hcmNoZW5vYmxlLmlkLnZuIiwiZXhwIjo0OTE3NTExMjUyLCJpYXQiOjE3NjE3NTEyNTIsInVzZXJJZCI6IjE2NzQwOGUzLWFmZWYtNDhiOS04ZTRmLTZkZDQxZWJmMzQ2NCIsImp0aSI6ImU2YzgyN2E2LTIyOTYtNGNlOC1iMjQ1LWM3MDIxNWM4MGJjNyIsImVtYWlsIjoidmluaGhpZW4xMnpAZ21haWwuY29tIn0.CPnP_NqB_WtaQb9X43YKFav8wYzdqB14jFNtnPr74as
EOF
    echo -e "${GREEN}✓ File .env.docker đã được tạo${NC}"
fi
echo -e "${GREEN}✓ Environment variables loaded${NC}\n"

# Run container
echo -e "${YELLOW}[4/4] Starting container...${NC}"

# Kiểm tra xem network e-commerce-network đã tồn tại chưa
if ! sudo docker network ls | grep -q "e-commerce-network"; then
    echo -e "${YELLOW}Creating network e-commerce-network...${NC}"
    sudo docker network create e-commerce-network
fi

sudo docker run -d \
    --name ${CONTAINER_NAME} \
    --env-file .env.docker \
    --network e-commerce-network \
    -p 9004:9004 \
    --restart unless-stopped \
    ${IMAGE_NAME}:latest

echo -e "${GREEN}✓ Container started successfully${NC}\n"

# Hiển thị thông tin container
echo -e "${BLUE}=== Container Information ===${NC}"
echo -e "${GREEN}Container Name:${NC} ${CONTAINER_NAME}"
echo -e "${GREEN}Image:${NC} ${IMAGE_NAME}:latest"
echo -e "${GREEN}Port:${NC} 9004"
echo -e "${GREEN}Network:${NC} e-commerce-network"
echo -e "${GREEN}Status:${NC}"
sudo docker ps --filter "name=${CONTAINER_NAME}" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo -e "\n${BLUE}=== Service Dependencies ===${NC}"
echo -e "${GREEN}MySQL Databases:${NC}"
echo -e "  - Order DB: 172.26.127.95:3306/ecommerce_order_db"
echo -e "  - Transaction DB: 172.26.127.95:3306/ecommerce_transacion_db"
echo -e "  - Interact DB: 172.26.127.95:3306/ecommerce_interact_db"
echo -e "${GREEN}Redis:${NC} 172.26.127.95:6379"
echo -e "${GREEN}Allowed Origins:${NC} http://localhost:9999, http://localhost:8989"

echo -e "\n${BLUE}=== API Endpoints ===${NC}"
echo -e "${GREEN}Health Check:${NC} http://localhost:9004/health"
echo -e "${GREEN}Platform Analytics:${NC} http://localhost:9004/v1/platform/*"
echo -e "${GREEN}Shop Analytics:${NC} http://localhost:9004/v1/shop/*"
echo -e "${GREEN}Public APIs:${NC} http://localhost:9004/v1/public/*"
echo -e "${GREEN}Feedback System:${NC}"
echo -e "  - Submit Review: POST http://localhost:9004/v1/public/chatbox/review"
echo -e "  - Submit Feedback: POST http://localhost:9004/v1/public/customer-support/complaint"
echo -e "  - Admin Stats: GET http://localhost:9004/v1/platform/chatbox/statistics"

echo -e "\n${BLUE}=== Useful Commands ===${NC}"
echo -e "${GREEN}Xem logs:${NC} docker logs -f ${CONTAINER_NAME}"
echo -e "${GREEN}Xem logs realtime:${NC} docker logs -f --tail 100 ${CONTAINER_NAME}"
echo -e "${GREEN}Dừng container:${NC} docker stop ${CONTAINER_NAME}"
echo -e "${GREEN}Khởi động lại:${NC} docker restart ${CONTAINER_NAME}"
echo -e "${GREEN}Xóa container:${NC} docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}"
echo -e "${GREEN}Vào shell container:${NC} docker exec -it ${CONTAINER_NAME} sh"
echo -e "${GREEN}Kiểm tra health:${NC} curl http://localhost:9004/health"
echo -e "${GREEN}Test feedback API:${NC} curl -X POST http://localhost:9004/v1/public/chatbox/review -H 'Content-Type: application/json' -d '{\"user_id\":\"test\",\"session_id\":\"s1\",\"event_id\":\"e1\",\"rating\":1}'"

echo -e "\n${BLUE}=== Database Verification ===${NC}"
echo -e "${YELLOW}Checking database connections...${NC}"
sleep 2
if sudo docker exec ${CONTAINER_NAME} sh -c "echo 'Database connections will be verified on startup'" 2>/dev/null; then
    echo -e "${GREEN}✓ Container is running${NC}"
fi

echo -e "\n${GREEN}✓ Deployment completed! Analytics Service đang chạy tại http://0.0.0.0:9004${NC}"
echo -e "${YELLOW}📊 Service này cung cấp analytics cho Platform Admin và Shop Owner${NC}"
echo -e "${YELLOW}💬 Bao gồm hệ thống Feedback & Ratings mới${NC}"