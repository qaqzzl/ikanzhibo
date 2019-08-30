package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// 系统状态监控
type SystemMonitor struct {
	TotalRequestNum   				int     `json:"total_request_num"`   		// 总处理请求数
	Tps          					float64 `json:"tps"`          				// 系统吞出量

	RedisOnlineList   				int     `json:"redis_online_list"`   				// 在线队列数量
	RedisOnlineSet   				int     `json:"redis_online_set"`   				// 在线集合数量
	RedisNotFollowOfflineList   	int     `json:"redis_notFollow_offline_list"`   	// 关注不在线队列数量
	RedisNotFollowOfflineSet   		int     `json:"redis_notFollow_offline_set"`   		// 关注不在线集合数量
	RedisFollowOfflineList   		int     `json:"redis_follow_offline_list"`   		// 未关注不在线队列数量
	RedisFollowOffSet   			int     `json:"redis_follow_offline_set"`   		// 未关注不在线集合数量
	RedisListList   				int     `json:"redis_list_list"`   					// 发现任务队列数量
	RedisListOnceSet   				int     `json:"redis_list_once_set"`   				// 发现任务集合数量
	RedisInfoOnceSet   				int     `json:"redis_info_once_set"`   				// 开播通知队列数量

	UsedMemoryHuman          		string	`json:"used_memory_human"`          // 应用使用内存, 20.00M

	ChanParsersNum  				int     `json:"channel_parsers_num"`  		// 等待解析 channel 数量
	ChanProduceListNum 				int     `json:"channel_produceList_num"` 	// 发现任务 channel 数量
	ChanWriteInfoNum 				int     `json:"channel_writeInfo_num"` 		// 写入数据 channel 数量

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
	ticker := time.NewTicker(time.Second * 1)
	go func() { //协程
		for {
			<-ticker.C
			//tps
			m.TpsSli = append(m.TpsSli, m.Data.TotalRequestNum)
			if len(m.TpsSli) > 2 {
				m.TpsSli = m.TpsSli[1:]
			}
		}
		for {
			//在线队列数量
			m.Data.RedisOnlineList = 0;					//在线队列数量
			m.Data.RedisOnlineSet = 0;					//在线集合数量
			m.Data.RedisNotFollowOfflineList = 0;		//关注不在线队列数量
			m.Data.RedisNotFollowOfflineSet = 0;		//关注不在线集合数量
			m.Data.RedisFollowOfflineList = 0;			//未关注不在线队列数量
			m.Data.RedisFollowOffSet = 0;				//未关注不在线集合数量
			m.Data.RedisListList = 0;					//发现任务队列数量
			m.Data.RedisListOnceSet = 0;				//发现任务集合数量
			m.Data.RedisInfoOnceSet = 0;				//开播通知队列数量
		}

	}()
}


func (m *Monitor) Start(s *Spider) {
	m.StatusRta(s)
	//http 服务 可以保持服务永久运行
	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		writer.Header().Set("content-type", "application/json")             //返回数据格式是json

		m.Data.RunTime = time.Now().Sub(m.StartTime).String()
		if len(m.TpsSli) >= 2 {
			m.Data.Tps = float64(m.TpsSli[1]-m.TpsSli[0])
		}
		ret, _ := json.MarshalIndent(m.Data, "", "\t")

		io.WriteString(writer, string(ret))
	})

	http.ListenAndServe(":1415", nil)
}
