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
func (spider *Spider) EventOnlineNotice(w *WriteInfo, rconn redis.Conn)  {
	if w.TableLive.Live_is_online != "yes" {
		return
	}

	str,_ := json.Marshal(w)
	if _, err := rconn.Do("RPUSH", db.RedisOnlineNotice, str); err != nil {
		log.Println(err.Error())
	}
}


/**
 * 主播开播通知
*/
func (spider *Spider) SendOnlineNotice() {


}