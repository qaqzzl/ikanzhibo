package main

import (
	"encoding/json"
	"fmt"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (spider *Spider) Downloader()  {
	go spider.downloaderFollowOffline()
	go spider.downloaderNotFollowOffline()
	go spider.downloaderOnline()
	go spider.downloaderTotalPlatform()
}

//被关注&&不在线 下载器
func (spider *Spider) downloaderFollowOffline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	var queue db.Queue
	for {
		v, err := rconn.Do("RPOP", db.RedisFollowOfflineList)
		if err != nil {
			log.Panicln(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			<-time.Tick(time.Second * 5)
			continue
		}
		body, err := downloaders(v, &queue)
		if err != nil {
			log.Println(err.Error())
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
	for {
		v, err := rconn.Do("RPOP", db.RedisNotFollowOfflineList)
		if err != nil {
			log.Panicln(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			<-time.Tick(time.Second * 5)
			continue
		}
		body, err := downloaders(v, &queue)
		if err != nil {
			log.Println(err.Error())
		}

		spider.ChanParsers <- &Parser{
			Body:body,
			Queue:queue,
		}

	}
}

//在线直播间
func (spider *Spider) downloaderOnline() {
	rconn := redis.GetConn()
	defer rconn.Close()
	var queue db.Queue
	for {
		v, err := rconn.Do("RPOP", db.RedisOnlineList)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			<-time.Tick(time.Second * 5)
			continue
		}
		body, err := downloaders(v, &queue)
		if err != nil {
			log.Println(err.Error())
			return
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
	var queue db.Queue
	for {
		v, err := rconn.Do("RPOP", db.RedisListList)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			<-time.Tick(time.Second * 5)
			continue
		}

		body, err := downloaders(v, &queue)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		spider.ChanParsers <- &Parser{
			Body:body,
			Queue:queue,
		}
	}
}

func downloaders(v interface{}, queue *db.Queue) (body []byte, err error)  {
	if err = json.Unmarshal(v.([]byte), &queue); err != nil {
		return body, err
	}


	client := &http.Client{}

	request, err := http.NewRequest("GET", queue.Uri, nil)
	if err != nil {
		return body, err
	}
	//增加header选项
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")

	response, err := client.Do(request)
	// 下面这句导致内存泄露  - 原因:资源还需要使用, 但是被close回收了
	defer response.Body.Close()


	body, err = ioutil.ReadAll(response.Body)
	fmt.Println(queue.Uri)
	return body, err
}