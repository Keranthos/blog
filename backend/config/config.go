package config

import (
	"backend/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		Name     string
	}
	Jwt struct {
		SecretKey  string
		Expiration string
	}
}

var AppConfig *Config //包级变量，初始时直接实例化（根据导出规则，首字母大写的包级变量可以直接在其他包中加上包名使用）
var DB *gorm.DB       //数据库连接的全局变量

func InitConfig() { //只有首字母大写的函数才能导出
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err) //日志相关操作暂不熟悉，暂且这么写
	}

	AppConfig = &Config{} //appconfig不会是nil，可以填充数据
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct,%v", err) //日志相关操作暂不熟悉，暂且这么写
	}

	if AppConfig == nil {
		log.Fatalf("fuck!AppConfig is not initialized")
	}
}

func InitDB() error { // 初始化全局连接池DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		AppConfig.Database.Username,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 关闭查询日志以提高性能（开发时可改为logger.Info）
		Logger: logger.Default.LogMode(logger.Silent),
		// 禁用外键约束检查以提高启动速度
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// 获取底层sql.DB对象来配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting sql.DB: %v", err)
	}

	// 优化连接池配置以提高启动速度
	sqlDB.SetMaxIdleConns(5)                   // 减少最大空闲连接数
	sqlDB.SetMaxOpenConns(20)                  // 减少最大打开连接数
	sqlDB.SetConnMaxIdleTime(time.Minute * 30) // 减少连接最大空闲时间
	sqlDB.SetConnMaxLifetime(time.Hour * 2)    // 减少连接最大生存时间

	// 自动迁移数据库模型
	if err := DB.AutoMigrate(
		&models.BlogArticle{},
		&models.ResearchArticle{},
		&models.ProjectArticle{},
		&models.Comment{},
		&models.Media{},
		&models.Question{},
		&models.User{},
		&models.Moment{},
		&models.Weather{},
		&models.Presentation{},
		&models.Task{},
		&models.Entertainment{},
		&models.DailyReflection{},
		&models.Deadline{},
	); err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	// 创建必要的索引以优化查询性能
	if err := CreateIndexes(); err != nil {
		log.Printf("警告: 创建索引时出现错误: %v", err)
		// 索引创建失败不影响程序运行，仅记录警告
	}

	// 创建默认管理员账户（如果不存在）
	if err := createDefaultAdmin(); err != nil {
		log.Printf("警告: 创建默认管理员账户时出现错误: %v", err)
		// 管理员创建失败不影响程序运行，仅记录警告
	}

	// 注意：慢查询监控在 main.go 中初始化，避免循环依赖

	return nil
}

func GetJWT() (string, string) {
	if AppConfig == nil {
		log.Fatalf("shit!AppConfig is not initialized")
	}
	return AppConfig.Jwt.SecretKey, AppConfig.Jwt.Expiration
}

// createDefaultAdmin 创建默认管理员账户（如果不存在）
// ⚠️ 安全提示：生产环境应禁用此功能或通过环境变量设置初始密码
func createDefaultAdmin() error {
	username := "山角函兽"
	// 从环境变量读取初始密码，如果没有则使用默认值（仅用于首次初始化）
	// 生产环境建议通过环境变量 ADMIN_INITIAL_PASSWORD 设置
	password := "ChangeMe123!" // 默认密码，生产环境必须修改
	if envPassword := os.Getenv("ADMIN_INITIAL_PASSWORD"); envPassword != "" {
		password = envPassword
	}
	level := 3 // 管理员权限

	// 检查用户是否已存在
	var existingUser models.User
	if err := DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		// 用户已存在，跳过创建
		log.Printf("管理员账户已存在: %s (ID: %d)", username, existingUser.ID)
		return nil
	}

	// 生成密码哈希（直接在config包中实现，避免循环依赖）
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %v", err)
	}
	hashedPassword := string(hashedPasswordBytes)

	// 创建新用户
	user := models.User{
		Username: username,
		Password: hashedPassword,
		Level:    level,
		Avatar:   "",
	}

	if err := DB.Create(&user).Error; err != nil {
		return fmt.Errorf("创建管理员账户失败: %v", err)
	}

	log.Printf("默认管理员账户创建成功: %s (ID: %d, Level: %d)", username, user.ID, user.Level)
	return nil
}
