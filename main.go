package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Auth struct {
		GenerateTokens []string `yaml:"generate_tokens"`
		ModelTokens    []string `yaml:"model_tokens"`
	} `yaml:"auth"`
	Service struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"service"`
}

const (
	configFile = "config.yaml"
)

func main() {
	// 读取配置文件
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	// 添加全局异常处理
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 500,
					"msg":  "服务器异常",
				})
				c.Abort()
			}
		}()
		c.Next()
	})

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 系统管理路由组，用于模型管理
	sys := r.Group("/sys", authMiddleware(*config))
	{
		// 模型相关接口
		sys.GET("/tags", proxyOllama("/api/tags"))
		sys.POST("/pull", proxyOllama("/api/pull"))
		sys.DELETE("/delete", proxyOllama("/api/delete"))
		sys.POST("/copy", proxyOllama("/api/copy"))
		sys.POST("/push", proxyOllama("/api/push"))
		sys.GET("/show", proxyOllama("/api/show"))
	}

	// API路由组，用于生成相关功能
	api := r.Group("/api", authMiddleware(*config))
	{
		// 生成相关接口
		api.POST("/generate", proxyOllama("/api/generate"))
		api.POST("/chat", proxyOllama("/api/chat"))
		api.POST("/embeddings", proxyOllama("/api/embeddings"))
	}

	r.Run(":8080")
}

func loadConfig() (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 解析YAML配置
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// authMiddleware 认证中间件
func authMiddleware(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "未提供认证token",
			})
			c.Abort()
			return
		}

		// 根据路由组判断使用哪种token验证
		validToken := false
		if strings.HasPrefix(c.Request.URL.Path, "/sys") {
			// 系统管理接口使用模型管理token
			for _, allowedToken := range config.Auth.ModelTokens {
				if token == allowedToken {
					validToken = true
					break
				}
			}
		} else if strings.HasPrefix(c.Request.URL.Path, "/api") {
			// 生成相关接口使用生成token
			for _, allowedToken := range config.Auth.GenerateTokens {
				if token == allowedToken {
					validToken = true
					break
				}
			}
		}

		if !validToken {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "非授权访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// proxyOllama 创建一个代理处理函数
func proxyOllama(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := loadConfig()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "内部服务器错误",
				"data":    nil,
			})
			return
		}

		// 创建代理请求
		baseURL := "http://localhost:11434"
		if config.Service.BaseURL != "" {
			baseURL = config.Service.BaseURL
		}
		req, err := http.NewRequest(c.Request.Method, baseURL+path, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "内部服务器错误",
				"data":    nil,
			})
			return
		}

		// 复制原始请求的header
		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// 发送请求到Ollama服务
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "内部服务器错误",
				"data":    nil,
			})
			return
		}
		defer resp.Body.Close()

		// 复制响应header
		for name, values := range resp.Header {
			for _, value := range values {
				c.Header(name, value)
			}
		}

		// 设置响应状态码
		c.Status(resp.StatusCode)

		// 复制响应体
		c.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, resp.Body)
			return err == nil
		})
	}
}
