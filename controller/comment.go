package controller

import (
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	claims := c.MustGet("claims").(*token.Claims)
	videoIdStr := c.Query("video_id")
	//userIdStr := c.Query("user_id")
	actionType := c.Query("action_type")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 36)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}

	//userId, err := strconv.ParseInt(userIdStr, 10, 36)
	//if err != nil {
	//	c.JSON(http.StatusOK, CommentActionResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: err.Error()},
	//	})
	//}

	comment := repository.Comment{
		VideoId:    videoId,
		UserId:     claims.UserId,
		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	if actionType == "1" {
		commentText := c.Query("comment_text")
		comment.Content = commentText
		if err := comment.CreateComment(); err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Success"},
			Comment: Comment{
				Id: comment.Id,
				User: User{
					Id:            claims.UserId,
					Name:          claims.Username,
					FollowCount:   0,
					FollowerCount: 0,
					IsFollow:      false,
				},
				Content:    comment.Content,
				CreateDate: comment.CreateDate,
			},
		})
	} else if actionType == "2" {
		commentIdStr := c.Query("comment_id")

		commentId, err := strconv.ParseInt(commentIdStr, 10, 36)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		comment.Id = commentId
		err = comment.DeleteComment()
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
		})
	}
}

func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 36)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	comment := repository.Comment{
		VideoId: videoId,
	}
	commentListRepo, err := comment.FindCommentsByVideoId()
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	var comments []Comment

	for _, commentRepo := range commentListRepo {
		comment := Comment{
			Id: commentRepo.Id,
			User: User{
				Id:            commentRepo.UserId,
				Name:          commentRepo.User.Username,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
			},
			Content:    commentRepo.Content,
			CreateDate: commentRepo.CreateDate,
		}
		comments = append(comments, comment)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: comments,
	})
}
