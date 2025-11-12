# PowerShell script to setup virtual environment and run Vietnamese Vector API

Write-Host "=== Vietnamese Vector API Setup ===" -ForegroundColor Green

# Set virtual environment name
$VENV_NAME = "uv"

# Check if virtual environment exists
if (-Not (Test-Path $VENV_NAME)) {
    Write-Host "Creating virtual environment: $VENV_NAME" -ForegroundColor Yellow
    python -m venv $VENV_NAME
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Failed to create virtual environment" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Virtual environment created successfully" -ForegroundColor Green
} else {
    Write-Host "Virtual environment already exists" -ForegroundColor Green
}

# Activate virtual environment
Write-Host "Activating virtual environment..." -ForegroundColor Yellow
& "$VENV_NAME\Scripts\Activate.ps1"

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Failed to activate virtual environment" -ForegroundColor Red
    Write-Host "Note: You may need to run 'Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser' first" -ForegroundColor Yellow
    exit 1
}

Write-Host "Virtual environment activated" -ForegroundColor Green

# Install/upgrade pip
Write-Host "Upgrading pip..." -ForegroundColor Yellow
python -m pip install --upgrade pip

# Install requirements
Write-Host "Installing requirements..." -ForegroundColor Yellow
pip install -r vietnamese_vector_api\requirements.txt

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Failed to install requirements" -ForegroundColor Red
    exit 1
}

Write-Host "Requirements installed successfully" -ForegroundColor Green

# Run the application
Write-Host "`n=== Starting Vietnamese Vector API ===" -ForegroundColor Green
Write-Host "Running on http://0.0.0.0:9101" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop`n" -ForegroundColor Yellow

uvicorn vietnamese_vector_api.main:app --host 0.0.0.0 --port 9101 --reload
