package routers

import (
	"go-tile-cache/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.TileController{})
    beego.Router("/test", &controllers.TestController{})
}
