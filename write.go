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

	data := []*WriteInfo{}
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 30;	//控制更新数据库写入频率
	for v := range spider.ChanWriteInfo {
		//30秒 || 数据大于20 -> 更新数据
		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime <= currentTime || len(data) > 20 {
			//写入 code ...
			if len(data) > 0 {
				writeLiveInfos(data)
			}
			//清空
			data = []*WriteInfo{}

			//初始化时间
			currentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
			endTime = currentTime + 30
		} else {
			data = append(data,v)
		}

		//事件i
		event := strings.Split(v.Queue.Event, ",")
		for i:=0; i<len(event); i++ {
			if event[i] == "online_notice" {		//发送开播通知
				spider.EventOnlineNotice(v, rconn)
			}
			if event[i] == "send_barrage" {			//发送弹幕
			}
			if event[i] == "listener_barrage" {		//监听弹幕
			}
		}


	}
}

func writeLiveInfos(info []*WriteInfo)  {
	rconn := redis.GetConn()
	defer rconn.Close()
	//var mysqls string
	sql := "INSERT INTO `live` (live_title,live_anchortv_name,live_anchortv_photo,live_anchortv_sex,live_cover,live_play,live_class,live_tag,live_introduction," +
		"live_online_user,live_follow,live_uri,live_type_id,live_type_name,live_platform,live_is_online,created_at,updated_at) VALUES "
	for i:=0; i<len(info); i++ {
		data := info[i].TableLive
		queue := info[i].Queue
		//redis -> 在播添加 , 不在播删除
		str,_ := json.Marshal(queue)
		switch data.Live_is_online {
		case "yes":
			rconn.Do("SADD", db.RedisOnlineSet, str)		//被关注&&不在线直播间集合
			rconn.Do("SREM", db.RedisNotFollowOfflineSet, str)	//未关注&&不在线直播间集合
		case "no":
			rconn.Do("SREM", db.RedisOnlineSet, str)		//被关注&&不在线直播间集合
			rconn.Do("SADD", db.RedisNotFollowOfflineSet, str)	//未关注&&不在线直播间集合
		case "vio":
			rconn.Do("SREM", db.RedisOnlineSet, str)		//被关注&&不在线直播间集合
			rconn.Do("SREM", db.RedisNotFollowOfflineSet, str)	//未关注&&不在线直播间集合
		case "del":
			rconn.Do("SREM", db.RedisOnlineSet, str)		//被关注&&不在线直播间集合
			rconn.Do("SREM", db.RedisNotFollowOfflineSet, str)	//未关注&&不在线直播间集合
		}

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
			data.Live_pull_url+"','" +
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
		"live_anchortv_sex=VALUES(live_anchortv_sex)," +
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
		"live_pull_url=VALUES(live_pull_url)," +
		"updated_at=VALUES(updated_at);"

	err := mysql.Conn().InsertSql(sql);

	if err != nil {
		log.Printf(err.Error())
	}

}