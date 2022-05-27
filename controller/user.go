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
	ok, err := user.FindUserByName()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	if ok == true {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}

	pwd := utils.MakeSha1(password)
	user = repository.User{
		Username: username,
		Password: pwd,
	}
	err = user.CreateUser()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
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

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := repository.User{
		Username: username,
	}

	ok, err := user.FindUserByName()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	if ok == false {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User does not exist"},
		})
		return
	}
	if user.Password != utils.MakeSha1(password) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password is not correct"},
		})
		return
	}
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

func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")

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
	err = user.FindUserById()
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	userRes := User{
		Id:            userId,
		Name:          user.Username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		User:     userRes,
	})
}
