package db

import (
	"encoding/json"
	"ikanzhibo/db/mysql"
	"ikanzhibo/db/redis"
	"log"
	"reflect"
	"strconv"
	"time"
)

//redis表名
var (
	RedisFollowOfflineList			= "live_follow_offline_list"				//被关注&&不在线直播间队列
	RedisFollowOffSet				= "live_follow_offline_set"					//被关注&&不在线直播间集合,定时跟数据库做一致性同步	model

	RedisNotFollowOfflineList		= "live_not_follow_offline_list"			//未关注&&不在线直播间队列
	RedisNotFollowOffSet			= "live_not_follow_offline_set"				//未关注&&不在线直播间集合,定时跟数据库做一致性同步

	RedisOnlineList					= "live_online_list"						//在线直播间队列
	RedisOnlineSet					= "live_online_set"							//在线直播间集合

	RedisListList		 			= "live_list_list"							//生产任务队列 - 未爬取
	RedisListOnceSet				= "live_list_once_set"						//生产任务集合 - 已爬取

	RedisInfoOnceSet				= "live_info_once_set"						//总info set


	RedisOnlineNotice				= "event_online_notice_list"				//事件 - 开播通知

)

//平台表结构体
type Platform struct {
	PlatformId		int
	Mark			string
	Name			string
	Domain			string
	PullUrl			string
	Status			int
	DomainUrl		string
}
var platforms []Platform




//任务队列结构体
type Queue struct {
	LiveId	 	string		//直播ID
	Platform 	string		//所属平台
	Uri 		string
	Type 		string		//任务类型 , live_info:直播间数据, live_list:直播列表
	Event		string		//触发事件, online_notice:开播通知, send_barrage:发送弹幕 , listener_barrage:监听弹幕 多个事件用逗号隔开
}

// 直播间数据结构体
type TableLive struct {
	Live_title				string			//标题
	Live_anchortv_name		string			//主播名称
	Live_anchortv_photo		string			//主播头像
	Live_anchortv_sex		string			//主播性别 0-保密 1-女 2-男
	Live_cover				string			//直播间封面
	Live_play				string			//播放地址
	Live_class				string			//平台直播间分类
	Live_tag				string			//直播间标签
	Live_introduction		string			//直播间简介
	Live_online_user		string			//直播间在线人数
	Live_follow				string			//被关注人数
	Live_uri				string			//直播间地址
	Live_pull_url			string			//抓取url
	Live_type_id			string			//自定义分类ID
	Live_type_name			string			//自定义分类ID
	Live_platform			string			//所属平台
	Live_is_online			string			//直播间是否在播 ,yes|no
	Created_at				string
	Updated_at				string
	Queue_id				string			//队列ID
}


//获取所有未开播但有人订阅的直播间地址
func GetFollowOffline() (l []Queue, err error) {
	rconn := redis.GetConn()
	defer rconn.Close()

	//查询redis
	rlist, err := rconn.Do("SMEMBERS", RedisFollowOffSet)
	if err != nil {
		log.Println("Redis SMEMBERS error",RedisFollowOffSet)
		return l, err
	}
	v := reflect.ValueOf(rlist)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	len := v.Len()
	if len != 0 {
		for i := 0; i < len; i++ {
			queue := Queue{}
			json.Unmarshal(v.Index(i).Interface().([]byte), &queue)
			vo := Queue{
				LiveId: queue.LiveId,
				Platform: queue.Platform,
				Uri:      queue.Uri,
			}
			l = append(l, vo)
		}
		return l,err
	}

	//redis不存在 , 查询mysql
	list, err := mysql.Conn().QueryAll("select l.live_id,l.live_uri,l.live_platform from live as l JOIN live_user_follow as luf ON luf.live_id = l.live_id WHERE luf.`status`=1 AND luf.is_notice=1 AND l.live_is_online='no'")
	if err != nil {
		log.Println("MySql error", RedisFollowOffSet)
		return l, err
	}
	for _, v := range list {
		vo := Queue{
			LiveId: v["live_id"],
			Platform: v["live_platform"],
			Uri:      "https://www.huya.com/"+v["live_uri"],
		}
		l = append(l, vo)

		str,_ := json.Marshal(vo)
		rconn.Do("SADD", RedisFollowOffSet, str)
	}
	return l,err
}

//获取未开播 && 未关注
func GetNotFollowOffline() (l []Queue, err error) {
	rconn := redis.GetConn()
	defer rconn.Close()

	rlist, err := rconn.Do("SMEMBERS", RedisNotFollowOffSet)
	if err != nil {
		log.Println("Redis SMEMBERS error",RedisNotFollowOffSet)
		return l, err
	}
	v := reflect.ValueOf(rlist)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	len := v.Len()
	if len != 0 {
		for i := 0; i < len; i++ {
			queue := Queue{}
			json.Unmarshal(v.Index(i).Interface().([]byte), &queue)
			vo := Queue{
				LiveId: queue.LiveId,
				Platform: queue.Platform,
				Uri:      queue.Uri,
			}
			l = append(l, vo)
		}
		return l,err
	}

	//redis不存在 , 查询mysql

	/**
	select l.live_id,l.live_uri,l.live_platform from live as l
	LEFT JOIN (select * from live_user_follow) as luf ON l.live_id=luf.live_id
	where l.live_is_online = 'no' and  luf.live_id is null
	*/

	/**
	select live_id,live_uri,live_platform from live
	WHERE live_is_online = 'no' and live_id not in ( select live_id from live_user_follow)
	*/

	list, err := mysql.Conn().QueryAll(`select l.live_id,l.live_uri,l.live_platform from live as l
	LEFT JOIN (select * from live_user_follow) as luf ON l.live_id=luf.live_id
	where l.live_is_online = 'no' and  luf.live_id is null`)
	if err != nil {
		log.Println("MySql error", RedisNotFollowOffSet)
		return l, err
	}
	for _, v := range list {
		vo := Queue{
			Platform: v["live_platform"],
			Uri:      "https://www.huya.com/"+v["live_uri"],
		}
		l = append(l, vo)

		str,_ := json.Marshal(vo)
		rconn.Do("SADD", RedisNotFollowOffSet, str)
	}
	return l,err
}

func GetOnline() (l []Queue, err error)  {
	rconn := redis.GetConn()
	rlist, err := rconn.Do("SMEMBERS", RedisOnlineSet)
	if err != nil {
		log.Println("Redis SMEMBERS error",RedisOnlineSet)
		return l, err
	}
	v := reflect.ValueOf(rlist)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	len := v.Len()
	if len != 0 {
		for i := 0; i < len; i++ {
			queue := Queue{}
			//json.Unmarshal(.([]byte), &queue)
			vo := Queue{
				LiveId: "",
				Platform: queue.Platform,
				Uri:      v.Index(i).String(),
			}
			l = append(l, vo)
		}
	}

	return l,err
}

//获取平台数据
func GetPlatforms() (p []Platform, err error) {
	if platforms != nil {
		return platforms,err
	}
	if res,err := mysql.Table("live_platform").Where("status=1").Get(); err == nil {
		for _, v := range res {
			PlatformId, _ := strconv.Atoi(v["platform_id"])
			vs := Platform{
				PlatformId: PlatformId,
				Mark:       v["mark"],
				Name:       v["name"],
				Domain:     v["domain"],
				PullUrl:    v["pull_url"],
				Status:     0,
				DomainUrl:  v["domain_url"],
			}
			p = append(p, vs)
		}
	}

	platforms = p
	return p,err
}
var Tests int
func Test()  {
	go func() {
		for true  {
			<-time.Tick(time.Second * 1)
			Tests++
		}
	}()
}