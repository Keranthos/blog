package middlewares

import (
	"backend/utils"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SlowQueryThreshold 慢查询阈值（默认 500ms）
var SlowQueryThreshold = 500 * time.Millisecond

// SlowQueryMonitor GORM 插件：监控慢查询
// 注意：这个函数需要在使用数据库之前调用
func SlowQueryMonitor(db *gorm.DB) {
	// 在查询前记录开始时间
	err := db.Callback().Query().Before("gorm:query").Register("record_query_start_time", func(db *gorm.DB) {
		if db.Statement != nil {
			ctx := db.Statement.Context
			if ctx == nil {
				ctx = context.Background()
			}
			ctx = context.WithValue(ctx, "query_start_time", time.Now())
			db.Statement.Context = ctx
		}
	})
	if err != nil {
		utils.LogError("注册查询开始时间记录失败: %v", err)
		return
	}

	// 在查询后检查耗时
	err = db.Callback().Query().After("gorm:query").Register("slow_query_monitor", func(db *gorm.DB) {
		if db.Statement != nil && db.Statement.Context != nil {
			// 从上下文中获取查询开始时间
			if startTime, ok := db.Statement.Context.Value("query_start_time").(time.Time); ok {
				duration := time.Since(startTime)
				if duration > SlowQueryThreshold {
					// 记录慢查询
					var query string
					if db.Statement.SQL.String() != "" {
						query = db.Statement.SQL.String()
					} else {
						query = "N/A"
					}
					utils.LogSlowQuery(query, duration, SlowQueryThreshold)
				}
			}
		}
	})
	if err != nil {
		utils.LogError("注册慢查询监控失败: %v", err)
	}
}

// RequestLogger 中间件：记录请求日志和慢请求
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 计算请求耗时
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 记录请求日志
		if latency > 1*time.Second {
			// 慢请求警告
			utils.LogWarn("慢请求: %s %s | 状态: %d | 耗时: %v | IP: %s",
				method, path, statusCode, latency, c.ClientIP())
		} else {
			// 正常请求信息
			utils.LogInfo("%s %s | 状态: %d | 耗时: %v | IP: %s",
				method, path, statusCode, latency, c.ClientIP())
		}
	}
}

