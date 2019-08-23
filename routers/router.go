// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"napi-server/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/news",
			beego.NSInclude(
				&controllers.NewsController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/data",
			beego.NSInclude(
				&controllers.DataController{},
			),
		),
		beego.NSNamespace("/group",
			beego.NSInclude(
				&controllers.GroupController{},
			),
		),
		beego.NSNamespace("/invitation",
			beego.NSInclude(
				&controllers.InvitationController{},
			),
		),
		beego.NSNamespace("/mall",
			beego.NSInclude(
				&controllers.MallController{},
			)),
	)
	beego.AddNamespace(ns)
}
