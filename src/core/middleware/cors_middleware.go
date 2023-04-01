package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	config "github.com/kemlee/go-rest-api-practise/config"
)

func CorsMiddleware(apiConfig *config.ApiConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", strings.Join(apiConfig.Cors.AllowOrigin, ","))
		ctx.Header("Access-Control-Request-Method", strings.Join(apiConfig.Cors.AllowMethod, ","))
		ctx.Header("Access-Control-Request-Headers", strings.Join(apiConfig.Cors.AllowHeader, ","))
	}
}
