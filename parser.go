package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (spider *Spider) Parsers()  {
	for v := range spider.ChanParsers {
		switch v.Queue.QueueSet.Platform {
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


/**
 * 通过url获取uri
 *
 */
func urlGetUri(url string) string {
	regexpUrl := regexp.MustCompile(`[a-zA-z]+://[^\s]*/(.+)`)
	uris := regexpUrl.FindSubmatch([]byte(url))
	return string(uris[1])
}

/**
 * 过滤字符串 , 防止sql注入跟其他错误
 */
func liveReplaceSql(data string) (ret string) {
	ret = strings.Replace(data, "'", "\"", -1)
	ret = strings.Replace(ret, "\\", "\\\\", -1)
	//ret = strings.Replace(ret, " ", "_", -1)

	return ret
}


/**
 * 通过平台分类查找本地分类ID
 * @param string 平台分类
 * @return string 自定义分类ID string
 */
func platformTypeToLocal(live_class string) (live_type_id string,live_type_name string) {
	live_type_id = "0"
	live_type_name = ""
	for i := 0; i<len(LiveMyTypeData); i++ {
		if strings.Contains(LiveMyTypeData[i]["subset"],"#"+live_class+"#") {
			live_type_id = LiveMyTypeData[i]["type_id"]
			live_type_name = LiveMyTypeData[i]["name"]
			break
		}
	}
	return live_type_id,live_type_name
}

/**
 * 查询分类权重
*/
func getLocalTypeIdWeight(live_id string) (weight int, err error) {
	for i := 0; i<len(LiveMyTypeData); i++ {
		if LiveMyTypeData[i]["type_id"] == live_id {
			weight,err = strconv.Atoi(LiveMyTypeData[i]["weight"])
			break
		}
	}
	return weight,err
}

/**
 * 通过平台关注人数计算权重
*/
func platformFollowToWeight(followSum int, platform string) (weight int) {
	if platform == "huya" {
		platform_weight := 50;
		weight = (followSum * platform_weight / 100)
	}
	if platform == "douyu" {
		platform_weight := 100;
		weight = (followSum * platform_weight / 100)
	}
	return weight
}