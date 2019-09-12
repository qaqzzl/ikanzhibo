package main

import (
	"encoding/json"
	"fmt"
	"ikanzhibo/db"
	"ikanzhibo/db/mysql"
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

type SpiderInterface interface {
	handlerFollowOffline()
}

type Spider struct {
	ChanParsers 		chan *Parser		//解析任务 chan
	ChanProduceList 	chan *db.Queue		//发现任务列表 chan
	ChanWriteInfo		chan *db.Queue		//写入数据 chan
}

//调度器
func Master(spider *Spider)  {
	//初始化数据
	InitLive()

	go spider.handlerTotalPlatforms()				//发现任务

	go spider.handlerFollowOffline()				//关注&&不在线直播间

	go spider.handlerNotFollowOffline()				//未关注&&不在线直播间

	go spider.handlerOnline()						//在线直播间

	for i:=0; i<5; i++ {
		go spider.Downloader()						//下载器
	}

	go spider.Parsers() //解析器

	go spider.UniqueList()						//去重器

	go spider.WriteLiveInfo()					//写入器

	go Crontab()
}

//被关注&&不在线
func (spider *Spider) handlerFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	//控制抓取频率
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 5;
	ticker := time.NewTicker(time.Second * 2)
	for {
		<-ticker.C
		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime > currentTime {
			continue
		}

		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisFollowOfflineList); err == nil {
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
			v.WriteEvent = "online_notice"
			str,_ := json.Marshal(v)
			//加入任务到redis队列
			if _, err := rconn.Do("RPUSH", db.RedisFollowOfflineList, str); err != nil {
				log.Println(err.Error())
			}
		}
	}

}

//无关注&&不在线
func (Spider *Spider) handlerNotFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	//控制抓取频率
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 5;
	ticker := time.NewTicker(time.Second * 2)
	for {
		<-ticker.C

		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime > currentTime {
			continue
		}

		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisNotFollowOfflineList); err == nil {
			if queueCounts.(int64) != 0 {
				continue
			}
		} else {
			log.Println(err.Error())
			continue
		}
		//初始化抓取频率时间
		currentTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		endTime = currentTime + 3000;	//300秒-> 5分

		//获取被关注过但不在线的直播间
		l, err := db.GetNotFollowOffline();
		if err != nil {
			log.Panicln(err)
		}

		for _, v := range l {
			v.WriteEvent = "online_notice"
			str,_ := json.Marshal(v)
			//加入任务到redis队列
			if _, err := rconn.Do("RPUSH", db.RedisNotFollowOfflineList, str); err != nil {
				log.Println(err.Error())
			}

		}
	}
}

//在线
func (Spider *Spider) handlerOnline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	//控制抓取频率
	initTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	endTime := initTime + 5;
	ticker := time.NewTicker(time.Second * 2)
	for {
		//<-time.Tick(time.Second * 1)		//暂停, 单位 / 秒
		<-ticker.C
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
			log.Println(err)
			continue
		}

		for _, v := range l {
			str,_ := json.Marshal(v)
			//加入任务到redis队列
			if _, err := rconn.Do("RPUSH", db.RedisOnlineList, str); err != nil {
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
	ticker := time.NewTicker(time.Second * 2)
	for {
		//<-time.Tick(time.Second * 1)	//暂停, 单位 / 秒
		<-ticker.C
		currentTime, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if endTime > currentTime {
			continue
		}
		//判断队列是否空
		if queueCounts, err := rconn.Do("LLEN", db.RedisListList); err == nil {
			if queueCounts.(int64) != int64(0) {
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
			val := db.Queue{
				QueueSet:db.QueueSet{
					Request:     db.Request{
						Url: v.PullUrl,
					},
					QueueType: "live_list",
					Platform: v.Mark,
				},
			}
			rconn.Do("SADD", db.RedisListOnceSet, v.PullUrl)
			str,_ := json.Marshal(val)
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
	go func() {		//关注 && 不在线
		ticker := time.NewTicker(time.Second * 60)
		for true {
			//<-time.Tick(time.Second * 60)	//60秒清除一次 被关注&&不在线直播间集合,定时跟数据库做一致性同步
			<-ticker.C
			rconn.Do("del", db.RedisFollowOfflineSet)	//清空
		}
	}()

	go func() {		//未关注 && 不在线
		newtime := time.Now().Unix()
		//获取本地location
		timeStr := time.Now().Format("2006-01-02")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 03:00:00", time.Local)
		time3hour := t.Unix() + 86400
		tickertime := time3hour - newtime
		ticker := time.NewTicker(time.Second * time.Duration(tickertime))
		for true {
			<-ticker.C
			rconn.Do("del", db.RedisNotFollowOfflineSet)	//清空
			db.NotFollowOfflineEmpty = 0
			db.GetNotFollowOffline()	//初始化直播未关注,不在播数据
		}
	}()

	go func() {		//在线 3点同步 , 策略,删除并读取数据库在线直播间进行同步
		newtime := time.Now().Unix()
		//获取本地location
		timeStr := time.Now().Format("2006-01-02")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 03:00:00", time.Local)
		time3hour := t.Unix() + 86400
		tickertime := time3hour - newtime
		ticker := time.NewTicker(time.Second * time.Duration(tickertime))
		for true {
			<-ticker.C
			db.OnlineEmpty = 0
			rconn.Do("del", db.RedisOnlineSet)	//清空
		}
	}()

	go func() {		//全任务发现 3点清除
		newtime := time.Now().Unix()
		//获取本地location
		timeStr := time.Now().Format("2006-01-02")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 03:00:00", time.Local)
		time3hour := t.Unix() + 86400
		tickertime := time3hour - newtime
		ticker := time.NewTicker(time.Second * time.Duration(tickertime))
		for true {
			<-ticker.C
			rconn.Do("del", db.RedisInfoOnceSet)	//清空
			rconn.Do("del", db.RedisListOnceSet)	//清空
		}
	}()
}

func InitLive() {
	initLiveMyType()		//初始化分类数据
	initLiveFollow()		//初始化直播关注数据
}

//初始化分类映射数据
var	LiveMyTypeData []map[string]string
func initLiveMyType() (err error) {
	fmt.Println("init type data")
	if LiveMyTypeData, err = mysql.Table("live_type").Select("type_id,name,subset,weight,weight_addition").Order("`order` asc").Get(); err != nil {
		panic("初始化失败 . 分类映射数据出错 ,"+err.Error())
	}
	return err
}


//被关注主播数据
func initLiveFollow() (err error) {
	rconn := redis.GetConn()
	defer rconn.Close()
	fmt.Println("init live follow data")

	list, err := mysql.Conn().QueryAll("select l.live_id,l.spider_pull_url,l.live_platform from live as l JOIN live_user_follow as luf ON luf.live_id = l.live_id WHERE luf.`status`=1 AND luf.is_notice=1")
	var l []db.Queue
	for _, v := range list {
		vo := db.Queue{
			QueueSet: db.QueueSet{
				Request:     db.Request{
					Url: v["spider_pull_url"],
				},
				QueueType: "live_info",
				Platform: v["live_platform"],
			},
		}
		l = append(l, vo)

		str,_ := json.Marshal(vo.QueueSet)
		rconn.Do("SADD", db.RedisFollowSet, str)
	}

	return err
}

