@echo off
chcp 65001 >nul
title Blog Backend Server - Development Mode

echo.
echo ================================
echo   Blog Backend Development Server
echo ================================
echo.
echo [INFO] Using Air hot-reload tool
echo [INFO] Press Ctrl+C to stop server
echo.

REM Check if Air is installed
where air >nul 2>&1
if %errorlevel% neq 0 (
    echo [!] Air not installed, installing...
    go install github.com/air-verse/air@latest
    if %errorlevel% neq 0 (
        echo [X] Air installation failed
        pause
        exit /b 1
    )
    echo [OK] Air installed successfully
)

echo [OK] Starting Air hot-reload mode...
echo.

REM Start Air
air
