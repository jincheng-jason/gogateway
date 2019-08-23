package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"napi-server/utils"
	"net/url"
)

// Operations about Users
type NewsController struct {
	beego.Controller
}

// @router /* [get]
func (u *NewsController) Get() {
	//请求路径
	u.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	reqstr := u.Ctx.Input.Param(":splat")
	//	cachekey := utils.RealNewsProxy.GetProxyAddress() + reqstr + "?" + u.Ctx.Input.Request.URL.RawQuery
	proxyurl, _ := url.Parse(utils.RealNewsProxy.GetProxyAddress() + reqstr)
	fmt.Println(proxyurl)
	ok := false
	//	bodybytes, ok := utils.CacheGet(cachekey)
	if !ok {
		fmt.Println("Cache missed,read from the service.")
		proxy := utils.NewSingleHostReverseProxy(proxyurl)
		proxy.HandleRequest(u.Ctx.Request)

		if err := proxy.Request(); err != nil {
			u.Data["json"] = utils.JsonError
			u.ServeJson()
			return
		}
		bytebody, err := ioutil.ReadAll(proxy.Outres.Body)
		defer proxy.Outres.Body.Close()
		if err != nil {
			u.Data["json"] = utils.JsonError
			u.ServeJson()
			return
		}
		//		utils.CacheSet(cachekey, string(bytebody))
		proxy.HandleResponse(u.Ctx.ResponseWriter, bytebody)
		return
	}
	//	u.Ctx.Output.Body(bodybytes)

}

// @router /* [post]
func (u *NewsController) Post() {
	//请求路径
	reqstr := u.Ctx.Input.Param(":splat")
	proxyurl, _ := url.Parse(utils.RealNewsProxy.GetProxyAddress() + reqstr)
	proxy := utils.NewSingleHostReverseProxy(proxyurl)
	u.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	proxy.ServeHTTP(u.Ctx.ResponseWriter, u.Ctx.Request)
}

// @router /* [put]
func (u *NewsController) Put() {
	//请求路径
	reqstr := u.Ctx.Input.Param(":splat")
	proxyurl, _ := url.Parse(utils.RealNewsProxy.GetProxyAddress() + reqstr)
	proxy := utils.NewSingleHostReverseProxy(proxyurl)
	u.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	proxy.ServeHTTP(u.Ctx.ResponseWriter, u.Ctx.Request)
}

// @router /* [delete]
func (u *NewsController) Delete() {
	//请求路径
	reqstr := u.Ctx.Input.Param(":splat")
	proxyurl, _ := url.Parse(utils.RealNewsProxy.GetProxyAddress() + reqstr)
	proxy := utils.NewSingleHostReverseProxy(proxyurl)
	u.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	proxy.ServeHTTP(u.Ctx.ResponseWriter, u.Ctx.Request)
}
