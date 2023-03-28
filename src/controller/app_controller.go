package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppController struct{}

func (contr *AppController) CheckHeath(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success",
	})
}

func NewAppController() *AppController {
	return &AppController{}
}

func AppControllerRegister(router *gin.Engine) {
	appContr := NewAppController()
	router.GET("/check-heath", appContr.CheckHeath)
}
