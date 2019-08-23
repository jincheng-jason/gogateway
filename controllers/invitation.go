package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	//	"io/ioutil"
	"net/http"
)

type InvitationController struct {
	beego.Controller
}

type Invresp struct {
	OperCode int      `json:"oper_code"`
	Message  string   `json:"message"`
	ErrCode  string   `json:"err_code"`
	Data     RespData `json:"data"`
}

type RespData struct {
	Praisedusers praisedusers `json:"praisedusers"`
	Replys       replys       `json:"replys"`
	Topic        topic        `json:"topic"`
}

type praisedusers struct {
}

type replys struct {
	ReplyData []reply `json:"data"`
	ReplyPage page    `json:"page"`
}

type topic struct {
	AuthorId int    `json:"authorId"`
	Content  string `json:"content"`
	IsBest   int    `json:"isBest"`
	IsTop    int    `json:"isTop"`
	PostedAt string `json:"postedAt"`
	Title    string `json:"title"`
	User     user   `json:"user"`
}

type reply struct {
	Content  string `json:"content"`
	PostedAt string `json:"postedAt"`
	User     user   `json:"user"`
	Floor    int    `json:"floor"`
}

type user struct {
	Avatar   string `json:"avatar"`
	NickName string `json:"nickName"`
}

type page struct {
	TotalPage   int  `json:"total-page"`
	TotalRow    int  `json:"total-row"`
	PageHasNext bool `json:"page-hasnext"`
	PageNum     int  `json:"page-num"`
	PageRow     int  `json:"page-row"`
}

//url= ?sectionId=3&topicId=123
// @router /share [get]
func (inv *InvitationController) Get() {
	fmt.Println("------------GET SHARED PAGE-----------")
	sectionid, err := inv.GetInt("sectionId")
	if err != nil {
		fmt.Println(err)
		sectionid = 0
	}
	topicid, err := inv.GetInt("topicId")
	if err != nil {
		fmt.Println(err)
		topicid = 0
	} // /v1/group/sections/3/topics/123
	rurl := fmt.Sprintf("http://localhost:%d/v1/group/sections/%d/topics/%d", beego.HttpPort, sectionid, topicid)
	fmt.Println(rurl)
	//获取页面所需数据 /v1/group/sections/3/topics/123
	resp, err := http.Get(rurl)
	if err != nil || resp.StatusCode >= 400 {
		fmt.Println(err)
		inv.TplNames = "404.html"
		return
	}
	defer resp.Body.Close()
	//	rsbyte, err := ioutil.ReadAll(resp.Body)
	result := &Invresp{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&result)
	//	json.Unmarshal(rsbyte, &result)
	inv.Data["topic"] = result.Data.Topic
	inv.Data["replys"] = result.Data.Replys.ReplyData
	inv.TplNames = "share.html"
}

// @router /* [post]
func (inv *InvitationController) Post() {
}

// @router /* [put]
func (inv *InvitationController) Put() {
}

// @router /* [delete]
func (inv *InvitationController) Delete() {
}
