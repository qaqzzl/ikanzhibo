package main

import (
	"ikanzhibo/db"
	"log"
	"os"
)

func init() {
	file := "./" +"message"+ ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}



func main() {
	//db.Test()
	//
	//go func() {
	//	for true  {
	//		<-time.Tick(time.Second * 1)
	//		fmt.Println(db.Tests)
	//	}
	//}()
	//<-time.Tick(time.Second * 50000)

	spider := Spider{
		ChanParsers: make(chan *Parser, 1000),
		ChanProduceList: make(chan *db.Queue, 1000),
		ChanWriteInfo: make(chan *WriteInfo, 1000),
	}

	Master(&spider)

	Monitor := Monitor{}
	Monitor.Start(&spider)
}