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

	// 请求日志中间件（记录所有请求和慢请求）
	r.Use(middlewares.RequestLogger())

	// 创建控制器实例
	weatherController := controllers.NewWeatherController(config.DB)
	presentationController := controllers.NewPresentationController(config.DB)
	plannerController := controllers.NewPlannerController(config.DB)

	// 自定义 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3002",
			"https://keranthos.me",
			"https://www.keranthos.me",
		}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"}, // 允许的请求头
		AllowCredentials: true,                                                                              // 允许携带凭证（如 cookies 或授权头）
		MaxAge:           12 * 60 * 60,                                                                      // 预检请求缓存时间
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 公共路由，不需要验证
		api.POST("/auth/login", controllers.Login)           // 登录
		api.POST("/auth/register", controllers.Register)     // 注册
		api.GET("/image-apis", controllers.GetAvailableAPIs) // 获取可用的图片API列表（公开）

		// 需要验证的路由
		// 用户更新自己的信息，权限: 1 (普通用户及以上)
		api.POST("/auth/updateUser", middlewares.Auth(1), controllers.UpdateUser)

		// 文章相关路由，权限: 0 (所有用户)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/count", controllers.GetArticleCount)
		api.GET("/articles/:id", controllers.GetArticleById)
		api.GET("/articles/top", controllers.GetTopArticles)                                        // 获取置顶文章
		api.POST("/articles", middlewares.Auth(3), controllers.CreateArticle)                       // 创建文章需要管理员权限
		api.PUT("/articles/:id", middlewares.Auth(3), controllers.UpdateArticle)                    // 更新文章需要管理员权限
		api.DELETE("/articles/:id", middlewares.Auth(3), controllers.DeleteArticle)                 // 删除文章需要管理员权限
		api.POST("/articles/migrate-images", middlewares.Auth(3), controllers.MigrateArticleImages) // 批量本地化历史封面并替换指定文章

		// 评论相关路由，权限: 0 (所有用户)
		api.GET("/comments", controllers.GetAllComments) // 获取所有评论
		api.GET("/comments/:blogID", controllers.GetComments)
		api.POST("/comments", controllers.CreateComment)
		api.DELETE("/comments/:commentId", middlewares.Auth(3), controllers.DeleteComment) // 删除评论需要管理员权限

		// 媒体相关路由，权限: 3 (管理员)
		api.GET("/media", controllers.GetMedia)
		api.GET("/media/:mediaId", controllers.GetMediaByID)
		api.POST("/media", middlewares.Auth(3), controllers.CreateMedia)
		api.PUT("/media/:mediaId", middlewares.Auth(3), controllers.UpdateMedia)
		api.DELETE("/media/:mediaId", middlewares.Auth(3), controllers.DeleteMedia)            // 删除媒体需要管理员权限
		api.POST("/media/migrate-images", middlewares.Auth(3), controllers.MigrateMediaImages) // 批量本地化媒体封面

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
		api.GET("/images", middlewares.Auth(3), controllers.GetImages)                       // 获取图片列表需要管理员权限
		api.GET("/images/usages", middlewares.Auth(3), controllers.GetImageUsages)           // 获取图片使用位置需要管理员权限
		api.DELETE("/images", middlewares.Auth(3), controllers.DeleteImage)                  // 删除图片需要管理员权限
		api.POST("/random-image", middlewares.Auth(1), controllers.GetRandomImage)           // 获取随机图片需要用户权限
		api.POST("/upload/custom-image", middlewares.Auth(1), controllers.CustomImageUpload) // 上传自定义图片需要用户权限

		// 天气相关路由
		api.GET("/weather", weatherController.GetCurrentWeather)                                    // 获取当前天气（公开）
		api.POST("/weather/update-location", middlewares.Auth(3), weatherController.UpdateLocation) // 更新位置（需要管理员权限）

		// 讲演相关路由
		api.GET("/presentations", presentationController.GetPresentations)                                // 获取讲演列表（公开）
		api.GET("/presentations/:id", presentationController.GetPresentation)                             // 获取讲演详情（公开）
		api.GET("/presentations/:id/download", presentationController.ServePresentationFile)              // 下载讲演文件（公开）
		api.GET("/presentations/:id/preview", presentationController.PreviewPresentationFile)             // 预览讲演文件（公开）
		api.GET("/presentations/:id/info", presentationController.GetPresentationInfo)                    // 获取讲演文件信息（公开）
		api.POST("/presentations/:id/verify-password", presentationController.VerifyPresentationPassword) // 验证讲演密码（公开）
		api.POST("/presentations", middlewares.Auth(1), presentationController.CreatePresentation)        // 创建讲演需要用户权限
		api.PUT("/presentations/:id", middlewares.Auth(1), presentationController.UpdatePresentation)     // 更新讲演需要用户权限
		api.DELETE("/presentations/:id", middlewares.Auth(1), presentationController.DeletePresentation)  // 删除讲演需要用户权限

		// 计划管理相关路由（需要管理员权限）
		// 任务相关
		api.GET("/planner/tasks", middlewares.Auth(3), plannerController.GetTasks)          // 获取任务列表
		api.POST("/planner/tasks", middlewares.Auth(3), plannerController.CreateTask)       // 创建任务
		api.PUT("/planner/tasks/:id", middlewares.Auth(3), plannerController.UpdateTask)    // 更新任务
		api.DELETE("/planner/tasks/:id", middlewares.Auth(3), plannerController.DeleteTask) // 删除任务

		// 娱乐活动相关
		api.GET("/planner/entertainments", middlewares.Auth(3), plannerController.GetEntertainments)          // 获取娱乐活动列表
		api.POST("/planner/entertainments", middlewares.Auth(3), plannerController.CreateEntertainment)       // 创建娱乐活动
		api.PUT("/planner/entertainments/:id", middlewares.Auth(3), plannerController.UpdateEntertainment)    // 更新娱乐活动
		api.DELETE("/planner/entertainments/:id", middlewares.Auth(3), plannerController.DeleteEntertainment) // 删除娱乐活动

		// 每日反思相关
		api.GET("/planner/reflection", middlewares.Auth(3), plannerController.GetReflection)             // 获取每日反思
		api.POST("/planner/reflection", middlewares.Auth(3), plannerController.CreateOrUpdateReflection) // 创建或更新每日反思
		api.GET("/planner/reflections", middlewares.Auth(3), plannerController.GetReflections)           // 获取反思列表

		// 截止日期相关
		api.GET("/planner/deadlines", middlewares.Auth(3), plannerController.GetDeadlines)          // 获取截止日期列表
		api.POST("/planner/deadlines", middlewares.Auth(3), plannerController.CreateDeadline)       // 创建截止日期
		api.PUT("/planner/deadlines/:id", middlewares.Auth(3), plannerController.UpdateDeadline)    // 更新截止日期
		api.DELETE("/planner/deadlines/:id", middlewares.Auth(3), plannerController.DeleteDeadline) // 删除截止日期
	}

	// 静态文件服务（提供上传的图片访问）
	r.Static("/uploads", "./uploads")

	return r
}
