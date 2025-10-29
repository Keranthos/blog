package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WeatherController struct {
	DB *gorm.DB
}

func NewWeatherController(db *gorm.DB) *WeatherController {
	return &WeatherController{DB: db}
}

// SeniverseResponse 心知天气API响应结构
type SeniverseResponse struct {
	Results []struct {
		Location struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			Country        string `json:"country"`
			Path           string `json:"path"`
			Timezone       string `json:"timezone"`
			TimezoneOffset string `json:"timezone_offset"`
		} `json:"location"`
		Now struct {
			Text        string `json:"text"`        // 天气现象文字
			Code        string `json:"code"`        // 天气现象代码
			Temperature string `json:"temperature"` // 温度
			FeelsLike   string `json:"feels_like"`  // 体感温度
			Pressure    string `json:"pressure"`    // 气压
			Humidity    string `json:"humidity"`    // 相对湿度
			Visibility  string `json:"visibility"`  // 能见度
			Wind        struct {
				Direction string `json:"direction"` // 风向
				Speed     string `json:"speed"`     // 风速
			} `json:"wind"`
		} `json:"now"`
		LastUpdate string `json:"last_update"` // 数据更新时间
	} `json:"results"`
}

// SeniverseForecastResponse 心知天气预报API响应结构
type SeniverseForecastResponse struct {
	Results []struct {
		Location struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			Country        string `json:"country"`
			Path           string `json:"path"`
			Timezone       string `json:"timezone"`
			TimezoneOffset string `json:"timezone_offset"`
		} `json:"location"`
		Daily []struct {
			Date                string `json:"date"`                  // 日期
			TextDay             string `json:"text_day"`              // 白天天气现象文字
			CodeDay             string `json:"code_day"`              // 白天天气现象代码
			TextNight           string `json:"text_night"`            // 晚间天气现象文字
			CodeNight           string `json:"code_night"`            // 晚间天气现象代码
			High                string `json:"high"`                  // 最高温度
			Low                 string `json:"low"`                   // 最低温度
			Rainfall            string `json:"rainfall"`              // 降雨量
			Precip              string `json:"precip"`                // 降水概率
			WindDirection       string `json:"wind_direction"`        // 风向
			WindDirectionDegree string `json:"wind_direction_degree"` // 风向角度
			WindSpeed           string `json:"wind_speed"`            // 风速
			WindScale           string `json:"wind_scale"`            // 风力等级
			Humidity            string `json:"humidity"`              // 湿度
		} `json:"daily"`
		LastUpdate string `json:"last_update"` // 数据更新时间
	} `json:"results"`
}

// SeniverseLifeResponse 心知天气生活指数API响应结构
type SeniverseLifeResponse struct {
	Results []struct {
		Location struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			Country        string `json:"country"`
			Path           string `json:"path"`
			Timezone       string `json:"timezone"`
			TimezoneOffset string `json:"timezone_offset"`
		} `json:"location"`
		Suggestion []struct {
			Name   string `json:"name"`   // 指数名称
			Brief  string `json:"brief"`  // 指数简介
			Detail string `json:"detail"` // 指数详情
		} `json:"suggestion"`
		LastUpdate string `json:"last_update"` // 数据更新时间
	} `json:"results"`
}

// GetCurrentWeather 获取当前天气信息
func (wc *WeatherController) GetCurrentWeather(c *gin.Context) {
	var weather models.Weather

	// 查找最新的天气记录
	err := wc.DB.Order("last_updated DESC").First(&weather).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有天气记录，返回默认信息
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"location":     "北京",
					"temperature":  "22°",
					"weather":      "多云",
					"wind":         "2级",
					"humidity":     "65%",
					"air_quality":  "良",
					"last_updated": time.Now().Format("2006-01-02 15:04:05"),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	// 检查是否需要更新天气数据
	if weather.IsExpired() {
		// 天气数据过期，异步更新
		go wc.updateWeatherData(weather.City, weather.Latitude, weather.Longitude)
	}

	// 返回当前天气信息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"location":     weather.City,
			"temperature":  fmt.Sprintf("%.0f°", weather.Temperature),
			"weather":      weather.Weather,
			"description":  weather.Description,
			"wind":         fmt.Sprintf("%.1fm/s", weather.WindSpeed),
			"humidity":     fmt.Sprintf("%d%%", weather.Humidity),
			"air_quality":  weather.AirQuality,
			"last_updated": weather.LastUpdated.Format("2006-01-02 15:04:05"),
			"tomorrow":     weather.TomorrowForecast, // 使用数据库中的预报数据
			"life_index":   weather.LifeIndex,        // 使用数据库中的生活指数
		},
	})
}

// UpdateLocation 更新博主位置
func (wc *WeatherController) UpdateLocation(c *gin.Context) {
	var request struct {
		City      string  `json:"city" binding:"required"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 更新位置并获取新天气数据
	err := wc.updateWeatherData(request.City, request.Latitude, request.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新天气数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "位置更新成功",
		"userToast": true, // 需要看板娘提示
	})
}

// updateWeatherData 更新天气数据
func (wc *WeatherController) updateWeatherData(city string, lat, lon float64) error {
	// 获取心知天气API密钥
	publicKey := os.Getenv("SENIVERSE_PUBLIC_KEY")
	privateKey := os.Getenv("SENIVERSE_PRIVATE_KEY")

	if publicKey == "" || privateKey == "" {
		// 如果没有配置环境变量，使用你提供的密钥
		publicKey = "PN5XMAu0u11Nwk1Zu"
		privateKey = "S5PhMOwI5CpN7aXXH"
	}

	// 城市名称映射（中文到英文）
	cityMap := map[string]string{
		"北京": "beijing",
		"上海": "shanghai",
		"广州": "guangzhou",
		"深圳": "shenzhen",
		"杭州": "hangzhou",
		"南京": "nanjing",
		"武汉": "wuhan",
		"成都": "chengdu",
		"西安": "xian",
		"重庆": "chongqing",
		"珠海": "zhuhai",
	}

	// 获取英文城市名称
	englishCity := cityMap[city]
	if englishCity == "" {
		englishCity = city // 如果没有映射，使用原名称
	}

	// 调用心知天气API - 获取实况天气（使用私钥）
	nowURL := fmt.Sprintf("https://api.seniverse.com/v3/weather/now.json?key=%s&location=%s&language=zh-Hans&unit=c",
		privateKey, englishCity)

	// TODO: 后续可以添加预报和生活指数功能
	// forecastURL := fmt.Sprintf("https://api.seniverse.com/v3/weather/daily.json?key=%s&location=%s&language=zh-Hans&unit=c&start=0&days=3",
	//		publicKey, city)
	// lifeURL := fmt.Sprintf("https://api.seniverse.com/v3/life/suggestion.json?key=%s&location=%s&language=zh-Hans",
	//		publicKey, city)

	resp, err := http.Get(nowURL)
	if err != nil {
		return fmt.Errorf("天气API调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("天气API返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取API响应失败: %v", err)
	}

	var weatherResp SeniverseResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return fmt.Errorf("解析API响应失败: %v", err)
	}

	// 检查API响应是否有效
	if len(weatherResp.Results) == 0 {
		return fmt.Errorf("心知天气API返回空结果")
	}

	result := weatherResp.Results[0]

	// 解析温度
	temperature, err := strconv.ParseFloat(result.Now.Temperature, 64)
	if err != nil {
		temperature = 0
	}

	// 心知天气免费版不包含湿度、气压、能见度、风速等详细信息
	// 使用默认值或模拟数据
	humidity := 65   // 默认湿度
	pressure := 1013 // 默认气压
	visibility := 10 // 默认能见度
	windSpeed := 5.0 // 默认风速

	// 获取3天预报数据
	forecastURL := fmt.Sprintf("https://api.seniverse.com/v3/weather/daily.json?key=%s&location=%s&language=zh-Hans&unit=c&start=0&days=3",
		privateKey, englishCity)

	var tomorrowForecast string
	var lifeIndex string

	// 获取预报数据
	forecastResp, err := http.Get(forecastURL)
	if err == nil && forecastResp.StatusCode == http.StatusOK {
		forecastBody, err := io.ReadAll(forecastResp.Body)
		forecastResp.Body.Close()

		if err == nil {
			var forecastData SeniverseForecastResponse
			if err := json.Unmarshal(forecastBody, &forecastData); err == nil && len(forecastData.Results) > 0 && len(forecastData.Results[0].Daily) > 1 {
				// 获取明天的数据（索引1，因为索引0是今天）
				tomorrow := forecastData.Results[0].Daily[1]
				tomorrowForecast = fmt.Sprintf("%s %s°C/%s°C", tomorrow.TextDay, tomorrow.High, tomorrow.Low)
			}
		}
	}

	// 如果获取预报失败，使用默认值
	if tomorrowForecast == "" {
		tomorrowForecast = "晴 27°C/33°C" // 珠海明天的默认值
	}

	// 设置默认生活指数
	lifeIndex = "适宜出行"

	// 创建或更新天气记录
	weather := models.Weather{
		City:        result.Location.Name,
		Country:     result.Location.Country,
		Latitude:    lat, // 使用传入的坐标
		Longitude:   lon, // 使用传入的坐标
		Temperature: temperature,
		Weather:     result.Now.Text,
		Description: result.Now.Text,
		WindSpeed:   windSpeed,
		Humidity:    humidity,
		Pressure:    pressure,
		Visibility:  visibility,
		AirQuality:  "良好", // 心知天气免费版不包含空气质量
		LastUpdated: time.Now(),
		// 新增字段
		TomorrowForecast: tomorrowForecast,
		LifeIndex:        lifeIndex,
	}

	// 删除旧的天气记录（只保留最新的）
	wc.DB.Where("1 = 1").Delete(&models.Weather{})

	// 创建新记录
	if err := wc.DB.Create(&weather).Error; err != nil {
		return fmt.Errorf("保存天气数据失败: %v", err)
	}

	return nil
}
