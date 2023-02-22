# tiktok-simple

## 上手指南

### 配置要求

- [Go v1.17](https://go.dev/)
- [Mysql v8.0.28](https://www.mysql.com/)
- [redis v5.0.14](https://redis.io/)
- [ffmpeg](https://ffmpeg.org/)

### **运行步骤**

项目已经部署到服务器上
客户端连接ip地址：http://49.235.232.7:8989/

## 文件说明

```shell
filetree
├─cache（缓存）
│  ├─count.go（评论总数，点赞总数相关）
│  └─favourite.go（点赞逻辑相关）
├─conf
│  ├─conf.go（viper配置）
│  └─conf.yaml（配置文件）
├─controller(视图层实现)
├─model（mysql数据库模型）
├─pkg（扩展包）
│  ├─middleware(中间件)
│    └─jwtAuth.go
│  ├─token(token鉴权相关)
│    └─jwt.go(jwt)
│  └─utils(公共包)
│    ├─cron_task.go(定时任务)
│    ├─ffmpeg.go(截取视频为帧)
│    └─makesha1.go(加密)
├─public（本地存放视频）
│  ├─cover(封面地址)
│  └─video(视频地址)
├─service（服务层实现）
│  ├─comment_list.go(评论列表)
│  ├─common.go(公共)
│  ├─favourite_list.go(喜好列表)
│  ├─feed_videoList.go(视频流)
│  ├─follow_list.go(关注列表)
│  ├─follower_list.go(粉丝列表)
│  ├─post_comment.go(评论操作)
│  ├─post_follow_service.go(关注操作)
│  ├─publish_list.go(发布列表)
│  └─user_service.go(用户登陆注册相关)
```
