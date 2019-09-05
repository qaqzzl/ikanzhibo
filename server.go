package main

import (
	"encoding/json"
	"ikanzhibo/db"
	"ikanzhibo/db/redis"
	"io"
	"net/http"
	"time"
)

// 系统状态监控
type SystemMonitor struct {
	CurrentTime						int64	`json:"current_time"`
	TotalRequestNum   				int     `json:"total_request_num"`   		// 总处理请求数
	Tps          					float64 `json:"tps"`          				// 系统吞出量

	RedisOnlineList   				int64     `json:"redis_online_list"`   				// 在线队列数量
	RedisOnlineSet   				int64     `json:"redis_online_set"`   				// 在线集合数量
	RedisNotFollowOfflineList   	int64     `json:"redis_notFollow_offline_list"`   	// 关注不在线队列数量
	RedisNotFollowOfflineSet   		int64     `json:"redis_notFollow_offline_set"`   		// 关注不在线集合数量
	RedisFollowOfflineList   		int64     `json:"redis_follow_offline_list"`   		// 未关注不在线队列数量
	RedisFollowOfflineSet   		int64     `json:"redis_follow_offline_set"`   		// 未关注不在线集合数量
	RedisListList   				int64     `json:"redis_list_list"`   					// 发现任务队列数量
	RedisListOnceSet   				int64     `json:"redis_list_once_set"`   				// 发现任务集合数量
	RedisInfoOnceSet   				int64     `json:"redis_info_once_set"`   				// 开播通知队列数量

	UsedMemoryHuman          		string	`json:"used_memory_human"`          // 应用使用内存, 20.00M

	ChanParsersNum  				int     `json:"channel_parsers_num"`  		// 等待解析 channel 数量
	ChanProduceListNum 				int     `json:"channel_produce_list_num"` 	// 发现任务 channel 数量
	ChanWriteInfoNum 				int     `json:"channel_write_info_num"` 		// 写入数据 channel 数量

	UptimeInSeconds      			string  `json:"uptime_in_seconds"`      	// 运行时间
	ErrNum       					int     `json:"err_num"`       				// 错误数 , 数据库错误, 下载错误
	ExceptionNum       				int     `json:"exception_num"`       		// 异常数 , redis错误
	WarningNum       				int     `json:"warning_num"`       			// 警告数 , 解析错误
}

const (
	TypeRequestNum 		= 0		//请求数量
	TypeErrNum     		= 1		//错误数
	TypeExceptionNum   	= 2		//异常数
	TypeWarningNum     	= 3		//警告数
)

var TypeMonitorChan = make(chan int, 200)

type Monitor struct {
	StartTime time.Time
	Data      SystemMonitor
	TpsSli    []int
}

// 应用 处理数量|错误数量
func (m *Monitor) StatusRta(s *Spider) {
	go func() {
		rconn := redis.GetConn()
		defer rconn.Close()
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.Data.ErrNum += 1
			case TypeRequestNum:
				m.Data.TotalRequestNum += 1
			case TypeExceptionNum:
				m.Data.ExceptionNum += 1
			case TypeWarningNum:
				m.Data.WarningNum += 1
			}
		}
	}()

	// 应用 Tps
	//ticker := time.NewTicker(time.Second * 1)
	go func() { //协程
		for {
			<-time.Tick(time.Second * 1)
			//<-ticker.C
			//tps
			m.TpsSli = append(m.TpsSli, m.Data.TotalRequestNum)
			if len(m.TpsSli) > 2 {
				m.TpsSli = m.TpsSli[1:]
			}
		}
	}()

	go func() {
		rconn := redis.GetConn()
		defer rconn.Close()
		for {
			<-ticker.C
			RedisOnlineList,_ := rconn.Do("LLEN", db.RedisOnlineList)
			m.Data.RedisOnlineList = RedisOnlineList.(int64);					//在线队列数量

			RedisOnlineSet,_ := rconn.Do("SCARD", db.RedisOnlineSet)
			m.Data.RedisOnlineSet = RedisOnlineSet.(int64);					//在线集合数量

			RedisNotFollowOfflineList,_ := rconn.Do("LLEN", db.RedisNotFollowOfflineList)
			m.Data.RedisNotFollowOfflineList = RedisNotFollowOfflineList.(int64);		//关注不在线队列数量

			RedisNotFollowOfflineSet,_ := rconn.Do("SCARD", db.RedisNotFollowOfflineSet)
			m.Data.RedisNotFollowOfflineSet = RedisNotFollowOfflineSet.(int64);		//关注不在线集合数量

			RedisFollowOfflineList,_ := rconn.Do("LLEN", db.RedisFollowOfflineList)
			m.Data.RedisFollowOfflineList = RedisFollowOfflineList.(int64);			//未关注不在线队列数量

			RedisFollowOfflineSet,_ := rconn.Do("SCARD", db.RedisFollowOfflineSet)
			m.Data.RedisFollowOfflineSet = RedisFollowOfflineSet.(int64);				//未关注不在线集合数量

			RedisListList,_ := rconn.Do("LLEN", db.RedisListList)
			m.Data.RedisListList = RedisListList.(int64);					//发现任务队列数量

			RedisListOnceSet,_ := rconn.Do("SCARD", db.RedisListOnceSet)
			m.Data.RedisListOnceSet = RedisListOnceSet.(int64);				//发现任务集合数量

			RedisInfoOnceSet,_ := rconn.Do("SCARD", db.RedisInfoOnceSet)
			m.Data.RedisInfoOnceSet = RedisInfoOnceSet.(int64);				//开播通知队列数量

			m.Data.ChanParsersNum	= len(s.ChanParsers);				//等待解析 channel 数量
			m.Data.ChanProduceListNum = len(s.ChanProduceList);				//发现任务 channel 数量
			m.Data.ChanWriteInfoNum = len(s.ChanWriteInfo);				//入数据 channel 数量
		}

	}()
}


func (m *Monitor) Start(s *Spider) {
	go m.StatusRta(s)
	//http 服务 可以保持服务永久运行
	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		writer.Header().Set("content-type", "application/json")             //返回数据格式是json

		m.Data.UptimeInSeconds = time.Now().Sub(m.StartTime).String()
		if len(m.TpsSli) >= 2 {
			m.Data.Tps = float64(m.TpsSli[1]-m.TpsSli[0])
		}
		m.Data.CurrentTime = time.Now().Unix()
		ret, _ := json.MarshalIndent(m.Data, "", "\t")

		io.WriteString(writer, string(ret))
	})
	err := http.ListenAndServe(":1415", nil)
	if (  err != nil ) {
		panic(err.Error())
	}
}
