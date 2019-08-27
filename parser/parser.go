package parser

import (
	"ikanzhibo/db"
	"log"
)

//待解析数据 chan
type Parser struct {
	Body	[]byte
	Queue	db.Queue
}
var ChanParsers = make(chan *Parser, 1000)

//抓取列表 chan
var ChanProduceList = make(chan *db.Queue, 1000)

//直播间数据 chan
type ProduceLiveInfo struct {
	TableLive	db.TableLive
	Queue		db.Queue
}
var ChanProduceLiveInfo = make(chan *ProduceLiveInfo, 1000)


func Parsers()  {
	for v := range ChanParsers {
		switch v.Queue.Platform {
		case "huya":
			huYaParser(v)
		case "douyu":
			douYuParser(v)
		default:
			log.Println("未知平台")
		}
	}
}