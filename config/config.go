package config

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies"`
}

type SuConfig struct {
	ConverterAPI string       `yaml:"converterAPI"`
	SubURL       string       `yaml:"subURL"`
	LocalFile    bool         `yaml:"localFile"`
	Token        string       `yaml:"token"`
	MaxConn      int          `yaml:"maxConn"`
	GistUrl      string       `yaml:"gistUrl,omitempty"`
	Internal     int          `yaml:"internal"`
	LogLevel     log.LogLevel `yaml:"log_level"`
}

func Init(cfgPath *string) (s *SuConfig) {
	//initial config.yaml
	var buf []byte
	if *cfgPath != "" {
		buf, _ = ioutil.ReadFile(*cfgPath)
	} else {
		_, err := os.Stat("config.yaml")
		if err != nil {
			b, _ := ioutil.ReadFile("config.example.yaml")
			_ = ioutil.WriteFile("config.yaml", b, 644)
		}
		buf, _ = ioutil.ReadFile("config.yaml")
	}
	var cfg SuConfig
	_ = yaml.Unmarshal(buf, &cfg)
	return &cfg
}

func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = ioutil.WriteFile("proxies.yaml", nil, 0644)
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}
	return data, err
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{}
	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}
	return rawCfg, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, err error) {
	proxies = make(map[string]C.Proxy)
	proxiesConfig := cfg.Proxy

	for idx, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			return nil, fmt.Errorf("proxy %d: %w", idx, err)
		}
		if _, exist := proxies[proxy.Name()]; exist {
			return nil, fmt.Errorf("proxy %s is the duplicate name", proxy.Name())
		}
		proxies[proxy.Name()] = proxy
	}
	return proxies, err
}

func GenerateProxies(sCfg *SuConfig) (proxies map[string]C.Proxy, cfg *RawConfig, err error) {
	var data []byte
	if sCfg.LocalFile {
		configFile := "proxies.yaml"
		currentDir, _ := os.Getwd()
		configFile = filepath.Join(currentDir, configFile)
		data, err = readConfig(configFile)
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Infoln("Converting from API server.")
		data = convertAPI(sCfg)
	}
	cfg, _ = UnmarshalRawConfig(data)
	proxies, err = parseProxies(cfg)
	return proxies, cfg, err
}

func convertAPI(sCfg *SuConfig) (p []byte) {
	baseUrl, err := url.Parse(sCfg.ConverterAPI)
	baseUrl.Path += "sub"
	params := url.Values{}
	params.Add("target", "clash")
	params.Add("list", "true")
	params.Add("url", sCfg.SubURL)
	baseUrl.RawQuery = params.Encode()
	reqs, err := http.Get(baseUrl.String())
	if err != nil {
		log.Errorln(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(reqs.Body)
	p, _ = ioutil.ReadAll(reqs.Body)
	if strings.Contains(string(p), "The following link doesn't contain any valid node info") {
		log.Errorln("The following link doesn't contain any valid node info.")
		panic("Invalid link.")
	}
	return
}
