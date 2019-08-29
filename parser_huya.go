package main

import (
	"encoding/json"
	"fmt"

	//uuid "github.com/satori/go.uuid"
	"ikanzhibo/db"
	"ikanzhibo/db/mysql"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//huya直播间的json数据
type hyPlayerConfig struct {
	WEBYYFROM string `json:"WEBYYFROM"`
	WEBYYHOST string `json:"WEBYYHOST"`
	WEBYYSWF  string `json:"WEBYYSWF"`
	HTML5     int    `json:"html5"`
	Stream    struct {
		Count int `json:"count"`
		Data  []struct {
			GameLiveInfo struct {
				ActivityCount      string      `json:"activityCount"`
				ActivityID         string      `json:"activityId"`
				AttendeeCount      interface{} `json:"attendeeCount"`
				Avatar180          string      `json:"avatar180"`
				BitRate            string      `json:"bitRate"`
				BussType           string      `json:"bussType"`
				CameraOpen         string      `json:"cameraOpen"`
				Channel            string      `json:"channel"`
				CodecType          string      `json:"codecType"`
				GameFullName       string      `json:"gameFullName"`
				GameHostName       string      `json:"gameHostName"`
				GameType           interface{} `json:"gameType"`
				Gid                string      `json:"gid"`
				Introduction       string      `json:"introduction"`
				IsSecret           string      `json:"isSecret"`
				Level              string      `json:"level"`
				LiveChannel        string      `json:"liveChannel"`
				LiveCompatibleFlag string      `json:"liveCompatibleFlag"`
				LiveID             string      `json:"liveId"`
				LiveSourceType     string      `json:"liveSourceType"`
				MultiStreamFlag    string      `json:"multiStreamFlag"`
				Nick               string      `json:"nick"`
				PrivateHost        string      `json:"privateHost"`
				ProfileHomeHost    string      `json:"profileHomeHost"`
				ProfileRoom        string      `json:"profileRoom"`
				RecommendStatus    string      `json:"recommendStatus"`
				RoomName           string      `json:"roomName"`
				ScreenType         string      `json:"screenType"`
				Screenshot         string      `json:"screenshot"`
				Sex                string      `json:"sex"`
				ShortChannel       string      `json:"shortChannel"`
				StartTime          string      `json:"startTime"`
				TotalCount         string      `json:"totalCount"`
				UID                string      `json:"uid"`
				Yyid               string      `json:"yyid"`
			} `json:"gameLiveInfo"`
			GameStreamInfoList []struct {
				IIsMaster           int           `json:"iIsMaster"`
				IIsMultiStream      int           `json:"iIsMultiStream"`
				IIsP2PSupport       int           `json:"iIsP2PSupport"`
				ILineIndex          int           `json:"iLineIndex"`
				IMobilePriorityRate int           `json:"iMobilePriorityRate"`
				IPCPriorityRate     int           `json:"iPCPriorityRate"`
				IWebPriorityRate    int           `json:"iWebPriorityRate"`
				LChannelID          int64         `json:"lChannelId"`
				LFreeFlag           int           `json:"lFreeFlag"`
				LPresenterUID       int64         `json:"lPresenterUid"`
				LSubChannelID       int64         `json:"lSubChannelId"`
				NewCFlvAntiCode     string        `json:"newCFlvAntiCode"`
				SCdnType            string        `json:"sCdnType"`
				SFlvAntiCode        string        `json:"sFlvAntiCode"`
				SFlvURL             string        `json:"sFlvUrl"`
				SFlvURLSuffix       string        `json:"sFlvUrlSuffix"`
				SHlsAntiCode        string        `json:"sHlsAntiCode"`
				SHlsURL             string        `json:"sHlsUrl"`
				SHlsURLSuffix       string        `json:"sHlsUrlSuffix"`
				SP2pAntiCode        string        `json:"sP2pAntiCode"`
				SP2pURL             string        `json:"sP2pUrl"`
				SP2pURLSuffix       string        `json:"sP2pUrlSuffix"`
				SStreamName         string        `json:"sStreamName"`
				VFlvIPList          []interface{} `json:"vFlvIPList"`
			} `json:"gameStreamInfoList"`
		} `json:"data"`
		IWebDefaultBitRate int    `json:"iWebDefaultBitRate"`
		Msg                string `json:"msg"`
		Status             int    `json:"status"`
		VMultiStreamInfo   []struct {
			IBitRate     int    `json:"iBitRate"`
			SDisplayName string `json:"sDisplayName"`
		} `json:"vMultiStreamInfo"`
	} `json:"stream"`
	Vappid int `json:"vappid"`
}

//不在播放抓取策略数据结构体 - 类型变化解决
type tT_PROFILE_INFO struct {
	Aid         int    			`json:"aid"`
	Avatar      string 			`json:"avatar"`
	Fans        int    			`json:"fans"`
	Host        string 			`json:"host"`
	//Lp          int    `json:"lp"`
	Nick        string 			`json:"nick"`
	ProfileRoom string 			`json:"profileRoom"`
	Sex         interface{}    	`json:"sex"`
	//Yyid        int    `json:"yyid"`
}
//func JsonDecodeTT_PROFILE_INFO(t string,data TT_PROFILE_INFO) (ret interface{},err error) {
//	err = json.Unmarshal([]byte(t), &data)
//	if err != nil {
//		return ret,err
//	}
//}


//huya解析方法
func (spider *Spider) huYaParser(p *Parser)  {
	switch p.Queue.Type {
	case "live_info":
		spider.huYaLiveInfo(p)
	case "live_list":
		spider.huYaLiveList(p)
	}
}

func (spider *Spider) huYaLiveInfo(p *Parser) {
	hyPlayerConfig := hyPlayerConfig{}
	Live := db.TableLive{}
	//哎呀，虎牙君找不到这个主播，要不搜索看看？
	if strings.Contains(string(p.Body),"哎呀，虎牙君找不到这个主播，要不搜索看看？") {
		Live.Live_uri = urlGetUri(p.Queue.Uri)
		Live.Live_pull_url= p.Queue.Uri
		Live.Live_is_online = "del"
		Live.Live_anchortv_sex = "0"
		Live.Live_online_user = "0"
		Live.Live_follow = "0"
		Live.Live_type_id = "0"
		Live.Created_at = strconv.FormatInt(time.Now().Unix(),10)
		Live.Updated_at = strconv.FormatInt(time.Now().Unix(),10)
		spider.WriteInfo <- &WriteInfo{
			TableLive:Live,
			Queue: p.Queue,
		}
		return
	}
	//该主播涉嫌违规，正在整改中……
	if strings.Contains(string(p.Body),"该主播涉嫌违规，正在整改中……") {
		Live.Live_uri = urlGetUri(p.Queue.Uri)
		Live.Live_pull_url= p.Queue.Uri
		Live.Live_is_online = "vio"
		Live.Live_anchortv_sex = "0"
		Live.Live_online_user = "0"
		Live.Live_follow = "0"
		Live.Live_type_id = "0"
		Live.Created_at = strconv.FormatInt(time.Now().Unix(),10)
		Live.Updated_at = strconv.FormatInt(time.Now().Unix(),10)
		spider.WriteInfo <- &WriteInfo{
			TableLive:Live,
			Queue: p.Queue,
		}
		return
	}

	//(?<=var hyPlayerConfig = )\{[\s\S]* \}(?=;) 不支持,坑货 fuck
	regexp_hyPlayerConfig := regexp.MustCompile(`var hyPlayerConfig = (\{[\s\S]* \});`);
	josn_hyPlayerConfig := regexp_hyPlayerConfig.FindSubmatch(p.Body)
	//fmt.Println(string(josn_hyPlayerConfig[1]))
	if josn_hyPlayerConfig != nil {
		err := json.Unmarshal(josn_hyPlayerConfig[1],&hyPlayerConfig)
		if err !=nil {
			log.Println("解析JSON失败."+err.Error())
			return
		}
	} else {
		log.Println("解析josn_hyPlayerConfig 为空.")
		//log.Println(string(p.Body))
		return
	}
	//.Live_is_online - 判断是在播 . 如果不再播 . 使用新策略爬取直播间基本数据
	if len(hyPlayerConfig.Stream.Data) == 0 {
		spider.huyaLive_is_online_no(p)
		return
	} else {
		Live.Live_is_online = "yes"
	}


	//.Live_pull_url #
	Live.Live_uri = liveReplaceSql(urlGetUri(p.Queue.Uri))

	//.Live_platform #
	Live.Live_platform = "huya"

	//.Live_title #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Introduction == "" {
		log.Println("Live_title:NULL")
		return
	}
	Live.Live_title = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Introduction)

	//.Live_anchortv_name #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Nick == "" {
		log.Println("Live_anchortv_name:NULL")
		return
	}
	Live.Live_anchortv_name = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Nick)

	//.Live_anchortv_photo #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Avatar180 == "" {
		log.Println("Live_anchortv_photo:NULL")
		return
	}
	Live.Live_anchortv_photo = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Avatar180)

	//.Live_cover #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Screenshot == "" {
		log.Println("Live_cover:NULL")
		return
	}
	Live.Live_cover = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Screenshot)

	//.Live_play #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.ProfileRoom == "" {
		log.Println("Live_play:NULL")
		return
	}
	Live.Live_play = liveReplaceSql("https://liveshare.huya.com/iframe/"+hyPlayerConfig.Stream.Data[0].GameLiveInfo.ProfileRoom)

	//.Live_class # gameFullName
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.GameFullName == "" {
		log.Println("Live_class:NULL")
		return
	}
	Live.Live_class = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.GameFullName)

	//.Live_online_user # totalCount
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.TotalCount == "" {
		log.Println("Live_online_user:NULL")
		return
	}
	Live.Live_online_user = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.TotalCount)

	//.Live_follow # activityCount
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.ActivityCount == "" {
		log.Println("Live_follow:NULL")
		return
	}
	Live.Live_follow = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.ActivityCount)

	//.Live_tag
	Live.Live_tag = ""

	//.Live_introduction - 这个不对
	Live.Live_introduction = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Introduction)

	//.live_anchortv_sex
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Sex != "" {
		Live.Live_anchortv_sex = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Sex)
	} else {
		Live.Live_anchortv_sex = "0"
	}

	//.Live_type_id
	//.Live_type_name
	Live.Live_type_id,Live.Live_type_name = liveGetMyTypeId(Live.Live_class)

	Live.Created_at = strconv.FormatInt(time.Now().Unix(),10)
	Live.Updated_at = strconv.FormatInt(time.Now().Unix(),10)
	Live.Live_pull_url= p.Queue.Uri

	spider.WriteInfo <- &WriteInfo{
		TableLive:Live,
		Queue: p.Queue,
	}
	return
}

//huya不在播抓取策略
func (spider *Spider) huyaLive_is_online_no(p *Parser) {
	TT_PROFILE_INFO := tT_PROFILE_INFO{}
	Live := db.TableLive{}

	regexp_TT_PROFILE_INFO := regexp.MustCompile(`var TT_PROFILE_INFO = (\{.+});var TT_PLAYER_CFG = `);
	josn_TT_PROFILE_INFO := regexp_TT_PROFILE_INFO.FindSubmatch(p.Body)
	if josn_TT_PROFILE_INFO != nil {
		err := json.Unmarshal(josn_TT_PROFILE_INFO[1],&TT_PROFILE_INFO)
		if err !=nil{
			log.Println("解析JSON失败."+err.Error())
			return
		}
	} else {
		log.Println("not 抓取失败")
		return
	}

	//14.Live_is_online
	Live.Live_is_online = "no"

	//1.Live_pull_url #
	//if TT_PROFILE_INFO.Host == "" {
	//	log.Println()("Live_pull_url:NULL")
	//	return
	//}
	Live.Live_uri = liveReplaceSql(urlGetUri(p.Queue.Uri))

	//2.Live_platform #
	Live.Live_platform = "huya"

	//3.Live_title #
	Live.Live_title = ""

	//4.Live_anchortv_name #
	if TT_PROFILE_INFO.Nick == "" {
		log.Println("Live_anchortv_name:NULL")
		return
	}
	Live.Live_anchortv_name = liveReplaceSql(TT_PROFILE_INFO.Nick)

	//5.Live_anchortv_photo #
	if TT_PROFILE_INFO.Avatar != "" {
		Live.Live_anchortv_photo = liveReplaceSql(TT_PROFILE_INFO.Avatar)
	} else {
		log.Panicln("头像为空")
		Live.Live_anchortv_photo = ""
	}


	//6.Live_cover #
	Live.Live_cover = ""

	//7.Live_play #
	if TT_PROFILE_INFO.ProfileRoom == "" {
		log.Println("Live_play:NULL")
		return
	}
	Live.Live_play = liveReplaceSql("https://liveshare.huya.com/iframe/"+TT_PROFILE_INFO.ProfileRoom)

	//8.Live_class #
	Live.Live_class = ""

	//9.Live_online_user #
	Live.Live_online_user = "0"

	//10.Live_follow #
	Live.Live_follow = "0"

	//11.Live_tag
	Live.Live_tag = ""

	//12.Live_introduction
	Live.Live_introduction = ""

	//13.live_anchortv_sex
	if TT_PROFILE_INFO.Sex != nil {
		if reflect.TypeOf(TT_PROFILE_INFO.Sex).String() == "float64" {
			Live.Live_anchortv_sex = liveReplaceSql(strconv.FormatFloat(TT_PROFILE_INFO.Sex.(float64), 'E', -1, 64))
		} else {
			Live.Live_anchortv_sex = liveReplaceSql(TT_PROFILE_INFO.Sex.(string))
		}
	} else {
		Live.Live_anchortv_sex = "0"
	}

	//15.Live_type_id
	//16.Live_type_name
	Live.Live_type_id = "0"
	Live.Live_type_name = ""

	Live.Created_at = strconv.FormatInt(time.Now().Unix(),10)
	Live.Updated_at = strconv.FormatInt(time.Now().Unix(),10)
	Live.Live_pull_url= p.Queue.Uri

	spider.WriteInfo <- &WriteInfo{
		TableLive:Live,
		Queue: p.Queue,
	}
	return
}


type huYaLiveListStruct struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data struct {
		Page int `json:"page"`				//当前页
		PageSize int `json:"pageSize"`		//每页数量
		TotalPage int `json:"totalPage"`	//总页数
		TotalCount int `json:"totalCount"`	//总数据
		Datas []struct {
			GameFullName string `json:"gameFullName"`					//分类名称
			GameHostName string `json:"gameHostName"`					//分类url名称
			BoxDataInfo interface{} `json:"boxDataInfo"`				//
			TotalCount string `json:"totalCount"`						//直播间人数?
			RoomName string `json:"roomName"`							//直播间标题
			BussType string `json:"bussType"`							//当前页位置?
			Screenshot string `json:"screenshot"`						//直播间封面
			PrivateHost string `json:"privateHost"`						//直播间URI
			Nick string `json:"nick"`									//主播昵称
			Avatar180 string `json:"avatar180"`							//主播头像
			Gid string `json:"gid"`										//
			Introduction string `json:"introduction"`					//直播间简介
			RecommendStatus string `json:"recommendStatus"`				//推荐状态?
			RecommendTagName string `json:"recommendTagName"`			//推荐TAG标签
			IsBluRay string `json:"isBluRay"`							//
			BluRayMBitRate string `json:"bluRayMBitRate"`				//清晰度,如10M
			ScreenType string `json:"screenType"`						//屏幕类型
			LiveSourceType string `json:"liveSourceType"`				//直播屏幕类型
			UID string `json:"uid"`										//用户ID
			Channel string `json:"channel"`								//渠道?
			LiveChannel string `json:"liveChannel"`						//直播渠道?
			ImgRecInfo interface{} `json:"imgRecInfo"`					//
			AliveNum string `json:"aliveNum"`
			Attribute interface{} `json:"attribute"`					//属性
			ProfileRoom string `json:"profileRoom"`						//房间号
			IsRoomPay int `json:"isRoomPay"`
			RoomPayTag string `json:"roomPayTag"`
		} `json:"datas"`
		Time int `json:"time"`
	} `json:"data"`
}
func (spider *Spider) huYaLiveList(p *Parser)  {
	huYaLiveListStruct := huYaLiveListStruct{}
	if err := json.Unmarshal(p.Body, &huYaLiveListStruct); err != nil {
		log.Panicln(err.Error())
		return
	}
	if huYaLiveListStruct.Status != 200 {
		return
	}

	//判断当前页是否大于总页数
	if huYaLiveListStruct.Data.Page > huYaLiveListStruct.Data.TotalPage {
		return
	}
	//更多列表
	for i:=1; i<=huYaLiveListStruct.Data.TotalPage; i++ {
		spider.ChanProduceList <- &db.Queue{
			Platform: p.Queue.Platform,
			Uri:      "https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&tagAll=0&page="+strconv.Itoa(i),
			Type:     "live_list",
			Event:    p.Queue.Event,
		}

	}

	for _, v := range huYaLiveListStruct.Data.Datas {
		spider.ChanProduceList <- &db.Queue{
			Platform: p.Queue.Platform,
			Uri:      "https://www.huya.com/"+v.ProfileRoom,
			Type:     "live_info",
			Event:    "",
		}
	}
}


/**
 * 通过url获取uri
 *
 */
func urlGetUri(url string) string {
	regexpUrl := regexp.MustCompile(`[a-zA-z]+://[^\s]*/(.+)`)
	uris := regexpUrl.FindSubmatch([]byte(url))
	return string(uris[1])
}

/**
 * 过滤字符串 , 防止sql注入跟其他错误
 */
func liveReplaceSql(data string) (ret string) {
	ret = strings.Replace(data, "'", "\"", -1)
	ret = strings.Replace(ret, "\\", "\\\\", -1)
	//ret = strings.Replace(ret, " ", "_", -1)

	return ret
}


/**
 * 获取自定义直播分类
 * @param string 平台分类
 * @return string 自定义分类ID string
 */
func liveGetMyTypeId(live_class string) (live_type_id string,live_type_name string) {
	live_type_id = "0"
	live_type_name = ""
	for i := 0; i<len(LiveMyTypeData); i++ {
		if strings.Contains(LiveMyTypeData[i]["subset"],"#"+live_class+"#") {
			live_type_id = LiveMyTypeData[i]["type_id"]
			live_type_name = LiveMyTypeData[i]["name"]
			break
		}
	}

	return live_type_id,live_type_name
}

//初始化分类映射数据
func InitLiveMyType() (err error) {
	fmt.Println("init type data")
	if LiveMyTypeData, err = mysql.Table("live_type").Select("type_id,name,subset").Order("`order` asc").Get(); err != nil {
		panic("初始化失败 . 分类映射数据出错 ,"+err.Error())
	}
	return err
}
var	LiveMyTypeData []map[string]string