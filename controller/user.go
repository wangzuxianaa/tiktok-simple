package controller

import (
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/pkg/utils"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]UserMessage{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//
// UserLoginResponse
// @Description: 用户登陆响应
//
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

//
// UserResponse
// @Description: 用户信息的响应
//
type UserResponse struct {
	Response
	User UserMessage `json:"user"`
}

//
// Register
// @Description: 注册功能
// @param c
//
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := repository.User{
		Username: username,
	}

	// 查找用户是否存在
	ok, err := user.FindUserByName()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	// 如果用户存在
	if ok == true {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}

	// 密码加密
	pwd := utils.MakeSha1(password)
	user = repository.User{
		Username: username,
		Password: pwd,
	}
	// 创建一条新的用户记录
	err = user.CreateUser()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	// 生成Token
	generateToken, err := token.GenerateToken(user.Id, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Success"},
		UserId:   user.Id,
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

	user := repository.User{
		Username: username,
	}

	// 查找用户是否存在
	ok, err := user.FindUserByName()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	// 用户不存在
	if ok == false {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User does not exist"},
		})
		return
	}
	// 密码校验
	if user.Password != utils.MakeSha1(password) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password is not correct"},
		})
		return
	}

	// 生成token
	generateToken, err := token.GenerateToken(user.Id, username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		UserId:   user.Id,
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

	// string 转 int
	userId, err := strconv.ParseInt(userIdStr, 10, 36)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	user := repository.User{
		Id: userId,
	}
	// 查找用户信息
	err = user.FindUserById()
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	userRes := UserMessage{
		Id:            userId,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		User:     userRes,
	})
}
