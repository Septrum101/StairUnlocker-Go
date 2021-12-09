package main

import (
	"StairUnlocker-Go/config"
	"StairUnlocker-Go/utils"
	"flag"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	tg          utils.TgBot
	start       time.Time
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
	start = time.Now()
	netflixList := utils.BatchCheck(proxiesList, connNum)
	report := fmt.Sprintf("Total %d nodes test completed, %d unlock nodes, Elapsed time: %s", len(proxiesList), len(netflixList), time.Now().Sub(start).String())
	log.Warnln(report)
	if daemon && su.EnableTelegram {
		telegramReport := fmt.Sprintf("%s, Timestamp: %s", report, time.Now().Format("2006/01/02 15:04:05.000"))
		tg.SendMessage = telegramReport
		_, _ = tg.Bot.Send(tgBot.NewMessage(su.Telegram.ChatID, telegramReport))
	}

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

func daemonRun() {
	start = time.Now()
	for {
		resp, _ := http.Get("https://www.netflix.com/title/70143836")
		err := resp.Body.Close()
		if err != nil {
			log.Errorln(err.Error())
		}
		if resp.StatusCode != 200 {
			log.Errorln("Cannot access NETFLIX, Retesting all nodes.")
			// 清空 proxiesList 切片
			start = time.Now()
			proxiesList = proxiesList[:0]
			run()
		} else {
			log.Infoln("Stream Media is unlocking.")
			if time.Now().Sub(start) > 12*time.Hour {
				// 每12小时强制更新
				log.Infoln("Force re-testing all nodes.")
				start = time.Now()
				proxiesList = proxiesList[:0]
				run()
			}
		}
		time.Sleep(time.Duration(su.Internal) * time.Second)
	}
}

func main() {
	versionStr := fmt.Sprintf("StairUnlock %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime)
	//command-line
	if ver {
		fmt.Printf(versionStr)
		return
	}
	if help {
		fmt.Printf(versionStr)
		flag.PrintDefaults()
		return
	}
	fmt.Printf(versionStr)
	log.SetLevel(su.LogLevel)
	fmt.Printf("Log Level: %s\n", su.LogLevel)
	if su.LocalFile {
		fmt.Println("Local file mode: on")
	} else {
		fmt.Println("Gist mode: on")
	}
	// 初始化信息
	if daemon {
		fmt.Println("Daemon mode: on")
		fmt.Printf("Check internal: %ds\n", su.Internal)
		// 初始化telegramBot
		if su.EnableTelegram {
			fmt.Println("Telegram Bot: on")
			tg.NewBot(su)
		}
	} else {
		fmt.Println("Daemon mode: off")
	}
	run()

	if daemon {
		ch := make(chan bool, 1)
		if su.EnableTelegram {
			go func() { tg.TelegramUpdates(&ch) }()
		} else {
			close(ch)
		}
		go func() { daemonRun() }()
		for check := range ch {
			if check {
				log.Infoln("Telegram: Force re-testing all nodes.")
				start = time.Now()
				proxiesList = proxiesList[:0]
				run()
				ch <- false
			}
		}
	}
}
