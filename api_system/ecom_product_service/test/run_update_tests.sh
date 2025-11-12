#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║         TEST UPDATE PRODUCT - E-COMMERCE SERVICE          ║${NC}"
echo -e "${CYAN}╚═══════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${YELLOW}⏳ Đang kiểm tra kết nối server...${NC}"
if curl -s -o /dev/null -w "%{http_code}" http://172.26.127.95:9001/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Server đang chạy${NC}"
else
    echo -e "${RED}❌ Không thể kết nối tới server. Vui lòng kiểm tra lại!${NC}"
    exit 1
fi

echo ""
cd "$(dirname "$0")"

if [ ! -z "$1" ]; then
    echo -e "${YELLOW}Đang chạy test: $1${NC}"
    echo ""
    go test -v -run "$1" -timeout 5m
else
    echo -e "${YELLOW}Chọn test để chạy:${NC}"
    echo "1) TestUpdateProductName (Cập nhật chỉ tên)"
    echo "2) TestUpdateProductWithImage (Cập nhật với ảnh mới)"
    echo "3) TestUpdateProductWithMedia (Cập nhật với media files)"
    echo "4) TestUpdateProductOptions (Cập nhật options - thêm option mới)"
    echo "5) TestUpdateProductComplete (Cập nhật toàn bộ thông tin)"
    echo "6) TestUpdateProductDescription (Cập nhật chỉ mô tả)"
    echo "7) TestUpdateProductBrandCategory (Cập nhật brand và category)"
    echo "8) TestUpdateProductMinimal (Test với dữ liệu tối thiểu)"
    echo "9) TestUpdateProductWithNewOption (Thêm option mới - Size)"
    echo "10) TestUpdateProductNameAndDescription (Cập nhật tên và mô tả)"
    echo "11) Chạy TẤT CẢ các test"
    echo "12) Chạy test nhanh (test 1, 2, 6)"
    echo ""
    read -p "Nhập lựa chọn (1-12): " choice

    case $choice in
        1)
            echo -e "${YELLOW}========== TEST 1: CẬP NHẬT CHỈ TÊN ===========${NC}"
            go test -v -run TestUpdateProductName$ -timeout 2m
            ;;
        2)
            echo -e "${YELLOW}========== TEST 2: CẬP NHẬT VỚI ẢNH MỚI ===========${NC}"
            go test -v -run TestUpdateProductWithImage$ -timeout 2m
            ;;
        3)
            echo -e "${YELLOW}========== TEST 3: CẬP NHẬT VỚI MEDIA ===========${NC}"
            go test -v -run TestUpdateProductWithMedia$ -timeout 2m
            ;;
        4)
            echo -e "${YELLOW}========== TEST 4: CẬP NHẬT OPTIONS ===========${NC}"
            go test -v -run TestUpdateProductOptions$ -timeout 2m
            ;;
        5)
            echo -e "${YELLOW}========== TEST 5: CẬP NHẬT TOÀN BỘ ===========${NC}"
            go test -v -run TestUpdateProductComplete$ -timeout 2m
            ;;
        6)
            echo -e "${YELLOW}========== TEST 6: CẬP NHẬT MÔ TẢ ===========${NC}"
            go test -v -run TestUpdateProductDescription$ -timeout 2m
            ;;
        7)
            echo -e "${YELLOW}========== TEST 7: CẬP NHẬT BRAND & CATEGORY ===========${NC}"
            go test -v -run TestUpdateProductBrandCategory$ -timeout 2m
            ;;
        8)
            echo -e "${YELLOW}========== TEST 8: DỮ LIỆU TỐI THIỂU ===========${NC}"
            go test -v -run TestUpdateProductMinimal$ -timeout 2m
            ;;
        9)
            echo -e "${YELLOW}========== TEST 9: THÊM OPTION MỚI ===========${NC}"
            go test -v -run TestUpdateProductWithNewOption$ -timeout 2m
            ;;
        10)
            echo -e "${YELLOW}========== TEST 10: CẬP NHẬT TÊN & MÔ TẢ ===========${NC}"
            go test -v -run TestUpdateProductNameAndDescription$ -timeout 2m
            ;;
        11)
            echo -e "${YELLOW}========== CHẠY TẤT CẢ CÁC TEST ===========${NC}"
            go test -v -timeout 10m
            ;;
        12)
            echo -e "${YELLOW}========== CHẠY TEST NHANH ===========${NC}"
            echo -e "${GREEN}Test 1: Cập nhật tên${NC}"
            go test -v -run TestUpdateProductName$ -timeout 2m
            echo ""
            echo -e "${GREEN}Test 2: Cập nhật với ảnh${NC}"
            go test -v -run TestUpdateProductWithImage$ -timeout 2m
            echo ""
            echo -e "${GREEN}Test 6: Cập nhật mô tả${NC}"
            go test -v -run TestUpdateProductDescription$ -timeout 2m
            ;;
        *)
            echo -e "${RED}Lựa chọn không hợp lệ${NC}"
            exit 1
            ;;
    esac
fi

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    HOÀN THÀNH TEST!                       ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
