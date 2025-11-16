set -e

# Màu sắc cho output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Ecommerce Payment Service - Docker Deployment ===${NC}\n"

# Tên image và container
IMAGE_NAME="ecom-payment-service"
CONTAINER_NAME="ecom-payment-container"

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
    exit 1
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
    -p 9003:9003 \
    --restart unless-stopped \
    ${IMAGE_NAME}:latest

echo -e "${GREEN}✓ Container started successfully${NC}\n"

# Hiển thị thông tin container
echo -e "${BLUE}=== Container Information ===${NC}"
echo -e "${GREEN}Container Name:${NC} ${CONTAINER_NAME}"
echo -e "${GREEN}Image:${NC} ${IMAGE_NAME}:latest"
echo -e "${GREEN}Port:${NC} 9003"
echo -e "${GREEN}Network:${NC} e-commerce-network"
echo -e "${GREEN}Status:${NC}"
sudo docker ps --filter "name=${CONTAINER_NAME}" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo -e "\n${BLUE}=== Service Dependencies ===${NC}"
echo -e "${GREEN}MySQL:${NC} 172.26.127.95:3306 (ecommerce_transacion_db)"
echo -e "${GREEN}Redis:${NC} 172.26.127.95:6379"
echo -e "${GREEN}Kafka:${NC} 172.26.127.95:9092"
echo -e "${GREEN}Product Service:${NC} 172.26.127.95:9001"
echo -e "${GREEN}Order Service:${NC} 172.26.127.95:9002"

echo -e "\n${BLUE}=== Payment Gateway Configuration ===${NC}"
echo -e "${GREEN}MoMo Endpoint:${NC} https://test-payment.momo.vn/v2/gateway/api/create"
echo -e "${GREEN}Redirect URL:${NC} http://localhost:9999/vi/dat-hang-thanh-cong"
echo -e "${GREEN}IPN URL:${NC} https://51c3b9baa7ac.ngrok-free.app/v1/transaction/callback"
echo -e "${GREEN}Public ID:${NC} https://51c3b9baa7ac.ngrok-free.app/v1"

echo -e "\n${BLUE}=== Email Configuration ===${NC}"
echo -e "${GREEN}Email Service:${NC} Brevo API"
echo -e "${GREEN}Sender Email:${NC} hienlazada1912@gmail.com"
echo -e "${GREEN}Sender Name:${NC} lemarchenoble"

echo -e "\n${BLUE}=== Useful Commands ===${NC}"
echo -e "${GREEN}Xem logs:${NC} docker logs -f ${CONTAINER_NAME}"
echo -e "${GREEN}Xem logs realtime:${NC} docker logs -f --tail 100 ${CONTAINER_NAME}"
echo -e "${GREEN}Dừng container:${NC} docker stop ${CONTAINER_NAME}"
echo -e "${GREEN}Khởi động lại:${NC} docker restart ${CONTAINER_NAME}"
echo -e "${GREEN}Xóa container:${NC} docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}"
echo -e "${GREEN}Vào shell container:${NC} docker exec -it ${CONTAINER_NAME} sh"
echo -e "${GREEN}Kiểm tra health:${NC} curl http://localhost:9003/health"

echo -e "\n${YELLOW}⚠️  LƯU Ý:${NC}"
echo -e "${YELLOW}1. Cập nhật IPNURL khi thay đổi ngrok domain${NC}"
echo -e "${YELLOW}2. Đảm bảo MySQL, Redis, Kafka đang chạy${NC}"
echo -e "${YELLOW}3. Kafka consumer group: ecom-payment-service-group${NC}"

echo -e "\n${GREEN}✓ Deployment completed! Payment Service đang chạy tại http://0.0.0.0:9003${NC}"
