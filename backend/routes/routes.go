package routes

import (
	"backend/config"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 创建天气控制器实例
	weatherController := controllers.NewWeatherController(config.DB)

	// 自定义 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3002"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // 允许的请求头，包括 Authorization
		AllowCredentials: true,                                      // 允许携带凭证（如 cookies 或授权头）
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 公共路由，不需要验证
		api.POST("/auth/login", controllers.Login)       // 登录
		api.POST("/auth/register", controllers.Register) // 注册

		// 需要验证的路由
		// 用户更新自己的信息，权限: 1 (普通用户及以上)
		api.POST("/auth/updateUser", middlewares.Auth(1), controllers.UpdateUser)

		// 文章相关路由，权限: 0 (所有用户)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/count", controllers.GetArticleCount)
		api.GET("/articles/:id", controllers.GetArticleById)
		api.GET("/articles/top", controllers.GetTopArticles)                        // 获取置顶文章
		api.POST("/articles", middlewares.Auth(3), controllers.CreateArticle)       // 创建文章需要管理员权限
		api.PUT("/articles/:id", middlewares.Auth(3), controllers.UpdateArticle)    // 更新文章需要管理员权限
		api.DELETE("/articles/:id", middlewares.Auth(3), controllers.DeleteArticle) // 删除文章需要管理员权限

		// 评论相关路由，权限: 0 (所有用户)
		api.GET("/comments/:blogID", controllers.GetComments)
		api.POST("/comments", controllers.CreateComment)
		api.DELETE("/comments/:commentId", middlewares.Auth(3), controllers.DeleteComment) // 删除评论需要管理员权限

		// 媒体相关路由，权限: 3 (管理员)
		api.GET("/media", controllers.GetMedia)
		api.GET("/media/:mediaId", controllers.GetMediaByID)
		api.POST("/media", middlewares.Auth(3), controllers.CreateMedia)
		api.PUT("/media/:mediaId", middlewares.Auth(3), controllers.UpdateMedia)
		api.DELETE("/media/:mediaId", middlewares.Auth(3), controllers.DeleteMedia) // 删除媒体需要管理员权限

		// 问题相关路由，权限: 0 (所有用户)
		api.GET("/questions", controllers.GetQuestions)
		api.POST("/questions", controllers.CreateQuestion)
		api.POST("/questions/:questionId/answer", middlewares.Auth(3), controllers.AnswerQuestion) // 回答问题需要编辑者权限

		// 碎碎念相关路由
		api.GET("/moments", controllers.GetMoments)
		api.GET("/moments/count", controllers.GetMomentsCount)
		api.GET("/moments/:id", controllers.GetMoment)
		api.POST("/moments", middlewares.Auth(3), controllers.CreateMoment)       // 创建碎碎念需要管理员权限
		api.PUT("/moments/:id", middlewares.Auth(3), controllers.UpdateMoment)    // 更新碎碎念需要管理员权限
		api.DELETE("/moments/:id", middlewares.Auth(3), controllers.DeleteMoment) // 删除碎碎念需要管理员权限

		// 文件上传相关路由
		api.POST("/upload/image", middlewares.Auth(3), controllers.UploadImage) // 上传图片需要管理员权限

		// 图片管理相关路由
		api.GET("/images", middlewares.Auth(3), controllers.GetImages)             // 获取图片列表需要管理员权限
		api.DELETE("/images", middlewares.Auth(3), controllers.DeleteImage)        // 删除图片需要管理员权限
		api.POST("/random-image", middlewares.Auth(1), controllers.GetRandomImage) // 获取随机图片需要用户权限

		// RSS相关路由
		api.GET("/rss", controllers.GenerateRSS)             // RSS订阅源
		api.GET("/rss/stats", controllers.GetRSSStats)       // RSS统计信息
		api.GET("/sitemap.xml", controllers.GenerateSitemap) // Sitemap

		// 天气相关路由
		api.GET("/weather", weatherController.GetCurrentWeather)                                    // 获取当前天气（公开）
		api.POST("/weather/update-location", middlewares.Auth(3), weatherController.UpdateLocation) // 更新位置（需要管理员权限）
	}

	// 静态文件服务（提供上传的图片访问）
	r.Static("/uploads", "./uploads")

	// RSS订阅源（直接访问）
	r.GET("/rss", controllers.GenerateRSS)

	// Sitemap（直接访问）
	r.GET("/sitemap.xml", controllers.GenerateSitemap)

	return r
}
