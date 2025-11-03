# Cloudflare + Nginx 配置步骤

## 1. 创建 Cloudflare IP 配置文件

在服务器上执行：

```bash
vim /etc/nginx/cloudflare.conf
```

将以下内容复制进去：

```
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

保存退出（`:wq`）

## 2. 修改 Nginx 主配置文件

```bash
vim /etc/nginx/nginx.conf
```

在 `http {` 块中添加：

```
include /etc/nginx/cloudflare.conf;
```

例如：

```
http {
    include /etc/nginx/cloudflare.conf;
    
    # ... 其他配置 ...
}
```

保存退出

## 3. 更新站点配置文件

```bash
vim /etc/nginx/conf.d/keranthos.me.conf
```

替换为以下内容：

```
server {
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
        
        # 使用 Cloudflare 真实 IP
        proxy_set_header X-Real-IP $http_cf_connecting_ip;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Cloudflare 相关头部
        proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
        proxy_set_header CF-Ray $http_cf_ray;
        proxy_set_header CF-Visitor $http_cf_visitor;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 后端静态上传文件
    location /uploads/ {
        alias /opt/blog/backend/uploads/;
        
        # 缓存设置
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

    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/keranthos.me/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/keranthos.me/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
}

server {
    if ($host = www.keranthos.me) {
        return 301 https://$host$request_uri;
    }

    if ($host = keranthos.me) {
        return 301 https://$host$request_uri;
    }

    listen 80;
    server_name keranthos.me www.keranthos.me;
    return 404;
}
```

保存退出

## 4. 测试并重新加载

```bash
# 测试配置
nginx -t

# 重新加载
systemctl reload nginx
```

完成！

