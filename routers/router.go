package routers

import (
	"image-bg-remover/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/upload", &controllers.MainController{}, "post:UploadImage")

}
