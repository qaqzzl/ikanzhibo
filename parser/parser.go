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

//直播数据 chan

func Parsers()  {
	for v := range ChanParsers {
		switch v.Queue.Platform {
		case "huya":
			huYaParser(v)
		case "douyu":
			douYuParser(v)
		default:
			log.Panicln("未知平台")
		}
	}
}