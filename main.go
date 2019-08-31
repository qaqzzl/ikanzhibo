package main

import (
	"github.com/antchfx/htmlquery"
	"ikanzhibo/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

	spider := Spider{
		ChanParsers: make(chan *Parser, 1000),
		ChanProduceList: make(chan *db.Queue, 1000),
		ChanWriteInfo: make(chan *WriteInfo, 1000),
	}

	go Master(&spider)


	<-time.Tick(time.Second * 60000)
	//Monitor := Monitor{}
	//Monitor.Start(&spider)
}

func test()  {

	resp,_ := http.Get("https://live.kuaishou.com/profile/saoyu2002")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	Live := db.TableLive{}

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println("htmlquery ERR:" + err.Error())
		return
	}
	//.Live_is_online - 判断是在播
	Live_is_online := htmlquery.FindOne(doc, "//div[@class='live-card']")
	if Live_is_online == nil {
		Live.Live_is_online = "no"
	} else {
		Live.Live_is_online = "yes"
	}


	//.Live_uri #
	Live.Live_uri = "p.Queue.Uri"

	//.Live_platform #
	Live.Live_platform = "kuaishou"

}