package main

import (
	"fmt"
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
	//test()

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

	client := &http.Client{}

	request, err := http.NewRequest("GET", "https://live.kuaishou.com/profile/Sanguo12138", nil)

	//增加header选项
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")

	response, err := client.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	Live := db.TableLive{}

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println("htmlquery ERR:" + err.Error())
		return
	}
	Live_title := htmlquery.FindOne(doc, "//a[@class='router-link-exact-active router-link-active live-card-following-info-title']")
	Live.Live_title = htmlquery.SelectAttr(Live_title, "title")
	fmt.Println(Live.Live_title)
	Live_anchortv_name := htmlquery.FindOne(doc, "//p[@class='user-info-name']")
	if Live_anchortv_name == nil {
		fmt.Println("昵称查找失败 \n")
		return
	}
	Live.Live_anchortv_name = htmlquery.InnerText(Live_anchortv_name)
	// 去除空格
	Live.Live_anchortv_name = strings.Replace(Live.Live_anchortv_name, " ", "", -1)
	// 去除换行符
	Live.Live_anchortv_name = strings.Replace(Live.Live_anchortv_name, "\n", "", -1)
	// 去除字符 (举报)
	Live.Live_anchortv_name = strings.Replace(Live.Live_anchortv_name, "举报", "", -1)
	fmt.Println(Live.Live_anchortv_name)

	//.Live_uri #
	Live.Live_uri = "p.Queue.Uri"

	//.Live_platform #
	Live.Live_platform = "kuaishou"

}