@echo off
chcp 65001 >nul
title 博客后端服务器

echo.
echo ================================
echo   博客后端快速启动工具
echo ================================
echo.

REM 检查是否存在编译好的文件
if exist blog.exe (
    echo [√] 发现已编译文件，直接启动...
    echo.
    blog.exe
) else (
    echo [!] 未找到编译文件，开始编译...
    echo [提示] 首次编译需要 10-30 秒，请耐心等待
    echo.
    go build -o blog.exe main.go
    if %errorlevel% equ 0 (
        echo.
        echo [√] 编译成功！启动服务器...
        echo [提示] 下次启动将秒开！
        echo.
        blog.exe
    ) else (
        echo.
        echo [×] 编译失败，请检查错误信息
        pause
    )
)


