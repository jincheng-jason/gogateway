package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	yaml "gopkg.in/v2/yaml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	//	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Rinfo *RootInfo

type RootInfo struct {
	GetInfos    []*RequestInfo `yaml:"get"`
	PostInfos   []*RequestInfo `yaml:"post"`
	PutInfos    []*RequestInfo `yaml:"put"`
	DeleteInfos []*RequestInfo `yaml:"delete"`
}

type RequestInfo struct {
	ApiUrl          string              `yaml:"key-url"`
	PathParamsIndex map[string]string   `yaml:"path-params-index"`
	Params          map[string]string   `yaml:"params"`
	ProxyRequest    []*ProxyRequestInfo `yaml:"proxy-request"`
}

type ProxyRequestInfo struct {
	ReturnType    string   `yaml:"return-type"`
	ProxyUrl      string   `yaml:"proxy-url"`
	AccessName    string   `yaml:"access-name"`
	HasPageSize   bool     `yaml:"haspagesize"`
	SubAccessName string   `yaml:"sub-access-name"`
	SubQuery      string   `yaml:"sub-query"`
	SubRefParams  []string `yaml:"sub-ref-params"`
	ProxyNewUrl   string
}

func (p *ProxyRequestInfo) GetProxyUrl(params map[string]string) string {
	p.ProxyNewUrl = p.ProxyUrl
	for k, v := range params {
		splitstr := fmt.Sprintf("{%s}", k)
		p.ProxyNewUrl = strings.Replace(p.ProxyNewUrl, splitstr, v, -1)
	}
	//	fmt.Println("======================")
	//	fmt.Println(p.ProxyNewUrl)
	//	fmt.Println("======================")
	return p.ProxyNewUrl
}

func init() {
	initRuleFromFile("conf/urls.yaml")
}

func initRuleFromFile(path string) {
	if Rinfo == nil {
		Rinfo = &RootInfo{}
	}
	getInfoFromFile(path)
}

//从配置文解析路由规则到RMap
func getInfoFromFile(filepath string) {
	fi, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer fi.Close()
	byteyaml, err := ioutil.ReadAll(fi)
	yaml.Unmarshal(byteyaml, &Rinfo)
	fmt.Println(Rinfo)
	for _, v := range Rinfo.PostInfos {
		fmt.Println(v)
	}
}

func (r *RequestInfo) SetParams(pms map[string]string) {
	for k, _ := range r.Params {
		r.Params[k] = pms[r.PathParamsIndex[k]]
	}
	//	fmt.Println(r.Params)
}

func GetRequestRule(url string, method Method) *RequestInfo {
	var rinfo *RequestInfo
	switch method {
	case GET:
		rinfo = getRequestInfo(url, Rinfo.GetInfos)
	case POST:
		rinfo = getRequestInfo(url, Rinfo.PostInfos)
	case PUT:
		rinfo = getRequestInfo(url, Rinfo.PutInfos)
	case DELETE:
		rinfo = getRequestInfo(url, Rinfo.DeleteInfos)
	}

	return rinfo
}

func getRequestInfo(url string, rp []*RequestInfo) *RequestInfo {
	var rinfo *RequestInfo
	for _, v := range rp {
		//		fmt.Println(v.ApiUrl)
		//		fmt.Println(url)
		matched, err := regexp.MatchString(v.ApiUrl, url)
		//		fmt.Println(matched)
		//		fmt.Println(err)
		if matched && err == nil {
			fmt.Println(v.ApiUrl)
			fmt.Println(url)
			rinfo = v
			break
		}
	}
	return rinfo
}

//-------------------------------------------------method----------------------//
type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

//-------------------------------------------------response--------------------//
type MSResponseInfo struct {
	OperCode int         `json:"oper_code"`
	Message  string      `json:"message"`
	ErrCode  string      `json:"err_code"`
	Data     interface{} `json:"data"`
}

func NewMsResponseInfo() *MSResponseInfo {
	return &MSResponseInfo{
		OperCode: 1,
		Message:  "",
		ErrCode:  "",
	}
}

//-----------------------------------MSproxy---------------------------------//

type MSProxy struct {
	Director      func(*http.Request)
	Transport     http.RoundTripper
	FlushInterval time.Duration
	ErrorLog      *log.Logger
	Msreq         *http.Request  //微服务请求
	Msres         *http.Response //微服务应答
	MsRawQuery    string
}

//创建一个微服务代理
func NewMsProxy(target *url.URL) *MSProxy {
	targetQuery := target.RawQuery // 请求参数
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.RequestURI = target.String()
		//		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		fmt.Println(req.URL.RawQuery)
	}
	return &MSProxy{Director: director}
}

//通过代理处理请求
func (p *MSProxy) HandleProxyRequest(req *http.Request, result *MSResponseInfo, reqinfo *ProxyRequestInfo) error {
	if p.Transport == nil {
		p.Transport = http.DefaultTransport
	}

	p.Msreq = new(http.Request)
	p.MsRawQuery = req.URL.RawQuery
	req.URL.RawQuery = ""
	*(p.Msreq) = *req //复制请求

	p.Director(p.Msreq) //导向处理
	p.Msreq.Proto = "HTTP/1.1"
	p.Msreq.ProtoMajor = 1
	p.Msreq.ProtoMinor = 1
	p.Msreq.Close = false

	copiedHeaders := false
	for _, h := range msHopHeaders { //如果未拷贝头，先拷贝头，再删除逐跳头
		if p.Msreq.Header.Get(h) != "" {
			if !copiedHeaders {
				p.Msreq.Header = make(http.Header)
				copyHeader(p.Msreq.Header, req.Header)
				copiedHeaders = true
			}
			p.Msreq.Header.Del(h)
		}
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := p.Msreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		p.Msreq.Header.Set("X-Forwarded-For", clientIP)
	}
	req.ParseForm()
	pageNum := req.Form.Get("page-num")
	pageRow := req.Form.Get("page-row")

	//	fmt.Println(req.Form)
	//	fmt.Println("-------===-------")
	//	fmt.Println(pageNum)
	//	fmt.Println(pageRow)
	//	fmt.Println("-------===-------")

	if pageNum != "" {
		p.Msreq.Header.Set("X-Page-Num", pageNum)
	}
	if pageRow != "" {
		p.Msreq.Header.Set("X-Page-row", pageRow)
	}
	fmt.Println(p.Msreq.URL)
	res, err := p.Transport.RoundTrip(p.Msreq) //一次请求
	if err != nil {
		p.logf("http: proxy error: %v", err)
		result.OperCode = 0
		result.Message = ErrMap[res.StatusCode]
		return err
	}
	p.Msres = res
	defer p.Msres.Body.Close()
	//处理一次状态码
	if p.Msres.StatusCode >= 400 && p.Msres.StatusCode < 500 {
		result.OperCode = 0
		msg, err := url.QueryUnescape(p.Msres.Header.Get("X-Err-Message"))
		result.Message = msg
		//		log.Println(url.QueryUnescape(result.Message))
		if result.Message == "" || err != nil {
			result.Message = ErrMap[p.Msres.StatusCode]
		}
		return errors.New(ErrMap[p.Msres.StatusCode])
	}

	var bodyrs interface{}
	decoder := json.NewDecoder(p.Msres.Body)
	decoder.UseNumber()
	err = decoder.Decode(&bodyrs)
	//	fmt.Println(bodyrs)
	if err != nil {
		result.OperCode = 0
		result.Message = "json序列化错误"
		return err
	}
	result.OperCode = 1
	//	fmt.Println(reflect.TypeOf(result.Data))

	if reqinfo.AccessName != "" && result.Data == nil {
		result.Data = make(map[string]interface{})
	}
	if reqinfo.AccessName != "" {
		if reqinfo.HasPageSize {
			tempmap := make(map[string]interface{})
			tempmap["data"] = bodyrs
			tempmap["page"] = p.handlerPageInfo()
			p.handleSubQuery(bodyrs.([]interface{}), reqinfo)
			result.Data.(map[string]interface{})[reqinfo.AccessName] = tempmap
		} else {
			result.Data.(map[string]interface{})[reqinfo.AccessName] = bodyrs
		}

	} else {
		result.Data = bodyrs
	}
	return nil
}

func (p *MSProxy) handleSubQuery(result []interface{}, reqinfo *ProxyRequestInfo) error {
	//	temparr := make([]interface{})
	if reqinfo.SubQuery != "" {
		for _, v := range result {
			realstr := RealGroupProxy.GetProxyHostUrl() + reqinfo.SubQuery
			for _, vv := range reqinfo.SubRefParams {
				splitstr := fmt.Sprintf("{%s}", vv)
				realstr = strings.Replace(realstr, splitstr, string(v.(map[string]interface{})[vv].(json.Number)), -1)
				//				realstr = RealGroupProxy.GetProxyHostUrl() + reqinfo.SubQuery
				//
				//fmt.Println(realstr)
				//				fmt.Println(v.(map[string]interface{})[vv])
				//
			}
			//			tempproxyurl := url.Parse(realstr)
			tempresp, temperr := http.Get(realstr)
			fmt.Println("Get : " + realstr)
			if temperr != nil {
				return nil
			}
			defer tempresp.Body.Close()
			var bodyrs interface{}
			decoder := json.NewDecoder(tempresp.Body)
			decoder.UseNumber()
			err := decoder.Decode(&bodyrs)
			if reqinfo.SubAccessName == "" {
				fmt.Println(err)
				return nil
			}
			v.(map[string]interface{})[reqinfo.SubAccessName] = bodyrs
		}
	}
	return nil
}

func (p *MSProxy) handleErrInfo() map[string]interface{} {
	info := make(map[string]interface{})
	for _, h := range msPageInfoHeaders {
		switch h {
		case "X-Err-Message":
			info["message"] = p.Msres.Header.Get(h)
		case "X-Err-Code":
			info["err_code"] = p.Msres.Header.Get(h)
		}
	}
	return info
}

//当前请求的页码信息
func (p *MSProxy) handlerPageInfo() *PageInfo {
	page := new(PageInfo)
	existPageInfo := true
	//	fmt.Println(p.Msres.Header)
	for _, h := range msPageInfoHeaders {
		switch h {
		case "X-Total-Page":
			tp := p.Msres.Header.Get(h)
			if tp == "" {
				existPageInfo = false
			} else {
				//				fmt.Println(tp)
				tpint, err := strconv.Atoi(tp)
				if err == nil {
					page.TotalPage = tpint
				} else {
					page.TotalPage = 0
				}
			}
		case "X-Total-Row":
			tr := p.Msres.Header.Get(h)
			if tr == "" {
				existPageInfo = false
			} else {
				//				fmt.Println(tr)
				trint, err := strconv.Atoi(tr)
				if err == nil {
					page.TotalRow = trint
					continue
				} else {
					page.TotalRow = 0
				}
			}
		case "X-Page-Row":
			pr := p.Msres.Header.Get(h)
			if pr == "" {
				existPageInfo = false
			} else {

				pageint, err := strconv.Atoi(pr)
				if err == nil {
					page.PageRow = pageint
				} else {
					page.PageRow = 0
				}
			}
		case "X-Page-Num":
			pn := p.Msres.Header.Get(h)
			if pn == "" {
				existPageInfo = false
			} else {
				pnint, err := strconv.Atoi(pn)
				if err == nil {
					page.PageNum = pnint
				} else {
					page.PageNum = 0
				}
			}

		case "X-Page-Hasnext":
			ph := p.Msres.Header.Get(h)
			if ph == "" {
				existPageInfo = false
			} else {
				phbool := false
				if ph == "true" {
					phbool = true
				}
				page.PageHasNext = phbool

			}

		default:
		}
	}
	//	fmt.Println(p.Msres.Header)
	//	fmt.Println(page)
	if existPageInfo {
		return page
	}
	return nil
}

func (p *MSProxy) HandlerResponse() {
	//检测服务器状态码
	if p.Msres.StatusCode >= 200 && p.Msres.StatusCode < 400 { //正常服务

	} else if p.Msres.StatusCode >= 400 && p.Msres.StatusCode < 500 { //参数一场

	} else if p.Msres.StatusCode >= 500 { //服务器异常

	}
	//	temprs := make(map[string]interface{})
	//	pinfo := p.handlerPageInfo()
	//	if pinfo != nil {
	//		temprs["page"] = pinfo
	//	}

	//bodybyte, err := ioutil.ReadAll(p.Msres.Body)
	//	if err != nil {
	//		return nil
	//	}
	//	fmt.Println(string(bodybyte))
	//	var bodyrs interface{}
	//	decoder := json.NewDecoder(p.Msres.Body)
	//	decoder.UseNumber()
	//	err := decoder.Decode(&bodyrs)
	//	err = json.Unmarshal(bodybyte, &bodyrs)
	//	if err == nil {
	//		temprs[accessName] = bodyrs
	//		return temprs
	//	}
	//	fmt.Println(err)

}

//func handleServerError(errcode int) string{
//	errmsg := ""
//switch errcode{
//	case 200:
//	errmsg ="当GET请求成功完成，DELETE或者PATCH请求同步完成。"
//case	201:
//errmsg ="创建数据成功，同步方式成功完成POST请求。",
//	case 202:
//	errmsg ="POST，DELETE或者PATCH请求提交成功，稍后将异步的进行处理。",
//	case 204:
//	errmsg ="无内容，资源有空表示。",
//	case 206:
//	errmsg ="GET请求成功完成，但只返回了部分数据。(目前可不实现该约定)。"
//	case 301:
//	errmsg ="Moved Permanently，资源的URI已被更新(转移)"
//	case 303:
//	errmsg ="See Other,其他（如，负载均衡）"
//	case 304:
//	errmsg ="Not Modified, 资源未更改（缓存）"
//	case 400:
//	errmsg ="Bad Request: 请求格式错误，请求参数有误，不被支持。"
//	case 401:
//	errmsg ="Unauthorized: 请求失败，因为用户没有进行认证。"
//	case 403:
//	errmsg ="Forbidden: 请求失败，因为用户被认定没有访问特定资源的权限。"
//	case 404:
//	errmsg ="Not Found: 未找到该资源地址，服务器不能接受该url的请求。"
//	case 422:
//	errmsg ="Unprocessable Entity: 你的请求服务器可以理解，但是其中包含了不合法的参数。"
//	case 429:
//	errmsg ="Too Many Requests: 请求频率超配，稍后再试。"
//	case 500:
//	errmsg ="Internal Server Error: 服务器出错了，检查网站的状态，或者报告问题"
//	case 503:
//	errmsg ="Service Unavailable: 服务端当前无法处理请求"
//}
//return errmsg
//}

func (p *MSProxy) logf(format string, args ...interface{}) {
	if p.ErrorLog != nil {
		p.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

//自定义分页信息头
var msPageInfoHeaders = []string{
	"X-Total-Page",
	"X-Total-Row",
	"X-Page-Row",
	"X-Page-Num",
	"X-Page-Hasnext",
}

type PageInfo struct {
	TotalPage   int  `json:"total-page"`
	TotalRow    int  `json:"total-row"`
	PageHasNext bool `json:"page-hasnext"`
	PageNum     int  `json:"page-num"`
	PageRow     int  `json:"page-row"`
}

// 逐跳头. 向后端发送的时候删除掉.
var msHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

//微服务系统状态码
var ErrMap = map[int]string{
	200: "当GET请求成功完成，DELETE或者PATCH请求同步完成。",
	201: "创建数据成功，同步方式成功完成POST请求。",
	202: "POST，DELETE或者PATCH请求提交成功，稍后将异步的进行处理。",
	204: "无内容，资源有空表示。",
	206: "GET请求成功完成，但只返回了部分数据。(目前可不实现该约定)。",
	301: "Moved Permanently，资源的URI已被更新(转移)",
	303: "See Other,其他（如，负载均衡）",
	304: "Not Modified, 资源未更改（缓存）",
	400: "Bad Request: 请求格式错误，请求参数有误，不被支持。",
	401: "Unauthorized: 请求失败，因为用户没有进行认证。",
	403: "Forbidden: 请求失败，因为用户被认定没有访问特定资源的权限。",
	404: "Not Found: 未找到该资源地址，服务器不能接受该url的请求。",
	422: "Unprocessable Entity: 你的请求服务器可以理解，但是其中包含了不合法的参数。",
	429: "Too Many Requests: 请求频率超配，稍后再试。",
	500: "Internal Server Error: 服务器出错了，检查网站的状态，或者报告问题",
	503: "Service Unavailable: 服务端当前无法处理请求",
}
