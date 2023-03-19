package handler

import (
	"encoding/json"
	mongodbDriver "go-practise/driver"
	"go-practise/model"
	userRepo "go-practise/repository/repoimpl"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

var jwtKey = []byte("abcdefghijklmnopq")

func CheckHealthHandler(res http.ResponseWriter, req *http.Request) {
	json_data, _ := json.Marshal(model.Error{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	})
	res.Write(json_data)
}

func Register(res http.ResponseWriter, req *http.Request) {
	user_register_data := model.UserRegister{}
	err := json.NewDecoder(req.Body).Decode(&user_register_data)
	if err != nil {
		ResponseErr(res, http.StatusBadRequest)
		return
	}
	mongodb, err := mongodbDriver.GetMongoDb()
	if err != nil {
		panic(err)
	}
	userRepo := userRepo.NewUserRepo(mongodb.Client.Database("gpd"))
	_, err = userRepo.FindUserByEmail(user_register_data.Email)
	if err == nil {
		ResponseErr(res, http.StatusConflict)
		return
	}
	user := model.User{
		Name:     user_register_data.Name,
		Email:    user_register_data.Email,
		Password: user_register_data.Password,
	}
	err = userRepo.Insert(user)
	if err != nil {
		ResponseErr(res, http.StatusInternalServerError)
		return
	}
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	access_token := req.Header.Get("Authorization")
	if access_token == "" {
		ResponseErr(res, http.StatusUnauthorized)
		return
	}

	tk := &Claims{}
	token, err := jwt.ParseWithClaims(access_token, tk, func(jwt_token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		ResponseErr(res, http.StatusUnauthorized)
		return
	}
	if token.Valid {
		ResponseOk(res, token.Claims)
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	user_login_data := model.UserLogin{}
	err := json.NewDecoder(req.Body).Decode(&user_login_data)
	if err != nil {
		ResponseErr(res, http.StatusBadRequest)
		return
	}

	mongodb, err := mongodbDriver.GetMongoDb()
	if err != nil {
		panic(err)
	}
	userRepo := userRepo.NewUserRepo(mongodb.Client.Database("gpd"))
	user, err := userRepo.CheckUserLogin(user_login_data.Email, user_login_data.Password)
	if err != nil {
		ResponseErr(res, http.StatusUnprocessableEntity)
		return
	}
	expire_time := time.Now().Add(1800 * time.Second)
	token, err := GetToken(user, expire_time)
	if err != nil {
		ResponseErr(res, http.StatusInternalServerError)
		return
	}

	ResponseOk(res, model.TokenResponse{
		Token:    token,
		ExpireAt: expire_time.Unix(),
	})
	return
}

func ResponseErr(res http.ResponseWriter, status_code int) {
	json_data, err := json.Marshal(model.Error{
		Status:  status_code,
		Message: http.StatusText(status_code),
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-type", "application/json")
	res.Write(json_data)
}

func ResponseOk(res http.ResponseWriter, data interface{}) {
	if data == nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	json_data, err := json.Marshal(data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-type", "application/json")
	res.Write(json_data)
	return
}

func GetToken(user model.User, expire_time time.Time) (string, error) {
	claims := &Claims{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire_time.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
