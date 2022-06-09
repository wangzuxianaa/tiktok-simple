package main

import (
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/pkg/utils"
	"github.com/gin-gonic/gin"
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

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
