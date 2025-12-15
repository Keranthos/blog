package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("注册函数发生panic: %v, 堆栈: %s", r, getStackTrace())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("服务器内部错误: %v", r)})
		}
	}()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("注册请求JSON绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误: " + err.Error()})
		return
	}

	log.Printf("注册请求 - 用户名: %s, 密码长度: %d, Avatar: %s", user.Username, len(user.Password), user.Avatar)

	// 去除用户名首尾空格（验证函数内部也会处理，但这里提前处理确保一致性）
	user.Username = strings.TrimSpace(user.Username)

	// 使用新的验证函数
	if err := utils.ValidateUsername(user.Username); err != nil {
		log.Printf("用户名验证失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidatePassword(user.Password); err != nil {
		log.Printf("密码验证失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		log.Printf("用户名已存在: %s", user.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器加密错误: " + err.Error()})
		return
	}
	
	// 创建新用户对象，确保所有字段都正确设置
	newUser := models.User{
		Username: user.Username,
		Password: hashedPassword,
		Level:    2, // 默认用户等级（普通用户）
		Avatar:   user.Avatar, // 保持原值，如果为空就是空字符串
	}

	log.Printf("准备创建用户 - Username: %s, Level: %d, Avatar: %s", newUser.Username, newUser.Level, newUser.Avatar)

	if err := config.DB.Create(&newUser).Error; err != nil {
		log.Printf("创建用户失败 - 数据库错误: %v, 用户数据: Username=%s, Level=%d, Avatar=%s", err, newUser.Username, newUser.Level, newUser.Avatar)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败: " + err.Error()})
		return
	}

	log.Printf("用户创建成功 - ID: %d, Username: %s", newUser.ID, newUser.Username)

	token, err := utils.GenerateJWT(newUser.Username, 2)
	if err != nil {
		log.Printf("Token生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token生成失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"isLogged": true,
			"name":     newUser.Username,
			"level":    2,
			"avatar":   newUser.Avatar,
		},
		"jwt": token,
	})
}

// getStackTrace 获取堆栈跟踪信息（简化版）
func getStackTrace() string {
	return "stack trace unavailable"
}

func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"` //Username必须用大写才可以导出匹配，所以在后面加上json:"username"，注意这两者都是大小写敏感的
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "传入格式不对"})
		return
	}

	var dbUser models.User
	if err := config.DB.Where("username = ?", loginData.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名错误"})
		return
	}

	if !utils.CheckPasswordHash(loginData.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token, err := utils.GenerateJWT(dbUser.Username, dbUser.Level)
	log.Println(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"isLogged": true,
			"name":     dbUser.Username,
			"level":    dbUser.Level,
			"avatar":   dbUser.Avatar != "",
		},
		"jwt": token,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		dbUser.Password = hashedPassword
	}

	if user.Avatar != "" {
		dbUser.Avatar = user.Avatar
	}

	if err := config.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	token, err := utils.GenerateJWT(dbUser.Username, dbUser.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"isLogged": true,
			"name":     dbUser.Username,
			"level":    dbUser.Level,
			"avatar":   dbUser.Avatar,
		},
		"jwt": token,
	})
}
