#!/bin/bash

# ==============================================================================
# RUN API TESTS
# ==============================================================================

cd "$(dirname "$0")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}E-Commerce Analytics API Test Runner${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Show menu
show_menu() {
    echo "Chọn test suite để chạy:"
    echo "  1. Test tất cả APIs (all)"
    echo "  2. Test Shop APIs only (shop)"
    echo "  3. Test Platform APIs only (platform)"
    echo "  4. Test Edge Cases only (edge)"
    echo "  0. Exit"
    echo ""
}

# Run test
run_test() {
    local suite=$1
    echo -e "${YELLOW}Running test suite: $suite${NC}"
    echo ""
    go run api_tests.go "$suite"
    exit_code=$?
    echo ""
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✓ Test suite '$suite' completed successfully!${NC}"
    else
        echo -e "${RED}✗ Test suite '$suite' failed!${NC}"
    fi
    
    return $exit_code
}

# Main
if [ $# -eq 0 ]; then
    # Interactive mode
    while true; do
        show_menu
        read -p "Nhập lựa chọn: " choice
        echo ""
        
        case $choice in
            1)
                run_test "all"
                ;;
            2)
                run_test "shop"
                ;;
            3)
                run_test "platform"
                ;;
            4)
                run_test "edge"
                ;;
            0)
                echo "Goodbye!"
                exit 0
                ;;
            *)
                echo -e "${RED}Lựa chọn không hợp lệ!${NC}"
                ;;
        esac
        
        echo ""
        read -p "Press Enter để tiếp tục..."
        clear
    done
else
    # Command line mode
    run_test "$1"
fi
