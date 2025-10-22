@echo off
echo 正在编译...
go build -o blog.exe main.go
if %errorlevel% neq 0 (
    echo 编译失败！
    pause
    exit /b %errorlevel%
)
echo 编译成功！启动服务器...
blog.exe


