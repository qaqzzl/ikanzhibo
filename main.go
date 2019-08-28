package main

import (
	"log"
)

func init() {
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


	Master()
}