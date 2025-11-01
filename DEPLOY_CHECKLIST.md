# 部署前检查清单

## ✅ 已完成的清理工作

### 1. 废弃代码清理
- ✅ 已删除 `BooksView.vue`、`MoviesView.vue`、`TimelineView.vue`、`CommentsView.vue`
- ✅ 已删除测试页面 `LoadingTestView.vue` 及相关路由
- ✅ 已删除 `waifu-tips.js` 中关于 `/fragments/books`、`/fragments/novels`、`/fragments/movies` 的消息
- ✅ 已同步 `dist` 目录下的 `waifu-tips.js` 与 `public` 目录一致

### 2. 控制台日志
- ✅ 已确认：大部分 `console.log` 已注释
- ✅ 保留 `console.error` 和 `console.warn` 用于生产环境错误追踪
- ✅ 错误上报系统已配置（`errorReporter.js`）

### 3. 看板娘消息优化
- ✅ 首次打开博客时不显示欢迎语
- ✅ 回到首页时不显示消息
- ✅ 三个媒体界面消息已删除

## ⚠️ 部署前需要处理的配置

### 1. 后端配置（backend/config/config.yml）
**⚠️ 安全警告：配置文件中包含敏感信息（数据库密码、JWT密钥）**

**建议操作：**
- [ ] 生产环境使用环境变量替代硬编码密码
- [ ] 确保 `config.yml` 在 `.gitignore` 中（已检查：不在忽略列表）
- [ ] 或创建 `config.yml.example` 模板，生产环境单独配置

**当前配置内容：**
```yaml
database:
  password: "Wan05609"  # ⚠️ 需要改为环境变量
jwt:
  secretkey: "wwj_secret_key"  # ⚠️ 需要改为更强的密钥
```

### 2. CORS 配置（backend/routes/routes.go）
**当前只允许本地开发地址，需要添加生产域名：**

```go
AllowOrigins: []string{
    "http://localhost:3000", 
    "http://localhost:3002",
    // ⚠️ 需要添加生产环境域名，例如：
    // "https://your-domain.com",
}
```

### 3. 前端环境变量（frontend/src/config/api.js）
**已支持环境变量，但需要配置：**

生产环境构建前需要设置：
```bash
export VUE_APP_API_BASE_URL=https://your-domain.com
export VUE_APP_API_PREFIX=/api
```

或创建 `.env.production`：
```
VUE_APP_API_BASE_URL=https://your-domain.com
VUE_APP_API_PREFIX=/api
```

### 4. Gin 生产模式
**确保启动时设置：**
```bash
export GIN_MODE=release
# 或
./main -mode release
```

已在 `DEPLOYMENT.txt` 中说明使用 systemd 时设置环境变量。

## ✅ 已优化的功能

### 1. 代码分割
- ✅ Webpack 已配置代码分割（marked、DOMPurify、Font Awesome、Vue Router/Vuex 单独打包）

### 2. 图片优化
- ✅ 图片懒加载（`loading="lazy" decoding="async"`）
- ✅ 图片错误处理（隐藏而非显示错误图片）
- ✅ 图片尺寸限制（文章：80%宽、750px高；媒体卡片：80%宽、400px高）
- ✅ 图片本地化（外链转本地存储）

### 3. 数据库优化
- ✅ 已添加必要索引（`created_at`、`type` 等）
- ✅ 分页上限校验（最多100条）
- ✅ 慢查询监控

### 4. 安全性
- ✅ 图片上传安全（文件大小、扩展名、MIME类型验证）
- ✅ JWT 认证
- ✅ 全局错误处理
- ✅ Markdown XSS 防护（DOMPurify）

### 5. 用户体验
- ✅ 无限滚动（带防抖和重试）
- ✅ Skeleton 加载效果
- ✅ 错误重试提示
- ✅ Live2D 看板娘
- ✅ 响应式设计

## 📋 部署步骤检查

### 前端构建
```bash
cd frontend
npm ci  # 确保使用 package-lock.json
npm run build
# 将 dist/ 目录部署到服务器
```

### 后端部署
1. [ ] 更新 `backend/config/config.yml` 中的数据库连接信息
2. [ ] 更新 JWT 密钥为更安全的随机字符串
3. [ ] 更新 CORS 允许的域名列表
4. [ ] 编译：`go build -o main .`
5. [ ] 设置 `GIN_MODE=release`
6. [ ] 确保 `uploads/` 目录可写
7. [ ] 启动服务（systemd 或直接运行）

### Nginx 配置
参考 `DEPLOYMENT.txt` 中的配置示例：
- [ ] 前端静态资源路由
- [ ] `/api` 代理到后端
- [ ] `/uploads` 代理到后端
- [ ] HTTPS 配置（推荐使用 Let's Encrypt）

### 数据库
- [ ] 创建生产数据库
- [ ] 导入表结构（GORM AutoMigrate 或手动迁移）
- [ ] 备份策略

### 监控与日志
- [ ] Nginx 日志配置
- [ ] 后端日志路径配置
- [ ] 错误监控（可选：接入 Sentry 等）

## 🔒 安全建议

1. **敏感信息：**
   - 使用环境变量存储数据库密码、JWT密钥
   - 不要在代码仓库中提交 `config.yml`（如果包含真实密码）

2. **HTTPS：**
   - 必须启用 HTTPS（使用 Let's Encrypt 免费证书）
   - 配置 HSTS 头

3. **数据库：**
   - 使用专用数据库用户（非 root）
   - 限制数据库访问 IP
   - 定期备份

4. **文件上传：**
   - 已实现的验证：文件大小、类型、MIME
   - 考虑添加病毒扫描（可选）

5. **访问控制：**
   - JWT 密钥足够复杂（至少32字符随机字符串）
   - 定期轮换密钥（可选）

## 📝 待优化功能（非必须，可后续添加）

1. CDN 集成：图片资源通过 CDN 加速
2. 对象存储：大文件存储到 OSS/COS（用户提到会使用）
3. 图片压缩：上传时自动压缩
4. 缩略图生成：文章列表使用缩略图（用户提到因服务器内存限制暂缓）
5. API 限流：防止恶意请求
6. 日志分析：接入 ELK 或类似系统

## ✅ 部署清单总结

- [x] 废弃代码已清理
- [x] 测试文件已删除
- [x] 看板娘消息已优化
- [x] 代码分割已配置
- [x] 图片优化已完成
- [x] 数据库索引已添加
- [x] 错误处理已完善
- [ ] **配置环境变量**（重要）
- [ ] **更新 CORS 域名**（重要）
- [ ] **更新数据库配置**（重要）
- [ ] **配置 HTTPS**（推荐）
- [ ] **设置备份策略**（推荐）

---

**最后更新：** 2025-01-XX
**状态：** 已准备就绪，需完成配置更新后部署

