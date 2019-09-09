package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"log"
)

/**
  * 主播开播事件
*/
func (spider *Spider) EventOnlineNotice(q *db.Queue, rconn redis.Conn)  {
	if q.LiveData.Live_is_online != "yes" {
		return
	}

	setStr,_ := json.Marshal(q.QueueSet)
	set, err := rconn.Do("SISMEMBER", db.RedisFollowSet, setStr)

	if err != nil {
		log.Println(err.Error())
		return
	}

	if set.(int64) != int64(1) { //不存在
		return
	}

	str,_ := json.Marshal(q)
	if _, err := rconn.Do("RPUSH", db.RedisOnlineNoticeList, str); err != nil {
		log.Println(err.Error())
	}
}


/**
 * 主播开播通知
*/
func (spider *Spider) SendOnlineNotice() {


}