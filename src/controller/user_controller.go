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

type userLoginDto struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}
type ResponseData struct {
	status  bool
	message string
	data    interface{}
}
type getUserSerialization struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserController struct {
	userService userService.IUserService
}

func (userContr *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	exist, err := userContr.userService.CheckIdExist(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !exist {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "The user not found",
			"data":    nil,
		})
		return
	}

	user, err := userContr.userService.GetUserById(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userSerialization := &getUserSerialization{
		Email: user.Email,
		Name:  user.Name,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Retrieve successfully",
		"data":    userSerialization,
	})
}

func (userContr *UserController) ListUser(ctx *gin.Context) {
	users, err := userContr.userService.GetUserList(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var usersSerialization []*getUserSerialization
	for _, user := range users {
		usersSerialization = append(usersSerialization, &getUserSerialization{
			Email: user.Email,
			Name:  user.Name,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Retrieve successfully",
		"data":    usersSerialization,
	})
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

func (contr *UserController) SignIn(ctx *gin.Context) {
	input := &userLoginDto{}
	if err := ctx.BindJSON(input); err != nil {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	token, err := contr.userService.SignIn(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Successfully",
		"token":   token,
	})
}

func NewUserController(userSer userService.IUserService) *UserController {
	return &UserController{
		userService: userSer,
	}
}

func UserControllerRegister(router *gin.Engine, userService userService.IUserService) {
	controller := NewUserController(userService)
	endPoint := router.Group("/user")
	{
		endPoint.POST("/sign-up", controller.SignUp)
		endPoint.POST("/login", controller.SignIn)
		endPoint.GET("/list", controller.ListUser)
		endPoint.GET("/get/:id", controller.GetUser)
	}
}
