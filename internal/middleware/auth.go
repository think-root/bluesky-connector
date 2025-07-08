package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/think-root/bluesky-connector/internal/config"
	"github.com/think-root/bluesky-connector/internal/logger"
)

func APIKeyMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		
		if apiKey == "" {
			logger.Warn("Request without API key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing API key",
			})
			c.Abort()
			return
		}

		if apiKey != cfg.Server.APIKey {
			logger.Warnf("Invalid API key: %s", apiKey)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		logger.Debug("API key validated successfully")
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Infof("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		return ""
	})
}