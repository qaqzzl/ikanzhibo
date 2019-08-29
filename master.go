package main

import (
	"encoding/json"
	"fmt"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"log"
	"strconv"
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

	go spider.handlerFollowOffline()				//关注&&不在线直播间

	//go HandlerOnline()					//在线直播间
	for i := 0; i < 5; i++ {
		go spider.Downloader()						//下载器
	}

	go spider.Parsers() //解析器

	go spider.UniqueList()						//去重器

	go spider.WriteLiveInfo()					//写入器

	go Crontab()

	<-time.Tick(time.Second * 50000)
}

//被关注&&不在线
func (spider *Spider) handlerFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	//控制抓取频率
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 5;
	for {
		<-time.Tick(time.Second * 1)		//暂停, 单位 / 秒

		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime > currentTime {
			continue
		}

		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisFollowOffLine); err == nil {
			if queueCounts.(int64) != 0 {
				continue
			}
		} else {
			log.Println(err.Error())
			continue
		}
		//初始化抓取频率时间
		currentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		endTime = currentTime + 300;	//300秒-> 5分

		//获取被关注过但不在线的直播间
		l, err := db.GetFollowOffline();
		if err != nil {
			log.Panicln(err)
		}

		for _, v := range l {
			v.Type = "live_info"
			v.Event = "online_notice"
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
	rconn := redis.GetConn()
	defer rconn.Close()
	//控制抓取频率
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 5;
	for {
		<-time.Tick(time.Second * 1)		//暂停, 单位 / 秒

		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime > currentTime {
			continue
		}

		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisOnlineList); err == nil {
			if queueCounts.(int64) != 0 {
				continue
			}
		} else {
			log.Println(err.Error())
			continue
		}
		//初始化抓取频率时间
		currentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		endTime = currentTime + 600;	//600秒-> 10分

		//获取在线直播间
		l, err := db.GetOnline();
		if err != nil {
			log.Panicln(err)
		}

		for _, v := range l {
			v.Type = "live_info"
			v.Event = "online_notice"
			str,_ := json.Marshal(v)
			//加入任务到redis队列
			if _, err := rconn.Do("RPUSH", db.RedisFollowOffLine, str); err != nil {
				log.Println(err.Error())
			}

		}
	}
}

//全平台任务发现
func (spider *Spider) handlerTotalPlatforms() {
	rconn := redis.GetConn()
	defer rconn.Close()

	//控制抓取频率
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 5;
	for {
		<-time.Tick(time.Second * 1)	//暂停, 单位 / 秒

		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime > currentTime {
			continue
		}
		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisListList); err == nil {
			if queueCounts.(int64) != int64(0) {
				fmt.Println(queueCounts.(int64))
				fmt.Println("判断队列是否空")
				continue
			}
		} else {
			log.Println(err.Error())
			continue
		}
		rconn.Do("del", db.RedisListOnceSet)	//清空已经爬取的任务
		//初始化时间
		currentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		endTime = currentTime + 1800;	//1800秒-> 30分

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


//定时任务
func Crontab()  {
	rconn := redis.GetConn()
	defer rconn.Close()
	go func() {
		for true {
			<-time.Tick(time.Second * 60)	//60秒清除一次 被关注&&不在线直播间集合,定时跟数据库做一致性同步

			rconn.Do("del", db.RedisModelFollowOffSet)	//清空
		}
	}()
}