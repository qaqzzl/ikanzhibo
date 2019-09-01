package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"log"
	"strconv"
	"time"
)

func (spider *Spider) douYuParser(p *Parser)  {
	switch p.Queue.Type {
	case "live_info":
		spider.douYuLiveInfo(p)
	case "live_list":
		spider.douYuLiveList(p)
	}
}

type douYuLiveInfo struct {
	VoddUploadURL           string        `json:"VoddUploadUrl"`
	BarragePraise           int           `json:"barrage_praise"`
	BarrageTimeoutDowngrade int           `json:"barrage_timeout_downgrade"`
	BindVodCateURL          string        `json:"bind_vodCateUrl"`
	Black                   []interface{} `json:"black"`
	CacheTime               int           `json:"cache_time"`
	CanSendGift             string        `json:"can_send_gift"`
	CateID                  int           `json:"cate_id"`
	DefaultRankName string `json:"defaultRankName"`
	FaceList        string `json:"faceList"`
	H5Default   int           `json:"h5_default"`
	H5GuardJS   []interface{} `json:"h5_guardJS"`
	HomeAdInfo  []interface{} `json:"home_ad_info"`
	HotPostStatus        string `json:"hot_post_status"`
	IsNewbie             int    `json:"is_newbie"`
	NearShowTime  interface{} `json:"near_show_time"`
	PageURL       string      `json:"page_url"`
	PlayerBarrage int         `json:"player_barrage"`
	PostList      []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"post_list"`
	Room    struct {
		Avatar struct {
			Big    string `json:"big"`
			Middle string `json:"middle"`
			Small  string `json:"small"`
		} `json:"avatar"`
		AvatarMid   string `json:"avatar_mid"`
		AvatarSmall string `json:"avatar_small"`
		BanDisplay  int    `json:"ban_display"`
		//BgimgSrc    string `json:"bgimg_src"`
		CanSendGift string `json:"can_send_gift"`
		Cate1ID     string `json:"cate1_id"`
		Cate2ID     string `json:"cate2_id"`
		Cate3ID     string `json:"cate3_id"`
		CateID      int    `json:"cate_id"`
		CateLimit   struct {
			LimitNum       int `json:"limit_num"`
			LimitThreshold int `json:"limit_threshold"`
			LimitTime      int `json:"limit_time"`
			LimitType      int `json:"limit_type"`
		} `json:"cate_limit"`
		CategoryID        string        `json:"category_id"`
		CfmGiftList       []interface{} `json:"cfmGiftList"`
		ChatAgeLimit      string        `json:"chat_age_limit"`
		ChatCdFactor      string        `json:"chat_cd_factor"`
		ChatGroup         string        `json:"chat_group"`
		ChatLevel         string        `json:"chat_level"`
		ChildID           string        `json:"child_id"`
		Cityname          string        `json:"cityname"`
		CoverSrc          string        `json:"coverSrc"`
		Cq                string        `json:"cq"`
		DefaultSrc        string        `json:"defaultSrc"`
		DetailsData       struct{}      `json:"detailsData"`
		EffectInfo        []interface{} `json:"effectInfo"`
		EmperorPush       []interface{} `json:"emperorPush"`
		EndTime           string        `json:"end_time"`
		Eticket           []interface{} `json:"eticket"`
		H5wsproxy  []struct {
			Domain string `json:"domain"`
			Port   string `json:"port"`
		} `json:"h5wsproxy"`
		IconEndTime            string `json:"icon_end_time"`
		IconID                 string `json:"icon_id"`
		IconStartTime          string `json:"icon_start_time"`
		IsDefaultAvatar        int    `json:"isDefaultAvatar"`
		IsNzRoom               int    `json:"isNzRoom"`
		IsPubgmRoom            int    `json:"isPubgmRoom"`
		IsVertical             int    `json:"isVertical"`
		IsDiy                  string `json:"is_diy"`
		IsHighGame             int    `json:"is_high_game"`
		IsMultibit             string `json:"is_multibit"`
		IsPassword             int    `json:"is_password"`
		IsSetFansBadge         int    `json:"is_set_fans_badge"`
		IsShowRankList         string `json:"is_show_rank_list"`
		IsVideoHighQualityTime int    `json:"is_video_high_quality_time"`
		IsVr                   int    `json:"is_vr"`
		Isvertival             int    `json:"isvertival"`
		LevelInfo              struct {
			EndTime        string  `json:"end_time"`
			ExpInc         float64 `json:"exp_inc"`
			Experience     float64 `json:"experience"`
			IsKeepTaskComp bool    `json:"isKeepTaskComp"`
			IsMaxed        bool    `json:"isMaxed"`
			Level          string  `json:"level"`
			MinExp         int     `json:"min_exp"`
			NextLevel      int     `json:"next_level"`
			Progress       float64 `json:"progress"`
		} `json:"levelInfo"`
		Multirates     []struct {
			Name string `json:"name"`
			Type int    `json:"type"`
		} `json:"multirates"`
		Music struct {
			DmSp int `json:"dm_sp"`
			DmSt int `json:"dm_st"`
			DmUm int `json:"dm_um"`
		} `json:"music"`
		Nickname    string `json:"nickname"`
		Nowtime        int `json:"nowtime"`
		OfficialAnchor struct {
			Image    string `json:"image"`
			Ioa      int    `json:"ioa"`
			JumpType int    `json:"jumpType"`
			Od       string `json:"od"`
			URL      string `json:"url"`
		} `json:"officialAnchor"`
		OpenFullScreen int    `json:"open_full_screen"`
		OwnerAvatar    string `json:"owner_avatar"`
		OwnerName      string `json:"owner_name"`
		OwnerUID       int    `json:"owner_uid"`
		P2pSetting     struct {
			MDm         int `json:"m_dm"`
			NameID      int `json:"name_id"`
			OnlineLimit int `json:"online_limit"`
			PlanID      int `json:"plan_id"`
			Player      int `json:"player"`
			WDm         int `json:"w_dm"`
		} `json:"p2p_setting"`
		Pwd          string `json:"pwd"`
		RoomID   int `json:"room_id"`
		RoomIdle struct {
			Active      int `json:"active"`
			MinuteLimit int `json:"minute_limit"`
		} `json:"room_idle"`
		RoomLabelRightFlag int    `json:"room_label_right_flag"`
		RoomName           string `json:"room_name"`
		RoomPic            string `json:"room_pic"`
		RoomPlugin         string `json:"room_plugin"`
		RoomSrc            string `json:"room_src"`
		RoomURL            string `json:"room_url"`
		SecondLvlName      string `json:"second_lvl_name"`
		Share              struct {
			Common string `json:"common"`
			Flash  string `json:"flash"`
			Video  string `json:"video"`
		} `json:"share"`
		ShowDetails          string `json:"show_details"`
		ShowID               int    `json:"show_id"`
		ShowStatus           int    `json:"show_status"`
		ShowTime             int    `json:"show_time"`
		SimplifyBulletScreen struct {
			Condition struct {
				MinNum int `json:"minNum"`
			} `json:"condition"`
			Rule struct {
				Level   int `json:"level"`
				MaxNum  int `json:"maxNum"`
				Percent int `json:"percent"`
			} `json:"rule"`
		} `json:"simplifyBulletScreen"`
		SpeakSet struct {
			OnlyAdmin  string `json:"onlyAdmin"`
			SpeakCd    string `json:"speakCd"`
			SpeakLv    string `json:"speakLv"`
			Speaklimit int    `json:"speaklimit"`
		} `json:"speakSet"`
		St         int    `json:"st"`
		StsignRoom struct {
			Ctime string `json:"ctime"`
			State struct {
				Mobile int `json:"mobile"`
				Yzxx   int `json:"yzxx"`
			} `json:"state"`
		} `json:"stsign_room"`
		//Tags                       string   `json:"tags"`
		UpID                       string   `json:"up_id"`
		VideoHighQualityNum        string   `json:"video_high_quality_num"`
		VideoHighQualityResolution string   `json:"video_high_quality_resolution"`
		Videop                     string   `json:"videop"`
		VipID                      int      `json:"vipId"`

		YubaJumpURL string `json:"yuba_jump_url"`
	} `json:"room"`
	RoomArgs struct {
		NoHome       int    `json:"no_home"`
		NoHomeTime   string `json:"no_home_time"`
		ResPath      string `json:"res_path"`
		RPCSwitch    int    `json:"rpc_switch"`
		ServerConfig string `json:"server_config"`
		SwfURL       string `json:"swf_url"`
	} `json:"room_args"`
	SeoInfo struct {
		SeoDescription string `json:"seo_description"`
		SeoKeyword     string `json:"seo_keyword"`
		SeoTitle       string `json:"seo_title"`
	} `json:"seo_info"`
	ServiceSwitch struct {
		BarrageReply     int `json:"barrageReply"`
		FastBarrage      int `json:"fastBarrage"`
		OffLineFriendRec int `json:"offLineFriendRec"`
	} `json:"serviceSwitch"`
	ShareSwfURL string `json:"share_swf_url"`
	SwfURL      string `json:"swf_url"`
	VarIsYz     bool   `json:"var_is_yz"`
	VarYzPkName string `json:"var_yz_pk_name"`
	VideoTitle      string `json:"video_title"`
}

func (spider *Spider) douYuLiveInfo(p *Parser) {
	Live := db.TableLive{}
	douYuLiveInfo := douYuLiveInfo{}

	if err := json.Unmarshal(p.Body, &douYuLiveInfo); err != nil {
		log.Println(err.Error()+"\n"+p.Queue.Uri)
		return
	}

	//.Live_is_online - 判断是在播
	if douYuLiveInfo.Room.ShowStatus == 1 {
		Live.Live_is_online = "yes"
	}

	//.Live_uri #
	Live.Live_uri = liveReplaceSql(urlGetUri(p.Queue.Uri))

	//.Live_platform #
	Live.Live_platform = p.Queue.Platform

	//.Live_title #
	Live.Live_title = liveReplaceSql(douYuLiveInfo.Room.RoomName)

	//.Live_anchortv_name #
	Live.Live_anchortv_name = liveReplaceSql(douYuLiveInfo.Room.Nickname)

	//.Live_anchortv_photo #
	Live.Live_anchortv_photo = liveReplaceSql(douYuLiveInfo.Room.Avatar.Big)

	//.Live_cover #
	Live.Live_cover = liveReplaceSql(douYuLiveInfo.Room.RoomPic)

	//.Live_play #
	Live.Live_play = liveReplaceSql(douYuLiveInfo.Room.RoomURL)

	//.Live_class # gameFullName
	Live.Live_class = liveReplaceSql(douYuLiveInfo.Room.SecondLvlName)

	//.Live_online_user # totalCount
	Live.Live_online_user = "0"

	//.Live_follow # activityCount
	Live.Live_follow = "0"

	//.Live_tag
	Live.Live_tag = ""

	//.Live_introduction
	Live.Live_introduction = liveReplaceSql(douYuLiveInfo.Room.ShowDetails)

	//.live_anchortv_sex
	Live.Live_anchortv_sex = "0"

	//.Live_type_id
	//.Live_type_name
	Live.Live_type_id,Live.Live_type_name = liveGetMyTypeId(Live.Live_class)

	Live.Created_at = strconv.FormatInt(time.Now().Unix(),10)
	Live.Updated_at = strconv.FormatInt(time.Now().Unix(),10)

	spider.ChanWriteInfo <- &WriteInfo{
		TableLive:Live,
		Queue: p.Queue,
	}
	return
}



type douYuLiveLists struct {
	Code int `json:"code"`
	Data struct {
		Ct struct {
			Iv    int    `json:"iv"`
			Ivcv  int    `json:"ivcv"`
			Tag   int    `json:"tag"`
			Tn    string `json:"tn"`
			Vmcm  string `json:"vmcm"`
			Vmcrr int    `json:"vmcrr"`
		} `json:"ct"`
		Pgcnt int `json:"pgcnt"`
		Rl    []struct {
			Av     string `json:"av"`
			Bid    int    `json:"bid"`
			C2name string `json:"c2name"`
			C2url  string `json:"c2url"`
			Chanid int    `json:"chanid"`
			Cid1   int    `json:"cid1"`
			Cid2   int    `json:"cid2"`
			Cid3   int    `json:"cid3"`
			Clis   int    `json:"clis"`
			Dot    int    `json:"dot"`
			Gldid  int    `json:"gldid"`
			Icdata struct {
				Six00 struct {
					H   int    `json:"h"`
					URL string `json:"url"`
					W   int    `json:"w"`
				} `json:"600"`
			} `json:"icdata"`
			Icv1  [][]interface{} `json:"icv1"`
			Ioa   int             `json:"ioa"`
			Iv    int             `json:"iv"`
			Nn    string          `json:"nn"`
			Od    string          `json:"od"`
			Ol    int             `json:"ol"`
			Ot    int             `json:"ot"`
			Rgrpt int             `json:"rgrpt"`
			Rid   int             `json:"rid"`
			Rkic  string          `json:"rkic"`
			Rn    string          `json:"rn"`
			Rpos  int             `json:"rpos"`
			Rs1   string          `json:"rs1"`
			Rs16  string          `json:"rs16"`
			Rt    int             `json:"rt"`
			Subrt int             `json:"subrt"`
			Topid int             `json:"topid"`
			UID   int             `json:"uid"`
			URL   string          `json:"url"`
			Utag  interface{}     `json:"utag"`
		} `json:"rl"`
	} `json:"data"`
	Msg string `json:"msg"`
}
func (spider *Spider) douYuLiveList(p *Parser)  {
	douYuLiveLists := douYuLiveLists{}
	if err := json.Unmarshal(p.Body, &douYuLiveLists); err != nil {
		log.Println(err.Error()+"\n"+p.Queue.Uri)
		return
	}
	if douYuLiveLists.Code != 0 {
		log.Println("douYuLiveLists.Status != 200\n"+p.Queue.Uri)
		return
	}

	//判断当前页是否大于总页数
	if douYuLiveLists.Data.Pgcnt > douYuLiveLists.Data.Pgcnt {
		return
	}
	//更多列表
	for i:=1; i<=douYuLiveLists.Data.Pgcnt; i++ {
		spider.ChanProduceList <- &db.Queue{
			Platform: p.Queue.Platform,
			Uri:      "https://www.douyu.com/gapi/rkc/directory/0_0/"+strconv.Itoa(i),
			Type:     "live_list",
			Event:    p.Queue.Event,
		}

	}

	for _, v := range douYuLiveLists.Data.Rl {
		spider.ChanProduceList <- &db.Queue{
			Platform: p.Queue.Platform,
			Uri:      "https://www.douyu.com/betard/"+strconv.Itoa(v.Rid),
			Type:     "live_info",
			Event:    "",
		}
	}
}
