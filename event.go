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

	str,_ := json.Marshal(q)
	if _, err := rconn.Do("RPUSH", db.RedisOnlineNotice, str); err != nil {
		log.Println(err.Error())
	}
}


/**
 * 主播开播通知
*/
func (spider *Spider) SendOnlineNotice() {


}