package main

import (
	"log"
)

func (spider *Spider) Parsers()  {
	for v := range spider.ChanParsers {
		switch v.Queue.LiveData.Live_platform {
		case "huya":
			spider.huYaParser(v)
		case "douyu":
			spider.douYuParser(v)
		case "kuaishou":
			spider.kuaiShouParser(v)
		default:
			log.Println("未知平台")
		}
	}
}