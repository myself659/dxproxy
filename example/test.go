package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	url := "http://tvp.daxiangdaili.com/ip/?tid=559836760324732&num=10000&operator=1&area=杭州&&protocol=http"
	resp, err := http.Get(url)
	fmt.Println(err)
	fmt.Println(resp.Body)
	// 如何打印一个网页
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(nil)
		return
	}
	// fmt.Println()
	// 从这里提取
	bs := string(b)
	abs := strings.Split(bs, "\r")
	fmt.Println(abs)
	fmt.Println(len(abs))

}
