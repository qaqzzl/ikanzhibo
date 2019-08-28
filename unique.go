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
		switch v.Type {
		case "live_list":
			spider.uniqueLiveList(v, rconn)
		case "live_info":
			spider.uniqueLiveInfo(v, rconn)
		}
	}
}

func (spider *Spider) uniqueLiveList(v *db.Queue, rconn redis.Conn) {
	//加入已爬取集合(set) 如果存在会返回 0 ,加入成功返回 1
	set, err := rconn.Do("SADD", db.RedisListOnceSet, v.Uri)
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
	//查看在线直播间, 被关注&&不在线, 未关注&&不在线 集合是否存在
	OnlineSet, err := rconn.Do("SISMEMBER", db.RedisOnlineSet, v.Uri)
	if err != nil {
		log.Println(err.Error())
		return
	}
	FollowOffSet, err := rconn.Do("SISMEMBER", db.RedisFollowOffSet, v.Uri)
	if err != nil {
		log.Println(err.Error())
		return
	}
	NotFollowOffSet, err := rconn.Do("SISMEMBER", db.RedisNotFollowOffSet, v.Uri)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if FollowOffSet.(int64) == int64(1) || OnlineSet.(int64) == int64(1) || NotFollowOffSet.(int64) == int64(1) { //存在
		return
	}
	//加入 在线直播间队列(list)
	str,_ := json.Marshal(v)
	_, err = rconn.Do("LPUSH", db.RedisListList, str)
	if err != nil { //写入错误
		log.Println(err.Error())
	}
	//加入 在线直播间集合(set)
	rconn.Do("SADD", db.RedisOnlineSet, v.Uri)
}