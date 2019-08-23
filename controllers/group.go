package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"napi-server/utils"
	"net/url"
)

// Operations about Users
type GroupController struct {
	beego.Controller
}

// @router /* [get]
func (g *GroupController) Get() {
	fmt.Println("------------METHOD GET-----------")
	g.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	result := utils.NewMsResponseInfo()
	rule := utils.GetRequestRule(g.Ctx.Input.Url(), utils.GET)
	fmt.Println(g.Ctx.Input.Url())
	if rule != nil {
		rule.SetParams(g.Ctx.Input.Params)
		g.proxyRequest(rule, result)
	} else {
		result.OperCode = 0
		result.Message = "请求的资源不存在"
	}
	respbyte, _ := json.Marshal(result)
	g.Ctx.Output.Body(respbyte)
}

// @router /* [post]
func (g *GroupController) Post() {
	fmt.Println("------------METHOD POST-----------")
	g.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	result := utils.NewMsResponseInfo()
	rule := utils.GetRequestRule(g.Ctx.Input.Url(), utils.POST)
	if rule != nil {
		rule.SetParams(g.Ctx.Input.Params)
		g.proxyRequest(rule, result)

	} else {
		result.OperCode = 0
		result.Message = "请求的资源不存在"
	}
	respbyte, _ := json.Marshal(result)
	g.Ctx.Output.Body(respbyte)
}

// @router /* [put]
func (g *GroupController) Put() {
	fmt.Println("------------METHOD PUT-----------")
	g.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	result := utils.NewMsResponseInfo()
	rule := utils.GetRequestRule(g.Ctx.Input.Url(), utils.PUT)
	if rule != nil {
		rule.SetParams(g.Ctx.Input.Params)
		g.proxyRequest(rule, result)
	} else {
		result.OperCode = 0
		result.Message = "请求的资源不存在"
	}
	respbyte, _ := json.Marshal(result)
	g.Ctx.Output.Body(respbyte)
}

// @router /* [delete]
func (g *GroupController) Delete() {
	fmt.Println("------------METHOD DELETE-----------")
	g.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	result := utils.NewMsResponseInfo()
	rule := utils.GetRequestRule(g.Ctx.Input.Url(), utils.DELETE)
	if rule != nil {
		rule.SetParams(g.Ctx.Input.Params)
		g.proxyRequest(rule, result)
	} else {
		result.OperCode = 0
		result.Message = "请求的资源不存在"
	}
	respbyte, _ := json.Marshal(result)
	g.Ctx.Output.Body(respbyte)
}

//----------------------------------Private Method-----------------------------//
//代理请求POST,PUT,DELETE
func (g *GroupController) proxyRequest(rule *utils.RequestInfo, res *utils.MSResponseInfo) error {
	//	if len(rule.ProxyRequest) > 1 {
	//		res.Data = make(map[string]interface{})
	//	}
	rawQuery := g.Ctx.Input.Request.URL.RawQuery
	for _, v := range rule.ProxyRequest {
		//真实请求地址
		realstr := utils.RealGroupProxy.GetProxyHostUrl() + v.GetProxyUrl(rule.Params)
		if rawQuery != "" {
			realstr += "?" + rawQuery
		}
		fmt.Println("realstr :" + realstr)
		proxyurl, _ := url.Parse(realstr)

		msproxy := utils.NewMsProxy(proxyurl)
		msproxy.HandleProxyRequest(g.Ctx.Request, res, v)
		//		fmt.Println(proxyurl)
		//		result := msproxy.HandlerResponseWithOutAccessName()
		//		fmt.Println(result)
	}
	return nil
}
