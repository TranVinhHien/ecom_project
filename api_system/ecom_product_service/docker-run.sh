#!/bin/bash

# Script để build và run Docker container với các biến môi trường

set -e

# Màu sắc cho output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Ecommerce Product Service - Docker Deployment ===${NC}\n"

# Tên image và container
IMAGE_NAME="ecom-product-service"
CONTAINER_NAME="ecom-product-container"

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
    echo -e "${YELLOW}Tạo .env.docker hoặc cung cấp file trước khi chạy.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Environment variables loaded${NC}\n"

# Run container
echo -e "${YELLOW}[4/4] Starting container...${NC}"
sudo docker run -d \
    --name ${CONTAINER_NAME} \
    --env-file .env.docker \
    -p 9001:9001 \
    -v $(pwd)/images:/app/images \
    --restart unless-stopped \
    ${IMAGE_NAME}:latest

echo -e "${GREEN}✓ Container started successfully${NC}\n"

# Hiển thị thông tin container
echo -e "${BLUE}=== Container Information ===${NC}"
echo -e "${GREEN}Container Name:${NC} ${CONTAINER_NAME}"
echo -e "${GREEN}Image:${NC} ${IMAGE_NAME}:latest"
echo -e "${GREEN}Port:${NC} 9001"
echo -e "${GREEN}Status:${NC}"
sudo docker ps --filter "name=${CONTAINER_NAME}" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo -e "\n${BLUE}=== Useful Commands ===${NC}"
echo -e "${GREEN}Xem logs:${NC} docker logs -f ${CONTAINER_NAME}"
echo -e "${GREEN}Dừng container:${NC} docker stop ${CONTAINER_NAME}"
echo -e "${GREEN}Khởi động lại:${NC} docker restart ${CONTAINER_NAME}"
echo -e "${GREEN}Xóa container:${NC} docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}"
echo -e "${GREEN}Vào shell container:${NC} docker exec -it ${CONTAINER_NAME} sh"

echo -e "\n${GREEN}✓ Deployment completed! Service đang chạy tại http://0.0.0.0:9001${NC}"
