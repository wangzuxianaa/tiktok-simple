package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []UserMessage `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	claims := c.MustGet("claims").(*token.Claims)
	fansId := claims.UserId
	followIdStr := c.Query("to_user_id")
	isFollowStr := c.Query("action_type")

	followId, err := strconv.ParseInt(followIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	var isFollow bool
	if comp := strings.Compare(isFollowStr, "1"); comp == 0 {
		isFollow = true
	} else {
		isFollow = false
	}

	followUser := repository.User{
		Id: followId,
	}
	fans := repository.User{
		Id: fansId,
	}

	err = followUser.FindUserById()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	err = fans.FindUserById()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	follow := repository.Follow{
		FollowId: followId,
		FansId:   fansId,
	}
	if isFollow {
		exist, err := follow.FindFollow()
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		if !exist {
			//创建新的关注列
			follow.IsFollow = true
			err = follow.CreateFollow()
			if err != nil {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: err.Error()},
				})
				return
			}
			err = follow.UpdateFollow(true)
			if err != nil {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: err.Error()},
				})
				return
			}
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "Follow success",
		})
	} else {
		exist, err := follow.FindFollow()
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		if exist {
			//删除一条关注记录
			err = follow.DeleteFollow()
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				})
				return
			}
			err = follow.UpdateFollow(false)
			if err != nil {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: err.Error()},
				})
				return
			}
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "Unfollow success",
		})
	}
}

// FollowList
func FollowList(c *gin.Context) {
	fansIdStr := c.Query("user_id")
	fansId, _ := strconv.ParseInt(fansIdStr, 10, 32)
	follow := repository.Follow{
		FansId: fansId,
	}
	followIdList, err := follow.AllFollow()
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: []UserMessage{},
		})
	}
	user := repository.User{
		Id: fansId,
	}
	followList, err := user.GetFollowList(followIdList)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: []UserMessage{},
		})
	}
	var usermessages []UserMessage
	for _, value := range followList {
		usermessage := UserMessage{
			Id:            value.Id,
			Name:          value.Username,
			FollowCount:   value.FollowCount,
			FollowerCount: value.FollowerCount,
			IsFollow:      true,
		}
		usermessages = append(usermessages, usermessage)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "FollowList",
		},
		UserList: usermessages,
	})
}

// FollowerList
func FollowerList(c *gin.Context) {
	followIdStr := c.Query("user_id")
	followId, _ := strconv.ParseInt(followIdStr, 10, 32)
	follow := repository.Follow{
		FollowId: followId,
	}
	fansIdList, err := follow.AllFans()
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: []UserMessage{},
		})
	}
	user := repository.User{
		Id: followId,
	}
	fansList, err := user.GetFollowList(fansIdList)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: []UserMessage{},
		})
	}
	var usermessages []UserMessage
	for _, value := range fansList {
		usermessage := UserMessage{
			Id:            value.Id,
			Name:          value.Username,
			FollowCount:   value.FollowCount,
			FollowerCount: value.FollowerCount,
		}
		tempFollow := repository.Follow{
			FollowId: value.Id,
			FansId:   followId,
		}
		is_follow, err := tempFollow.FindFollow()
		if err != nil {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
				UserList: []UserMessage{},
			})
			return
		}
		usermessage.IsFollow = is_follow
		usermessages = append(usermessages, usermessage)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "FollowList",
		},
		UserList: usermessages,
	})
}
