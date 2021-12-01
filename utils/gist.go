package utils

import (
	"StairUnlocker-Go/config"
	"bytes"
	"encoding/json"
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
)

type gistCtx struct {
	cfg  *config.SuConfig
	ctx  []byte
	body io.ReadCloser
}

func (ths *gistCtx) upload() (*http.Response, error) {
	c, _ := json.Marshal(map[string]map[string]map[string]string{"files": {"stairunlocker": {"content": string(ths.ctx)}}})
	client := &http.Client{}
	if ths.cfg.GistUrl == "" {
		ths.cfg.GistUrl = "https://api.github.com/gists"
	}
	req, err := http.NewRequest(http.MethodPost, ths.cfg.GistUrl, bytes.NewReader(c))
	req.Header = map[string][]string{
		"Authorization": {"token " + ths.cfg.Token},
		"Accept":        {"application/vnd.github.v3+json"}}
	reqs, err := client.Do(req)
	return reqs, err
}
func (ths *gistCtx) create() {
	b, _ := ioutil.ReadAll(ths.body)
	var url map[string]string
	_ = yaml.Unmarshal(b, &url)
	ths.cfg.GistUrl = url["url"]
	tmp, _ := yaml.Marshal(ths.cfg)
	_ = ioutil.WriteFile("config.yaml", tmp, 0644)
}

func Gist(data []byte, cfg *config.SuConfig) (err error) {
	ctx := gistCtx{cfg, data, nil}
	reqs, err := ctx.upload()
	ctx.body = reqs.Body
	if reqs.StatusCode == 200 {
		log.Infoln("Update gist success! Please visit: https://gist.github.com for details.")
	} else if reqs.StatusCode == 201 {
		ctx.create()
		log.Infoln("Create gist success! Please visit: https://gist.github.com for details.")
	} else if reqs.StatusCode == 404 {
		ctx.cfg.GistUrl = ""
		reqs, err := ctx.upload()
		if err != nil {
			return err
		}
		ctx.body = reqs.Body
		ctx.create()
		log.Errorln("The gist is not exist! A new gist will be created.")
	}
	return
}
