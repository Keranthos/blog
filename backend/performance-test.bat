@echo off
chcp 65001 >nul
title 博客后端性能测试

echo.
echo ================================
echo   博客后端启动性能测试
echo ================================
echo.

echo [测试1] 直接运行已编译文件...
if exist main.exe (
    echo 开始时间: %time%
    start /wait main.exe
    echo 结束时间: %time%
) else (
    echo [×] main.exe 不存在
)

echo.
echo [测试2] 使用 go run 编译运行...
echo 开始时间: %time%
timeout /t 2 >nul
go run main.go
echo 结束时间: %time%

echo.
echo [提示] 性能对比：
echo   - 直接运行可执行文件：秒开
echo   - go run 编译运行：需要10-30秒
echo.
pause
