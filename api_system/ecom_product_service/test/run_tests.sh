#!/bin/bash

# Script to run product creation tests
# Usage: ./run_tests.sh [test_name]

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}   Product Creation Test Runner${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Check if server is running
echo -e "${YELLOW}Checking if server is running...${NC}"
if curl -s -o /dev/null -w "%{http_code}" http://172.26.127.95:9001/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Server is running${NC}"
else
    echo -e "${RED}✗ Warning: Server might not be running at http://172.26.127.95:9001${NC}"
    echo -e "${YELLOW}  Continue anyway? (y/n)${NC}"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""

# Change to test directory
cd "$(dirname "$0")"

# If test name is provided, run specific test
if [ ! -z "$1" ]; then
    echo -e "${YELLOW}Running test: $1${NC}"
    echo ""
    go test -v -run "$1" -timeout 5m
else
    echo -e "${YELLOW}Select test to run:${NC}"
    echo "1) TestCreateProduct (Single product with 10 SKUs)"
    echo "2) TestCreateMultipleProducts (3 products batch)"
    echo "3) TestCreateProductWithoutOptionImages (Simple product)"
    echo "4) TestCreateProductWithManyOptions (4 colors x 3 sizes)"
    echo "5) All tests"
    echo "6) Benchmark"
    echo ""
    read -p "Enter choice (1-6): " choice

    case $choice in
        1)
            echo -e "${YELLOW}Running TestCreateProduct...${NC}"
            go test -v -run TestCreateProduct$ -timeout 2m
            ;;
        2)
            echo -e "${YELLOW}Running TestCreateMultipleProducts...${NC}"
            go test -v -run TestCreateMultipleProducts -timeout 5m
            ;;
        3)
            echo -e "${YELLOW}Running TestCreateProductWithoutOptionImages...${NC}"
            go test -v -run TestCreateProductWithoutOptionImages -timeout 2m
            ;;
        4)
            echo -e "${YELLOW}Running TestCreateProductWithManyOptions...${NC}"
            go test -v -run TestCreateProductWithManyOptions -timeout 2m
            ;;
        5)
            echo -e "${YELLOW}Running all tests...${NC}"
            go test -v -timeout 10m
            ;;
        6)
            echo -e "${YELLOW}Running benchmark...${NC}"
            go test -v -bench=BenchmarkCreateProduct -benchmem -timeout 10m
            ;;
        *)
            echo -e "${RED}Invalid choice${NC}"
            exit 1
            ;;
    esac
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   Tests completed!${NC}"
echo -e "${GREEN}========================================${NC}"
