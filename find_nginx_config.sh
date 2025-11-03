#!/bin/bash
# 查找 Nginx 配置文件位置

echo "查找 Nginx 主配置文件..."
nginx -t 2>&1 | grep "configuration file"

echo ""
echo "查找所有 Nginx 配置文件..."
find /etc -name "nginx.conf" 2>/dev/null

echo ""
echo "检查 Nginx 配置目录结构..."
ls -la /etc/nginx/ 2>/dev/null || echo "目录不存在"

echo ""
echo "检查 Nginx 是否安装..."
which nginx || rpm -qa | grep nginx || yum list installed | grep nginx

