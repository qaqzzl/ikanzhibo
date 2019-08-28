package main

import (
	"encoding/json"
	"fmt"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"log"
	"time"
)

//待解析数据结构体
type Parser struct {
	Body	[]byte
	Queue	db.Queue
}

//直播间数据 chan
type WriteInfo struct {
	TableLive	db.TableLive
	Queue		db.Queue
}

type SpiderInterface interface {
	handlerFollowOffline()
}

type Spider struct {
	ChanParsers 		chan *Parser
	ChanProduceList 	chan *db.Queue		//待抓取列表 chan
	WriteInfo			chan *WriteInfo		//吸入数据 chan
}

//调度器
func Master()  {
	spider := Spider{
		ChanParsers: make(chan *Parser, 1000),
		ChanProduceList: make(chan *db.Queue, 1000),
		WriteInfo: make(chan *WriteInfo, 1000),
	}
	go spider.handlerTotalPlatforms()				//发现任务

	//go HandlerFollowOffline()				//关注&&不在线直播间

	//go HandlerOnline()					//在线直播间
	for i := 0; i < 5; i++ {
		go spider.Downloader()						//下载器
	}

	go spider.Parsers() //解析器

	go spider.UniqueList()						//去重器

	go spider.WriteLiveInfo()					//写入器

	<-time.Tick(time.Second * 50000)
}

//被关注&&不在线
func (spider *Spider) handlerFollowOffline() {
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
func (Spider *Spider) handlerNotFollowOffline() {

}

//在线
func (Spider *Spider) handlerOnline() {

}

//全平台任务发现
func (spider *Spider) handlerTotalPlatforms() {
	rconn := redis.GetConn()
	defer rconn.Close()
	myChan := time.NewTicker(time.Second * 10) 	//抓取频率控制, 单位 / 秒
	i := 0
	for {
		<- myChan.C
		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisListList); err == nil {
			if queueCounts.(int64) != int64(0) {
				fmt.Println(queueCounts.(int64))
				fmt.Println("判断队列是否空")
				continue
			}
			fmt.Println("I",i)
			if i == 1 {
				log.Fatalln("结束")
				continue
			}
			i++
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