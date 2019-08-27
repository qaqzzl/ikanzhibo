package main

import (
	"ikanzhibo/db"
	"ikanzhibo/db/mysql"
	"ikanzhibo/db/redis"
	"ikanzhibo/parser"
	"log"
	"strconv"
	"strings"
	"time"
)

func WriteLiveInfo()  {
	data := []*db.TableLive{}
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 30;
	for v := range parser.ChanProduceLiveInfo {
		//30秒 || 大于10 -> 更新数据
		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime <= currentTime || len(data) > 10 {
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