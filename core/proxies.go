package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kgretzky/evilginx2/log"
	"h12.io/socks"
)

func ReadProxyList(path string) []string {
	resp_chan := make(chan QueryResp, 10)
	dat, err := ioutil.ReadFile(path + "/proxies.txt")
	if err != nil {
		log.Fatal("read proxy list: %v", err)
	}

	dats := strings.Split(strings.TrimSuffix(string(dat), "\n"), "\n")
	wg := sync.WaitGroup{}
	runtime.GOMAXPROCS(4)

	for _, addr := range dats {
		go query(addr, resp_chan, &wg)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	close(resp_chan)

	var proxies []string

	for r := range resp_chan {
		if r.Err == nil {
			proxies = append(proxies, r.Addr)
		}
	}
	return proxies
}

type QueryResp struct {
	Addr string
	Err  error
}

type IP struct {
	IP string
}

func query(host string, c chan QueryResp, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	var timeout = time.Duration(5 * time.Second)
	httpClient := &http.Client{Timeout: timeout}
	url_proxy, err := url.Parse(host)
	if err != nil {
		log.Error("couldn't parse proxy address")
		c <- QueryResp{Addr: host, Err: err}
		return
	}
	var ip IP
	var resp *http.Response
	if url_proxy.Scheme == "socks5" {

		dialSocksProxy := socks.Dial(host)
		tr := &http.Transport{Dial: dialSocksProxy}
		httpClient.Transport = tr

		resp, err = httpClient.Get("http://api.ipify.org?format=json")
		if err != nil {
			c <- QueryResp{Addr: host, Err: err}
			return
		}
	} else {
		req, err := http.NewRequest("GET", "https://api.ipify.org/?format=json", nil)
		if err != nil {
			c <- QueryResp{Addr: host, Err: err}
			return
		}
		httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(url_proxy), ProxyConnectHeader: req.Header}

		resp, err = httpClient.Do(req)
		if err != nil {
			c <- QueryResp{Addr: host, Err: err}
			return
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- QueryResp{Addr: host, Err: err}
		return
	}
	defer resp.Body.Close()

	json.Unmarshal([]byte(body), &ip)
	sp := strings.Split(url_proxy.Host, ":")
	respIp := sp[0]

	if ip.IP == respIp {
		log.Info("finished checking proxy: %v", host)
		c <- QueryResp{Addr: host, Err: nil}
		return
	}

	c <- QueryResp{Addr: host, Err: fmt.Errorf("failed checking proxy, unknown error")}
}
