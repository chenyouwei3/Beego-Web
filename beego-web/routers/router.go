package routers

import (
	"beego-web/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/account", &controllers.AccountController{}, "post:Post")

}
