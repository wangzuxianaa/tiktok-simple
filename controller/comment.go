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
	videoIdStr := c.Query("video_id") //视频id
	//userIdStr := c.Query("user_id")
	actionType := c.Query("action_type") //评论/删除评论
	//评论列表应该是与视频绑定的

	videoId, err := strconv.ParseInt(videoIdStr, 10, 36) //将视频Id转换为int
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}

	comment := repository.Comment{ //构建一个评论列表
		VideoId:    videoId,
		UserId:     claims.UserId,
		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	if actionType == "1" { //发布评论
		commentText := c.Query("comment_text") //评论内容
		comment.Content = commentText
		if err := comment.CreateComment(); err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{ //返回评论信息
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
	} else if actionType == "2" { //删除评论
		commentIdStr := c.Query("comment_id")                    //获取要删除的评论的id信息
		commentId, err := strconv.ParseInt(commentIdStr, 10, 36) //要删除的评论对应的id
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		comment.Id = commentId
		err = comment.DeleteComment() //通过要删除的评论对应的Id将评论删除掉
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

func CommentList(c *gin.Context) { //显示评论列表
	videoIdStr := c.Query("video_id") //获取视频id

	videoId, err := strconv.ParseInt(videoIdStr, 10, 36)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	comment := repository.Comment{ //建立一个评论 并设置这是哪个视频的评论
		VideoId: videoId,
	}
	commentListRepo, err := comment.FindCommentsByVideoId() //通过视频的id找到对应的评论
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
