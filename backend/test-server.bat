@echo off
chcp 65001 >nul
title Test Server Connection

echo.
echo ================================
echo   Test Server Connection
echo ================================
echo.

echo [INFO] Testing server connection...
echo [INFO] Server should be running on http://localhost:3000
echo.

REM Test if server is responding
curl -s http://localhost:3000 >nul 2>&1
if %errorlevel% equ 0 (
    echo [OK] Server is running and responding!
    echo [INFO] Opening browser...
    start http://localhost:3000
) else (
    echo [X] Server is not responding
    echo [INFO] Make sure the server is running
)

echo.
pause
