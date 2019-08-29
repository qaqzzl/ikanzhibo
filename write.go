package main

import (
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

	data := []*db.TableLive{}
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 30;	//控制更新数据库写入频率
	for v := range spider.WriteInfo {
		//30秒 || 数据大于20 -> 更新数据
		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime <= currentTime || len(data) > 20 {
			//写入 code ...
			if len(data) > 0 {
				writeLiveInfos(data)
			}
			//清空
			data = []*db.TableLive{}

			//初始化时间
			currentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
			endTime = currentTime + 30
		} else {
			data = append(data,&v.TableLive)
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

func writeLiveInfos(data []*db.TableLive)  {
	rconn := redis.GetConn()
	defer rconn.Close()
	//var mysqls string
	sql := "INSERT INTO `live` (live_title,live_anchortv_name,live_anchortv_photo,live_anchortv_sex,live_cover,live_play,live_class,live_tag,live_introduction," +
		"live_online_user,live_follow,live_uri,live_type_id,live_type_name,live_platform,live_is_online,created_at,updated_at) VALUES "
	for i:=0; i<len(data); i++ {
		//redis -> 在播添加 , 不在播删除
		if data[i].Live_is_online == "yes" {
			rconn.Do("SADD", db.RedisListOnceSet, data[i].Live_pull_url)
		} else { //删除
			rconn.Do("SREM", db.RedisListOnceSet, data[i].Live_pull_url)
		}

		//sql
		sql += "('"+data[i].Live_title+"','" +
			data[i].Live_anchortv_name+"','" +
			data[i].Live_anchortv_photo+"','" +
			data[i].Live_anchortv_sex+"','" +
			data[i].Live_cover+"','" +
			data[i].Live_play+"','" +
			data[i].Live_class+"','" +
			data[i].Live_tag+"','" +
			data[i].Live_introduction+"','" +
			data[i].Live_online_user+"','" +
			data[i].Live_follow+"','" +
			data[i].Live_uri+"','" +
			data[i].Live_type_id+"','" +
			data[i].Live_type_name+"','" +
			data[i].Live_platform+"','" +
			data[i].Live_is_online+"','" +
			data[i].Created_at+"','" +
			data[i].Updated_at+
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
		"updated_at=VALUES(updated_at);"

	err := mysql.Conn().InsertSql(sql);

	if err != nil {
		log.Printf(err.Error())
	}

}