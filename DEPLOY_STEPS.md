# 部署步骤指南（方案A：本地构建 + 服务器部署文件）

## 服务器信息
- **域名**: keranthos.me
- **服务器IP**: 47.242.6.37
- **云服务商**: 阿里云
- **CDN**: Cloudflare（已配置）
- **系统**: Alibaba Cloud Linux / CentOS（使用yum）

## 📋 方案说明

**本方案采用本地构建、服务器仅部署文件的方式：**
- ✅ 服务器不需要安装 Node.js、npm、Go 编译器
- ✅ 节省服务器内存（约节省800MB-1.7GB）
- ✅ 安全性更高（生产环境无开发工具）
- ✅ 构建速度更快（本地通常比服务器快）

---

## 第一部分：服务器环境准备

### 步骤1：连接服务器

```bash
ssh root@47.242.6.37
# 或使用密钥（如果配置了）
ssh -i ~/.ssh/your_key root@47.242.6.37
```

### 步骤2：更新系统

```bash
yum update -y
```

### 步骤3：安装基础工具

```bash
yum install -y curl wget git vim unzip
```

### 步骤4：安装 Nginx

```bash
yum install -y nginx
systemctl start nginx
systemctl enable nginx
```

验证 Nginx 是否运行：
```bash
systemctl status nginx
curl http://localhost
# 应该能看到 Nginx 默认页面
```

### 步骤5：安装数据库

**说明：** 
- 后端代码使用 `gorm.io/driver/mysql`（MySQL驱动）
- MariaDB 完全兼容 MySQL 协议和语法，可以直接使用
- 在 CentOS/Alibaba Cloud Linux 中，默认 yum 源提供的是 **MariaDB**（而不是 MySQL）
- 两者在代码层面**完全相同**，无需修改任何代码

**选择方案A：安装 MariaDB（推荐，默认源，安装简单）**

```bash
# 安装 MariaDB（MySQL的替代，完全兼容）
yum install -y mariadb-server mariadb
systemctl start mariadb
systemctl enable mariadb

# 初始化数据库（设置root密码）
mysql_secure_installation
```

**选择方案B：安装 MySQL（如果需要官方 MySQL）**

```bash
# 添加 MySQL 官方仓库
wget https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
rpm -ivh mysql80-community-release-el7-3.noarch.rpm

# 安装 MySQL 8.0
yum install -y mysql-community-server
systemctl start mysqld
systemctl enable mysqld

# 获取临时密码
grep 'temporary password' /var/log/mysqld.log

# 使用临时密码登录并修改
mysql -u root -p
ALTER USER 'root'@'localhost' IDENTIFIED BY '你的新密码';
# MySQL 8.0 需要设置密码策略
SET GLOBAL validate_password.policy=LOW;
SET GLOBAL validate_password.length=8;
```

**按提示操作（mysql_secure_installation）：**
- 设置 root 密码（请记住这个密码）
- 移除匿名用户：Y
- 禁止 root 远程登录：Y（如果只需要本地登录）
- 移除测试数据库：Y
- 重新加载权限表：Y

**注意：** MariaDB 和 MySQL 使用的命令相同（都是 `mysql`），后续步骤完全一致。

### 步骤6：创建数据库和用户

```bash
mysql -u root -p
```

在 MySQL 中执行：
```sql
CREATE DATABASE my_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'blog_user'@'localhost' IDENTIFIED BY '你的数据库密码';
GRANT ALL PRIVILEGES ON my_blog.* TO 'blog_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

测试连接：
```bash
mysql -u blog_user -p my_blog
# 如果成功连接说明配置正确
EXIT;
```

---

## 第二部分：本地构建（Windows）

### 步骤1：构建前端

在本地 Windows PowerShell 中执行：

```powershell
cd D:\blog\frontend

# 安装依赖（如果还没安装）
npm install

# 构建生产版本（使用相对路径，通过 Nginx 代理访问 API）
npm run build
```

**注意：** 
- 前端已配置为生产环境自动使用相对路径 `/api`（而不是 `http://localhost:3000/api`）
- 这样可以通过 Nginx 代理统一访问后端 API
- 无需设置环境变量，代码会自动检测生产环境

构建完成后，`frontend/dist/` 目录包含所有静态文件。

### 步骤2：编译后端

```powershell
cd D:\blog\backend

# 编译为Linux可执行文件（交叉编译）
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o main.exe main.go

# 或使用 Git Bash（推荐）
# GOOS=linux GOARCH=amd64 go build -o main main.go
```

编译完成后，`backend/main.exe`（或 `backend/main`）就是 Linux 可执行文件。

**验证编译结果：**
```powershell
file backend/main.exe  # 如果是Git Bash
# 应该显示：ELF 64-bit LSB executable, x86-64
```

---

## 第三部分：准备配置文件

### 步骤1：创建服务器端配置文件

在本地创建 `backend/config/config.yml`（参考 `config.yml.example`）：

```yaml
app:
  name: "WWJ_Blog"
  port: "3000"

database:
  host: "localhost"
  port: "3306"
  username: "blog_user"
  password: "你的数据库密码"  # 改为步骤1.6中设置的密码
  name: "my_blog"

jwt:
  secretkey: "你的JWT密钥"  # 建议使用 openssl rand -hex 32 生成
  expiration: 86400
```

**生成 JWT 密钥（在 Git Bash 中）：**
```bash
openssl rand -hex 32
# 将输出复制到 config.yml 的 secretkey 字段
```

### 步骤2：创建服务器目录结构

在服务器上执行：

```bash
# 创建项目目录
mkdir -p /opt/blog/backend
mkdir -p /opt/blog/backend/uploads/images/content
mkdir -p /opt/blog/backend/uploads/images/media
mkdir -p /var/www/blog

# 设置权限
chmod -R 755 /opt/blog/backend/uploads
```

---

## 第四部分：上传文件到服务器

### 方式A：使用 SCP（推荐）

在本地 Windows PowerShell 或 Git Bash 中执行：

**首次部署或完整更新：**

```powershell
# 1. 清理服务器上的旧文件（可选，推荐在首次部署或重大更新时）
ssh root@47.242.6.37 "rm -rf /var/www/blog/*"

# 2. 上传前端构建文件
scp -r frontend/dist/* root@47.242.6.37:/var/www/blog/

# 3. 上传后端二进制文件
scp backend/main root@47.242.6.37:/opt/blog/backend/
# 或如果文件名为 main.exe
scp backend/main.exe root@47.242.6.37:/opt/blog/backend/main

# 4. 上传配置文件（如果需要更新）
scp backend/config/config.yml root@47.242.6.37:/opt/blog/backend/config/
```

**增量更新（只更新前端，不清除旧文件）：**

```powershell
# 直接上传新文件（会覆盖同名文件，但不会删除旧文件）
scp -r frontend/dist/* root@47.242.6.37:/var/www/blog/
```

**注意：**
- 首次部署或重大更新时，建议先清理旧文件避免缓存问题
- 如果只是代码小更新，可以直接覆盖，不需要清理
- 清理命令会删除 `/var/www/blog/` 下所有文件，确保重要文件已备份

**如果上传失败，检查：**
1. 阿里云安全组是否开放 22 端口（SSH）
2. SSH 是否可用：`ssh root@47.242.6.37`

### 方式B：使用阿里云控制台（适合大文件）

1. 在本地压缩文件：
   ```powershell
   # 压缩前端
   Compress-Archive -Path frontend/dist/* -DestinationPath dist.zip
   
   # 压缩后端（包含二进制和配置）
   # 先创建临时目录
   mkdir temp_backend
   copy backend\main.exe temp_backend\main
   copy backend\config\config.yml temp_backend\config.yml
   Compress-Archive -Path temp_backend/* -DestinationPath backend.zip
   ```

2. 在阿里云控制台使用 **Workbench** 或 **文件管理** 功能上传

3. 在服务器上解压：
   ```bash
   cd /var/www/blog
   unzip dist.zip
   
   cd /opt/blog/backend
   unzip backend.zip
   ```

### 步骤：设置后端文件权限

在服务器上执行：

```bash
cd /opt/blog/backend
chmod +x main  # 给可执行文件添加执行权限
```

---

## 第五部分：初始化数据库

**⚠️ 重要：表结构是在后端程序首次运行时自动创建的，而不是手动创建。必须先运行后端程序！**

**🔒 数据安全说明：**
- `AutoMigrate` **不会删除**已有数据和表
- 如果表已存在，只会**添加缺失的列**，不会删除列或数据
- 如果表不存在，才会创建新表
- 每次启动时，只是检查和更新表结构（增量更新）
- **你的文章、评论、媒体等所有数据都是安全的**

### 步骤1：检查必要文件是否存在

在服务器上执行：

```bash
cd /opt/blog/backend

# 检查文件是否存在
ls -la main                    # 应该看到可执行文件
ls -la config/config.yml       # 应该看到配置文件

# 如果文件不存在，说明还没有上传，请先完成"第四部分：上传文件到服务器"
```

### 步骤2：检查配置文件

```bash
# 查看配置文件内容（确认数据库配置正确）
cat config/config.yml
```

确认以下配置正确：
- `database.host`: 应该是 `localhost`
- `database.port`: 应该是 `3306`
- `database.username`: 应该是 `blog_user`
- `database.password`: 应该是你创建的密码
- `database.name`: 应该是 `my_blog`

### 步骤3：首次运行后端以创建表结构

在服务器上执行：

```bash
cd /opt/blog/backend

# 确保文件有执行权限
chmod +x main

# 运行后端程序（会输出日志）
GIN_MODE=release ./main
```

**预期输出（成功的情况）：**
```
Created directory: uploads
Created directory: uploads/images
开始创建数据库索引...
数据库索引创建完成
```

**注意：** 
- 程序会一直运行（这是正常的！），不会自动退出
- 在生产模式下（`GIN_MODE=release`），路由信息不会显示在日志中
- 看到"数据库索引创建完成"说明**程序已成功启动，表已经创建**
- **首次运行会创建表，后续运行只会检查并更新表结构，不会删除任何数据**

**验证程序正在运行：**
- 在**另一个终端窗口**（或新开一个SSH连接）执行：
  ```bash
  # ⚠️ 注意：URL 必须用引号括起来，否则 & 会被 shell 解析
  curl "http://localhost:3000/api/articles?type=blog&page=1&limit=10"
  ```
  - 如果返回 JSON 数据，说明程序正常运行 ✅
  - 如果返回错误信息（如 `{"error":"..."}`），说明程序在运行，只是参数问题
  - 如果返回 `curl: (7) Failed to connect`，说明程序可能没有启动

**确认程序正常运行后，在原终端按 `Ctrl+C` 停止程序。**

**如果出现错误：**

**错误1：配置文件不存在或读取失败**
```
Error reading config file: ...
```
- 检查 `config/config.yml` 是否存在
- 检查配置文件格式是否正确（YAML格式）

**错误2：数据库连接失败**
```
error connecting to database: ...
```
- 检查数据库服务是否运行：`systemctl status mariadb`
- 检查数据库用户名、密码是否正确
- 测试数据库连接：`mysql -u blog_user -p my_blog`

**错误3：权限不足**
```
permission denied
```
- 检查文件权限：`chmod +x main`
- 检查目录权限：`chmod -R 755 /opt/blog/backend`

**运行几秒后，按 `Ctrl+C` 停止后端程序。**

### 步骤4：验证表是否创建

```bash
mysql -u blog_user -p my_blog
```

在 MySQL 中执行：
```sql
SHOW TABLES;
```

**预期结果（应该看到以下表）：**
```
+--------------------------+
| Tables_in_my_blog         |
+--------------------------+
| blog_articles            |
| comments                 |
| media                    |
| moments                  |
| presentations            |
| project_articles         |
| questions                |
| research_articles        |
| users                    |
| weather                  |
+--------------------------+
```

**如果 `SHOW TABLES;` 后没有任何显示（Empty set）：**

说明表还没有创建，可能的原因和解决方法：

1. **后端程序还没有运行**
   - 返回步骤3，运行 `./main` 并确认没有错误

2. **后端程序运行失败**
   - 检查后端程序输出的错误信息
   - 重新运行并查看完整日志：`GIN_MODE=release ./main 2>&1 | tee startup.log`

3. **数据库连接失败**
   ```bash
   # 测试数据库连接
   mysql -u blog_user -p my_blog
   # 如果连接失败，检查：
   # - 数据库服务：systemctl status mariadb
   # - 配置文件中的密码是否正确
   ```

4. **配置文件路径错误**
   ```bash
   # 确认配置文件在正确的位置
   ls -la /opt/blog/backend/config/config.yml
   # 确认后端程序从正确的工作目录运行
   cd /opt/blog/backend && pwd
   ```

5. **数据库用户权限不足**
   ```bash
   # 检查用户权限
   mysql -u root -p
   SHOW GRANTS FOR 'blog_user'@'localhost';
   # 如果没有权限，重新授权：
   GRANT ALL PRIVILEGES ON my_blog.* TO 'blog_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

**验证成功后：**
```sql
EXIT;
```

### 步骤5：再次验证（可选）

如果想确认表结构是否正确：

```bash
mysql -u blog_user -p my_blog
```

```sql
-- 查看某个表的结构
DESCRIBE blog_articles;

-- 或查看所有表
SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'my_blog';

EXIT;
```

---

## 第六部分：配置 Nginx

### 步骤1：创建 Nginx 配置文件

```bash
vim /etc/nginx/conf.d/keranthos.me.conf
```

配置内容：
```nginx
# 如果使用 Cloudflare 代理，需要先配置真实 IP 获取
# Cloudflare IP 段（用于信任 CF 的真实 IP 头部）
# 在 http 块中添加（如果不在 http 块中，也可以放在 server 块之前）

# 在 /etc/nginx/nginx.conf 的 http 块中添加以下配置（如果使用 CF）
# set_real_ip_from 103.21.244.0/22;
# set_real_ip_from 103.22.200.0/22;
# set_real_ip_from 103.31.4.0/22;
# set_real_ip_from 104.16.0.0/13;
# set_real_ip_from 104.24.0.0/14;
# set_real_ip_from 108.162.192.0/18;
# set_real_ip_from 131.0.72.0/22;
# set_real_ip_from 141.101.64.0/18;
# set_real_ip_from 162.158.0.0/15;
# set_real_ip_from 172.64.0.0/13;
# set_real_ip_from 173.245.48.0/20;
# set_real_ip_from 188.114.96.0/20;
# set_real_ip_from 190.93.240.0/20;
# set_real_ip_from 197.234.240.0/22;
# set_real_ip_from 198.41.128.0/17;
# set_real_ip_from 2400:cb00::/32;
# set_real_ip_from 2606:4700::/32;
# set_real_ip_from 2803:f800::/32;
# set_real_ip_from 2405:b500::/32;
# set_real_ip_from 2405:8100::/32;
# set_real_ip_from 2c0f:f248::/32;
# set_real_ip_from 2a06:98c0::/29;
# real_ip_header CF-Connecting-IP;

server {
    listen 80;
    server_name keranthos.me www.keranthos.me;

    # 前端静态资源
    root /var/www/blog;
    index index.html;

    # 前端路由历史模式支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:3000/;
        proxy_set_header Host $host;
        
        # 真实 IP 获取（如果使用 Cloudflare 代理）
        # 优先使用 CF-Connecting-IP，如果不存在则使用 X-Forwarded-For
        proxy_set_header X-Real-IP $http_cf_connecting_ip;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Cloudflare 相关头部（如果使用 CF 代理）
        proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
        proxy_set_header CF-Ray $http_cf_ray;
        proxy_set_header CF-Visitor $http_cf_visitor;
        
        # WebSocket 支持（如果需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 后端静态上传文件（使用 ^~ 确保优先级高于正则 location）
    location ^~ /uploads/ {
        alias /opt/blog/backend/uploads/;
        
        # 缓存设置（可选）
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    # 静态资源缓存（CSS、JS）- /uploads/ 已用 ^~ 排除，不会被此规则匹配
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

**⚠️ 如果使用 Cloudflare 代理，还需要修改 Nginx 主配置：**

编辑 `/etc/nginx/nginx.conf`，在 `http` 块中添加：

```nginx
http {
    # ... 其他配置 ...
    
    # Cloudflare 真实 IP 配置（如果使用 CF 代理）
    # 信任 Cloudflare IP 段，从 CF-Connecting-IP 获取真实 IP
    include /etc/nginx/cloudflare.conf;  # 或者直接在下面添加 IP 段
}
```

创建 Cloudflare IP 配置文件（可选，但推荐）：
```bash
vim /etc/nginx/cloudflare.conf
```

内容：
```nginx
# Cloudflare IP 段（2024年更新）
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 131.0.72.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
# IPv6
set_real_ip_from 2400:cb00::/32;
set_real_ip_from 2606:4700::/32;
set_real_ip_from 2803:f800::/32;
set_real_ip_from 2405:b500::/32;
set_real_ip_from 2405:8100::/32;
set_real_ip_from 2c0f:f248::/32;
set_real_ip_from 2a06:98c0::/29;

# 从 CF-Connecting-IP 头部获取真实 IP
real_ip_header CF-Connecting-IP;
```

**然后在 `/etc/nginx/nginx.conf` 的 `http` 块中添加：**
```nginx
include /etc/nginx/cloudflare.conf;
```

**注意：**
- 如果**不使用** Cloudflare 代理（仅DNS模式），则不需要添加上述配置
- 如果**使用** Cloudflare 代理（橙色云朵开启），则建议添加上述配置以正确获取用户真实 IP

### 步骤2：配置 Cloudflare 真实 IP（如果使用 CF 代理）

**如果使用 Cloudflare 代理（橙色云朵开启）：**

1. 创建 Cloudflare IP 配置文件：
```bash
vim /etc/nginx/cloudflare.conf
```

内容：
```nginx
# Cloudflare IP 段（2024年更新）
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 131.0.72.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
# IPv6
set_real_ip_from 2400:cb00::/32;
set_real_ip_from 2606:4700::/32;
set_real_ip_from 2803:f800::/32;
set_real_ip_from 2405:b500::/32;
set_real_ip_from 2405:8100::/32;
set_real_ip_from 2c0f:f248::/32;
set_real_ip_from 2a06:98c0::/29;

# 从 CF-Connecting-IP 头部获取真实 IP
real_ip_header CF-Connecting-IP;
```

2. 在主配置文件中引入：
```bash
vim /etc/nginx/nginx.conf
```

在 `http` 块中添加：
```nginx
http {
    # ... 其他配置 ...
    
    # Cloudflare 真实 IP 配置（如果使用 CF 代理）
    include /etc/nginx/cloudflare.conf;
    
    # ... 其他配置 ...
}
```

**如果只使用 Cloudflare DNS（不开启代理），可以跳过此步骤。**

### 步骤3：测试并重新加载 Nginx

```bash
# 测试配置
nginx -t

# 重新加载
systemctl reload nginx
```

### 步骤4：测试 HTTP 访问

在浏览器访问 `http://47.242.6.37` 或 `http://keranthos.me`，应该能看到前端页面（如果后端已启动）。

**注意：** 如果使用 Cloudflare 代理，访问域名会通过 CF，速度可能更快。直接访问 IP 则不会经过 CF。

---

## 第七部分：配置后端服务（Systemd）

### 步骤1：创建 Systemd 服务文件

```bash
vim /etc/systemd/system/blog-backend.service
```

内容：
```ini
[Unit]
Description=Blog Backend Service
After=network.target mariadb.service

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

### 步骤2：启用并启动服务

```bash
# 重新加载 systemd
systemctl daemon-reload

# 启用服务（开机自启）
systemctl enable blog-backend

# 启动服务
systemctl start blog-backend

# 查看状态
systemctl status blog-backend

# 查看日志（实时）
journalctl -u blog-backend -f
```

### 步骤3：测试服务

```bash
# 检查服务是否运行
systemctl status blog-backend

# 测试 API
curl http://localhost:3000/api/articles?type=blog&page=1&limit=10
```

---

## 第八部分：配置 HTTPS（SSL证书）

### 步骤1：安装 Certbot

```bash
# CentOS/Alibaba Cloud Linux 需要先安装 EPEL
yum install -y epel-release
yum install -y certbot python3-certbot-nginx
```

### 步骤2：配置 Cloudflare DNS

在 Cloudflare 控制台确认：
- A记录：`keranthos.me` → `47.242.6.37`，**代理状态：仅DNS（橙色云朵关闭）**
- A记录：`www.keranthos.me` → `47.242.6.37`，**代理状态：仅DNS（橙色云朵关闭）**

**⚠️ 重要：申请证书时必须关闭 Cloudflare 代理（仅DNS模式），证书申请成功后再开启。**

### 步骤3：申请 SSL 证书

```bash
certbot --nginx -d keranthos.me -d www.keranthos.me
```

按提示操作：
- 输入邮箱地址
- 同意服务条款（A）
- 选择是否共享邮箱（N）
- 选择是否重定向 HTTP 到 HTTPS（2，推荐选择重定向）

### 步骤4：测试证书自动续期

```bash
certbot renew --dry-run
```

### 步骤5：开启 Cloudflare 代理

证书申请成功后，在 Cloudflare 控制台将 DNS 记录的代理状态改为：**已代理（橙色云朵开启）**

---

## 第九部分：配置防火墙和安全组

### 步骤1：阿里云安全组配置（必须！）

**⚠️ 必须在阿里云控制台配置安全组规则，否则无法访问！**

1. 登录 [阿里云控制台](https://ecs.console.aliyun.com)
2. 进入 **云服务器ECS** → **实例** → 找到你的服务器
3. 点击 **安全组** → **配置规则** → **入方向** → **添加安全组规则**
4. 添加以下规则：
   - **22/tcp** (SSH) - 授权对象：你的IP或 0.0.0.0/0（不推荐，建议限制IP）
   - **80/tcp** (HTTP) - 授权对象：0.0.0.0/0
   - **443/tcp** (HTTPS) - 授权对象：0.0.0.0/0
   - **3000/tcp** (后端服务，仅限内网) - 授权对象：127.0.0.1/32

### 步骤2：系统防火墙（可选）

**检查系统防火墙状态：**
```bash
systemctl status firewalld
```

**如果启用了防火墙，配置规则：**
```bash
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --permanent --add-service=ssh
firewall-cmd --reload
```

**如果不需要系统防火墙（阿里云安全组已提供保护）：**
```bash
systemctl stop firewalld
systemctl disable firewalld
```

**注意：** 阿里云的安全组已经提供了网络层防火墙，系统防火墙通常是关闭的。建议保持系统防火墙关闭，避免配置冲突。

---

## 第十部分：验证部署

### 验证清单

- [ ] **前端验证**
  - 访问 https://keranthos.me 能正常显示
  - 路由跳转正常（博客、随笔、书影集等）
  - 图片能正常加载

- [ ] **后端验证**
  - API 请求正常：`curl https://keranthos.me/api/articles?type=blog&page=1&limit=10`
  - 登录功能正常
  - 图片上传功能正常

- [ ] **数据库验证**
  ```bash
  mysql -u blog_user -p my_blog
  SHOW TABLES;
  # 应该能看到 blog_articles, research_articles 等表
  EXIT;
  ```

- [ ] **服务状态验证**
  ```bash
  systemctl status blog-backend
  systemctl status nginx
  # MariaDB 使用 mariadb，MySQL 使用 mysqld
  systemctl status mariadb  # 或 systemctl status mysqld
  ```

---

## 第十一部分：后续维护

### 更新代码流程

**当需要更新代码时：**

#### 方式一：使用自动化脚本（推荐）⚡

**后端更新（一键完成）：**

在项目根目录下执行：
```powershell
cd backend; .\..\update-backend.ps1
```

或者在 backend 目录下直接执行：
```powershell
.\update-backend.ps1
```

**脚本功能：**
- ✅ 自动编译 Linux 版本
- ✅ 自动停止服务器上的后端服务
- ✅ 自动上传文件（使用临时文件名避免冲突）
- ✅ 自动替换文件并设置权限
- ✅ 自动启动服务并检查状态

**脚本参数（可选）：**
```powershell
# 指定服务器IP和用户
.\update-backend.ps1 -ServerIP "47.242.6.37" -ServerUser "root"
```

---

#### 方式二：手动更新

**当需要更新代码时：**

1. **本地构建（Windows）：**
   ```powershell
   # 更新代码
   cd D:\blog
   git pull
   
   # 构建前端
   cd frontend
   npm install  # 如果 package.json 有更新
   npm run build
   
   # 编译后端
   cd ../backend
   $env:GOOS="linux"
   $env:GOARCH="amd64"
   go build -o main.exe main.go
   ```

2. **上传到服务器：**
   ```powershell
   # 方式A：增量更新（推荐，只覆盖同名文件）
   # 上传前端
   scp -r frontend/dist/* root@47.242.6.37:/var/www/blog/
   
   # 上传后端（先停止服务避免文件锁定）
   ssh root@47.242.6.37 "systemctl stop blog-backend"
   scp backend/main.exe root@47.242.6.37:/opt/blog/backend/main.new
   ssh root@47.242.6.37 "cd /opt/blog/backend && mv -f main.new main && chmod +x main && systemctl start blog-backend"
   
   # 方式B：完全清理后更新（如果有缓存问题）
   # 清理旧文件
   ssh root@47.242.6.37 "rm -rf /var/www/blog/*"
   # 上传新文件
   scp -r frontend/dist/* root@47.242.6.37:/var/www/blog/
   ```

3. **重启服务（服务器上）：**
   ```bash
   systemctl restart blog-backend
   systemctl reload nginx
   ```

### 备份数据库

创建备份脚本：
```bash
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

### 查看日志

```bash
# 后端日志
journalctl -u blog-backend -f

# Nginx 访问日志
tail -f /var/log/nginx/access.log

# Nginx 错误日志
tail -f /var/log/nginx/error.log
```

---

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
- 检查后端可执行文件权限：`ls -la /opt/blog/backend/main`

### 问题3：前端页面空白或无法访问
- 检查 Nginx 配置：`nginx -t`
- 检查 dist 目录权限：`ls -la /var/www/blog`
- 查看 Nginx 错误日志：`tail -f /var/log/nginx/error.log`
- **检查阿里云安全组**：确保80和443端口已开放

### 问题3.1：前端控制台报错 `ERR_CONNECTION_REFUSED` 或 `ERR_CONNECTION_RESET`
**错误现象：**
```
Failed to load resource: net::ERR_CONNECTION_REFUSED
localhost:3000/api/articles?page=1&limit=3&type=blog:1
```

**原因：** 前端在尝试连接 `localhost:3000`，但在生产环境中应该通过 Nginx 代理访问 API。

**解决方法：**
1. **检查前端构建是否正确：**
   - 前端代码已修复，生产环境会自动使用相对路径 `/api`
   - 如果使用的是旧版本构建，需要重新构建：
   ```powershell
   cd D:\blog\frontend
   npm run build
   ```
   - 重新上传 `frontend/dist/*` 到服务器

2. **检查 Nginx 代理配置：**
   ```bash
   # 确认 Nginx 配置中有 API 代理
   cat /etc/nginx/conf.d/keranthos.me.conf | grep -A 10 "location /api/"
   ```
   - 应该看到 `proxy_pass http://127.0.0.1:3000/;`

3. **检查后端服务是否运行：**
   ```bash
   systemctl status blog-backend
   curl http://localhost:3000/api/articles?type=blog&page=1&limit=10
   ```

4. **清除浏览器缓存：**
   - 按 `Ctrl+Shift+R` 强制刷新
   - 或在浏览器中清除缓存后重新访问

5. **检查浏览器控制台网络请求：**
   - 打开浏览器开发者工具（F12）
   - 查看 Network 标签
   - 检查 API 请求的 URL 是否为 `https://keranthos.me/api/...`（正确）
   - 如果是 `http://localhost:3000/api/...`（错误），说明前端需要重新构建

### 问题4：数据库连接失败
- 检查 MySQL/MariaDB 服务：
  - MariaDB: `systemctl status mariadb`
  - MySQL: `systemctl status mysqld`
- 测试连接：`mysql -u blog_user -p my_blog`
- 检查配置文件：`cat /opt/blog/backend/config/config.yml`
- 检查数据库用户权限：`mysql -u root -p -e "SHOW GRANTS FOR 'blog_user'@'localhost';"`

### 问题4.1：SHOW TABLES 没有显示任何表（Empty set）
**原因：** 后端程序还没有运行或者运行失败，表还没有创建。

**解决方法：**
1. **确认后端程序已运行**：
   ```bash
   cd /opt/blog/backend
   GIN_MODE=release ./main
   # 看到 "Listening and serving HTTP on :3000" 说明启动成功
   # 运行几秒后按 Ctrl+C 停止
   ```

2. **检查后端程序运行日志**：
   ```bash
   # 查看是否有错误信息
   cd /opt/blog/backend
   ./main 2>&1 | tee startup.log
   # 查看日志文件
   cat startup.log
   ```

3. **手动检查数据库连接**：
   ```bash
   # 测试连接
   mysql -u blog_user -p my_blog
   # 如果连接失败，检查配置文件中的密码
   ```

4. **检查配置文件**：
   ```bash
   cat /opt/blog/backend/config/config.yml
   # 确认所有配置项都正确，特别是数据库密码
   ```

5. **重新授权数据库用户**（如果权限问题）：
   ```bash
   mysql -u root -p
   GRANT ALL PRIVILEGES ON my_blog.* TO 'blog_user'@'localhost';
   FLUSH PRIVILEGES;
   EXIT;
   ```

### 问题5：图片上传失败
- 检查 uploads 目录权限：`ls -la /opt/blog/backend/uploads`
- 确保目录可写：`chmod -R 755 /opt/blog/backend/uploads`
- 检查磁盘空间：`df -h`

### 问题6：后端编译失败（交叉编译）
- 确保 Go 版本 >= 1.21：`go version`
- 检查环境变量：`echo $GOOS $GOARCH`（应该在 Git Bash 中显示 `linux amd64`）
- 如果使用 PowerShell，确保正确设置：`$env:GOOS="linux"; $env:GOARCH="amd64"`

### 问题7：SSL证书申请失败
- **检查 Cloudflare 代理状态**：申请证书时必须关闭（仅DNS模式）
- 检查域名DNS解析：`nslookup keranthos.me`
- 确保80端口已开放（Let's Encrypt需要验证）
- 检查 Nginx 配置是否正确：`nginx -t`

---

## 完成 ✅

部署完成后，你的博客应该可以通过 https://keranthos.me 访问了！

**建议的后续优化：**
1. 配置 Cloudflare 缓存规则
2. 启用 Cloudflare 的 Auto Minify（HTML/CSS/JS压缩）
3. 配置图片 CDN（如使用对象存储）
4. 设置监控告警（可选）

---

## 快速参考命令

```bash
# 服务管理
systemctl start blog-backend    # 启动后端
systemctl stop blog-backend     # 停止后端
systemctl restart blog-backend   # 重启后端
systemctl status blog-backend   # 查看状态
journalctl -u blog-backend -f   # 查看日志

# Nginx
systemctl reload nginx          # 重新加载配置
nginx -t                        # 测试配置
tail -f /var/log/nginx/error.log # 查看错误日志

# 数据库
mysql -u blog_user -p my_blog   # 连接数据库
systemctl status mariadb        # 查看 MariaDB 状态（或 mysqld 如果是 MySQL）
```
