package service

import (
	"errors"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"github.com/wangzuxianaa/tiktok-simple/pkg/token"
	"github.com/wangzuxianaa/tiktok-simple/pkg/utils"
)

//
// Register
// @Description: 注册功能，需要用户名和密码
// @param username
// @param password
// @return int64
// @return string
// @return error
//
func Register(username, password string) (int64, string, error) {
	// 根据username来查询用户信息
	_, ok, err := model.NewUserDaoInstance().FindUserByName(username)
	if err != nil {
		return UserIdNotFound, NoGenerateToken, err
	}
	// 查询到用户
	if ok == true {
		return UserIdNotFound, NoGenerateToken, errors.New("user already exists")
	}

	// 加密
	pwd := utils.MakeSha1(password)

	user := model.User{
		Username: username,
		Password: pwd,
	}
	// 创建用户
	if err := model.NewUserDaoInstance().CreateUser(&user); err != nil {
		return -1, "", err
	}
	// 生成token
	generateToken, err := token.GenerateToken(user.Id, user.Username)
	if err != nil {
		return user.Id, NoGenerateToken, errors.New("generating token fails")
	}
	return user.Id, generateToken, nil
}

//
// Login
// @Description: 登陆功能，需要账户名和密码
// @param username
// @param password
// @return int64
// @return string
// @return error
//
func Login(username, password string) (int64, string, error) {
	user, ok, err := model.NewUserDaoInstance().FindUserByName(username)
	if err != nil {
		return UserIdNotFound, NoGenerateToken, err
	}
	// 没有找到用户
	if ok == false {
		return UserIdNotFound, NoGenerateToken, errors.New("user does not exist")
	}

	// 密码不正确
	if user.Password != utils.MakeSha1(password) {
		return user.Id, NoGenerateToken, errors.New("password is not correct")
	}

	// 生成token
	generateToken, err := token.GenerateToken(user.Id, username)
	if err != nil {
		return user.Id, NoGenerateToken, errors.New("generating token fails")
	}
	return user.Id, generateToken, err
}

//
// GetUserInfo
// @Description: 根据Id号获取用户信息
// @param Id
// @return *UserMessage
// @return error
//
func GetUserInfo(userId int64, userIdFromToken int64) (*UserMessage, error) {
	// 根据用户id查找用户
	user, err := model.NewUserDaoInstance().FindUserById(userId)
	if err != nil {
		return nil, err
	}
	isFollow, err := model.NewFollowDaoInstance().FindFollow(userId, userIdFromToken)
	if err != nil {
		return nil, err
	}
	var userMessage = UserMessage{
		Id:            userId,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	return &userMessage, err
}
