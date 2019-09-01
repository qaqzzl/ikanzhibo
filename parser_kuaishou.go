package main

import (
	"github.com/antchfx/htmlquery"
	"ikanzhibo/db"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//huya解析方法
func (spider *Spider) kuaiShouParser(p *Parser)  {
	switch p.Queue.QueueType {
	case "live_info":
		spider.kuaiShouLiveInfo(p)
	case "live_list":
		spider.kuaiShouLiveList(p)
	}
}

type userLive struct {
	Typename     string `json:"__typename"`
	//BannedStatus struct {
	//	Generated bool   `json:"generated"`
	//	ID        string `json:"id"`
	//	Type      string `json:"type"`
	//	Typename  string `json:"typename"`
	//} `json:"bannedStatus"`
	CityName      string `json:"cityName"`
	Constellation string `json:"constellation"`
	//CountsInfo    struct {
	//	Generated bool   `json:"generated"`
	//	ID        string `json:"id"`
	//	Type      string `json:"type"`
	//	Typename  string `json:"typename"`
	//} `json:"countsInfo"`
	Description    string      `json:"description"`
	Eid            string      `json:"eid"`
	Feeds          interface{} `json:"feeds"`
	ID             string      `json:"id"`
	IsNew          bool        `json:"isNew"`
	KwaiID         string      `json:"kwaiId"`
	Living         bool        `json:"living"`
	Name           string      `json:"name"`
	PrincipalID    string      `json:"principalId"`
	Privacy        bool        `json:"privacy"`
	Profile        string      `json:"profile"`
	Sex            string      `json:"sex"`
	UserID         string      `json:"userId"`
	//VerifiedStatus struct {
	//	Generated bool   `json:"generated"`
	//	ID        string `json:"id"`
	//	Type      string `json:"type"`
	//	Typename  string `json:"typename"`
	//} `json:"verifiedStatus"`
	//WatchingCount interface{} `json:"watchingCount"`
}

//直播间解析方法
func (spider *Spider) kuaiShouLiveInfo(p *Parser) {
	Live := db.TableLive{}

	doc, err := htmlquery.Parse(strings.NewReader(string(p.Body)))
	if err != nil {
		log.Println("htmlquery ERR:" + err.Error())
		return
	}
	//.Live_is_online - 判断是在播
	Live_is_online := htmlquery.FindOne(doc, "//div[@class='live-card']")
	if Live_is_online == nil {
		Live.Live_is_online = "no"
	} else {
		Live.Live_is_online = "yes"
	}

	//.Live_uri #
	Live.Live_uri = p.Queue.Request.Url

	//.Live_platform #
	Live.Live_platform = "kuaishou"

	//.Live_title #
	Live_title := htmlquery.FindOne(doc, "//a[@class='router-link-exact-active router-link-active live-card-following-info-title']")
	Live.Live_title = htmlquery.SelectAttr(Live_title, "title")

	//.Live_anchortv_name #
	Live_anchortv_name := htmlquery.FindOne(doc, "//p[@class='user-info-name']")
	if Live_anchortv_name == nil {
		log.Println("昵称查找失败 \n" + p.Queue.Request.Url + "\n" +string(p.Body))
		return
	}
	Live.Live_anchortv_name = htmlquery.InnerText(Live_anchortv_name)
	// 去除空格
	Live.Live_anchortv_name = strings.Replace(Live.Live_anchortv_name, " ", "", -1)
	// 去除换行符
	Live.Live_anchortv_name = strings.Replace(Live.Live_anchortv_name, "\n", "", -1)
	// 去除字符 (举报)
	Live.Live_anchortv_name = strings.Replace(Live.Live_anchortv_name, "举报", "", -1)

	//.Live_anchortv_photo #
	Live_anchortv_photo := htmlquery.FindOne(doc, "//div[@class='avatar user-info-avatar']/img")
	if Live_anchortv_photo == nil {
		log.Println("头像查找失败\n" + p.Queue.Request.Url)
		return
	}
	Live.Live_anchortv_name = htmlquery.SelectAttr(Live_anchortv_photo, "src")

	//.Live_cover #
	Live.Live_cover = ""

	//.Live_play #
	Live.Live_play = htmlquery.SelectAttr(Live_title, "href")

	//.Live_class #
	Live_class := htmlquery.FindOne(doc, "//span[@class='game-name']")
	if Live_class != nil {
		Live.Live_class = htmlquery.InnerText(Live_class)
	}

	//.Live_online_user #

	//.Live_follow # 被加密了,先不搞
	Live.Live_class = "0"

	//.Live_tag

	//.Live_introduction
	Live_introduction := htmlquery.FindOne(doc, "//p[@class='user-info-description']")
	if Live_introduction != nil {
		Live.Live_introduction = htmlquery.InnerText(Live_introduction)
	}

	//.live_anchortv_sex #


	//.Live_type_id

	//.Live_type_name


	Live.Created_at = strconv.FormatInt(time.Now().Unix(),10)
	Live.Updated_at = strconv.FormatInt(time.Now().Unix(),10)

	spider.ChanWriteInfo <- &p.Queue
	return
}

//直播发现
func (spider *Spider) kuaiShouLiveList(p *Parser) {
	doc, err := htmlquery.Parse(strings.NewReader(string(p.Body)))
	if err != nil {
		log.Println("htmlquery ERR:" + err.Error())
		return
	}

	//分类直播列表页
	regexps := regexp.MustCompile(`<a href="(/cate/[/0-9a-zA-Z]+)" class="category-card-preview"`)
	t := regexps.FindAllSubmatch(p.Body, -1)
	for i:=0; i<len(t); i++ {
		spider.ChanProduceList <- &db.Queue{
			Request:     db.Request{
				Url: "https://live.kuaishou.com"+string(t[i][1]),
			},
			LiveData:    db.TableLive{
				Live_platform: "kuaishou",
			},
			QueueType:     "live_list",
		}
	}

	//直播info详情
	live_info_url := htmlquery.Find(doc, "//a[@class='user-info']")
	for _, n := range live_info_url {
		url := htmlquery.SelectAttr(n, "href")
		spider.ChanProduceList <- &db.Queue{
			Request:     db.Request{
				Url: "https://live.kuaishou.com"+url,
			},
			LiveData:    db.TableLive{
				Live_platform: "kuaishou",
			},
			QueueType:     "live_info",
		}
	}
	//regexps = regexp.MustCompile(`<a href="(/profile/[0-9a-zA-Z]+)" title="[\S]+" target="_blank" class="user-info"`)
	//t = regexps.FindAllSubmatch(p.Body, -1)
	//for i:=0; i<len(t); i++ {
	//	spider.ChanProduceList <- &db.Queue{
	//		Platform: p.Queue.Platform,
	//		Uri:      "https://live.kuaishou.com"+string(t[i][1]),
	//		Type:     "live_info",
	//	}
	//}

	//下一页
	var page int
	regexps = regexp.MustCompile(`<li class="pl-pagination-list-item" data-v-[a-zA-Z0-9]+>[\s\n]+([0-9]+)[\s\n]+</li>`)
	t = regexps.FindAllSubmatch(p.Body, -1)

	regexps = regexp.MustCompile(`(/cate/[/0-9a-zA-Z]+)`)
	uri := regexps.FindSubmatch([]byte(p.Queue.Request.Url))
	for i:=0; i<len(t); i++ {
		page, _ = strconv.Atoi(string(t[i][1]))
	}
	if page != 0 {
		for i:=1; i<page+1; i++ {
			spider.ChanProduceList <- &db.Queue{
				Request:     db.Request{
					Url: "https://live.kuaishou.com"+string(uri[1])+"/?page="+strconv.Itoa(i),
				},
				LiveData:    db.TableLive{
					Live_platform: "kuaishou",
				},
				QueueType:     "live_list",
			}
		}
	}

}