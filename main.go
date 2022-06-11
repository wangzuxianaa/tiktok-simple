package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wangzuxianaa/tiktok-simple/conf"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"github.com/wangzuxianaa/tiktok-simple/pkg/utils"
	"log"
	"os"
)

func Init() error {
	var err error
	// 读取配置文件
	if err = conf.Config(); err != nil {
		log.Print(err)
		return err
	}
	// mysql数据库初始化
	if err = model.MysqlInit(); err != nil {
		log.Print(err)
		return err
	}
	// redis初始化
	if err = model.RedisInit(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func main() {
	if err := Init(); err != nil {
		os.Exit(-1)
	}
	// 定期任务1h一次
	go utils.ExecuteCron()
	r := gin.Default()

	initRouter(r)

	err := r.Run()
	if err != nil {
		os.Exit(-1)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
