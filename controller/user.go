package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wangzuxianaa/tiktok-simple/pkg/token"
	"github.com/wangzuxianaa/tiktok-simple/service"
	"net/http"
	"strconv"
)

//
// UserLoginResponse
// @Description: 用户登陆响应
//
type UserLoginResponse struct {
	service.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

//
// UserResponse
// @Description: 用户信息的响应
//
type UserResponse struct {
	service.Response
	User service.UserMessage `json:"user"`
}

//
// Register
// @Description: 注册功能
// @param c
//
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 注册
	userId, generateToken, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: service.Response{StatusCode: 0, StatusMsg: "Success"},
		UserId:   userId,
		Token:    generateToken,
	})
}

//
// Login
// @Description: 登陆功能
// @param c
//
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 登陆
	userId, generateToken, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: service.Response{StatusCode: 0, StatusMsg: "success"},
		UserId:   userId,
		Token:    generateToken,
	})
}

//
// UserInfo
// @Description: 用户信息
// @param c
//
func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	claims := c.MustGet("claims").(*token.Claims)
	var err error
	var userId int64
	userId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	var userRes *service.UserMessage
	// 获取用户信息
	userRes, err = service.GetUserInfo(userId, claims.UserId)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: service.Response{StatusCode: 0, StatusMsg: "success"},
		User:     *userRes,
	})
}
