package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"ikanzhibo/parser"
	"log"
	"time"
)

//调度器
func Master()  {
	go HandlerFollowOffline()			//关注&&不在线

	Downloader()						//下载器

	parser.Parsers()					//解析器

	<-time.Tick(time.Second * 50)
}

//被关注&&不在线
func HandlerFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	myChan := time.NewTicker(time.Second * 1) 	//单位 / 秒
	for {
		<-myChan.C	//程序等待

		//判断队列是否空(live_follow_offline_list)
		if queueCounts, err := rconn.Do("LLEN", db.RedisFollowOffline); err == nil {
			if queueCounts.(int64) != 0 {
				continue
			}
		} else {
			log.Println(err.Error())
			continue
		}
		//获取被关注过但不在线的直播间
		l, err := db.GetFollowOffline();
		if err != nil {
			log.Panicln(err)
		}

		for _, v := range l {
			str,_ := json.Marshal(v)
			//加入redis任务到队列
			if _, err := rconn.Do("RPUSH", db.RedisFollowOffline, str); err != nil {
				log.Println(err.Error())
			}

		}
	}

}

//无关注&&不在线
func HandlerNotFollowOffline() {

}

//在线
func HandlerOnline() {

}