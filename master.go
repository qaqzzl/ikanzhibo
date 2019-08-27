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
	go HandlerTotalPlatforms()

	//go HandlerFollowOffline()				//关注&&不在线直播间

	//go HandlerOnline()					//在线直播间
	for i := 0; i < 5; i++ {
		go Downloader()						//下载器
	}

	go parser.Parsers()					//解析器

	go UniqueList()						//去重器

	go WriteLiveInfo()					//写入器

	<-time.Tick(time.Second * 50000)
}

//被关注&&不在线
func HandlerFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	myChan := time.NewTicker(time.Second * 5) 	//抓取频率控制, 单位 / 秒
	for {
		<-myChan.C	//程序等待

		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisFollowOffLine); err == nil {
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
			//加入任务到redis队列
			if _, err := rconn.Do("RPUSH", db.RedisFollowOffLine, str); err != nil {
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

//全平台任务发现
func HandlerTotalPlatforms() {
	rconn := redis.GetConn()
	defer rconn.Close()
	myChan := time.NewTicker(time.Second * 10) 	//抓取频率控制, 单位 / 秒
	for {
		<- myChan.C
		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisListList); err == nil {
			if queueCounts.(int64) != 0 {
				continue
			}
			rconn.Do("del", db.RedisListOnceSet)	//清空已经爬取的任务
		} else {
			log.Println(err.Error())
			continue
		}

		//获取任务数据
		p, err := db.GetPlatforms()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		for _, v := range p {
			vs := db.Queue{
				Platform: v.Mark,
				Uri:      v.PullUrl,
				Type:     "live_list",
			}
			rconn.Do("SADD", db.RedisListOnceSet, vs.Uri)
			str,_ := json.Marshal(vs)
			//加入任务到redis队列
			if _, err := rconn.Do("RPUSH", db.RedisListList, str); err != nil {
				log.Println(err.Error())
			}

		}
	}

}