package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
	"web_penetration/internal/model"
)

// 响应体写入器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// 日志中间件
func Logging(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 记录日志
		log := &model.APILog{
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Query:        c.Request.URL.RawQuery,
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Status:       c.Writer.Status(),
			Error:        c.Errors.String(),
			Latency:      latency,
			RequestAt:    start,
			RequestBody:  string(requestBody),
			ResponseBody: w.body.String(),
		}

		// 获取用户ID
		if userID, exists := c.Get("user_id"); exists {
			log.UserID = userID.(uint)
		}

		// 异步保存日志
		go saveAPILog(log, db)
	}
}

// 保存API日志
func saveAPILog(log *model.APILog, db *gorm.DB) {
	if err := db.Create(log).Error; err != nil {
		fmt.Printf("Failed to save API log: %v\n", err)
	}
}
