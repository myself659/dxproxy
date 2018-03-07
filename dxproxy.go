package dxproxy

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//http request

/*
User-Agent:Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36

UserAgent = [
'Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0)',
'Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.2)',
'Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)',
'Mozilla/5.0 (Windows; U; Windows NT 5.2) Gecko/2008070208 Firefox/3.0.1',
'Mozilla/5.0 (Windows; U; Windows NT 5.1) Gecko/20070803 Firefox/1.5.0.12',
'Mozilla/5.0 (Macintosh; PPC Mac OS X; U; en) Opera 8.0',
'Opera/8.0 (Macintosh; PPC Mac OS X; U; en)',
'Opera/9.27 (Windows NT 5.2; U; zh-cn)',
'Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/525.13 (KHTML, like Gecko) Chrome/0.2.149.27 Safari/525.13',
'Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.12) Gecko/20080219 Firefox/2.0.0.12 Navigator/9.0.0.6',
'Mozilla/5.0 (iPhone; U; CPU like Mac OS X) AppleWebKit/420.1 (KHTML, like Gecko) Version/3.0 Mobile/4A93 Safari/419.3',
'Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/525.13 (KHTML, like Gecko) Version/3.1 Safari/525.13'
]


*/
/*
透明代理
代理的策略 边测边验证的方式先验证再循环的循环模式
*/

var uas = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	"Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/525.13 (KHTML, like Gecko) Version/3.1 Safari/525.13",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1",
}

var ualen int = len(uas)

type proxyInfo struct {
	url string
	ua  string
}

type ProxyPool struct {
	urls     []string
	addurls  chan []string
	delurl   chan string
	getproxy chan *proxyInfo
	cur      *proxyInfo
	limit    int
}

func genProxyUrl(kind, ip, port string) string {
	if strings.Contains(kind, "https") || strings.Contains(kind, "HTTPS") {
		return "https://" + ip + ":" + port
	}

	return "http://" + ip + ":" + port
}

func NewProxyPool() *ProxyPool {
	pool := new(ProxyPool)

	pool.getproxy = make(chan *proxyInfo, 32)
	pool.urls = make([]string, 0)
	pool.addurls = make(chan []string)
	pool.delurl = make(chan string)
	pool.limit = 0
	pool.run()

	return pool
}

func (self *ProxyPool) fetchdx() {
	urls := []string{"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&operator=1&area=天津&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&operator=1&area=北京&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&operator=1&area=北京&protocol=http", "http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&operator=1&area=河北&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=30&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=20000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=20000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=20000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=20000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=20000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=20000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=2000&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=200&protocol=https"}
	url := urls[rand.Intn(len(urls))]
	resp, err := http.Get(url)
	if err == nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("fetchdx:", err)
			return
		}
		as := strings.Split(string(b), "\r")
		if len(as) < 1 {
			return
		}
		proxyurls := make([]string, len(as))
		for k, item := range as {
			httpsurl := "https://" + strings.TrimSpace(strings.Trim(item, "\r"))
			//fmt.Println("url:", httpsurl)
			proxyurls[k] = httpsurl
			//fmt.Println("url:", httpurl)
		}

		self.addurls <- proxyurls
		//fmt.Println(proxyurls)

	} else {
		fmt.Println("fetchdx:", err)
	}

}

func (self *ProxyPool) fetchdxhttp() {
	urls := []string{"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=1000&operator=1&area=天津&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&operator=1&area=北京&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=1000&operator=1&area=北京&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&operator=1&area=河北&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=3&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=3&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=30&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=3&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=3&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=3&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http",
		"http://tvp.daxiangdaili.com/ip/?tid=555783631065241&num=100&protocol=http"}
	url := urls[rand.Intn(len(urls))]
	resp, err := http.Get(url)
	if err == nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("fetchdx:", err)
			return
		}
		as := strings.Split(string(b), "\r")
		if len(as) < 1 {
			return
		}
		proxyurls := make([]string, len(as))
		for k, item := range as {
			httpsurl := "http://" + strings.TrimSpace(strings.Trim(item, "\r"))
			//fmt.Println("url:", httpsurl)
			proxyurls[k] = httpsurl
			//fmt.Println("url:", httpurl)
		}

		self.addurls <- proxyurls
		//fmt.Println(proxyurls)

	} else {
		fmt.Println("fetchdx:", err)
	}

}

func (self *ProxyPool) fetch() {
	var url string
	rn := rand.Intn(4)
	switch rn {
	case 0:
		{
			url = "http://www.xicidaili.com/nn"
		}
	case 1:
		{
			url = "http://www.xicidaili.com/nn"
		}
	case 2:
		{
			url = "http://www.xicidaili.com/nn/2"
		}
	case 3:
		{
			url = "http://www.xicidaili.com/nn/3"
		}
	case 4:
		{
			url = "http://www.xicidaili.com/nn/4"
		}
	}

	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	request := gorequest.New()
	resp, _, errs := request.Get(url).Set("User-Agent", ua).End()
	if errs != nil {
		fmt.Println(errs)
		return
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
	 <table id="ip_list">
	*/
	tables := doc.Find("table")
	tableslen := tables.Length()
	//fmt.Println(tableslen)
	for i := 0; i < tableslen; i++ {
		table := tables.Eq(i)
		_, ok := table.Attr("id")
		//fmt.Println(idattr, ok)
		if ok == true {
			trs := table.Find("tr")
			trslen := trs.Length()
			urls := make([]string, trslen)
			for i := 1; i < trslen; i++ {
				tr := trs.Eq(i)
				tds := tr.Find("td")
				kind := tds.Eq(5).Text()
				ip := tds.Eq(1).Text()
				port := tds.Eq(2).Text()
				j := i - 1
				urls[j] = genProxyUrl(kind, ip, port)
			}
			fmt.Println("send urls")
			//fmt.Println(urls)
			self.addurls <- urls
			break

		}
	}
}

func (self *ProxyPool) run() {

	go func() {

		var temp proxyInfo
		for {

			if len(self.urls) < 50 {
				go self.fetchdx()
				go self.fetchdxhttp()

				<-time.After(10 * time.Second)
				// continue
			}

			if len(self.urls) > 0 {
				urli := rand.Intn(len(self.urls))
				uai := rand.Intn(ualen)
				temp.url = self.urls[urli]
				temp.ua = uas[uai]
			}

			select {

			case urls := <-self.addurls:
				{

					fmt.Println("recv addurls:", len(urls))
					/*
						fmt.Println(urls)
					*/
					for _, url := range urls {
						self.urls = append(self.urls, url)
					}
					//fmt.Println(self.urls)
				}
			case url := <-self.delurl:
				{
					// del invalid  url
					for i := 0; i < len(self.urls); i++ {
						if self.urls[i] == url {
							copy(self.urls[i:], self.urls[i+1:])
							self.urls = self.urls[:len(self.urls)-1]
							//self.urls = append(self.urls[:i], self.urls[i+1:])
						}
					}

				}
			case self.getproxy <- &temp:
				{
					//do nothing
					//fmt.Println("send:", temp)
				}
			case <-time.After(100 * time.Second):
				{
					// time out
					fmt.Println("urls len:", len(self.urls))
				}
			}
		}
	}()

}

func (self *ProxyPool) Get(url string) (*http.Response, error) {
	//return http.Get(url)

	var proxyurl string
	var uas string
	for {
		if self.limit <= 0 {
			pinfo := <-self.getproxy
			proxyurl = pinfo.url
			if "" == proxyurl {
				self.delurl <- proxyurl
				continue
			}
			if strings.Contains(proxyurl, "ERROR") {
				self.delurl <- proxyurl
				continue
			}
			fmt.Println("proxyurl:", proxyurl, url)
			uas = pinfo.ua
			self.cur = pinfo
			self.limit = 60
		} else {
			proxyurl = self.cur.url
			uas = self.cur.ua
		}
		//fmt.Println(pinfo)

		// 后续可以考虑Proxy复用，这样对GC友好
		// 时间减少,快速处理 提高速度 现在应该不缺ip
		request := gorequest.New().Proxy(proxyurl).Set("User-Agent", uas).Timeout(10 * time.Second)
		resp, _, errs := request.Get(url).End()
		if errs == nil {
			self.limit--
			fmt.Println(url, "done")
			return resp, nil
		}
		fmt.Println(url, errs)
		// 有效提高速度,一定要将害群之马清除
		self.delurl <- proxyurl
		//fmt.Println(resp)
		self.limit = 0

	}
}
