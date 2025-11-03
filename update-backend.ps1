# 后端增量更新脚本
# 用途：快速更新服务器上的后端程序

param(
    [string]$ServerIP = "47.242.6.37",
    [string]$ServerUser = "root"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "后端增量更新脚本" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查是否在 backend 目录
if (-not (Test-Path "main.go")) {
    Write-Host "错误: 请在 backend 目录下运行此脚本" -ForegroundColor Red
    Write-Host "当前目录: $(Get-Location)" -ForegroundColor Yellow
    exit 1
}

# 步骤1: 编译 Linux 版本
Write-Host "[1/5] 编译 Linux 版本..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o main.exe main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ 编译失败" -ForegroundColor Red
    exit 1
}
Write-Host "✓ 编译成功" -ForegroundColor Green
Write-Host ""

# 步骤2: 停止服务器上的后端服务
Write-Host "[2/5] 停止服务器后端服务..." -ForegroundColor Yellow
try {
    ssh "${ServerUser}@${ServerIP}" "systemctl stop blog-backend"
    Write-Host "✓ 服务已停止" -ForegroundColor Green
} catch {
    Write-Host "✗ 停止服务失败（可能服务未运行）: $_" -ForegroundColor Yellow
}
Write-Host ""

# 步骤3: 上传文件（使用临时文件名）
Write-Host "[3/5] 上传文件到服务器..." -ForegroundColor Yellow
Write-Host "  服务器: ${ServerUser}@${ServerIP}:/opt/blog/backend/" -ForegroundColor Gray
try {
    scp main.exe "${ServerUser}@${ServerIP}:/opt/blog/backend/main.new"
    Write-Host "✓ 文件上传成功" -ForegroundColor Green
} catch {
    Write-Host "✗ 文件上传失败: $_" -ForegroundColor Red
    Write-Host "  提示: 请确保 SSH 密钥已配置，或手动上传文件" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# 步骤4: 替换文件并设置权限
Write-Host "[4/5] 替换文件并设置权限..." -ForegroundColor Yellow
try {
    ssh "${ServerUser}@${ServerIP}" "cd /opt/blog/backend && mv -f main.new main && chmod +x main"
    Write-Host "✓ 文件替换成功" -ForegroundColor Green
} catch {
    Write-Host "✗ 文件替换失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 步骤5: 启动服务并检查状态
Write-Host "[5/5] 启动后端服务..." -ForegroundColor Yellow
try {
    ssh "${ServerUser}@${ServerIP}" "systemctl start blog-backend"
    Start-Sleep -Seconds 2
    ssh "${ServerUser}@${ServerIP}" "systemctl status blog-backend --no-pager -l"
    Write-Host "✓ 服务启动完成" -ForegroundColor Green
} catch {
    Write-Host "✗ 服务启动失败: $_" -ForegroundColor Red
    Write-Host "  提示: 请手动检查: ssh ${ServerUser}@${ServerIP} 'systemctl status blog-backend'" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "更新完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "查看实时日志: ssh ${ServerUser}@${ServerIP} 'journalctl -u blog-backend -f'" -ForegroundColor Cyan
Write-Host ""

