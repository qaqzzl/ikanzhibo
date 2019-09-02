package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"log"
)

//发现任务列表去除重复
func (spider *Spider) UniqueList()  {
	rconn := redis.GetConn()
	defer rconn.Close()
	for v := range spider.ChanProduceList {
		switch v.QueueSet.QueueType {
		case "live_list":
			spider.uniqueLiveList(v, rconn)
		case "live_info":
			spider.uniqueLiveInfo(v, rconn)
		}
	}
}

func (spider *Spider) uniqueLiveList(v *db.Queue, rconn redis.Conn) {
	//加入已爬取集合(set) 如果存在会返回 0 ,加入成功返回 1
	set, err := rconn.Do("SADD", db.RedisListOnceSet, v.QueueSet.Request.Url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if set.(int64) == int64(0) { //存在
		return
	}
	//加入 等待爬取列表队列(list)
	str,_ := json.Marshal(v)
	_, err = rconn.Do("LPUSH", db.RedisListList, str)
	if err != nil { //写入错误
		log.Println(err.Error())
	}
}

func (spider *Spider) uniqueLiveInfo(v *db.Queue, rconn redis.Conn)  {
	setStr,_ := json.Marshal(v.QueueSet)
	//查看info总集合是否已经存在
	set, err := rconn.Do("SISMEMBER", db.RedisInfoOnceSet, setStr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if set.(int64) == int64(1) { //存在
		return
	}
	//加入 在线直播间队列(list)
	listStr,_ := json.Marshal(v)
	_, err = rconn.Do("LPUSH", db.RedisListList, listStr)
	if err != nil { //写入错误
		log.Println(err.Error())
	}
	//加入 在线直播间集合(set)
	rconn.Do("SADD", db.RedisOnlineSet, setStr)
	rconn.Do("SADD", db.RedisInfoOnceSet, setStr)
}