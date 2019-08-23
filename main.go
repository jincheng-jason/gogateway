package main

import (
	_ "napi-server/docs"
	_ "napi-server/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "sbwagger"
		beego.SetStaticPath("/images", "statics/images")
		beego.SetStaticPath("/css", "statics/css")
		beego.SetStaticPath("/js", "statics/js")
	}
	beego.SetStaticPath("/images", "statics/images")
        beego.SetStaticPath("/css", "statics/css")
        beego.SetStaticPath("/js", "statics/js")
	beego.Run()
}
