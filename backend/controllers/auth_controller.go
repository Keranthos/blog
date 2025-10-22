package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 使用新的验证函数
	if err := utils.ValidateUsername(user.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidatePassword(user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器加密错误"})
		return
	}
	user.Password = hashedPassword

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	token, err := utils.GenerateJWT(user.Username, 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"isLogged": true,
			"name":     user.Username,
			"level":    2,
			"avatar":   user.Avatar,
		},
		"jwt": token,
	})
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
