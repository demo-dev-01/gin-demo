package main

import (
	"github.com/godemo-dev/gin-demo/model"
	"github.com/godemo-dev/gin-demo/routers"
)

func main() {
	model.Setup()
	routers.InitRouters()
}
