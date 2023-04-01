package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	app "github.com/kemlee/go-rest-api-practise/app"
	config "github.com/kemlee/go-rest-api-practise/config"
	collectionHelper "github.com/kemlee/go-rest-api-practise/core/collection_helper"
)

func LanguageMiddleware(apiConfig *config.ApiConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.Request
		requestedLanguageStr := request.Header.Get(app.LANGUAGE_HEADER_PARAM)
		appLanguages := apiConfig.DefaultLanguages
		if requestedLanguageStr != "" {
			var setAppLanguage []string
			requestedLanguages := strings.Split(requestedLanguageStr, ",")
			requestedLanguages = collectionHelper.Unique(requestedLanguages)
			for _, lang := range requestedLanguages {
				if lang == string(app.ENGLISH) || lang == string(app.CHINA) {
					setAppLanguage = append(setAppLanguage, lang)
				}
			}
			appLanguages = strings.Join(setAppLanguage, ",")
		}
		request.Header.Set(app.LANGUAGE_HEADER_PARAM, appLanguages)
		ctx.Next()
	}
}
