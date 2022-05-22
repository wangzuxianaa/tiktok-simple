package main

import (
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
	"os"
)

func Init() error {
	if err := repository.Init(); err != nil {
		return err
	}
	return nil
}
func main() {
	if err := Init(); err != nil {
		os.Exit(-1)
	}
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
