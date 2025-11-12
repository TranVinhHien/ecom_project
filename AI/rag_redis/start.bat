@echo off
REM Batch script to setup virtual environment and run Vietnamese Vector API

echo === Vietnamese Vector API Setup ===

REM Set virtual environment name
set VENV_NAME=uv

REM Check if virtual environment exists
if not exist "%VENV_NAME%" (
    echo Creating virtual environment: %VENV_NAME%
    python -m venv %VENV_NAME%
    
    if errorlevel 1 (
        echo Error: Failed to create virtual environment
        exit /b 1
    )
    
    echo Virtual environment created successfully
) else (
    echo Virtual environment already exists
)

REM Activate virtual environment
echo Activating virtual environment...
call %VENV_NAME%\Scripts\activate.bat

if errorlevel 1 (
    echo Error: Failed to activate virtual environment
    exit /b 1
)

echo Virtual environment activated

REM Install/upgrade pip
echo Upgrading pip...
python -m pip install --upgrade pip

REM Install requirements
echo Installing requirements...
pip install -r vietnamese_vector_api\requirements.txt

if errorlevel 1 (
    echo Error: Failed to install requirements
    exit /b 1
)

echo Requirements installed successfully

REM Run the application
echo.
echo === Starting Vietnamese Vector API ===
echo Running on http://0.0.0.0:9101
echo Press Ctrl+C to stop
echo.

uvicorn vietnamese_vector_api.main:app --host 0.0.0.0 --port 9101 --reload
