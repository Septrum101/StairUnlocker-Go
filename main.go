package main

import (
	"StairUnlocker-Go/config"
	"StairUnlocker-Go/utils"
	"flag"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

var (
	proxiesList []C.Proxy
	su          *config.SuConfig
	ver         bool
	help        bool
	daemon      bool
	configFile  string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&ver, "v", false, "show current ver of StairUnlock")
	flag.BoolVar(&daemon, "D", false, "Daemon mode")
	flag.StringVar(&configFile, "f", "", "specify configuration file")
	flag.Parse()
	su = config.Init(&configFile)
	flag.StringVar(&su.SubURL, "u", su.SubURL, "Load node from subscription url")
	flag.StringVar(&su.Token, "t", su.Token, "The github token")
	flag.StringVar(&su.GistUrl, "g", su.GistUrl, "The gist api URL")
	flag.Parse()
	fmt.Printf("StairUnlock-Go %s %s %s with %s\n", utils.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
	log.SetLevel(su.LogLevel)
	fmt.Printf("Log Level: %s\n", su.LogLevel)
	if su.LocalFile {
		fmt.Println("Local file mode: on")
	} else {
		fmt.Println("Gist mode: on")
	}
	if daemon {
		fmt.Println("Daemon mode: on")
		fmt.Printf("Check internal: %ds\n", su.Internal)
	} else {
		fmt.Println("Daemon mode: off")
	}
}

func run() {
	proxies, cfg, _ := config.GenerateProxies(su)
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	//同时连接数
	connNum := su.MaxConn
	if i := len(proxiesList); i < connNum {
		connNum = i
	}
	start := time.Now()
	netflixList := utils.BatchCheck(proxiesList, connNum)
	log.Warnln("Total %d nodes test completed, %d unlock nodes, Elapsed time: %s", len(proxiesList), len(netflixList), time.Now().Sub(start).String())
	marshal, _ := yaml.Marshal(config.NETFLIXFilter(netflixList, cfg))

	if su.LocalFile {
		_ = ioutil.WriteFile("netflix.yaml", marshal, 0644)
		log.Infoln("Written to netflix.yaml.")
	} else {
		err := utils.Gist(marshal, su)
		if err != nil {
			return
		}
	}
}

func main() {
	//command-line
	if ver {
		fmt.Printf("StairUnlock %s %s %s with %s\n", utils.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
		return
	}
	if help {
		fmt.Printf("StairUnlock %s %s %s with %s\n", utils.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
		flag.PrintDefaults()
		return
	}
	run()

	if daemon {
		start := time.Now()
		for {
			resp, _ := http.Get("https://www.netflix.com/title/70143836")
			err := resp.Body.Close()
			if err != nil {
				return
			}
			if resp.StatusCode != 200 {
				log.Errorln("Cannot access NETFLIX, Retesting all nodes.")
				// 清空 proxiesList 切片
				proxiesList = proxiesList[:0]
				run()
			} else {
				log.Infoln("Stream Media is unlocking.")
				if time.Now().Sub(start) > 3*time.Hour {
					// 每3小时强制更新
					start = time.Now()
					proxiesList = proxiesList[:0]
					run()
				}
			}
			time.Sleep(time.Duration(su.Internal) * time.Second)
		}
	}
}
