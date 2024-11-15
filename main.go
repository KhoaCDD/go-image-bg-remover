package main

import (
	_ "image-bg-remover/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.Run()
}
