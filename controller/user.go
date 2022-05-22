package controller

import (
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/pkg/utils"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := repository.User{
		Username: username,
	}

	// 查找用户是否存在
	userLoginInfo, err := user.FindUserByName(username)
	if userLoginInfo != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	pwd := utils.MakeSha1(password)
	userLoginInfo, err = user.CreateUser(username, pwd)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	generateToken, err := token.GenerateToken(userLoginInfo.Id, userLoginInfo.Username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Success"},
		UserId:   userLoginInfo.Id,
		Token:    generateToken,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := repository.User{
		Username: username,
	}

	signInUser, err := user.FindUserByName(username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	if signInUser == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User does not exist"},
		})
		return
	}
	if signInUser.Password != utils.MakeSha1(password) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password is not correct"},
		})
		return
	}
	generateToken, err := token.GenerateToken(signInUser.Id, signInUser.Username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		UserId:   signInUser.Id,
		Token:    generateToken,
	})
}

func UserInfo(c *gin.Context) {
	//token := c.Query("token")
	claims := c.MustGet("claims").(*token.Claims)
	var user repository.User
	userId := claims.UserId

	userInfo, err := user.FindUserById(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	userRes := User{
		Id:            userId,
		Name:          userInfo.Username,
		FollowCount:   userInfo.FollowCount,
		FollowerCount: userInfo.FollowerCount,
		IsFollow:      userInfo.IsFollow,
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		User:     userRes,
	})
}
