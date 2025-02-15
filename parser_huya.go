package main

import (
	"encoding/json"
	//uuid "github.com/satori/go.uuid"
	"ikanzhibo/db"
	"log"
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
	switch p.Queue.QueueSet.QueueType {
	case "live_info":
		spider.huYaLiveInfo(p)
	case "live_list":
		spider.huYaLiveList(p)
	}
}

func (spider *Spider) huYaLiveInfo(p *Parser) {
	hyPlayerConfig := hyPlayerConfig{}
	p.Queue.LiveData.Spider_pull_time = strconv.FormatInt(time.Now().Unix(),10)	// *
	p.Queue.LiveData.Spider_pull_url = p.Queue.QueueSet.Request.Url				// *
	p.Queue.LiveData.Live_platform = "huya"			//Live_platform #
	//哎呀，虎牙君找不到这个主播，要不搜索看看？
	if strings.Contains(string(p.Body),"哎呀，虎牙君找不到这个主播，要不搜索看看？") {
		p.Queue.LiveData.Live_uri = urlGetUri(p.Queue.QueueSet.Request.Url)
		p.Queue.LiveData.Live_is_online = "del"
		p.Queue.LiveData.Live_anchortv_sex = "0"
		p.Queue.LiveData.Live_online_user = "0"
		p.Queue.LiveData.Live_follow = "0"
		p.Queue.LiveData.Live_type_id = "0"
		p.Queue.LiveData.Created_at = strconv.FormatInt(time.Now().Unix(),10)
		p.Queue.LiveData.Updated_at = strconv.FormatInt(time.Now().Unix(),10)
		p.Queue.LiveData.Live_play_time = "0"
		p.Queue.LiveData.Live_play_end_time = "0"
		spider.ChanWriteInfo <- &p.Queue
		return
	}
	//该主播涉嫌违规，正在整改中……
	if strings.Contains(string(p.Body),"该主播涉嫌违规，正在整改中……") {
		p.Queue.LiveData.Live_uri = urlGetUri(p.Queue.QueueSet.Request.Url)
		p.Queue.LiveData.Live_is_online = "vio"
		p.Queue.LiveData.Live_anchortv_sex = "0"
		p.Queue.LiveData.Live_online_user = "0"
		p.Queue.LiveData.Live_follow = "0"
		p.Queue.LiveData.Live_type_id = "0"
		p.Queue.LiveData.Created_at = strconv.FormatInt(time.Now().Unix(),10)
		p.Queue.LiveData.Updated_at = strconv.FormatInt(time.Now().Unix(),10)
		p.Queue.LiveData.Live_play_time = "0"
		p.Queue.LiveData.Live_play_end_time = "0"
		spider.ChanWriteInfo <- &p.Queue
		return
	}

	//(?<=var hyPlayerConfig = )\{[\s\S]* \}(?=;) 不支持,坑货 fuck
	regexp_hyPlayerConfig := regexp.MustCompile(`var hyPlayerConfig = (\{[\s\S]* \});`);
	josn_hyPlayerConfig := regexp_hyPlayerConfig.FindSubmatch(p.Body)
	//fmt.Println(string(josn_hyPlayerConfig[1]))
	if josn_hyPlayerConfig != nil {
		err := json.Unmarshal(josn_hyPlayerConfig[1],&hyPlayerConfig)
		if err !=nil {
			log.Println("解析JSON失败. ERR: "+err.Error() +"\n"+p.Queue.QueueSet.Request.Url)
			return
		}
	} else {
		log.Println("解析josn_hyPlayerConfig 为空.\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	//.Live_is_online - 判断是在播 . 如果不再播 . 使用新策略爬取直播间基本数据
	if len(hyPlayerConfig.Stream.Data) == 0 {
		spider.huyaLive_is_online_no(p)
		return
	} else {
		p.Queue.LiveData.Live_is_online = "yes"
	}

	//Live_play_time
	p.Queue.LiveData.Live_play_time = hyPlayerConfig.Stream.Data[0].GameLiveInfo.StartTime
	if (p.Queue.LiveData.Live_play_time == "") {
		p.Queue.LiveData.Live_play_time = "0";
	}
	p.Queue.LiveData.Live_play_end_time = "0";

	//.Live_uri #
	p.Queue.LiveData.Live_uri = hyPlayerConfig.Stream.Data[0].GameLiveInfo.ProfileHomeHost

	p.Queue.LiveData.Platform_room_id = hyPlayerConfig.Stream.Data[0].GameLiveInfo.ProfileRoom

	//.Live_title #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Introduction == "" {
		log.Println("Live_title:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_title = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Introduction)

	//.Live_anchortv_name #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Nick == "" {
		log.Println("Live_anchortv_name:NULL \n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_anchortv_name = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Nick)

	//.Live_anchortv_photo #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Avatar180 == "" {
		log.Println("Live_anchortv_photo:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_anchortv_photo = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Avatar180)

	//.Live_cover #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Screenshot == "" {
		log.Println("Live_cover:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_cover = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Screenshot)

	//.Live_play #
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.ProfileRoom == "" {
		log.Println("Live_play:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_play = liveReplaceSql("https://liveshare.huya.com/iframe/"+hyPlayerConfig.Stream.Data[0].GameLiveInfo.ProfileRoom)

	//.Live_class # gameFullName
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.GameFullName == "" {
		log.Println("Live_class:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_class = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.GameFullName)

	//.Live_online_user # totalCount
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.TotalCount == "" {
		log.Println("Live_online_user:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_online_user = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.TotalCount)

	//.Live_follow # activityCount
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.ActivityCount == "" {
		log.Println("Live_follow:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_follow = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.ActivityCount)

	//.Live_tag
	p.Queue.LiveData.Live_tag = ""

	//.Live_introduction - 这个不对
	p.Queue.LiveData.Live_introduction = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Introduction)

	//.live_anchortv_sex
	if hyPlayerConfig.Stream.Data[0].GameLiveInfo.Sex != "" {
		p.Queue.LiveData.Live_anchortv_sex = liveReplaceSql(hyPlayerConfig.Stream.Data[0].GameLiveInfo.Sex)
	} else {
		p.Queue.LiveData.Live_anchortv_sex = "0"
	}

	//.Live_type_id
	//.Live_type_name
	p.Queue.LiveData.Live_type_id,p.Queue.LiveData.Live_type_name = platformTypeToLocal(p.Queue.LiveData.Live_class)

	p.Queue.LiveData.Created_at = strconv.FormatInt(time.Now().Unix(),10)
	p.Queue.LiveData.Updated_at = strconv.FormatInt(time.Now().Unix(),10)

	spider.ChanWriteInfo <- &p.Queue
	return
}

//huya不在播抓取策略
func (spider *Spider) huyaLive_is_online_no(p *Parser) {
	TT_PROFILE_INFO := tT_PROFILE_INFO{}

	regexp_TT_PROFILE_INFO := regexp.MustCompile(`var TT_PROFILE_INFO = (\{.+});var TT_PLAYER_CFG = `);
	josn_TT_PROFILE_INFO := regexp_TT_PROFILE_INFO.FindSubmatch(p.Body)
	if josn_TT_PROFILE_INFO != nil {
		err := json.Unmarshal(josn_TT_PROFILE_INFO[1],&TT_PROFILE_INFO)
		if err !=nil{
			log.Println("解析JSON失败."+err.Error() +"\n"+p.Queue.QueueSet.Request.Url)
			return
		}
	} else {
		log.Println("not 抓取失败\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Spider_pull_time = strconv.FormatInt(time.Now().Unix(),10)
	p.Queue.LiveData.Spider_pull_url = p.Queue.QueueSet.Request.Url
	p.Queue.LiveData.Live_anchortv_sex = "0"
	p.Queue.LiveData.Live_online_user = "0"
	p.Queue.LiveData.Live_follow = "0"
	p.Queue.LiveData.Live_type_id = "0"
	p.Queue.LiveData.Live_play_time = "0"
	p.Queue.LiveData.Live_play_end_time = "0"
	p.Queue.LiveData.Created_at = strconv.FormatInt(time.Now().Unix(),10)
	p.Queue.LiveData.Updated_at = strconv.FormatInt(time.Now().Unix(),10)

	//14.Live_is_online
	p.Queue.LiveData.Live_is_online = "no"

	p.Queue.LiveData.Live_uri = liveReplaceSql(urlGetUri(p.Queue.QueueSet.Request.Url))
	p.Queue.LiveData.Platform_room_id = TT_PROFILE_INFO.ProfileRoom
	//2.Live_platform #
	p.Queue.LiveData.Live_platform = "huya"

	//4.Live_anchortv_name #
	if TT_PROFILE_INFO.Nick == "" {
		log.Println("Live_anchortv_name:NULL\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	p.Queue.LiveData.Live_anchortv_name = liveReplaceSql(TT_PROFILE_INFO.Nick)

	p.Queue.LiveData.Live_play_time = "0";
	p.Queue.LiveData.Live_play_end_time = "0";

	//5.Live_anchortv_photo #
	if TT_PROFILE_INFO.Avatar != "" {
		p.Queue.LiveData.Live_anchortv_photo = liveReplaceSql(TT_PROFILE_INFO.Avatar)
	} else {
		log.Println("头像为空\n"+p.Queue.QueueSet.Request.Url)
	}

	//10.Live_follow #
	p.Queue.LiveData.Live_follow = strconv.Itoa(TT_PROFILE_INFO.Fans)

	//12.Live_introduction
	//p.Queue.LiveData.Live_introduction = ""

	p.Queue.LiveData.Updated_at = strconv.FormatInt(time.Now().Unix(),10)

	spider.ChanWriteInfo <- &p.Queue
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
		log.Println(err.Error()+"\n"+p.Queue.QueueSet.Request.Url)
		return
	}
	if huYaLiveListStruct.Status != 200 {
		log.Println("huYaLiveListStruct.Status != 200\n"+p.Queue.QueueSet.Request.Url)
		return
	}

	//判断当前页是否大于总页数
	if huYaLiveListStruct.Data.Page > huYaLiveListStruct.Data.TotalPage {
		return
	}
	//更多列表
	for i:=1; i<=huYaLiveListStruct.Data.TotalPage; i++ {
		spider.ChanProduceList <- &db.Queue{
			QueueSet:db.QueueSet{
				Request:     db.Request{
					Url: "https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&tagAll=0&page="+strconv.Itoa(i),
				},
				QueueType: "live_list",
				Platform: "huya",
			},
		}

	}

	for _, v := range huYaLiveListStruct.Data.Datas {
		spider.ChanProduceList <- &db.Queue{
			QueueSet:db.QueueSet{
				Request:     db.Request{
					Url: "https://www.huya.com/"+v.ProfileRoom,
				},
				QueueType: "live_info",
				Platform: "huya",
			},
		}
	}
}
