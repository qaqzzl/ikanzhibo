package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/mysql"
	"ikanzhibo/db/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

func (spider *Spider) WriteLiveInfo()  {
	rconn := redis.GetConn()
	defer rconn.Close()

	offline_data := []*db.Queue{}
	online_data := []*db.Queue{}
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	offline_endTime := initTime + 30;	//控制更新数据库写入频率
	online_endTime := initTime + 30;	//控制更新数据库写入频率
	for v := range spider.ChanWriteInfo {
		//30秒 || 数据大于20 -> 更新数据 , 在线
		onlineCurrentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if online_endTime <= onlineCurrentTime || len(online_data) > 20 {
			if len(online_data) > 0 {
				writeOnlineLiveInfos(online_data)		//写入 code ...
			}
			//清空
			online_data = []*db.Queue{}
			//初始化时间
			onlineCurrentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
			online_endTime = onlineCurrentTime + 30
		}
		//30秒 || 数据大于20 -> 更新数据 , 离线
		offlineCurrentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if offline_endTime <= offlineCurrentTime || len(offline_data) > 20 {
			if len(offline_data) > 0 {
				writeOfflineLiveInfos(offline_data)		//写入 code ...
			}
			offline_data = []*db.Queue{}	//清空
			//初始化时间
			offlineCurrentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
			offline_endTime = offlineCurrentTime + 30
		}

		if v.LiveData.Live_is_online == "no" {	//离线
			offline_data = append(offline_data,v)
		} else {
			online_data = append(online_data,v)
		}

		//事件
		event := strings.Split(v.WriteEvent, ",")
		for i:=0; i<len(event); i++ {
			if event[i] == "online_notice" {		//发送开播通知
				spider.EventOnlineNotice(v, rconn)
			}
			if event[i] == "send_barrage" {			//发送弹幕
			}
			if event[i] == "listener_barrage" {		//监听弹幕
			}
		}

		//redis -> 在播添加 , 不在播删除
		setStr,_ := json.Marshal(v.QueueSet)
		switch v.LiveData.Live_is_online {
		case "yes":
			rconn.Do("SADD", db.RedisOnlineSet, setStr)		//被关注&&不在线直播间集合
			rconn.Do("SREM", db.RedisNotFollowOfflineSet, setStr)	//未关注&&不在线直播间集合
		case "no":
			rconn.Do("SREM", db.RedisOnlineSet, setStr)		//被关注&&不在线直播间集合
			rconn.Do("SADD", db.RedisNotFollowOfflineSet, setStr)	//未关注&&不在线直播间集合
		case "vio":
			rconn.Do("SREM", db.RedisOnlineSet, setStr)		//被关注&&不在线直播间集合
			rconn.Do("SREM", db.RedisNotFollowOfflineSet, setStr)	//未关注&&不在线直播间集合
		case "del":
			rconn.Do("SREM", db.RedisOnlineSet, setStr)		//被关注&&不在线直播间集合
			rconn.Do("SREM", db.RedisNotFollowOfflineSet, setStr)	//未关注&&不在线直播间集合
		}

	}
}

func writeOnlineLiveInfos(info []*db.Queue)  {
	//var mysqls string
	sql := "INSERT INTO `live` (live_title,live_anchortv_name,live_anchortv_photo,live_anchortv_sex,live_cover,live_play,live_class,live_tag,live_introduction," +
		"live_online_user,live_follow,live_uri,live_type_id,live_type_name,live_platform,live_is_online," +
		"spider_pull_url,platform_room_id,spider_pull_time,live_play_time,live_play_end_time,created_at,updated_at) VALUES "
	for i:=0; i<len(info); i++ {
		data := info[i].LiveData
		//sql
		sql += "('"+data.Live_title+"','" +
			data.Live_anchortv_name+"','" +
			data.Live_anchortv_photo+"','" +
			data.Live_anchortv_sex+"','" +
			data.Live_cover+"','" +
			data.Live_play+"','" +
			data.Live_class+"','" +
			data.Live_tag+"','" +
			data.Live_introduction+"','" +
			data.Live_online_user+"','" +
			data.Live_follow+"','" +
			data.Live_uri+"','" +
			data.Live_type_id+"','" +
			data.Live_type_name+"','" +
			data.Live_platform+"','" +
			data.Live_is_online+"','" +
			data.Spider_pull_url+"','" +
			data.Platform_room_id+"','" +
			data.Spider_pull_time+"','" +
			data.Live_play_time+"','" +
			data.Live_play_end_time+"','" +
			data.Created_at+"','" +
			data.Updated_at+
			"'),"
	}
	//sql
	sql = strings.Trim(sql,",")
	sql += " ON DUPLICATE KEY UPDATE " +
		"live_title=VALUES(live_title)," +
		"live_anchortv_name=VALUES(live_anchortv_name)," +
		"live_anchortv_photo=VALUES(live_anchortv_photo)," +
		"live_cover=VALUES(live_cover)," +
		"live_play=VALUES(live_play)," +
		"live_class=VALUES(live_class)," +
		"live_tag=VALUES(live_tag)," +
		"live_introduction=VALUES(live_introduction)," +
		"live_online_user=VALUES(live_online_user)," +
		"live_follow=VALUES(live_follow)," +
		"live_type_id=VALUES(live_type_id)," +
		"live_type_name=VALUES(live_type_name)," +
		"live_is_online=VALUES(live_is_online)," +
		"spider_pull_url=VALUES(spider_pull_url)," +
		"live_uri=VALUES(live_uri)," +
		"spider_pull_time=VALUES(spider_pull_time)," +
		"live_play_time=VALUES(live_play_time)," +
		"live_play_end_time=VALUES(live_play_end_time)," +
		"updated_at=VALUES(updated_at);"

	err := mysql.Conn().InsertSql(sql);

	if err != nil {
		log.Printf(err.Error())
	}
}

func writeOfflineLiveInfos(info []*db.Queue)  {
	//var mysqls string
	sql := "INSERT INTO `live` (live_title,live_anchortv_name,live_anchortv_photo,live_anchortv_sex,live_cover,live_play,live_class,live_tag,live_introduction," +
		"live_online_user,live_follow,live_uri,live_type_id,live_type_name,live_platform,live_is_online," +
		"spider_pull_url,platform_room_id,spider_pull_time,live_play_time,live_play_end_time,created_at,updated_at) VALUES "
	for i:=0; i<len(info); i++ {
		data := info[i].LiveData
		//sql
		sql += "('"+data.Live_title+"','" +
			data.Live_anchortv_name+"','" +
			data.Live_anchortv_photo+"','" +
			data.Live_anchortv_sex+"','" +
			data.Live_cover+"','" +
			data.Live_play+"','" +
			data.Live_class+"','" +
			data.Live_tag+"','" +
			data.Live_introduction+"','" +
			data.Live_online_user+"','" +
			data.Live_follow+"','" +
			data.Live_uri+"','" +
			data.Live_type_id+"','" +
			data.Live_type_name+"','" +
			data.Live_platform+"','" +
			data.Live_is_online+"','" +
			data.Spider_pull_url+"','" +
			data.Platform_room_id+"','" +
			data.Spider_pull_time+"','" +
			data.Live_play_time+"','" +
			data.Live_play_end_time+"','" +
			data.Created_at+"','" +
			data.Updated_at+
			"'),"
	}
	//sql
	sql = strings.Trim(sql,",")
	sql += " ON DUPLICATE KEY UPDATE " +
		"live_title=VALUES(live_title)," +
		//"live_introduction=VALUES(live_introduction)," +
		"live_follow=VALUES(live_follow)," +
		"live_is_online=VALUES(live_is_online)," +
		"spider_pull_url=VALUES(spider_pull_url)," +
		"live_uri=VALUES(live_uri)," +
		"spider_pull_time=VALUES(spider_pull_time)," +
		"live_play_time=VALUES(live_play_time)," +
		"live_play_end_time=VALUES(live_play_end_time)," +
		"updated_at=VALUES(updated_at);"

	err := mysql.Conn().InsertSql(sql);

	if err != nil {
		log.Printf(err.Error())
	}
}