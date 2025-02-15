package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"io/ioutil"
	"net/http"
	"time"
)

func (spider *Spider) Downloader()  {
	go spider.downloaderFollowOffline()
	go spider.downloaderTotalPlatform()
	go spider.downloaderNotFollowOffline()
	go spider.downloaderOnline()
	go spider.downloaderRecommend()
}

//被关注&&不在线 下载器
func (spider *Spider) downloaderFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	var queue db.Queue
	ticker := time.NewTicker(time.Second * 5)
	for {
		v, err := rconn.Do("LPOP", db.RedisFollowOfflineList)
		if err != nil {
			ErrLog(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			//<-time.Tick(time.Second * 5)
			<-ticker.C
			continue
		}
		body, err := downloaders(v, &queue)
		if body == nil {
			continue
		}
		if err != nil {
			ErrLog(err.Error())
			continue
		}

		spider.ChanParsers <- &Parser{
			Body:body,
			Queue:queue,
		}

	}
}

//未关注&&不在线 下载器
func (spider *Spider) downloaderNotFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	var queue db.Queue
	ticker := time.NewTicker(time.Second * 5)
	for {
		v, err := rconn.Do("LPOP", db.RedisNotFollowOfflineList)
		if err != nil {
			ErrLog(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			//<-time.Tick(time.Second * 5)
			<-ticker.C
			continue
		}
		body, err := downloaders(v, &queue)
		if body == nil {
			continue
		}
		if err != nil {
			ErrLog(err.Error())
			continue
		}

		spider.ChanParsers <- &Parser{
			Body:body,
			Queue: queue,
		}

	}
}

//在线直播间
func (spider *Spider) downloaderOnline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	var queue db.Queue
	ticker := time.NewTicker(time.Second * 5)
	for {
		v, err := rconn.Do("LPOP", db.RedisOnlineList)
		if err != nil {
			ErrLog(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			//<-time.Tick(time.Second * 5)
			<-ticker.C
			continue
		}
		body, err := downloaders(v, &queue)
		if body == nil {
			ErrLog("download body nil")
			continue
		}
		if err != nil {
			ErrLog(err.Error())
			continue
		}

		spider.ChanParsers <- &Parser{
			Body:body,
			Queue:queue,
		}

	}
}

//全平台任务发现下载
func (spider *Spider) downloaderTotalPlatform()  {
	rconn := redis.GetConn()
	defer rconn.Close()
	queue := db.Queue{}
	ticker := time.NewTicker(time.Second * 5)
	for {
		v, err := rconn.Do("LPOP", db.RedisListList)
		if err != nil {
			ErrLog(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			//<-time.Tick(time.Second * 5)
			<-ticker.C
			continue
		}

		body, err := downloaders(v, &queue)
		if err != nil {
			ErrLog(err.Error())
			continue
		}
		if body == nil {
			continue
		}
		spider.ChanParsers <- &Parser{
			Body:body,
			Queue:queue,
		}
	}
}

//推荐直播间
func (spider *Spider) downloaderRecommend()  {

}

//var client = &http.Client{}
func downloaders(v interface{}, queue *db.Queue) (body []byte, err error)  {
	if err = json.Unmarshal(v.([]byte), &queue); err != nil {
		return body, err
	}

	tr := http.Transport{DisableKeepAlives: true}
	client := http.Client{Transport: &tr}
	//client := &http.Client{}

	request, err := http.NewRequest("GET", queue.QueueSet.Request.Url, nil)

	TypeMonitorChan <- TypeRequestNum

	//request.Close = true

	if err != nil {
		return body, err
	}
	//增加header选项
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")

	response, err := client.Do(request)
	//response,err := http.Get( queue.QueueSet.Request.Url)
	if err != nil {
		ErrLog(err.Error()+"\n"+queue.QueueSet.Request.Url)
		return body, err
	}
	if response == nil {
		ErrLog("err response")
		return body, err
	}
	// 可以不回收 close , 因为在for里 所以资源还在被使用?
	//defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		return body, err
	}

	//fmt.Println(queue.QueueSet.Request.Url)
	return body, err
}