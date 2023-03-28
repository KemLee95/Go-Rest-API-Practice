package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userService "github.com/kemlee/go-rest-api-practise/user"
)

type userSignUpDto struct {
	Email    string `json:"email" binding:"email,required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserController struct {
	userService userService.IUserService
}

func (userContr *UserController) GetUser() {

}

func (userContr *UserController) SignUp(ctx *gin.Context) {
	input := new(userSignUpDto)

	if err := ctx.BindJSON(input); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	exist, err := userContr.userService.CheckEmailExist(ctx.Request.Context(), input.Email)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if exist {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "The email already exists",
		})
		return
	}

	err = userContr.userService.CreateUser(ctx.Request.Context(), &userService.CreateUserRequest{
		Email:    input.Email,
		Name:     input.Name,
		Password: input.Password,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Successfully",
	})
}

func UserControllerRegister(router *gin.Engine, userService userService.IUserService) {
	controller := NewUserController(userService)
	endPoint := router.Group("/user")
	{
		endPoint.POST("/sign-up", controller.SignUp)
	}
}

func NewUserController(userSer userService.IUserService) *UserController {
	return &UserController{
		userService: userSer,
	}
}
