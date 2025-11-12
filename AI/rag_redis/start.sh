#!/bin/bash

# Script to setup virtual environment and run Vietnamese Vector API

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Vietnamese Vector API Setup ===${NC}"

# Set virtual environment name
VENV_NAME="uv"

# Check if virtual environment exists
if [ ! -d "$VENV_NAME" ]; then
    echo -e "${YELLOW}Creating virtual environment: $VENV_NAME${NC}"
    python3 -m venv $VENV_NAME
    if [ $? -ne 0 ]; then
        echo "Error: Failed to create virtual environment"
        exit 1
    fi
    echo -e "${GREEN}Virtual environment created successfully${NC}"
else
    echo -e "${GREEN}Virtual environment already exists${NC}"
fi

# Activate virtual environment
echo -e "${YELLOW}Activating virtual environment...${NC}"
source $VENV_NAME/bin/activate

if [ $? -ne 0 ]; then
    echo "Error: Failed to activate virtual environment"
    exit 1
fi

echo -e "${GREEN}Virtual environment activated${NC}"

# Install/upgrade pip
echo -e "${YELLOW}Upgrading pip...${NC}"
pip install --upgrade pip

# Install requirements
echo -e "${YELLOW}Installing requirements...${NC}"
pip install -r vietnamese_vector_api/requirements.txt

if [ $? -ne 0 ]; then
    echo "Error: Failed to install requirements"
    exit 1
fi

echo -e "${GREEN}Requirements installed successfully${NC}"

# Run the application
echo -e "${GREEN}=== Starting Vietnamese Vector API ===${NC}"
echo -e "${YELLOW}Running on http://0.0.0.0:9101${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"

uvicorn vietnamese_vector_api.main:app --host 0.0.0.0 --port 9101 --reload
