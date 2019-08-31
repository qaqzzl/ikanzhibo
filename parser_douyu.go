package main

import (
	"fmt"
)

func douYuParser(p *Parser)  {
	switch p.Queue.Type {
	case "live_info":
		douYuLiveInfo(p)
	case "live_list":
		douYuLiveInfo(p)
	}
}

func douYuLiveInfo(p *Parser) (l interface{}, err error) {
	fmt.Println("进入huya")
	return l,err
}
