@echo off
chcp 65001 >nul
title 博客后端服务器 - 优化版

echo.
echo ================================
echo   博客后端快速启动工具 v2.0
echo ================================
echo.

REM 检查是否存在编译好的文件
if exist main.exe (
    echo [√] 发现已编译文件，直接启动...
    echo [提示] 跳过编译步骤，秒开启动！
    echo.
    main.exe
) else if exist blog.exe (
    echo [√] 发现备用编译文件，直接启动...
    echo.
    blog.exe
) else (
    echo [!] 未找到编译文件，开始编译...
    echo [提示] 首次编译需要 10-30 秒，请耐心等待
    echo [提示] 编译完成后将自动启动，下次启动将秒开！
    echo.
    
    REM 使用更快的编译选项
    go build -ldflags="-s -w" -o main.exe main.go
    if %errorlevel% equ 0 (
        echo.
        echo [√] 编译成功！启动服务器...
        echo [提示] 下次启动将秒开！
        echo.
        main.exe
    ) else (
        echo.
        echo [×] 编译失败，请检查错误信息
        echo [提示] 请确保：
        echo   - Go环境已正确安装
        echo   - 数据库服务正在运行
        echo   - 配置文件存在且正确
        pause
    )
)

echo.
echo 服务器已停止运行
pause
