package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"ikanzhibo/parser"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Downloader()  {
	go DownloaderFollowOffline()
}

//被关注&&不在线 下载器
func DownloaderFollowOffline()  {
	rconn := redis.GetConn()
	defer rconn.Close()
	var queue db.Queue
	for {
		v, err := rconn.Do("RPOP", db.RedisFollowOffline)
		if err != nil {
			log.Panicln(err.Error())
			continue
		}
		if v == nil {
			//暂停 5 秒
			<-time.Tick(time.Second * 5)
			continue
		}
		if err = json.Unmarshal(v.([]byte), &queue); err != nil {
			log.Println(err.Error())
			continue
		}
		resp,err := http.Get(queue.Uri)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		// 下面这句导致内存泄露  - 原因:资源还需要使用, 但是被close回收了
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		parser.ChanParsers <- &parser.Parser{
			Body:body,
			Queue:queue,
		}

	}
}