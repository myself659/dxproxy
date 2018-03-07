package main

import (
	"fmt"
	"github.com/myself659/dxproxy"
	"time"
)

func main() {
	xcproxy := dxproxy.NewProxyPool()
	<-time.After(5 * time.Second)
	url := "https://www.baidu.com"
	resp, err := xcproxy.Get(url)
	fmt.Println(resp)

	fmt.Println(err)
	<-time.After(60 * time.Second)
}
