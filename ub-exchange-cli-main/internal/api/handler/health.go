package handler

import (
	"context"
	"net/http"
	"time"

	"exchange-go/internal/platform"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthCheck is a liveness probe — returns 200 if the process is alive.
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "OK",
		})
	}
}

// ReadinessCheck is a readiness probe — returns 200 only if DB and Redis are reachable.
func ReadinessCheck(db *gorm.DB, rc platform.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		checks := make(map[string]string)

		// Check database
		sqlDB, err := db.DB()
		if err != nil {
			checks["db"] = err.Error()
		} else if err := sqlDB.PingContext(ctx); err != nil {
			checks["db"] = err.Error()
		} else {
			checks["db"] = "ok"
		}

		// Check Redis
		if err := rc.Set(ctx, "_health:ping", "1"); err != nil {
			checks["redis"] = err.Error()
		} else {
			checks["redis"] = "ok"
		}

		if checks["db"] != "ok" || checks["redis"] != "ok" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "not ready",
				"data":    checks,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "ready",
			"data":    checks,
		})
	}
}
