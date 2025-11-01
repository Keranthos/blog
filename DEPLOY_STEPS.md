# 部署步骤指南（针对 keranthos.me - 阿里云服务器）

## 服务器信息
- **域名**: keranthos.me
- **服务器IP**: 47.242.6.37
- **云服务商**: 阿里云
- **CDN**: Cloudflare（已配置）

## 第一步：连接服务器并初始化环境

### 1.1 连接到服务器
```bash
ssh root@47.242.6.37
# 或使用密钥（如果配置了）
ssh -i ~/.ssh/your_key root@47.242.6.37
```

### 1.2 检查系统类型
```bash
# 查看系统信息
cat /etc/os-release
# 或
uname -a
```

**常见阿里云系统类型：**
- **Alibaba Cloud Linux** 或 **Aliyun Linux**（基于CentOS，使用yum）
- **Ubuntu**（使用apt）
- **CentOS**（使用yum）

### 1.3 更新系统

**如果是 Aliyun Linux / CentOS（使用yum）：**
```bash
yum update -y
# 阿里云通常已经配置好国内镜像源，速度较快
```

**如果是 Ubuntu（使用apt）：**
```bash
apt update && apt upgrade -y
# 如果速度慢，可以切换到阿里云镜像源
```

### 1.4 安装必要软件

#### 安装基础工具

**Aliyun Linux / CentOS:**
```bash
yum install -y curl wget git vim
```

**Ubuntu:**
```bash
apt install -y curl wget git vim
```

#### 安装 Nginx

**Aliyun Linux / CentOS:**
```bash
yum install -y nginx
systemctl start nginx
systemctl enable nginx
```

**Ubuntu:**
```bash
apt install -y nginx
systemctl start nginx
systemctl enable nginx
```

#### 安装 MySQL（推荐使用 MariaDB，兼容MySQL）

**Aliyun Linux / CentOS:**
```bash
# 安装 MariaDB（MySQL的替代，完全兼容）
yum install -y mariadb-server mariadb
systemctl start mariadb
systemctl enable mariadb
mysql_secure_installation  # 按提示设置root密码
```

**Ubuntu:**
```bash
apt install -y mysql-server
systemctl start mysql
systemctl enable mysql
mysql_secure_installation  # 按提示设置root密码
```

**注意：** MariaDB的命令是 `mysql` 而不是 `mysqld`，其他使用方式相同。

#### 安装 Go

```bash
# 下载 Go 1.21+ (根据你的go.mod版本，从阿里云镜像下载更快)
# 方式1：使用阿里云镜像（推荐，速度快）
wget https://mirrors.aliyun.com/golang/go1.21.5.linux-amd64.tar.gz

# 方式2：官方源（如果镜像不可用）
# wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# 添加到 PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证
go version

# 配置 Go 代理（使用阿里云镜像，加速依赖下载）
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
go env -w GOSUMDB=off  # 可选：关闭校验（如果网络不稳定）
```

#### 安装 Node.js 和 npm

**Aliyun Linux / CentOS:**
```bash
# 使用 NodeSource 仓库（需要先安装 EPEL）
yum install -y epel-release

# 安装 Node.js 18 LTS
curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
yum install -y nodejs

# 验证
node -v
npm -v

# 配置 npm 镜像（使用阿里云镜像，加速下载）
npm config set registry https://registry.npmmirror.com
```

**Ubuntu:**
```bash
# 使用 NodeSource 仓库安装 Node.js 18 LTS
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs

# 验证
node -v
npm -v

# 配置 npm 镜像（使用阿里云镜像，加速下载）
npm config set registry https://registry.npmmirror.com
```

## 第二步：配置数据库

### 2.1 创建数据库和用户

**如果安装的是 MariaDB（Aliyun Linux/CentOS）：**
```bash
mysql -u root -p
```

**如果安装的是 MySQL（Ubuntu）：**
```bash
sudo mysql -u root -p
# 或者
mysql -u root -p
```

在 MySQL/MariaDB 中执行：
```sql
CREATE DATABASE my_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'blog_user'@'localhost' IDENTIFIED BY '你的数据库密码';
GRANT ALL PRIVILEGES ON my_blog.* TO 'blog_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 2.2 测试数据库连接
```bash
mysql -u blog_user -p my_blog
# 如果成功连接说明配置正确
EXIT;
```

**注意：** MariaDB 和 MySQL 的使用方式完全相同，代码中无需修改。

## 第三步：上传代码到服务器

### 3.1 创建项目目录
```bash
mkdir -p /opt/blog
cd /opt/blog
```

### 3.2 阿里云特殊提示

**推荐方式：使用阿里云的文件上传功能**
1. 在阿里云控制台使用 **文件管理** 或 **Workbench** 功能直接上传
2. 或使用 `scp` 从本地直接上传（如果服务器安全组已开放22端口）

### 3.3 上传代码（选择一种方式）

**方式A：使用 Git（推荐，如果代码已推送到远程仓库）**
```bash
git clone https://github.com/your-username/your-blog-repo.git .
# 或
git clone https://github.com/your-username/your-blog-repo.git /opt/blog
```

**方式B：使用 SCP 从本地上传（Windows PowerShell）**
```powershell
# 在本地 Windows PowerShell 执行
cd D:\blog

# Windows 使用 Compress-Archive（如果 tar 不可用）
# 或者先安装 Git Bash，使用其中的 tar 命令

# 如果有 tar 命令（Git Bash 或 WSL）
tar -czf blog.tar.gz --exclude='node_modules' --exclude='frontend/node_modules' --exclude='frontend/dist' --exclude='backend/uploads' --exclude='backend/main.exe' --exclude='backend/test-*' --exclude='.git' backend/ frontend/ .gitignore DEPLOYMENT.txt DEPLOY_CHECKLIST.md DEPLOY_STEPS.md

# 上传到服务器
scp blog.tar.gz root@47.242.6.37:/opt/blog/

# 然后在服务器上解压
ssh root@47.242.6.37
cd /opt/blog
tar -xzf blog.tar.gz
```

**如果 Windows 没有 tar 命令：**
```powershell
# 方法1：使用 Git Bash（如果有 Git）
# 打开 Git Bash，然后执行上面的 tar 命令

# 方法2：使用 PowerShell 压缩（但无法精确排除，需要手动排除）
Compress-Archive -Path backend,frontend -DestinationPath blog.zip -Force
scp blog.zip root@47.242.6.37:/opt/blog/

# 然后在服务器上（需要安装 unzip）
yum install -y unzip  # Aliyun Linux/CentOS
# 或
apt install -y unzip  # Ubuntu
cd /opt/blog
unzip blog.zip
```

**方式C：使用 rsync（如果有rsync）**
```bash
# 在本地执行
rsync -avz --exclude 'node_modules' --exclude 'dist' --exclude 'uploads' \
  ./ root@47.242.6.37:/opt/blog/
```

## 第四步：配置后端

### 4.1 配置后端配置文件
```bash
cd /opt/blog/backend
vim config/config.yml
```

更新为生产环境配置：
```yaml
app:
  name: "WWJ_Blog"
  port: "3000" 

database:
  host: "localhost"
  port: "3306"
  username: "blog_user"  # 改为你创建的用户名
  password: "你的数据库密码"  # 改为你设置的密码
  name: "my_blog"

jwt:
  secretkey: "你的JWT密钥_建议使用32位以上随机字符串"  # 改为更强的密钥
  expiration: 86400
```

**⚠️ 重要：生成强密钥**
```bash
# 生成随机密钥
openssl rand -hex 32
# 将输出复制到 config.yml 的 secretkey 字段
```

### 4.2 创建必要的目录
```bash
cd /opt/blog/backend
mkdir -p uploads/images/content
mkdir -p uploads/images/media
chmod -R 755 uploads
```

### 4.3 初始化数据库表
```bash
cd /opt/blog/backend

# 先确保 go.mod 依赖已安装
go mod download

# 编译并运行一次以触发 AutoMigrate（或创建迁移脚本）
go run main.go
# 运行几秒后按 Ctrl+C 停止，表结构应该已经创建

# 或者直接编译
go build -o /opt/blog/backend/main .
```

### 4.4 测试后端运行
```bash
cd /opt/blog/backend
GIN_MODE=release ./main
# 在另一个终端测试
curl http://localhost:3000/api/articles?type=blog&page=1&limit=10
# 如果返回数据说明后端正常，按 Ctrl+C 停止
```

## 第五步：构建前端

### 5.1 安装前端依赖
```bash
cd /opt/blog/frontend
npm ci  # 使用 package-lock.json，确保依赖版本一致
```

### 5.2 配置前端环境变量（如果需要）
```bash
# 如果前后端同源部署（推荐），则不需要修改
# 前端会使用相对路径 /api

# 如果异源部署，创建 .env.production
echo 'VUE_APP_API_BASE_URL=https://keranthos.me' > .env.production
echo 'VUE_APP_API_PREFIX=/api' >> .env.production
```

### 5.3 构建前端
```bash
cd /opt/blog/frontend
npm run build
# 构建完成后，dist/ 目录包含静态文件
```

## 第六步：配置 Nginx

### 6.1 创建 Nginx 配置
```bash
vim /etc/nginx/sites-available/keranthos.me
```

配置内容：
```nginx
server {
    listen 80;
    server_name keranthos.me www.keranthos.me;

    # 前端静态资源
    root /opt/blog/frontend/dist;
    index index.html;

    # 前端路由历史模式支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:3000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持（如果需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 后端静态上传文件
    location /uploads/ {
        proxy_pass http://127.0.0.1:3000/;
        proxy_set_header Host $host;
        
        # 缓存设置（可选）
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    # 静态资源缓存（CSS、JS）
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;
}
```

### 6.2 启用站点
```bash
# 创建软链接
ln -s /etc/nginx/sites-available/keranthos.me /etc/nginx/sites-enabled/

# 删除默认站点（可选）
rm /etc/nginx/sites-enabled/default

# 测试配置
nginx -t

# 重新加载 Nginx
systemctl reload nginx
```

### 6.3 测试 HTTP 访问
在浏览器访问 `http://keranthos.me`，应该能看到前端页面（如果后端已启动）。

## 第七步：配置 HTTPS（SSL证书）

### 7.1 安装 Certbot
```bash
apt install -y certbot python3-certbot-nginx
```

### 7.2 配置 Cloudflare DNS
在 Cloudflare 控制台确认：
- A记录：`keranthos.me` → `47.242.6.37`，代理状态：仅DNS（橙色云朵关闭）
- A记录：`www.keranthos.me` → `47.242.6.37`，代理状态：仅DNS（橙色云朵关闭）

**⚠️ 重要：申请证书时需要关闭 Cloudflare 代理（仅DNS模式），证书申请成功后再开启。**

### 7.3 申请 SSL 证书
```bash
certbot --nginx -d keranthos.me -d www.keranthos.me
```

按提示操作：
- 输入邮箱地址
- 同意服务条款
- 选择是否共享邮箱（建议选择否）
- 选择是否重定向 HTTP 到 HTTPS（建议选择是）

### 7.4 测试证书自动续期
```bash
certbot renew --dry-run
```

### 7.5 开启 Cloudflare 代理
证书申请成功后，在 Cloudflare 控制台将 DNS 记录的代理状态改为：已代理（橙色云朵开启）

## 第八步：配置后端服务（Systemd）

### 8.1 创建 Systemd 服务文件
```bash
vim /etc/systemd/system/blog-backend.service
```

内容：
```ini
[Unit]
Description=Blog Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/blog/backend
ExecStart=/opt/blog/backend/main
Restart=always
RestartSec=5s
Environment=GIN_MODE=release

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=blog-backend

[Install]
WantedBy=multi-user.target
```

### 8.2 启用并启动服务
```bash
# 重新加载 systemd
systemctl daemon-reload

# 启用服务（开机自启）
systemctl enable blog-backend

# 启动服务
systemctl start blog-backend

# 查看状态
systemctl status blog-backend

# 查看日志
journalctl -u blog-backend -f
```

### 8.3 测试服务
```bash
# 检查服务是否运行
systemctl status blog-backend

# 测试 API
curl https://keranthos.me/api/articles?type=blog&page=1&limit=10
```

## 第九步：配置防火墙和安全组

### 9.1 阿里云安全组配置（重要！）

**⚠️ 必须先在阿里云控制台配置安全组规则，否则无法访问！**

1. 登录 [阿里云控制台](https://ecs.console.aliyun.com)
2. 进入 **云服务器ECS** → **实例** → 找到你的服务器
3. 点击 **安全组** → **配置规则**
4. 添加以下入方向规则：
   - **22/tcp** (SSH) - 来源：你的IP或 0.0.0.0/0（不推荐，建议限制IP）
   - **80/tcp** (HTTP) - 来源：0.0.0.0/0
   - **443/tcp** (HTTPS) - 来源：0.0.0.0/0
   - **3000/tcp** (后端服务，仅限内网) - 来源：127.0.0.1/32

### 9.2 系统防火墙配置（可选，如果启用了系统防火墙）

**Aliyun Linux / CentOS（firewalld）：**
```bash
# 检查防火墙状态
systemctl status firewalld

# 如果启用，配置规则
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --permanent --add-service=ssh
firewall-cmd --reload

# 如果不需要系统防火墙（阿里云安全组已提供保护），可以关闭
systemctl stop firewalld
systemctl disable firewalld
```

**Ubuntu（UFW）：**
```bash
# 检查防火墙状态
ufw status

# 配置规则
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw enable

# 如果不需要系统防火墙（阿里云安全组已提供保护），可以关闭
ufw disable
```

**注意：** 阿里云的安全组已经提供了网络层防火墙，系统防火墙通常是关闭的。建议保持系统防火墙关闭，避免配置冲突。

## 第十步：更新后端 CORS 配置

### 10.1 更新允许的域名
```bash
vim /opt/blog/backend/routes/routes.go
```

修改 CORS 配置：
```go
AllowOrigins: []string{
    "http://localhost:3000",
    "http://localhost:3002",
    "https://keranthos.me",        // 添加
    "https://www.keranthos.me",    // 添加
},
```

### 10.2 重新编译并重启
```bash
cd /opt/blog/backend
go build -o main .
systemctl restart blog-backend
```

## 第十一步：验证部署

### 11.1 前端验证
- [ ] 访问 https://keranthos.me 能正常显示
- [ ] 路由跳转正常（博客、随笔、书影集等）
- [ ] 图片能正常加载

### 11.2 后端验证
- [ ] API 请求正常：`curl https://keranthos.me/api/articles?type=blog&page=1&limit=10`
- [ ] 登录功能正常
- [ ] 图片上传功能正常

### 11.3 数据库验证
```bash
mysql -u blog_user -p my_blog
SHOW TABLES;
# 应该能看到 blog_articles, research_articles 等表
EXIT;
```

## 第十二步：后续维护

### 12.1 更新代码
```bash
# 如果使用 Git
cd /opt/blog
git pull

# 重新构建前端
cd frontend
npm ci
npm run build

# 重新编译后端
cd ../backend
go build -o main .

# 重启服务
systemctl restart blog-backend
systemctl reload nginx
```

### 12.2 备份数据库
```bash
# 创建备份脚本
vim /opt/blog/backup-db.sh
```

内容：
```bash
#!/bin/bash
BACKUP_DIR="/opt/blog/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
mysqldump -u blog_user -p你的密码 my_blog > $BACKUP_DIR/backup_$DATE.sql
# 保留最近7天的备份
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete
```

设置执行权限：
```bash
chmod +x /opt/blog/backup-db.sh
```

设置定时任务（每天凌晨2点备份）：
```bash
crontab -e
# 添加
0 2 * * * /opt/blog/backup-db.sh
```

**阿里云建议：** 除了本地备份，还可以使用阿里云的快照功能定期备份整个服务器，或使用OSS对象存储备份数据库文件。

### 12.3 监控日志
```bash
# 后端日志
journalctl -u blog-backend -f

# Nginx 日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

## 常见问题排查

### 问题1：无法SSH连接到服务器
- **检查阿里云安全组**：确保22端口已开放，来源IP已添加
- 检查服务器是否运行：在阿里云控制台查看实例状态
- 尝试使用阿里云控制台的 **Workbench** 或 **VNC** 连接

### 问题2：502 Bad Gateway
- 检查后端服务是否运行：`systemctl status blog-backend`
- 检查端口是否被占用：`netstat -tlnp | grep 3000` 或 `ss -tlnp | grep 3000`
- 查看后端日志：`journalctl -u blog-backend -n 50`
- **检查阿里云安全组**：确保3000端口已对127.0.0.1开放（内网访问）

### 问题3：前端页面空白或无法访问
- 检查 Nginx 配置：`nginx -t`
- 检查 dist 目录权限：`ls -la /opt/blog/frontend/dist`
- 查看 Nginx 错误日志：`tail -f /var/log/nginx/error.log`
- **检查阿里云安全组**：确保80和443端口已开放

### 问题4：数据库连接失败
- 检查 MySQL/MariaDB 服务：
  - Aliyun Linux/CentOS: `systemctl status mariadb`
  - Ubuntu: `systemctl status mysql`
- 测试连接：`mysql -u blog_user -p my_blog`
- 检查配置文件：`cat /opt/blog/backend/config/config.yml`
- 检查数据库用户权限：`mysql -u root -p -e "SHOW GRANTS FOR 'blog_user'@'localhost';"`

### 问题5：图片上传失败
- 检查 uploads 目录权限：`ls -la /opt/blog/backend/uploads`
- 确保目录可写：`chmod -R 755 /opt/blog/backend/uploads`
- 检查磁盘空间：`df -h`

### 问题6：npm/yum 安装速度慢
- **npm 镜像**：已配置阿里云镜像 `https://registry.npmmirror.com`
- **Go 代理**：已配置阿里云镜像 `https://mirrors.aliyun.com/goproxy/`
- **yum 源**：阿里云系统默认已配置阿里云镜像源，通常速度较快

### 问题7：SSL证书申请失败
- **检查 Cloudflare 代理状态**：申请证书时必须关闭（仅DNS模式）
- 检查域名DNS解析：`nslookup keranthos.me`
- 确保80端口已开放（Let's Encrypt需要验证）
- 检查 Nginx 配置是否正确：`nginx -t`

## 完成 ✅

部署完成后，你的博客应该可以通过 https://keranthos.me 访问了！

**建议的后续优化：**
1. 配置 Cloudflare 缓存规则
2. 启用 Cloudflare 的 Auto Minify（HTML/CSS/JS压缩）
3. 配置图片 CDN（如使用对象存储）
4. 设置监控告警（可选）


