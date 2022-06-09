package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static/video", "./public/video")
	r.Static("/static/cover", "./public/cover")

	apiRouter := r.Group("/douyin")

	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	apiRouter.Use(middleware.TokenChecker())
	{
		apiRouter.GET("/user/", controller.UserInfo)
		apiRouter.GET("/publish/list/", controller.PublishList)
		apiRouter.POST("/publish/action/", controller.Publish)

		apiRouter.POST("/favorite/action/", controller.FavouriteAction)
		apiRouter.GET("/favorite/list/", controller.FavouriteList)
		apiRouter.POST("/comment/action/", controller.CommentAction)
		apiRouter.GET("/comment/list/", controller.CommentList)

		apiRouter.POST("/relation/action/", controller.RelationAction)
		apiRouter.GET("/relation/follow/list/", controller.FollowList)
		apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	}
}
