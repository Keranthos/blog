# 本地部署脚本（方案A：本地构建 + 服务器部署）
# 使用方法：.\deploy.ps1

param(
    [string]$ServerIP = "47.242.6.37",
    [string]$ServerUser = "root",
    [switch]$SkipBuild = $false
)

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "博客部署脚本 (方案A: 本地构建)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查必要工具
Write-Host "[1/7] 检查环境..." -ForegroundColor Yellow
if (-not (Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 node 命令，请先安装 Node.js" -ForegroundColor Red
    exit 1
}
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 go 命令，请先安装 Go" -ForegroundColor Red
    exit 1
}
Write-Host "✓ 环境检查通过" -ForegroundColor Green
Write-Host ""

# 构建前端
if (-not $SkipBuild) {
    Write-Host "[2/7] 构建前端..." -ForegroundColor Yellow
    Push-Location frontend
    try {
        # 检查是否需要安装依赖
        if (-not (Test-Path "node_modules")) {
            Write-Host "  安装前端依赖..." -ForegroundColor Gray
            npm install
        }
        Write-Host "  正在构建..." -ForegroundColor Gray
        npm run build
        if ($LASTEXITCODE -ne 0) {
            throw "前端构建失败"
        }
        Write-Host "✓ 前端构建完成" -ForegroundColor Green
    }
    catch {
        Write-Host "✗ 前端构建失败: $_" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    finally {
        Pop-Location
    }
    Write-Host ""
} else {
    Write-Host "[2/7] 跳过前端构建" -ForegroundColor Gray
    Write-Host ""
}

# 编译后端
if (-not $SkipBuild) {
    Write-Host "[3/7] 编译后端..." -ForegroundColor Yellow
    Push-Location backend
    try {
        Write-Host "  交叉编译为 Linux 可执行文件..." -ForegroundColor Gray
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        $env:CGO_ENABLED = "0"
        go build -o main.exe main.go
        if ($LASTEXITCODE -ne 0) {
            throw "后端编译失败"
        }
        
        # 重命名为 main（Linux 文件名）
        if (Test-Path "main") {
            Remove-Item "main"
        }
        Rename-Item "main.exe" "main"
        
        Write-Host "✓ 后端编译完成" -ForegroundColor Green
    }
    catch {
        Write-Host "✗ 后端编译失败: $_" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    finally {
        Pop-Location
    }
    Write-Host ""
} else {
    Write-Host "[3/7] 跳过后端编译" -ForegroundColor Gray
    Write-Host ""
}

# 检查文件是否存在
Write-Host "[4/7] 检查构建产物..." -ForegroundColor Yellow
if (-not (Test-Path "frontend/dist/index.html")) {
    Write-Host "✗ 前端构建产物不存在: frontend/dist/index.html" -ForegroundColor Red
    exit 1
}
if (-not (Test-Path "backend/main")) {
    Write-Host "✗ 后端可执行文件不存在: backend/main" -ForegroundColor Red
    exit 1
}
if (-not (Test-Path "backend/config/config.yml")) {
    Write-Host "⚠ 警告: 配置文件不存在: backend/config/config.yml" -ForegroundColor Yellow
    Write-Host "  请确保服务器上已有配置文件，或手动上传" -ForegroundColor Yellow
}
Write-Host "✓ 构建产物检查通过" -ForegroundColor Green
Write-Host ""

# 上传前端
Write-Host "[5/7] 上传前端文件到服务器..." -ForegroundColor Yellow
Write-Host "  服务器: $ServerUser@$ServerIP:/var/www/blog/" -ForegroundColor Gray
$confirm = Read-Host "  确认上传? (Y/N)"
if ($confirm -ne "Y" -and $confirm -ne "y") {
    Write-Host "  已取消上传" -ForegroundColor Yellow
    exit 0
}

try {
    # 使用 scp 上传（需要 SSH 配置）
    $distPath = Resolve-Path "frontend/dist"
    Write-Host "  正在上传前端文件..." -ForegroundColor Gray
    scp -r "$distPath/*" "${ServerUser}@${ServerIP}:/var/www/blog/"
    Write-Host "✓ 前端文件上传完成" -ForegroundColor Green
}
catch {
    Write-Host "✗ 前端文件上传失败: $_" -ForegroundColor Red
    Write-Host "  提示: 请确保 SSH 密钥已配置，或手动使用 SCP 上传" -ForegroundColor Yellow
}
Write-Host ""

# 上传后端
Write-Host "[6/7] 上传后端文件到服务器..." -ForegroundColor Yellow
Write-Host "  服务器: $ServerUser@$ServerIP:/opt/blog/backend/" -ForegroundColor Gray
$confirm = Read-Host "  确认上传? (Y/N)"
if ($confirm -ne "Y" -and $confirm -ne "y") {
    Write-Host "  已取消上传" -ForegroundColor Yellow
    exit 0
}

try {
    Write-Host "  正在上传后端可执行文件..." -ForegroundColor Gray
    scp "backend/main" "${ServerUser}@${ServerIP}:/opt/blog/backend/main"
    
    # 上传配置文件（如果存在）
    if (Test-Path "backend/config/config.yml") {
        Write-Host "  正在上传配置文件..." -ForegroundColor Gray
        ssh "${ServerUser}@${ServerIP}" "mkdir -p /opt/blog/backend/config"
        scp "backend/config/config.yml" "${ServerUser}@${ServerIP}:/opt/blog/backend/config/config.yml"
    }
    
    Write-Host "✓ 后端文件上传完成" -ForegroundColor Green
}
catch {
    Write-Host "✗ 后端文件上传失败: $_" -ForegroundColor Red
    Write-Host "  提示: 请确保 SSH 密钥已配置，或手动使用 SCP 上传" -ForegroundColor Yellow
}
Write-Host ""

# 重启服务（可选）
Write-Host "[7/7] 重启服务..." -ForegroundColor Yellow
$confirm = Read-Host "  是否在服务器上重启后端服务? (Y/N)"
if ($confirm -eq "Y" -or $confirm -eq "y") {
    try {
        Write-Host "  正在重启后端服务..." -ForegroundColor Gray
        ssh "${ServerUser}@${ServerIP}" "chmod +x /opt/blog/backend/main && systemctl restart blog-backend && systemctl reload nginx"
        Write-Host "✓ 服务重启完成" -ForegroundColor Green
    }
    catch {
        Write-Host "✗ 服务重启失败: $_" -ForegroundColor Red
        Write-Host "  提示: 请手动在服务器上执行: systemctl restart blog-backend" -ForegroundColor Yellow
    }
} else {
    Write-Host "  跳过服务重启" -ForegroundColor Gray
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "部署完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "访问: https://keranthos.me" -ForegroundColor Cyan
Write-Host ""

