package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["napi-server/controllers:DataController"] = append(beego.GlobalControllerRouter["napi-server/controllers:DataController"],
		beego.ControllerComments{
			"Get",
			`/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:DataController"] = append(beego.GlobalControllerRouter["napi-server/controllers:DataController"],
		beego.ControllerComments{
			"Post",
			`/*`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:DataController"] = append(beego.GlobalControllerRouter["napi-server/controllers:DataController"],
		beego.ControllerComments{
			"Put",
			`/*`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:DataController"] = append(beego.GlobalControllerRouter["napi-server/controllers:DataController"],
		beego.ControllerComments{
			"Delete",
			`/*`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:GroupController"] = append(beego.GlobalControllerRouter["napi-server/controllers:GroupController"],
		beego.ControllerComments{
			"Get",
			`/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:GroupController"] = append(beego.GlobalControllerRouter["napi-server/controllers:GroupController"],
		beego.ControllerComments{
			"Post",
			`/*`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:GroupController"] = append(beego.GlobalControllerRouter["napi-server/controllers:GroupController"],
		beego.ControllerComments{
			"Put",
			`/*`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:GroupController"] = append(beego.GlobalControllerRouter["napi-server/controllers:GroupController"],
		beego.ControllerComments{
			"Delete",
			`/*`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:InvitationController"] = append(beego.GlobalControllerRouter["napi-server/controllers:InvitationController"],
		beego.ControllerComments{
			"Get",
			`/share`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:InvitationController"] = append(beego.GlobalControllerRouter["napi-server/controllers:InvitationController"],
		beego.ControllerComments{
			"Post",
			`/*`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:InvitationController"] = append(beego.GlobalControllerRouter["napi-server/controllers:InvitationController"],
		beego.ControllerComments{
			"Put",
			`/*`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:InvitationController"] = append(beego.GlobalControllerRouter["napi-server/controllers:InvitationController"],
		beego.ControllerComments{
			"Delete",
			`/*`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:NewsController"] = append(beego.GlobalControllerRouter["napi-server/controllers:NewsController"],
		beego.ControllerComments{
			"Get",
			`/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:NewsController"] = append(beego.GlobalControllerRouter["napi-server/controllers:NewsController"],
		beego.ControllerComments{
			"Post",
			`/*`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:NewsController"] = append(beego.GlobalControllerRouter["napi-server/controllers:NewsController"],
		beego.ControllerComments{
			"Put",
			`/*`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:NewsController"] = append(beego.GlobalControllerRouter["napi-server/controllers:NewsController"],
		beego.ControllerComments{
			"Delete",
			`/*`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:UserController"] = append(beego.GlobalControllerRouter["napi-server/controllers:UserController"],
		beego.ControllerComments{
			"Get",
			`/*`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:UserController"] = append(beego.GlobalControllerRouter["napi-server/controllers:UserController"],
		beego.ControllerComments{
			"Post",
			`/*`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:UserController"] = append(beego.GlobalControllerRouter["napi-server/controllers:UserController"],
		beego.ControllerComments{
			"Put",
			`/*`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["napi-server/controllers:UserController"] = append(beego.GlobalControllerRouter["napi-server/controllers:UserController"],
		beego.ControllerComments{
			"Delete",
			`/*`,
			[]string{"delete"},
			nil})

}
