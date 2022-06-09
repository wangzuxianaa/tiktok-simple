package controller

import (
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//
// CommentListResponse
// @Description: 评论列表的响应
//
type CommentListResponse struct {
	service.Response
	CommentList []service.CommentMessage `json:"comment_list,omitempty"`
}

//
// CommentActionResponse
// @Description: 评论操作的响应
//
type CommentActionResponse struct {
	service.Response
	Comment service.CommentMessage `json:"comment,omitempty"`
}

//
// CommentAction
// @Description: 评论操作
// @param c
//
func CommentAction(c *gin.Context) {
	// 获取claims，里面包含用户id和用户名
	claims := c.MustGet("claims").(*token.Claims)
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")

	// videoId string 转 int
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	// 创建评论
	if actionType == "1" {
		// 获取评论的内容
		commentText := c.Query("comment_text")
		comment, err := service.PublishComment(videoId, claims.UserId, commentText, actionType)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: service.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: service.Response{StatusCode: 0, StatusMsg: "Success"},
			Comment: service.CommentMessage{
				Id: comment.Id,
				User: service.UserMessage{
					Id:            claims.UserId,
					Name:          comment.User.Username,
					FollowCount:   comment.User.FollowCount,
					FollowerCount: comment.User.FollowerCount,
					IsFollow:      comment.User.IsFollow,
				},
				Content:    comment.Content,
				CreateDate: comment.CreateDate,
			},
		})
	} else if actionType == "2" { // 删除评论
		commentIdStr := c.Query("comment_id")
		commentId, err := strconv.ParseInt(commentIdStr, 10, 36)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}

		// 删除一条评论信息
		if err = service.DeleteComment(videoId, commentId, actionType); err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: service.Response{StatusCode: 0, StatusMsg: "success"},
		})
	}
}

//
// CommentList
// @Description: 评论列表
// @param c
//
func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")

	var err error
	var videoId int64
	videoId, err = strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	var commentList *[]service.CommentMessage
	// 获取评论列表
	commentList, err = service.GetCommentList(videoId)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    service.Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: *commentList,
	})
}
